package keys

import (
	"math/big"
)

var _ iCurve = &secp256k1Curve{}

type secp256k1Curve struct{}

func order() *big.Int {
	i, _ := big.NewInt(0).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 16)
	return i
}

func maxS() *big.Int {
	i, _ := big.NewInt(0).SetString("7FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF5D576E7357A4501DDFE92F46681B20A0", 16)
	return i
}

func (e *secp256k1Curve) addressPrefix() []byte {
	return []byte{6, 161, 161}
}

func (e *secp256k1Curve) publicKeyPrefix() []byte {
	return []byte{3, 254, 226, 86}
}

func (e *secp256k1Curve) privateKeyPrefix() []byte {
	return []byte{17, 162, 224, 201}
}

func (e *secp256k1Curve) signaturePrefix() []byte {
	return []byte{13, 115, 101, 19, 63}
}

func (e *secp256k1Curve) getECKind() ECKind {
	return Secp256k1
}

func (e *secp256k1Curve) getPrivateKey(v []byte) []byte {
	return v[:32]
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
