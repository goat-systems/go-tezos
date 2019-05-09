package gotezos

import (
	"encoding/json"
	"testing"
)

func TestCreateWalletWithMnemonic(t *testing.T) {
	t.Log("Create new wallet using Alphanet faucet account")

	var cases = []struct {
		mnemonic string
		password string
		email    string
		result   Wallet
	}{
		{"normal dash crumble neutral reflect parrot know stairs culture fault check whale flock dog scout",
			"PYh8nXDQLB",
			"vksbjweo.qsrgfvbw@tezos.example.org",
			Wallet{Sk: "edskRxB2DmoyZSyvhsqaJmw5CK6zYT7dbkUfEVSiQeWU1gw3ZMnC99QMMXru3imsbUrLhvuHktrymvNqhMxkhz7Y4LJAtevW5V",
				Pk:      "edpkvEoAbkdaGALxi2FfeefB8hUkMZ4J1UVwkzyumx2GvbVpkYUHnm",
				Address: "tz1Qny7jVMGiwRrP9FikRK95jTNbJcffTpx1",
			},
		},
	}

	gt, _ := NewGoTezos("http://127.0.0.1:8732")

	for _, c := range cases {
		myWallet, err := gt.Account.CreateWallet(c.mnemonic, c.email+c.password)
		if err != nil {
			t.Errorf("Unable to create wallet from Mnemonic: %s", err)
		}

		if myWallet.Address != c.result.Address || myWallet.Pk != c.result.Pk || myWallet.Sk != c.result.Sk {
			t.Errorf("Created wallet values do not match known answers")
		}
	}
}

func TestImportWalletFullSk(t *testing.T) {

	t.Log("Import existing wallet using complete secret key")

	var cases = []struct {
		pkh string
		pk  string
		sk  string
	}{
		{
			"tz1fYvVTsSQWkt63P5V8nMjW764cSTrKoQKK",
			"edpkvH3h91QHjKtuR45X9BJRWJJmK7s8rWxiEPnNXmHK67EJYZF75G",
			"edskSA4oADtx6DTT6eXdBc6Pv5MoVBGXUzy8bBryi6D96RQNQYcRfVEXd2nuE2ZZPxs4YLZeM7KazUULFT1SfMDNyKFCUgk6vR",
		},
	}

	gt, _ := NewGoTezos("http://127.0.0.1:8732")

	for _, c := range cases {
		myWallet, err := gt.Account.ImportWallet(c.pkh, c.pk, c.sk)
		if err != nil {
			t.Errorf("%s", err)
		}

		if myWallet.Address != c.pkh || myWallet.Pk != c.pk || myWallet.Sk != c.sk {
			t.Errorf("Created wallet values do not match known answers")
		}
	}
}

func TestImportWalletSeedSk(t *testing.T) {

	t.Log("Import existing wallet using seed-secret key")

	var cases = []struct {
		pkh string
		pk  string
		sk  string
		sks string
	}{
		{
			"tz1U8sXoQWGUMQrfZeAYwAzMZUvWwy7mfpPQ",
			"edpkunwa7a3Y5vDr9eoKy4E21pzonuhqvNjscT9XG27aQV4gXq4dNm",
			"edskRjBSseEx9bSRSJJpbypJe5ZXucTtApb6qjechMB1BzEYwcEZyfLooo22Nwk33mPPJ3xZniFoa3o8Js7nNXDdqK9nNjFDi7",
			"edsk362Ypv3qLgbnGvZK7JwqNbwiLGe18XhTMFQY4gUonqnaCPiT6X",
		},
	}

	gt, _ := NewGoTezos("http://127.0.0.1:8732")

	for _, c := range cases {
		myWallet, err := gt.Account.ImportWallet(c.pkh, c.pk, c.sks)
		if err != nil {
			t.Errorf("%s", err)
		}

		if myWallet.Address != c.pkh || myWallet.Pk != c.pk || myWallet.Sk != c.sk {
			t.Errorf("Created wallet values do not match known answers")
		}
	}
}

func TestImportEncryptedSecret(t *testing.T) {

	t.Log("Import wallet using password and encrypted key")

	var cases = []struct {
		pw  string
		sk  string
		pk  string
		pkh string
	}{
		{
			"password12345##",
			"edesk1fddn27MaLcQVEdZpAYiyGQNm6UjtWiBfNP2ZenTy3CFsoSVJgeHM9pP9cvLJ2r5Xp2quQ5mYexW1LRKee2",
			"edpkuHMDkMz46HdRXYwom3xRwqk3zQ5ihWX4j8dwo2R2h8o4gPcbN5",
			"tz1L8fUQLuwRuywTZUP5JUw9LL3kJa8LMfoo",
		},
	}

	gt, _ := NewGoTezos("http://127.0.0.1:8732")

	for _, c := range cases {

		myWallet, err := gt.Account.ImportEncryptedWallet(c.pw, c.sk)
		if err != nil {
			t.Errorf("%s", err)
		}

		if myWallet.Address != c.pkh || myWallet.Pk != c.pk {
			t.Errorf("Imported encrypted wallet does not match known answers")
		}
	}
}
func TestGetSnapShot(t *testing.T) {
	var cases = []struct {
		in  int
		out SnapShot
	}{
		{171, SnapShot{Cycle: 171, AssociatedBlock: 340992, AssociatedHash: "BLV6XGmLgvkNi7BgCCbLjD3mdQ3LaZQCLcZa6aFzPsDuu3ySQvU"}},
		{132, SnapShot{Cycle: 132, AssociatedBlock: 261376, AssociatedHash: "BMUaWotqn6icj8Wk1ERJJLGVvdLMc75fUrPkG7dLhAMYFcWBYfe"}},
	}

	gt, err := NewGoTezos("http://127.0.0.1:8732")
	if err != nil {
		t.Errorf("could not connect to network")
	}

	for _, c := range cases {
		snapshot, err := gt.SnapShot.Get(c.in)
		if err != nil {
			t.Error(err)
		}
		if c.out.AssociatedBlock != snapshot.AssociatedBlock || c.out.AssociatedHash != snapshot.AssociatedHash || c.out.Cycle != snapshot.Cycle {
			t.Errorf("Snap Shot %v, does not match the snapshot queryied: %v", c.out, snapshot)
		}
	}
}

func TestBlockGetHead(t *testing.T) {
	gt, err := NewGoTezos("http://127.0.0.1:8732")
	if err != nil {
		t.Errorf("could not connect to network: %v", err)
	}

	_, err = gt.Block.GetHead()
	if err != nil {
		t.Errorf("%v", err)
	}
}

func TestNetworkGetConstants(t *testing.T) {
	gt, err := NewGoTezos("http://127.0.0.1:8732")
	if err != nil {
		t.Errorf("could not connect to network")
	}

	_, err = gt.Network.GetConstants()
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestGetNetworkVersions(t *testing.T) {
	gt, err := NewGoTezos("http://127.0.0.1:8732")
	if err != nil {
		t.Errorf("could not connect to network")
	}

	_, err = gt.Network.GetVersions()
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestGetBlockOperationHashesHead(t *testing.T) {
	gt, err := NewGoTezos("http://127.0.0.1:8732")
	if err != nil {
		t.Errorf("could not connect to network")
	}

	_, err = gt.Operation.GetBlockOperationHashes(100000)
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestGetBlockOperationHashes(t *testing.T) {
	gt, err := NewGoTezos("http://127.0.0.1:8732")
	if err != nil {
		t.Errorf("could not connect to network")
	}

	ops, err := gt.Operation.GetBlockOperationHashes("BMWQkPagYkqzT5sj7tjuBwwDgJRLfKBLBbdhVqqJhmgwiRgQBuk")
	if err != nil {
		t.Errorf("%s", err)
	}

	if len(ops) != 10 {
		t.Errorf("%d", len(ops))
	}

}

func TestGetAccountBalanceAtSnapshot(t *testing.T) {
	gt, err := NewGoTezos("http://127.0.0.1:8732")
	if err != nil {
		t.Errorf("could not connect to network")
	}

	_, err = gt.Account.GetBalanceAtSnapshot("tz3gN8NTLNLJg5KRsUU47NHNVHbdhcFXjjaB", 15)
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestGetAccountBalance(t *testing.T) {
	gt, err := NewGoTezos("http://127.0.0.1:8732")
	if err != nil {
		t.Errorf("could not connect to network")
	}

	_, err = gt.Account.GetBalance("tz3gN8NTLNLJg5KRsUU47NHNVHbdhcFXjjaB")
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestDelegateGetStakingBalance(t *testing.T) {
	gt, err := NewGoTezos("http://127.0.0.1:8732")
	if err != nil {
		t.Errorf("could not connect to network")
	}

	_, err = gt.Delegate.GetStakingBalance("tz3gN8NTLNLJg5KRsUU47NHNVHbdhcFXjjaB", 15)
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestCycleGetCurrent(t *testing.T) {
	gt, err := NewGoTezos("http://127.0.0.1:8732")
	if err != nil {
		t.Errorf("could not connect to network")
	}

	_, err = gt.Cycle.GetCurrent()
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestAccountGetBalanceAtBlock(t *testing.T) {
	gt, err := NewGoTezos("http://127.0.0.1:8732")
	if err != nil {
		t.Errorf("could not connect to network")
	}

	_, err = gt.Account.GetBalanceAtBlock("tz3gN8NTLNLJg5KRsUU47NHNVHbdhcFXjjaB", 100000)
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestDelegateGetDelegations(t *testing.T) {
	gt, err := NewGoTezos("http://127.0.0.1:8732")
	if err != nil {
		t.Errorf("could not connect to network")
	}

	_, err = gt.Delegate.GetDelegations("tz3gN8NTLNLJg5KRsUU47NHNVHbdhcFXjjaB")
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestGetDelegationsAtCycle(t *testing.T) {
	gt, err := NewGoTezos("http://127.0.0.1:8732")
	if err != nil {
		t.Errorf("could not connect to network")
	}

	_, err = gt.Delegate.GetDelegationsAtCycle("tz3gN8NTLNLJg5KRsUU47NHNVHbdhcFXjjaB", 12)
	if err != nil {
		t.Errorf("%s", err)
	}

}

func TestDelegateGetReport(t *testing.T) {
	gt, err := NewGoTezos("http://127.0.0.1:8732")
	if err != nil {
		t.Errorf("could not connect to network")
	}

	report, err := gt.Delegate.GetReport("tz1T8UYSbVuRm6CdhjvwCfXsKXb4yL9ai9Q3", 172, 0.05)
	if err != nil {
		t.Errorf("%s", err)
	}

	t.Log(PrettyReport(report))

}

func TestGetCycleRewards(t *testing.T) {
	gt, err := NewGoTezos("http://127.0.0.1:8732")
	if err != nil {
		t.Errorf("could not connect to network")
	}

	_, err = gt.Delegate.GetRewards("tz3gN8NTLNLJg5KRsUU47NHNVHbdhcFXjjaB", 60)
	if err != nil {
		t.Errorf("%s", err)
	}

}

func TestGetDelegate(t *testing.T) {
	gt, err := NewGoTezos("http://127.0.0.1:8732")
	if err != nil {
		t.Errorf("could not connect to network")
	}

	_, err = gt.Delegate.GetDelegate("tz3gN8NTLNLJg5KRsUU47NHNVHbdhcFXjjaB")
	if err != nil {
		t.Errorf("%s", err)
	}

}

func TestGetBakingRights(t *testing.T) {
	gt, err := NewGoTezos("http://127.0.0.1:8732")
	if err != nil {
		t.Errorf("could not connect to network")
	}

	_, err = gt.Delegate.GetBakingRights(60)
	if err != nil {
		t.Errorf("%s", err)
	}

}

func TestGetBakingRightsForDelegate(t *testing.T) {
	gt, err := NewGoTezos("http://127.0.0.1:8732")
	if err != nil {
		t.Errorf("could not connect to network")
	}

	_, err = gt.Delegate.GetBakingRightsForDelegate(60, "tz3gN8NTLNLJg5KRsUU47NHNVHbdhcFXjjaB", 2)
	if err != nil {
		t.Errorf("%s", err)
	}

}

func TestGetEndorsingRightsForDelegate(t *testing.T) {
	gt, err := NewGoTezos("http://127.0.0.1:8732")
	if err != nil {
		t.Errorf("could not connect to network")
	}

	_, err = gt.Delegate.GetEndorsingRightsForDelegate(60, "tz3gN8NTLNLJg5KRsUU47NHNVHbdhcFXjjaB")
	if err != nil {
		t.Errorf("%s", err)
	}

}

func TestGetEndorsingRights(t *testing.T) {
	gt, err := NewGoTezos("http://127.0.0.1:8732")
	if err != nil {
		t.Errorf("could not connect to network")
	}

	_, err = gt.Delegate.GetEndorsingRights(60)
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestGetAllDelegatesByHash(t *testing.T) {
	gt, err := NewGoTezos("http://127.0.0.1:8732")
	if err != nil {
		t.Errorf("could not connect to network")
	}

	_, err = gt.Delegate.GetAllDelegatesByHash("BMWQkPagYkqzT5sj7tjuBwwDgJRLfKBLBbdhVqqJhmgwiRgQBuk")
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestGetAllDelegates(t *testing.T) {
	gt, err := NewGoTezos("http://127.0.0.1:8732")
	if err != nil {
		t.Errorf("could not connect to network")
	}

	_, err = gt.Delegate.GetAllDelegates()
	if err != nil {
		t.Errorf("%s", err)
	}

}

//Takes an interface v and returns a pretty json string.
func PrettyReport(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		return string(b)
	}
	return ""
}
