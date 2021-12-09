package rpc_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/completium/go-tezos/v4/rpc"
	"github.com/stretchr/testify/assert"
)

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
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regBallotList, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get ballot list",
				rpc.BallotList{},
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regBallotList, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get ballot list: failed to parse json",
				rpc.BallotList{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regBallotList, readResponse(ballotList)}, blankHandler)),
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

			r, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, ballotList, err := r.BallotList(&rpc.BlockIDHead{})
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
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regBallots, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get ballots",
				rpc.Ballots{},
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regBallots, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get ballots: failed to parse json",
				rpc.Ballots{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regBallots, readResponse(ballots)}, blankHandler)),
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

			r, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, ballots, err := r.Ballots(&rpc.BlockIDHead{})
			checkErr(t, tt.wantErr, tt.containsErr, err)
			assert.Equal(t, tt.want.ballots, ballots)
		})
	}
}

func Test_CurrentPeriod(t *testing.T) {
	type want struct {
		wantErr       bool
		containsErr   string
		currentPeriod rpc.VotingPeriod
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"handles RPC error",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regCurrentPeriod, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get current period",
				rpc.VotingPeriod{},
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regCurrentPeriod, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get current period: failed to parse json",
				rpc.VotingPeriod{},
			},
		},
		// TODO: was unable to get real mock data
		// {
		// 	"is successful",
		// 	gtGoldenHTTPMock(mockHandler(&requestResultPair{regCurrentPeriod, []byte(`"promotion_vote"`)}, blankHandler)),
		// 	want{
		// 		false,
		// 		"",
		// 		"promotion_vote",
		// 	},
		// },
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			r, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, currentPeriodKind, err := r.CurrentPeriod(&rpc.BlockIDHead{})
			checkErr(t, tt.wantErr, tt.containsErr, err)
			assert.Equal(t, tt.want.currentPeriod, currentPeriodKind)
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
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regCurrentPeriodKind, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get current period kind",
				"",
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regCurrentPeriodKind, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get current period kind: failed to parse json",
				"",
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regCurrentPeriodKind, []byte(`"promotion_vote"`)}, blankHandler)),
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

			r, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, currentPeriodKind, err := r.CurrentPeriodKind(&rpc.BlockIDHead{})
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
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regCurrentProposal, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get current proposal",
				"",
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regCurrentProposal, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get current proposal: failed to parse json",
				"",
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regCurrentProposal, []byte(`"PtEdoTezd3RHSC31mpxxo1npxFjoWWcFgQtxapi51Z8TLu6v6Uq"`)}, blankHandler)),
			want{
				false,
				"",
				"PtEdoTezd3RHSC31mpxxo1npxFjoWWcFgQtxapi51Z8TLu6v6Uq",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			r, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, currentProposal, err := r.CurrentProposal(&rpc.BlockIDHead{})
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
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regCurrentQuorum, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get current quorum",
				0,
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regCurrentQuorum, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get current quorum: failed to parse json",
				0,
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regCurrentQuorum, []byte(`7470`)}, blankHandler)),
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

			r, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, currentQuorum, err := r.CurrentQuorum(&rpc.BlockIDHead{})
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
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regVoteListings, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get listings",
				rpc.Listings{},
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regVoteListings, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get listings: failed to parse json",
				rpc.Listings{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regVoteListings, readResponse(voteListings)}, blankHandler)),
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

			r, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, voteListings, err := r.Listings(&rpc.BlockIDHead{})
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
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regProposals, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get proposals",
				rpc.Proposals{},
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regProposals, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get proposals: failed to parse json",
				rpc.Proposals{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regProposals, readResponse(proposals)}, blankHandler)),
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

			r, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, proposals, err := r.Proposals(&rpc.BlockIDHead{})
			checkErr(t, tt.wantErr, tt.containsErr, err)
			assert.Equal(t, tt.want.proposals, proposals)
		})
	}
}

func Test_SuccessorPeriod(t *testing.T) {
	type want struct {
		wantErr         bool
		containsErr     string
		successorPeriod rpc.VotingPeriod
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"handles RPC error",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regSuccessorPeriod, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get successor period",
				rpc.VotingPeriod{},
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regSuccessorPeriod, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get successor period: failed to parse json",
				rpc.VotingPeriod{},
			},
		},
		// TODO: was unable to get real mock data
		// {
		// 	"is successful",
		// 	gtGoldenHTTPMock(mockHandler(&requestResultPair{regCurrentPeriod, []byte(`"promotion_vote"`)}, blankHandler)),
		// 	want{
		// 		false,
		// 		"",
		// 		"promotion_vote",
		// 	},
		// },
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			r, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, successorPeriod, err := r.SuccessorPeriod(&rpc.BlockIDHead{})
			checkErr(t, tt.wantErr, tt.containsErr, err)
			assert.Equal(t, tt.want.successorPeriod, successorPeriod)
		})
	}
}

func Test_TotalVotingPower(t *testing.T) {
	type want struct {
		wantErr     bool
		containsErr string
		votingPower int
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"handles RPC error",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regTotalVotingPower, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get total voting power",
				0,
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regTotalVotingPower, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get total voting power: failed to parse json",
				0,
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regTotalVotingPower, []byte(`7470`)}, blankHandler)),
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

			r, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, votingPower, err := r.TotalVotingPower(&rpc.BlockIDHead{})
			checkErr(t, tt.wantErr, tt.containsErr, err)
			assert.Equal(t, tt.want.votingPower, votingPower)
		})
	}
}
