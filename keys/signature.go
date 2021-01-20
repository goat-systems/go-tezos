package keys

import (
	"encoding/hex"
	"fmt"

	"github.com/goat-systems/go-tezos/v4/internal/crypto"
)

// Signature represents the signature of an operation
type Signature struct {
	Bytes  []byte
	prefix []byte
}

// ToBytes returns the signature as bytes
func (s *Signature) ToBytes() []byte {
	return s.Bytes
}

// ToBase58 returns the signature as a base58 encoded string with the correct prefix
func (s *Signature) ToBase58() string {
	return crypto.B58cencode(s.Bytes, s.prefix)
}

// ToHex returns the signature encoded to hex
func (s *Signature) ToHex() string {
	return hex.EncodeToString(s.Bytes)
}

// AppendToHex takes a hex encoded message and adds the signature to it for injection
func (s *Signature) AppendToHex(msg string) string {
	return fmt.Sprintf("%s%s", msg, s.ToHex())
}

// AppendToBytes takes a bytes message and adds the signature to it for injection
func (s *Signature) AppendToBytes(msg []byte) []byte {
	return append(msg, s.Bytes...)
}
