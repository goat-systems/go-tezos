package goTezos

import (
	"testing"
)

func TestNewRPCClient(t *testing.T) {
	
	t.Log("RPC client connecting to a Tezos node over localhost, port 8732")
	
	gtClient := NewTezosRPCClient("localhost", "8732")
	gt := NewGoTezos()
	gt.AddNewClient(gtClient)
	
	if ! gtClient.Healthcheck() {
		t.Errorf("Unable to query RPC on 'localhost:8732'. Check that a node is accessible.")
	}
}


func TestNewWebClient(t *testing.T) {
	
	t.Log("Web-based RPC client using https://rpc.tzbeta.net")
	
	gtClient := NewTezosRPCClient("rpc.tzbeta.net", "443")
	gtClient.IsWebClient(true)
	
	gt := NewGoTezos()
	gt.AddNewClient(gtClient)
	
	if ! gtClient.Healthcheck() {
		t.Errorf("Unable to query RPC at 'https://rpc.tzbeta.net'.")
	}
}


func TestCreateWalletWithMnemonic(t *testing.T) {
	
	gt := NewGoTezos()
	
	t.Log("Create new wallet using Alphanet faucet account")
	
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
	
	t.Log("Import existing wallet using complete secret key")
	
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
	
	t.Log("Import existing wallet using seed-secret key")
	
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


func TestImportEncryptedSecret(t *testing.T) {
	
	gt := NewGoTezos()
	
	t.Log("Import wallet using password and encrypted key")
	
	pw := "password12345##"
	sk := "edesk1fddn27MaLcQVEdZpAYiyGQNm6UjtWiBfNP2ZenTy3CFsoSVJgeHM9pP9cvLJ2r5Xp2quQ5mYexW1LRKee2"
	
	// known answers for testing
	pk := "edpkuHMDkMz46HdRXYwom3xRwqk3zQ5ihWX4j8dwo2R2h8o4gPcbN5"
	pkh := "tz1L8fUQLuwRuywTZUP5JUw9LL3kJa8LMfoo"
	
	myWallet, err := gt.ImportEncryptedWallet(pw, sk)
	if err != nil {
		t.Errorf("%s", err)
	}
	
	if myWallet.Address != pkh || myWallet.Pk != pk {
		t.Errorf("Imported encrypted wallet does not match known answers")
	}
}
