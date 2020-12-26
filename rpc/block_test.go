package rpc_test

import (
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

func Test_OperationHashes(t *testing.T) {
	goldenOperationHashses := getResponse(operationhashes).([][]string)

	type want struct {
		wantErr             bool
		containsErr         string
		wantOperationHashes [][]string
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"failed to unmarshal",
			gtGoldenHTTPMock(operationHashesHandlerMock([]byte(`junk`), blankHandler)),
			want{
				true,
				"could not unmarshal operation hashes",
				[][]string{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(operationHashesHandlerMock(readResponse(operationhashes), blankHandler)),
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

			rpc, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, operationHashes, err := rpc.OperationHashes("BLzGD63HA4RP8Fh5xEtvdQSMKa2WzJMZjQPNVUc4Rqy8Lh5BEY1")
			checkErr(t, tt.wantErr, tt.containsErr, err)
			assert.Equal(t, tt.want.wantOperationHashes, operationHashes)
		})
	}
}

func Test_BallotList(t *testing.T) {
	goldenBallotList := getResponse(ballotList).(*rpc.BallotList)

	type want struct {
		wantErr     bool
		containsErr string
		ballotList  rpc.BallotList
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"handles RPC error",
			gtGoldenHTTPMock(ballotListHandlerMock(readResponse(rpcerrors), blankHandler)),
			want{
				true,
				"failed to get ballot list",
				rpc.BallotList{},
			},
		},
		{
			"failed to unmarshal",
			gtGoldenHTTPMock(ballotListHandlerMock([]byte(`junk`), blankHandler)),
			want{
				true,
				"failed to unmarshal ballot list",
				rpc.BallotList{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(ballotListHandlerMock(readResponse(ballotList), blankHandler)),
			want{
				false,
				"",
				*goldenBallotList,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			rpc, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, ballotList, err := rpc.BallotList("BLzGD63HA4RP8Fh5xEtvdQSMKa2WzJMZjQPNVUc4Rqy8Lh5BEY1")
			checkErr(t, tt.wantErr, tt.containsErr, err)
			assert.Equal(t, tt.want.ballotList, ballotList)
		})
	}
}

func Test_Ballots(t *testing.T) {
	goldenBallots := getResponse(ballots).(*rpc.Ballots)

	type want struct {
		wantErr     bool
		containsErr string
		ballots     rpc.Ballots
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"handles RPC error",
			gtGoldenHTTPMock(ballotsHandlerMock(readResponse(rpcerrors), blankHandler)),
			want{
				true,
				"failed to get ballots",
				rpc.Ballots{},
			},
		},
		{
			"failed to unmarshal",
			gtGoldenHTTPMock(ballotsHandlerMock([]byte(`junk`), blankHandler)),
			want{
				true,
				"failed to unmarshal ballots",
				rpc.Ballots{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(ballotsHandlerMock(readResponse(ballots), blankHandler)),
			want{
				false,
				"",
				*goldenBallots,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			rpc, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, ballots, err := rpc.Ballots("BLzGD63HA4RP8Fh5xEtvdQSMKa2WzJMZjQPNVUc4Rqy8Lh5BEY1")
			checkErr(t, tt.wantErr, tt.containsErr, err)
			assert.Equal(t, tt.want.ballots, ballots)
		})
	}
}

func Test_CurrentPeriodKind(t *testing.T) {
	type want struct {
		wantErr           bool
		containsErr       string
		currentPeriodKind string
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"handles RPC error",
			gtGoldenHTTPMock(currentPeriodKindHandlerMock(readResponse(rpcerrors), blankHandler)),
			want{
				true,
				"failed to get current period kind",
				"",
			},
		},
		{
			"failed to unmarshal",
			gtGoldenHTTPMock(currentPeriodKindHandlerMock([]byte(`junk`), blankHandler)),
			want{
				true,
				"failed to unmarshal current period kind",
				"",
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(currentPeriodKindHandlerMock([]byte(`"promotion_vote"`), blankHandler)),
			want{
				false,
				"",
				"promotion_vote",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			rpc, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, currentPeriodKind, err := rpc.CurrentPeriodKind("BLzGD63HA4RP8Fh5xEtvdQSMKa2WzJMZjQPNVUc4Rqy8Lh5BEY1")
			checkErr(t, tt.wantErr, tt.containsErr, err)
			assert.Equal(t, tt.want.currentPeriodKind, currentPeriodKind)
		})
	}
}

func Test_CurrentProposal(t *testing.T) {
	type want struct {
		wantErr         bool
		containsErr     string
		currentProposal string
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"handles RPC error",
			gtGoldenHTTPMock(currentProposalHandlerMock(readResponse(rpcerrors), blankHandler)),
			want{
				true,
				"failed to get current proposal",
				"",
			},
		},
		{
			"failed to unmarshal",
			gtGoldenHTTPMock(currentProposalHandlerMock([]byte(`junk`), blankHandler)),
			want{
				true,
				"failed to unmarshal current proposal",
				"",
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(currentProposalHandlerMock([]byte(`"promotion_vote"`), blankHandler)),
			want{
				false,
				"",
				"promotion_vote",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			rpc, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, currentProposal, err := rpc.CurrentProposal("BLzGD63HA4RP8Fh5xEtvdQSMKa2WzJMZjQPNVUc4Rqy8Lh5BEY1")
			checkErr(t, tt.wantErr, tt.containsErr, err)
			assert.Equal(t, tt.want.currentProposal, currentProposal)
		})
	}
}

func Test_CurrentQuorum(t *testing.T) {
	type want struct {
		wantErr       bool
		containsErr   string
		currentQuorum int
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"handles RPC error",
			gtGoldenHTTPMock(currentQuorumHandlerMock(readResponse(rpcerrors), blankHandler)),
			want{
				true,
				"failed to get current quorum",
				0,
			},
		},
		{
			"failed to unmarshal",
			gtGoldenHTTPMock(currentQuorumHandlerMock([]byte(`junk`), blankHandler)),
			want{
				true,
				"failed to unmarshal current quorum",
				0,
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(currentQuorumHandlerMock([]byte(`7470`), blankHandler)),
			want{
				false,
				"",
				7470,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			rpc, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, currentQuorum, err := rpc.CurrentQuorum("BLzGD63HA4RP8Fh5xEtvdQSMKa2WzJMZjQPNVUc4Rqy8Lh5BEY1")
			checkErr(t, tt.wantErr, tt.containsErr, err)
			assert.Equal(t, tt.want.currentQuorum, currentQuorum)
		})
	}
}

func Test_VoteListings(t *testing.T) {
	goldenVoteListings := getResponse(voteListings).(rpc.Listings)

	type want struct {
		wantErr      bool
		containsErr  string
		voteListings rpc.Listings
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"handles RPC error",
			gtGoldenHTTPMock(voteListingsHandlerMock(readResponse(rpcerrors), blankHandler)),
			want{
				true,
				"failed to get listings",
				rpc.Listings{},
			},
		},
		{
			"failed to unmarshal",
			gtGoldenHTTPMock(voteListingsHandlerMock([]byte(`junk`), blankHandler)),
			want{
				true,
				"failed to unmarshal listings",
				rpc.Listings{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(voteListingsHandlerMock(readResponse(voteListings), blankHandler)),
			want{
				false,
				"",
				goldenVoteListings,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			rpc, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, voteListings, err := rpc.VoteListings("BLzGD63HA4RP8Fh5xEtvdQSMKa2WzJMZjQPNVUc4Rqy8Lh5BEY1")
			checkErr(t, tt.wantErr, tt.containsErr, err)
			assert.Equal(t, tt.want.voteListings, voteListings)
		})
	}
}

func Test_Proposals(t *testing.T) {
	goldenProposals := getResponse(proposals).(rpc.Proposals)

	type want struct {
		wantErr     bool
		containsErr string
		proposals   rpc.Proposals
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"handles RPC error",
			gtGoldenHTTPMock(proposalsHandlerMock(readResponse(rpcerrors), blankHandler)),
			want{
				true,
				"failed to get proposals",
				rpc.Proposals{},
			},
		},
		{
			"failed to unmarshal",
			gtGoldenHTTPMock(proposalsHandlerMock([]byte(`junk`), blankHandler)),
			want{
				true,
				"failed to unmarshal proposals",
				rpc.Proposals{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(proposalsHandlerMock(readResponse(proposals), blankHandler)),
			want{
				false,
				"",
				goldenProposals,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			rpc, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, proposals, err := rpc.Proposals("BLzGD63HA4RP8Fh5xEtvdQSMKa2WzJMZjQPNVUc4Rqy8Lh5BEY1")
			checkErr(t, tt.wantErr, tt.containsErr, err)
			assert.Equal(t, tt.want.proposals, proposals)
		})
	}
}
