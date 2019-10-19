package block

import (
	"encoding/json"
	"testing"

	"gotest.tools/assert"

	tezc "github.com/DefinitelyNotAGoat/go-tezos/v2/client"
)

func Test_Get(t *testing.T) {
	cases := []struct {
		id       interface{}
		want     []byte
		wantErr  bool
		tzclient tezc.TezosClient
	}{
		{
			id:      "BLTGSUUjDpaHe7BYZa1zsrccJ7skurNiHZ1mpCz3cak9GnDfRoT",
			want:    goldenBlock,
			wantErr: false,
			tzclient: &client{
				ReturnBody: goldenBlock,
			},
		},
		{
			id:      "BLTGSUUjDpaHe7BYZa1zsrccJ7skurNiHZ1mpCz3cak9GnDfRoT",
			want:    goldenBlock,
			wantErr: true,
			tzclient: &client{
				ReturnBody: []byte("malformed response"),
			},
		},
	}

	for _, tc := range cases {
		blockService := NewBlockService(tc.tzclient)

		block, err := blockService.Get(tc.id)
		if !tc.wantErr {
			assert.NilError(t, err)
			blockwant := Block{}
			blockwant, err = blockwant.unmarshalJSON(goldenBlock)
			assert.NilError(t, err)

			jsonHave, _ := json.Marshal(block)
			jsonWant, _ := json.Marshal(blockwant)

			assert.Equal(t, string(jsonHave), string(jsonWant))
		} else {
			assert.Assert(t, err != nil)
		}
	}
}

func Test_GetHead(t *testing.T) {
	cases := []struct {
		want     []byte
		tzclient tezc.TezosClient
		wantErr  bool
	}{
		{
			want: goldenBlock,
			tzclient: &client{
				ReturnBody: goldenBlock,
			},
			wantErr: false,
		},
		{
			want: goldenBlock,
			tzclient: &client{
				ReturnBody: []byte("malformed response"),
			},
			wantErr: true,
		},
	}

	for _, tc := range cases {
		blockService := NewBlockService(tc.tzclient)

		block, err := blockService.GetHead()
		if !tc.wantErr {
			assert.NilError(t, err)
			blockwant := Block{}
			blockwant, err = blockwant.unmarshalJSON(goldenBlock)
			assert.NilError(t, err)

			jsonHave, _ := json.Marshal(block)
			jsonWant, _ := json.Marshal(blockwant)

			assert.Equal(t, string(jsonHave), string(jsonWant))
		} else {
			assert.Assert(t, err != nil)
		}

	}
}
