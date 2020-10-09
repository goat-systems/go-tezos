package keys

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/go-playground/validator"
	tzcrypt "github.com/goat-systems/go-tezos/v3/crypto"
	"github.com/pkg/errors"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/nacl/secretbox"
	"golang.org/x/crypto/pbkdf2"
)

type fromMnemonicInput struct {
	Mnemonic string `validate:"required"`
	Email    string
	Password string
	Kind     ECKind `validate:"required"`
}

/*
NewKeyInput is the input for the keys.NewKey function.

Note:
	You can generate a key with the following combinations:
		* Bytes
		* EncodedString: hex, base64, or base58
		* Esk & Password
		* Mnemonic
		* Mnemonic & Password
		* Mnemonic & Password & Email

Function:
	func NewKey(input NewKeyInput) (Key, error) {}
*/
type NewKeyInput struct {
	Bytes         []byte
	EncodedString string
	Esk           string
	Password      string
	Mnemonic      string
	Email         string

	Kind ECKind `validate:"required"`
}

/*
VerifyInput is the input for the key.Verify function.

Note:
	You can verify with the following combinations:
		* BytesData & BytesSignature
		* Data & Signature

Function:
	func Verify(input VerifyInput) bool {}
*/
type VerifyInput struct {
	BytesData      []byte
	BytesSignature []byte
	Data           string
	Signature      string
}

/*
SignInput is the input for the key.Sign function.

Note:
	You can sign with the following combinations:
		* Message
			or
		* Bytes

Function:
	func FromMnemonic(input FromMnemonicInput) (Key, error) {}
*/
type SignInput struct {
	Message string
	Bytes   []byte
}

// Key is the cryptographic key to a Tezos Wallet
type Key struct {
	curve   iCurve
	privKey []byte
	PubKey  PubKey
}

/*
GenerateKey returns a new cryptographic key based on the kind of elliptical curve passed
	* Ed25519
	* Secp256k1
	* NistP256
*/
func GenerateKey(kind ECKind) (Key, error) {
	token := make([]byte, 32)
	rand.Read(token)
	return key(token, kind)
}

// NewKey gets a new Key based on the input passed
func NewKey(input NewKeyInput) (Key, error) {
	err := validator.New().Struct(input)
	if err != nil {
		return Key{}, errors.Wrap(err, "invalid input")
	}

	if input.Bytes != nil {
		return key(input.Bytes, input.Kind)
	}

	if input.EncodedString != "" {

		if v, err := hex.DecodeString(input.EncodedString); err == nil {
			return key(v, input.Kind)
		}

		if v, err := base64.StdEncoding.DecodeString(input.EncodedString); err == nil {
			return key(v, input.Kind)
		}

		fmt.Println("HERE")

		//base58
		if curve, err := getCurveByPrefix(input.EncodedString[0:4]); err == nil {
			return key(tzcrypt.B58cdecode(input.EncodedString, curve.privateKeyPrefix()), curve.getECKind())
		}
	}

	// esk with password
	if input.Esk != "" && input.Password != "" {
		return fromEsk(input.Esk, input.Password)
	}

	if input.Mnemonic != "" {
		return fromMnemonic(fromMnemonicInput{
			Mnemonic: input.Mnemonic,
			Email:    input.Email,
			Password: input.Password,
			Kind:     input.Kind,
		})
	}

	return Key{}, errors.New("unsupported")
}

func key(v []byte, kind ECKind) (Key, error) {
	curve := getCurve(kind)
	pubKey, err := newPubKey(curve.getPrivateKey(v), kind)
	if err != nil {
		return Key{}, err
	}

	return Key{
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

// Sign will either sign a hex encoded string or bytes with Key
func (k *Key) Sign(input SignInput) (Signature, error) {
	if input.Bytes != nil {
		return k.curve.sign(checkAndAddWaterMark(input.Bytes), k.privKey)
	} else if input.Message != "" {
		bytes, err := hex.DecodeString(input.Message)
		if err != nil {
			return Signature{}, errors.Wrap(err, "failed to hex decode message")
		}

		return k.curve.sign(checkAndAddWaterMark(bytes), k.privKey)
	}

	return Signature{}, errors.New("missing Bytes or Message in input")
}

func checkAndAddWaterMark(v []byte) []byte {
	if v != nil {
		if v[0] != byte(3) {
			v = append([]byte{3}, v...)
		}
	}

	return v
}

// Verify will verify the authenticity of the public key, signature and data.
func (k *Key) Verify(input VerifyInput) bool {
	return k.PubKey.Verify(input)
}

func fromEsk(esk string, password string) (Key, error) {
	curve, err := getCurveByPrefix(esk[:5])
	if err != nil {
		return Key{}, err
	}

	// Convert key from base58 to []byte
	b58c, err := tzcrypt.Decode(esk)
	if err != nil {
		return Key{}, errors.Wrap(err, "encrypted key is not base58")
	}

	// Strip off prefix and extract parts
	esb := b58c[5:]
	salt := esb[:8]
	esm := esb[8:] // encrypted key

	// Convert string pw to []byte
	passWd := []byte(password)

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
		return Key{}, errors.New("invalid password")
	}

	return key(unencSecret, curve.getECKind())
}

// fromMnemonic generates a new Key based off the mnemonic input passed.
func fromMnemonic(input fromMnemonicInput) (Key, error) {
	err := validator.New().Struct(input)
	if err != nil {
		return Key{}, errors.Wrap(err, "invalid input")
	}

	seed, err := bip39.NewSeedWithErrorChecking(input.Mnemonic, fmt.Sprintf("%s%s", input.Email, input.Password))
	if err != nil {
		return Key{}, err
	}

	return key(seed, input.Kind)
}
