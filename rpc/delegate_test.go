package rpc_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goat-systems/go-tezos/v4/rpc"
	"github.com/stretchr/testify/assert"
)

func Test_BakingRights(t *testing.T) {

	var goldenBakingRights rpc.BakingRights
	json.Unmarshal(readResponse(bakingrights), &goldenBakingRights)

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
			"returns rpc error",
			gtGoldenHTTPMock(bakingRightsHandlerMock(readResponse(rpcerrors), blankHandler)),
			want{
				true,
				"could not get baking rights",
				&rpc.BakingRights{},
			},
		},
		{
			"fails to unmarshal",
			gtGoldenHTTPMock(bakingRightsHandlerMock([]byte(`junk`), blankHandler)),
			want{
				true,
				"could not unmarshal baking rights",
				&rpc.BakingRights{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(bakingRightsHandlerMock(readResponse(bakingrights), blankHandler)),
			want{
				false,
				"",
				&goldenBakingRights,
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
				BlockHash: mockBlockHash,
			})
			checkErr(t, tt.wantErr, tt.containsErr, err)

			assert.Equal(t, tt.want.wantBakingRights, bakingRights)
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
			"returns rpc error",
			gtGoldenHTTPMock(endorsingRightsHandlerMock(readResponse(rpcerrors), blankHandler)),
			want{
				true,
				"could not get endorsing rights",
				&rpc.EndorsingRights{},
			},
		},
		{
			"fails to unmarshal",
			gtGoldenHTTPMock(endorsingRightsHandlerMock([]byte(`junk`), blankHandler)),
			want{
				true,
				"could not unmarshal endorsing rights",
				&rpc.EndorsingRights{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(endorsingRightsHandlerMock(readResponse(endorsingrights), blankHandler)),
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
				BlockHash: mockBlockHash,
			})
			checkErr(t, tt.wantErr, tt.containsErr, err)

			assert.Equal(t, tt.want.wantEndorsingRights, endorsingRights)
		})
	}
}
