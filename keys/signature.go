package keys

import (
	"encoding/hex"

	"github.com/utdrmac/go-tezos/v3/crypto"
)

// Signature represents the signature of an operation
type Signature struct {
	Bytes  []byte
	Prefix []byte
}

// ToBytes returns the signature as bytes
func (s *Signature) ToBytes() []byte {
	return s.Bytes
}

// ToBase58 returns the signature as a base58 encoded string with the correct prefix
func (s *Signature) ToBase58() string {
	return crypto.B58cencode(s.Bytes, s.Prefix)
}

// ToHex returns the signature encoded to hex
func (s *Signature) ToHex() string {
	return hex.EncodeToString(s.Bytes)
}
