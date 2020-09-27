package rpc

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_PreapplyOperation(t *testing.T) {
	type input struct {
		handler                 http.Handler
		preapplyOperationsInput PreapplyOperationsInput
	}

	type want struct {
		err         bool
		errContains string
		operations  []Operations
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"handles invalid input",
			input{
				gtGoldenHTTPMock(
					preapplyOperationsHandlerMock(
						readResponse(rpcerrors),
						blankHandler,
					),
				),
				PreapplyOperationsInput{},
			},
			want{
				true,
				"invalid input: Key: 'PreapplyOperationsInput.Blockhash'",
				nil,
			},
		},
		{
			"handles rpc error",
			input{
				gtGoldenHTTPMock(
					preapplyOperationsHandlerMock(
						readResponse(rpcerrors),
						blankHandler,
					),
				),
				PreapplyOperationsInput{
					Blockhash: "some_hash",
					Operations: []Operations{
						{
							Protocol:  "some_protocol",
							Signature: "some_sig",
							Contents:  Contents{},
						},
					},
				},
			},
			want{
				true,
				"failed to preapply operation",
				nil,
			},
		},
		{
			"handles failure to unmarshal",
			input{
				gtGoldenHTTPMock(
					preapplyOperationsHandlerMock(
						[]byte("junk"),
						blankHandler,
					),
				),
				PreapplyOperationsInput{
					Blockhash: "some_hash",
					Operations: []Operations{
						{
							Protocol:  "some_protocol",
							Signature: "some_sig",
							Contents:  Contents{},
						},
					},
				},
			},
			want{
				true,
				"failed to unmarshal operation",
				nil,
			},
		},
		{
			"is successful",
			input{
				gtGoldenHTTPMock(
					preapplyOperationsHandlerMock(
						readResponse(preapplyOperations),
						blankHandler,
					),
				),
				PreapplyOperationsInput{
					Blockhash: "some_hash",
					Operations: []Operations{
						{
							Protocol:  "some_protocol",
							Signature: "some_sig",
							Contents:  Contents{},
						},
					},
				},
			},
			want{
				false,
				"",
				[]Operations{
					{
						Contents: Contents{
							{
								Kind:         "transaction",
								Source:       "tz1W3HW533csCBLor4NPtU79R2TT2sbKfJDH",
								Fee:          "3000",
								Counter:      "1263232",
								GasLimit:     "20000",
								StorageLimit: "0",
								Amount:       "50",
								Destination:  "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
								Metadata: &ContentsMetadata{
									BalanceUpdates: []BalanceUpdates{
										{
											Kind:     "contract",
											Contract: "tz1W3HW533csCBLor4NPtU79R2TT2sbKfJDH",
											Change:   "-3000",
										},
										{
											Kind:     "freezer",
											Category: "fees",
											Delegate: "tz1Ke2h7sDdakHJQh8WX4Z372du1KChsksyU",
											Cycle:    229,
											Change:   "3000",
										},
									},
									OperationResults: &OperationResults{
										Status: "applied",
										BalanceUpdates: []BalanceUpdates{
											{
												Kind:     "contract",
												Contract: "tz1W3HW533csCBLor4NPtU79R2TT2sbKfJDH",
												Change:   "-50",
											},
											{
												Kind:     "contract",
												Contract: "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
												Change:   "50",
											},
										},
										ConsumedGas: "10207",
									},
								},
							},
						},
						Signature: "edsig...."},
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input.handler)
			defer server.Close()

			rpc, err := New(server.URL)
			assert.Nil(t, err)

			operations, err := rpc.PreapplyOperations(tt.input.preapplyOperationsInput)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.operations, operations)
		})
	}
}

func Test_InjectOperation(t *testing.T) {
	goldenOp := "a732d3520eeaa3de98d78e5e5cb6c85f72204fd46feb9f76853841d4a701add36c0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e00b960000008ba0cb2fad622697145cf1665124096d25bc31e006c0008ba0cb2fad622697145cf1665124096d25bc31ed3e7bd1008d3bb0300b1a803000008ba0cb2fad622697145cf1665124096d25bc31e00"
	goldenHash := []byte(`"oopfasdfadjkfalksj"`)

	type input struct {
		handler http.Handler
	}

	type want struct {
		err         bool
		errContains string
		result      string
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"returns rpc error",
			input{
				gtGoldenHTTPMock(
					injectionOperationHandlerMock(
						readResponse(rpcerrors),
						blankHandler,
					),
				),
			},
			want{
				true,
				"failed to inject operation",
				"",
			},
		},
		{
			"handles failure to unmarshal",
			input{
				gtGoldenHTTPMock(
					injectionOperationHandlerMock(
						[]byte("junk"),
						blankHandler,
					),
				),
			},
			want{
				true,
				"failed to unmarshal operation",
				"",
			},
		},
		{
			"is successful",
			input{
				gtGoldenHTTPMock(
					injectionOperationHandlerMock(
						goldenHash,
						blankHandler,
					),
				),
			},
			want{
				false,
				"",
				"oopfasdfadjkfalksj",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input.handler)
			defer server.Close()

			rpc, err := New(server.URL)
			assert.Nil(t, err)

			result, err := rpc.InjectionOperation(InjectionOperationInput{
				Operation: goldenOp,
			})
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.result, result)
		})
	}
}

func Test_InjectBlock(t *testing.T) {
	goldenRPCError := readResponse(rpcerrors)
	goldenHash := []byte("some_hash")
	type input struct {
		handler http.Handler
	}

	type want struct {
		err         bool
		errContains string
		result      []byte
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"returns rpc error",
			input{
				gtGoldenHTTPMock(
					injectionBlockHandlerMock(
						readResponse(rpcerrors),
						blankHandler,
					),
				),
			},
			want{
				true,
				"failed to inject block",
				goldenRPCError,
			},
		},
		{
			"is successful",
			input{
				gtGoldenHTTPMock(
					injectionBlockHandlerMock(
						goldenHash,
						blankHandler,
					),
				),
			},
			want{
				false,
				"",
				goldenHash,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input.handler)
			defer server.Close()

			rpc, err := New(server.URL)
			assert.Nil(t, err)

			result, err := rpc.InjectionBlock(InjectionBlockInput{
				Block: &Block{},
			})
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.result, result)
		})
	}
}

func Test_Counter(t *testing.T) {
	goldenCounter := 10
	goldenRPCError := readResponse(rpcerrors)
	type input struct {
		handler http.Handler
	}

	type want struct {
		err         bool
		errContains string
		counter     int
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"failed to unmarshal counter",
			input{
				gtGoldenHTTPMock(
					counterHandlerMock(
						[]byte(`bad_counter_data`),
						blankHandler,
					),
				),
			},
			want{
				true,
				"failed to unmarshal counter",
				0,
			},
		},
		{
			"returns rpc error",
			input{
				gtGoldenHTTPMock(
					counterHandlerMock(
						goldenRPCError,
						blankHandler,
					),
				),
			},
			want{
				true,
				"failed to get counter",
				0,
			},
		},
		{
			"is successful",
			input{
				gtGoldenHTTPMock(
					counterHandlerMock(
						readResponse(counter),
						blankHandler,
					),
				),
			},
			want{
				false,
				"",
				goldenCounter,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input.handler)
			defer server.Close()

			rpc, err := New(server.URL)
			assert.Nil(t, err)

			counter, err := rpc.Counter(CounterInput{
				mockBlockHash,
				mockAddressTz1,
			})
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.counter, counter)
		})
	}
}

func Test_StripBranchFromForgedOperation(t *testing.T) {
	op := "a732d3520eeaa3de98d78e5e5cb6c85f72204fd46feb9f76853841d4a701add36d0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e00928fe29c01ff0008ba0cb2fad622697145cf1665124096d25bc31e000000c602000000c105000764085e036c055f036d0000000325646f046c000000082564656661756c740501035d050202000000950200000012020000000d03210316051f02000000020317072e020000006a0743036a00000313020000001e020000000403190325072c020000000002000000090200000004034f0327020000000b051f02000000020321034c031e03540348020000001e020000000403190325072c020000000002000000090200000004034f0327034f0326034202000000080320053d036d03420000001a0a000000150008ba0cb2fad622697145cf1665124096d25bc31e"
	branch, _, err := stripBranchFromForgedOperation(op, false)
	assert.Nil(t, err)
	assert.Equal(t, "BLyvCRkxuTXkx1KeGvrcEXiPYj4p1tFxzvFDhoHE7SFKtmP1rbk", branch)
}

func Test_UnForgeOperationWithRPC(t *testing.T) {
	type input struct {
		inputHandler http.Handler
		operation    UnforgeOperationInput
	}

	type want struct {
		err         bool
		errContains string
		operations  []Operations
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"handles invalid input",
			input{
				gtGoldenHTTPMock(unforgeOperationWithRPCMock(readResponse(rpcerrors), blankHandler)),
				UnforgeOperationInput{},
			},
			want{
				true,
				"invalid input: Key: 'UnforgeOperationInput.Blockhash'",
				[]Operations{},
			},
		},
		{
			"handles rpc error",
			input{
				gtGoldenHTTPMock(unforgeOperationWithRPCMock(readResponse(rpcerrors), blankHandler)),
				UnforgeOperationInput{
					Blockhash: "some_hash",
					Operations: []UnforgeOperation{
						{
							Data:   "some_data",
							Branch: "some_branch",
						},
					},
				},
			},
			want{
				true,
				"failed to unforge forge operations with RPC",
				[]Operations{},
			},
		},
		{
			"handles failure to unmarshal",
			input{
				gtGoldenHTTPMock(unforgeOperationWithRPCMock([]byte(`junk`), blankHandler)),
				UnforgeOperationInput{
					Blockhash: "some_hash",
					Operations: []UnforgeOperation{
						{
							Data:   "some_data",
							Branch: "some_branch",
						},
					},
				},
			},
			want{
				true,
				"failed to unforge forge operations with RPC: invalid character",
				[]Operations{},
			},
		},
		{
			"is successful",
			input{
				gtGoldenHTTPMock(unforgeOperationWithRPCMock(readResponse(parseOperations), blankHandler)),
				UnforgeOperationInput{
					Blockhash: "some_hash",
					Operations: []UnforgeOperation{
						{
							Data:   "some_data",
							Branch: "some_branch",
						},
					},
				},
			},
			want{
				false,
				"",
				[]Operations{
					{
						Protocol: "",
						ChainID:  "",
						Hash:     "",
						Branch:   "BLz6yCE4BUL4ppo1zsEWdK9FRCt15WAY7ECQcuK9RtWg4xeEVL7",
						Contents: Contents{
							{
								Kind:         "transaction",
								Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
								Fee:          "10100",
								Counter:      "10",
								GasLimit:     "10100",
								StorageLimit: "0",
								Amount:       "12345",
								Destination:  "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
								Delegate:     "",
								Secret:       "",
								Level:        0,
								Period:       0,
								Proposal:     "",
								Proposals:    []string(nil),
								Ballot:       "",
								Metadata:     nil,
							},
						},
						Signature: "edsigtXomBKi5CTRf5cjATJWSyaRvhfYNHqSUGrn4SdbYRcGwQrUGjzEfQDTuqHhuA8b2d8NarZjz8TRf65WkpQmo423BtomS8Q",
					},
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input.inputHandler)
			defer server.Close()

			c, err := New(server.URL)
			assert.Nil(t, err)

			//"BLyvCRkxuTXkx1KeGvrcEXiPYj4p1tFxzvFDhoHE7SFKtmP1rbk", tt.input.operation
			op, err := c.UnforgeOperation(tt.input.operation)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.operations, op)
		})
	}
}

// func Test_ForgeOperationWithRPC(t *testing.T) {
// 	type input struct {
// 		inputHandler               http.Handler
// 		forgeOperationWithRPCInput ForgeOperationWithRPCInput
// 	}

// 	type want struct {
// 		err         bool
// 		errContains string
// 		operation   string
// 	}

// 	cases := []struct {
// 		name  string
// 		input input
// 		want  want
// 	}{
// 		{
// 			"handles invalid input",
// 			input{
// 				gtGoldenHTTPMock(forgeOperationWithRPCMock(readResponse(rpcerrors), blankHandler)),
// 				ForgeOperationWithRPCInput{},
// 			},
// 			want{
// 				true,
// 				"invalid input: Key: 'ForgeOperationWithRPCInput.Blockhash'",
// 				"",
// 			},
// 		},
// 		{
// 			"handles rpc error",
// 			input{
// 				gtGoldenHTTPMock(forgeOperationWithRPCMock(readResponse(rpcerrors), blankHandler)),
// 				ForgeOperationWithRPCInput{
// 					Blockhash: "some_hash",
// 					Branch:    "some_branch",
// 					Contents:  []Contents{},
// 				},
// 			},
// 			want{
// 				true,
// 				"failed to forge operation: rpc error (somekind)",
// 				"",
// 			},
// 		},
// 		{
// 			"handles failure to unmarshal",
// 			input{
// 				gtGoldenHTTPMock(forgeOperationWithRPCMock([]byte(`junk`), blankHandler)),
// 				ForgeOperationWithRPCInput{
// 					Blockhash: "some_hash",
// 					Branch:    "some_branch",
// 					Contents:  []Contents{},
// 				},
// 			},
// 			want{
// 				true,
// 				"failed to forge operation: invalid character",
// 				"",
// 			},
// 		},
// 		{
// 			"handles failure to strip operation branch",
// 			input{
// 				gtGoldenHTTPMock(forgeOperationWithRPCMock([]byte(`"some_junk_op_string"`), unforgeOperationWithRPCMock(readResponse(rpcerrors), blankHandler))),
// 				ForgeOperationWithRPCInput{
// 					Blockhash: "some_hash",
// 					Branch:    "some_branch",
// 					Contents:  []Contents{},
// 				},
// 			},
// 			want{
// 				true,
// 				"failed to forge operation: unable to verify rpc returned a valid contents",
// 				"some_junk_op_string",
// 			},
// 		},
// 		{
// 			"handles failure to parse forged operation",
// 			input{
// 				gtGoldenHTTPMock(forgeOperationWithRPCMock([]byte(`"some_operation_string"`), unforgeOperationWithRPCMock(readResponse(rpcerrors), blankHandler))),
// 				ForgeOperationWithRPCInput{
// 					Blockhash: "some_hash",
// 					Branch:    "some_branch",
// 					Contents:  []Contents{},
// 				},
// 			},
// 			want{
// 				true,
// 				"failed to forge operation: unable to verify rpc returned a valid contents",
// 				"some_operation_string",
// 			},
// 		},
// 		{
// 			"handles failure to match forge with expected contents",
// 			input{
// 				gtGoldenHTTPMock(forgeOperationWithRPCMock([]byte(`"a79ec80dba1f8ddb2cde90b8f12f7c62fdc36556030281ff8904a3d0df82cddc08000008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e00b960000008ba0cb2fad622697145cf1665124096d25bc31e00"`), unforgeOperationWithRPCMock(readResponse(parseOperations), blankHandler))),
// 				ForgeOperationWithRPCInput{
// 					Blockhash: "some_hash",
// 					Branch:    "some_branch",
// 					Contents: []Contents{
// 						{
// 							Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
// 							Fee:          NewInt(100),
// 							Counter:      NewInt(10),
// 							GasLimit:     NewInt(10100),
// 							StorageLimit: NewInt(0),
// 							Amount:       NewInt(12345),
// 							Destination:  "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
// 							Kind:         TRANSACTIONOP,
// 						},
// 					},
// 				},
// 			},
// 			want{
// 				true,
// 				"failed to forge operation: alert rpc returned invalid contents",
// 				"a79ec80dba1f8ddb2cde90b8f12f7c62fdc36556030281ff8904a3d0df82cddc08000008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e00b960000008ba0cb2fad622697145cf1665124096d25bc31e00",
// 			},
// 		},
// 		{
// 			"is successful",
// 			input{
// 				gtGoldenHTTPMock(forgeOperationWithRPCMock([]byte(`"a79ec80dba1f8ddb2cde90b8f12f7c62fdc36556030281ff8904a3d0df82cddc08000008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e00b960000008ba0cb2fad622697145cf1665124096d25bc31e00"`), unforgeOperationWithRPCMock(readResponse(parseOperations), blankHandler))),
// 				ForgeOperationWithRPCInput{
// 					Blockhash: "some_hash",
// 					Branch:    "some_branch",
// 					Contents: []Contents{
// 						{
// 							Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
// 							Fee:          NewInt(10100),
// 							Counter:      NewInt(10),
// 							GasLimit:     NewInt(10100),
// 							StorageLimit: NewInt(0),
// 							Amount:       NewInt(12345),
// 							Destination:  "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
// 							Kind:         TRANSACTIONOP,
// 						},
// 					},
// 				},
// 			},
// 			want{
// 				false,
// 				"",
// 				"a79ec80dba1f8ddb2cde90b8f12f7c62fdc36556030281ff8904a3d0df82cddc08000008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e00b960000008ba0cb2fad622697145cf1665124096d25bc31e00",
// 			},
// 		},
// 	}

// 	for _, tt := range cases {
// 		t.Run(tt.name, func(t *testing.T) {
// 			server := httptest.NewServer(tt.input.inputHandler)
// 			defer server.Close()

// 			rpc, err := New(server.URL)
// 			assert.Nil(t, err)

// 			op, err := rpc.ForgeOperationWithRPC(tt.input.forgeOperationWithRPCInput)
// 			checkErr(t, tt.want.err, tt.want.errContains, err)
// 			assert.Equal(t, tt.want.operation, op)
// 		})
// 	}
// }
