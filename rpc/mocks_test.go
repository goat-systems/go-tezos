package rpc_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"testing"

	"github.com/completium/go-tezos/v4/rpc"
	"github.com/stretchr/testify/assert"
)

type responseKey string

func (r *responseKey) String() string {
	return string(*r)
}

const (
	activechains            responseKey = ".test-fixtures/active_chains.json"
	bakingrights            responseKey = ".test-fixtures/baking_rights.json"
	balance                 responseKey = ".test-fixtures/balance.json"
	ballotList              responseKey = ".test-fixtures/ballot_list.json"
	ballots                 responseKey = ".test-fixtures/ballots.json"
	block                   responseKey = ".test-fixtures/block.json"
	blocks                  responseKey = ".test-fixtures/blocks.json"
	chainid                 responseKey = ".test-fixtures/chain_id.json"
	contract                responseKey = ".test-fixtures/contract.json"
	connections             responseKey = ".test-fixtures/connections.json"
	constants               responseKey = ".test-fixtures/constants.json"
	contractEntrypoints     responseKey = ".test-fixtures/entrypoints.json"
	counter                 responseKey = ".test-fixtures/counter.json"
	currentLevel            responseKey = ".test-fixtures/current_level.json"
	cycle                   responseKey = ".test-fixtures/cycle.json"
	delegate                responseKey = ".test-fixtures/delegate.json"
	delegatedcontracts      responseKey = ".test-fixtures/delegated_contracts.json"
	endorsingrights         responseKey = ".test-fixtures/endorsing_rights.json"
	frozenbalance           responseKey = ".test-fixtures/frozen_balance.json"
	frozenbalanceByCycle    responseKey = ".test-fixtures/frozen_balance_by_cycle.json"
	header                  responseKey = ".test-fixtures/header.json"
	headerShell             responseKey = ".test-fixtures/header_shell.json"
	liveBlocks              responseKey = ".test-fixtures/live_blocks.json"
	metadata                responseKey = ".test-fixtures/metadata.json"
	operations              responseKey = ".test-fixtures/operations.json"
	operationhashes         responseKey = ".test-fixtures/operation_hashes.json"
	operationMetaDataHashes responseKey = ".test-fixtures/operation_metadata_hashes.json"
	parseOperations         responseKey = ".test-fixtures/parse_operations.json"
	preapplyOperations      responseKey = ".test-fixtures/preapply_operations.json"
	proposals               responseKey = ".test-fixtures/proposals.json"
	protocols               responseKey = ".test-fixtures/protocols.json"
	protocolData            responseKey = ".test-fixtures/protocol_data.json"
	rpcerrors               responseKey = ".test-fixtures/rpc_errors.json"
	voteListings            responseKey = ".test-fixtures/vote_listings.json"
)

func readResponse(key responseKey) []byte {
	f, _ := ioutil.ReadFile(key.String())
	return f
}

func getResponse(key responseKey) interface{} {
	switch key {
	case activechains:
		f := readResponse(key)
		var out rpc.ActiveChains
		json.Unmarshal(f, &out)
		return out
	case bakingrights:
		f := readResponse(key)
		var out []rpc.BakingRights
		json.Unmarshal(f, &out)
		return out
	case balance:
		f := readResponse(key)
		var out int
		json.Unmarshal(f, &out)
		return &out
	case ballotList:
		f := readResponse(key)
		var out rpc.BallotList
		json.Unmarshal(f, &out)
		return &out
	case ballots:
		f := readResponse(key)
		var out rpc.Ballots
		json.Unmarshal(f, &out)
		return &out
	case block:
		f := readResponse(key)
		var out rpc.Block
		json.Unmarshal(f, &out)
		return &out
	case blocks:
		f := readResponse(key)
		var out [][]string
		json.Unmarshal(f, &out)
		return out
	case chainid:
		f := readResponse(key)
		var out string
		json.Unmarshal(f, &out)
		return out
	case contract:
		f := readResponse(key)
		var out rpc.Contract
		json.Unmarshal(f, &out)
		return out
	case connections:
		f := readResponse(key)
		var out rpc.Connections
		json.Unmarshal(f, &out)
		return out
	case constants:
		f := readResponse(key)
		var out rpc.Constants
		json.Unmarshal(f, &out)
		return out
	case counter:
		f := readResponse(key)
		var out int
		json.Unmarshal(f, &out)
		return out
	case currentLevel:
		f := readResponse(key)
		var out rpc.CurrentLevel
		json.Unmarshal(f, &out)
		return out
	case cycle:
		f := readResponse(key)
		var out rpc.Cycle
		json.Unmarshal(f, &out)
		return out
	case delegate:
		f := readResponse(key)
		var out rpc.Delegate
		json.Unmarshal(f, &out)
		return out
	case delegatedcontracts:
		f := readResponse(key)
		var out []string
		json.Unmarshal(f, &out)
		return out
	case endorsingrights:
		f := readResponse(key)
		var out []rpc.EndorsingRights
		json.Unmarshal(f, &out)
		return out
	case frozenbalanceByCycle:
		f := readResponse(key)
		var out []rpc.FrozenBalanceByCycle
		json.Unmarshal(f, &out)
		return out
	case header:
		f := readResponse(key)
		var out rpc.Header
		json.Unmarshal(f, &out)
		return out
	case headerShell:
		f := readResponse(key)
		var out rpc.HeaderShell
		json.Unmarshal(f, &out)
		return out
	case liveBlocks:
		f := readResponse(key)
		var out []string
		json.Unmarshal(f, &out)
		return out
	case metadata:
		f := readResponse(key)
		var out rpc.Metadata
		json.Unmarshal(f, &out)
		return out
	case operations:
		f := readResponse(key)
		var out rpc.FlattenedOperations
		json.Unmarshal(f, &out)
		return out
	case operationhashes:
		f := readResponse(key)
		var out rpc.OperationHashes
		json.Unmarshal(f, &out)
		return out
	case operationMetaDataHashes:
		f := readResponse(key)
		var out rpc.OperationMetadataHashes
		json.Unmarshal(f, &out)
		return out
	case parseOperations:
		f := readResponse(key)
		var out []rpc.Operations
		json.Unmarshal(f, &out)
		return out
	case preapplyOperations:
		f := readResponse(key)
		var out []rpc.Operations
		json.Unmarshal(f, &out)
		return out
	case proposals:
		f := readResponse(key)
		var out rpc.Proposals
		json.Unmarshal(f, &out)
		return out
	case protocols:
		f := readResponse(key)
		var out rpc.Protocols
		json.Unmarshal(f, &out)
		return out
	case protocolData:
		f := readResponse(key)
		var out rpc.ProtocolData
		json.Unmarshal(f, &out)
		return out
	case rpcerrors:
		f := readResponse(key)
		var out rpc.Errors
		json.Unmarshal(f, &out)
		return out
	case voteListings:
		f := readResponse(key)
		var out rpc.Listings
		json.Unmarshal(f, &out)
		return out
	default:
		return nil
	}
}

// The below variables contain mocks that are unmarshaled.
var (
	mockAddressTz1 = "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc"
	mockBlockHash  = "BLzGD63HA4RP8Fh5xEtvdQSMKa2WzJMZjQPNVUc4Rqy8Lh5BEY1"
)

// Regexes to allow the capture of custom handlers for unit testing.
var (
	regActiveChains                 = regexp.MustCompile(`\/monitor\/active_chains`)
	regBakingRights                 = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/helpers\/baking_rights`)
	regContractBalance              = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/contracts\/[A-z0-9]+\/balance`)
	regBallotList                   = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/votes\/ballot_list`)
	regBallots                      = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/votes\/ballots`)
	regBlock                        = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+`)
	regBigMap                       = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/big_maps\/[0-9]+\/[A-z0-9]+`)
	regContracts                    = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/contracts`)
	regContract                     = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/contracts\/[A-z0-9]+`)
	regConnections                  = regexp.MustCompile(`\/network\/connections`)
	regConstants                    = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/constants`)
	regContractDelegate             = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/contracts\/[A-z0-9]+\/delegate`)
	regContractEntrypoints          = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/contracts\/[A-z0-9]+\/entrypoints`)
	regContractEntrypoint           = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/contracts\/[A-z0-9]+\/entrypoints\/[A-z]+`)
	regContractManagerKey           = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/contracts\/[A-z0-9]+\/manager_key`)
	regContractScript               = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/contracts\/[A-z0-9]+\/script`)
	regContractCounter              = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/contracts\/[A-z0-9]+\/counter`)
	regCurrentPeriodKind            = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/votes\/current_period_kind`)
	regCurrentPeriod                = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/votes\/current_period_kind`)
	regSuccessorPeriod              = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/votes\/successor_period`)
	regCurrentLevel                 = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/helpers\/current_level`)
	regCurrentProposal              = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/votes\/current_proposal`)
	regCurrentQuorum                = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/votes\/current_quorum`)
	regTotalVotingPower             = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/votes\/total_voting_power`)
	regCycle                        = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/raw\/json\/cycle\/[0-9]+`)
	regDelegate                     = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/delegates\/[A-z0-9]+`)
	regDelegates                    = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/delegates`)
	regDelegateDelegatedContracts   = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/delegates\/[A-z0-9]+\/delegated_contracts`)
	regDelegateBalance              = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/delegates\/[A-z0-9]+\/balance`)
	regDelegateDeactivated          = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/delegates\/[A-z0-9]+\/deactivated`)
	regDelegateDelegatedBalance     = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/delegates\/[A-z0-9]+\/delegated_balance`)
	regDelegateFrozenBalance        = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/delegates\/[A-z0-9]+\/frozen_balance`)
	regDelegateFrozenBalanceByCycle = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/delegates\/[A-z0-9]+\/frozen_balance_by_cycle`)
	regDelegateGracePeriod          = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/delegates\/[A-z0-9]+\/grace_period`)
	regDelegateVotingPower          = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/delegates\/[A-z0-9]+\/voting_power`)
	regEndorsingRights              = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/helpers\/endorsing_rights`)
	regEndorsingPower               = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/endorsing_power`)
	regEntrypoint                   = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/helpers/scripts/entrypoint`)
	regEntrypoints                  = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/helpers/scripts/entrypoints`)
	regForgeOperationWithRPC        = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/helpers\/forge\/operations`)
	regForgeBlockHeader             = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/helpers\/forge_block_header`)
	regHash                         = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/hash`)
	regHeader                       = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/header`)
	regHeaderRaw                    = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/header\/raw`)
	regHeaderShell                  = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/header\/shell`)
	regHeaderProtocolData           = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/header\/protocol_data`)
	regHeaderProtocolDataRaw        = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/header\/protocol_data/raw`)
	regLiveBlocks                   = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/live_blocks`)
	regInjectionBlock               = regexp.MustCompile(`\/injection\/block`)
	regInjectionOperation           = regexp.MustCompile(`\/injection\/operation`)
	regLevelsInCurrentCycle         = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/helpers\/levels_in_current_cycle`)
	regMetadata                     = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/metadata`)
	regMetadataHash                 = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/metadata_hash`)
	regMinimalValidTime             = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/minimal_valid_time`)
	regNonces                       = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/nonces\/[0-9]+`)
	regOperations                   = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/operations`)
	regOperationsMetadataHash       = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/operations_metadata_hash`)
	regOperationHashes              = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/operation_hashes`)
	regOperationMetadataHashes      = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/operation_metadata_hashes`)
	regPackData                     = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/helpers/scripts/pack_data`)
	regParseBlock                   = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/helpers\/parse\/block`)
	regParseOperations              = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/helpers\/parse\/operations`)
	regPreapplyBlock                = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/helpers\/preapply\/block`)
	regPreapplyOperations           = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/helpers\/preapply\/operations`)
	regProposals                    = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/votes\/proposals`)
	regProtocols                    = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/protocols`)
	regRawBytes                     = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/raw\/bytes`)
	regRunCode                      = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/helpers/scripts/run_code`)
	regRunOperation                 = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/helpers\/scripts\/run_operation`)
	regRequiredEndorsements         = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/required_endorsements`)
	regSeed                         = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/seed`)
	regTraceCode                    = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/helpers\/scripts\/trace_code`)
	regTypecheckCode                = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/helpers\/scripts\/typecheck_code`)
	regTypecheckData                = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/helpers\/scripts\/typecheck_data`)
	regDelegateStakingBalance       = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/delegates\/[A-z0-9]+\/staking_balance`)
	regContractStorage              = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/contracts\/[A-z0-9]+\/storage`)
	regVoteListings                 = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/votes\/listings`)
)

// blankHandler handles the end of a http test handler chain
var blankHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

// ----------------------------------------- //
// Mock Handlers
// The below handlers are to simulate the Tezos RPC server for unit testing.
// ----------------------------------------- //

type requestResultPair struct {
	requestPath *regexp.Regexp
	resp        []byte
}

func mockHandler(pair *requestResultPair, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if pair.requestPath.MatchString(r.URL.String()) {
			w.Write(pair.resp)
			return
		}

		next.ServeHTTP(w, r)
	})
}

type blockHandlerMock struct {
	used bool
}

func newBlockMock() *blockHandlerMock {
	return &blockHandlerMock{}
}

func (b *blockHandlerMock) handler(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regBlock.MatchString(r.URL.String()) && !b.used {
			w.Write(resp)
			b.used = true
			return
		}

		next.ServeHTTP(w, r)
	})
}

type constantsHandlerMock struct {
	used bool
}

func newConstantsMock() *constantsHandlerMock {
	return &constantsHandlerMock{}
}

func (c *constantsHandlerMock) handler(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regConstants.MatchString(r.URL.String()) && !c.used {
			w.Write(resp)
			c.used = true
			return
		}

		next.ServeHTTP(w, r)
	})
}

func cycleHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regCycle.MatchString(r.URL.String()) {
			w.Write(resp)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func runOperationHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regRunOperation.MatchString(r.URL.String()) {
			w.Write(resp)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func checkErr(t *testing.T, wantErr bool, errContains string, err error) {
	if wantErr {
		assert.Error(t, err)
		if err != nil {
			assert.Contains(t, err.Error(), errContains)
		}
	} else {
		assert.Nil(t, err)
	}
}
