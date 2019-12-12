package gotezos

// import (
// 	"testing"

// 	"gotest.tools/assert"

// 	"github.com/DefinitelyNotAGoat/go-tezos/v2/account"
// 	"github.com/DefinitelyNotAGoat/go-tezos/v2/client"

// 	"github.com/DefinitelyNotAGoat/go-tezos/v2/delegate"
// )

// func Test_CreateBatchPayment(t *testing.T) {
// 	cases := []struct {
// 		payments   []delegate.Payment
// 		wallet     account.Wallet
// 		paymentFee int
// 		gasLimit   int
// 		batchSize  int
// 		tzclient   client.TezosClient
// 	}{
// 		{
// 			payments: []delegate.Payment{
// 				{
// 					Address: "KT1VLb6tJLgmcWTSx7ud4U2n3cNHRBQgxa1t",
// 					Amount:  500000,
// 				},
// 				{
// 					Address: "KT1KwbBJhDvrUkVpu2ZzX1mTTWW8GDwoBRMZ",
// 					Amount:  30823,
// 				},
// 				{
// 					Address: "KT1F3Dwm4j8eZPLNGXUJvb1t3iHouVqdpix8",
// 					Amount:  1423241,
// 				},
// 				{
// 					Address: "KT1CpHMRYfbfMLpLnbisamVEfJo9UZ1KJZu9",
// 					Amount:  2134131423,
// 				},
// 			},
// 			wallet: account.Wallet{
// 				Address: "tz1Qny7jVMGiwRrP9FikRK95jTNbJcffTpx1",
// 				Mnemonic: "normal dash crumble neutral reflect parrot know stairs culture fault check whale flock dog scout",
// 				Seed: []byte{154, 23, 28, 173, 27, 109, 145, 229, 148, 175, 251, 182, 67, 184, 0, 156, 79, 78, 167, 25, 57, 185, 88, 43, 24, 183, 6, 57, 206, 75, 229, 254, 210, 50, 185, 117, 148, 62, 160, 43, 103, 145, 99, 96, 9, 180, 147, 219, 6, 84, 16, 110, 100, 250, 66, 200 77 217 158 70 87 156 200 123] [210 50 185 117 148 62 160 43 103 145 99 96 9 180 147 219 6 84 16 110 100 250 66 200 77 217 158 70 87 156 200 123]}
// 			}
// 		},
// 	}

// 	for _, tc := range cases {
// 		opServ := NewOperationService(&blockServiceMock{}, tc.tzclient)
// 		batchpayments, err := opServ.CreateBatchPayment(tc.payments, tc.wallet, tc.paymentFee, tc.gasLimit, tc.batchSize)
// 		assert.NilError(t, err)
// 	}
// }
