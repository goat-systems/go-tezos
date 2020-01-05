package gotezos

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Versions(t *testing.T) {

	var goldenVersions Versions
	json.Unmarshal(mockVersions, &goldenVersions)

	type want struct {
		wantErr      bool
		containsErr  string
		wantVersions Versions
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"returns rpc error",
			gtGoldenHTTPMock(versionsHandlerMock(mockRPCError, blankHandler)),
			want{
				true,
				"could not get network versions",
				Versions{},
			},
		},
		{
			"fails to unmarshal",
			gtGoldenHTTPMock(versionsHandlerMock([]byte(`junk`), blankHandler)),
			want{
				true,
				"could not unmarshal network versions",
				Versions{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(versionsHandlerMock(mockVersions, blankHandler)),
			want{
				false,
				"",
				goldenVersions,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			gt, err := New(server.URL)
			assert.Nil(t, err)

			versions, err := gt.Versions()
			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tt.want.containsErr)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, tt.want.wantVersions, versions)
		})
	}
}

func Test_Constants(t *testing.T) {

	var goldenConstants Constants
	json.Unmarshal(mockConstants, &goldenConstants)

	type want struct {
		wantErr       bool
		containsErr   string
		wantConstants Constants
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"returns rpc error",
			gtGoldenHTTPMock(newConstantsMock().handler(mockRPCError, blankHandler)),
			want{
				true,
				"could not get network constants",
				Constants{},
			},
		},
		{
			"fails to unmarshal",
			gtGoldenHTTPMock(newConstantsMock().handler([]byte(`junk`), blankHandler)),
			want{
				true,
				"could not unmarshal network constants",
				Constants{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(newConstantsMock().handler(mockConstants, blankHandler)),
			want{
				false,
				"",
				goldenConstants,
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

func Test_ChainID(t *testing.T) {

	var goldenChainID string
	json.Unmarshal(mockChainID, &goldenChainID)

	type want struct {
		wantErr     bool
		containsErr string
		wantChainID string
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"returns rpc error",
			gtGoldenHTTPMock(chainIDHandlerMock(mockRPCError, blankHandler)),
			want{
				true,
				"could not get chain id",
				"",
			},
		},
		{
			"fails to unmarshal",
			gtGoldenHTTPMock(chainIDHandlerMock([]byte(`junk`), blankHandler)),
			want{
				true,
				"could unmarshal chain id",
				"",
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(chainIDHandlerMock(mockChainID, blankHandler)),
			want{
				false,
				"",
				goldenChainID,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			gt, err := New(server.URL)
			assert.Nil(t, err)

			chainID, err := gt.ChainID()
			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tt.want.containsErr)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, tt.want.wantChainID, chainID)
		})
	}
}

func Test_Connections(t *testing.T) {

	var goldenConnections Connections
	json.Unmarshal(mockConnections, &goldenConnections)

	type want struct {
		wantErr         bool
		containsErr     string
		wantConnections Connections
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"returns rpc error",
			gtGoldenHTTPMock(connectionsHandlerMock(mockRPCError, blankHandler)),
			want{
				true,
				"could not get network connections",
				Connections{},
			},
		},
		{
			"fails to unmarshal",
			gtGoldenHTTPMock(connectionsHandlerMock([]byte(`junk`), blankHandler)),
			want{
				true,
				"could not unmarshal network connections",
				Connections{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(connectionsHandlerMock(mockConnections, blankHandler)),
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
	json.Unmarshal(mockBootstrap, &goldenBootstrap)

	type want struct {
		wantErr       bool
		containsErr   string
		wantBootstrap Bootstrap
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"returns rpc error",
			gtGoldenHTTPMock(bootstrapHandlerMock(mockRPCError, blankHandler)),
			want{
				true,
				"could not get bootstrap",
				Bootstrap{},
			},
		},
		{
			"fails to unmarshal",
			gtGoldenHTTPMock(bootstrapHandlerMock([]byte(`junk`), blankHandler)),
			want{
				true,
				"could not unmarshal bootstrap",
				Bootstrap{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(bootstrapHandlerMock(mockBootstrap, blankHandler)),
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
	json.Unmarshal(mockCommit, &goldenCommit)

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
			gtGoldenHTTPMock(commitHandlerMock(mockRPCError, blankHandler)),
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
			gtGoldenHTTPMock(commitHandlerMock(mockCommit, blankHandler)),
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
