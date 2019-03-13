// +build cgo

package goTezos

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Messer4/base58check"
	"github.com/jamesruan/sodium"
	"golang.org/x/crypto/pbkdf2"
)

//Forges batch payments and returns them ready to inject to a Tezos RPC. PaymentFee must be expressed in mutez.
func (this *GoTezos) CreateBatchPayment(payments []Payment, wallet Wallet, paymentFee int) ([]string, error) {

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
		operationBytes, operationContents, newCounter, err := this.forgeOperationBytes(blockHead.Hash, counter, wallet, batches[k], paymentFee)
		if err != nil {
			return operationSignatures, err
		}
		counter = newCounter

		// Sign this batch of operations with the secret key; return that signature
		edsig, err := this.signOperationBytes(operationBytes, wallet)
		if err != nil {
			return operationSignatures, err
		}

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

func (this *GoTezos) CreateWallet(mnemonic, password string) (Wallet, error) {

	var signSecretKey sodium.SignSecretKey
	var wallet Wallet

	// Copied from https://github.com/tyler-smith/go-bip39/blob/dbb3b84ba2ef14e894f5e33d6c6e43641e665738/bip39.go#L268
	seed := pbkdf2.Key([]byte(mnemonic), []byte("mnemonic"+password), 2048, 64, sha512.New)
	signSecretKey.Bytes = []byte(seed)
	signSeed := signSecretKey.Seed()
	signKP := sodium.SeedSignKP(signSeed)

	// Generate public address from public key
	generatedAddress, err := this.generatePublicHash(signKP)
	if err != nil {
		return wallet, fmt.Errorf("CreateWallet Error: %s", err)
	}

	// Construct wallet
	wallet = Wallet{
		Address:  generatedAddress,
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

	wallet.Kp = signKP

	// Generate public address from public key
	generatedAddress, err := this.generatePublicHash(signKP)
	if err != nil {
		return wallet, fmt.Errorf("Import Wallet Error: %s", err)
	}

	if generatedAddress != address {
		return wallet, fmt.Errorf("Import Wallet Error: Reconstructed address '%s' and provided address '%s' do not match.", generatedAddress, address)
	}
	wallet.Address = generatedAddress

	// Genrate and check public key
	generatedPublicKey := this.b58cencode(signKP.PublicKey.Bytes, edpk)
	if generatedPublicKey != public {
		return wallet, fmt.Errorf("Import Wallet Error: Reconstructed Pkh '%s' and provided Pkh '%s' do not match.", generatedPublicKey, public)
	}
	wallet.Pk = generatedPublicKey

	return wallet, nil
}

// Import an encrypted wallet using password provided by caller.
// Caller should remove any 'encrypted:' scheme prefix.
func (this *GoTezos) ImportEncryptedWallet(pw, encKey string) (Wallet, error) {

	var wallet Wallet

	// Check if user copied 'encrypted:' scheme prefix
	if encKey[:5] != "edesk" || len(encKey) != 88 {
		return wallet, fmt.Errorf("ImportEncryptedWallet: Encrypted secret key does not conform to known patterns.")
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

	// No nonce used
	emptyNonceBytes := make([]byte, 24)
	boxNonce := sodium.SecretBoxNonce{emptyNonceBytes}

	// Create box and key object
	var box sodium.Bytes = esm
	boxKey := sodium.SecretBoxKey{key}

	// Decrypt. Returns bytes for a SignSecretKey
	unencSecret, err := box.SecretBoxOpen(boxNonce, boxKey)
	if err != nil {
		return wallet, fmt.Errorf("Incorrect password for encrypted key")
	}
	signSeed := sodium.SignSeed{unencSecret}

	// Create key-pair from signing seed
	signKP := sodium.SeedSignKP(signSeed)

	// public key & secret key
	wallet.Kp = signKP
	wallet.Sk = this.b58cencode(signKP.SecretKey.Bytes, edsk)
	wallet.Pk = this.b58cencode(signKP.PublicKey.Bytes, edpk)

	// Generate public address from public key
	generatedAddress, err := this.generatePublicHash(signKP)
	if err != nil {
		return wallet, fmt.Errorf("ImportEncryptedWallet: %s", err)
	}
	wallet.Address = generatedAddress

	return wallet, nil
}

//Sign previously forged Operation bytes using secret key of wallet
func (this *GoTezos) signOperationBytes(operation_bytes string, wallet Wallet) (string, error) {

	//Prefixes
	edsigByte := []byte{9, 245, 205, 134, 18}
	watermark := []byte{3}

	opBytes, err := hex.DecodeString(operation_bytes)
	if err != nil {
		return "", fmt.Errorf("Unable to sign operation bytes: %s", err)
	}
	opBytes = append(watermark, opBytes...)

	// Generic hash of 32 bytes
	genericHash := sodium.NewGenericHash(32)

	// Write operation bytes to hash
	i, err := genericHash.Write(opBytes)
	if i != len(opBytes) || err != nil {
		return "", fmt.Errorf("Unable to write operations to generic hash")
	}
	finalHash := genericHash.Sum([]byte{})

	// Sign the finalized generic hash of operations and b58 encode
	sig := sodium.Bytes(finalHash).SignDetached(wallet.Kp.SecretKey)
	edsig := this.b58cencode(sig.Bytes, edsigByte)

	return edsig, nil
}

//Helper function to generate public key hash
func (this *GoTezos) generatePublicHash(kp sodium.SignKP) (string, error) {

	// Generic hash of 20 bytes
	genericHash := sodium.NewGenericHash(20)

	// Write public key
	i, err := genericHash.Write(kp.PublicKey.Bytes)
	if i != 32 || err != nil {
		return "", fmt.Errorf("Unable to write public key to generic hash")
	}

	// "Sum" up the hash calculation and return encoded hash
	return this.b58cencode(genericHash.Sum([]byte{}), tz1), nil
}

func (this *GoTezos) forgeOperationBytes(branch_hash string, counter int, wallet Wallet, batch []Payment, paymentFee int) (string, Conts, int, error) {

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
