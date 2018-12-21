package goTezos

import (
	"testing"
)

func TestCreateWalletWithMnemonic(t *testing.T) {
	
	gt := NewGoTezos()
	
	// Test Wallet from Alphanet faucet:
	mnemonic := "normal dash crumble neutral reflect parrot know stairs culture fault check whale flock dog scout"
	password := "PYh8nXDQLB"
	email := "vksbjweo.qsrgfvbw@tezos.example.org"
	
	// These values were gathered after manually importing above mnemonic into CLI wallet
	pkh := "tz1Qny7jVMGiwRrP9FikRK95jTNbJcffTpx1"
	pk := "edpkvEoAbkdaGALxi2FfeefB8hUkMZ4J1UVwkzyumx2GvbVpkYUHnm"
	sk := "edskRxB2DmoyZSyvhsqaJmw5CK6zYT7dbkUfEVSiQeWU1gw3ZMnC99QMMXru3imsbUrLhvuHktrymvNqhMxkhz7Y4LJAtevW5V"
	
	// Alphanet 'password' is email & password concatenated together
	myWallet, err := gt.CreateWallet(mnemonic, email+password)
	if err != nil {
		t.Errorf("Unable to create wallet from Mnemonic: %s", err)
	}
	
	if myWallet.Address != pkh || myWallet.Pk != pk || myWallet.Sk != sk {
		t.Errorf("Created wallet values do not match known answers")
	}
}
	
func TestImportWalletFullSk(t *testing.T) {
	
	gt := NewGoTezos()
	
	pkh := "tz1fYvVTsSQWkt63P5V8nMjW764cSTrKoQKK"
	pk := "edpkvH3h91QHjKtuR45X9BJRWJJmK7s8rWxiEPnNXmHK67EJYZF75G"
	sk := "edskSA4oADtx6DTT6eXdBc6Pv5MoVBGXUzy8bBryi6D96RQNQYcRfVEXd2nuE2ZZPxs4YLZeM7KazUULFT1SfMDNyKFCUgk6vR"
	
	myWallet, err := gt.ImportWallet(pkh, pk, sk)
	if err != nil {
		t.Errorf("%s", err)
	}
	
	if myWallet.Address != pkh || myWallet.Pk != pk || myWallet.Sk != sk {
		t.Errorf("Created wallet values do not match known answers")
	}
}

func TestImportWalletSeedSk(t *testing.T) {
	
	gt := NewGoTezos()
	
	pkh := "tz1U8sXoQWGUMQrfZeAYwAzMZUvWwy7mfpPQ"
	pk  := "edpkunwa7a3Y5vDr9eoKy4E21pzonuhqvNjscT9XG27aQV4gXq4dNm"
	sks := "edsk362Ypv3qLgbnGvZK7JwqNbwiLGe18XhTMFQY4gUonqnaCPiT6X"
	sk  := "edskRjBSseEx9bSRSJJpbypJe5ZXucTtApb6qjechMB1BzEYwcEZyfLooo22Nwk33mPPJ3xZniFoa3o8Js7nNXDdqK9nNjFDi7"
	
	myWallet, err := gt.ImportWallet(pkh, pk, sks)
	if err != nil {
		t.Errorf("%s", err)
	}
	
	if myWallet.Address != pkh || myWallet.Pk != pk || myWallet.Sk != sk {
		t.Errorf("Created wallet values do not match known answers")
	}
}

