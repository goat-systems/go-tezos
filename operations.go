package gotezos

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"strconv"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/nacl/secretbox"

	"github.com/Messer4/base58check"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/pbkdf2"
)

var (
	// How many Transactions per batch are injected. I recommend 100. Now 30 for easier testing
	batchSize = 100

	// For (de)constructing addresses
	tz1   = []byte{6, 161, 159}
	edsk  = []byte{43, 246, 78, 7}
	edsk2 = []byte{13, 15, 58, 7}
	edpk  = []byte{13, 15, 37, 217}
	edesk = []byte{7, 90, 60, 179, 41}
)

// CreateBatchPayment forges batch payments and returns them ready to inject to a Tezos RPC. PaymentFee must be expressed in mutez.
func (gt *GoTezos) CreateBatchPayment(payments []Payment, wallet Wallet, paymentFee int, gaslimit int) ([]string, error) {

	var operationSignatures []string

	// Get current branch head
	blockHead, err := gt.GetChainHead()
	if err != nil {
		return operationSignatures, err
	}

	// Get the counter for the payment address and increment it
	counter, err := gt.getAddressCounter(wallet.Address)
	if err != nil {
		return operationSignatures, err
	}
	counter++

	// Split our slice of []Payment into batches
	batches := gt.splitPaymentIntoBatches(payments)
	operationSignatures = make([]string, len(batches))

	for k := range batches {

		// Convert (ie: forge) each 'Payment' into an actual Tezos transfer operation
		operationBytes, operationContents, newCounter, err := gt.forgeOperationBytes(blockHead.Hash, counter, wallet, batches[k], paymentFee, gaslimit)
		if err != nil {
			return operationSignatures, err
		}
		counter = newCounter

		// Sign gt batch of operations with the secret key; return that signature
		edsig, err := gt.signOperationBytes(operationBytes, wallet)
		if err != nil {
			return operationSignatures, err
		}

		// Extract and decode the bytes of the signature
		decodedSignature := gt.decodeSignature(edsig)
		decodedSignature = decodedSignature[10:(len(decodedSignature))]

		// The signed bytes of gt batch
		fullOperation := operationBytes + decodedSignature

		// We can validate gt batch against the node for any errors
		if err := gt.preApplyOperations(operationContents, edsig, blockHead); err != nil {
			return operationSignatures, fmt.Errorf("CreateBatchPayment failed to Pre-Apply: %s", err)
		}

		// Add the signature (raw operation bytes & signature of operations) of gt batch of transfers to the returnning slice
		// gt will be used to POST to /injection/operation
		operationSignatures[k] = fullOperation

	}

	return operationSignatures, nil
}

// CreateWallet returns Wallet with the mnemonic and password provided
func (gt *GoTezos) CreateWallet(mnenomic string, password string) (Wallet, error) {
	// Copied from https://github.com/tyler-smith/go-bip39/blob/dbb3b84ba2ef14e894f5e33d6c6e43641e665738/bip39.go#L268
	seed := pbkdf2.Key([]byte(mnenomic), []byte("mnemonic"+password), 2048, 32, sha512.New)
	privKey := ed25519.NewKeyFromSeed(seed)
	pubKey := privKey.Public().(ed25519.PublicKey)
	pubKeyBytes := []byte(pubKey)
	signKp := KeyPair{PrivKey: privKey, PubKey: pubKeyBytes}

	address, err := gt.generatePublicHash(pubKeyBytes)
	if err != nil {
		return Wallet{}, err
	}

	wallet := Wallet{
		Address:  address,
		Mnemonic: mnenomic,
		Kp:       signKp,
		Seed:     seed,
		Sk:       gt.b58cencode(privKey, edsk),
		Pk:       gt.b58cencode(pubKeyBytes, edpk),
	}

	return wallet, nil
}

// ImportWallet returns an imported Wallet
func (gt *GoTezos) ImportWallet(address, public, secret string) (Wallet, error) {

	var wallet Wallet
	var signKP KeyPair

	// Sanity check
	secretLength := len(secret)
	if secret[:4] != "edsk" || (secretLength != 98 && secretLength != 54) {
		return wallet, fmt.Errorf("import Wallet Error: The provided secret does not conform to known patterns")
	}

	// Determine if 'secret' is an actual secret key or a seed
	if secretLength == 98 {

		// A full secret key
		decodedSecretKey := gt.b58cdecode(secret, edsk)

		// Public key is last 32 of decoded secret, re-encoded as edpk
		publicKey := decodedSecretKey[32:]

		signKP.PubKey = []byte(publicKey)
		signKP.PrivKey = []byte(secret)

		wallet.Sk = secret

	} else if secretLength == 54 {

		// "secret" is actually a seed
		decodedSeed := gt.b58cdecode(secret, edsk2)

		//signSeed := sodium.SignSeed{Bytes: decodedSeed}

		// Reconstruct keypair from seed
		privKey := ed25519.NewKeyFromSeed(decodedSeed)
		pubKey := privKey.Public().(ed25519.PublicKey)
		signKP.PrivKey = privKey
		signKP.PubKey = []byte(pubKey)

		wallet.Sk = gt.b58cencode(signKP.PrivKey, edsk)

	} else {

		return wallet, fmt.Errorf("import Wallet Error: Secret key is not the correct length")
	}

	wallet.Kp = signKP

	// Generate public address from public key
	generatedAddress, err := gt.generatePublicHash(signKP.PubKey)
	if err != nil {
		return wallet, fmt.Errorf("Import Wallet Error: %s", err)
	}

	if generatedAddress != address {
		return wallet, fmt.Errorf("import Wallet Error: Reconstructed address '%s' and provided address '%s' do not match", generatedAddress, address)
	}
	wallet.Address = generatedAddress

	// Genrate and check public key
	generatedPublicKey := gt.b58cencode(signKP.PubKey, edpk)
	if generatedPublicKey != public {
		return wallet, fmt.Errorf("import Wallet Error: Reconstructed Pkh '%s' and provided Pkh '%s' do not match", generatedPublicKey, public)
	}
	wallet.Pk = generatedPublicKey

	return wallet, nil
}

// ImportEncryptedWallet imports an encrypted wallet using password provided by caller.
// Caller should remove any 'encrypted:' scheme prefix.
func (gt *GoTezos) ImportEncryptedWallet(pw, encKey string) (Wallet, error) {

	var wallet Wallet

	// Check if user copied 'encrypted:' scheme prefix
	if encKey[:5] != "edesk" || len(encKey) != 88 {
		return wallet, fmt.Errorf("importEncryptedWallet: encrypted secret key does not conform to known patterns")
	}

	// Convert key from base58 to []byte
	b58c, err := base58check.Decode(encKey)
	if err != nil {
		return wallet, err
	}

	// Strip off prefix and extract parts
	esb := b58c[len(edesk):]
	salt := esb[:8]
	esm := esb[8:] // encrypted key

	// Convert string pw to []byte
	passWd := []byte(pw)

	// Derive a key from password, salt and number of iterations
	key := pbkdf2.Key(passWd, salt, 32768, 32, sha512.New)
	var byteKey [32]byte
	for i := range key {
		byteKey[i] = key[i]
	}

	var out []byte
	var emptyNonceBytes [24]byte

	unencSecret, ok := secretbox.Open(out, esm, &emptyNonceBytes, &byteKey)
	if !ok {
		return wallet, fmt.Errorf("incorrect password for encrypted key")
	}

	privKey := ed25519.NewKeyFromSeed(unencSecret)
	pubKey := privKey.Public().(ed25519.PublicKey)
	pubKeyBytes := []byte(pubKey)
	signKP := KeyPair{PrivKey: privKey, PubKey: pubKeyBytes}

	// public key & secret key
	wallet.Kp = signKP
	wallet.Sk = gt.b58cencode(signKP.PrivKey, edsk)
	wallet.Pk = gt.b58cencode(signKP.PubKey, edpk)

	// Generate public address from public key
	generatedAddress, err := gt.generatePublicHash(signKP.PubKey)
	if err != nil {
		return wallet, fmt.Errorf("importEncryptedWallet: %s", err)
	}
	wallet.Address = generatedAddress

	return wallet, nil
}

//Sign previously forged Operation bytes using secret key of wallet
func (gt *GoTezos) signOperationBytes(operationBytes string, wallet Wallet) (string, error) {

	//Prefixes
	edsigByte := []byte{9, 245, 205, 134, 18}
	watermark := []byte{3}

	opBytes, err := hex.DecodeString(operationBytes)
	if err != nil {
		return "", fmt.Errorf("Unable to sign operation bytes: %s", err)
	}
	opBytes = append(watermark, opBytes...)

	// Generic hash of 32 bytes
	genericHash, err := blake2b.New(32, []byte{})

	// Write operation bytes to hash
	i, err := genericHash.Write(opBytes)
	if i != len(opBytes) || err != nil {
		return "", fmt.Errorf("Unable to write operations to generic hash")
	}
	finalHash := genericHash.Sum([]byte{})

	// Sign the finalized generic hash of operations and b58 encode
	sig := ed25519.Sign(wallet.Kp.PrivKey, finalHash)
	//sig := sodium.Bytes(finalHash).SignDetached(wallet.Kp.PrivKey)
	edsig := gt.b58cencode(sig, edsigByte)

	return edsig, nil
}

func (gt *GoTezos) generatePublicHash(publicKey []byte) (string, error) {
	hash, err := blake2b.New(20, []byte{})
	hash.Write(publicKey)
	if err != nil {
		return "", fmt.Errorf("Unable to write public key to generic hash: %v", err)
	}
	return gt.b58cencode(hash.Sum(nil), tz1), nil
}

func (gt *GoTezos) forgeOperationBytes(branchHash string, counter int, wallet Wallet, batch []Payment, paymentFee int, gaslimit int) (string, Conts, int, error) {

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
				Fee:          strconv.Itoa(paymentFee),
				GasLimit:     strconv.Itoa(gaslimit),
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
	contents.Branch = branchHash

	var opBytes string

	forge := "/chains/main/blocks/head/helpers/forge/operations"
	output, err := gt.PostResponse(forge, contents.String())
	if err != nil {
		return "", contents, counter, fmt.Errorf("POST-Forge Operation Error: %s", err)
	}

	err = json.Unmarshal(output.Bytes, &opBytes)
	if err != nil {
		return "", contents, counter, fmt.Errorf("Forge Operation Error: %s", err)
	}

	return opBytes, contents, counter, nil
}

// Pre-apply an operation, or batch of operations, to a Tezos node to ensure correctness
func (gt *GoTezos) preApplyOperations(paymentOperations Conts, signature string, blockHead Block) error {

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

	if gt.debug {
		fmt.Println("\n== preApplyOperations Submit:", string(transfersOp))
	}

	// POST the JSON to the RPC
	preApplyResp, err := gt.PostResponse("/chains/main/blocks/head/helpers/preapply/operations", string(transfersOp))
	if err != nil {
		return err
	}

	if gt.debug {
		fmt.Println("\n== preApplyOperations Result:", string(preApplyResp.Bytes))
	}

	return nil
}

// InjectOperation injects an signed operation string and returns the response
func (gt *GoTezos) InjectOperation(op string) ([]byte, error) {
	post := "/injection/operation"
	jsonBytes, err := json.Marshal(op)
	if err != nil {
		return nil, err
	}
	resp, err := gt.PostResponse(post, string(jsonBytes))
	if err != nil {
		return resp.Bytes, err
	}
	return resp.Bytes, nil
}

//Getting the Counter of an address from the RPC
func (gt *GoTezos) getAddressCounter(address string) (int, error) {
	rpc := "/chains/main/blocks/head/context/contracts/" + address + "/counter"
	resp, err := gt.GetResponse(rpc, "{}")
	if err != nil {
		return 0, err
	}
	rtnStr, err := unmarshalString(resp.Bytes)
	if err != nil {
		return 0, err
	}
	counter, err := strconv.Atoi(rtnStr)
	return counter, err
}

func (gt *GoTezos) splitPaymentIntoBatches(rewards []Payment) [][]Payment {
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

//Helper function to return the decoded signature
func (gt *GoTezos) decodeSignature(sig string) string {
	decBytes, err := base58check.Decode(sig)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	return hex.EncodeToString(decBytes)
}

//Helper Function to get the right format for wallet.
func (gt *GoTezos) b58cencode(payload []byte, prefix []byte) string {
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

func (gt *GoTezos) b58cdecode(payload string, prefix []byte) []byte {
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
