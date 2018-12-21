package goTezos

import (
	"testing"
)

func TestCreateWalletWithMnemonic(t *testing.T) {
	
	gt := NewGoTezos()
	
	// Test Wallet from Alphanet faucet:
	//   Address: tz1Qny7jVMGiwRrP9FikRK95jTNbJcffTpx1
	//   Public Key: edpkvEoAbkdaGALxi2FfeefB8hUkMZ4J1UVwkzyumx2GvbVpkYUHnm
	//   Secret Key: edskRxB2DmoyZSyvhsqaJmw5CK6zYT7dbkUfEVSiQeWU1gw3ZMnC99QMMXru3imsbUrLhvuHktrymvNqhMxkhz7Y4LJAtevW5V

	mnemonic := "normal dash crumble neutral reflect parrot know stairs culture fault check whale flock dog scout"
	password := "PYh8nXDQLB"
	email := "vksbjweo.qsrgfvbw@tezos.example.org"
	
	bakerWallet, err := gt.CreateWallet(mnemonic, email+password)
	if err != nil {
		t.Errorf("Unable to create wallet from Mnemonic: %s", err)
	}
	
	if bakerWallet.Address != "tz1Qny7jVMGiwRrP9FikRK95jTNbJcffTpx1" ||
	   bakerWallet.Pk != "edpkvEoAbkdaGALxi2FfeefB8hUkMZ4J1UVwkzyumx2GvbVpkYUHnm" ||
	   bakerWallet.Sk != "edskRxB2DmoyZSyvhsqaJmw5CK6zYT7dbkUfEVSiQeWU1gw3ZMnC99QMMXru3imsbUrLhvuHktrymvNqhMxkhz7Y4LJAtevW5V" {
		t.Errorf("Created wallet values do not match known answers")
	}
}
	
func TestImportWallet(t *testing.T) {
	
	gt := NewGoTezos()
	
	anotherWallet, err := gt.ImportWallet(
		"tz1fYvVTsSQWkt63P5V8nMjW764cSTrKoQKK",
		"edpkvH3h91QHjKtuR45X9BJRWJJmK7s8rWxiEPnNXmHK67EJYZF75G",
		"edskSA4oADtx6DTT6eXdBc6Pv5MoVBGXUzy8bBryi6D96RQNQYcRfVEXd2nuE2ZZPxs4YLZeM7KazUULFT1SfMDNyKFCUgk6vR")
	if err != nil {
		t.Errorf("%s", err)
	}
	
	if anotherWallet.Address != "tz1fYvVTsSQWkt63P5V8nMjW764cSTrKoQKK" ||
	   anotherWallet.Pk != "edpkvH3h91QHjKtuR45X9BJRWJJmK7s8rWxiEPnNXmHK67EJYZF75G" ||
	   anotherWallet.Sk != "edskSA4oADtx6DTT6eXdBc6Pv5MoVBGXUzy8bBryi6D96RQNQYcRfVEXd2nuE2ZZPxs4YLZeM7KazUULFT1SfMDNyKFCUgk6vR" {
		t.Errorf("Created wallet values do not match known answers")
	}
}
