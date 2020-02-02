package gotezos

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Blocks(t *testing.T) {
	type input struct {
		handler http.Handler
	}

	type want struct {
		err         bool
		errContains string
		blocks      [][]string
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"returns rpc error",
			input{
				gtGoldenHTTPMock(blocksHandlerMock(mockRPCError, blankHandler)),
			},
			want{
				true,
				"failed to get blocks",
				[][]string{},
			},
		},
		{
			"fails to unmarshal",
			input{
				gtGoldenHTTPMock(blocksHandlerMock([]byte(`junk`), blankHandler)),
			},
			want{
				true,
				"failed to unmarshal blocks",
				[][]string{},
			},
		},
		{
			"is successful",
			input{
				gtGoldenHTTPMock(blocksHandlerMock(mockBlocks, blankHandler)),
			},
			want{
				false,
				"",
				[][]string{[]string{"BLUdLeoqJtswBAmboRjokR8bM8aiD22FzfM2LVVp5NR8sxLt15r"}},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input.handler)
			defer server.Close()

			gt, err := New(server.URL)
			assert.Nil(t, err)

			blocks, err := gt.Blocks()
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.blocks, blocks)
		})
	}
}

func Test_ChainID(t *testing.T) {

	var goldenChainID string
	json.Unmarshal(mockChainID, &goldenChainID)

	type input struct {
		handler http.Handler
	}

	type want struct {
		err         bool
		errContains string
		chainID     string
	}

	cases := []struct {
		name  string
		input input
		want
	}{
		{
			"returns rpc error",
			input{
				gtGoldenHTTPMock(chainIDHandlerMock(mockRPCError, blankHandler)),
			},
			want{
				true,
				"failed to get chain id",
				"",
			},
		},
		{
			"fails to unmarshal",
			input{
				gtGoldenHTTPMock(chainIDHandlerMock([]byte(`junk`), blankHandler)),
			},
			want{
				true,
				"failed to unmarshal chain id",
				"",
			},
		},
		{
			"is successful",
			input{
				gtGoldenHTTPMock(chainIDHandlerMock(mockChainID, blankHandler)),
			},
			want{
				false,
				"",
				goldenChainID,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input.handler)
			defer server.Close()

			gt, err := New(server.URL)
			assert.Nil(t, err)

			chainID, err := gt.ChainID()
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.chainID, chainID)
		})
	}
}

func Test_Checkpoint(t *testing.T) {
	var goldenCheckpoint Checkpoint
	json.Unmarshal(mockCheckpoint, &goldenCheckpoint)

	type input struct {
		handler http.Handler
	}

	type want struct {
		err         bool
		errContains string
		checkpoint  Checkpoint
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"returns rpc error",
			input{
				gtGoldenHTTPMock(checkpointHandlerMock(mockRPCError, blankHandler)),
			},
			want{
				true,
				"failed to get checkpoint",
				Checkpoint{},
			},
		},
		{
			"fails to unmarshal",
			input{
				gtGoldenHTTPMock(checkpointHandlerMock([]byte(`junk`), blankHandler)),
			},
			want{
				true,
				"failed to unmarshal checkpoint",
				Checkpoint{},
			},
		},
		{
			"is successful",
			input{
				gtGoldenHTTPMock(checkpointHandlerMock(mockCheckpoint, blankHandler)),
			},
			want{
				false,
				"",
				goldenCheckpoint,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input.handler)
			defer server.Close()

			gt, err := New(server.URL)
			assert.Nil(t, err)

			c, err := gt.Checkpoint()
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.checkpoint, c)
		})
	}
}
