package rpc_test

import (
	"net/http/httptest"
	"testing"

	"github.com/goat-systems/go-tezos/v4/rpc"
	"github.com/stretchr/testify/assert"
)

func Test_iface(t *testing.T) {
	server := httptest.NewServer(gtGoldenHTTPMock(blankHandler))
	defer server.Close()

	var r rpc.IFace
	var err error
	r, err = rpc.New(server.URL)
	assert.Nil(t, err)
	assert.NotNil(t, r)
}
