package rpc_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goat-systems/go-tezos/v4/rpc"
	"github.com/stretchr/testify/assert"
)

func Test_Head(t *testing.T) {
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
				"could not get head block: invalid character",
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

			rpc, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, block, err := rpc.Head()
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
				"could not get block '50': invalid character",
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

			rpc, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, block, err := rpc.Block(50)
			checkErr(t, tt.wantErr, tt.containsErr, err)
			assert.Equal(t, tt.want.wantBlock, block)
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
