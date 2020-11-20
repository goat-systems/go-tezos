package keys

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"math/big"
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

	pref := []byte{}
	if privKey.PublicKey.Y.Bytes()[31]%2 == 0 {
		pref = []byte{2}
	} else {
		pref = []byte{3}
	}

	// 32 padded 0's
	pad := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	pad = append(pad, privKey.PublicKey.X.Bytes()...)

	return append(pref, pad[len(pad)-32:]...), nil
}

func (n *nistP256Curve) sign(msg []byte, privateKey []byte) (Signature, error) {
	return Signature{}, nil
}

func (n *nistP256Curve) verify(v []byte, signature []byte, pubKey []byte) bool {
	return false
}
