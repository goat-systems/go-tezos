package rpc

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
				gtGoldenHTTPMock(blocksHandlerMock(readResponse(rpcerrors), blankHandler)),
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
				gtGoldenHTTPMock(blocksHandlerMock(readResponse(blocks), blankHandler)),
			},
			want{
				false,
				"",
				[][]string{{"BLUdLeoqJtswBAmboRjokR8bM8aiD22FzfM2LVVp5NR8sxLt15r"}},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input.handler)
			defer server.Close()

			rpc, err := New(server.URL)
			assert.Nil(t, err)

			blocks, err := rpc.Blocks(BlocksInput{})
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.blocks, blocks)
		})
	}
}

func Test_ChainID(t *testing.T) {

	var goldenChainID string
	json.Unmarshal(readResponse(chainid), &goldenChainID)

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
				gtGoldenHTTPMock(chainIDHandlerMock(readResponse(rpcerrors), blankHandler)),
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
				gtGoldenHTTPMock(chainIDHandlerMock(readResponse(chainid), blankHandler)),
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

			rpc, err := New(server.URL)
			assert.Nil(t, err)

			chainID, err := rpc.ChainID()
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.chainID, chainID)
		})
	}
}

func Test_Checkpoint(t *testing.T) {
	var goldenCheckpoint Checkpoint
	json.Unmarshal(readResponse(checkpoint), &goldenCheckpoint)

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
				gtGoldenHTTPMock(checkpointHandlerMock(readResponse(rpcerrors), blankHandler)),
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
				gtGoldenHTTPMock(checkpointHandlerMock(readResponse(checkpoint), blankHandler)),
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

			rpc, err := New(server.URL)
			assert.Nil(t, err)

			c, err := rpc.Checkpoint()
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.checkpoint, c)
		})
	}
}

func Test_InvalidBlocks(t *testing.T) {
	//goldenInvalidBlocks := getResponse(invalidblocks).(*InvalidBlock)

	type input struct {
		handler http.Handler
	}

	type want struct {
		err           bool
		errContains   string
		invalidBlocks []InvalidBlock
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"returns rpc error",
			input{
				gtGoldenHTTPMock(invalidBlocksHandlerMock(readResponse(rpcerrors), blankHandler)),
			},
			want{
				true,
				"failed to get invalid blocks",
				[]InvalidBlock{},
			},
		},
		{
			"fails to unmarshal",
			input{
				gtGoldenHTTPMock(invalidBlocksHandlerMock([]byte(`junk`), blankHandler)),
			},
			want{
				true,
				"failed to unmarshal invalid blocks",
				[]InvalidBlock{},
			},
		},
		// {
		// 	"is successful",
		// 	input{
		// 		gtGoldenHTTPMock(invalidBlocksHandlerMock(mockInvalidBlocksResp, blankHandler)),
		// 	},
		// 	want{
		// 		false,
		// 		"",
		// 		&goldenInvalidBlocks,
		// 	},
		// },
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input.handler)
			defer server.Close()

			rpc, err := New(server.URL)
			assert.Nil(t, err)

			blocks, err := rpc.InvalidBlocks()
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.invalidBlocks, blocks)
		})
	}
}

func Test_InvalidBlock(t *testing.T) {
	//goldenInvalidBlock := getResponse(invalidblock).(*InvalidBlock)

	type input struct {
		handler http.Handler
	}

	type want struct {
		err          bool
		errContains  string
		invalidBlock InvalidBlock
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"returns rpc error",
			input{
				gtGoldenHTTPMock(invalidBlocksHandlerMock(readResponse(rpcerrors), blankHandler)),
			},
			want{
				true,
				"failed to get invalid blocks",
				InvalidBlock{},
			},
		},
		{
			"fails to unmarshal",
			input{
				gtGoldenHTTPMock(invalidBlocksHandlerMock([]byte(`junk`), blankHandler)),
			},
			want{
				true,
				"failed to unmarshal invalid blocks",
				InvalidBlock{},
			},
		},
		// {
		// 	"is successful",
		// 	input{
		// 		gtGoldenHTTPMock(invalidBlocksHandlerMock(mockInvalidBlockResp, blankHandler)),
		// 	},
		// 	want{
		// 		false,
		// 		"",
		// 		&goldenInvalidBlock,
		// 	},
		// },
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input.handler)
			defer server.Close()

			rpc, err := New(server.URL)
			assert.Nil(t, err)

			block, err := rpc.InvalidBlock(mockBlockHash)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.invalidBlock, block)
		})
	}
}
