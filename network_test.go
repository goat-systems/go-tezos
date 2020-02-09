package gotezos

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Version(t *testing.T) {

	var goldenVersion Version
	json.Unmarshal(mockVersionResp, &goldenVersion)

	type want struct {
		wantErr     bool
		containsErr string
		wantVersion *Version
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"returns rpc error",
			gtGoldenHTTPMock(versionsHandlerMock(mockRPCErrorResp, blankHandler)),
			want{
				true,
				"could not get network version",
				&Version{},
			},
		},
		{
			"fails to unmarshal",
			gtGoldenHTTPMock(versionsHandlerMock([]byte(`junk`), blankHandler)),
			want{
				true,
				"could not unmarshal network version",
				&Version{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(versionsHandlerMock(mockVersionResp, blankHandler)),
			want{
				false,
				"",
				&goldenVersion,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			gt, err := New(server.URL)
			assert.Nil(t, err)

			version, err := gt.Version()
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

func Test_Constants(t *testing.T) {

	var goldenConstants Constants
	json.Unmarshal(mockConstantsResp, &goldenConstants)

	type want struct {
		wantErr       bool
		containsErr   string
		wantConstants *Constants
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"returns rpc error",
			gtGoldenHTTPMock(newConstantsMock().handler(mockRPCErrorResp, blankHandler)),
			want{
				true,
				"could not get network constants",
				&Constants{},
			},
		},
		{
			"fails to unmarshal",
			gtGoldenHTTPMock(newConstantsMock().handler([]byte(`junk`), blankHandler)),
			want{
				true,
				"could not unmarshal network constants",
				&Constants{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(newConstantsMock().handler(mockConstantsResp, blankHandler)),
			want{
				false,
				"",
				&goldenConstants,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			gt, err := New(server.URL)
			assert.Nil(t, err)

			constants, err := gt.Constants("BLzGD63HA4RP8Fh5xEtvdQSMKa2WzJMZjQPNVUc4Rqy8Lh5BEY1")
			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tt.want.containsErr)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, tt.want.wantConstants, constants)
		})
	}
}

func Test_Connections(t *testing.T) {

	var goldenConnections Connections
	json.Unmarshal(mockConnectionsResp, &goldenConnections)

	type want struct {
		wantErr         bool
		containsErr     string
		wantConnections *Connections
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"returns rpc error",
			gtGoldenHTTPMock(connectionsHandlerMock(mockRPCErrorResp, blankHandler)),
			want{
				true,
				"could not get network connections",
				&Connections{},
			},
		},
		{
			"fails to unmarshal",
			gtGoldenHTTPMock(connectionsHandlerMock([]byte(`junk`), blankHandler)),
			want{
				true,
				"could not unmarshal network connections",
				&Connections{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(connectionsHandlerMock(mockConnectionsResp, blankHandler)),
			want{
				false,
				"",
				&goldenConnections,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			gt, err := New(server.URL)
			assert.Nil(t, err)

			connections, err := gt.Connections()
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
	var goldenBootstrap Bootstrap
	json.Unmarshal(mockBootstrapResp, &goldenBootstrap)

	type want struct {
		wantErr       bool
		containsErr   string
		wantBootstrap *Bootstrap
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"returns rpc error",
			gtGoldenHTTPMock(bootstrapHandlerMock(mockRPCErrorResp, blankHandler)),
			want{
				true,
				"could not get bootstrap",
				&Bootstrap{},
			},
		},
		{
			"fails to unmarshal",
			gtGoldenHTTPMock(bootstrapHandlerMock([]byte(`junk`), blankHandler)),
			want{
				true,
				"could not unmarshal bootstrap",
				&Bootstrap{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(bootstrapHandlerMock(mockBootstrapResp, blankHandler)),
			want{
				false,
				"",
				&goldenBootstrap,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			gt, err := New(server.URL)
			assert.Nil(t, err)

			bootstrap, err := gt.Bootstrap()
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
	json.Unmarshal(mockCommitResp, &goldenCommit)

	type want struct {
		wantErr     bool
		containsErr string
		wantCommit  *string
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"returns rpc error",
			gtGoldenHTTPMock(commitHandlerMock(mockRPCErrorResp, blankHandler)),
			want{
				true,
				"could not get commit",
				nil,
			},
		},
		{
			"fails to unmarshal",
			gtGoldenHTTPMock(commitHandlerMock([]byte(`junk`), blankHandler)),
			want{
				true,
				"could unmarshal commit",
				nil,
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(commitHandlerMock(mockCommitResp, blankHandler)),
			want{
				false,
				"",
				&goldenCommit,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			gt, err := New(server.URL)
			assert.Nil(t, err)

			commit, err := gt.Commit()
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
		cycle       *Cycle
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
				&Cycle{},
			},
		},
		{
			"failed to get cycle because cycle is in the future",
			input{
				gtGoldenHTTPMock(newBlockMock().handler(mockBlockResp, blankHandler)),
				300,
			},
			want{
				true,
				"request is in the future",
				&Cycle{},
			},
		},
		{
			"failed to get block less than cycle",
			input{
				gtGoldenHTTPMock(
					newBlockMock().handler(
						mockBlockResp,
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
				&Cycle{},
			},
		},
		{
			"failed to unmarshal cycle",
			input{
				gtGoldenHTTPMock(
					cycleHandlerMock(
						[]byte(`bad_cycle_data`),
						newBlockMock().handler(
							mockBlockResp,
							newBlockMock().handler(
								mockBlockResp,
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
				&Cycle{},
			},
		},
		{
			"failed to get cycle block level",
			input{
				gtGoldenHTTPMock(
					cycleHandlerMock(
						mockCycleResp,
						newBlockMock().handler(
							mockBlockResp,
							newBlockMock().handler(
								mockBlockResp,
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
				&Cycle{
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
				&Cycle{
					RandomSeed:   "04dca5c197fc2e18309b60844148c55fc7ccdbcb498bd57acd4ac29f16e22846",
					RollSnapshot: 4,
					BlockHash:    "BLzGD63HA4RP8Fh5xEtvdQSMKa2WzJMZjQPNVUc4Rqy8Lh5BEY1",
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input.handler)
			defer server.Close()

			gt, err := New(server.URL)
			assert.Nil(t, err)

			cycle, err := gt.Cycle(tt.input.cycle)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.cycle, cycle)
		})
	}
}

func mockCycleSuccessful(next http.Handler) http.Handler {
	var blockmock blockHandlerMock
	var oldHTTPBlock blockHandlerMock
	var blockAtLevel blockHandlerMock
	return cycleHandlerMock(
		mockCycleResp,
		blockmock.handler(
			mockBlockResp,
			oldHTTPBlock.handler(
				mockBlockResp,
				blockAtLevel.handler(
					mockBlockResp,
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
			mockBlockResp,
			oldHTTPBlock.handler(
				mockBlockResp,
				blankHandler,
			),
		),
	)
}
