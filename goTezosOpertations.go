package goTezos

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strconv"
	generichash "github.com/GoKillers/libsodium-go/cryptogenerichash"
	"github.com/Messer4/base58check"
	encoding "github.com/anaskhan96/base58check"
	"github.com/jamesruan/sodium"
	"github.com/tyler-smith/go-bip39"
)

var (
	// How many Transactions per batch are injected. I recommend 100. Now 30 for easier testing
	batchSize = 30
	
	// For (de)constructing addresses
	tz1   = []byte{6, 161, 159}
	edsk  = []byte{43, 246, 78, 7}
	edsk2 = []byte{13, 15, 58, 7}
	edpk  = []byte{13, 15, 37, 217}
)

//Forges batch payments and returns them ready to inject to an tezos rpc
func (this *GoTezos) CreateBatchPayment(payments []Payment, wallet Wallet) ([]string, error) {
	
	var dec_sigs []string
	
	//Get current branch head
	blockHead, err := this.GetChainHead()
	if err != nil {
		return dec_sigs, err
	}
	
	//get the counter for the wallet && increment it
	counter, err := this.getAddressCounter(wallet.Address)
	if err != nil {
		return dec_sigs, err
	}
	counter++
	
	batches := this.splitPaymentIntoBatches(payments)
	dec_sigs = make([]string, len(batches))
	
	for k := range batches {
		
		operation_bytes, _, newCounter := this.forgeOperationBytes(blockHead.Hash, counter, wallet, batches[k])
		counter = newCounter

		signed_operation_bytes := this.signOperationBytes(operation_bytes, wallet)
		
		//TODO: Here we could preapply, but eg. tezrpc is not supporting it
		dec_sig := this.decodeSignature(signed_operation_bytes, operation_bytes)
		dec_sigs[k] = dec_sig
	}
	
	return dec_sigs, nil
}


func (this *GoTezos) CreateWallet(mnemonic, password string) (Wallet, error) {
	
	var signSecretKey sodium.SignSecretKey
	var wallet Wallet
	
	seed := bip39.NewSeed(mnemonic, password)
	signSecretKey.Bytes = []byte(seed)
	signSeed := signSecretKey.Seed()
	signKP := sodium.SeedSignKP(signSeed)
	key := sodium.GenericHashKey{signKP.PublicKey.Bytes}
	genericHash, _ := generichash.CryptoGenericHash(20, key.Bytes, nil)
	
	wallet = Wallet{
		Address: this.b58cencode(genericHash, tz1),
		Mnemonic: mnemonic,
		Seed: seed,
		Kp: signKP,
		Sk: this.b58cencode(signKP.SecretKey.Bytes, edsk),
		Pk: this.b58cencode(signKP.PublicKey.Bytes, edpk),
	}
	
	return wallet, nil
}

func (this *GoTezos) ImportWallet(address, public, secret string) (Wallet, error) {
	
	var wallet Wallet
	var signKP sodium.SignKP
	
	// Sanity check
	secretLength := len(secret)
	if secret[:4] != "edsk" || (secretLength != 98 && secretLength != 54) {
		return wallet, fmt.Errorf("Import Wallet Error: The provided secret does not conform to known patterns.")
	}
	
	// Determine if 'secret' is an actual secret key or a seed
	if secretLength == 98 {
		
		// A full secret key
		decodedSecretKey := this.b58cdecode(secret, edsk)
		
		// Public key is last 32 of decoded secret, re-encoded as edpk
		publicKey := decodedSecretKey[32:]
		
		signKP.PublicKey = sodium.SignPublicKey{[]byte(publicKey)}
		signKP.SecretKey = sodium.SignSecretKey{[]byte(secret)}
		
		wallet.Sk = secret
		
	} else if secretLength == 54 {
		
		// "secret" is actually a seed
		decodedSeed := this.b58cdecode(secret, edsk2)
		
		signSeed := sodium.SignSeed{decodedSeed}
		
		// Reconstruct keypair from seed
		signKP = sodium.SeedSignKP(signSeed)
		
		wallet.Sk = this.b58cencode(signKP.SecretKey.Bytes, edsk)
		
	} else {
		
		return wallet, fmt.Errorf("Import Wallet Error: Secret key is not the correct length.")
	}
	
	// Generate public address from public key
	hashKey := sodium.GenericHashKey{signKP.PublicKey.Bytes}
	addressHash, _ := generichash.CryptoGenericHash(20, hashKey.Bytes, nil)
	generatedAddress := this.b58cencode(addressHash, tz1)
	
	// Populate Wallet
	wallet.Address = generatedAddress
	wallet.Kp = signKP
	wallet.Pk = this.b58cencode(signKP.PublicKey.Bytes, edpk)

	// Couple more sanity checks
	if generatedAddress != address {
		return wallet, fmt.Errorf("Import Wallet Error: Reconstructed address '%s' and provided address '%s' do not match.", generatedAddress, address)
	}
	
	if wallet.Pk != public {
		return wallet, fmt.Errorf("Import Wallet Error: Reconstructed Pkh '%s' and provided Pkh '%s' do not match.", wallet.Pk, public)
	}
	
	return wallet, nil
}

//Getting the Counter of an address from the RPC
func (this *GoTezos) getAddressCounter(address string) (int, error) {
	rpc := "/chains/main/blocks/head/context/contracts/" + address + "/counter"
	resp, err := this.GetResponse(rpc,"{}")
	if err != nil {
		return 0, err
	}
	rtnStr, err := unMarshelString(resp.Bytes)
	if err != nil {
		return 0, err
	}
	counter, err := strconv.Atoi(rtnStr)
	return counter, err
}

func (this *GoTezos) splitPaymentIntoBatches(rewards []Payment) [][]Payment {
	var batches [][]Payment
	for i := 0; i < len(rewards); i += batchSize {
		end := i + batchSize
		if end > len(rewards) {
			end = len(rewards)
		}
		batches = append(batches, rewards[i:end])
	}
	return batches
}

func (this *GoTezos) forgeOperationBytes(branch_hash string, counter int, wallet Wallet, batch []Payment) (string, Conts, int) {

	var contents Conts
	var combinedOps []TransOp

	//left here to display how to reveal a new wallet (needs funds to be revealed!)
	/**
	  combinedOps = append(combinedOps, TransOp{Kind: "reveal", PublicKey: wallet.pk , Source: wallet.address, Fee: "0", GasLimit: "127", StorageLimit: "0", Counter: strCounter})
	  counter++
	**/

	for k := range batch {
		
		if batch[k].Amount > 0 {
		
			operation := TransOp{
				Kind: "transaction",
				Source: wallet.Address,
				Fee: "1420",
				GasLimit: "11000",
				StorageLimit: "0",
				Amount: strconv.FormatFloat(roundPlus(batch[k].Amount, 0), 'f', -1, 64),
				Destination: batch[k].Address,
				Counter: strconv.Itoa(counter),
			}
			combinedOps = append(combinedOps, operation)
			counter++
		}
	}
	contents.Contents = combinedOps
	contents.Branch = branch_hash
	
	var opBytes string
	
	forge := "/chains/main/blocks/head/helpers/forge/operations"
	output, err := this.PostResponse(forge, contents.String())
	if err != nil {
		return "", contents, counter
	}
	
	err = json.Unmarshal(output.Bytes, &opBytes)
	if err != nil {
		log.Println("Could not unmarshal to string " + err.Error())
		return "", contents, counter
	}
	
	return opBytes, contents, counter
}

//Sign previously forged Operation bytes using secret key of wallet
func (this *GoTezos) signOperationBytes(operation_bytes string, wallet Wallet) string {
	//Prefixes
	edsigByte := []byte{9, 245, 205, 134, 18}
	watermark := []byte{3}

	op, err := hex.DecodeString(operation_bytes)
	if err != nil {
		log.Fatal(err)
	}
	op = append(watermark, op...)
	genericHash, _ := generichash.CryptoGenericHash(32, op, nil)
	sig := sodium.Bytes(genericHash).SignDetached(wallet.Kp.SecretKey)
	edsig := this.b58cencode(sig.Bytes, edsigByte)
	return edsig
}

func (this *GoTezos) decodeSignature(sig string, operation_bytes string) (dec_sig string) {
	dec_bytes, err := encoding.Decode(sig)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	dec_sig = string(dec_bytes)
	//no need to cut the 8 last hexbyte, that's already done
	dec_sig = dec_sig[10:(len(dec_sig))]
	dec_sig = operation_bytes + dec_sig
	return
}

//Helper Function to get the right format for wallet.
func (this *GoTezos) b58cencode(payload []byte, prefix []byte) string {
	n := make([]byte, (len(prefix) + len(payload)))
	for k := range prefix {
		n[k] = prefix[k]
	}
	for l := range payload {
		n[l+len(prefix)] = payload[l]
	}
	b58c := base58check.Encode(n)
	return b58c
}

func (this *GoTezos) b58cdecode(payload string, prefix []byte) []byte {
	b58c, _ := base58check.Decode(payload)
	return b58c[len(prefix):]
}

//Helper Functions to round float64
func roundPlus(f float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return round(f*shift) / shift
}

func round(f float64) float64 {
	return math.Floor(f + .5)
}
