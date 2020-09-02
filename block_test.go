package gotezos

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Head(t *testing.T) {
	goldenBlock := getResponse(block).(*Block)
	type want struct {
		wantErr     bool
		containsErr string
		wantBlock   *Block
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
				&Block{},
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
	goldenBlock := getResponse(block).(*Block)
	type want struct {
		wantErr     bool
		containsErr string
		wantBlock   *Block
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
				&Block{},
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

			gt, err := New(server.URL)
			assert.Nil(t, err)

			block, err := gt.Block(50)
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

			gt, err := New(server.URL)
			assert.Nil(t, err)

			operationHashes, err := gt.OperationHashes("BLzGD63HA4RP8Fh5xEtvdQSMKa2WzJMZjQPNVUc4Rqy8Lh5BEY1")
			checkErr(t, tt.wantErr, tt.containsErr, err)
			assert.Equal(t, tt.want.wantOperationHashes, operationHashes)
		})
	}
}

func Test_BallotList(t *testing.T) {
	goldenBallotList := getResponse(ballotList).(*BallotList)

	type want struct {
		wantErr     bool
		containsErr string
		ballotList  BallotList
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
				BallotList{},
			},
		},
		{
			"failed to unmarshal",
			gtGoldenHTTPMock(ballotListHandlerMock([]byte(`junk`), blankHandler)),
			want{
				true,
				"failed to unmarshal ballot list",
				BallotList{},
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

			gt, err := New(server.URL)
			assert.Nil(t, err)

			ballotList, err := gt.BallotList("BLzGD63HA4RP8Fh5xEtvdQSMKa2WzJMZjQPNVUc4Rqy8Lh5BEY1")
			checkErr(t, tt.wantErr, tt.containsErr, err)
			assert.Equal(t, tt.want.ballotList, ballotList)
		})
	}
}

func Test_Ballots(t *testing.T) {
	goldenBallots := getResponse(ballots).(*Ballots)

	type want struct {
		wantErr     bool
		containsErr string
		ballots     Ballots
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
				Ballots{},
			},
		},
		{
			"failed to unmarshal",
			gtGoldenHTTPMock(ballotsHandlerMock([]byte(`junk`), blankHandler)),
			want{
				true,
				"failed to unmarshal ballots",
				Ballots{},
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

			gt, err := New(server.URL)
			assert.Nil(t, err)

			ballots, err := gt.Ballots("BLzGD63HA4RP8Fh5xEtvdQSMKa2WzJMZjQPNVUc4Rqy8Lh5BEY1")
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

			gt, err := New(server.URL)
			assert.Nil(t, err)

			currentPeriodKind, err := gt.CurrentPeriodKind("BLzGD63HA4RP8Fh5xEtvdQSMKa2WzJMZjQPNVUc4Rqy8Lh5BEY1")
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

			gt, err := New(server.URL)
			assert.Nil(t, err)

			currentProposal, err := gt.CurrentProposal("BLzGD63HA4RP8Fh5xEtvdQSMKa2WzJMZjQPNVUc4Rqy8Lh5BEY1")
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

			gt, err := New(server.URL)
			assert.Nil(t, err)

			currentQuorum, err := gt.CurrentQuorum("BLzGD63HA4RP8Fh5xEtvdQSMKa2WzJMZjQPNVUc4Rqy8Lh5BEY1")
			checkErr(t, tt.wantErr, tt.containsErr, err)
			assert.Equal(t, tt.want.currentQuorum, currentQuorum)
		})
	}
}

func Test_VoteListings(t *testing.T) {
	goldenVoteListings := getResponse(voteListings).(Listings)

	type want struct {
		wantErr      bool
		containsErr  string
		voteListings Listings
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
				Listings{},
			},
		},
		{
			"failed to unmarshal",
			gtGoldenHTTPMock(voteListingsHandlerMock([]byte(`junk`), blankHandler)),
			want{
				true,
				"failed to unmarshal listings",
				Listings{},
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

			gt, err := New(server.URL)
			assert.Nil(t, err)

			voteListings, err := gt.VoteListings("BLzGD63HA4RP8Fh5xEtvdQSMKa2WzJMZjQPNVUc4Rqy8Lh5BEY1")
			checkErr(t, tt.wantErr, tt.containsErr, err)
			assert.Equal(t, tt.want.voteListings, voteListings)
		})
	}
}

func Test_Proposals(t *testing.T) {
	goldenProposals := getResponse(proposals).(Proposals)

	type want struct {
		wantErr     bool
		containsErr string
		proposals   Proposals
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
				Proposals{},
			},
		},
		{
			"failed to unmarshal",
			gtGoldenHTTPMock(proposalsHandlerMock([]byte(`junk`), blankHandler)),
			want{
				true,
				"failed to unmarshal proposals",
				Proposals{},
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

			gt, err := New(server.URL)
			assert.Nil(t, err)

			proposals, err := gt.Proposals("BLzGD63HA4RP8Fh5xEtvdQSMKa2WzJMZjQPNVUc4Rqy8Lh5BEY1")
			checkErr(t, tt.wantErr, tt.containsErr, err)
			assert.Equal(t, tt.want.proposals, proposals)
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
			checkErr(t, tt.wantErr, "", err)
			assert.Equal(t, tt.wantID, id)
		})
	}
}

func Test_BigMapDiff(t *testing.T) {
	cases := []struct {
		name       string
		bigMapDiff []byte
		wantErr    bool
		want       BigMapDiff
	}{
		{
			"successfully unmarshals BigMapDiff action update",
			[]byte(`[
						{
						"action":"update",
						"big_map":"52",
						"key_hash":"exprta5PGni3vkj7z6B5CHRELDe796kyPq7q9qAqzadnm3fr4AvNhJ",
						"key":{
							"string":"6238d74df3089fe8b263422eea4f35101aa2b8bb50687aa98bdb15e1111b909d"
						},
						"value":{
							"prim":"Pair",
							"args":[
								{
									"int":"1593806466"
								},
								{
									"bytes":"00004cc5b68779c9166b20f6aed04f6fb7b01929ab9a"
								}
							]
						}
						}
			 		]`),
			false,
			BigMapDiff{
				Updates: []BigMapDiffUpdate{
					{
						Action:  "update",
						BigMap:  52,
						KeyHash: "exprta5PGni3vkj7z6B5CHRELDe796kyPq7q9qAqzadnm3fr4AvNhJ",
						Key: &MichelineExpression{
							Object: &Micheline{
								Int:    "",
								String: "6238d74df3089fe8b263422eea4f35101aa2b8bb50687aa98bdb15e1111b909d",
								Bytes:  "",
							},
						},
						Value: &MichelineExpression{
							Object: &Micheline{
								Prim: "Pair",
								Args: &MichelineArgs{
									Array: []Micheline{
										{
											Int: "1593806466",
										},
										{
											Bytes: "00004cc5b68779c9166b20f6aed04f6fb7b01929ab9a",
										},
									},
								},
							},
						},
					},
				},
				Removals: nil,
				Copies:   nil,
				Alloc:    nil,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			var bigMapDiff BigMapDiff
			err := bigMapDiff.UnmarshalJSON(tt.bigMapDiff)
			checkErr(t, tt.wantErr, "", err)

			assert.Equal(t, tt.want, bigMapDiff)

			v, err := bigMapDiff.MarshalJSON()
			checkErr(t, tt.wantErr, "", err)

			b, err := tt.want.MarshalJSON()
			checkErr(t, tt.wantErr, "", err)

			assert.Equal(t, string(v), string(b))
		})
	}
}

func Test_Contents(t *testing.T) {
	cases := []struct {
		name     string
		contents []byte
		wantErr  bool
		want     Contents
	}{
		{
			"successfully unmarshals and marshals endorsement",
			[]byte(`[{"kind":"endorsement","level":839680,"metadata":{"balance_updates":[{"kind":"contract","contract":"tz1iZEKy4LaAjnTmn2RuGDf2iqdAQKnRi8kY","change":"-64000000"},{"kind":"freezer","category":"deposits","delegate":"tz1iZEKy4LaAjnTmn2RuGDf2iqdAQKnRi8kY","cycle":204,"change":"64000000"},{"kind":"freezer","category":"rewards","delegate":"tz1iZEKy4LaAjnTmn2RuGDf2iqdAQKnRi8kY","cycle":204,"change":"2000000"}],"delegate":"tz1iZEKy4LaAjnTmn2RuGDf2iqdAQKnRi8kY","slots":[1]}}]`),
			false,
			Contents{
				Endorsements: []Endorsement{
					{
						Kind:  "endorsement",
						Level: 839680,
						Metadata: &EndorsementMetadata{
							BalanceUpdates: []BalanceUpdates{
								{
									Kind:     "contract",
									Contract: "tz1iZEKy4LaAjnTmn2RuGDf2iqdAQKnRi8kY",
									Change:   -64000000,
								},
								{
									Kind:     "freezer",
									Category: "deposits",
									Delegate: "tz1iZEKy4LaAjnTmn2RuGDf2iqdAQKnRi8kY",
									Cycle:    204,
									Change:   64000000,
								},
								{
									Kind:     "freezer",
									Category: "rewards",
									Delegate: "tz1iZEKy4LaAjnTmn2RuGDf2iqdAQKnRi8kY",
									Cycle:    204,
									Change:   2000000,
								},
							},
							Delegate: "tz1iZEKy4LaAjnTmn2RuGDf2iqdAQKnRi8kY",
							Slots:    []int{1},
						},
					},
				},
			},
		},
		{
			"successfully unmarshals and marshals transaction",
			[]byte(`[{"kind":"transaction","source":"tz1Vyuu4EJ5Nym4JcrfRLnp3hpaq1DSEp1Ke","fee":"1792","counter":"929880","gas_limit":"15385","storage_limit":"0","amount":"2176730","destination":"KT1G7uE8NW2Jg7mA13gzFtLbX2zw3X7N1uYG","metadata":{"balance_updates":[{"kind":"contract","contract":"tz1Vyuu4EJ5Nym4JcrfRLnp3hpaq1DSEp1Ke","change":"-1792"},{"kind":"freezer","category":"fees","delegate":"tz1Vc9XAD7iphycJoRwE1Nxx5krB9C7XyBu5","cycle":260,"change":"1792"}],"operation_result":{"status":"applied","storage":{"bytes":"002e130ee23658766386fa47d81ca5f727129f2c72"},"balance_updates":[{"kind":"contract","contract":"tz1Vyuu4EJ5Nym4JcrfRLnp3hpaq1DSEp1Ke","change":"-2176730"},{"kind":"contract","contract":"KT1G7uE8NW2Jg7mA13gzFtLbX2zw3X7N1uYG","change":"2176730"}],"consumed_gas":"15285","storage_size":"232"}}}]`),
			false,
			Contents{
				Transactions: []Transaction{
					{
						Kind:         "transaction",
						Source:       "tz1Vyuu4EJ5Nym4JcrfRLnp3hpaq1DSEp1Ke",
						Fee:          1792,
						Counter:      929880,
						GasLimit:     15385,
						StorageLimit: 0,
						Amount:       2176730,
						Destination:  "KT1G7uE8NW2Jg7mA13gzFtLbX2zw3X7N1uYG",
						Metadata: &TransactionMetadata{
							BalanceUpdates: []BalanceUpdates{
								{
									Kind:     "contract",
									Contract: "tz1Vyuu4EJ5Nym4JcrfRLnp3hpaq1DSEp1Ke",
									Change:   -1792,
								},
								{
									Kind:     "freezer",
									Category: "fees",
									Delegate: "tz1Vc9XAD7iphycJoRwE1Nxx5krB9C7XyBu5",
									Cycle:    260,
									Change:   1792,
								},
							},
							OperationResult: OperationResultTransfer{
								Status: "applied",
								Storage: &MichelineExpression{
									Object: &Micheline{
										Bytes: "002e130ee23658766386fa47d81ca5f727129f2c72",
									},
								},
								BalanceUpdates: []BalanceUpdates{
									{
										Kind:     "contract",
										Contract: "tz1Vyuu4EJ5Nym4JcrfRLnp3hpaq1DSEp1Ke",
										Change:   -2176730,
									},
									{
										Kind:     "contract",
										Contract: "KT1G7uE8NW2Jg7mA13gzFtLbX2zw3X7N1uYG",
										Change:   2176730,
									},
								},
								ConsumedGas: 15285,
								StorageSize: 232,
							},
						},
					},
				},
			},
		},
		{
			"successfully unmarshals and marshals delegation",
			[]byte(`[{"kind":"delegation","source":"tz1Qf7Eyq2S74oRKSEnT3GqNdy7op2Jkc8Vz","fee":"30000","counter":"2612138","gas_limit":"18136","storage_limit":"257","metadata":{"balance_updates":[{"kind":"contract","contract":"tz1Qf7Eyq2S74oRKSEnT3GqNdy7op2Jkc8Vz","change":"-30000"},{"kind":"freezer","category":"fees","delegate":"tz1Vc9XAD7iphycJoRwE1Nxx5krB9C7XyBu5","cycle":260,"change":"30000"}],"operation_result":{"status":"applied","consumed_gas":"10000"}}}]`),
			false,
			Contents{
				Delegations: []Delegation{
					{
						Kind:         "delegation",
						Source:       "tz1Qf7Eyq2S74oRKSEnT3GqNdy7op2Jkc8Vz",
						Fee:          30000,
						Counter:      2612138,
						GasLimit:     18136,
						StorageLimit: 257,
						Metadata: &DelegationMetadata{
							BalanceUpdates: []BalanceUpdates{
								{
									Kind:     "contract",
									Contract: "tz1Qf7Eyq2S74oRKSEnT3GqNdy7op2Jkc8Vz",
									Change:   -30000,
								},
								{
									Kind:     "freezer",
									Category: "fees",
									Delegate: "tz1Vc9XAD7iphycJoRwE1Nxx5krB9C7XyBu5",
									Cycle:    260,
									Change:   30000,
								},
							},
							OperationResults: OperationResultDelegation{
								Status:      "applied",
								ConsumedGas: 10000,
							},
						},
					},
				},
			},
		},
		{
			"successfully unmarshals and marshals reveal",
			[]byte(`[{"kind":"reveal","source":"tz1SiokJ4WgfHBxRKwPtTHPqWiRqg9njweba","fee":"1300","counter":"5735942","gas_limit":"10000","storage_limit":"0","public_key":"edpktwt2W2eLfYppE6E7vdNdUmkpGX6BgpBFxvwzEMg87oTauJbUju","metadata":{"balance_updates":[{"kind":"contract","contract":"tz1SiokJ4WgfHBxRKwPtTHPqWiRqg9njweba","change":"-1300"},{"kind":"freezer","category":"fees","delegate":"tz1S8MNvuFEUsWgjHvi3AxibRBf388NhT1q2","cycle":260,"change":"1300"}],"operation_result":{"status":"applied","consumed_gas":"10000"}}}]`),
			false,
			Contents{
				Reveals: []Reveal{
					{
						Kind:         "reveal",
						Source:       "tz1SiokJ4WgfHBxRKwPtTHPqWiRqg9njweba",
						Fee:          1300,
						Counter:      5735942,
						GasLimit:     10000,
						StorageLimit: 0,
						PublicKey:    "edpktwt2W2eLfYppE6E7vdNdUmkpGX6BgpBFxvwzEMg87oTauJbUju",
						Metadata: &RevealMetadata{
							BalanceUpdates: []BalanceUpdates{
								{
									Kind:     "contract",
									Contract: "tz1SiokJ4WgfHBxRKwPtTHPqWiRqg9njweba",
									Change:   -1300,
								},
								{
									Kind:     "freezer",
									Category: "fees",
									Delegate: "tz1S8MNvuFEUsWgjHvi3AxibRBf388NhT1q2",
									Cycle:    260,
									Change:   1300,
								},
							},
							OperationResult: OperationResultReveal{
								Status:      "applied",
								ConsumedGas: 10000,
							},
						},
					},
				},
			},
		},
		{
			"successfully unmarshals and marshals account_activation",
			[]byte(`[{"kind":"activate_account","pkh":"tz1iTTEtNCQfm2hXqiuDoCQZPEUHF6J5bwDU","secret":"df6fc7828243ce0edebbad09357edee0471cc617","metadata":{"balance_updates":[{"kind":"contract","contract":"tz1iTTEtNCQfm2hXqiuDoCQZPEUHF6J5bwDU","change":"2427770280"}]}}]`),
			false,
			Contents{
				AccountActivations: []AccountActivation{
					{
						Kind:   "activate_account",
						Pkh:    "tz1iTTEtNCQfm2hXqiuDoCQZPEUHF6J5bwDU",
						Secret: "df6fc7828243ce0edebbad09357edee0471cc617",
						Metadata: &AccountActivationMetadata{
							BalanceUpdates: []BalanceUpdates{
								{
									Kind:     "contract",
									Contract: "tz1iTTEtNCQfm2hXqiuDoCQZPEUHF6J5bwDU",
									Change:   2427770280,
								},
							},
						},
					},
				},
			},
		},
		{
			"successfully unmarshals and marshals proposals",
			[]byte(`[{"kind":"proposals","source":"tz1fNdh4YftsUasbB1BWBpqDmr4sFZaPNZVL","period":10,"proposals":["Pt24m4xiPbLDhVgVfABUjirbmda3yohdN82Sp9FeuAXJ4eV9otd","Psd1ynUBhMZAeajwcZJAeq5NrxorM6UCU4GJqxZ7Bx2e9vUWB6z"],"metadata":{}}]`),
			false,
			Contents{
				Proposals: []Proposal{
					{
						Kind:   "proposals",
						Source: "tz1fNdh4YftsUasbB1BWBpqDmr4sFZaPNZVL",
						Period: 10,
						Proposals: []string{
							"Pt24m4xiPbLDhVgVfABUjirbmda3yohdN82Sp9FeuAXJ4eV9otd",
							"Psd1ynUBhMZAeajwcZJAeq5NrxorM6UCU4GJqxZ7Bx2e9vUWB6z",
						},
					},
				},
			},
		},
		{
			"successfully unmarshals and marshals ballots",
			[]byte(`[{"kind":"ballot","source":"tz1bf816tUrSLYsWkUrsDFH9kkbpg3oXjriR","period":11,"proposal":"Pt24m4xiPbLDhVgVfABUjirbmda3yohdN82Sp9FeuAXJ4eV9otd","ballot":"yay","metadata":{}}]`),
			false,
			Contents{
				Ballots: []Ballot{
					{
						Kind:     "ballot",
						Source:   "tz1bf816tUrSLYsWkUrsDFH9kkbpg3oXjriR",
						Period:   11,
						Proposal: "Pt24m4xiPbLDhVgVfABUjirbmda3yohdN82Sp9FeuAXJ4eV9otd",
						Ballot:   "yay",
					},
				},
			},
		},
		{
			"successfully unmarshals and marshals double baking evidence",
			[]byte(`[{"kind":"double_baking_evidence","bh1":{"level":32958,"proto":2,"predecessor":"BMVMTfywir5E2QucCxwv2DkM2CYtnnqesriXJtegEGqfLRt9Sxh","timestamp":"2018-07-25T01:36:57Z","validation_pass":4,"operations_hash":"LLoZfBngaArquerV55FPZsTMbBALJ5Pfrh9KAUzHKEJPYAutMDLcV","fitness":["00","000000000009e668"],"context":"CoVauxcNbkGuCVE9ajsbVhdDGg9zs7CwvLAjtJ6zpGLaMcWhtWRq","priority":0,"proof_of_work_nonce":"9f41eeca8bab4056","signature":"sigceQeBtuHbhcvpTM8QXAztU4KzeVKQCNapSA5mqCBc6tVC5FYQBhTAyr75cfFe9J3yX7YH1Yy769tHyHJLG4NCY4Azoujb"},"bh2":{"level":32958,"proto":2,"predecessor":"BMVMTfywir5E2QucCxwv2DkM2CYtnnqesriXJtegEGqfLRt9Sxh","timestamp":"2018-07-25T01:36:57Z","validation_pass":4,"operations_hash":"LLoZfBngaArquerV55FPZsTMbBALJ5Pfrh9KAUzHKEJPYAutMDLcV","fitness":["00","000000000009e668"],"context":"CoVauxcNbkGuCVE9ajsbVhdDGg9zs7CwvLAjtJ6zpGLaMcWhtWRq","priority":0,"proof_of_work_nonce":"fa0bedc335d24740","signature":"sigYdNVjxow5vDvNrNSwBvaWBeYV8gPGvDP5GGmEzamEwUGyNJuNUNyh19Feq1EJPMTdiNEXCU6MoD8WghWQwYVgYoY3Zzmi"},"metadata":{"balance_updates":[{"kind":"freezer","category":"deposits","delegate":"tz1gRa9cb1RKEm2B5HoctkBbovFVLbtwV25f","level":8,"change":"-64000000"},{"kind":"freezer","category":"rewards","delegate":"tz1gRa9cb1RKEm2B5HoctkBbovFVLbtwV25f","level":8,"change":"-16000000"},{"kind":"freezer","category":"rewards","delegate":"tz3VEZ4k6a4Wx42iyev6i2aVAptTRLEAivNN","level":8,"change":"32000000"}]}}]`),
			false,
			Contents{
				DoubleBakingEvidence: []DoubleBakingEvidence{
					{
						Kind: "double_baking_evidence",
						Bh1: &BlockHeader{
							Level:          32958,
							Proto:          2,
							Predecessor:    "BMVMTfywir5E2QucCxwv2DkM2CYtnnqesriXJtegEGqfLRt9Sxh",
							Timestamp:      timeFromStr("2018-07-25T01:36:57Z"),
							ValidationPass: 4,
							OperationsHash: "LLoZfBngaArquerV55FPZsTMbBALJ5Pfrh9KAUzHKEJPYAutMDLcV",
							Fitness: []string{
								"00",
								"000000000009e668",
							},
							Context:          "CoVauxcNbkGuCVE9ajsbVhdDGg9zs7CwvLAjtJ6zpGLaMcWhtWRq",
							Priority:         0,
							ProofOfWorkNonce: "9f41eeca8bab4056",
							Signature:        "sigceQeBtuHbhcvpTM8QXAztU4KzeVKQCNapSA5mqCBc6tVC5FYQBhTAyr75cfFe9J3yX7YH1Yy769tHyHJLG4NCY4Azoujb",
						},
						Bh2: &BlockHeader{
							Level:          32958,
							Proto:          2,
							Predecessor:    "BMVMTfywir5E2QucCxwv2DkM2CYtnnqesriXJtegEGqfLRt9Sxh",
							Timestamp:      timeFromStr("2018-07-25T01:36:57Z"),
							ValidationPass: 4,
							OperationsHash: "LLoZfBngaArquerV55FPZsTMbBALJ5Pfrh9KAUzHKEJPYAutMDLcV",
							Fitness: []string{
								"00",
								"000000000009e668",
							},
							Context:          "CoVauxcNbkGuCVE9ajsbVhdDGg9zs7CwvLAjtJ6zpGLaMcWhtWRq",
							Priority:         0,
							ProofOfWorkNonce: "fa0bedc335d24740",
							Signature:        "sigYdNVjxow5vDvNrNSwBvaWBeYV8gPGvDP5GGmEzamEwUGyNJuNUNyh19Feq1EJPMTdiNEXCU6MoD8WghWQwYVgYoY3Zzmi",
						},
						Metadata: &DoubleBakingEvidenceMetadata{
							BalanceUpdates: []BalanceUpdates{
								{
									Kind:     "freezer",
									Category: "deposits",
									Delegate: "tz1gRa9cb1RKEm2B5HoctkBbovFVLbtwV25f",
									Level:    8,
									Change:   -64000000,
								},
								{
									Kind:     "freezer",
									Category: "rewards",
									Delegate: "tz1gRa9cb1RKEm2B5HoctkBbovFVLbtwV25f",
									Level:    8,
									Change:   -16000000,
								},
								{
									Kind:     "freezer",
									Category: "rewards",
									Delegate: "tz3VEZ4k6a4Wx42iyev6i2aVAptTRLEAivNN",
									Level:    8,
									Change:   32000000,
								},
							},
						},
					},
				},
			},
		},
		{
			"successfully unmarshals and marshals double endorsing evidence",
			[]byte(`[{"kind":"double_endorsement_evidence","op1":{"branch":"BLyQHMFeNzZEKHmKgfD9imcowLm8hc4aUo16QtYZcS5yvx7RFqQ","operations":{"kind":"endorsement","level":554811},"signature":"sigqgQgW5qQCsuHP5HhMhAYR2HjcChUE7zAczsyCdF681rfZXpxnXFHu3E6ycmz4pQahjvu3VLfa7FMCxZXmiMiuZFQS4MHy"},"op2":{"branch":"BLTfU3iAfPFMuHTmC1F122AHqdhqnFTfkxBmzYCWtCkBMpYNjxw","operations":{"kind":"endorsement","level":554811},"signature":"sigPwkrKhsDdEidvvUgEEtsaVhyiGmzhCYqCJGKqbYMtH8KxkrFds2HmpDCpRxSTnehKoSC8XKCs9eej6PEzcZoy6fqRAPEZ"},"metadata":{"balance_updates":[{"kind":"freezer","category":"deposits","delegate":"tz1PeZx7FXy7QRuMREGXGxeipb24RsMMzUNe","cycle":135,"change":"-38656000000"},{"kind":"freezer","category":"fees","delegate":"tz1PeZx7FXy7QRuMREGXGxeipb24RsMMzUNe","cycle":135,"change":"-87580"},{"kind":"freezer","category":"rewards","delegate":"tz1PeZx7FXy7QRuMREGXGxeipb24RsMMzUNe","cycle":135,"change":"-1190166666"},{"kind":"freezer","category":"rewards","delegate":"tz1gk3TDbU7cJuiBRMhwQXVvgDnjsxuWhcEA","cycle":135,"change":"19328043790"}]}}]`),
			false,
			Contents{
				DoubleEndorsementEvidence: []DoubleEndorsementEvidence{
					{
						Kind: "double_endorsement_evidence",
						Op1: &InlinedEndorsement{
							Branch: "BLyQHMFeNzZEKHmKgfD9imcowLm8hc4aUo16QtYZcS5yvx7RFqQ",
							Operations: &InlinedEndorsementOperations{
								Kind:  "endorsement",
								Level: 554811,
							},
							Signature: "sigqgQgW5qQCsuHP5HhMhAYR2HjcChUE7zAczsyCdF681rfZXpxnXFHu3E6ycmz4pQahjvu3VLfa7FMCxZXmiMiuZFQS4MHy",
						},
						Op2: &InlinedEndorsement{
							Branch: "BLTfU3iAfPFMuHTmC1F122AHqdhqnFTfkxBmzYCWtCkBMpYNjxw",
							Operations: &InlinedEndorsementOperations{
								Kind:  "endorsement",
								Level: 554811,
							},
							Signature: "sigPwkrKhsDdEidvvUgEEtsaVhyiGmzhCYqCJGKqbYMtH8KxkrFds2HmpDCpRxSTnehKoSC8XKCs9eej6PEzcZoy6fqRAPEZ",
						},
						Metadata: &DoubleEndorsementEvidenceMetadata{
							BalanceUpdates: []BalanceUpdates{
								{
									Kind:     "freezer",
									Category: "deposits",
									Delegate: "tz1PeZx7FXy7QRuMREGXGxeipb24RsMMzUNe",
									Cycle:    135,
									Change:   -38656000000,
								},
								{
									Kind:     "freezer",
									Category: "fees",
									Delegate: "tz1PeZx7FXy7QRuMREGXGxeipb24RsMMzUNe",
									Cycle:    135,
									Change:   -87580,
								},
								{
									Kind:     "freezer",
									Category: "rewards",
									Delegate: "tz1PeZx7FXy7QRuMREGXGxeipb24RsMMzUNe",
									Cycle:    135,
									Change:   -1190166666,
								},
								{
									Kind:     "freezer",
									Category: "rewards",
									Delegate: "tz1gk3TDbU7cJuiBRMhwQXVvgDnjsxuWhcEA",
									Cycle:    135,
									Change:   19328043790,
								},
							},
						},
					},
				},
			},
		},
		{
			"successfully unmarshals and marshals seed nonce revelation",
			[]byte(`[{"kind":"seed_nonce_revelation","level":3872,"nonce":"8202758af02e67400ffa9fa00b673876fd1797e18d17a73fe2c207f8623a3258","metadata":{"balance_updates":[{"kind":"freezer","category":"rewards","delegate":"tz3WMqdzXqRWXwyvj5Hp2H7QEepaUuS7vd9K","level":0,"change":"125000"}]}}]`),
			false,
			Contents{
				SeedNonceRevelations: []SeedNonceRevelation{
					{
						Kind:  "seed_nonce_revelation",
						Level: 3872,
						Nonce: "8202758af02e67400ffa9fa00b673876fd1797e18d17a73fe2c207f8623a3258",

						Metadata: &SeedNonceRevelationMetadata{
							BalanceUpdates: []BalanceUpdates{
								{
									Kind:     "freezer",
									Category: "rewards",
									Delegate: "tz3WMqdzXqRWXwyvj5Hp2H7QEepaUuS7vd9K",
									Level:    0,
									Change:   125000,
								},
							},
						},
					},
				},
			},
		},
		{
			"successfully unmarshals and marshals origination",
			[]byte(`[{"kind":"origination","source":"tz1Wit2PqodvPeuRRhdQXmkrtU8e8bRYZecd","fee":"0","counter":"24","gas_limit":"0","storage_limit":"0","managerPubkey":"tz1Wit2PqodvPeuRRhdQXmkrtU8e8bRYZecd","balance":"10000000","delegate":"tz1Wit2PqodvPeuRRhdQXmkrtU8e8bRYZecd","metadata":{"balance_updates":[],"operation_result":{"status":"applied","balance_updates":[{"kind":"contract","contract":"tz1Wit2PqodvPeuRRhdQXmkrtU8e8bRYZecd","change":"-257000"},{"kind":"contract","contract":"tz1Wit2PqodvPeuRRhdQXmkrtU8e8bRYZecd","change":"-10000000"},{"kind":"contract","contract":"KT1EC9cuswn2quyhM77iWpngTju1BjqLSd5u","change":"10000000"}],"originated_contracts":["KT1EC9cuswn2quyhM77iWpngTju1BjqLSd5u"]}}}]`),
			false,
			Contents{
				Originations: []Origination{
					{
						Kind:          "origination",
						Source:        "tz1Wit2PqodvPeuRRhdQXmkrtU8e8bRYZecd",
						Fee:           0,
						Counter:       24,
						GasLimit:      0,
						StorageLimit:  0,
						ManagerPubkey: "tz1Wit2PqodvPeuRRhdQXmkrtU8e8bRYZecd",
						Balance:       10000000,
						Delegate:      "tz1Wit2PqodvPeuRRhdQXmkrtU8e8bRYZecd",
						Metadata: &OriginationMetadata{
							BalanceUpdates: []BalanceUpdates{},
							OperationResults: OperationResultOrigination{
								Status: "applied",
								BalanceUpdates: []BalanceUpdates{
									{
										Kind:     "contract",
										Contract: "tz1Wit2PqodvPeuRRhdQXmkrtU8e8bRYZecd",
										Change:   -257000,
									},
									{
										Kind:     "contract",
										Contract: "tz1Wit2PqodvPeuRRhdQXmkrtU8e8bRYZecd",
										Change:   -10000000,
									},
									{
										Kind:     "contract",
										Contract: "KT1EC9cuswn2quyhM77iWpngTju1BjqLSd5u",
										Change:   10000000,
									},
								},
								OriginatedContracts: []string{
									"KT1EC9cuswn2quyhM77iWpngTju1BjqLSd5u",
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			var contents Contents
			err := contents.UnmarshalJSON(tt.contents)
			checkErr(t, tt.wantErr, "", err)
			assert.Equal(t, tt.want, contents)

			v, err := contents.MarshalJSON()
			checkErr(t, tt.wantErr, "", err)

			b, err := tt.want.MarshalJSON()
			checkErr(t, tt.wantErr, "", err)

			assert.Equal(t, string(b), string(v))
		})
	}
}

func Test_MichelineMichelsonV1Expression_JSON(t *testing.T) {
	cases := []struct {
		name      string
		expresion []byte
		wantErr   bool
		want      MichelineExpression
	}{
		{
			"is successful with string",
			[]byte(`{"string": "Test"}`),
			false,
			MichelineExpression{
				Object: &Micheline{
					String: "Test",
				},
			},
		},
		{
			"is successful with nested expression",
			[]byte(`[
				{
				  "prim": "parameter",
				  "args": [
					{
					  "prim": "string"
					}
				  ]
				},
				{
				  "prim": "storage",
				  "args": [
					{
					  "prim": "string"
					}
				  ]
				}]`),
			false,
			MichelineExpression{
				Array: []Micheline{
					{
						Prim: "parameter",
						Args: &MichelineArgs{
							Array: []Micheline{
								{
									Prim: "string",
								},
							},
						},
					},
					{
						Prim: "storage",
						Args: &MichelineArgs{
							Array: []Micheline{
								{
									Prim: "string",
								},
							},
						},
					},
				},
			},
		},
		{
			"is successful with complex micheline",
			[]byte(`{
				"prim": "parameter",
				"args": [
				  {
					"prim": "or",
					"args": [
					  {
						"prim": "or",
						"args": [
						  {
							"prim": "unit",
							"annots": [
							  "%payoff"
							]
						  },
						  {
							"prim": "unit",
							"annots": [
							  "%refund"
							]
						  }
						]
					  },
					  {
						"prim": "unit",
						"annots": [
						  "%sendFund"
						]
					  }
					]
				  }
				]
			  }`),
			false,
			MichelineExpression{
				Object: &Micheline{
					Prim: "parameter",
					Args: &MichelineArgs{
						Array: []Micheline{
							{
								Prim: "or",
								Args: &MichelineArgs{
									Array: []Micheline{
										{
											Prim: "or",
											Args: &MichelineArgs{
												Array: []Micheline{
													{
														Prim: "unit",
														Annots: []string{
															"%payoff",
														},
													},
													{
														Prim: "unit",
														Annots: []string{
															"%refund",
														},
													},
												},
											},
										},
										{
											Prim: "unit",
											Annots: []string{
												"%sendFund",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			"is successful with more complex micheline",
			[]byte(`[
				{
				  "prim": "parameter",
				  "args": [
					{
					  "prim": "string"
					}
				  ]
				},
				{
				  "prim": "storage",
				  "args": [
					{
					  "prim": "string"
					}
				  ]
				},
				{
				  "prim": "code",
				  "args": [
					[
					  {
						"prim": "CAR"
					  },
					  {
						"prim": "NIL",
						"args": [
						  {
							"prim": "operation"
						  }
						]
					  },
					  {
						"prim": "PAIR"
					  }
					]
				  ]
				}
			  ]`),
			false,
			MichelineExpression{
				Array: []Micheline{
					{
						Prim: "parameter",
						Args: &MichelineArgs{
							Array: []Micheline{
								{
									Prim: "string",
								},
							},
						},
					},
					{
						Prim: "storage",
						Args: &MichelineArgs{
							Array: []Micheline{
								{
									Prim: "string",
								},
							},
						},
					},
					{
						Prim: "code",
						Args: &MichelineArgs{
							MultiArray: [][]Micheline{
								{
									{
										Prim: "CAR",
									},
									{
										Prim: "NIL",
										Args: &MichelineArgs{
											Array: []Micheline{
												{
													Prim: "operation",
												},
											},
										},
									},
									{
										Prim: "PAIR",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			exp := MichelineExpression{}
			err := json.Unmarshal(tt.expresion, &exp)
			checkErr(t, false, "", err)
			assert.Equal(t, tt.want, exp)
			v, _ := exp.MarshalJSON()
			assert.Equal(t, stripString(string(tt.expresion)), string(v))
		})
	}
}

func strToPointer(str string) *string {
	return &str
}

func intToPointer(i int) *int {
	return &i
}

func stripString(str string) string {
	str = strings.Replace(string(str), "\t", "", -1)
	str = strings.Replace(string(str), "\n", "", -1)
	str = strings.Replace(string(str), " ", "", -1)
	return str
}

func decodeHexString(str string) []byte {
	v, _ := hex.DecodeString(str)
	return v
}

func timeFromStr(str string) time.Time {
	t, _ := time.Parse("2006-01-02T15:04:05Z", "2018-07-25T01:36:57Z")
	return t
}
