package rpc_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goat-systems/go-tezos/v4/rpc"
	"github.com/stretchr/testify/assert"
)

func Test_BakingRights(t *testing.T) {
	goldenBakingRights := getResponse(bakingrights).([]rpc.BakingRights)

	type want struct {
		wantErr          bool
		containsErr      string
		wantBakingRights []rpc.BakingRights
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regBakingRights, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get baking rights",
				[]rpc.BakingRights{},
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regBakingRights, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get baking rights: failed to parse json",
				[]rpc.BakingRights{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regBakingRights, readResponse(bakingrights)}, blankHandler)),
			want{
				false,
				"",
				goldenBakingRights,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			r, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, bakingRights, err := r.BakingRights(rpc.BakingRightsInput{
				BlockID: &rpc.BlockIDHead{},
			})
			checkErr(t, tt.wantErr, tt.containsErr, err)
			assert.Equal(t, tt.want.wantBakingRights, bakingRights)
		})
	}
}

func Test_CompletePrefix(t *testing.T) {
	goldenCurrentLevel := getResponse(currentLevel).(rpc.CurrentLevel)

	type want struct {
		wantErr          bool
		containsErr      string
		wantCurrentLevel rpc.CurrentLevel
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regCurrentLevel, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get current level",
				rpc.CurrentLevel{},
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regCurrentLevel, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get current level: failed to parse json",
				rpc.CurrentLevel{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regCurrentLevel, readResponse(currentLevel)}, blankHandler)),
			want{
				false,
				"",
				goldenCurrentLevel,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			r, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, currentLevel, err := r.CurrentLevel(rpc.CurrentLevelInput{
				BlockID: &rpc.BlockIDHead{},
			})
			checkErr(t, tt.wantErr, tt.containsErr, err)
			assert.Equal(t, tt.want.wantCurrentLevel, currentLevel)
		})
	}
}

func Test_EndorsingRights(t *testing.T) {
	goldenEndorsingRights := getResponse(endorsingrights).([]rpc.EndorsingRights)

	type want struct {
		wantErr             bool
		containsErr         string
		wantEndorsingRights []rpc.EndorsingRights
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regEndorsingRights, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get endorsing rights",
				[]rpc.EndorsingRights{},
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regEndorsingRights, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get endorsing rights: failed to parse json",
				[]rpc.EndorsingRights{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regEndorsingRights, readResponse(endorsingrights)}, blankHandler)),
			want{
				false,
				"",
				goldenEndorsingRights,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			r, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, endorsingRights, err := r.EndorsingRights(rpc.EndorsingRightsInput{
				BlockID: &rpc.BlockIDHead{},
			})
			checkErr(t, tt.wantErr, tt.containsErr, err)
			assert.Equal(t, tt.want.wantEndorsingRights, endorsingRights)
		})
	}
}

func Test_ForgeOperations(t *testing.T) {
	goldenHash := rpc.BlockIDHash(mockBlockHash)
	goldenOperationBytes := []byte(`"a79ec80dba1f8ddb2cde90b8f12f7c62fdc36556030281ff8904a3d0df82cddc08000008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e00b960000008ba0cb2fad622697145cf1665124096d25bc31e00"`)

	type want struct {
		err         bool
		errContains string
		operation   string
	}

	type input struct {
		handler http.Handler
		i       rpc.ForgeOperationsInput
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"handles rpc failure",
			input{
				gtGoldenHTTPMock(mockHandler(&requestResultPair{regForgeOperationWithRPC, readResponse(rpcerrors)}, blankHandler)),
				rpc.ForgeOperationsInput{
					BlockIDHash: goldenHash,
					Branch:      "some_branch",
					Contents:    rpc.Contents{},
				},
			},
			want{
				true,
				"failed to forge operation: rpc error (somekind)",
				"",
			},
		},
		{
			"handles failure to unmarshal",
			input{
				gtGoldenHTTPMock(mockHandler(&requestResultPair{regForgeOperationWithRPC, []byte(`junk`)}, blankHandler)),
				rpc.ForgeOperationsInput{
					BlockIDHash: goldenHash,
					Branch:      "some_branch",
					Contents:    rpc.Contents{},
				},
			},
			want{
				true,
				"failed to forge operation: invalid character",
				"",
			},
		},
		{
			"handles failure to strip operation branch",
			input{
				gtGoldenHTTPMock(mockHandler(&requestResultPair{regForgeOperationWithRPC, []byte(`"some_junk_op_string"`)}, mockHandler(&requestResultPair{regParseOperations, readResponse(rpcerrors)}, blankHandler))),
				rpc.ForgeOperationsInput{
					BlockIDHash: goldenHash,
					Branch:      "some_branch",
					Contents:    rpc.Contents{},
				},
			},
			want{
				true,
				"failed to forge operation: unable to verify rpc returned a valid contents",
				"some_junk_op_string",
			},
		},
		{
			"handles failure to parse forged operation",
			input{
				gtGoldenHTTPMock(mockHandler(&requestResultPair{regForgeOperationWithRPC, []byte(`"some_operation_string"`)}, mockHandler(&requestResultPair{regParseOperations, readResponse(rpcerrors)}, blankHandler))),
				rpc.ForgeOperationsInput{
					BlockIDHash: goldenHash,
					Branch:      "some_branch",
					Contents:    rpc.Contents{},
				},
			},
			want{
				true,
				"failed to forge operation: unable to verify rpc returned a valid contents",
				"some_operation_string",
			},
		},
		{
			"handles failure to match forge with expected contents",
			input{
				gtGoldenHTTPMock(mockHandler(&requestResultPair{regForgeOperationWithRPC, goldenOperationBytes}, mockHandler(&requestResultPair{regParseOperations, readResponse(parseOperations)}, blankHandler))),
				rpc.ForgeOperationsInput{
					BlockIDHash: goldenHash,
					Branch:      "some_branch",
					Contents: rpc.Contents{
						{
							Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
							Fee:          "100",
							Counter:      "10",
							GasLimit:     "10100",
							StorageLimit: "0",
							Amount:       "12345",
							Destination:  "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
							Kind:         rpc.TRANSACTION,
						},
					},
				},
			},
			want{
				true,
				"failed to forge operation: alert rpc returned invalid contents",
				"a79ec80dba1f8ddb2cde90b8f12f7c62fdc36556030281ff8904a3d0df82cddc08000008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e00b960000008ba0cb2fad622697145cf1665124096d25bc31e00",
			},
		},
		{
			"is successful",
			input{
				gtGoldenHTTPMock(mockHandler(&requestResultPair{regForgeOperationWithRPC, goldenOperationBytes}, mockHandler(&requestResultPair{regParseOperations, readResponse(parseOperations)}, blankHandler))),
				rpc.ForgeOperationsInput{
					BlockIDHash: goldenHash,
					Branch:      "some_branch",
					Contents: rpc.Contents{
						{
							Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
							Fee:          "10100",
							Counter:      "10",
							GasLimit:     "10100",
							StorageLimit: "0",
							Amount:       "12345",
							Destination:  "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
							Kind:         rpc.TRANSACTION,
						},
					},
				},
			},
			want{
				false,
				"",
				"a79ec80dba1f8ddb2cde90b8f12f7c62fdc36556030281ff8904a3d0df82cddc08000008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e00b960000008ba0cb2fad622697145cf1665124096d25bc31e00",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input.handler)
			defer server.Close()

			r, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, op, err := r.ForgeOperations(tt.input.i)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.operation, op)
		})
	}
}

func Test_ForgeBlockHeader(t *testing.T) {
	type want struct {
		wantErr              bool
		containsErr          string
		wantForgeBlockHeader rpc.ForgeBlockHeader
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regForgeBlockHeader, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to forge block header",
				rpc.ForgeBlockHeader{},
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regForgeBlockHeader, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to forge block header: failed to parse json",
				rpc.ForgeBlockHeader{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regForgeBlockHeader, []byte(`{ "block": "fdasjlfh" }`)}, blankHandler)),
			want{
				false,
				"",
				rpc.ForgeBlockHeader{
					Block: "fdasjlfh",
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			r, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, blockHeader, err := r.ForgeBlockHeader(rpc.ForgeBlockHeaderInput{
				BlockID:     &rpc.BlockIDHead{},
				BlockHeader: rpc.ForgeBlockHeaderBody{},
			})
			checkErr(t, tt.wantErr, tt.containsErr, err)
			assert.Equal(t, tt.want.wantForgeBlockHeader, blockHeader)
		})
	}
}

func Test_LevelsInCurrentCycle(t *testing.T) {
	type want struct {
		wantErr                  bool
		containsErr              string
		wantLevelsInCurrentCycle rpc.LevelsInCurrentCycle
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regLevelsInCurrentCycle, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get levels in current cycle",
				rpc.LevelsInCurrentCycle{},
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regLevelsInCurrentCycle, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get levels in current cycle: failed to parse json",
				rpc.LevelsInCurrentCycle{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regLevelsInCurrentCycle, []byte(`{"first": 1273857,"last": 1277952}`)}, blankHandler)),
			want{
				false,
				"",
				rpc.LevelsInCurrentCycle{
					First: 1273857,
					Last:  1277952,
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			r, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, levelsInCurrentCycle, err := r.LevelsInCurrentCycle(rpc.LevelsInCurrentCycleInput{
				BlockID: &rpc.BlockIDHead{},
			})
			checkErr(t, tt.wantErr, tt.containsErr, err)
			assert.Equal(t, tt.want.wantLevelsInCurrentCycle, levelsInCurrentCycle)
		})
	}
}

func Test_ParseBlock(t *testing.T) {
	type want struct {
		wantErr                       bool
		containsErr                   string
		wantBlockHeaderSignedContents rpc.BlockHeaderSignedContents
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regParseBlock, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to parse block",
				rpc.BlockHeaderSignedContents{},
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regParseBlock, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to parse block: failed to parse json",
				rpc.BlockHeaderSignedContents{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regParseBlock, []byte(`{ "priority": 100, "proof_of_work_nonce": "some_nonce", "seed_nonce_hash": "some_nonce", "signature": "sig" }`)}, blankHandler)),
			want{
				false,
				"",
				rpc.BlockHeaderSignedContents{
					Priority:         100,
					ProofOfWorkNonce: "some_nonce",
					SeedNonceHash:    "some_nonce",
					Signature:        "sig",
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			r, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, parsedBlock, err := r.ParseBlock(rpc.ParseBlockInput{
				BlockID:     &rpc.BlockIDHead{},
				BlockHeader: rpc.ForgeBlockHeaderBody{},
			})
			checkErr(t, tt.wantErr, tt.containsErr, err)
			assert.Equal(t, tt.want.wantBlockHeaderSignedContents, parsedBlock)
		})
	}
}

func Test_ParseOperations(t *testing.T) {
	type input struct {
		handler http.Handler
		i       rpc.ParseOperationsInput
	}

	type want struct {
		err         bool
		errContains string
		operations  []rpc.Operations
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"handles rpc error",
			input{
				gtGoldenHTTPMock(mockHandler(&requestResultPair{regParseOperations, readResponse(rpcerrors)}, blankHandler)),
				rpc.ParseOperationsInput{
					BlockID: &rpc.BlockIDHead{},
					Operations: []rpc.ParseOperationsBody{
						{
							Data:   "some_data",
							Branch: "some_branch",
						},
					},
				},
			},
			want{
				true,
				"failed to parse operations",
				[]rpc.Operations{},
			},
		},
		{
			"handles failure to unmarshal",
			input{
				gtGoldenHTTPMock(mockHandler(&requestResultPair{regParseOperations, []byte(`junk`)}, blankHandler)),
				rpc.ParseOperationsInput{
					BlockID: &rpc.BlockIDHead{},
					Operations: []rpc.ParseOperationsBody{
						{
							Data:   "some_data",
							Branch: "some_branch",
						},
					},
				},
			},
			want{
				true,
				"failed to parse operations: failed to parse json",
				[]rpc.Operations{},
			},
		},
		{
			"is successful",
			input{
				gtGoldenHTTPMock(mockHandler(&requestResultPair{regParseOperations, readResponse(parseOperations)}, blankHandler)),
				rpc.ParseOperationsInput{
					BlockID: &rpc.BlockIDHead{},
					Operations: []rpc.ParseOperationsBody{
						{
							Data:   "some_data",
							Branch: "some_branch",
						},
					},
				},
			},
			want{
				false,
				"",
				[]rpc.Operations{
					{
						Protocol: "",
						ChainID:  "",
						Hash:     "",
						Branch:   "BLz6yCE4BUL4ppo1zsEWdK9FRCt15WAY7ECQcuK9RtWg4xeEVL7",
						Contents: rpc.Contents{
							{
								Kind:         "transaction",
								Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
								Fee:          "10100",
								Counter:      "10",
								GasLimit:     "10100",
								StorageLimit: "0",
								Amount:       "12345",
								Destination:  "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
								Delegate:     "",
								Secret:       "",
								Level:        0,
								Period:       0,
								Proposal:     "",
								Proposals:    []string(nil),
								Ballot:       "",
								Metadata:     nil,
							},
						},
						Signature: "edsigtXomBKi5CTRf5cjATJWSyaRvhfYNHqSUGrn4SdbYRcGwQrUGjzEfQDTuqHhuA8b2d8NarZjz8TRf65WkpQmo423BtomS8Q",
					},
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input.handler)
			defer server.Close()

			r, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, op, err := r.ParseOperations(tt.input.i)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.operations, op)
		})
	}
}

func Test_PreapplyBlock(t *testing.T) {
	type want struct {
		wantErr             bool
		containsErr         string
		wantPreappliedBlock rpc.PreappliedBlock
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regPreapplyBlock, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to preapply block",
				rpc.PreappliedBlock{},
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regPreapplyBlock, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to preapply block: failed to parse json",
				rpc.PreappliedBlock{},
			},
		},
		// TODO Get mock data for successful run
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			r, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, preappliedBlock, err := r.PreapplyBlock(rpc.PreapplyBlockInput{
				BlockID: &rpc.BlockIDHead{},
				Block:   rpc.PreapplyBlockBody{},
			})
			checkErr(t, tt.wantErr, tt.containsErr, err)
			assert.Equal(t, tt.want.wantPreappliedBlock, preappliedBlock)
		})
	}
}

func Test_PreapplyOperation(t *testing.T) {
	type input struct {
		handler http.Handler
		i       rpc.PreapplyOperationsInput
	}

	type want struct {
		err         bool
		errContains string
		operations  []rpc.Operations
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"handles rpc error",
			input{
				gtGoldenHTTPMock(mockHandler(&requestResultPair{regPreapplyOperations, readResponse(rpcerrors)}, blankHandler)),
				rpc.PreapplyOperationsInput{
					BlockID: &rpc.BlockIDHead{},
					Operations: []rpc.Operations{
						{
							Protocol:  "some_protocol",
							Signature: "some_sig",
							Contents:  rpc.Contents{},
						},
					},
				},
			},
			want{
				true,
				"failed to preapply operations",
				nil,
			},
		},
		{
			"handles failure to unmarshal",
			input{
				gtGoldenHTTPMock(mockHandler(&requestResultPair{regPreapplyOperations, []byte("junk")}, blankHandler)),
				rpc.PreapplyOperationsInput{
					BlockID: &rpc.BlockIDHead{},
					Operations: []rpc.Operations{
						{
							Protocol:  "some_protocol",
							Signature: "some_sig",
							Contents:  rpc.Contents{},
						},
					},
				},
			},
			want{
				true,
				"failed to preapply operations: failed to parse json",
				nil,
			},
		},
		{
			"is successful",
			input{
				gtGoldenHTTPMock(mockHandler(&requestResultPair{regPreapplyOperations, readResponse(preapplyOperations)}, blankHandler)),
				rpc.PreapplyOperationsInput{
					BlockID: &rpc.BlockIDHead{},
					Operations: []rpc.Operations{
						{
							Protocol:  "some_protocol",
							Signature: "some_sig",
							Contents:  rpc.Contents{},
						},
					},
				},
			},
			want{
				false,
				"",
				[]rpc.Operations{
					{
						Contents: rpc.Contents{
							{
								Kind:         "transaction",
								Source:       "tz1W3HW533csCBLor4NPtU79R2TT2sbKfJDH",
								Fee:          "3000",
								Counter:      "1263232",
								GasLimit:     "20000",
								StorageLimit: "0",
								Amount:       "50",
								Destination:  "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
								Metadata: &rpc.ContentsMetadata{
									BalanceUpdates: []rpc.BalanceUpdates{
										{
											Kind:     "contract",
											Contract: "tz1W3HW533csCBLor4NPtU79R2TT2sbKfJDH",
											Change:   "-3000",
										},
										{
											Kind:     "freezer",
											Category: "fees",
											Delegate: "tz1Ke2h7sDdakHJQh8WX4Z372du1KChsksyU",
											Cycle:    229,
											Change:   "3000",
										},
									},
									OperationResults: &rpc.OperationResults{
										Status: "applied",
										BalanceUpdates: []rpc.BalanceUpdates{
											{
												Kind:     "contract",
												Contract: "tz1W3HW533csCBLor4NPtU79R2TT2sbKfJDH",
												Change:   "-50",
											},
											{
												Kind:     "contract",
												Contract: "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
												Change:   "50",
											},
										},
										ConsumedGas: "10207",
									},
								},
							},
						},
						Signature: "edsig...."},
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input.handler)
			defer server.Close()

			r, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, operations, err := r.PreapplyOperations(tt.input.i)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.operations, operations)
		})
	}
}

func Test_Entrypoint(t *testing.T) {
	goldenEntrypointData := []byte(`{"int": 10}`)
	m := &json.RawMessage{}
	*m = goldenEntrypointData

	type want struct {
		wantErr        bool
		containsErr    string
		wantEntrypoint rpc.Entrypoint
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regEntrypoint, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get entrypoint type:",
				rpc.Entrypoint{},
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regEntrypoint, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get entrypoint type: failed to parse json",
				rpc.Entrypoint{},
			},
		},
		{
			// TODO: Get real mock data
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regEntrypoint, []byte(`{"entrypoint_type": {"int": 10}}`)}, blankHandler)),
			want{
				false,
				"",
				rpc.Entrypoint{
					EntrypointType: m,
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			r, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, entrypoint, err := r.Entrypoint(rpc.EntrypointInput{
				BlockID:    &rpc.BlockIDHead{},
				Entrypoint: rpc.EntrypointBody{},
			})
			checkErr(t, tt.wantErr, tt.containsErr, err)
			assert.Equal(t, tt.want.wantEntrypoint, entrypoint)
		})
	}
}

func Test_Entrypoints(t *testing.T) {
	goldenEntrypointData := []byte(`{"int": 10}`)
	m := &json.RawMessage{}
	*m = goldenEntrypointData

	type want struct {
		wantErr         bool
		containsErr     string
		wantEntrypoints rpc.Entrypoints
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regEntrypoints, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get entrypoints:",
				rpc.Entrypoints{},
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regEntrypoints, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get entrypoints: failed to parse json",
				rpc.Entrypoints{},
			},
		},
		{
			// TODO: Get real mock data
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regEntrypoints, []byte(`{"entrypoints": {"int": 10}}`)}, blankHandler)),
			want{
				false,
				"",
				rpc.Entrypoints{
					EntrypointsFromScript: m,
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			r, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, entrypoints, err := r.Entrypoints(rpc.EntrypointsInput{
				BlockID:     &rpc.BlockIDHead{},
				Entrypoints: rpc.EntrypointsBody{},
			})
			checkErr(t, tt.wantErr, tt.containsErr, err)
			assert.Equal(t, tt.want.wantEntrypoints, entrypoints)
		})
	}
}

func Test_PackData(t *testing.T) {
	type want struct {
		wantErr        bool
		containsErr    string
		wantPackedData rpc.PackedData
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regPackData, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to pack data:",
				rpc.PackedData{},
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regPackData, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to pack data: failed to parse json",
				rpc.PackedData{},
			},
		},
		{
			// TODO: Get real mock data
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regPackData, []byte(`{"packed": "exprupozG51AtT7yZUy5sg6VbJQ4b9omAE1PKD2PXvqi2YBuZqoKG3", "gas":"1000"}`)}, blankHandler)),
			want{
				false,
				"",
				rpc.PackedData{
					Packed: "exprupozG51AtT7yZUy5sg6VbJQ4b9omAE1PKD2PXvqi2YBuZqoKG3",
					Gas:    "1000",
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			r, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, packedData, err := r.PackData(rpc.PackDataInput{
				BlockID: &rpc.BlockIDHead{},
				Data:    rpc.PackDataBody{},
			})
			checkErr(t, tt.wantErr, tt.containsErr, err)
			assert.Equal(t, tt.want.wantPackedData, packedData)
		})
	}
}

func Test_RunCode(t *testing.T) {
	type want struct {
		wantErr     bool
		containsErr string
		wantRanCode rpc.RanCode
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regRunCode, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to run code:",
				rpc.RanCode{},
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regRunCode, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to run code: failed to parse json",
				rpc.RanCode{},
			},
		},
		// TODO: Get real mock data
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			r, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, rancode, err := r.RunCode(rpc.RunCodeInput{
				BlockID: &rpc.BlockIDHead{},
				Code:    rpc.RunCodeBody{},
			})
			checkErr(t, tt.wantErr, tt.containsErr, err)
			assert.Equal(t, tt.want.wantRanCode, rancode)
		})
	}
}

func Test_RunOperation(t *testing.T) {
	type want struct {
		wantErr        bool
		containsErr    string
		wantOperations rpc.Operations
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regRunOperation, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to run operation:",
				rpc.Operations{},
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regRunOperation, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to run operation: failed to parse json",
				rpc.Operations{},
			},
		},
		// TODO: Get real mock data
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			r, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, operations, err := r.RunOperation(rpc.RunOperationInput{
				BlockID: &rpc.BlockIDHead{},
				Operation: rpc.RunOperation{
					ChainID:   "some_chain_id",
					Operation: rpc.Operations{},
				},
			})
			checkErr(t, tt.wantErr, tt.containsErr, err)
			assert.Equal(t, tt.want.wantOperations, operations)
		})
	}
}

func Test_TraceCode(t *testing.T) {
	type want struct {
		wantErr        bool
		containsErr    string
		wantTracedCode rpc.TracedCode
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regTraceCode, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to trace code",
				rpc.TracedCode{},
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regTraceCode, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to trace code: failed to parse json",
				rpc.TracedCode{},
			},
		},
		// TODO: Get real mock data
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			r, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, tracedCode, err := r.TraceCode(rpc.TraceCodeInput{
				BlockID: &rpc.BlockIDHead{},
				Code:    rpc.RunCodeBody{},
			})
			checkErr(t, tt.wantErr, tt.containsErr, err)
			assert.Equal(t, tt.want.wantTracedCode, tracedCode)
		})
	}
}

func Test_TypecheckCode(t *testing.T) {
	type want struct {
		wantErr             bool
		containsErr         string
		wantTypecheckedCode rpc.TypecheckedCode
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regTypecheckCode, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to typecheck code",
				rpc.TypecheckedCode{},
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regTypecheckCode, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to typecheck code: failed to parse json",
				rpc.TypecheckedCode{},
			},
		},
		// TODO: Get real mock data
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			r, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, typecheckedCode, err := r.TypecheckCode(rpc.TypeCheckcodeInput{
				BlockID: &rpc.BlockIDHead{},
				Code:    rpc.TypecheckCodeBody{},
			})
			checkErr(t, tt.wantErr, tt.containsErr, err)
			assert.Equal(t, tt.want.wantTypecheckedCode, typecheckedCode)
		})
	}
}

func Test_TypecheckData(t *testing.T) {
	type want struct {
		wantErr             bool
		containsErr         string
		wantTypecheckedData rpc.TypecheckedData
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regTypecheckData, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to typecheck data",
				rpc.TypecheckedData{},
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regTypecheckData, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to typecheck data: failed to parse json",
				rpc.TypecheckedData{},
			},
		},
		// TODO: Get real mock data
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			r, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, typecheckedData, err := r.TypecheckData(rpc.TypecheckDataInput{
				BlockID: &rpc.BlockIDHead{},
				Data:    rpc.TypecheckDataBody{},
			})
			checkErr(t, tt.wantErr, tt.containsErr, err)
			assert.Equal(t, tt.want.wantTypecheckedData, typecheckedData)
		})
	}
}
