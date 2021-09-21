//go:build integration
// +build integration

package keys

import (
	"os"
	"strconv"
	"testing"

	"github.com/completium/go-tezos/v4/forge"
	"github.com/completium/go-tezos/v4/internal/testutils"
	"github.com/completium/go-tezos/v4/rpc"
)

func Test_OperationWithKey(t *testing.T) {
	type input struct {
		sk   string
		kind ECKind
	}

	cases := []struct {
		name    string
		input   input
		wantErr bool
	}{
		{
			"is successful Ed25519",
			input{
				sk:   "edsk2oJWw5CX7Fh3g8QDqtK9CmrvRDDDSeHAPvWPnm7CwD3RfQ1KbK",
				kind: Ed25519,
			},
			false,
		},
		{
			"is successful Secp256k1",
			input{
				sk:   "spsk1WCtWP1fEc4RaE63YK6oUEmbjLK2aTe7LevYSb9Z3zDdtq58wS",
				kind: Secp256k1,
			},
			false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			rpchost := os.Getenv("GOTEZOS_TEST_RPC_HOST")
			r, _ := rpc.New(rpchost)

			key, err := FromBase58(tt.input.sk, tt.input.kind)
			testutils.CheckErr(t, tt.wantErr, "", err)

			head, _ := r.Head()

			counter, _ := r.Counter(rpc.CounterInput{
				Blockhash: head.Hash,
				Address:   key.PubKey.address,
			})

			transaction := rpc.Transaction{
				Kind:         rpc.TRANSACTION,
				Source:       key.PubKey.address,
				Fee:          "2941",
				Counter:      strconv.Itoa((counter + 1)),
				GasLimit:     "26283",
				Amount:       "1",
				StorageLimit: "0",
				Destination:  "tz1RomaiWJV3NFDZWTMVR2aEeHknsn3iF5Gi",
			}

			op, err := forge.Encode(head.Hash, transaction.ToContent())
			testutils.CheckErr(t, tt.wantErr, "", err)

			sig, err := key.SignHex(op)
			testutils.CheckErr(t, tt.wantErr, "", err)

			_, err = r.PreapplyOperations(rpc.PreapplyOperationsInput{
				Blockhash: head.Hash,
				Operations: []rpc.Operations{
					{
						Protocol:  head.Protocol,
						Branch:    head.Hash,
						Contents:  rpc.Contents{transaction.ToContent()},
						Signature: sig.ToBase58(),
					},
				},
			})
			testutils.CheckErr(t, tt.wantErr, "", err)
		})
	}
}
