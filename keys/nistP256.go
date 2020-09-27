package keys

import (
	"crypto/ed25519"
)

var _ iCurve = &nistP256Curve{}

type nistP256Curve struct{}

func (e *nistP256Curve) addressPrefix() []byte {
	return []byte{6, 161, 161}
}

func (e *nistP256Curve) publicKeyPrefix() []byte {
	return []byte{3, 254, 226, 86}
}

func (e *nistP256Curve) privateKeyPrefix() []byte {
	return []byte{17, 162, 224, 201}
}

func (e *nistP256Curve) signaturePrefix() []byte {
	return []byte{13, 115, 101, 19, 63}
}

func (e *nistP256Curve) getECKind() ECKind {
	return Ed25519
}

func (e *nistP256Curve) getPrivateKey(v []byte) []byte {
	ed25519.GenerateKey(nil)
	return []byte{}
}

func (e *nistP256Curve) getPublicKey(privateKey []byte) ([]byte, error) {
	return []byte{}, nil
}

func (e *nistP256Curve) sign(msg []byte, privateKey []byte) (Signature, error) {
	return Signature{}, nil
}

func (e *nistP256Curve) verify(v []byte, signature []byte, pubKey []byte) bool {
	return false
}
