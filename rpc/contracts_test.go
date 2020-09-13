package rpc

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ContractStorage(t *testing.T) {
	storageJSON := []byte(`[
		{
		  "prim": "parameter",
		  "args": [
			{
			  "prim": "unit",
			  "annots": [
				"%abc"
			  ]
			}
		  ]
		},
		{
		  "prim": "storage",
		  "args": [
			{
			  "prim": "unit"
			}
		  ]
		},
		{
		  "prim": "code",
		  "args": [
			[
			  {
				"prim": "CDR"
			  },
			  {
				"prim": "NIL",
				"args": [
				  {
					"prim": "operation"
				  }
				]
			  },
			  {
				"prim": "PAIR"
			  }
			]
		  ]
		}
	  ]`)

	type want struct {
		err         bool
		containsErr string
		micheline   []byte
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"returns rpc error",
			gtGoldenHTTPMock(storageHandlerMock(readResponse(rpcerrors), blankHandler)),
			want{
				true,
				"could not get storage",
				[]byte{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(storageHandlerMock(storageJSON, blankHandler)),
			want{
				false,
				"",
				storageJSON,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			gt, err := New(server.URL)
			assert.Nil(t, err)

			micheline, err := gt.ContractStorage("BLzGD63HA4RP8Fh5xEtvdQSMKa2WzJMZjQPNVUc4Rqy8Lh5BEY1", "KT1LfoE9EbpdsfUzowRckGUfikGcd5PyVKg")
			checkErr(t, tt.want.err, tt.containsErr, err)
			assert.Equal(t, tt.want.micheline, micheline)
		})
	}
}

func Test_ForgeScriptExpressionForAddress(t *testing.T) {
	val, err := ForgeScriptExpressionForAddress(`tz1S82rGFZK8cVbNDpP1Hf9VhTUa4W8oc2WV`)
	checkErr(t, false, "", err)
	assert.Equal(t, ScriptExpression("expru1LH1CafV3yYgs9BkbrMWWfAE9ye3RdWwyndr9MKYN8w5VQ7Rt"), val)
}
