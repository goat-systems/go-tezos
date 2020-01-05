package gotezos

import (
	"crypto/sha512"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/nacl/secretbox"
	"golang.org/x/crypto/pbkdf2"
)

//Wallet needed for signing operations
type Wallet struct {
	Address  string
	Mnemonic string
	Seed     []byte
	Kp       keyPair
	Sk       string
	Pk       string
}

// Key Pair Storage
type keyPair struct {
	PrivKey []byte
	PubKey  []byte
}

// Balance gets the balance of an address at the given block hash
func (t *GoTezos) Balance(blockhash, address string) (string, error) {
	query := fmt.Sprintf("/chains/main/blocks/%s/context/contracts/%s/balance", blockhash, address)
	resp, err := t.get(query)
	if err != nil {
		return "", errors.Wrap(err, "failed to get balance")
	}

	var balance string
	err = json.Unmarshal(resp, &balance)
	if err != nil {
		return "", errors.Wrap(err, "failed to unmarshal balance")
	}

	return balance, nil
}

// CreateWallet returns a Wallet with the mnemonic and password provided
func CreateWallet(mnenomic string, password string) (Wallet, error) {

	seed := pbkdf2.Key([]byte(mnenomic), []byte("mnemonic"+password), 2048, 32, sha512.New)
	privKey := ed25519.NewKeyFromSeed(seed)
	pubKey := privKey.Public().(ed25519.PublicKey)
	pubKeyBytes := []byte(pubKey)
	signKp := keyPair{PrivKey: privKey, PubKey: pubKeyBytes}

	address, err := generatePublicHash(pubKeyBytes)
	if err != nil {
		return Wallet{}, errors.Wrapf(err, "could not create wallet")
	}

	wallet := Wallet{
		Address:  address,
		Mnemonic: mnenomic,
		Kp:       signKp,
		Seed:     seed,
		Sk:       b58cencode(privKey, prefix_edsk),
		Pk:       b58cencode(pubKeyBytes, prefix_edpk),
	}

	return wallet, nil
}

// ImportWallet returns an imported Wallet
func ImportWallet(address, public, secret string) (Wallet, error) {

	var wallet Wallet
	var signKP keyPair

	// Sanity check
	secretLength := len(secret)
	if secretLength != 98 && secretLength != 54 {
		return wallet, errors.New("wallet prefix is not edsk")
	}

	if secret[:4] != "edsk" {
		return wallet, errors.New("wallet prefix is not edsk")
	}

	// Determine if 'secret' is an actual secret key or a seed
	if secretLength == 98 {

		// A full secret key
		decodedSecretKey := b58cdecode(secret, prefix_edsk)

		// Public key is last 32 of decoded secret, re-encoded as edpk
		publicKey := decodedSecretKey[32:]

		signKP.PubKey = []byte(publicKey)
		signKP.PrivKey = []byte(secret)

		wallet.Sk = secret

	} else if secretLength == 54 {

		// "secret" is actually a seed
		decodedSeed := b58cdecode(secret, prefix_edsk2)

		//signSeed := sodium.SignSeed{Bytes: decodedSeed}

		// Reconstruct keypair from seed
		privKey := ed25519.NewKeyFromSeed(decodedSeed)
		pubKey := privKey.Public().(ed25519.PublicKey)
		signKP.PrivKey = privKey
		signKP.PubKey = []byte(pubKey)

		wallet.Sk = b58cencode(signKP.PrivKey, prefix_edsk)

	} else {
		return wallet, errors.Errorf("wallet secret key length '%d' does not = '%d'", 54, secretLength)
	}

	wallet.Kp = signKP

	// Generate public address from public key
	generatedAddress, err := generatePublicHash(signKP.PubKey)
	if err != nil {
		return wallet, errors.Wrapf(err, "could not generate public hash")
	}

	if generatedAddress != address {
		return wallet, errors.Errorf("reconstructed address '%s' does not match provided address '%s'", generatedAddress, address)
	}

	wallet.Address = generatedAddress

	// Genrate and check public key
	generatedPublicKey := b58cencode(signKP.PubKey, prefix_edpk)
	if generatedPublicKey != public {
		return wallet, errors.Errorf("reconstructed pk '%s' does not match provided pk '%s'", generatedPublicKey, public)
	}
	wallet.Pk = generatedPublicKey

	return wallet, nil
}

// ImportEncryptedWallet imports an encrypted wallet using the password and encrypted key passed
func ImportEncryptedWallet(pw, encKey string) (Wallet, error) {

	var wallet Wallet
	// Check if user copied 'encrypted:' scheme prefix
	if len(encKey) != 88 {
		return wallet, errors.New("encrypted secret key does not 88 characters long")
	}
	if encKey[:5] != "edesk" {
		return wallet, errors.New("encrypted secret key does not prefix with edesk")
	}

	// Convert key from base58 to []byte
	b58c, err := decode(encKey)
	if err != nil {
		return wallet, errors.Wrap(err, "encrypted key is not base58")
	}

	// Strip off prefix and extract parts
	esb := b58c[len(prefix_edesk):]
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
		return wallet, errors.New("invalid password")
	}

	privKey := ed25519.NewKeyFromSeed(unencSecret)
	pubKey := privKey.Public().(ed25519.PublicKey)
	pubKeyBytes := []byte(pubKey)
	signKP := keyPair{PrivKey: privKey, PubKey: pubKeyBytes}

	// public key & secret key
	wallet.Kp = signKP
	wallet.Sk = b58cencode(signKP.PrivKey, prefix_edsk)
	wallet.Pk = b58cencode(signKP.PubKey, prefix_edpk)

	// Generate public address from public key
	generatedAddress, err := generatePublicHash(signKP.PubKey)
	if err != nil {
		return wallet, errors.Wrapf(err, "could not generate public hash")
	}
	wallet.Address = generatedAddress

	return wallet, nil
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
	return b58cencode(hash.Sum(nil), prefix_tz1), nil
}
