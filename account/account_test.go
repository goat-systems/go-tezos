package account

import (
	"testing"

	"gotest.tools/assert"

	tzc "github.com/DefinitelyNotAGoat/go-tezos/v2/client"
)

func Test_CreateWalletWithMnemonic(t *testing.T) {
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

	for _, c := range cases {
		accountService := NewAccountService(
			nil,
			nil,
			nil,
		)
		myWallet, err := accountService.CreateWallet(c.mnemonic, c.email+c.password)
		assert.NilError(t, err)
		assert.Assert(t, myWallet.Address == c.result.Address || myWallet.Pk == c.result.Pk || myWallet.Sk == c.result.Sk)
	}
}

func Test_ImportWalletFullSk(t *testing.T) {
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

	for _, c := range cases {
		accountService := NewAccountService(
			nil,
			nil,
			nil,
		)
		myWallet, err := accountService.ImportWallet(c.pkh, c.pk, c.sk)
		assert.NilError(t, err)
		assert.Assert(t, myWallet.Address == c.pkh || myWallet.Pk == c.pk || myWallet.Sk == c.sk)
	}
}

func Test_ImportWalletSeedSk(t *testing.T) {
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

	for _, c := range cases {
		accountService := NewAccountService(
			nil,
			nil,
			nil,
		)
		myWallet, err := accountService.ImportWallet(c.pkh, c.pk, c.sks)
		assert.NilError(t, err)
		assert.Assert(t, myWallet.Address == c.pkh || myWallet.Pk == c.pk || myWallet.Sk == c.sk)
	}
}

func Test_ImportEncryptedSecret(t *testing.T) {
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

	for _, c := range cases {
		accountService := NewAccountService(
			nil,
			nil,
			nil,
		)
		myWallet, err := accountService.ImportEncryptedWallet(c.pw, c.sk)
		assert.NilError(t, err)
		assert.Assert(t, myWallet.Address == c.pkh || myWallet.Pk == c.pk)
	}
}

func Test_GetBalanceAtSnapshot(t *testing.T) {
	var cases = []struct {
		address  string
		cycle    int
		tzclient tzc.TezosClient
		want     float64
	}{
		{
			address: "tz1U8sXoQWGUMQrfZeAYwAzMZUvWwy7mfpPQ",
			cycle:   9,
			tzclient: &clientMock{
				ReturnBody: []byte(`"450209832"`),
			},
			want: 450.209832,
		},
	}

	for _, tc := range cases {
		accountService := NewAccountService(tc.tzclient, &blockServiceMock{}, &snapshotServiceMock{})
		bal, err := accountService.GetBalanceAtSnapshot(tc.address, tc.cycle)
		assert.NilError(t, err)
		assert.Equal(t, bal, tc.want)
	}
}

func Test_GetBalance(t *testing.T) {
	var cases = []struct {
		address  string
		block    int
		tzclient tzc.TezosClient
		want     float64
	}{
		{
			address: "tz1U8sXoQWGUMQrfZeAYwAzMZUvWwy7mfpPQ",
			block:   100000,
			tzclient: &clientMock{
				ReturnBody: []byte(`"450209832"`),
			},
			want: 450.209832,
		},
	}

	for _, tc := range cases {
		accountService := NewAccountService(tc.tzclient, &blockServiceMock{}, &snapshotServiceMock{})
		bal, err := accountService.GetBalanceAtBlock(tc.address, tc.block)
		assert.NilError(t, err)
		assert.Equal(t, bal, tc.want)
	}
}
