package keys

import (
	"encoding/hex"

	"github.com/goat-systems/go-tezos/v3/internal/crypto"
)

type Signature struct {
	Bytes  []byte
	Prefix []byte
}

func (s *Signature) ToBytes() []byte {
	return s.Bytes
}

func (s *Signature) ToBase58() string {
	return crypto.B58cencode(s.Bytes, s.Prefix)
}

func (s *Signature) ToHex() string {
	return hex.EncodeToString(s.Bytes)
}
