package gotezos

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
	val, err := ForgeScriptExpressionForAddress(`050a000000160000b2e19a9e74440d86c59f13dab8a18ff873e889ea`)
	checkErr(t, false, "", err)
	assert.Equal(t, ScriptExpression("exprv6UsC1sN3Fk2XfgcJCL8NCerP5rCGy1PRESZAqr7L2JdzX55EN"), val)
}

func Test_pack(t *testing.T) {
	str, err := pack(`tz1bwsEWCwSEXdRvnJxvegQZKeX5dj6oKEys`)
	checkErr(t, false, "", err)
	assert.Equal(t, "050a000000160000b2e19a9e74440d86c59f13dab8a18ff873e889ea", str)

	_, err = pack(`junk_j4urjofpr`)
	checkErr(t, true, "", err)
}
