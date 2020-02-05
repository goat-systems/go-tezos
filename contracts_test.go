package gotezos

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ContractStorage(t *testing.T) {
	type want struct {
		err         bool
		containsErr string
		rpcerr      []byte
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"returns rpc error",
			gtGoldenHTTPMock(storageHandlerMock(mockRPCErrorResp, blankHandler)),
			want{
				true,
				"could not get storage",
				mockRPCErrorResp,
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(storageHandlerMock([]byte(`"Hello Tezos!"`), blankHandler)),
			want{
				false,
				"",
				[]byte(`"Hello Tezos!"`),
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			gt, err := New(server.URL)
			assert.Nil(t, err)

			rpcerr, err := gt.ContractStorage("BLzGD63HA4RP8Fh5xEtvdQSMKa2WzJMZjQPNVUc4Rqy8Lh5BEY1", "KT1LfoE9EbpdsfUzowRckGUfikGcd5PyVKg")
			if tt.want.err {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tt.want.containsErr)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, tt.want.rpcerr, rpcerr)
		})
	}
}
