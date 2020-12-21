package rpc

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_BakingRights(t *testing.T) {

	var goldenBakingRights BakingRights
	json.Unmarshal(readResponse(bakingrights), &goldenBakingRights)

	type want struct {
		wantErr          bool
		containsErr      string
		wantBakingRights *BakingRights
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
				&BakingRights{},
			},
		},
		{
			"fails to unmarshal",
			gtGoldenHTTPMock(bakingRightsHandlerMock([]byte(`junk`), blankHandler)),
			want{
				true,
				"could not unmarshal baking rights",
				&BakingRights{},
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

			rpc, err := New(server.URL)
			assert.Nil(t, err)

			bakingRights, err := rpc.BakingRights(BakingRightsInput{
				BlockHash: mockBlockHash,
			})
			checkErr(t, tt.wantErr, tt.containsErr, err)

			assert.Equal(t, tt.want.wantBakingRights, bakingRights)
		})
	}
}

func Test_EndorsingRights(t *testing.T) {
	goldenEndorsingRights := getResponse(endorsingrights).(*EndorsingRights)

	type want struct {
		wantErr             bool
		containsErr         string
		wantEndorsingRights *EndorsingRights
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
				&EndorsingRights{},
			},
		},
		{
			"fails to unmarshal",
			gtGoldenHTTPMock(endorsingRightsHandlerMock([]byte(`junk`), blankHandler)),
			want{
				true,
				"could not unmarshal endorsing rights",
				&EndorsingRights{},
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

			rpc, err := New(server.URL)
			assert.Nil(t, err)

			endorsingRights, err := rpc.EndorsingRights(EndorsingRightsInput{
				BlockHash: mockBlockHash,
			})
			checkErr(t, tt.wantErr, tt.containsErr, err)

			assert.Equal(t, tt.want.wantEndorsingRights, endorsingRights)
		})
	}
}
