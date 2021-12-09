package keys

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"math/big"

	"github.com/pkg/errors"
	"golang.org/x/crypto/blake2b"
)

var _ iCurve = &nistP256Curve{}

type nistP256Curve struct{}

func (n *nistP256Curve) addressPrefix() []byte {
	return []byte{6, 161, 164}
}

func (n *nistP256Curve) publicKeyPrefix() []byte {
	return []byte{3, 178, 139, 127}
}

func (n *nistP256Curve) privateKeyPrefix() []byte {
	return []byte{16, 81, 238, 189}
}

func (n *nistP256Curve) signaturePrefix() []byte {
	return []byte{54, 240, 44, 52}
}

func (n *nistP256Curve) getECKind() ECKind {
	return NistP256
}

func (n *nistP256Curve) getPrivateKey(v []byte) []byte {
	return v[:32]
}

func (n *nistP256Curve) getPublicKey(privateKey []byte) ([]byte, error) {
	var privKey ecdsa.PrivateKey
	privKey.D = new(big.Int).SetBytes(privateKey)
	privKey.PublicKey.Curve = elliptic.P256()
	privKey.PublicKey.X, privKey.PublicKey.Y = privKey.PublicKey.Curve.ScalarBaseMult(privKey.D.Bytes())

	pubKeyBytes := make([]byte, 33)

	bY := privKey.PublicKey.Y.Bytes()
	if bY[len(bY)-1]%2 == 0 {
		pubKeyBytes[0] = 2
	} else {
		pubKeyBytes[0] = 3
	}

	// Fill pubKeyBytes[1:] with 0-padded PublicKey.X
	privKey.PublicKey.X.FillBytes(pubKeyBytes[1:])

	return pubKeyBytes, nil
}

func (n *nistP256Curve) sign(msg []byte, privateKey []byte) (Signature, error) {
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

	var privKey ecdsa.PrivateKey
	privKey.D = new(big.Int).SetBytes(privateKey)
	privKey.PublicKey.Curve = elliptic.P256()
	privKey.PublicKey.X, privKey.PublicKey.Y = privKey.PublicKey.Curve.ScalarBaseMult(privKey.D.Bytes())

	r, ss, err := ecdsa.Sign(rand.Reader, &privKey, hash.Sum([]byte{}))
	if err != nil {
		return Signature{}, err
	}

	signature := append(r.Bytes(), ss.Bytes()...)

	return Signature{
		Bytes:  signature,
		prefix: n.signaturePrefix(),
	}, nil
}

func (n *nistP256Curve) checkSignature(pubKey []byte, hash []byte, signature []byte) (bool, error) {
	return false, errors.New("checkSignature nistP256: not implemented")
}
