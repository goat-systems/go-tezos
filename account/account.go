package account

import (
	"crypto/sha512"
	"encoding/json"
	"strconv"

	"github.com/DefinitelyNotAGoat/go-tezos/block"

	"github.com/pkg/errors"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/nacl/secretbox"
	"golang.org/x/crypto/pbkdf2"

	tzc "github.com/DefinitelyNotAGoat/go-tezos/client"
	"github.com/DefinitelyNotAGoat/go-tezos/snapshot"

	"github.com/DefinitelyNotAGoat/go-tezos/crypto"
)

const MUTEZ = 1000000

// AccountService is a struct wrapper for account functions
type AccountService struct {
	tzclient        tzc.TezosClient
	snapshotService snapshot.TezosSnapshotService
	blockService    block.TezosBlockService
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
func NewAccountService(tzclient tzc.TezosClient, blockService block.TezosBlockService, snapshotService snapshot.TezosSnapshotService) *AccountService {
	return &AccountService{
		tzclient:        tzclient,
		blockService:    blockService,
		snapshotService: snapshotService,
	}
}

// GetBalanceAtSnapshot gets the balance of a public key hash at a specific snapshot for a cycle.
func (s *AccountService) GetBalanceAtSnapshot(tezosAddr string, cycle int) (float64, error) {
	snapShot, err := s.snapshotService.Get(cycle)
	if err != nil {
		return 0, errors.Wrapf(err, "could not get balance for %s at snapshot at %d cycle", tezosAddr, cycle)
	}

	return s.GetBalanceAtBlock(tezosAddr, snapShot.AssociatedHash)
}

// GetBalance gets the balance of a public key hash at a specific snapshot for a cycle.
func (s *AccountService) GetBalance(tezosAddr string) (float64, error) {

	query := "/chains/main/blocks/head/context/contracts/" + tezosAddr + "/balance"
	resp, err := s.tzclient.Get(query, nil)
	if err != nil {
		return 0, errors.Wrapf(err, "could not get balance '%s'", query)
	}

	strBalance, err := unmarshalString(resp)
	if err != nil {
		return 0, errors.Wrapf(err, "could not get balance '%s'", query)
	}

	floatBalance, err := strconv.ParseFloat(strBalance, 64)
	if err != nil {
		return 0, errors.Wrapf(err, "could not get balance '%s'", query)
	}

	return floatBalance / MUTEZ, nil
}

// GetBalanceAtBlock get the balance of an address at a specific hash
func (s *AccountService) GetBalanceAtBlock(tezosAddr string, id interface{}) (float64, error) {
	blockID, err := s.blockService.IDToString(id)
	if err != nil {
		return 0, errors.Wrapf(err, "could not get balance at block %v", id)
	}

	query := "/chains/main/blocks/" + blockID + "/context/contracts/" + tezosAddr + "/balance"
	resp, err := s.tzclient.Get(query, nil)
	if err != nil {
		return 0, errors.Wrapf(err, "could not get balance at snapshot '%s'", query)
	}

	strBalance, err := unmarshalString(resp)
	if err != nil {
		return 0, errors.Wrapf(err, "could not get balance at snapshot '%s'", query)
	}

	floatBalance, err := strconv.ParseFloat(strBalance, 64)
	if err != nil {
		return 0, errors.Wrapf(err, "could not get balance at snapshot '%s'", query)
	}

	return floatBalance / MUTEZ, nil
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
		return Wallet{}, errors.Wrapf(err, "could not create wallet")
	}

	wallet := Wallet{
		Address:  address,
		Mnemonic: mnenomic,
		Kp:       signKp,
		Seed:     seed,
		Sk:       crypto.B58cencode(privKey, crypto.Prefix_edsk),
		Pk:       crypto.B58cencode(pubKeyBytes, crypto.Prefix_edpk),
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
		return wallet, errors.New("could not import wallet, prefix edsk not found")
	}

	// Determine if 'secret' is an actual secret key or a seed
	if secretLength == 98 {

		// A full secret key
		decodedSecretKey := crypto.B58cdecode(secret, crypto.Prefix_edsk)

		// Public key is last 32 of decoded secret, re-encoded as edpk
		publicKey := decodedSecretKey[32:]

		signKP.PubKey = []byte(publicKey)
		signKP.PrivKey = []byte(secret)

		wallet.Sk = secret

	} else if secretLength == 54 {

		// "secret" is actually a seed
		decodedSeed := crypto.B58cdecode(secret, crypto.Prefix_edsk2)

		//signSeed := sodium.SignSeed{Bytes: decodedSeed}

		// Reconstruct keypair from seed
		privKey := ed25519.NewKeyFromSeed(decodedSeed)
		pubKey := privKey.Public().(ed25519.PublicKey)
		signKP.PrivKey = privKey
		signKP.PubKey = []byte(pubKey)

		wallet.Sk = crypto.B58cencode(signKP.PrivKey, crypto.Prefix_edsk)

	} else {
		return wallet, errors.Errorf("could not import wallet, secret key  length '%d' does not = '%d'", 54, secretLength)
	}

	wallet.Kp = signKP

	// Generate public address from public key
	generatedAddress, err := s.generatePublicHash(signKP.PubKey)
	if err != nil {
		return wallet, errors.Wrapf(err, "could not import wallet, failed to generate public hash")
	}

	if generatedAddress != address {
		return wallet, errors.Errorf("could not import wallet, reconstructed address '%s' does not match provided address '%s'", generatedAddress, address)
	}

	wallet.Address = generatedAddress

	// Genrate and check public key
	generatedPublicKey := crypto.B58cencode(signKP.PubKey, crypto.Prefix_edpk)
	if generatedPublicKey != public {
		return wallet, errors.Errorf("could not import wallet, reconstructed phk '%s' does not match provided phk '%s'", generatedPublicKey, public)
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
		return wallet, errors.New("could not encrypted import wallet, encrypted secret key does not prefix with edesk")
	}

	// Convert key from base58 to []byte
	b58c, err := crypto.Decode(encKey)
	if err != nil {
		return wallet, errors.Wrap(err, "could not encrypted import wallet, encrypted key is not base58")
	}

	// Strip off prefix and extract parts
	esb := b58c[len(crypto.Prefix_edesk):]
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
		return wallet, errors.New("could not encrypted import wallet, invalid password")
	}

	privKey := ed25519.NewKeyFromSeed(unencSecret)
	pubKey := privKey.Public().(ed25519.PublicKey)
	pubKeyBytes := []byte(pubKey)
	signKP := keyPair{PrivKey: privKey, PubKey: pubKeyBytes}

	// public key & secret key
	wallet.Kp = signKP
	wallet.Sk = crypto.B58cencode(signKP.PrivKey, crypto.Prefix_edsk)
	wallet.Pk = crypto.B58cencode(signKP.PubKey, crypto.Prefix_edpk)

	// Generate public address from public key
	generatedAddress, err := s.generatePublicHash(signKP.PubKey)
	if err != nil {
		return wallet, errors.Wrapf(err, "could not import encrypted wallet, failed to generate public hash")
	}
	wallet.Address = generatedAddress

	return wallet, nil
}

func (s *AccountService) generatePublicHash(publicKey []byte) (string, error) {
	hash, err := blake2b.New(20, []byte{})
	hash.Write(publicKey)
	if err != nil {
		return "", errors.Wrapf(err, "could not generate public hash from public key %s", string(publicKey))
	}
	return crypto.B58cencode(hash.Sum(nil), crypto.Prefix_tz1), nil
}

// unmarshalString unmarshals the bytes received as a parameter, into the type string.
func unmarshalString(v []byte) (string, error) {
	var str string
	err := json.Unmarshal(v, &str)
	if err != nil {
		return str, errors.Wrap(err, "could not unmarshal bytes to string")
	}
	return str, nil
}
