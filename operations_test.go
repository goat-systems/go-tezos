package gotezos

import (
	"math/big"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ForgeOperation(t *testing.T) {
	type input struct {
		handler  http.Handler
		contents []Contents
		branch   string
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
			"is successful transaction",
			input{
				gtGoldenHTTPMock(blankHandler),
				[]Contents{
					Contents{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          BigInt{*big.NewInt(10100)},
						Counter:      BigInt{*big.NewInt(10)},
						GasLimit:     BigInt{*big.NewInt(10100)},
						StorageLimit: BigInt{big.Int{}},
						Amount:       BigInt{*big.NewInt(12345)},
						Destination:  "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Kind:         "transaction",
					},
					Contents{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          BigInt{*big.NewInt(34567123)},
						Counter:      BigInt{*big.NewInt(8)},
						GasLimit:     BigInt{*big.NewInt(56787)},
						StorageLimit: BigInt{big.Int{}},
						Amount:       BigInt{*big.NewInt(54321)},
						Destination:  "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Kind:         "transaction",
					},
				},
				"BLyvCRkxuTXkx1KeGvrcEXiPYj4p1tFxzvFDhoHE7SFKtmP1rbk",
			},
			want{
				false,
				"",
				"a732d3520eeaa3de98d78e5e5cb6c85f72204fd46feb9f76853841d4a701add36c0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e00b960000008ba0cb2fad622697145cf1665124096d25bc31e006c0008ba0cb2fad622697145cf1665124096d25bc31ed3e7bd1008d3bb0300b1a803000008ba0cb2fad622697145cf1665124096d25bc31e00",
			},
		},
		{
			"is successful reveal",
			input{
				gtGoldenHTTPMock(blankHandler),
				[]Contents{
					Contents{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          BigInt{*big.NewInt(10100)},
						Counter:      BigInt{*big.NewInt(10)},
						GasLimit:     BigInt{*big.NewInt(10100)},
						StorageLimit: BigInt{big.Int{}},
						Phk:          "edpktnktxAzmXPD9XVNqAvdCFb76vxzQtkbVkSEtXcTz33QZQdb4JQ",
						Kind:         string(REVEAL),
					},
					Contents{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          BigInt{*big.NewInt(34567123)},
						Counter:      BigInt{*big.NewInt(8)},
						GasLimit:     BigInt{*big.NewInt(56787)},
						StorageLimit: BigInt{big.Int{}},
						Phk:          "edpktnktxAzmXPD9XVNqAvdCFb76vxzQtkbVkSEtXcTz33QZQdb4JQ",
						Kind:         string(REVEAL),
					},
				},
				"BLyvCRkxuTXkx1KeGvrcEXiPYj4p1tFxzvFDhoHE7SFKtmP1rbk",
			},
			want{
				false,
				"",
				"a732d3520eeaa3de98d78e5e5cb6c85f72204fd46feb9f76853841d4a701add36b0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e0000136083897bc97879c53e3e7855838fbbc87303ddd376080fc3d3e136b55d028b6b0008ba0cb2fad622697145cf1665124096d25bc31ed3e7bd1008d3bb030000136083897bc97879c53e3e7855838fbbc87303ddd376080fc3d3e136b55d028b",
			},
		},
		{
			"is successful origination",
			input{
				gtGoldenHTTPMock(blankHandler),
				[]Contents{
					Contents{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          BigInt{*big.NewInt(10100)},
						Counter:      BigInt{*big.NewInt(10)},
						GasLimit:     BigInt{*big.NewInt(10100)},
						StorageLimit: BigInt{big.Int{}},
						Kind:         string(ORIGINATION),
						Balance:      BigInt{*big.NewInt(328763282)},
						Delegate:     "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					},
				},
				"BLyvCRkxuTXkx1KeGvrcEXiPYj4p1tFxzvFDhoHE7SFKtmP1rbk",
			},
			want{
				false,
				"",
				"a732d3520eeaa3de98d78e5e5cb6c85f72204fd46feb9f76853841d4a701add36d0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e00928fe29c01ff0008ba0cb2fad622697145cf1665124096d25bc31e000000c602000000c105000764085e036c055f036d0000000325646f046c000000082564656661756c740501035d050202000000950200000012020000000d03210316051f02000000020317072e020000006a0743036a00000313020000001e020000000403190325072c020000000002000000090200000004034f0327020000000b051f02000000020321034c031e03540348020000001e020000000403190325072c020000000002000000090200000004034f0327034f0326034202000000080320053d036d03420000001a0a000000150008ba0cb2fad622697145cf1665124096d25bc31e",
			},
		},
		{
			"is successful delegation",
			input{
				gtGoldenHTTPMock(blankHandler),
				[]Contents{
					Contents{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          BigInt{*big.NewInt(10100)},
						Counter:      BigInt{*big.NewInt(10)},
						GasLimit:     BigInt{*big.NewInt(10100)},
						StorageLimit: BigInt{big.Int{}},
						Kind:         string(ORIGINATION),
						Balance:      BigInt{*big.NewInt(328763282)},
						Delegate:     "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					},
				},
				"BLyvCRkxuTXkx1KeGvrcEXiPYj4p1tFxzvFDhoHE7SFKtmP1rbk",
			},
			want{
				false,
				"",
				"a732d3520eeaa3de98d78e5e5cb6c85f72204fd46feb9f76853841d4a701add36e0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e00ff0008ba0cb2fad622697145cf1665124096d25bc31e",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			gt := testGoTezos(t, tt.input.handler)
			operation, err := gt.ForgeOperation(tt.input.branch, tt.input.contents...)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.operation, operation)
		})
	}
}

func Test_UnforgeOperation(t *testing.T) {
	type input struct {
		handler   http.Handler
		hexString string
		signed    bool
	}

	type want struct {
		err         bool
		errContains string
		contents    []Contents
		branch      string
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
				"a732d3520eeaa3de98d78e5e5cb6c85f72204fd46feb9f76853841d4a701add36c0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e00b960000008ba0cb2fad622697145cf1665124096d25bc31e006c0008ba0cb2fad622697145cf1665124096d25bc31ed3e7bd1008d3bb0300b1a803000008ba0cb2fad622697145cf1665124096d25bc31e00",
				false,
			},
			want{
				false,
				"",
				[]Contents{
					Contents{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          BigInt{*big.NewInt(10100)},
						Counter:      BigInt{*big.NewInt(10)},
						GasLimit:     BigInt{*big.NewInt(10100)},
						StorageLimit: BigInt{big.Int{}},
						Amount:       BigInt{*big.NewInt(12345)},
						Destination:  "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Kind:         "transaction",
					},
					Contents{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          BigInt{*big.NewInt(34567123)},
						Counter:      BigInt{*big.NewInt(8)},
						GasLimit:     BigInt{*big.NewInt(56787)},
						StorageLimit: BigInt{big.Int{}},
						Amount:       BigInt{*big.NewInt(54321)},
						Destination:  "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Kind:         "transaction",
					},
				},
				"BLyvCRkxuTXkx1KeGvrcEXiPYj4p1tFxzvFDhoHE7SFKtmP1rbk",
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
				[]Contents{
					Contents{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          BigInt{*big.NewInt(10100)},
						Counter:      BigInt{*big.NewInt(10)},
						GasLimit:     BigInt{*big.NewInt(10100)},
						StorageLimit: BigInt{big.Int{}},
						Phk:          "edpktnktxAzmXPD9XVNqAvdCFb76vxzQtkbVkSEtXcTz33QZQdb4JQ",
						Kind:         string(REVEAL),
					},
					Contents{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          BigInt{*big.NewInt(34567123)},
						Counter:      BigInt{*big.NewInt(8)},
						GasLimit:     BigInt{*big.NewInt(56787)},
						StorageLimit: BigInt{big.Int{}},
						Phk:          "edpktnktxAzmXPD9XVNqAvdCFb76vxzQtkbVkSEtXcTz33QZQdb4JQ",
						Kind:         string(REVEAL),
					},
				},
				"BLyvCRkxuTXkx1KeGvrcEXiPYj4p1tFxzvFDhoHE7SFKtmP1rbk",
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
				[]Contents{
					Contents{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          BigInt{*big.NewInt(10100)},
						Counter:      BigInt{*big.NewInt(10)},
						GasLimit:     BigInt{*big.NewInt(10100)},
						StorageLimit: BigInt{big.Int{}},
						Kind:         string(ORIGINATION),
						Balance:      BigInt{*big.NewInt(328763282)},
						Delegate:     "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					},
				},
				"BLyvCRkxuTXkx1KeGvrcEXiPYj4p1tFxzvFDhoHE7SFKtmP1rbk",
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
