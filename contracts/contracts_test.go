package contracts

import (
	"testing"

	tzc "github.com/DefinitelyNotAGoat/go-tezos/client"
	"gotest.tools/assert"
)

func Test_GetStorage(t *testing.T) {
	cases := []struct {
		input    string
		tzclient tzc.TezosClient
		wantErr  bool
	}{
		{
			input: "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
			tzclient: &clientMock{
				ReturnBody: goldenStorage,
			},
			wantErr: false,
		},
	}

	for _, tc := range cases {

		contractService := NewContractService(tc.tzclient)

		storage, err := contractService.GetStorage(tc.input)
		if !tc.wantErr {
			assert.NilError(t, err)
			assert.Equal(t, string(goldenStorage), string(storage))
		} else {
			assert.Assert(t, err != nil)
		}
	}
}
