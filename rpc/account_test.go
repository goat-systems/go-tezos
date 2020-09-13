package rpc

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Balance(t *testing.T) {
	type input struct {
		hash    string
		address string
		handler http.Handler
	}

	type want struct {
		wantErr     bool
		containsErr string
		balance     int
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
				0,
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
				"failed to get balance",
				0,
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
				1216660108948,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input.handler)
			defer server.Close()

			gt, err := New(server.URL)
			assert.Nil(t, err)

			balance, err := gt.Balance(BalanceInput{
				Address:   tt.input.address,
				Blockhash: tt.input.hash,
			})
			checkErr(t, tt.want.wantErr, tt.want.containsErr, err)
			assert.Equal(t, tt.want.balance, balance)
		})
	}
}

// func Test_SignOperation(t *testing.T) {
// 	type input struct {
// 		operation string
// 	}

// 	type want struct {
// 		wantErr     bool
// 		containsErr string
// 		sigop       SignOperationOutput
// 	}

// 	cases := []struct {
// 		name  string
// 		input input
// 		want  want
// 	}{
// 		{
// 			"is successful",
// 			input{
// 				"a732d3520eeaa3de98d78e5e5cb6c85f72204fd46feb9f76853841d4a701add36c0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e00b960000008ba0cb2fad622697145cf1665124096d25bc31e006c0008ba0cb2fad622697145cf1665124096d25bc31ed3e7bd1008d3bb0300b1a803018b88e99e66c1c2587f87118449f781cb7d44c9c40000a732d3520eeaa3de98d78e5e5cb6c85f72204fd46feb9f76853841d4a701add36c0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e00b960000008ba0cb2fad622697145cf1665124096d25bc31e006c0008ba0cb2fad622697145cf1665124096d25bc31ed3e7bd1008d3bb0300b1a803018b88e99e66c1c2587f87118449f781cb7d44c9c40000",
// 			},
// 			want{
// 				false,
// 				"",
// 				SignOperationOutput{
// 					"a732d3520eeaa3de98d78e5e5cb6c85f72204fd46feb9f76853841d4a701add36c0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e00b960000008ba0cb2fad622697145cf1665124096d25bc31e006c0008ba0cb2fad622697145cf1665124096d25bc31ed3e7bd1008d3bb0300b1a803018b88e99e66c1c2587f87118449f781cb7d44c9c40000a732d3520eeaa3de98d78e5e5cb6c85f72204fd46feb9f76853841d4a701add36c0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e00b960000008ba0cb2fad622697145cf1665124096d25bc31e006c0008ba0cb2fad622697145cf1665124096d25bc31ed3e7bd1008d3bb0300b1a803018b88e99e66c1c2587f87118449f781cb7d44c9c40000cab75865fb8ea042804d4f56d662794cb139564561b61717545a87d55b7b577fa325b875fbfdc999da37a360b50ae3bff129f557a95ec577a51c1aade6a34703",
// 					"cab75865fb8ea042804d4f56d662794cb139564561b61717545a87d55b7b577fa325b875fbfdc999da37a360b50ae3bff129f557a95ec577a51c1aade6a34703",
// 					"edsigtzKyzNAk8breKTYYyatMD1Yn2GhLLNjbpff6YzEjMQ93cwJANpQaRJAYZM5nWPpuKgFYjziULKqx1n7K5ALFY7rZf3k9Px",
// 				},
// 			},
// 		},
// 	}

// 	for _, tt := range cases {
// 		t.Run(tt.name, func(t *testing.T) {
// 			wallet, err := ImportEncryptedWallet("password12345##", "edesk1fddn27MaLcQVEdZpAYiyGQNm6UjtWiBfNP2ZenTy3CFsoSVJgeHM9pP9cvLJ2r5Xp2quQ5mYexW1LRKee2")
// 			assert.Nil(t, err)
// 			sigop, err := wallet.SignOperation(tt.input.operation)
// 			checkErr(t, tt.want.wantErr, tt.want.containsErr, err)
// 			assert.Equal(t, tt.want.sigop, sigop)
// 		})
// 	}
// }
// func Test_edsig(t *testing.T) {
// 	type input struct {
// 		operation string
// 	}

// 	type want struct {
// 		wantErr     bool
// 		containsErr string
// 		sigop       string
// 	}

// 	cases := []struct {
// 		name  string
// 		input input
// 		want  want
// 	}{
// 		{
// 			"is successful",
// 			input{
// 				"a732d3520eeaa3de98d78e5e5cb6c85f72204fd46feb9f76853841d4a701add36c0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e00b960000008ba0cb2fad622697145cf1665124096d25bc31e006c0008ba0cb2fad622697145cf1665124096d25bc31ed3e7bd1008d3bb0300b1a803018b88e99e66c1c2587f87118449f781cb7d44c9c40000a732d3520eeaa3de98d78e5e5cb6c85f72204fd46feb9f76853841d4a701add36c0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e00b960000008ba0cb2fad622697145cf1665124096d25bc31e006c0008ba0cb2fad622697145cf1665124096d25bc31ed3e7bd1008d3bb0300b1a803018b88e99e66c1c2587f87118449f781cb7d44c9c40000",
// 			},
// 			want{
// 				false,
// 				"",
// 				"edsigtzKyzNAk8breKTYYyatMD1Yn2GhLLNjbpff6YzEjMQ93cwJANpQaRJAYZM5nWPpuKgFYjziULKqx1n7K5ALFY7rZf3k9Px",
// 			},
// 		},
// 	}

// 	for _, tt := range cases {
// 		t.Run(tt.name, func(t *testing.T) {
// 			wallet, err := ImportEncryptedWallet("password12345##", "edesk1fddn27MaLcQVEdZpAYiyGQNm6UjtWiBfNP2ZenTy3CFsoSVJgeHM9pP9cvLJ2r5Xp2quQ5mYexW1LRKee2")
// 			assert.Nil(t, err)
// 			sigop, err := wallet.edsig(tt.input.operation)
// 			checkErr(t, tt.want.wantErr, tt.want.containsErr, err)
// 			assert.Equal(t, tt.want.sigop, sigop)
// 		})
// 	}
// }
