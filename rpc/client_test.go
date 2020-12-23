package rpc_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/goat-systems/go-tezos/v4/rpc"
	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	expectedConstants := expectedConstants(t)

	cases := []struct {
		name          string
		inputHandler  http.Handler
		wantErr       bool
		wantConstants *rpc.Constants
	}{
		{
			"Successful",
			gtGoldenHTTPMock(blankHandler),
			false,
			expectedConstants,
		},
		{
			"fails to fetch head block to get constants",
			http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				if strings.Contains(req.URL.String(), "/chains/main/blocks/head") {
					rw.Write([]byte(`some_junk_data`))
				}
			}),
			true,
			nil,
		},
		{
			"fails to fetch constants",
			http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				if strings.Contains(req.URL.String(), "/chains/main/blocks/head") {
					rw.Write(readResponse(block))
				}
				if strings.Contains(req.URL.String(), "/context/constants") {
					rw.Write([]byte(`some_junk_data`))
				}
			}),
			true,
			nil,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHandler)
			defer server.Close()

			r, err := rpc.New(server.URL)
			checkErr(t, tt.wantErr, "", err)

			assert.Equal(t, tt.wantConstants, r.CurrentContstants())
		})
	}
}
func Test_SetConstants(t *testing.T) {
	r := rpc.Client{}

	var constants rpc.Constants
	r.SetConstants(constants)

	assert.Equal(t, constants, r.CurrentContstants())
}

func gtGoldenHTTPMock(next http.Handler) http.Handler {
	var blockMock blockHandlerMock
	var constantsMock constantsHandlerMock
	return constantsMock.handler(
		readResponse(constants),
		blockMock.handler(
			readResponse(block),
			next,
		),
	)
}

func expectedConstants(t *testing.T) *rpc.Constants {
	var expectedConstants rpc.Constants
	err := json.Unmarshal(readResponse(constants), &expectedConstants)
	assert.Nilf(t, err, "could no unmarhsal mock constants")

	return &expectedConstants
}
