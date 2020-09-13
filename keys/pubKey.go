package keys

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"

	"github.com/go-playground/validator"
	tzcrypt "github.com/goat-systems/go-tezos/v3/internal/crypto"
	"github.com/pkg/errors"
	"golang.org/x/crypto/blake2b"
)

/*
NewPubKeyInput is the input for the keys.NewPubKey function.

Notes:
	The public key can be derived by Bytes, or a String that is encoded in hex, base64, or base58.

Function:
	func NewPubKey(input NewPubKeyInput) (PubKey, error) {}
*/
type NewPubKeyInput struct {
	Bytes  []byte
	String string
	Kind   ECKind `validate:"required"`
}

// PubKey is the public key to a Tezos Wallet
type PubKey struct {
	curve   iCurve
	pubKey  []byte
	address string
}

/*
GeneratePubKey returns a new cryptographic key based on the kind of elliptical curve passed
	* Ed25519
	* Secp256k1
	* NistP256
*/
func GeneratePubKey(kind ECKind) (PubKey, error) {
	token := make([]byte, 32)
	rand.Read(token)
	return newPubKey(token, kind)
}

// NewPubKey gets a new PubKey based on the input passed
func NewPubKey(input NewPubKeyInput) (PubKey, error) {
	err := validator.New().Struct(input)
	if err != nil {
		return PubKey{}, errors.Wrap(err, "invalid input")
	}

	if input.Bytes != nil {
		return newPubKey(input.Bytes, input.Kind)
	}

	if v, err := hex.DecodeString(input.String); err == nil {
		return newPubKey(v, input.Kind)
	}

	if v, err := base64.StdEncoding.DecodeString(input.String); err == nil {
		return newPubKey(v, input.Kind)
	}

	//base58
	if curve, err := getCurveByPrefix(input.String[0:4]); err == nil {
		newPubKey(tzcrypt.B58cdecode(input.String, curve.privateKeyPrefix()), curve.getECKind())
	}

	return PubKey{}, errors.New("unsupported encoding: not hex: not base64: not base58")
}

func newPubKey(v []byte, kind ECKind) (PubKey, error) {
	if len(v) < 32 {
		return PubKey{}, errors.New("invalid bytes length")
	}

	curve := getCurve(kind)
	pk, err := curve.getPublicKey(v)
	if err != nil {
		return PubKey{}, err
	}

	hash, err := blake2b.New(20, []byte{})
	if err != nil {
		return PubKey{}, errors.Wrapf(err, "could not generate public hash from public key %s", string(pk))
	}
	_, err = hash.Write(pk)
	if err != nil {
		return PubKey{}, errors.Wrapf(err, "could not generate public hash from public key %s", string(pk))
	}

	return PubKey{
		curve:   curve,
		pubKey:  pk,
		address: tzcrypt.B58cencode(hash.Sum(nil), curve.addressPrefix()),
	}, nil
}

// GetBytes will return the raw bytes of the public key
func (p *PubKey) GetBytes() []byte {
	return p.pubKey
}

/*
GetPublicKey will return the base58 encoded key with the private key prefix.
Example:
	edskRpfRbhVr7SjmVpMK1kzTDrSzuCKroxjQAsfJn94X7LgbpqJLvRDHfNHFT9KbCZAXVVhMmkQGz4APscezMbJFov5ZNPSY9H
*/
func (p *PubKey) GetPublicKey() string {
	return tzcrypt.B58cencode(p.pubKey, p.curve.publicKeyPrefix())
}

/*
GetPublicKeyHash will public key hash (address) of the public key.
Example:
	tz1L8fUQLuwRuywTZUP5JUw9LL3kJa8LMfoo
*/
func (p *PubKey) GetPublicKeyHash() string {
	return p.address
}

// Verify will verify the authenticity of the public key, signature and data.
func (p *PubKey) Verify(input VerifyInput) bool {
	if input.BytesData != nil && input.BytesSignature != nil {
		return p.curve.verify(input.BytesData, input.BytesSignature, p.pubKey)
	} else if input.Data != "" && input.Signature != "" {
		return p.curve.verify([]byte(input.Data), tzcrypt.B58cdecode(input.Signature, p.curve.signaturePrefix()), p.pubKey)
	}
	return false
}
