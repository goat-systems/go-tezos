package gotezos

import (
	"math/big"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// func Test_ForgeOperation(t *testing.T) {
// 	type input struct {
// 		handler  http.Handler
// 		contents []Contents
// 		branch   string
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
// 			"is successful",
// 			input{
// 				gtGoldenHTTPMock(blankHandler),
// 				[]Contents{
// 					Contents{
// 						Kind:         string(TRANSACTION),
// 						Fee:          BigInt{*big.NewInt(24000)},
// 						GasLimit:     BigInt{*big.NewInt(3000000)},
// 						StorageLimit: BigInt{*big.NewInt(82379423)},
// 						Source:       "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
// 						Destination:  "tz1W3HW533csCBLor4NPtU79R2TT2sbKfJDH",
// 					},
// 				},
// 				"BMUMXmpCL96m6zhMVD19TsaWsJUJYFoiLKw87n7GiPdetNbGvrK",
// 			},
// 			want{
// 				false,
// 				"",
// 				"e7c48a6630a00e276292bf5d0df311da1b72442b3e7982cb963904e8cd9d1ddf773e4d786c4b04ad1e57c2f13b61b3d2c95b3073d961a4132b93fb36890c00bb100c008d0b7109f0850a4027072172ff3890a9b72ad1a2271cde3593301e08924de5a137c00",
// 			},
// 		},
// 	}

// 	for _, tt := range cases {
// 		t.Run(tt.name, func(t *testing.T) {
// 			gt := testGoTezos(t, tt.input.handler)
// 			operation, err := gt.ForgeOperation(tt.input.branch, tt.input.contents...)
// 			checkErr(t, tt.want.err, tt.want.errContains, err)
// 			assert.Equal(t, tt.want.operation, operation)
// 		})
// 	}
// }

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
			"is successful",
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
