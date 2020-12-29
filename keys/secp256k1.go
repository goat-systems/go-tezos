package keys

import (
	"crypto/ecdsa"
	"crypto/rand"
	"math/big"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"golang.org/x/crypto/blake2b"
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

func (s *secp256k1Curve) addressPrefix() []byte {
	return []byte{6, 161, 161}
}

func (s *secp256k1Curve) publicKeyPrefix() []byte {
	return []byte{3, 254, 226, 86}
}

func (s *secp256k1Curve) privateKeyPrefix() []byte {
	return []byte{17, 162, 224, 201}
}

func (s *secp256k1Curve) signaturePrefix() []byte {
	return []byte{13, 115, 101, 19, 63}
}

func (s *secp256k1Curve) getECKind() ECKind {
	return Secp256k1
}

func (s *secp256k1Curve) getPrivateKey(v []byte) []byte {
	return v[:32]
}

func (s *secp256k1Curve) getPublicKey(privateKey []byte) ([]byte, error) {
	privKey, err := ethcrypto.ToECDSA(privateKey)
	if err != nil {
		return []byte{}, err
	}

	var pref []byte
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

func (s *secp256k1Curve) sign(msg []byte, privateKey []byte) (Signature, error) {
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

	privKey, err := ethcrypto.ToECDSA(privateKey)
	if err != nil {
		return Signature{}, err
	}

	r, ss, err := ecdsa.Sign(rand.Reader, privKey, hash.Sum([]byte{}))
	if err != nil {
		return Signature{}, err
	}

	if ss.Cmp(maxS()) > 0 {
		ss = big.NewInt(0).Sub(order(), ss)
	}

	signature := append(r.Bytes(), ss.Bytes()...)
	return Signature{
		Bytes:  signature,
		prefix: s.signaturePrefix(),
	}, nil
}
