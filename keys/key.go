package keys

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	tzcrypt "github.com/completium/go-tezos/v4/internal/crypto"
	"github.com/pkg/errors"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/nacl/secretbox"
	"golang.org/x/crypto/pbkdf2"
)

// Key is the cryptographic key to a Tezos Wallet
type Key struct {
	curve   iCurve
	privKey []byte
	PubKey  PubKey
}

/*
Generate returns a new cryptographic key based on the kind of elliptical curve passed
	* Ed25519
	* Secp256k1
	* NistP256
*/
func Generate(kind ECKind) (*Key, error) {
	token := make([]byte, 32)
	rand.Read(token)
	return key(token, kind)
}

// FromBytes returns a new key from a private key in byte form
func FromBytes(privKey []byte, kind ECKind) (*Key, error) {
	return key(privKey, kind)
}

// FromHex returns a new key from a private key in hex form
func FromHex(privKey string, kind ECKind) (*Key, error) {
	v, err := hex.DecodeString(privKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to import key")
	}

	return key(v, kind)
}

// FromBase64 returns a new key from a private key in base64 form
func FromBase64(privKey string, kind ECKind) (*Key, error) {
	v, err := base64.StdEncoding.DecodeString(privKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to import key")
	}
	return key(v, kind)
}

// FromBase58 returns a new key from a private key in base58 form
func FromBase58(privKey string, kind ECKind) (*Key, error) {
	if len(privKey) < 4 {
		return nil, errors.New("failed to import key: invalid key length")
	}

	curve, err := getCurveByPrefix(privKey[0:4])
	if err != nil {
		return nil, errors.Wrap(err, "failed to import key")
	}

	return key(tzcrypt.B58cdecode(privKey, curve.privateKeyPrefix()), curve.getECKind())
}

func FromBase58Pk(pubKey string, kind ECKind) (*PubKey, error) {
	if len(pubKey) < 4 {
		return nil, errors.New("failed to import pub key: invalid key length")
	}

	curve, err := getCurveByPrefix(pubKey[0:4])
	if err != nil {
		return nil, errors.Wrap(err, "failed to import key")
	}
	pk := tzcrypt.B58cdecode(pubKey, curve.publicKeyPrefix())

	hash, err := blake2b.New(20, []byte{})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to import pub key: failed to generate public hash from public key %s", string(pk))
	}
	_, err = hash.Write(pk)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to import pub key: failed to generate public hash from public key %s", string(pk))
	}

	return &PubKey{
		curve:   curve,
		pubKey:  pk,
		address: tzcrypt.B58cencode(hash.Sum(nil), curve.addressPrefix()),
	}, nil

}

// FromEncryptedSecret returns a new key from an encrypted private key
func FromEncryptedSecret(esk, passwd string) (*Key, error) {
	curve, err := getCurveByPrefix(esk[:5])
	if err != nil {
		return &Key{}, err
	}

	// Convert key from base58 to []byte
	b58c, err := tzcrypt.Decode(esk)
	if err != nil {
		return &Key{}, errors.Wrap(err, "failed to import key")
	}

	// Strip off prefix and extract parts
	esb := b58c[5:]
	salt := esb[:8]
	esm := esb[8:] // encrypted key

	// Convert string pw to []byte
	passWd := []byte(passwd)

	// Derive a key from password, salt and number of iterations
	pbkdf2key := pbkdf2.Key(passWd, salt, 32768, 32, sha512.New)
	var byteKey [32]byte
	for i := range pbkdf2key {
		byteKey[i] = pbkdf2key[i]
	}

	var out []byte
	var emptyNonceBytes [24]byte

	unencSecret, ok := secretbox.Open(out, esm, &emptyNonceBytes, &byteKey)
	if !ok {
		return &Key{}, errors.New("failed to import key: invalid password")
	}

	return key(unencSecret, curve.getECKind())
}

// FromMnemonic returns a new key from a mnemonic
func FromMnemonic(mnemonic, email, passwd string, kind ECKind) (*Key, error) {
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, fmt.Sprintf("%s%s", email, passwd))
	if err != nil {
		return &Key{}, err
	}

	return key(seed, kind)
}

func key(v []byte, kind ECKind) (*Key, error) {
	curve := getCurve(kind)
	pubKey, err := newPubKey(curve.getPrivateKey(v), kind)
	if err != nil {
		return nil, errors.Wrap(err, "failed to import key")
	}

	return &Key{
		curve:   curve,
		privKey: curve.getPrivateKey(v),
		PubKey:  pubKey,
	}, nil
}

// GetBytes will return the raw bytes of the private key.
func (k *Key) GetBytes() []byte {
	return k.privKey
}

/*
GetSecretKey will return the base58 encoded key with the secret key prefix. This is unencrypted.
Example: edskRpfRbhVr7SjmVpMK1kzTDrSzuCKroxjQAsfJn94X7LgbpqJLvRDHfNHFT9KbCZAXVVhMmkQGz4APscezMbJFov5ZNPSY9H
*/
func (k *Key) GetSecretKey() string {
	return tzcrypt.B58cencode(k.privKey, k.curve.privateKeyPrefix())
}

// SignHex will sign a hex encoded string for operation
func (k *Key) SignHex(msg string) (Signature, error) {
	bytes, err := hex.DecodeString(msg)
	if err != nil {
		return Signature{}, errors.Wrap(err, "failed to hex decode message")
	}

	return k.curve.sign(checkAndAddWaterMark(bytes), k.privKey)
}

// SignHex will sign a hex encoded string
func (k *Key) SignDataHex(msg string) (Signature, error) {
	bytes, err := hex.DecodeString(msg)
	if err != nil {
		return Signature{}, errors.Wrap(err, "failed to hex decode message")
	}

	return k.curve.sign(bytes, k.privKey)
}

// SignBytes will sign a byte message for operation
func (k *Key) SignBytes(msg []byte) (Signature, error) {
	return k.curve.sign(checkAndAddWaterMark(msg), k.privKey)
}

// SignBytes will sign a byte message
func (k *Key) SignDataBytes(msg []byte) (Signature, error) {
	return k.curve.sign(msg, k.privKey)
}

func checkAndAddWaterMark(v []byte) []byte {
	if v != nil {
		if v[0] != byte(3) {
			v = append([]byte{3}, v...)
		}
	}

	return v
}

func GetPkhFromBytes(b []byte) (string, error) {

	curve := iCurve(nil)
	if b[0] == 0 && b[1] == 0 {
		curve = getCurve(Ed25519)
	} else if b[0] == 0 && b[1] == 1 {
		curve = getCurve(Secp256k1)
	} else if b[0] == 0 && b[1] == 2 {
		curve = getCurve(NistP256)
	}

	if curve != nil {
		input := b[2:22]
		return tzcrypt.B58cencode(input, curve.addressPrefix()), nil
	} else if b[0] == 1 && b[21] == 0 {
		input := b[1:21]
		return tzcrypt.B58cencode(input, []byte{2, 90, 121}), nil
	}

	return "", errors.New("GetPkhFromBytes: Unknown hash")
}

// SignBytes will sign a byte message for operation
func (k *Key) CheckSignature(data string, signature string) (bool, error) {
	if len(signature) < 5 {
		return false, errors.New("failed to check signature: invalid signature length")
	}

	curve, err := getCurveByPrefix(signature[:5])
	if err != nil {
		return false, err
	}
	sig := tzcrypt.B58cdecode(signature, curve.signaturePrefix())

	msg, err := hex.DecodeString(data)
	if err != nil {
		return false, errors.New("CheckSignature: cannot decode data")
	}

	hash, err := blake2b.New(32, []byte{})
	if err != nil {
		return false, err
	}

	i, err := hash.Write(msg)
	if err != nil {
		return false, errors.Wrap(err, "failed to sign operation bytes")
	}
	if i != len(msg) {
		return false, errors.Errorf("failed to sign operation: generic hash length %d does not match bytes length %d", i, len(msg))
	}

	res, err := k.curve.checkSignature(k.PubKey.pubKey, hash.Sum(nil), sig)
	if err != nil {
		return false, err
	}

	return res, nil
}
