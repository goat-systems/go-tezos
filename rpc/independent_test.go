package rpc_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/completium/go-tezos/v4/rpc"
	"github.com/stretchr/testify/assert"
)

func Test_InjectOperation(t *testing.T) {
	goldenOp := "a732d3520eeaa3de98d78e5e5cb6c85f72204fd46feb9f76853841d4a701add36c0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e00b960000008ba0cb2fad622697145cf1665124096d25bc31e006c0008ba0cb2fad622697145cf1665124096d25bc31ed3e7bd1008d3bb0300b1a803000008ba0cb2fad622697145cf1665124096d25bc31e00"
	goldenHash := []byte(`"oopfasdfadjkfalksj"`)

	type want struct {
		err         bool
		errContains string
		result      string
	}

	cases := []struct {
		name  string
		input http.Handler
		want  want
	}{
		{
			"handles RPC error",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regInjectionOperation, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to inject operation",
				"",
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regInjectionOperation, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to inject operation: failed to parse json",
				"",
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regInjectionOperation, goldenHash}, blankHandler)),
			want{
				false,
				"",
				"oopfasdfadjkfalksj",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input)
			defer server.Close()

			r, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, result, err := r.InjectionOperation(rpc.InjectionOperationInput{
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

	type want struct {
		err         bool
		errContains string
		result      []byte
	}

	cases := []struct {
		name  string
		input http.Handler
		want  want
	}{
		{
			"handles RPC error",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regInjectionBlock, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to inject block",
				goldenRPCError,
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regInjectionBlock, goldenHash}, blankHandler)),
			want{
				false,
				"",
				goldenHash,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input)
			defer server.Close()

			r, err := rpc.New(server.URL)
			assert.Nil(t, err)

			result, err := r.InjectionBlock(rpc.InjectionBlockInput{
				Block: &rpc.Block{},
			})
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.result, result.Body())
		})
	}
}

func Test_Connections(t *testing.T) {
	var goldenConnections rpc.Connections
	json.Unmarshal(readResponse(connections), &goldenConnections)

	type want struct {
		wantErr         bool
		containsErr     string
		wantConnections rpc.Connections
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"handles RPC error",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regConnections, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get network connections",
				rpc.Connections{},
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regConnections, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get network connections: failed to parse json",
				rpc.Connections{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regConnections, readResponse(connections)}, blankHandler)),
			want{
				false,
				"",
				goldenConnections,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			rpc, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, connections, err := rpc.Connections()
			checkErr(t, tt.wantErr, tt.containsErr, err)
			assert.Equal(t, tt.want.wantConnections, connections)
		})
	}
}

func Test_ActiveChains(t *testing.T) {
	type want struct {
		err          bool
		errContains  string
		activeChains rpc.ActiveChains
	}

	cases := []struct {
		name  string
		input http.Handler
		want  want
	}{
		{
			"handles RPC error",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regActiveChains, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get active chains",
				rpc.ActiveChains{},
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regActiveChains, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get active chains: failed to parse json",
				rpc.ActiveChains{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regActiveChains, readResponse(activechains)}, blankHandler)),
			want{
				false,
				"",
				rpc.ActiveChains{
					{
						ChainID: "NetXdQprcVkpaWU",
					},
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input)
			defer server.Close()

			rpc, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, activeChains, err := rpc.ActiveChains()
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.activeChains, activeChains)
		})
	}
}
