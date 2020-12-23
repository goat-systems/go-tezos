package rpc_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goat-systems/go-tezos/v4/rpc"
	"github.com/stretchr/testify/assert"
)

func Test_Version(t *testing.T) {
	goldenVersion := getResponse(version).(rpc.Version)

	type want struct {
		wantErr     bool
		containsErr string
		wantVersion rpc.Version
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"returns rpc error",
			gtGoldenHTTPMock(versionsHandlerMock(readResponse(rpcerrors), blankHandler)),
			want{
				true,
				"could not get network version",
				rpc.Version{},
			},
		},
		{
			"fails to unmarshal",
			gtGoldenHTTPMock(versionsHandlerMock([]byte(`junk`), blankHandler)),
			want{
				true,
				"could not unmarshal network version",
				rpc.Version{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(versionsHandlerMock(readResponse(version), blankHandler)),
			want{
				false,
				"",
				goldenVersion,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			rpc, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, version, err := rpc.Version()
			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tt.want.containsErr)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, tt.want.wantVersion, version)
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
			"returns rpc error",
			gtGoldenHTTPMock(connectionsHandlerMock(readResponse(rpcerrors), blankHandler)),
			want{
				true,
				"could not get network connections",
				rpc.Connections{},
			},
		},
		{
			"fails to unmarshal",
			gtGoldenHTTPMock(connectionsHandlerMock([]byte(`junk`), blankHandler)),
			want{
				true,
				"could not unmarshal network connections",
				rpc.Connections{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(connectionsHandlerMock(readResponse(connections), blankHandler)),
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
			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tt.want.containsErr)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, tt.want.wantConnections, connections)
		})
	}
}

func Test_Bootsrap(t *testing.T) {
	var goldenBootstrap rpc.Bootstrap
	json.Unmarshal(readResponse(bootstrap), &goldenBootstrap)

	type want struct {
		wantErr       bool
		containsErr   string
		wantBootstrap rpc.Bootstrap
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"returns rpc error",
			gtGoldenHTTPMock(bootstrapHandlerMock(readResponse(rpcerrors), blankHandler)),
			want{
				true,
				"could not get bootstrap",
				rpc.Bootstrap{},
			},
		},
		{
			"fails to unmarshal",
			gtGoldenHTTPMock(bootstrapHandlerMock([]byte(`junk`), blankHandler)),
			want{
				true,
				"could not unmarshal bootstrap",
				rpc.Bootstrap{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(bootstrapHandlerMock(readResponse(bootstrap), blankHandler)),
			want{
				false,
				"",
				goldenBootstrap,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			rpc, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, bootstrap, err := rpc.Bootstrap()
			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tt.want.containsErr)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, tt.want.wantBootstrap, bootstrap)
		})
	}
}

func Test_Commit(t *testing.T) {
	var goldenCommit string
	json.Unmarshal(readResponse(commit), &goldenCommit)

	type want struct {
		wantErr     bool
		containsErr string
		wantCommit  string
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"returns rpc error",
			gtGoldenHTTPMock(commitHandlerMock(readResponse(rpcerrors), blankHandler)),
			want{
				true,
				"could not get commit",
				"",
			},
		},
		{
			"fails to unmarshal",
			gtGoldenHTTPMock(commitHandlerMock([]byte(`junk`), blankHandler)),
			want{
				true,
				"could unmarshal commit",
				"",
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(commitHandlerMock(readResponse(commit), blankHandler)),
			want{
				false,
				"",
				goldenCommit,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			rpc, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, commit, err := rpc.Commit()
			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tt.want.containsErr)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, tt.want.wantCommit, commit)
		})
	}
}

func Test_Cycle(t *testing.T) {
	type input struct {
		handler http.Handler
		cycle   int
	}

	type want struct {
		err         bool
		errContains string
		cycle       rpc.Cycle
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"failed to get head block",
			input{
				gtGoldenHTTPMock(newBlockMock().handler([]byte(`not_block_data`), blankHandler)),
				10,
			},
			want{
				true,
				"could not get cycle '10': could not get head block",
				rpc.Cycle{},
			},
		},
		{
			"failed to get cycle because cycle is in the future",
			input{
				gtGoldenHTTPMock(newBlockMock().handler(readResponse(block), blankHandler)),
				300,
			},
			want{
				true,
				"request is in the future",
				rpc.Cycle{},
			},
		},
		{
			"failed to get block less than cycle",
			input{
				gtGoldenHTTPMock(
					newBlockMock().handler(
						readResponse(block),
						newBlockMock().handler(
							[]byte(`not_block_data`),
							blankHandler,
						),
					),
				),
				2,
			},
			want{
				true,
				"could not get block",
				rpc.Cycle{},
			},
		},
		{
			"failed to unmarshal cycle",
			input{
				gtGoldenHTTPMock(
					cycleHandlerMock(
						[]byte(`bad_cycle_data`),
						newBlockMock().handler(
							readResponse(block),
							newBlockMock().handler(
								readResponse(block),
								blankHandler,
							),
						),
					),
				),
				2,
			},
			want{
				true,
				"could not unmarshal at cycle hash",
				rpc.Cycle{},
			},
		},
		{
			"failed to get cycle block level",
			input{
				gtGoldenHTTPMock(
					cycleHandlerMock(
						readResponse(cycle),
						newBlockMock().handler(
							readResponse(block),
							newBlockMock().handler(
								readResponse(block),
								newBlockMock().handler(
									[]byte(`not_block_data`),
									blankHandler,
								),
							),
						),
					),
				),
				2,
			},
			want{
				true,
				"could not get block",
				rpc.Cycle{
					LastRoll:     []string{},
					Nonces:       []string{},
					RandomSeed:   "04dca5c197fc2e18309b60844148c55fc7ccdbcb498bd57acd4ac29f16e22846",
					RollSnapshot: 4,
				},
			},
		},
		{
			"is successful",
			input{
				gtGoldenHTTPMock(mockCycleSuccessful(blankHandler)),
				2,
			},
			want{
				false,
				"",
				rpc.Cycle{
					LastRoll:     []string{},
					Nonces:       []string{},
					RandomSeed:   "04dca5c197fc2e18309b60844148c55fc7ccdbcb498bd57acd4ac29f16e22846",
					RollSnapshot: 4,
					BlockHash:    "BLBL72xDLHf4ffKu8NZhYnqy21DECDkZ3Vpjw7oZJDhbgySzwFT",
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input.handler)
			defer server.Close()

			rpc, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, cycle, err := rpc.Cycle(tt.input.cycle)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.cycle, cycle)
		})
	}
}

func Test_ActiveChains(t *testing.T) {
	type input struct {
		handler http.Handler
	}

	type want struct {
		err          bool
		errContains  string
		activeChains rpc.ActiveChains
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"returns rpc error",
			input{
				gtGoldenHTTPMock(activeChainsHandlerMock(readResponse(rpcerrors), blankHandler)),
			},
			want{
				true,
				"failed to get active chains",
				nil,
			},
		},
		{
			"fails to unmarshal",
			input{
				gtGoldenHTTPMock(activeChainsHandlerMock([]byte(`junk`), blankHandler)),
			},
			want{
				true,
				"failed to unmarshal active chains",
				nil,
			},
		},
		{
			"is successful",
			input{
				gtGoldenHTTPMock(activeChainsHandlerMock(readResponse(activechains), blankHandler)),
			},
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
			server := httptest.NewServer(tt.input.handler)
			defer server.Close()

			rpc, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, activeChains, err := rpc.ActiveChains()
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.activeChains, activeChains)
		})
	}
}

func mockCycleSuccessful(next http.Handler) http.Handler {
	var blockmock blockHandlerMock
	var oldHTTPBlock blockHandlerMock
	var blockAtLevel blockHandlerMock
	return cycleHandlerMock(
		readResponse(cycle),
		blockmock.handler(
			readResponse(block),
			oldHTTPBlock.handler(
				readResponse(block),
				blockAtLevel.handler(
					readResponse(block),
					next,
				),
			),
		),
	)
}

func mockCycleFailed(next http.Handler) http.Handler {
	var blockmock blockHandlerMock
	var oldHTTPBlock blockHandlerMock
	return cycleHandlerMock(
		[]byte(`bad_cycle_data`),
		blockmock.handler(
			readResponse(block),
			oldHTTPBlock.handler(
				readResponse(block),
				blankHandler,
			),
		),
	)
}
