package gotezos

import (
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CreateWallet(t *testing.T) {
	type input struct {
		mnemonic string
		password string
		email    string
	}

	type want struct {
		err         bool
		errContains string
		wallet      *Wallet
	}

	var cases = []struct {
		name  string
		input input
		want  want
	}{
		{
			"creates a new wallet with mnemonic",
			input{
				"normal dash crumble neutral reflect parrot know stairs culture fault check whale flock dog scout",
				"PYh8nXDQLB",
				"vksbjweo.qsrgfvbw@tezos.example.org",
			},
			want{
				false,
				"",
				&Wallet{
					Sk:      "edskRxB2DmoyZSyvhsqaJmw5CK6zYT7dbkUfEVSiQeWU1gw3ZMnC99QMMXru3imsbUrLhvuHktrymvNqhMxkhz7Y4LJAtevW5V",
					Pk:      "edpkvEoAbkdaGALxi2FfeefB8hUkMZ4J1UVwkzyumx2GvbVpkYUHnm",
					Address: "tz1Qny7jVMGiwRrP9FikRK95jTNbJcffTpx1",
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			wallet, err := CreateWallet(tt.input.mnemonic, tt.input.email+tt.input.password)
			if tt.want.err {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tt.want.errContains)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, tt.want.wallet.Sk, wallet.Sk)
			assert.Equal(t, tt.want.wallet.Pk, wallet.Pk)
			assert.Equal(t, tt.want.wallet.Address, wallet.Address)
		})
	}
}

func Test_ImportWallet(t *testing.T) {
	type input struct {
		pkh    string
		pk     string
		sk     string
		sks    string
		useSKS bool
	}

	type want struct {
		err         bool
		errContains string
		wallet      *Wallet
	}

	var cases = []struct {
		name  string
		input input
		want  want
	}{
		{
			"is successful with full sk",
			input{
				"tz1fYvVTsSQWkt63P5V8nMjW764cSTrKoQKK",
				"edpkvH3h91QHjKtuR45X9BJRWJJmK7s8rWxiEPnNXmHK67EJYZF75G",
				"edskSA4oADtx6DTT6eXdBc6Pv5MoVBGXUzy8bBryi6D96RQNQYcRfVEXd2nuE2ZZPxs4YLZeM7KazUULFT1SfMDNyKFCUgk6vR",
				"",
				false,
			},
			want{
				false,
				"",
				&Wallet{
					Address: "tz1fYvVTsSQWkt63P5V8nMjW764cSTrKoQKK",
					Pk:      "edpkvH3h91QHjKtuR45X9BJRWJJmK7s8rWxiEPnNXmHK67EJYZF75G",
					Sk:      "edskSA4oADtx6DTT6eXdBc6Pv5MoVBGXUzy8bBryi6D96RQNQYcRfVEXd2nuE2ZZPxs4YLZeM7KazUULFT1SfMDNyKFCUgk6vR",
				},
			},
		},
		{
			"is successful with seed sk (sks)",
			input{
				"tz1U8sXoQWGUMQrfZeAYwAzMZUvWwy7mfpPQ",
				"edpkunwa7a3Y5vDr9eoKy4E21pzonuhqvNjscT9XG27aQV4gXq4dNm",
				"",
				"edsk362Ypv3qLgbnGvZK7JwqNbwiLGe18XhTMFQY4gUonqnaCPiT6X",
				true,
			},
			want{
				false,
				"",
				&Wallet{
					Address: "tz1U8sXoQWGUMQrfZeAYwAzMZUvWwy7mfpPQ",
					Pk:      "edpkunwa7a3Y5vDr9eoKy4E21pzonuhqvNjscT9XG27aQV4gXq4dNm",
					Sk:      "edskRjBSseEx9bSRSJJpbypJe5ZXucTtApb6qjechMB1BzEYwcEZyfLooo22Nwk33mPPJ3xZniFoa3o8Js7nNXDdqK9nNjFDi7",
				},
			},
		},
		{
			"is missing pkh",
			input{
				"",
				"edpkvH3h91QHjKtuR45X9BJRWJJmK7s8rWxiEPnNXmHK67EJYZF75G",
				"edskSA4oADtx6DTT6eXdBc6Pv5MoVBGXUzy8bBryi6D96RQNQYcRfVEXd2nuE2ZZPxs4YLZeM7KazUULFT1SfMDNyKFCUgk6vR",
				"",
				false,
			},
			want{
				true,
				"reconstructed address 'tz1fYvVTsSQWkt63P5V8nMjW764cSTrKoQKK' does not match provided address",
				&Wallet{
					Sk: "edskSA4oADtx6DTT6eXdBc6Pv5MoVBGXUzy8bBryi6D96RQNQYcRfVEXd2nuE2ZZPxs4YLZeM7KazUULFT1SfMDNyKFCUgk6vR",
				},
			},
		},
		{
			"is missing pk",
			input{
				"tz1fYvVTsSQWkt63P5V8nMjW764cSTrKoQKK",
				"",
				"edskSA4oADtx6DTT6eXdBc6Pv5MoVBGXUzy8bBryi6D96RQNQYcRfVEXd2nuE2ZZPxs4YLZeM7KazUULFT1SfMDNyKFCUgk6vR",
				"",
				false,
			},
			want{
				true,
				"reconstructed pk 'edpkvH3h91QHjKtuR45X9BJRWJJmK7s8rWxiEPnNXmHK67EJYZF75G' does not match provided pk ''",
				&Wallet{
					Address: "tz1fYvVTsSQWkt63P5V8nMjW764cSTrKoQKK",
					Sk:      "edskSA4oADtx6DTT6eXdBc6Pv5MoVBGXUzy8bBryi6D96RQNQYcRfVEXd2nuE2ZZPxs4YLZeM7KazUULFT1SfMDNyKFCUgk6vR",
				},
			},
		},
		{
			"is missing sk",
			input{
				"tz1fYvVTsSQWkt63P5V8nMjW764cSTrKoQKK",
				"edpkvH3h91QHjKtuR45X9BJRWJJmK7s8rWxiEPnNXmHK67EJYZF75G",
				"",
				"",
				false,
			},
			want{
				true,
				"wallet prefix is not edsk",
				&Wallet{},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			var sk string
			if tt.input.useSKS {
				sk = tt.input.sks
			} else {
				sk = tt.input.sk
			}

			wallet, err := ImportWallet(tt.input.pkh, tt.input.pk, sk)
			if tt.want.err {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tt.want.errContains)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, tt.want.wallet.Address, wallet.Address)
			assert.Equal(t, tt.want.wallet.Pk, wallet.Pk)
			assert.Equal(t, tt.want.wallet.Sk, wallet.Sk)
		})
	}
}

func Test_ImportEncryptedSecret(t *testing.T) {
	type input struct {
		pw string
		sk string
	}

	type want struct {
		err         bool
		errContains string
		wallet      *Wallet
	}

	var cases = []struct {
		name  string
		input input
		want  want
	}{
		{
			"is successful",
			input{
				"password12345##",
				"edesk1fddn27MaLcQVEdZpAYiyGQNm6UjtWiBfNP2ZenTy3CFsoSVJgeHM9pP9cvLJ2r5Xp2quQ5mYexW1LRKee2",
			},
			want{
				false,
				"",
				&Wallet{
					Address: "tz1L8fUQLuwRuywTZUP5JUw9LL3kJa8LMfoo",
					Pk:      "edpkuHMDkMz46HdRXYwom3xRwqk3zQ5ihWX4j8dwo2R2h8o4gPcbN5",
					Sk:      "edskRsPBsKuULoLTEQV2R9UbvSZbzFqvoESvp1mYyQJU8xi9mJamt88r5uTXbWQpVHjSiPWWtnoyqTCuSLQLxbEKUXfwwTccsF",
				},
			},
		},
		{
			"is missing password",
			input{
				"",
				"edesk1fddn27MaLcQVEdZpAYiyGQNm6UjtWiBfNP2ZenTy3CFsoSVJgeHM9pP9cvLJ2r5Xp2quQ5mYexW1LRKee2",
			},
			want{
				true,
				"invalid password",
				&Wallet{},
			},
		},
		{
			"is missing esk",
			input{
				"password12345##",
				"",
			},
			want{
				true,
				"encrypted secret key does not 88 characters long",
				&Wallet{},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			wallet, err := ImportEncryptedWallet(tt.input.pw, tt.input.sk)
			if tt.want.err {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tt.want.errContains)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, tt.want.wallet.Address, wallet.Address)
			assert.Equal(t, tt.want.wallet.Pk, wallet.Pk)
			assert.Equal(t, tt.want.wallet.Sk, wallet.Sk)
		})
	}
}

func Test_Balance(t *testing.T) {
	goldenBalance := getResponse(balance).(*Int)

	type input struct {
		hash    string
		address string
		handler http.Handler
	}

	type want struct {
		wantErr     bool
		containsErr string
		balance     *big.Int
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"returns rpc error",
			input{
				mockBlockHash,
				"tz1U8sXoQWGUMQrfZeAYwAzMZUvWwy7mfpPQ",
				gtGoldenHTTPMock(balanceHandlerMock(readResponse(rpcerrors), blankHandler)),
			},
			want{
				true,
				"failed to get balance",
				big.NewInt(0),
			},
		},
		{
			"failed to unmarshal",
			input{
				mockBlockHash,
				"tz1U8sXoQWGUMQrfZeAYwAzMZUvWwy7mfpPQ",
				gtGoldenHTTPMock(balanceHandlerMock([]byte(`not_balance_data`), blankHandler)),
			},
			want{
				true,
				"failed to unmarshal balance",
				big.NewInt(0),
			},
		},
		{
			"is successful",
			input{
				mockBlockHash,
				"tz1U8sXoQWGUMQrfZeAYwAzMZUvWwy7mfpPQ",
				gtGoldenHTTPMock(balanceHandlerMock(readResponse(balance), blankHandler)),
			},
			want{
				false,
				"",
				goldenBalance.Big,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input.handler)
			defer server.Close()

			gt, err := New(server.URL)
			assert.Nil(t, err)

			balance, err := gt.Balance(tt.input.hash, tt.input.address)
			checkErr(t, tt.want.wantErr, tt.want.containsErr, err)
			assert.Equal(t, tt.want.balance, balance)
		})
	}
}
