package keys

import (
	tzcrypt "github.com/completium/go-tezos/v4/internal/crypto"
	"github.com/pkg/errors"
	"golang.org/x/crypto/blake2b"
)

// PubKey is the public key to a Tezos Wallet
type PubKey struct {
	curve   iCurve
	pubKey  []byte
	address string
}

func newPubKey(v []byte, kind ECKind) (PubKey, error) {
	if len(v) < 32 {
		return PubKey{}, errors.New("failed to import pub key")
	}

	curve := getCurve(kind)
	pk, err := curve.getPublicKey(v)
	if err != nil {
		return PubKey{}, errors.Wrap(err, "failed to import pub key")
	}

	hash, err := blake2b.New(20, []byte{})
	if err != nil {
		return PubKey{}, errors.Wrapf(err, "failed to import pub key: failed to generate public hash from public key %s", string(pk))
	}
	_, err = hash.Write(pk)
	if err != nil {
		return PubKey{}, errors.Wrapf(err, "failed to import pub key: failed to generate public hash from public key %s", string(pk))
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
GetAddress will public key hash (address) of the public key.
Example:
	tz1L8fUQLuwRuywTZUP5JUw9LL3kJa8LMfoo
*/
func (p *PubKey) GetAddress() string {
	return p.address
}
