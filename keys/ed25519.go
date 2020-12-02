package keys

import (
	"crypto/ed25519"

	"github.com/pkg/errors"
	"golang.org/x/crypto/blake2b"
)

var _ iCurve = &ed25519Curve{}

// https://tools.ietf.org/html/rfc8032
type ed25519Curve struct{}

func (e *ed25519Curve) addressPrefix() []byte {
	return []byte{6, 161, 159}
}

func (e *ed25519Curve) publicKeyPrefix() []byte {
	return []byte{13, 15, 37, 217}
}

func (e *ed25519Curve) privateKeyPrefix() []byte {
	return []byte{43, 246, 78, 7}
}

func (e *ed25519Curve) signaturePrefix() []byte {
	return []byte{9, 245, 205, 134, 18}
}

func (e *ed25519Curve) getECKind() ECKind {
	return Ed25519
}

func (e *ed25519Curve) getPrivateKey(v []byte) []byte {
	return ed25519.NewKeyFromSeed(v[:32])
}

func (e *ed25519Curve) getPublicKey(privateKey []byte) ([]byte, error) {
	pubKey, ok := ed25519.PrivateKey(privateKey).Public().(ed25519.PublicKey)
	if !ok {
		return []byte{}, errors.New("failed to cast crypto.PublicKey to ed25519.PublicKey")
	}

	return pubKey, nil
}

func (e *ed25519Curve) sign(msg []byte, privateKey []byte) (Signature, error) {
	hash, err := blake2b.New(32, []byte{})
	if err != nil {
		return Signature{}, err
	}

	i, err := hash.Write(msg)
	if err != nil {
		return Signature{}, errors.Wrap(err, "failed to sign operation bytes")
	}
	if i != len(msg) {
		return Signature{}, errors.Errorf("failed to sign operation: generic hash length %d does not match bytes length %d", i, len(msg))
	}

	return Signature{
		Bytes:  ed25519.Sign(ed25519.PrivateKey(privateKey), hash.Sum([]byte{})),
		prefix: e.signaturePrefix(),
	}, nil
}
