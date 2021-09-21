package forge

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/completium/go-tezos/v4/internal/testutils"
	"github.com/completium/go-tezos/v4/rpc"
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

	var transaction rpc.Content
	err := json.Unmarshal(transactionJSON, &transaction)
	testutils.CheckErr(t, false, "", err)

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

	var transaction2 rpc.Content
	err = json.Unmarshal(transactionJSON2, &transaction2)
	testutils.CheckErr(t, false, "", err)

	transactionJSON3 := []byte(`{
		"kind": "transaction",
		"source": "tz1f2MeahW6XMLcfHJSU5VH8USC4EuFiwdhx",
		"fee": "1188",
		"counter": "6",
		"gas_limit": "10307",
		"storage_limit": "0",
		"amount": "50000000000",
		"destination": "tz1aWXP237BLwNHJcCD4b3DutCevhqq2T1Z9"
	  }`)

	var transaction3 rpc.Content
	err = json.Unmarshal(transactionJSON3, &transaction3)
	testutils.CheckErr(t, false, "", err)

	revealJSON := []byte(`{
		"kind": "reveal",
		"source": "tz1f2MeahW6XMLcfHJSU5VH8USC4EuFiwdhx",
		"fee": "1257",
		"counter": "5",
		"gas_limit": "10000",
		"storage_limit": "0",
		"public_key": "edpkuEmaQSYKgDj5k9wfE3bTxjfjoG9k5YvRmYZsGf2bjEymZKkzNn"
	  }`)

	var reveal rpc.Content
	err = json.Unmarshal(revealJSON, &reveal)
	testutils.CheckErr(t, false, "", err)

	type input struct {
		branch   string
		contents []rpc.Content
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
				contents: []rpc.Content{
					transaction,
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
				contents: []rpc.Content{
					transaction2,
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
				contents: []rpc.Content{
					transaction,
					transaction2,
				},
			},
			want{
				false,
				"",
				"452b8599b0e4960b884d3ad61c89c594bc3348c798842651d3fa6cfafa77ce556c007fd82c06cf5a203f18faaf562447ed1efcc6c010830a07c350008090dfc04a0000a31e81ac3425310e3274a4698a793b2839dc0afa006c00490dc9520ec45270f240a3cc4f07aec76adc358d9617b693089fcd01000001fcc0bee1480bfca3a80481904cee4099400b1c8d00ff020000004f020000004a0358053d036d0743035d0100000024747a324c324875686161536e663653684544646854454172356a475057504e7770766342031e0743036a0002034f034d031b051f02000000020320",
			},
		},
		{
			"is successful with reveal and transaction",
			input{
				branch: "BLCFdxw2kWJfCk9TWQsYxrQd9CcPPs2YdbArbDDgL4GZTYvTfZN",
				contents: []rpc.Content{
					reveal,
					transaction3,
				},
			},
			want{
				false,
				"",
				"3f82cf0634a5965032d087daa63cf3603dd0f0325e2d670fee91b20486caa0d36b00d4a35d6c49ffbaa32b40e96c844dc485b0cdb5fae90905904e00004e7097e206a9afa864475095b58009014f9c24efd54c5d40240c1e807b4ab80c6c00d4a35d6c49ffbaa32b40e96c844dc485b0cdb5faa40906c3500080e8eda1ba010000a31e81ac3425310e3274a4698a793b2839dc0afa00",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			operation, err := Encode(tt.input.branch, tt.input.contents...)
			testutils.CheckErr(t, tt.want.err, tt.want.errContains, err)
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

	var origination rpc.Origination
	err := json.Unmarshal(originationJSON, &origination)
	testutils.CheckErr(t, false, "", err)

	originationJSON2 := []byte(`{
		"kind": "origination",
		"source": "tz1TJCwoX79reCZ8yccPeW8iB9Mba91v8H47",
		"fee": "2070",
		"counter": "307027",
		"gas_limit": "15919",
		"storage_limit": "526",
		"balance": "0",
		"script": {
		  "code": [
			{
			  "prim": "parameter",
			  "args": [
				{
				  "prim": "unit"
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
					"prim": "DUP"
				  },
				  {
					"prim": "DIP",
					"args": [
					  [
						{
						  "prim": "CDR"
						}
					  ]
					]
				  },
				  {
					"prim": "CAR"
				  },
				  {
					"prim": "PUSH",
					"args": [
					  {
						"prim": "address"
					  },
					  {
						"string": "KT1M8MStwA1R5SuGx2V6AMvgd8dGAFYcEUmu"
					  }
					]
				  },
				  {
					"prim": "CONTRACT",
					"args": [
					  {
						"prim": "or",
						"args": [
						  {
							"prim": "lambda",
							"args": [
							  {
								"prim": "unit"
							  },
							  {
								"prim": "list",
								"args": [
								  {
									"prim": "operation"
								  }
								]
							  }
							],
							"annots": [
							  "%do"
							]
						  },
						  {
							"prim": "unit",
							"annots": [
							  "%default"
							]
						  }
						]
					  }
					]
				  },
				  {
					"prim": "IF_NONE",
					"args": [
					  [
						{
						  "prim": "PUSH",
						  "args": [
							{
							  "prim": "string"
							},
							{
							  "string": "type mismatch"
							}
						  ]
						},
						{
						  "prim": "FAILWITH"
						}
					  ],
					  [
						[
						  {
							"prim": "DIP",
							"args": [
							  {
								"int": "2"
							  },
							  [
								{
								  "prim": "DUP"
								}
							  ]
							]
						  },
						  {
							"prim": "DIG",
							"args": [
							  {
								"int": "3"
							  }
							]
						  }
						],
						{
						  "prim": "NIL",
						  "args": [
							{
							  "prim": "operation"
							}
						  ]
						},
						[
						  {
							"prim": "DIP",
							"args": [
							  {
								"int": "2"
							  },
							  [
								{
								  "prim": "DUP"
								}
							  ]
							]
						  },
						  {
							"prim": "DIG",
							"args": [
							  {
								"int": "3"
							  }
							]
						  }
						],
						{
						  "prim": "DIP",
						  "args": [
							{
							  "int": "3"
							},
							[
							  {
								"prim": "DROP"
							  }
							]
						  ]
						},
						{
						  "prim": "PUSH",
						  "args": [
							{
							  "prim": "mutez"
							},
							{
							  "int": "1000000"
							}
						  ]
						},
						{
						  "prim": "UNIT"
						},
						{
						  "prim": "RIGHT",
						  "args": [
							{
							  "prim": "lambda",
							  "args": [
								{
								  "prim": "unit"
								},
								{
								  "prim": "list",
								  "args": [
									{
									  "prim": "operation"
									}
								  ]
								}
							  ]
							}
						  ]
						},
						{
						  "prim": "TRANSFER_TOKENS"
						},
						{
						  "prim": "CONS"
						},
						{
						  "prim": "PAIR"
						}
					  ]
					]
				  },
				  {
					"prim": "DIP",
					"args": [
					  [
						{
						  "prim": "DROP"
						},
						{
						  "prim": "DROP"
						}
					  ]
					]
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

	var origination2 rpc.Origination
	err = json.Unmarshal(originationJSON2, &origination2)
	testutils.CheckErr(t, false, "", err)

	type want struct {
		err         bool
		errContains string
		operation   string
	}

	cases := []struct {
		name  string
		input rpc.Origination
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
		{
			"is successful 2",
			origination2,
			want{
				false,
				"",
				"6d0054013ef6636fe99989a26006622bf270be0b14859610d3de12af7c8e040000000000ef02000000ea0500036c0501036c050202000000db0321051f0200000002031703160743036e01000000244b54314d384d5374774131523553754778325636414d7667643864474146596345556d7505550764085e036c055f036d0000000325646f046c000000082564656661756c74072f020000001807430368010000000d74797065206d69736d6174636803270200000051020000000f071f00020200000002032105700003053d036d020000000f071f00020200000002032105700003071f0003020000000203200743036a0080897a034f0544075e036c055f036d034d031b0342051f02000000040320032000000002030b",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			origination, err := forgeOrigination(tt.input)
			testutils.CheckErr(t, tt.want.err, tt.want.errContains, err)
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

	var transaction rpc.Transaction
	err := json.Unmarshal(transactionJSON, &transaction)
	testutils.CheckErr(t, false, "", err)

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

	var transaction2 rpc.Transaction
	err = json.Unmarshal(transactionJSON2, &transaction2)
	testutils.CheckErr(t, false, "", err)

	type want struct {
		err         bool
		errContains string
		operation   string
	}

	cases := []struct {
		name  string
		input rpc.Transaction
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
			transaction, err := forgeTransaction(tt.input)
			testutils.CheckErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.operation, hex.EncodeToString(transaction))
		})
	}
}

func Test_Forge_Reveal(t *testing.T) {
	revealJSON := []byte(`{
		"kind": "reveal",
		"source": "tz1f2MeahW6XMLcfHJSU5VH8USC4EuFiwdhx",
		"fee": "1257",
		"counter": "5",
		"gas_limit": "10000",
		"storage_limit": "0",
		"public_key": "edpkuEmaQSYKgDj5k9wfE3bTxjfjoG9k5YvRmYZsGf2bjEymZKkzNn"
	  }`)

	var reveal rpc.Reveal
	err := json.Unmarshal(revealJSON, &reveal)
	testutils.CheckErr(t, false, "", err)

	type want struct {
		err         bool
		errContains string
		operation   string
	}

	cases := []struct {
		name  string
		input rpc.Reveal
		want  want
	}{
		{
			"is successful",
			reveal,
			want{
				false,
				"",
				"6b00d4a35d6c49ffbaa32b40e96c844dc485b0cdb5fae90905904e00004e7097e206a9afa864475095b58009014f9c24efd54c5d40240c1e807b4ab80c",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			reveal, err := forgeReveal(tt.input)
			testutils.CheckErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.operation, hex.EncodeToString(reveal))
		})
	}
}

func Test_IntExpression(t *testing.T) {
	val, err := IntExpression(9)
	testutils.CheckErr(t, false, "", err)
	assert.Equal(t, "exprtvAzqNE9zfpBLL9nKEaY1Dd2rznyG9iTFtECJvDkuub1bj3XvW", val)

	val, err = IntExpression(-9)
	testutils.CheckErr(t, false, "", err)
	assert.Equal(t, "exprvH9jru3NJN4ZTNwwkCdC1PPLkWLWCoe6JxhcJ3a39mD5Bd4NH4", val)
}

func Test_NatExpression(t *testing.T) {
	val, err := NatExpression(9)
	testutils.CheckErr(t, false, "", err)
	assert.Equal(t, "exprtvAzqNE9zfpBLL9nKEaY1Dd2rznyG9iTFtECJvDkuub1bj3XvW", val)
}

func Test_AddressExpression(t *testing.T) {
	val, err := AddressExpression(`tz1S82rGFZK8cVbNDpP1Hf9VhTUa4W8oc2WV`)
	testutils.CheckErr(t, false, "", err)
	assert.Equal(t, "expruwEtkquVj9E92Wc7KTFMSnCqVGZ4KPngpspNmRTm6rX6KZbcvH", val)
}

func Test_StringExpression(t *testing.T) {
	val, err := StringExpression("Tezos Tacos Nachos")
	testutils.CheckErr(t, false, "", err)
	assert.Equal(t, "expruGmscHLuUazE7d79EepWCnDuPJreo8R87wsDGUgKAuH4E5ayEj", val)
}

func Test_KeyHashExpression(t *testing.T) {
	val, err := KeyHashExpression(`tz1eEnQhbwf6trb8Q8mPb2RaPkNk2rN7BKi8`)
	testutils.CheckErr(t, false, "", err)
	assert.Equal(t, "expruqnFVtyPKd2KcrjkiJTaqE1WU1fEf8K1ajHvzgKz5pcc5sZyjn", val)
}

func Test_BytesExpression(t *testing.T) {
	v, _ := hex.DecodeString(`0a0a0a`)
	val, err := BytesExpression(v)
	testutils.CheckErr(t, false, "", err)
	assert.Equal(t, "exprunb7V121UYKTTbQGj6UQrpgXcZE3F71TrNMUkw9WtARMzht9tN", val)
}

func Test_MichelineExpression(t *testing.T) {
	v := `{ "prim": "Pair", "args": [ { "int": "1" }, { "int": "12" } ] }`
	val, err := MichelineExpression(v)
	testutils.CheckErr(t, false, "", err)
	assert.Equal(t, "exprupozG51AtT7yZUy5sg6VbJQ4b9omAE1PKD2PXvqi2YBuZqoKG3", val)
}
