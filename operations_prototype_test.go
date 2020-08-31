package gotezos

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Forge_Origination(t *testing.T) {
	originationJSON := []byte(`{
		"kind": "origination",
		"source": "tz1ZmJch5fHBfgXf2YmGhvFEArH6my4JQUZd",
		"fee": "10000",
		"counter": "394934",
		"gas_limit": "10000",
		"storage_limit": "10000",
		"balance": "10000",
		"script": {
		  "code": [
			{
			  "prim": "parameter",
			  "args": [
				{
				  "prim": "string"
				}
			  ]
			},
			{
			  "prim": "storage",
			  "args": [
				{
				  "prim": "string"
				}
			  ]
			},
			{
			  "prim": "code",
			  "args": [
				[
				  {
					"prim": "CAR"
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
		  ],
		  "storage": {
			"string": "Test"
		  }
		}
	  }`)

	var origination Origination
	err := json.Unmarshal(originationJSON, &origination)
	checkErr(t, false, "", err)

	type want struct {
		err         bool
		errContains string
		operation   string
	}

	cases := []struct {
		name  string
		input Origination
		want  want
	}{
		{
			"is successful",
			origination,
			want{
				false,
				"",
				"",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			origination, err := tt.input.Forge_Prototype()
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.operation, string(origination))
		})
	}
}
