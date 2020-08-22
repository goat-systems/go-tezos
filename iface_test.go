package gotezos

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_iface(t *testing.T) {
	server := httptest.NewServer(gtGoldenHTTPMock(blankHandler))
	defer server.Close()

	var gt IFace
	var err error
	gt, err = New(server.URL)
	assert.Nil(t, err)
	assert.NotNil(t, gt)
}
