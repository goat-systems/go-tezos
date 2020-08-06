package gotezos

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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
						BigMap:  *NewInt(52),
						KeyHash: "exprta5PGni3vkj7z6B5CHRELDe796kyPq7q9qAqzadnm3fr4AvNhJ",
						Key: MichelineMichelsonV1Expression{
							Int:                            nil,
							String:                         strToPointer("6238d74df3089fe8b263422eea4f35101aa2b8bb50687aa98bdb15e1111b909d"),
							Bytes:                          nil,
							MichelineMichelsonV1Expression: nil,
						},
						Value: MichelineMichelsonV1Expression{
							Prim: "Pair",
							Args: []MichelineMichelsonV1Expression{
								{
									Int: strToPointer("1593806466"),
								},
								{
									Bytes: []byte("00004cc5b68779c9166b20f6aed04f6fb7b01929ab9a"),
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
									Change:   NewInt(-64000000),
								},
								{
									Kind:     "freezer",
									Category: "deposits",
									Delegate: "tz1iZEKy4LaAjnTmn2RuGDf2iqdAQKnRi8kY",
									Cycle:    204,
									Change:   NewInt(64000000),
								},
								{
									Kind:     "freezer",
									Category: "rewards",
									Delegate: "tz1iZEKy4LaAjnTmn2RuGDf2iqdAQKnRi8kY",
									Cycle:    204,
									Change:   NewInt(2000000),
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
			[]byte(`[{"kind":"transaction","source":"tz1Z1tMai15JWUWeN2PKL9faXXVPMuWamzJj","fee":"1792","counter":"84322","gas_limit":"15385","storage_limit":"0","amount":"24278768","destination":"KT1VFrVbFaK9YUy8ZDj49XmFkB3ZwvZeZTqi","metadata":{"balance_updates":[{"kind":"contract","contract":"tz1Z1tMai15JWUWeN2PKL9faXXVPMuWamzJj","change":"-1792"},{"kind":"freezer","category":"fees","delegate":"tz3adcvQaKXTCg12zbninqo3q8ptKKtDFTLv","cycle":205,"change":"1792"}],"operation_result":{"status":"applied","storage":{"bytes":"00984a5af599f114685322422940df9414ad551ee3"},"balance_updates":[{"kind":"contract","contract":"tz1Z1tMai15JWUWeN2PKL9faXXVPMuWamzJj","change":"-24278768"},{"kind":"contract","contract":"KT1VFrVbFaK9YUy8ZDj49XmFkB3ZwvZeZTqi","change":"24278768"}],"consumed_gas":"15285","storage_size":"232"}}}]`),
			false,
			Contents{
				Transactions: []Transaction{
					{
						Kind:         "transaction",
						Source:       "tz1Z1tMai15JWUWeN2PKL9faXXVPMuWamzJj",
						Fee:          NewInt(1792),
						Counter:      84322,
						GasLimit:     NewInt(15385),
						StorageLimit: NewInt(0),
						Amount:       NewInt(24278768),
						Destination:  "KT1VFrVbFaK9YUy8ZDj49XmFkB3ZwvZeZTqi",
						Metadata: &TransactionMetadata{
							BalanceUpdates: []BalanceUpdates{
								{
									Kind:     "contract",
									Contract: "tz1Z1tMai15JWUWeN2PKL9faXXVPMuWamzJj",
									Change:   NewInt(-1792),
								},
								{
									Kind:     "freezer",
									Category: "fees",
									Delegate: "tz3adcvQaKXTCg12zbninqo3q8ptKKtDFTLv",
									Cycle:    205,
									Change:   NewInt(1792),
								},
							},
							OperationResults: OperationResultTransfer{
								Status: "applied",
								Storage: &MichelineMichelsonV1Expression{
									Bytes: []byte(`00984a5af599f114685322422940df9414ad551ee3`),
								},
								BalanceUpdates: []BalanceUpdates{
									{
										Kind:     "contract",
										Contract: "tz1Z1tMai15JWUWeN2PKL9faXXVPMuWamzJj",
										Change:   NewInt(-24278768),
									},
									{
										Kind:     "contract",
										Contract: "KT1VFrVbFaK9YUy8ZDj49XmFkB3ZwvZeZTqi",
										Change:   NewInt(24278768),
									},
								},
								ConsumedGas: NewInt(15285),
								StorageSize: NewInt(232),
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

			v, err := contents.MarshalJSON()
			checkErr(t, tt.wantErr, "", err)

			b, err := tt.want.MarshalJSON()
			checkErr(t, tt.wantErr, "", err)

			assert.Equal(t, string(b), string(v))
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
