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
		balance     string
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
				"0",
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
				"0",
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
				"1216660108948",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input.handler)
			defer server.Close()

			rpc, err := New(server.URL)
			assert.Nil(t, err)

			balance, err := rpc.Balance(BalanceInput{
				Address:   tt.input.address,
				Blockhash: tt.input.hash,
			})
			checkErr(t, tt.want.wantErr, tt.want.containsErr, err)
			assert.Equal(t, tt.want.balance, balance)
		})
	}
}
