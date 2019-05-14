package gotezos

import (
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/Messer4/base58check"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/nacl/secretbox"
	"golang.org/x/crypto/pbkdf2"
)

// AccountService is a struct wrapper for account functions
type AccountService struct {
	gt *GoTezos
}

// FrozenBalance is representation of frozen balance on the Tezos network
type FrozenBalance struct {
	Deposits string `json:"deposits"`
	Fees     string `json:"fees"`
	Rewards  string `json:"rewards"`
}

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

// NewAccountService returns a new AccountService
func (gt *GoTezos) newAccountService() *AccountService {
	return &AccountService{gt: gt}
}

// GetBalanceAtSnapshot gets the balance of a public key hash at a specific snapshot for a cycle.
func (s *AccountService) GetBalanceAtSnapshot(tezosAddr string, cycle int) (float64, error) {
	snapShot, err := s.gt.SnapShot.Get(cycle)
	if err != nil {
		return 0, fmt.Errorf("could not get %s balance for snap shot at %d cycle: %v", tezosAddr, cycle, err)
	}

	query := "/chains/main/blocks/" + snapShot.AssociatedHash + "/context/contracts/" + tezosAddr + "/balance"
	resp, err := s.gt.Get(query, nil)
	if err != nil {
		return 0, fmt.Errorf("could not get %s balance for snap shot at %d cycle: %v", tezosAddr, cycle, err)
	}

	strBalance, err := unmarshalString(resp)
	if err != nil {
		return 0, fmt.Errorf("could not get %s balance for snap shot at %d cycle: %v", tezosAddr, cycle, err)
	}

	floatBalance, err := strconv.ParseFloat(strBalance, 64)
	if err != nil {
		return 0, fmt.Errorf("could not get %s balance for snap shot at %d cycle: %v", tezosAddr, cycle, err)
	}

	return floatBalance / MUTEZ, nil
}

// GetBalance gets the balance of a public key hash at a specific snapshot for a cycle.
func (s *AccountService) GetBalance(tezosAddr string) (float64, error) {

	query := "/chains/main/blocks/head/context/contracts/" + tezosAddr + "/balance"
	resp, err := s.gt.Get(query, nil)
	if err != nil {
		return 0, fmt.Errorf("could not get %s balance: %v", tezosAddr, err)
	}

	strBalance, err := unmarshalString(resp)
	if err != nil {
		return 0, fmt.Errorf("could not get %s balance: %v", tezosAddr, err)
	}

	floatBalance, err := strconv.ParseFloat(strBalance, 64)
	if err != nil {
		return 0, fmt.Errorf("could not get %s balance: %v", tezosAddr, err)
	}

	return floatBalance / MUTEZ, nil
}

// GetBalanceAtBlock get the balance of an address at a specific hash
func (s *AccountService) GetBalanceAtBlock(tezosAddr string, id interface{}) (int, error) {
	var balance string
	block, err := s.gt.Block.Get(id)
	if err != nil {
		return 0, fmt.Errorf("could not get balance at block %v: %v", id, err)
	}

	query := "/chains/main/blocks/" + block.Hash + "/context/contracts/" + tezosAddr + "/balance"

	resp, err := s.gt.Get(query, nil)
	if err != nil {
		return 0, err
	}
	balance, err = unmarshalString(resp)
	if err != nil {
		return 0, err
	}

	var returnBalance int
	if strings.Contains(balance, "No service found at gt URL") {
		returnBalance = 0
	}

	if len(balance) < 1 {
		returnBalance = 0
	} else {
		floatBalance, _ := strconv.Atoi(balance) //TODO error checking
		returnBalance = int(floatBalance)
	}

	return returnBalance, nil
}

// CreateWallet returns Wallet with the mnemonic and password provided
func (s *AccountService) CreateWallet(mnenomic string, password string) (Wallet, error) {

	seed := pbkdf2.Key([]byte(mnenomic), []byte("mnemonic"+password), 2048, 32, sha512.New)
	privKey := ed25519.NewKeyFromSeed(seed)
	pubKey := privKey.Public().(ed25519.PublicKey)
	pubKeyBytes := []byte(pubKey)
	signKp := keyPair{PrivKey: privKey, PubKey: pubKeyBytes}

	address, err := s.generatePublicHash(pubKeyBytes)
	if err != nil {
		return Wallet{}, fmt.Errorf("could not create wallet: %v", err)
	}

	wallet := Wallet{
		Address:  address,
		Mnemonic: mnenomic,
		Kp:       signKp,
		Seed:     seed,
		Sk:       b58cencode(privKey, edsk),
		Pk:       b58cencode(pubKeyBytes, edpk),
	}

	return wallet, nil
}

// ImportWallet returns an imported Wallet
func (s *AccountService) ImportWallet(address, public, secret string) (Wallet, error) {

	var wallet Wallet
	var signKP keyPair

	// Sanity check
	secretLength := len(secret)
	if secret[:4] != "edsk" || (secretLength != 98 && secretLength != 54) {
		return wallet, fmt.Errorf("could not import wallet: prefix edsk not found")
	}

	// Determine if 'secret' is an actual secret key or a seed
	if secretLength == 98 {

		// A full secret key
		decodedSecretKey := b58cdecode(secret, edsk)

		// Public key is last 32 of decoded secret, re-encoded as edpk
		publicKey := decodedSecretKey[32:]

		signKP.PubKey = []byte(publicKey)
		signKP.PrivKey = []byte(secret)

		wallet.Sk = secret

	} else if secretLength == 54 {

		// "secret" is actually a seed
		decodedSeed := b58cdecode(secret, edsk2)

		//signSeed := sodium.SignSeed{Bytes: decodedSeed}

		// Reconstruct keypair from seed
		privKey := ed25519.NewKeyFromSeed(decodedSeed)
		pubKey := privKey.Public().(ed25519.PublicKey)
		signKP.PrivKey = privKey
		signKP.PubKey = []byte(pubKey)

		wallet.Sk = b58cencode(signKP.PrivKey, edsk)

	} else {
		return wallet, fmt.Errorf("could not import wallet: secret key is not the correct length")
	}

	wallet.Kp = signKP

	// Generate public address from public key
	generatedAddress, err := s.generatePublicHash(signKP.PubKey)
	if err != nil {
		return wallet, fmt.Errorf("could not import wallet: %v", err)
	}

	if generatedAddress != address {
		return wallet, fmt.Errorf("could not import wallet: reconstructed address '%s' and provided address '%s' do not match", generatedAddress, address)
	}

	wallet.Address = generatedAddress

	// Genrate and check public key
	generatedPublicKey := b58cencode(signKP.PubKey, edpk)
	if generatedPublicKey != public {
		return wallet, fmt.Errorf("could not import wallet: reconstructed phk '%s' and provided phk '%s' do not match", generatedPublicKey, public)
	}
	wallet.Pk = generatedPublicKey

	return wallet, nil
}

// ImportEncryptedWallet imports an encrypted wallet using password provided by caller.
// Caller should remove any 'encrypted:' scheme prefix.
func (s *AccountService) ImportEncryptedWallet(pw, encKey string) (Wallet, error) {

	var wallet Wallet

	// Check if user copied 'encrypted:' scheme prefix
	if encKey[:5] != "edesk" || len(encKey) != 88 {
		return wallet, fmt.Errorf("could not import wallet: encrypted secret key does not prefix with edesk")
	}

	// Convert key from base58 to []byte
	b58c, err := base58check.Decode(encKey)
	if err != nil {
		return wallet, fmt.Errorf("could not import wallet: %v", err)
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
		return wallet, fmt.Errorf("could not import wallet: invalid password")
	}

	privKey := ed25519.NewKeyFromSeed(unencSecret)
	pubKey := privKey.Public().(ed25519.PublicKey)
	pubKeyBytes := []byte(pubKey)
	signKP := keyPair{PrivKey: privKey, PubKey: pubKeyBytes}

	// public key & secret key
	wallet.Kp = signKP
	wallet.Sk = b58cencode(signKP.PrivKey, edsk)
	wallet.Pk = b58cencode(signKP.PubKey, edpk)

	// Generate public address from public key
	generatedAddress, err := s.generatePublicHash(signKP.PubKey)
	if err != nil {
		return wallet, fmt.Errorf("could not import wallet: %v", err)
	}
	wallet.Address = generatedAddress

	return wallet, nil
}

func (s *AccountService) generatePublicHash(publicKey []byte) (string, error) {
	hash, err := blake2b.New(20, []byte{})
	hash.Write(publicKey)
	if err != nil {
		return "", fmt.Errorf("unable to write public key to generic hash: %v", err)
	}
	return b58cencode(hash.Sum(nil), tz1), nil
}

// unmarshalString unmarshals the bytes received as a parameter, into the type string.
func unmarshalString(v []byte) (string, error) {
	var str string
	err := json.Unmarshal(v, &str)
	if err != nil {
		log.Println("Could not unMarshal to string " + err.Error())
		return str, err
	}
	return str, nil
}
