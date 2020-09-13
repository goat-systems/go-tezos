package keys

import (
	"crypto/ed25519"
)

var _ iCurve = &secp256k1Curve{}

type secp256k1Curve struct{}

func (e *secp256k1Curve) addressPrefix() []byte {
	return []byte{6, 161, 159}
}

func (e *secp256k1Curve) publicKeyPrefix() []byte {
	return []byte{13, 15, 37, 217}
}

func (e *secp256k1Curve) privateKeyPrefix() []byte {
	return []byte{43, 246, 78, 7}
}

func (e *secp256k1Curve) signaturePrefix() []byte {
	return []byte{9, 245, 205, 134, 18}
}

func (e *secp256k1Curve) getECKind() ECKind {
	return Ed25519
}

func (e *secp256k1Curve) getPrivateKey(v []byte) []byte {
	ed25519.GenerateKey(nil)
	return []byte{}
}

func (e *secp256k1Curve) getPublicKey(privateKey []byte) ([]byte, error) {
	return []byte{}, nil
}

func (e *secp256k1Curve) sign(msg []byte, privateKey []byte) (Signature, error) {
	return Signature{}, nil
}

func (e *secp256k1Curve) verify(v []byte, signature []byte, pubKey []byte) bool {
	return false
}
