package rpc_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goat-systems/go-tezos/v4/rpc"
	"github.com/stretchr/testify/assert"
)

func Test_BakingRights(t *testing.T) {
	goldenBakingRights := getResponse(bakingrights).(*rpc.BakingRights)

	type want struct {
		wantErr          bool
		containsErr      string
		wantBakingRights *rpc.BakingRights
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regBakingRights, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get baking rights",
				&rpc.BakingRights{},
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regBakingRights, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get baking rights: failed to parse json",
				&rpc.BakingRights{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regBakingRights, readResponse(bakingrights)}, blankHandler)),
			want{
				false,
				"",
				goldenBakingRights,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			r, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, bakingRights, err := r.BakingRights(rpc.BakingRightsInput{
				BlockID: &rpc.BlockIDHead{},
			})
			checkErr(t, tt.wantErr, tt.containsErr, err)
			assert.Equal(t, tt.want.wantBakingRights, bakingRights)
		})
	}
}

func Test_CompletePrefix(t *testing.T) {
	goldenCurrentLevel := getResponse(currentLevel).(rpc.CurrentLevel)

	type want struct {
		wantErr          bool
		containsErr      string
		wantCurrentLevel rpc.CurrentLevel
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regCurrentLevel, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get current level",
				rpc.CurrentLevel{},
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regCurrentLevel, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get current level: failed to parse json",
				rpc.CurrentLevel{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regCurrentLevel, readResponse(currentLevel)}, blankHandler)),
			want{
				false,
				"",
				goldenCurrentLevel,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			r, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, currentLevel, err := r.CurrentLevel(rpc.CurrentLevelInput{
				BlockID: &rpc.BlockIDHead{},
			})
			checkErr(t, tt.wantErr, tt.containsErr, err)
			assert.Equal(t, tt.want.wantCurrentLevel, currentLevel)
		})
	}
}

func Test_EndorsingRights(t *testing.T) {
	goldenEndorsingRights := getResponse(endorsingrights).(*rpc.EndorsingRights)

	type want struct {
		wantErr             bool
		containsErr         string
		wantEndorsingRights *rpc.EndorsingRights
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regEndorsingRights, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get endorsing rights",
				&rpc.EndorsingRights{},
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regEndorsingRights, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get endorsing rights: failed to parse json",
				&rpc.EndorsingRights{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regEndorsingRights, readResponse(endorsingrights)}, blankHandler)),
			want{
				false,
				"",
				goldenEndorsingRights,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			r, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, endorsingRights, err := r.EndorsingRights(rpc.EndorsingRightsInput{
				BlockID: &rpc.BlockIDHead{},
			})
			checkErr(t, tt.wantErr, tt.containsErr, err)
			assert.Equal(t, tt.want.wantEndorsingRights, endorsingRights)
		})
	}
}

func Test_ForgeOperationWithRPC(t *testing.T) {
	type want struct {
		err         bool
		errContains string
		operation   string
	}

	cases := []struct {
		name  string
		input http.Handler
		want  want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(forgeOperationWithRPCMock(readResponse(rpcerrors), blankHandler)),
			want{
				true,
				"failed to forge operation: rpc error (somekind)",
				"",
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(forgeOperationWithRPCMock([]byte(`junk`), blankHandler)),
			want{
				true,
				"failed to forge operation: invalid character",
				"",
			},
		},
		{
			"handles failure to strip operation branch",
			gtGoldenHTTPMock(forgeOperationWithRPCMock([]byte(`"some_junk_op_string"`), unforgeOperationWithRPCMock(readResponse(rpcerrors), blankHandler))),
			want{
				true,
				"failed to forge operation: unable to verify rpc returned a valid contents",
				"some_junk_op_string",
			},
		},
		{
			"handles failure to parse forged operation",
			gtGoldenHTTPMock(forgeOperationWithRPCMock([]byte(`"some_operation_string"`), unforgeOperationWithRPCMock(readResponse(rpcerrors), blankHandler))),
			want{
				true,
				"failed to forge operation: unable to verify rpc returned a valid contents",
				"some_operation_string",
			},
		},
		{
			"handles failure to match forge with expected contents",
			gtGoldenHTTPMock(forgeOperationWithRPCMock([]byte(`"a79ec80dba1f8ddb2cde90b8f12f7c62fdc36556030281ff8904a3d0df82cddc08000008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e00b960000008ba0cb2fad622697145cf1665124096d25bc31e00"`), unforgeOperationWithRPCMock(readResponse(parseOperations), blankHandler))),
			want{
				true,
				"failed to forge operation: alert rpc returned invalid contents",
				"a79ec80dba1f8ddb2cde90b8f12f7c62fdc36556030281ff8904a3d0df82cddc08000008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e00b960000008ba0cb2fad622697145cf1665124096d25bc31e00",
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(forgeOperationWithRPCMock([]byte(`"a79ec80dba1f8ddb2cde90b8f12f7c62fdc36556030281ff8904a3d0df82cddc08000008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e00b960000008ba0cb2fad622697145cf1665124096d25bc31e00"`), unforgeOperationWithRPCMock(readResponse(parseOperations), blankHandler))),
			want{
				false,
				"",
				"a79ec80dba1f8ddb2cde90b8f12f7c62fdc36556030281ff8904a3d0df82cddc08000008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e00b960000008ba0cb2fad622697145cf1665124096d25bc31e00",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input.inputHandler)
			defer server.Close()

			rpc, err := New(server.URL)
			assert.Nil(t, err)

			op, err := rpc.ForgeOperationWithRPC(tt.input.forgeOperationWithRPCInput)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.operation, op)
		})
	}
}
