package goTezos

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	generichash "github.com/GoKillers/libsodium-go/cryptogenerichash"
	"github.com/Messer4/base58check"
	encoding "github.com/anaskhan96/base58check"
	"github.com/jamesruan/sodium"
	"github.com/tyler-smith/go-bip39"
	"log"
	"math"
	"strconv"
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
	
	var operationSignatures []string
	
	// Get current branch head
	blockHead, err := this.GetChainHead()
	if err != nil {
		return operationSignatures, err
	}
	
	// Get the counter for the payment address and increment it
	counter, err := this.getAddressCounter(wallet.Address)
	if err != nil {
		return operationSignatures, err
	}
	counter++
	
	// Split our slice of []Payment into batches
	batches := this.splitPaymentIntoBatches(payments)
	operationSignatures = make([]string, len(batches))
	
	for k := range batches {
		
		// Convert (ie: forge) each 'Payment' into an actual Tezos transfer operation
		operationBytes, operationContents, newCounter, err := this.forgeOperationBytes(blockHead.Hash, counter, wallet, batches[k])
		if err != nil {
			return operationSignatures, err
		}
		counter = newCounter
		
		// Sign this batch of operations with the secret key; return that signature
		edsig := this.signOperationBytes(operationBytes, wallet)
		
		// Extract and decode the bytes of the signature
		decodedSignature := this.decodeSignature(edsig)
		decodedSignature = decodedSignature[10:(len(decodedSignature))]
		
		// The signed bytes of this batch
		fullOperation := operationBytes + decodedSignature
		
		// We can validate this batch against the node for any errors
		if err := this.preApplyOperations(operationContents, edsig, blockHead); err != nil {
			return operationSignatures, fmt.Errorf("CreateBatchPayment failed to Pre-Apply: %s", err)
		}
		
		// Add the signature (raw operation bytes & signature of operations) of this batch of transfers to the returnning slice
		// This will be used to POST to /injection/operation
		operationSignatures[k] = fullOperation
		
	}
	
	return operationSignatures, nil
}


// Pre-apply an operation, or batch of operations, to a Tezos node to ensure correctness
func (this *GoTezos) preApplyOperations(paymentOperations Conts, signature string, blockHead Block) error {
	
	// Create a full transfer request
	var transfer Transfer
	transfer.Signature = signature
	transfer.Contents = paymentOperations.Contents
	transfer.Branch = blockHead.Hash
	transfer.Protocol = blockHead.Protocol
	
	// RPC says outer element must be JSON array
	var transfers = []Transfer{transfer}
	
	// Convert object to JSON string
	transfersOp, err := json.Marshal(transfers)
	if err != nil {
		return err
	}
	
	if this.debug {
		fmt.Println("\n== preApplyOperations Submit:", string(transfersOp))
	}
	
	// POST the JSON to the RPC
	preApplyResp, err := this.PostResponse("/chains/main/blocks/head/helpers/preapply/operations", string(transfersOp))
	if err != nil {
		return err
	}
	
	if this.debug {
		fmt.Println("\n== preApplyOperations Result:", string(preApplyResp.Bytes))
	}
	
	return nil
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
		Address:  this.b58cencode(genericHash, tz1),
		Mnemonic: mnemonic,
		Seed:     seed,
		Kp:       signKP,
		Sk:       this.b58cencode(signKP.SecretKey.Bytes, edsk),
		Pk:       this.b58cencode(signKP.PublicKey.Bytes, edpk),
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
	resp, err := this.GetResponse(rpc, "{}")
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

func (this *GoTezos) forgeOperationBytes(branch_hash string, counter int, wallet Wallet, batch []Payment) (string, Conts, int, error) {

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
				Kind:         "transaction",
				Source:       wallet.Address,
				Fee:          "1420",
				GasLimit:     "11000",
				StorageLimit: "0",
				Amount:       strconv.FormatFloat(roundPlus(batch[k].Amount, 0), 'f', -1, 64),
				Destination:  batch[k].Address,
				Counter:      strconv.Itoa(counter),
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
		return "", contents, counter, fmt.Errorf("POST-Forge Operation Error: %s", err)
	}

	err = json.Unmarshal(output.Bytes, &opBytes)
	if err != nil {
		return "", contents, counter, fmt.Errorf("Forge Operation Error: %s", err)
	}
	
	return opBytes, contents, counter, nil
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

//Helper function to return the decoded signature
func (this *GoTezos) decodeSignature(sig string) string {
	dec_bytes, err := encoding.Decode(sig)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	dec_sig := string(dec_bytes)
	return dec_sig
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
