package goTezos

import (
	"github.com/GoKillers/libsodium-go/cryptosign"
	"github.com/GoKillers/libsodium-go/cryptogenerichash"
	"gitlab.com/tulpenhaendler/hellotezos/base58check"
	"strings"
)

type KeyPair struct {
	Sk string
	Pk string
	Address string
}


func (this *GoTezos) addPrefix(b []byte, p []byte) []byte {
	p = append(p, b...)
	return p
}

func (this *GoTezos) GenerateAddress() KeyPair {
	var pkhr []byte
	sk, pk, _ := cryptosign.CryptoSignKeyPair()
	pkhr, _ = generichash.CryptoGenericHash(20, pk, []byte{})
	address := base58check.Encode("00", this.addPrefix(pkhr, []byte{6, 161, 159}))


	res := KeyPair{
		Sk:base58check.Encode("00", this.addPrefix(sk, []byte{43, 246, 78, 7})),
		Pk:base58check.Encode("00", this.addPrefix(pk, []byte{13, 15, 37, 217})),
		Address:address,
	}
	return res
}

func (this *GoTezos) VanityAddressPrefix(prefix string) KeyPair {
	for {
		addr := this.GenerateAddress()
		if strings.HasPrefix(addr.Address,prefix) {
			return addr
		}
	}
}