package rpc_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goat-systems/go-tezos/v4/rpc"
	"github.com/stretchr/testify/assert"
)

func Test_Block(t *testing.T) {
	goldenBlock := getResponse(block).(*rpc.Block)
	type want struct {
		wantErr     bool
		containsErr string
		wantBlock   *rpc.Block
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
				"failed to get block '50': failed to parse json",
				&rpc.Block{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(newBlockMock().handler(readResponse(block), blankHandler)),
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

			r, err := rpc.New(server.URL)
			assert.Nil(t, err)

			id := rpc.BlockIDLevel(50)
			_, block, err := r.Block(&id)
			checkErr(t, tt.wantErr, tt.containsErr, err)
			assert.Equal(t, tt.want.wantBlock, block)
		})
	}
}

func Test_EndorsingPower(t *testing.T) {
	type want struct {
		err         bool
		containsErr string
		result      int
	}

	cases := []struct {
		name  string
		input http.Handler
		want  want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regEndorsingPower, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get endorsing power",
				0,
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(newConstantsMock().handler([]byte(`junk`), blankHandler)),
			want{
				true,
				"failed to get endorsing power: failed to parse json",
				0,
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regEndorsingPower, []byte(`10`)}, blankHandler)),
			want{
				false,
				"",
				10,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input)
			defer server.Close()

			r, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, endorsingPower, err := r.EndorsingPower(rpc.EndorsingPowerInput{
				BlockID:        &rpc.BlockIDHead{},
				EndorsingPower: rpc.EndorsingPower{},
			})
			checkErr(t, tt.want.err, tt.want.containsErr, err)
			assert.Equal(t, tt.want.result, endorsingPower)
		})
	}
}

func Test_Hash(t *testing.T) {
	type want struct {
		err         bool
		containsErr string
		result      string
	}

	cases := []struct {
		name  string
		input http.Handler
		want  want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regHash, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get block 'head' hash",
				"",
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regHash, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get block 'head' hash: failed to parse json",
				"",
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regHash, []byte(`"BLzGD63HA4RP8Fh5xEtvdQSMKa2WzJMZjQPNVUc4Rqy8Lh5BEY1"`)}, blankHandler)),
			want{
				false,
				"",
				mockBlockHash,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input)
			defer server.Close()

			r, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, hash, err := r.Hash(&rpc.BlockIDHead{})
			checkErr(t, tt.want.err, tt.want.containsErr, err)
			assert.Equal(t, tt.want.result, hash)
		})
	}
}

func Test_Header(t *testing.T) {
	goldenHeader := getResponse(header).(rpc.Header)

	type want struct {
		err         bool
		containsErr string
		result      rpc.Header
	}

	cases := []struct {
		name  string
		input http.Handler
		want  want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regHeader, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get block 'head' header",
				rpc.Header{},
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regHeader, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get block 'head' header: failed to parse json",
				rpc.Header{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regHeader, readResponse(header)}, blankHandler)),
			want{
				false,
				"",
				goldenHeader,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input)
			defer server.Close()

			r, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, header, err := r.Header(&rpc.BlockIDHead{})
			checkErr(t, tt.want.err, tt.want.containsErr, err)
			assert.Equal(t, tt.want.result, header)
		})
	}
}

func Test_HeaderRaw(t *testing.T) {
	type want struct {
		err         bool
		containsErr string
		result      string
	}

	cases := []struct {
		name  string
		input http.Handler
		want  want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regHeaderRaw, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get block 'head' raw header",
				"",
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regHeaderRaw, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get block 'head' raw header: failed to parse json",
				"",
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regHeaderRaw, []byte(`"0000e69b63f1689f0200000b8ae69bfde78f9502f696a78bfe921d774acf67f0967275d391c9ace4a866a5356ef23ef826d3df7281ba4e5b57fad3a59d322554f18377d856dbc71d42c175"`)}, blankHandler)),
			want{
				false,
				"",
				"0000e69b63f1689f0200000b8ae69bfde78f9502f696a78bfe921d774acf67f0967275d391c9ace4a866a5356ef23ef826d3df7281ba4e5b57fad3a59d322554f18377d856dbc71d42c175",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input)
			defer server.Close()

			r, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, header, err := r.HeaderRaw(&rpc.BlockIDHead{})
			checkErr(t, tt.want.err, tt.want.containsErr, err)
			assert.Equal(t, tt.want.result, header)
		})
	}
}

func Test_HeaderShell(t *testing.T) {
	goldenHeaderShell := getResponse(headerShell).(rpc.HeaderShell)

	type want struct {
		err         bool
		containsErr string
		result      rpc.HeaderShell
	}

	cases := []struct {
		name  string
		input http.Handler
		want  want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regHeaderShell, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get block 'head' header shell",
				rpc.HeaderShell{},
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regHeaderShell, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get block 'head' header shell: failed to parse json",
				rpc.HeaderShell{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regHeaderShell, readResponse(headerShell)}, blankHandler)),
			want{
				false,
				"",
				goldenHeaderShell,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input)
			defer server.Close()

			r, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, headerShell, err := r.HeaderShell(&rpc.BlockIDHead{})
			checkErr(t, tt.want.err, tt.want.containsErr, err)
			assert.Equal(t, tt.want.result, headerShell)
		})
	}
}

func Test_HeaderProtocolData(t *testing.T) {
	goldenHeaderProtocolData := getResponse(protocolData).(rpc.ProtocolData)

	type want struct {
		err         bool
		containsErr string
		result      rpc.ProtocolData
	}

	cases := []struct {
		name  string
		input http.Handler
		want  want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regHeaderProtocolData, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get block 'head' protocol data",
				rpc.ProtocolData{},
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regHeaderProtocolData, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get block 'head' protocol data: failed to parse json",
				rpc.ProtocolData{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regHeaderProtocolData, readResponse(protocolData)}, blankHandler)),
			want{
				false,
				"",
				goldenHeaderProtocolData,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input)
			defer server.Close()

			r, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, protocolData, err := r.HeaderProtocolData(&rpc.BlockIDHead{})
			checkErr(t, tt.want.err, tt.want.containsErr, err)
			assert.Equal(t, tt.want.result, protocolData)
		})
	}
}

func Test_HeaderProtocolDataRaw(t *testing.T) {
	type want struct {
		err         bool
		containsErr string
		result      string
	}

	cases := []struct {
		name  string
		input http.Handler
		want  want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regHeaderProtocolDataRaw, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get block 'head' raw protocol data",
				"",
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regHeaderProtocolData, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get block 'head' raw protocol data: failed to parse json",
				"",
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regHeaderProtocolData, []byte(`"0000e69b63f1689f0200000b8ae69bfde78f9502f696a78bfe921d774acf67f0967275d391c9ace4a866a5356ef23ef826d3df7281ba4e5b57fad3a59d322554f18377d856dbc71d42c175"`)}, blankHandler)),
			want{
				false,
				"",
				"0000e69b63f1689f0200000b8ae69bfde78f9502f696a78bfe921d774acf67f0967275d391c9ace4a866a5356ef23ef826d3df7281ba4e5b57fad3a59d322554f18377d856dbc71d42c175",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input)
			defer server.Close()

			r, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, protocolData, err := r.HeaderProtocolDataRaw(&rpc.BlockIDHead{})
			checkErr(t, tt.want.err, tt.want.containsErr, err)
			assert.Equal(t, tt.want.result, protocolData)
		})
	}
}

func Test_LiveBlocks(t *testing.T) {
	goldenLiveBlocks := getResponse(liveBlocks).([]string)

	type want struct {
		err         bool
		containsErr string
		result      []string
	}

	cases := []struct {
		name  string
		input http.Handler
		want  want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regLiveBlocks, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get live blocks at 'head'",
				[]string{},
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regLiveBlocks, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get live blocks at 'head': failed to parse json",
				[]string{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regLiveBlocks, readResponse(liveBlocks)}, blankHandler)),
			want{
				false,
				"",
				goldenLiveBlocks,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input)
			defer server.Close()

			r, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, liveBlocks, err := r.LiveBlocks(&rpc.BlockIDHead{})
			checkErr(t, tt.want.err, tt.want.containsErr, err)
			assert.Equal(t, tt.want.result, liveBlocks)
		})
	}
}

func Test_Metadata(t *testing.T) {
	goldenMetadata := getResponse(metadata).(rpc.Metadata)

	type want struct {
		err         bool
		containsErr string
		result      rpc.Metadata
	}

	cases := []struct {
		name  string
		input http.Handler
		want  want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regMetadata, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get block 'head' metadata",
				rpc.Metadata{},
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regMetadata, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get block 'head' metadata: failed to parse json",
				rpc.Metadata{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regMetadata, readResponse(metadata)}, blankHandler)),
			want{
				false,
				"",
				goldenMetadata,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input)
			defer server.Close()

			r, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, metadata, err := r.Metadata(&rpc.BlockIDHead{})
			checkErr(t, tt.want.err, tt.want.containsErr, err)
			assert.Equal(t, tt.want.result, metadata)
		})
	}
}

func Test_MetadataHash(t *testing.T) {
	type want struct {
		err         bool
		containsErr string
		result      string
	}

	cases := []struct {
		name  string
		input http.Handler
		want  want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regMetadataHash, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get block 'head' metadata hash",
				"",
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regMetadataHash, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get block 'head' metadata hash: failed to parse json",
				"",
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regMetadataHash, []byte(`"some_hash"`)}, blankHandler)),
			want{
				false,
				"",
				"some_hash",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input)
			defer server.Close()

			r, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, metadataHash, err := r.MetadataHash(&rpc.BlockIDHead{})
			checkErr(t, tt.want.err, tt.want.containsErr, err)
			assert.Equal(t, tt.want.result, metadataHash)
		})
	}
}

func Test_MinimalValidTime(t *testing.T) {
	type want struct {
		err         bool
		containsErr string
		result      string
	}

	cases := []struct {
		name  string
		input http.Handler
		want  want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regMinimalValidTime, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get minimal valid time at 'head'",
				"0001-01-01 00:00:00 +0000 UTC",
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regMinimalValidTime, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get minimal valid time at 'head': failed to parse json",
				"0001-01-01 00:00:00 +0000 UTC",
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regMinimalValidTime, []byte(`"2020-12-29T00:07:52Z"`)}, blankHandler)),
			want{
				false,
				"",
				"2020-12-29 00:07:52 +0000 UTC",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input)
			defer server.Close()

			r, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, minimalValidTime, err := r.MinimalValidTime(rpc.MinimalValidTimeInput{
				BlockID: &rpc.BlockIDHead{},
			})
			checkErr(t, tt.want.err, tt.want.containsErr, err)
			assert.Equal(t, tt.want.result, minimalValidTime.String())
		})
	}
}
func Test_OperationHashes(t *testing.T) {
	goldenOperationHashses := getResponse(operationhashes).(rpc.OperationHashes)

	type want struct {
		wantErr             bool
		containsErr         string
		wantOperationHashes rpc.OperationHashes
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regOperationHashes, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get block 'head' operation hashes",
				rpc.OperationHashes{},
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regOperationHashes, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get block 'head' operation hashes",
				rpc.OperationHashes{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regOperationHashes, readResponse(operationhashes)}, blankHandler)),
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

			r, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, operationHashes, err := r.OperationHashes(rpc.OperationHashesInput{
				BlockID: &rpc.BlockIDHead{},
			})
			checkErr(t, tt.wantErr, tt.containsErr, err)
			assert.Equal(t, tt.want.wantOperationHashes, operationHashes)
		})
	}
}

func Test_OperationMetadataHashes(t *testing.T) {
	goldenOperationMetadataHashses := getResponse(operationMetaDataHashes).(rpc.OperationMetadataHashes)

	type want struct {
		wantErr                     bool
		containsErr                 string
		wantOperationMetadataHashes rpc.OperationMetadataHashes
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regOperationMetadataHashes, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get block 'head' operation metadata hashes",
				rpc.OperationMetadataHashes{},
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regOperationMetadataHashes, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get block 'head' operation metadata hashes",
				rpc.OperationMetadataHashes{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regOperationMetadataHashes, readResponse(operationMetaDataHashes)}, blankHandler)),
			want{
				false,
				"",
				goldenOperationMetadataHashses,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			r, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, operationMetadataHashes, err := r.OperationMetadataHashes(rpc.OperationMetadataHashesInput{
				BlockID: &rpc.BlockIDHead{},
			})
			checkErr(t, tt.wantErr, tt.containsErr, err)
			assert.Equal(t, tt.want.wantOperationMetadataHashes, operationMetadataHashes)
		})
	}
}

func Test_Operations(t *testing.T) {
	goldenOperations := getResponse(operations).(rpc.FlattenedOperations)

	type want struct {
		wantErr        bool
		containsErr    string
		wantOperations rpc.FlattenedOperations
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regOperations, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get block 'head' operations",
				rpc.FlattenedOperations{},
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regOperations, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get block 'head' operations",
				rpc.FlattenedOperations{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regOperations, readResponse(operations)}, blankHandler)),
			want{
				false,
				"",
				goldenOperations,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			r, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, operations, err := r.Operations(rpc.OperationsInput{
				BlockID: &rpc.BlockIDHead{},
			})
			checkErr(t, tt.wantErr, tt.containsErr, err)
			assert.Equal(t, tt.want.wantOperations, operations)
		})
	}
}

func Test_OperationsMetadataHash(t *testing.T) {
	type want struct {
		wantErr                    bool
		containsErr                string
		wantOperationsMetadataHash string
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regOperationsMetadataHash, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get block 'head' operations metadata hash",
				"",
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regOperationsMetadataHash, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get block 'head' operations metadata hash",
				"",
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regOperationsMetadataHash, []byte(`"some_hash"`)}, blankHandler)),
			want{
				false,
				"",
				"some_hash",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			r, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, operationsMetadataHash, err := r.OperationsMetadataHash(&rpc.BlockIDHead{})
			checkErr(t, tt.wantErr, tt.containsErr, err)
			assert.Equal(t, tt.want.wantOperationsMetadataHash, operationsMetadataHash)
		})
	}
}

func Test_Protocols(t *testing.T) {
	goldenProtocols := getResponse(protocols).(rpc.Protocols)

	type want struct {
		wantErr       bool
		containsErr   string
		wantProtocols rpc.Protocols
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regProtocols, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get block 'head' protocols",
				rpc.Protocols{},
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regProtocols, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get block 'head' protocols",
				rpc.Protocols{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regProtocols, readResponse(protocols)}, blankHandler)),
			want{
				false,
				"",
				goldenProtocols,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			r, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, protocols, err := r.Protocols(&rpc.BlockIDHead{})
			checkErr(t, tt.wantErr, tt.containsErr, err)
			assert.Equal(t, tt.want.wantProtocols, protocols)
		})
	}
}

func Test_RequiredEndorsements(t *testing.T) {
	type want struct {
		wantErr                  bool
		containsErr              string
		wantRequiredEndorsements int
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regRequiredEndorsements, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get block 'head' required endorsements",
				0,
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regRequiredEndorsements, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get block 'head' required endorsements",
				0,
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regRequiredEndorsements, []byte(`10`)}, blankHandler)),
			want{
				false,
				"",
				10,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			r, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, requiredEndorsements, err := r.RequiredEndorsements(rpc.RequiredEndorsementsInput{
				BlockID: &rpc.BlockIDHead{},
			})
			checkErr(t, tt.wantErr, tt.containsErr, err)
			assert.Equal(t, tt.want.wantRequiredEndorsements, requiredEndorsements)
		})
	}
}

func Test_TransactionEntrypoints(t *testing.T) {
	transactionJSON := []byte(`{
		"kind":"transaction",
		"source":"tz1Y7kGA8vLAwBkZfTgj4MJS4zgHXSmMn7tW",
		"fee":"20925",
		"counter":"6845930",
		"gas_limit":"205508",
		"storage_limit":"0",
		"amount":"0",
		"destination":"KT1PWx2mnDueood7fEmfbBDKx1D9BAnnXitn",
		"parameters":{
		   "entrypoint":"approve",
		   "value":{
			  "prim":"Pair",
			  "args":[
				 {
					"string":"KT1DrJV8vhkdLEj76h1H9Q4irZDqAkMPo1Qf"
				 },
				 {
					"int":"200000"
				 }
			  ]
		   }
		},
		"metadata":{
		   "balance_updates":[
			  {
				 "kind":"contract",
				 "contract":"tz1Y7kGA8vLAwBkZfTgj4MJS4zgHXSmMn7tW",
				 "change":"-20925"
			  },
			  {
				 "kind":"freezer",
				 "category":"fees",
				 "delegate":"tz1P7wwnURM4iccy9sSS6Bo2tay9JpuPHzf5",
				 "cycle":280,
				 "change":"20925"
			  }
		   ],
		   "operation_result":{
			  "status":"applied",
			  "storage":{
				 "prim":"Pair",
				 "args":[
					{
					   "int":"31"
					},
					{
					   "prim":"Pair",
					   "args":[
						  [
							 {
								"prim":"DUP"
							 },
							 {
								"prim":"CAR"
							 },
							 {
								"prim":"DIP",
								"args":[
								   [
									  {
										 "prim":"CDR"
									  }
								   ]
								]
							 },
							 {
								"prim":"DUP"
							 },
							 {
								"prim":"DUP"
							 },
							 {
								"prim":"CAR"
							 },
							 {
								"prim":"DIP",
								"args":[
								   [
									  {
										 "prim":"CDR"
									  }
								   ]
								]
							 },
							 {
								"prim":"DIP",
								"args":[
								   [
									  {
										 "prim":"DIP",
										 "args":[
											{
											   "int":"2"
											},
											[
											   {
												  "prim":"DUP"
											   }
											]
										 ]
									  },
									  {
										 "prim":"DIG",
										 "args":[
											{
											   "int":"2"
											}
										 ]
									  }
								   ]
								]
							 },
							 {
								"prim":"PUSH",
								"args":[
								   {
									  "prim":"string"
								   },
								   {
									  "string":"code"
								   }
								]
							 },
							 {
								"prim":"PAIR"
							 },
							 {
								"prim":"PACK"
							 },
							 {
								"prim":"GET"
							 },
							 {
								"prim":"IF_NONE",
								"args":[
								   [
									  {
										 "prim":"NONE",
										 "args":[
											{
											   "prim":"lambda",
											   "args":[
												  {
													 "prim":"pair",
													 "args":[
														{
														   "prim":"bytes"
														},
														{
														   "prim":"big_map",
														   "args":[
															  {
																 "prim":"bytes"
															  },
															  {
																 "prim":"bytes"
															  }
														   ]
														}
													 ]
												  },
												  {
													 "prim":"pair",
													 "args":[
														{
														   "prim":"list",
														   "args":[
															  {
																 "prim":"operation"
															  }
														   ]
														},
														{
														   "prim":"big_map",
														   "args":[
															  {
																 "prim":"bytes"
															  },
															  {
																 "prim":"bytes"
															  }
														   ]
														}
													 ]
												  }
											   ]
											}
										 ]
									  }
								   ],
								   [
									  {
										 "prim":"UNPACK",
										 "args":[
											{
											   "prim":"lambda",
											   "args":[
												  {
													 "prim":"pair",
													 "args":[
														{
														   "prim":"bytes"
														},
														{
														   "prim":"big_map",
														   "args":[
															  {
																 "prim":"bytes"
															  },
															  {
																 "prim":"bytes"
															  }
														   ]
														}
													 ]
												  },
												  {
													 "prim":"pair",
													 "args":[
														{
														   "prim":"list",
														   "args":[
															  {
																 "prim":"operation"
															  }
														   ]
														},
														{
														   "prim":"big_map",
														   "args":[
															  {
																 "prim":"bytes"
															  },
															  {
																 "prim":"bytes"
															  }
														   ]
														}
													 ]
												  }
											   ]
											}
										 ]
									  },
									  {
										 "prim":"IF_NONE",
										 "args":[
											[
											   {
												  "prim":"PUSH",
												  "args":[
													 {
														"prim":"string"
													 },
													 {
														"string":"UStore: failed to unpack code"
													 }
												  ]
											   },
											   {
												  "prim":"FAILWITH"
											   }
											],
											[
											   
											]
										 ]
									  },
									  {
										 "prim":"SOME"
									  }
								   ]
								]
							 },
							 {
								"prim":"IF_NONE",
								"args":[
								   [
									  {
										 "prim":"DROP"
									  },
									  {
										 "prim":"DIP",
										 "args":[
											[
											   {
												  "prim":"DUP"
											   },
											   {
												  "prim":"PUSH",
												  "args":[
													 {
														"prim":"bytes"
													 },
													 {
														"bytes":"05010000000866616c6c6261636b"
													 }
												  ]
											   },
											   {
												  "prim":"GET"
											   },
											   {
												  "prim":"IF_NONE",
												  "args":[
													 [
														{
														   "prim":"PUSH",
														   "args":[
															  {
																 "prim":"string"
															  },
															  {
																 "string":"UStore: no field fallback"
															  }
														   ]
														},
														{
														   "prim":"FAILWITH"
														}
													 ],
													 [
														
													 ]
												  ]
											   },
											   {
												  "prim":"UNPACK",
												  "args":[
													 {
														"prim":"lambda",
														"args":[
														   {
															  "prim":"pair",
															  "args":[
																 {
																	"prim":"pair",
																	"args":[
																	   {
																		  "prim":"string"
																	   },
																	   {
																		  "prim":"bytes"
																	   }
																	]
																 },
																 {
																	"prim":"big_map",
																	"args":[
																	   {
																		  "prim":"bytes"
																	   },
																	   {
																		  "prim":"bytes"
																	   }
																	]
																 }
															  ]
														   },
														   {
															  "prim":"pair",
															  "args":[
																 {
																	"prim":"list",
																	"args":[
																	   {
																		  "prim":"operation"
																	   }
																	]
																 },
																 {
																	"prim":"big_map",
																	"args":[
																	   {
																		  "prim":"bytes"
																	   },
																	   {
																		  "prim":"bytes"
																	   }
																	]
																 }
															  ]
														   }
														]
													 }
												  ]
											   },
											   {
												  "prim":"IF_NONE",
												  "args":[
													 [
														{
														   "prim":"PUSH",
														   "args":[
															  {
																 "prim":"string"
															  },
															  {
																 "string":"UStore: failed to unpack fallback"
															  }
														   ]
														},
														{
														   "prim":"FAILWITH"
														}
													 ],
													 [
														
													 ]
												  ]
											   },
											   {
												  "prim":"SWAP"
											   }
											]
										 ]
									  },
									  {
										 "prim":"PAIR"
									  },
									  {
										 "prim":"EXEC"
									  }
								   ],
								   [
									  {
										 "prim":"DIP",
										 "args":[
											[
											   {
												  "prim":"SWAP"
											   },
											   {
												  "prim":"DROP"
											   },
											   {
												  "prim":"PAIR"
											   }
											]
										 ]
									  },
									  {
										 "prim":"SWAP"
									  },
									  {
										 "prim":"EXEC"
									  }
								   ]
								]
							 }
						  ],
						  {
							 "prim":"Pair",
							 "args":[
								{
								   "int":"1"
								},
								{
								   "prim":"False"
								}
							 ]
						  }
					   ]
					}
				 ]
			  },
			  "big_map_diff":[
				 {
					"action":"update",
					"big_map":"31",
					"key_hash":"exprtu2J7H6H5ERrw6EwAccpqChEkUVBMYCLLEVJBtzpT7kVPh4wKy",
					"key":{
					   "bytes":"05070701000000066c65646765720a00000016000088df70e2e368821a08166a1ad762c9bc5bced878"
					},
					"value":{
					   "bytes":"05070700a0c6b60b020000002107040a000000160139c8ade2617663981fa2b87592c9ad92714d14c2000080b518"
					}
				 }
			  ],
			  "consumed_gas":"205408",
			  "storage_size":"24637"
		   }
		}
	 }`)

	var transaction rpc.Transaction
	err := json.Unmarshal(transactionJSON, &transaction)
	if assert.Nil(t, err) {
		assert.Equal(t, "approve", transaction.Parameters.Entrypoint)
	}
}
