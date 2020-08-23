package gotezos

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/Messer4/base58check"
	"github.com/pkg/errors"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/nacl/secretbox"
	"golang.org/x/crypto/pbkdf2"
)

/*
Wallet is a Tezos wallet.
*/
type Wallet struct {
	Address  string
	Mnemonic string
	Seed     []byte
	Kp       keyPair
	Sk       string
	Pk       string
}

type keyPair struct {
	PrivKey []byte
	PubKey  []byte
}

// SignOperationOutput contains an operation with the signature appended, and the signature
type SignOperationOutput struct {
	SignedOperation string
	Signature       string
	EDSig           string
}

/*
Balance gives access to the balance of a contract.

Path:
	../<block_id>/context/contracts/<contract_id>/balance (GET)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-context-contracts-contract-id-balance

Parameters:

	blockhash:
		The hash of block (height) of which you want to make the query.

	address:
		Any tezos public address.
*/
func (t *GoTezos) Balance(blockhash, address string) (*big.Int, error) {
	query := fmt.Sprintf("/chains/main/blocks/%s/context/contracts/%s/balance", blockhash, address)
	resp, err := t.get(query)
	if err != nil {
		return big.NewInt(0), errors.Wrap(err, "failed to get balance")
	}

	balance, err := newInt(resp)
	if err != nil {
		return big.NewInt(0), errors.Wrap(err, "failed to unmarshal balance")
	}

	return balance.Big, nil
}

/*
CreateWallet creates a new wallet.


Parameters:

	mnenomic:
		The seed phrase for the new wallet.

	password:
		The password for the wallet.
*/
func CreateWallet(mnenomic string, password string) (*Wallet, error) {

	seed := pbkdf2.Key([]byte(mnenomic), []byte("mnemonic"+password), 2048, 32, sha512.New)
	privKey := ed25519.NewKeyFromSeed(seed)
	pubKey := privKey.Public().(ed25519.PublicKey)
	pubKeyBytes := []byte(pubKey)
	signKp := keyPair{PrivKey: privKey, PubKey: pubKeyBytes}

	address, err := generatePublicHash(pubKeyBytes)
	if err != nil {
		return &Wallet{}, errors.Wrapf(err, "could not create wallet")
	}

	wallet := Wallet{
		Address:  address,
		Mnemonic: mnenomic,
		Kp:       signKp,
		Seed:     seed,
		Sk:       B58cencode(privKey, edskprefix),
		Pk:       B58cencode(pubKeyBytes, edpkprefix),
	}

	return &wallet, nil
}

/*
ImportWallet imports an unencrypted wallet.

Parameters:

	hash:
		The public key hash of the wallet (tz1, KT1).

	pk:
		The public key of the wallet (edpk).

	sk:
		The secret key of the wallet (edsk).
*/
func ImportWallet(hash, pk, sk string) (*Wallet, error) {

	var wallet Wallet
	var signKP keyPair

	// Sanity check
	secretLength := len(sk)
	if secretLength != 98 && secretLength != 54 {
		return &wallet, errors.New("wallet prefix is not edsk")
	}

	if sk[:4] != "edsk" {
		return &wallet, errors.New("wallet prefix is not edsk")
	}

	// Determine if 'secret' is an actual secret key or a seed
	if secretLength == 98 {

		// A full secret key
		decodedSecretKey := b58cdecode(sk, edskprefix)

		// Public key is last 32 of decoded secret, re-encoded as edpk
		publicKey := decodedSecretKey[32:]

		signKP.PubKey = publicKey
		signKP.PrivKey = decodedSecretKey

		wallet.Sk = sk

	} else if secretLength == 54 {

		// "secret" is actually a seed
		decodedSeed := b58cdecode(sk, edskprefix2)

		//signSeed := sodium.SignSeed{Bytes: decodedSeed}

		// Reconstruct keypair from seed
		privKey := ed25519.NewKeyFromSeed(decodedSeed)
		pubKey := privKey.Public().(ed25519.PublicKey)
		signKP.PrivKey = privKey
		signKP.PubKey = []byte(pubKey)

		wallet.Sk = B58cencode(signKP.PrivKey, edskprefix)

	} else {
		return &wallet, errors.Errorf("wallet secret key length '%d' does not = '%d'", 54, secretLength)
	}

	wallet.Kp = signKP

	// Generate public address from public key
	generatedAddress, err := generatePublicHash(signKP.PubKey)
	if err != nil {
		return &wallet, errors.Wrapf(err, "could not generate public hash")
	}

	if generatedAddress != hash {
		return &wallet, errors.Errorf("reconstructed address '%s' does not match provided address '%s'", generatedAddress, hash)
	}

	wallet.Address = generatedAddress

	// Genrate and check public key
	generatedPublicKey := B58cencode(signKP.PubKey, edpkprefix)
	if generatedPublicKey != pk {
		return &wallet, errors.Errorf("reconstructed pk '%s' does not match provided pk '%s'", generatedPublicKey, pk)
	}
	wallet.Pk = generatedPublicKey

	return &wallet, nil
}

/*
ImportEncryptedWallet imports an encrypted wallet.

Parameters:

	password:
		The password for the wallet.

	esk:
		The encrypted secret key of the wallet (encrypted:edesk).
*/
func ImportEncryptedWallet(password, esk string) (*Wallet, error) {

	var wallet Wallet
	// Check if user copied 'encrypted:' scheme prefix
	if len(esk) != 88 {
		return &wallet, errors.New("encrypted secret key does not 88 characters long")
	}
	if esk[:5] != "edesk" {
		return &wallet, errors.New("encrypted secret key does not prefix with edesk")
	}

	// Convert key from base58 to []byte
	b58c, err := decode(esk)
	if err != nil {
		return &wallet, errors.Wrap(err, "encrypted key is not base58")
	}

	// Strip off prefix and extract parts
	esb := b58c[len(edeskprefix):]
	salt := esb[:8]
	esm := esb[8:] // encrypted key

	// Convert string pw to []byte
	passWd := []byte(password)

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
		return &wallet, errors.New("invalid password")
	}

	privKey := ed25519.NewKeyFromSeed(unencSecret)
	pubKey := privKey.Public().(ed25519.PublicKey)
	pubKeyBytes := []byte(pubKey)
	signKP := keyPair{PrivKey: privKey, PubKey: pubKeyBytes}

	// public key & secret key
	wallet.Kp = signKP
	wallet.Sk = B58cencode(signKP.PrivKey, edskprefix)
	wallet.Pk = B58cencode(signKP.PubKey, edpkprefix)

	// Generate public address from public key
	generatedAddress, err := generatePublicHash(signKP.PubKey)
	if err != nil {
		return &wallet, errors.Wrapf(err, "could not generate public hash")
	}
	wallet.Address = generatedAddress

	return &wallet, nil
}

// SignOperation will return an operation string signed by wallet
func (w *Wallet) SignOperation(operation string) (SignOperationOutput, error) {

	// Generic watermark
	genericOperationWatermark := []byte{3}

	edsig, err := w.edsig(operation, genericOperationWatermark)
	if err != nil {
		return SignOperationOutput{}, errors.Wrap(err, "failed to sign operation")
	}

	decodedSig, err := decodeSignature(edsig)
	if err != nil {
		return SignOperationOutput{}, errors.Wrap(err, "failed to sign operation")
	}

	return SignOperationOutput{
		SignedOperation: fmt.Sprintf("%s%s", operation, decodedSig),
		Signature:       decodedSig,
		EDSig:           edsig,
	}, nil
}

func (w *Wallet) SignEndorsementOperation(operation, chainID string) (SignOperationOutput, error) {

	// Strip off the chainId prefix, and then base58 decode the chain id string (ie: NetXUdfLh6Gm88t)
	chainIdBytes := b58cdecode(chainID, chainidprefix)

	endorsementWatermark := append(endorsementprefix, chainIdBytes...)

	edsig, err := w.edsig(operation, endorsementWatermark)
	if err != nil {
		return SignOperationOutput{}, errors.Wrap(err, "failed to sign endorsement operation")
	}

	decodedSig, err := decodeSignature(edsig)
	if err != nil {
		return SignOperationOutput{}, errors.Wrap(err, "failed to sign operation")
	}

	return SignOperationOutput{
		SignedOperation: fmt.Sprintf("%s%s", operation, decodedSig),
		Signature:       decodedSig,
		EDSig:           edsig,
	}, nil
}

func (w *Wallet) edsig(operation string, watermark []byte) (string, error) {

	// Prefixes
	edsigByte := []byte{9, 245, 205, 134, 18}

	opBytes, err := hex.DecodeString(operation)
	if err != nil {
		return "", errors.Wrap(err, "failed to sign operation")
	}
	opBytes = append(watermark, opBytes...)

	// Generic hash of 32 bytes
	genericHash, err := blake2b.New(32, []byte{})
	if err != nil {
		return "", errors.Wrap(err, "failed to sign operation bytes")
	}

	// Write operation bytes to hash
	i, err := genericHash.Write(opBytes)

	if err != nil {
		return "", errors.Wrap(err, "failed to sign operation bytes")
	}
	if i != len(opBytes) {
		return "", errors.Errorf("failed to sign operation, generic hash length %d does not match bytes length %d", i, len(opBytes))
	}

	finalHash := genericHash.Sum([]byte{})

	// Sign the finalized generic hash of operations and b58 encode
	sig := ed25519.Sign(w.Kp.PrivKey, finalHash)
	//sig := sodium.Bytes(finalHash).SignDetached(wallet.Kp.PrivKey)
	edsig := B58cencode(sig, edsigByte)

	return edsig, nil
}

//Helper function to return the decoded signature
func decodeSignature(sig string) (string, error) {
	decBytes, err := base58check.Decode(sig)
	if err != nil {
		return "", errors.Wrap(err, "failed to decode signature")
	}

	decodedSigHex := hex.EncodeToString(decBytes)

	// sanity
	if len(decodedSigHex) > 10 {
		decodedSigHex = decodedSigHex[10:]
	} else {
		return "", errors.Wrap(err, "decoded signature is invalid length")
	}

	return decodedSigHex, nil
}

func generatePublicHash(publicKey []byte) (string, error) {
	hash, err := blake2b.New(20, []byte{})
	if err != nil {
		return "", errors.Wrapf(err, "could not generate public hash from public key %s", string(publicKey))
	}
	_, err = hash.Write(publicKey)
	if err != nil {
		return "", errors.Wrapf(err, "could not generate public hash from public key %s", string(publicKey))
	}
	return B58cencode(hash.Sum(nil), tz1prefix), nil
}
