package gotezos

import (
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_PreapplyOperation(t *testing.T) {
	goldenHash := []byte("some_hash")
	type input struct {
		handler http.Handler
	}

	type want struct {
		err         bool
		errContains string
		result      *[]byte
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"returns block RPC error",
			input{
				gtGoldenHTTPMock(
					preapplyOperationsHandlerMock(
						mockRPCErrorResp,
						mockRPCErrorResp,
						blankHandler,
					),
				),
			},
			want{
				true,
				"failed to preapply operation",
				nil,
			},
		},
		{
			"returns preapply rpc error",
			input{
				gtGoldenHTTPMock(
					preapplyOperationsHandlerMock(
						mockRPCErrorResp,
						mockBlockResp,
						blankHandler,
					),
				),
			},
			want{
				true,
				"failed to preapply operation",
				&mockRPCErrorResp,
			},
		},
		{
			"is successful",
			input{
				gtGoldenHTTPMock(
					preapplyOperationsHandlerMock(
						goldenHash,
						mockBlockResp,
						blankHandler,
					),
				),
			},
			want{
				false,
				"",
				&goldenHash,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input.handler)
			defer server.Close()

			gt, err := New(server.URL)
			assert.Nil(t, err)

			result, err := gt.PreapplyOperations(mockBlockHash, []Contents{}, "")
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.result, result)
		})
	}
}

func Test_InjectOperation(t *testing.T) {
	goldenOp := "a732d3520eeaa3de98d78e5e5cb6c85f72204fd46feb9f76853841d4a701add36c0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e00b960000008ba0cb2fad622697145cf1665124096d25bc31e006c0008ba0cb2fad622697145cf1665124096d25bc31ed3e7bd1008d3bb0300b1a803000008ba0cb2fad622697145cf1665124096d25bc31e00"
	goldenHash := []byte("some_hash")
	type input struct {
		handler http.Handler
	}

	type want struct {
		err         bool
		errContains string
		result      *[]byte
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
						mockRPCErrorResp,
						blankHandler,
					),
				),
			},
			want{
				true,
				"failed to inject operation",
				&mockRPCErrorResp,
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
				&goldenHash,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input.handler)
			defer server.Close()

			gt, err := New(server.URL)
			assert.Nil(t, err)

			result, err := gt.InjectionOperation(&InjectionOperationInput{
				Operation: &goldenOp,
			})
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.result, result)
		})
	}
}

func Test_InjectBlock(t *testing.T) {
	goldenHash := []byte("some_hash")
	type input struct {
		handler http.Handler
	}

	type want struct {
		err         bool
		errContains string
		result      *[]byte
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
						mockRPCErrorResp,
						blankHandler,
					),
				),
			},
			want{
				true,
				"failed to inject block",
				&mockRPCErrorResp,
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
				&goldenHash,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input.handler)
			defer server.Close()

			gt, err := New(server.URL)
			assert.Nil(t, err)

			result, err := gt.InjectionBlock(&InjectionBlockInput{
				Block: &Block{},
			})
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.result, result)
		})
	}
}

func Test_Counter(t *testing.T) {
	goldenCounter := 10
	type input struct {
		handler http.Handler
	}

	type want struct {
		err         bool
		errContains string
		counter     *int
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
				nil,
			},
		},
		{
			"returns rpc error",
			input{
				gtGoldenHTTPMock(
					counterHandlerMock(
						mockRPCErrorResp,
						blankHandler,
					),
				),
			},
			want{
				true,
				"failed to get counter",
				nil,
			},
		},
		{
			"is successful",
			input{
				gtGoldenHTTPMock(
					counterHandlerMock(
						mockCounterResp,
						blankHandler,
					),
				),
			},
			want{
				false,
				"",
				&goldenCounter,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input.handler)
			defer server.Close()

			gt, err := New(server.URL)
			assert.Nil(t, err)

			counter, err := gt.Counter(mockBlockHash, mockAddressTz1)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.counter, counter)
		})
	}
}

func Test_ForgeOperation(t *testing.T) {
	var (
		transactionOp = "a732d3520eeaa3de98d78e5e5cb6c85f72204fd46feb9f76853841d4a701add36c0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e00b960000008ba0cb2fad622697145cf1665124096d25bc31e006c0008ba0cb2fad622697145cf1665124096d25bc31ed3e7bd1008d3bb0300b1a803018b88e99e66c1c2587f87118449f781cb7d44c9c40000"
		revealOp      = "a732d3520eeaa3de98d78e5e5cb6c85f72204fd46feb9f76853841d4a701add36b0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e0000136083897bc97879c53e3e7855838fbbc87303ddd376080fc3d3e136b55d028b6b0008ba0cb2fad622697145cf1665124096d25bc31ed3e7bd1008d3bb030000136083897bc97879c53e3e7855838fbbc87303ddd376080fc3d3e136b55d028b"
		originationOp = "a732d3520eeaa3de98d78e5e5cb6c85f72204fd46feb9f76853841d4a701add36d0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e00928fe29c01ff0008ba0cb2fad622697145cf1665124096d25bc31e000000c602000000c105000764085e036c055f036d0000000325646f046c000000082564656661756c740501035d050202000000950200000012020000000d03210316051f02000000020317072e020000006a0743036a00000313020000001e020000000403190325072c020000000002000000090200000004034f0327020000000b051f02000000020321034c031e03540348020000001e020000000403190325072c020000000002000000090200000004034f0327034f0326034202000000080320053d036d03420000001a0a000000150008ba0cb2fad622697145cf1665124096d25bc31e"
		delegationOp  = "a732d3520eeaa3de98d78e5e5cb6c85f72204fd46feb9f76853841d4a701add36e0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e00ff0008ba0cb2fad622697145cf1665124096d25bc31e"
	)
	type input struct {
		contents []Contents
		branch   string
	}

	type want struct {
		err         bool
		errContains string
		operation   *string
	}
	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"is successful transaction",
			input{
				[]Contents{
					Contents{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          BigInt{*big.NewInt(10100)},
						Counter:      BigInt{*big.NewInt(10)},
						GasLimit:     BigInt{*big.NewInt(10100)},
						StorageLimit: BigInt{big.Int{}},
						Amount:       BigInt{*big.NewInt(12345)},
						Destination:  "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Kind:         TRANSACTIONOP,
					},
					Contents{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          BigInt{*big.NewInt(34567123)},
						Counter:      BigInt{*big.NewInt(8)},
						GasLimit:     BigInt{*big.NewInt(56787)},
						StorageLimit: BigInt{big.Int{}},
						Amount:       BigInt{*big.NewInt(54321)},
						Destination:  "KT1MJZWHKZU7ViybRLsphP3ppiiTc7myP2aj",
						Kind:         TRANSACTIONOP,
					},
				},
				"BLyvCRkxuTXkx1KeGvrcEXiPYj4p1tFxzvFDhoHE7SFKtmP1rbk",
			},
			want{
				false,
				"",
				&transactionOp,
			},
		},
		{
			"is successful reveal",
			input{
				[]Contents{
					Contents{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          BigInt{*big.NewInt(10100)},
						Counter:      BigInt{*big.NewInt(10)},
						GasLimit:     BigInt{*big.NewInt(10100)},
						StorageLimit: BigInt{big.Int{}},
						Phk:          "edpktnktxAzmXPD9XVNqAvdCFb76vxzQtkbVkSEtXcTz33QZQdb4JQ",
						Kind:         REVEALOP,
					},
					Contents{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          BigInt{*big.NewInt(34567123)},
						Counter:      BigInt{*big.NewInt(8)},
						GasLimit:     BigInt{*big.NewInt(56787)},
						StorageLimit: BigInt{big.Int{}},
						Phk:          "edpktnktxAzmXPD9XVNqAvdCFb76vxzQtkbVkSEtXcTz33QZQdb4JQ",
						Kind:         REVEALOP,
					},
				},
				"BLyvCRkxuTXkx1KeGvrcEXiPYj4p1tFxzvFDhoHE7SFKtmP1rbk",
			},
			want{
				false,
				"",
				&revealOp,
			},
		},
		{
			"is successful origination",
			input{
				[]Contents{
					Contents{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          BigInt{*big.NewInt(10100)},
						Counter:      BigInt{*big.NewInt(10)},
						GasLimit:     BigInt{*big.NewInt(10100)},
						StorageLimit: BigInt{big.Int{}},
						Kind:         ORIGINATIONOP,
						Balance:      BigInt{*big.NewInt(328763282)},
						Delegate:     "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					},
				},
				"BLyvCRkxuTXkx1KeGvrcEXiPYj4p1tFxzvFDhoHE7SFKtmP1rbk",
			},
			want{
				false,
				"",
				&originationOp,
			},
		},
		{
			"is successful delegation",
			input{
				[]Contents{
					Contents{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          BigInt{*big.NewInt(10100)},
						Counter:      BigInt{*big.NewInt(10)},
						GasLimit:     BigInt{*big.NewInt(10100)},
						StorageLimit: BigInt{big.Int{}},
						Kind:         DELEGATIONOP,
						Delegate:     "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					},
				},
				"BLyvCRkxuTXkx1KeGvrcEXiPYj4p1tFxzvFDhoHE7SFKtmP1rbk",
			},
			want{
				false,
				"",
				&delegationOp,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			gt := testGoTezos(t, gtGoldenHTTPMock(blankHandler))
			operation, err := gt.ForgeOperation(tt.input.branch, tt.input.contents...)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			if operation != nil {
				fmt.Println(*operation)
			}

			assert.Equal(t, tt.want.operation, operation)
		})
	}
}

func Test_forgeTransactionOperation(t *testing.T) {
	type input struct {
		contents Contents
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
			"works with tz1 addresses",
			input{
				Contents{
					Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					Fee:          BigInt{*big.NewInt(10100)},
					Counter:      BigInt{*big.NewInt(10)},
					GasLimit:     BigInt{*big.NewInt(10100)},
					StorageLimit: BigInt{big.Int{}},
					Amount:       BigInt{*big.NewInt(30)},
					Destination:  "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					Kind:         TRANSACTIONOP,
				},
			},
			want{
				false,
				"",
				"6c0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e001e000008ba0cb2fad622697145cf1665124096d25bc31e00",
			},
		},
		{
			"works with tz1 to kt",
			input{
				Contents{
					Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					Fee:          BigInt{*big.NewInt(10100)},
					Counter:      BigInt{*big.NewInt(10)},
					GasLimit:     BigInt{*big.NewInt(10100)},
					StorageLimit: BigInt{big.Int{}},
					Amount:       BigInt{*big.NewInt(30)},
					Destination:  "KT1MJZWHKZU7ViybRLsphP3ppiiTc7myP2aj",
					Kind:         TRANSACTIONOP,
				},
			},
			want{
				false,
				"",
				"6c0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e001e018b88e99e66c1c2587f87118449f781cb7d44c9c40000",
			},
		},
		{
			"handles common forge error",
			input{
				Contents{
					Source:       "LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					Fee:          BigInt{*big.NewInt(10100)},
					Counter:      BigInt{*big.NewInt(10)},
					GasLimit:     BigInt{*big.NewInt(10100)},
					StorageLimit: BigInt{big.Int{}},
					Amount:       BigInt{*big.NewInt(30)},
					Destination:  "KT1MJZWHKZU7ViybRLsphP3ppiiTc7myP2aj",
					Kind:         TRANSACTIONOP,
				},
			},
			want{
				true,
				"failed to remove tz1 from source prefix",
				"",
			},
		},
		{
			"handles failed to remove kt prefix from destination",
			input{
				Contents{
					Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					Fee:          BigInt{*big.NewInt(10100)},
					Counter:      BigInt{*big.NewInt(10)},
					GasLimit:     BigInt{*big.NewInt(10100)},
					StorageLimit: BigInt{big.Int{}},
					Amount:       BigInt{*big.NewInt(30)},
					Destination:  "KTJUNK",
					Kind:         TRANSACTIONOP,
				},
			},
			want{
				true,
				"could not forge transaction: provided destination is not a valid KT1 address",
				"",
			},
		},
		{
			"handles failed to remove tz1 prefix from destination",
			input{
				Contents{
					Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					Fee:          BigInt{*big.NewInt(10100)},
					Counter:      BigInt{*big.NewInt(10)},
					GasLimit:     BigInt{*big.NewInt(10100)},
					StorageLimit: BigInt{big.Int{}},
					Amount:       BigInt{*big.NewInt(30)},
					Destination:  "tz1JUNK",
					Kind:         TRANSACTIONOP,
				},
			},
			want{
				true,
				"could not forge transaction: provided destination is not a valid tz1 address",
				"",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			gt := testGoTezos(t, gtGoldenHTTPMock(blankHandler))
			operation, err := gt.forgeTransactionOperation(tt.input.contents)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			fmt.Println(operation)

			assert.Equal(t, tt.want.operation, operation)
		})
	}
}

func Test_forgeRevealOperation(t *testing.T) {
	type input struct {
		contents Contents
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
			"is successful",
			input{
				Contents{
					Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					Fee:          BigInt{*big.NewInt(10100)},
					Counter:      BigInt{*big.NewInt(10)},
					GasLimit:     BigInt{*big.NewInt(10100)},
					StorageLimit: BigInt{big.Int{}},
					Phk:          "edpktnktxAzmXPD9XVNqAvdCFb76vxzQtkbVkSEtXcTz33QZQdb4JQ",
					Kind:         REVEALOP,
				},
			},
			want{
				false,
				"",
				"6b0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e0000136083897bc97879c53e3e7855838fbbc87303ddd376080fc3d3e136b55d028b",
			},
		},
		{
			"handles failure to forge common",
			input{
				Contents{
					Source:       "tAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					Fee:          BigInt{*big.NewInt(10100)},
					Counter:      BigInt{*big.NewInt(10)},
					GasLimit:     BigInt{*big.NewInt(10100)},
					StorageLimit: BigInt{big.Int{}},
					Phk:          "edpktnktxAzmXPD9XVNqAvdCFb76vxzQtkbVkSEtXcTz33QZQdb4JQ",
					Kind:         REVEALOP,
				},
			},
			want{
				true,
				"failed to forge reveal operation: failed to remove tz1 from source prefix",
				"",
			},
		},
		{
			"handles failure to clean pub key",
			input{
				Contents{
					Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					Fee:          BigInt{*big.NewInt(10100)},
					Counter:      BigInt{*big.NewInt(10)},
					GasLimit:     BigInt{*big.NewInt(10100)},
					StorageLimit: BigInt{big.Int{}},
					Phk:          "tnktxAzm32--9XVNqAvdCFb76vxzQtkbVkSEtXcTz33QZQdb4JQ",
					Kind:         REVEALOP,
				},
			},
			want{
				true,
				"failed to forge reveal operation: failed to decode payload",
				"",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			gt := testGoTezos(t, gtGoldenHTTPMock(blankHandler))
			operation, err := gt.forgeRevealOperation(tt.input.contents)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			fmt.Println(operation)

			assert.Equal(t, tt.want.operation, operation)
		})
	}
}

func Test_forgeOriginationOperation(t *testing.T) {
	type input struct {
		contents Contents
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
			"is successful",
			input{
				Contents{
					Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					Fee:          BigInt{*big.NewInt(10100)},
					Counter:      BigInt{*big.NewInt(10)},
					GasLimit:     BigInt{*big.NewInt(10100)},
					StorageLimit: BigInt{big.Int{}},
					Kind:         ORIGINATIONOP,
					Balance:      BigInt{*big.NewInt(328763282)},
					Delegate:     "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
				},
			},
			want{
				false,
				"",
				"6d0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e00928fe29c01ff0008ba0cb2fad622697145cf1665124096d25bc31e000000c602000000c105000764085e036c055f036d0000000325646f046c000000082564656661756c740501035d050202000000950200000012020000000d03210316051f02000000020317072e020000006a0743036a00000313020000001e020000000403190325072c020000000002000000090200000004034f0327020000000b051f02000000020321034c031e03540348020000001e020000000403190325072c020000000002000000090200000004034f0327034f0326034202000000080320053d036d03420000001a0a000000150008ba0cb2fad622697145cf1665124096d25bc31e",
			},
		},
		{
			"handles failure to forge common",
			input{
				Contents{
					Source:       "tAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					Fee:          BigInt{*big.NewInt(10100)},
					Counter:      BigInt{*big.NewInt(10)},
					GasLimit:     BigInt{*big.NewInt(10100)},
					StorageLimit: BigInt{big.Int{}},
					Kind:         ORIGINATIONOP,
					Balance:      BigInt{*big.NewInt(328763282)},
					Delegate:     "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
				},
			},
			want{
				true,
				"failed to forge origination operation: failed to remove tz1 from source prefix",
				"",
			},
		},
		{
			"handles failure to clean delegate",
			input{
				Contents{
					Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					Fee:          BigInt{*big.NewInt(10100)},
					Counter:      BigInt{*big.NewInt(10)},
					GasLimit:     BigInt{*big.NewInt(10100)},
					StorageLimit: BigInt{big.Int{}},
					Kind:         ORIGINATIONOP,
					Balance:      BigInt{*big.NewInt(328763282)},
					Delegate:     "tz1LSAy890--cAVcNdYnXCy18bwVksXci8gUC2YpA",
				},
			},
			want{
				true,
				"failed to forge origination operation: failed to decode payload",
				"",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			gt := testGoTezos(t, gtGoldenHTTPMock(blankHandler))
			operation, err := gt.forgeOriginationOperation(tt.input.contents)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			fmt.Println(operation)

			assert.Equal(t, tt.want.operation, operation)
		})
	}
}

func Test_forgeDelegationOperation(t *testing.T) {
	type input struct {
		contents Contents
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
			"is successful",
			input{
				Contents{
					Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					Fee:          BigInt{*big.NewInt(10100)},
					Counter:      BigInt{*big.NewInt(10)},
					GasLimit:     BigInt{*big.NewInt(10100)},
					StorageLimit: BigInt{big.Int{}},
					Kind:         DELEGATIONOP,
					Delegate:     "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
				},
			},
			want{
				false,
				"",
				"6e0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e00ff0008ba0cb2fad622697145cf1665124096d25bc31e",
			},
		},
		{
			"handles failure to forge common",
			input{
				Contents{
					Source:       "tAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					Fee:          BigInt{*big.NewInt(10100)},
					Counter:      BigInt{*big.NewInt(10)},
					GasLimit:     BigInt{*big.NewInt(10100)},
					StorageLimit: BigInt{big.Int{}},
					Kind:         DELEGATIONOP,
					Delegate:     "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
				},
			},
			want{
				true,
				"failed to forge delegation operation: failed to remove tz1 from source prefix",
				"",
			},
		},
		{
			"handles failure to clean delegate tz1",
			input{
				Contents{
					Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					Fee:          BigInt{*big.NewInt(10100)},
					Counter:      BigInt{*big.NewInt(10)},
					GasLimit:     BigInt{*big.NewInt(10100)},
					StorageLimit: BigInt{big.Int{}},
					Kind:         ORIGINATIONOP,
					Balance:      BigInt{*big.NewInt(328763282)},
					Delegate:     "tz1LSAy890--cAVcNdYnXCy18bwVksXci8gUC2YpA",
				},
			},
			want{
				true,
				"failed to forge delegation operation: failed to decode payload",
				"",
			},
		},
		{
			"handles failure to clean delegate KT1",
			input{
				Contents{
					Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					Fee:          BigInt{*big.NewInt(10100)},
					Counter:      BigInt{*big.NewInt(10)},
					GasLimit:     BigInt{*big.NewInt(10100)},
					StorageLimit: BigInt{big.Int{}},
					Kind:         ORIGINATIONOP,
					Balance:      BigInt{*big.NewInt(328763282)},
					Delegate:     "KT1LSAy890--cAVcNdYnXCy18bwVksXci8gUC2YpA",
				},
			},
			want{
				true,
				"failed to forge delegation operation: failed to decode payload",
				"",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			gt := testGoTezos(t, gtGoldenHTTPMock(blankHandler))
			operation, err := gt.forgeDelegationOperation(tt.input.contents)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			fmt.Println(operation)

			assert.Equal(t, tt.want.operation, operation)
		})
	}
}

func Test_UnforgeOperation(t *testing.T) {
	mockHash := "BLyvCRkxuTXkx1KeGvrcEXiPYj4p1tFxzvFDhoHE7SFKtmP1rbk"
	type input struct {
		handler   http.Handler
		hexString string
		signed    bool
	}

	type want struct {
		err         bool
		errContains string
		contents    *[]Contents
		branch      *string
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"is successful transaction",
			input{
				gtGoldenHTTPMock(blankHandler),
				"a732d3520eeaa3de98d78e5e5cb6c85f72204fd46feb9f76853841d4a701add36c0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e00b960000008ba0cb2fad622697145cf1665124096d25bc31e006c0008ba0cb2fad622697145cf1665124096d25bc31ed3e7bd1008d3bb0300b1a803018b88e99e66c1c2587f87118449f781cb7d44c9c40000",
				false,
			},
			want{
				false,
				"",
				&[]Contents{
					Contents{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          BigInt{*big.NewInt(10100)},
						Counter:      BigInt{*big.NewInt(10)},
						GasLimit:     BigInt{*big.NewInt(10100)},
						StorageLimit: BigInt{big.Int{}},
						Amount:       BigInt{*big.NewInt(12345)},
						Destination:  "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Kind:         TRANSACTIONOP,
					},
					Contents{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          BigInt{*big.NewInt(34567123)},
						Counter:      BigInt{*big.NewInt(8)},
						GasLimit:     BigInt{*big.NewInt(56787)},
						StorageLimit: BigInt{big.Int{}},
						Amount:       BigInt{*big.NewInt(54321)},
						Destination:  "KT1MJZWHKZU7ViybRLsphP3ppiiTc7myP2aj",
						Kind:         TRANSACTIONOP,
					},
				},
				&mockHash,
			},
		},
		{
			"is successful reveal",
			input{
				gtGoldenHTTPMock(blankHandler),
				"a732d3520eeaa3de98d78e5e5cb6c85f72204fd46feb9f76853841d4a701add36b0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e0000136083897bc97879c53e3e7855838fbbc87303ddd376080fc3d3e136b55d028b6b0008ba0cb2fad622697145cf1665124096d25bc31ed3e7bd1008d3bb030000136083897bc97879c53e3e7855838fbbc87303ddd376080fc3d3e136b55d028ba",
				false,
			},
			want{
				false,
				"",
				&[]Contents{
					Contents{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          BigInt{*big.NewInt(10100)},
						Counter:      BigInt{*big.NewInt(10)},
						GasLimit:     BigInt{*big.NewInt(10100)},
						StorageLimit: BigInt{big.Int{}},
						Phk:          "edpktnktxAzmXPD9XVNqAvdCFb76vxzQtkbVkSEtXcTz33QZQdb4JQ",
						Kind:         REVEALOP,
					},
					Contents{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          BigInt{*big.NewInt(34567123)},
						Counter:      BigInt{*big.NewInt(8)},
						GasLimit:     BigInt{*big.NewInt(56787)},
						StorageLimit: BigInt{big.Int{}},
						Phk:          "edpktnktxAzmXPD9XVNqAvdCFb76vxzQtkbVkSEtXcTz33QZQdb4JQ",
						Kind:         REVEALOP,
					},
				},
				&mockHash,
			},
		},
		{
			"is successful origination",
			input{
				gtGoldenHTTPMock(blankHandler),
				"a732d3520eeaa3de98d78e5e5cb6c85f72204fd46feb9f76853841d4a701add36d0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e00928fe29c01ff0008ba0cb2fad622697145cf1665124096d25bc31e000000c602000000c105000764085e036c055f036d0000000325646f046c000000082564656661756c740501035d050202000000950200000012020000000d03210316051f02000000020317072e020000006a0743036a00000313020000001e020000000403190325072c020000000002000000090200000004034f0327020000000b051f02000000020321034c031e03540348020000001e020000000403190325072c020000000002000000090200000004034f0327034f0326034202000000080320053d036d03420000001a0a000000150008ba0cb2fad622697145cf1665124096d25bc31e",
				false,
			},
			want{
				false,
				"",
				&[]Contents{
					Contents{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          BigInt{*big.NewInt(10100)},
						Counter:      BigInt{*big.NewInt(10)},
						GasLimit:     BigInt{*big.NewInt(10100)},
						StorageLimit: BigInt{big.Int{}},
						Kind:         ORIGINATIONOP,
						Balance:      BigInt{*big.NewInt(328763282)},
						Delegate:     "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					},
				},
				&mockHash,
			},
		},
		{
			"is successful delegation",
			input{
				gtGoldenHTTPMock(blankHandler),
				"a732d3520eeaa3de98d78e5e5cb6c85f72204fd46feb9f76853841d4a701add36e0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e00ff0008ba0cb2fad622697145cf1665124096d25bc31e",
				false,
			},
			want{
				false,
				"",
				&[]Contents{
					Contents{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          BigInt{*big.NewInt(10100)},
						Counter:      BigInt{*big.NewInt(10)},
						GasLimit:     BigInt{*big.NewInt(10100)},
						StorageLimit: BigInt{big.Int{}},
						Kind:         DELEGATIONOP,
						Delegate:     "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					},
				},
				&mockHash,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			gt := testGoTezos(t, tt.input.handler)
			branch, contents, err := gt.UnforgeOperation(tt.input.hexString, tt.input.signed)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.branch, branch)
			assert.Equal(t, tt.want.contents, contents)
		})
	}
}

func Test_checkBoolean(t *testing.T) {
	type input struct {
		hexString string
	}

	type want struct {
		err         bool
		errContains string
		res         bool
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"is boolean",
			input{
				hexString: "ff",
			},
			want{
				false,
				"",
				true,
			},
		},
		{
			"is not boolean",
			input{
				hexString: "00",
			},
			want{
				false,
				"",
				false,
			},
		},
		{
			"is unkown",
			input{
				hexString: "dssdf",
			},
			want{
				true,
				"boolean value is invalid",
				false,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			res, err := checkBoolean(tt.input.hexString)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.res, res)
		})
	}
}

func Test_parseAddress(t *testing.T) {
	type input struct {
		hexString string
	}

	type want struct {
		err         bool
		errContains string
		res         string
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"is successful tz1",
			input{
				hexString: "000008ba0cb2fad622697145cf1665124096d25bc31e",
			},
			want{
				false,
				"",
				"tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
			},
		},
		{
			"is successful KT1",
			input{
				hexString: "018b88e99e66c1c2587f87118449f781cb7d44c9c400",
			},
			want{
				false,
				"",
				"KT1MJZWHKZU7ViybRLsphP3ppiiTc7myP2aj",
			},
		},
		{
			"handles junk",
			input{
				hexString: "e66c1c2587f87118449f781cb7d44c9c400",
			},
			want{
				true,
				"address format not supported",
				"",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			res, err := parseAddress(tt.input.hexString)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.res, res)
		})
	}
}

func Test_removeHexPrefix(t *testing.T) {
	type input struct {
		payload string
		prefix  prefix
	}

	type want struct {
		err         bool
		errContains string
		res         string
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"is successful tz1",
			input{
				"tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
				prefix_tz1,
			},
			want{
				false,
				"",
				"08ba0cb2fad622697145cf1665124096d25bc31e",
			},
		},
		{
			"is successful KT1",
			input{
				"KT1MJZWHKZU7ViybRLsphP3ppiiTc7myP2aj",
				prefix_tz1,
			},
			want{
				false,
				"",
				"8b88e99e66c1c2587f87118449f781cb7d44c9c4",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			res, err := removeHexPrefix(tt.input.payload, tt.input.prefix)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.res, res)
		})
	}
}
