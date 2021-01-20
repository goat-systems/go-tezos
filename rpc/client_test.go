package rpc_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goat-systems/go-tezos/v4/rpc"
	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	expectedConstants := getResponse(constants).(rpc.Constants)

	cases := []struct {
		name          string
		inputHandler  http.Handler
		wantErr       bool
		wantConstants *rpc.Constants
	}{
		{
			"Successful",
			newConstantsMock().handler(readResponse(constants), blankHandler),
			false,
			&expectedConstants,
		},
		{
			"fails to fetch constants",
			newConstantsMock().handler([]byte(`junk`), blankHandler),
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
			if err == nil {
				assert.Equal(t, *tt.wantConstants, r.CurrentContstants())
			}
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
	var constantsMock constantsHandlerMock
	return constantsMock.handler(
		readResponse(constants),
		next,
	)
}
