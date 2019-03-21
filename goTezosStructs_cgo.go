// +build cgo

package goTezos

import (
	"github.com/jamesruan/sodium"
)

//Wallet needed for signing operations
type Wallet struct {
	Address  string
	Mnemonic string
	Seed     []byte
	Kp       sodium.SignKP
	Sk       string
	Pk       string
}
