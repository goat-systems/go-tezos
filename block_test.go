package gotezos

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Head(t *testing.T) {
	var goldenBlock Block
	json.Unmarshal(mockBlockRandom, &goldenBlock)

	type want struct {
		wantErr     bool
		containsErr string
		wantBlock   Block
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"failed to unmarshal",
			gtGoldenHTTPMock(newBlockMock().handler([]byte(`not_block_data`), blankHandler)),
			want{
				true,
				"could not get head block: invalid character",
				Block{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(newBlockMock().handler(mockBlockRandom, blankHandler)),
			want{
				false,
				"",
				goldenBlock,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			gt, err := New(server.URL)
			assert.Nil(t, err)

			block, err := gt.Head()
			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tt.want.containsErr)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, tt.want.wantBlock, block)
		})
	}
}

func Test_Block(t *testing.T) {
	var goldenBlock Block
	json.Unmarshal(mockBlockRandom, &goldenBlock)

	type want struct {
		wantErr     bool
		containsErr string
		wantBlock   Block
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"failed to unmarshal",
			gtGoldenHTTPMock(newBlockMock().handler([]byte(`not_block_data`), blankHandler)),
			want{
				true,
				"could not get block '50': invalid character",
				Block{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(newBlockMock().handler(mockBlockRandom, blankHandler)),
			want{
				false,
				"",
				goldenBlock,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			gt, err := New(server.URL)
			assert.Nil(t, err)

			block, err := gt.Block(50)
			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tt.want.containsErr)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, tt.want.wantBlock, block)
		})
	}
}

func Test_OperationHashes(t *testing.T) {
	var goldenOperationHashses []string
	json.Unmarshal(mockOpHashes, &goldenOperationHashses)

	type want struct {
		wantErr             bool
		containsErr         string
		wantOperationHashes []string
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"failed to unmarshal",
			gtGoldenHTTPMock(opHashesHandlerMock([]byte(`junk`), blankHandler)),
			want{
				true,
				"could not unmarshal operation hashes",
				[]string{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(opHashesHandlerMock(mockOpHashes, blankHandler)),
			want{
				false,
				"",
				goldenOperationHashses,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			gt, err := New(server.URL)
			assert.Nil(t, err)

			operationHashes, err := gt.OperationHashes("BLzGD63HA4RP8Fh5xEtvdQSMKa2WzJMZjQPNVUc4Rqy8Lh5BEY1")
			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tt.want.containsErr)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, tt.want.wantOperationHashes, operationHashes)
		})
	}
}

func Test_idToString(t *testing.T) {
	cases := []struct {
		name    string
		input   interface{}
		wantErr bool
		wantID  string
	}{
		{
			"uses integer id",
			50,
			false,
			"50",
		},
		{
			"uses string id",
			"BLzGD63HA4RP8Fh5xEtvdQSMKa2WzJMZjQPNVUc4Rqy8Lh5BEY1",
			false,
			"BLzGD63HA4RP8Fh5xEtvdQSMKa2WzJMZjQPNVUc4Rqy8Lh5BEY1",
		},
		{
			"uses bad id type",
			45.433,
			true,
			"",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			id, err := idToString(tt.input)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)

			}
			assert.Equal(t, tt.wantID, id)
		})
	}
}
