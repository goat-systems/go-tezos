package gotezos

// import (
// 	"testing"

// 	"gotest.tools/assert"

// 	tzc "github.com/DefinitelyNotAGoat/go-tezos/v2/client"
// )

// func Test_CreateWallet(t *testing.T) {
// 	var cases = []struct {
// 		name     string
// 		mnemonic string
// 		password string
// 		email    string
// 		want     Wallet
// 	}{
// 		{
// 			name:     "creates a new wallet with mnemonic",
// 			mnemonic: "normal dash crumble neutral reflect parrot know stairs culture fault check whale flock dog scout",
// 			password: "PYh8nXDQLB",
// 			email:    "vksbjweo.qsrgfvbw@tezos.example.org",
// 			want: Wallet{
// 				Sk:      "edskRxB2DmoyZSyvhsqaJmw5CK6zYT7dbkUfEVSiQeWU1gw3ZMnC99QMMXru3imsbUrLhvuHktrymvNqhMxkhz7Y4LJAtevW5V",
// 				Pk:      "edpkvEoAbkdaGALxi2FfeefB8hUkMZ4J1UVwkzyumx2GvbVpkYUHnm",
// 				Address: "tz1Qny7jVMGiwRrP9FikRK95jTNbJcffTpx1",
// 			},
// 		},
// 	}

// 	for _, tc := range cases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			myWallet, err := CreateWallet(tc.mnemonic, tc.email+tc.password)
// 			assert.NilError(t, err)
// 			assert.Assert(t, myWallet.Address == tc.want.Address || myWallet.Pk == tc.want.Pk || myWallet.Sk == tc.want.Sk)
// 		})
// 	}
// }

// func Test_ImportWalletFullSk(t *testing.T) {
// 	var cases = []struct {
// 		pkh string
// 		pk  string
// 		sk  string
// 	}{
// 		{
// 			"tz1fYvVTsSQWkt63P5V8nMjW764cSTrKoQKK",
// 			"edpkvH3h91QHjKtuR45X9BJRWJJmK7s8rWxiEPnNXmHK67EJYZF75G",
// 			"edskSA4oADtx6DTT6eXdBc6Pv5MoVBGXUzy8bBryi6D96RQNQYcRfVEXd2nuE2ZZPxs4YLZeM7KazUULFT1SfMDNyKFCUgk6vR",
// 		},
// 	}

// 	for _, c := range cases {
// 		myWallet, err := ImportWallet(c.pkh, c.pk, c.sk)
// 		assert.NilError(t, err)
// 		assert.Assert(t, myWallet.Address == c.pkh || myWallet.Pk == c.pk || myWallet.Sk == c.sk)
// 	}
// }

// func Test_ImportWalletSeedSk(t *testing.T) {
// 	var cases = []struct {
// 		pkh string
// 		pk  string
// 		sk  string
// 		sks string
// 	}{
// 		{
// 			"tz1U8sXoQWGUMQrfZeAYwAzMZUvWwy7mfpPQ",
// 			"edpkunwa7a3Y5vDr9eoKy4E21pzonuhqvNjscT9XG27aQV4gXq4dNm",
// 			"edskRjBSseEx9bSRSJJpbypJe5ZXucTtApb6qjechMB1BzEYwcEZyfLooo22Nwk33mPPJ3xZniFoa3o8Js7nNXDdqK9nNjFDi7",
// 			"edsk362Ypv3qLgbnGvZK7JwqNbwiLGe18XhTMFQY4gUonqnaCPiT6X",
// 		},
// 	}

// 	for _, c := range cases {
// 		myWallet, err := ImportWallet(c.pkh, c.pk, c.sks)
// 		assert.NilError(t, err)
// 		assert.Assert(t, myWallet.Address == c.pkh || myWallet.Pk == c.pk || myWallet.Sk == c.sk)
// 	}
// }

// func Test_ImportEncryptedSecret(t *testing.T) {
// 	var cases = []struct {
// 		pw  string
// 		sk  string
// 		pk  string
// 		pkh string
// 	}{
// 		{
// 			"password12345##",
// 			"edesk1fddn27MaLcQVEdZpAYiyGQNm6UjtWiBfNP2ZenTy3CFsoSVJgeHM9pP9cvLJ2r5Xp2quQ5mYexW1LRKee2",
// 			"edpkuHMDkMz46HdRXYwom3xRwqk3zQ5ihWX4j8dwo2R2h8o4gPcbN5",
// 			"tz1L8fUQLuwRuywTZUP5JUw9LL3kJa8LMfoo",
// 		},
// 	}

// 	for _, c := range cases {
// 		myWallet, err := ImportEncryptedWallet(c.pw, c.sk)
// 		assert.NilError(t, err)
// 		assert.Assert(t, myWallet.Address == c.pkh || myWallet.Pk == c.pk)
// 	}
// }

// func Test_GetBalance(t *testing.T) {
// 	var cases = []struct {
// 		address   string
// 		blockhash string
// 		tzclient  tzc.TezosClient
// 		want      float64
// 	}{
// 		{
// 			address:   "tz1U8sXoQWGUMQrfZeAYwAzMZUvWwy7mfpPQ",
// 			blockhash: "somehash",
// 			tzclient: &clientMock{
// 				ReturnBody: []byte(`"450209832"`),
// 			},
// 			want: 450.209832,
// 		},
// 	}

// 	for _, tc := range cases {
// 		gt, err := New("")
// 		assert.NilError(t, err)
// 		bal, err := gt.GetBalance(tc.address, tc.blockhash)
// 		assert.NilError(t, err)
// 		assert.Equal(t, bal, tc.want)
// 	}
// }

// type clientMock struct {
// 	ReturnBody []byte
// }

// func (c *clientMock) Post(path, args string) ([]byte, error) {
// 	return c.ReturnBody, nil
// }

// func (c *clientMock) Get(path string, params map[string]string) ([]byte, error) {
// 	return c.ReturnBody, nil
// }
