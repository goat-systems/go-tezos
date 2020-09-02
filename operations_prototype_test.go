package gotezos

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ForgeOperation(t *testing.T) {
	transactionJSON := []byte(`{
		"kind": "transaction",
		"source": "tz1XJ1UNechmHKhQo4tvVX6qztnVuQuSFKgd",
		"fee": "1283",
		"counter": "7",
		"gas_limit": "10307",
		"storage_limit": "0",
		"amount": "20000000000",
		"destination": "tz1aWXP237BLwNHJcCD4b3DutCevhqq2T1Z9"
	  }`)

	var transaction Transaction
	err := json.Unmarshal(transactionJSON, &transaction)
	checkErr(t, false, "", err)

	transactionJSON2 := []byte(`{
		"kind": "transaction",
		"source": "tz1SJJY253HoEda8PS5vvfHVtyghgK3CTS2z",
		"fee": "2966",
		"counter": "133558",
		"gas_limit": "26271",
		"storage_limit": "0",
		"amount": "0",
		"destination": "KT1XdCkJncWfGvqf1NdbK2HBRTvRcHhJtNx5",
		"parameters": {
		  "entrypoint": "do",
		  "value": [
			{
			  "prim": "RENAME"
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
			  "prim": "PUSH",
			  "args": [
				{
				  "prim": "key_hash"
				},
				{
				  "string": "tz2L2HuhaaSnf6ShEDdhTEAr5jGPWPNwpvcB"
				}
			  ]
			},
			{
			  "prim": "IMPLICIT_ACCOUNT"
			},
			{
			  "prim": "PUSH",
			  "args": [
				{
				  "prim": "mutez"
				},
				{
				  "int": "2"
				}
			  ]
			},
			{
			  "prim": "UNIT"
			},
			{
			  "prim": "TRANSFER_TOKENS"
			},
			{
			  "prim": "CONS"
			},
			{
			  "prim": "DIP",
			  "args": [
				[
				  {
					"prim": "DROP"
				  }
				]
			  ]
			}
		  ]
		}
	  }`)

	var transaction2 Transaction
	err = json.Unmarshal(transactionJSON2, &transaction2)
	checkErr(t, false, "", err)

	type input struct {
		branch   string
		contents []OperationContents
	}

	type want struct {
		err         bool
		errContains string
		operation   string
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"is successful with transaction",
			input{
				branch: "BLQMkH2PSTuAJgVm6rGHshY5z6Z6SAmqXv6q1LDzhX6fchJ12Up",
				contents: []OperationContents{
					&transaction,
				},
			},
			want{
				false,
				"",
				"5aff622d53d32a8bae591627718c60a35b16737e301c57a13b6f1765483d88ff6c007fd82c06cf5a203f18faaf562447ed1efcc6c010830a07c350008090dfc04a0000a31e81ac3425310e3274a4698a793b2839dc0afa00",
			},
		},
		{
			"is successful with transaction 2",
			input{
				branch: "BLEkC1TqtP7DJjGnyxwhT8VDnEF75aNMMKS5qJXSTFmAKkV7Pch",
				contents: []OperationContents{
					&transaction2,
				},
			},
			want{
				false,
				"",
				"452b8599b0e4960b884d3ad61c89c594bc3348c798842651d3fa6cfafa77ce556c00490dc9520ec45270f240a3cc4f07aec76adc358d9617b693089fcd01000001fcc0bee1480bfca3a80481904cee4099400b1c8d00ff020000004f020000004a0358053d036d0743035d0100000024747a324c324875686161536e663653684544646854454172356a475057504e7770766342031e0743036a0002034f034d031b051f02000000020320",
			},
		},
		{
			"is successful with multiple transactions",
			input{
				branch: "BLEkC1TqtP7DJjGnyxwhT8VDnEF75aNMMKS5qJXSTFmAKkV7Pch",
				contents: []OperationContents{
					&transaction,
					&transaction2,
				},
			},
			want{
				false,
				"",
				"452b8599b0e4960b884d3ad61c89c594bc3348c798842651d3fa6cfafa77ce556c007fd82c06cf5a203f18faaf562447ed1efcc6c010830a07c350008090dfc04a0000a31e81ac3425310e3274a4698a793b2839dc0afa006c00490dc9520ec45270f240a3cc4f07aec76adc358d9617b693089fcd01000001fcc0bee1480bfca3a80481904cee4099400b1c8d00ff020000004f020000004a0358053d036d0743035d0100000024747a324c324875686161536e663653684544646854454172356a475057504e7770766342031e0743036a0002034f034d031b051f02000000020320",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			operation, err := ForgeOperation(tt.input.branch, tt.input.contents...)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.operation, operation)
		})
	}
}

func Test_Forge_Origination(t *testing.T) {
	originationJSON := []byte(`{
		"kind": "origination",
		"source": "tz1TJCwoX79reCZ8yccPeW8iB9Mba91v8H47",
		"fee": "1389",
		"counter": "307028",
		"gas_limit": "11140",
		"storage_limit": "323",
		"balance": "0",
		"script": {
		  "code": [
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
		  ],
		  "storage": {
			"prim": "Unit"
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
				"6d0054013ef6636fe99989a26006622bf270be0b1485ed0ad4de128457c302000000000024020000001f0500046c00000004256162630501036c050202000000080317053d036d034200000002030b",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			origination, err := tt.input.forge()
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.operation, hex.EncodeToString(origination))
		})
	}
}

func Test_Forge_Transaction(t *testing.T) {
	transactionJSON := []byte(`{
			"kind": "transaction",
			"source": "tz1NXjqkurAmpKJEF76T58oyNsy3hWK7mk8e",
			"fee": "22100",
			"counter": "377727",
			"gas_limit": "218465",
			"storage_limit": "668",
			"amount": "0",
			"destination": "KT1SkmB19o8nfhRvG9LL7TjDfX2Bm1nCuYoY"
		  }`)

	var transaction Transaction
	err := json.Unmarshal(transactionJSON, &transaction)
	checkErr(t, false, "", err)

	transactionJSON2 := []byte(`{
		"kind": "transaction",
		"source": "tz1SJJY253HoEda8PS5vvfHVtyghgK3CTS2z",
		"fee": "2966",
		"counter": "133558",
		"gas_limit": "26271",
		"storage_limit": "0",
		"amount": "0",
		"destination": "KT1XdCkJncWfGvqf1NdbK2HBRTvRcHhJtNx5",
		"parameters": {
		  "entrypoint": "do",
		  "value": [
			{
			  "prim": "RENAME"
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
			  "prim": "PUSH",
			  "args": [
				{
				  "prim": "key_hash"
				},
				{
				  "string": "tz2L2HuhaaSnf6ShEDdhTEAr5jGPWPNwpvcB"
				}
			  ]
			},
			{
			  "prim": "IMPLICIT_ACCOUNT"
			},
			{
			  "prim": "PUSH",
			  "args": [
				{
				  "prim": "mutez"
				},
				{
				  "int": "2"
				}
			  ]
			},
			{
			  "prim": "UNIT"
			},
			{
			  "prim": "TRANSFER_TOKENS"
			},
			{
			  "prim": "CONS"
			},
			{
			  "prim": "DIP",
			  "args": [
				[
				  {
					"prim": "DROP"
				  }
				]
			  ]
			}
		  ]
		}
	  }`)

	var transaction2 Transaction
	err = json.Unmarshal(transactionJSON2, &transaction2)
	checkErr(t, false, "", err)

	type want struct {
		err         bool
		errContains string
		operation   string
	}

	cases := []struct {
		name  string
		input Transaction
		want  want
	}{
		{
			"is successful json 1",
			transaction,
			want{
				false,
				"",
				"6c001fb7d0a599ddca61b88dc203eeefbac341422cdfd4ac01ff8617e1aa0d9c050001c756189bc655cc487d57e5fefe482449dbe00c390000",
			},
		},
		{
			"is successful json 2",
			transaction2,
			want{
				false,
				"",
				"6c00490dc9520ec45270f240a3cc4f07aec76adc358d9617b693089fcd01000001fcc0bee1480bfca3a80481904cee4099400b1c8d00ff020000004f020000004a0358053d036d0743035d0100000024747a324c324875686161536e663653684544646854454172356a475057504e7770766342031e0743036a0002034f034d031b051f02000000020320",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			transaction, err := tt.input.forge()
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.operation, hex.EncodeToString(transaction))
		})
	}
}
