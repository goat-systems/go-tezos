package rpc_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"testing"

	"github.com/goat-systems/go-tezos/v4/rpc"
	"github.com/stretchr/testify/assert"
)

type responseKey string

func (r *responseKey) String() string {
	return string(*r)
}

const (
	activechains         responseKey = ".test-fixtures/active_chains.json"
	bakingrights         responseKey = ".test-fixtures/baking_rights.json"
	balance              responseKey = ".test-fixtures/balance.json"
	ballotList           responseKey = ".test-fixtures/ballot_list.json"
	ballots              responseKey = ".test-fixtures/ballots.json"
	block                responseKey = ".test-fixtures/block.json"
	blocks               responseKey = ".test-fixtures/blocks.json"
	bootstrap            responseKey = ".test-fixtures/bootstrap.json"
	chainid              responseKey = ".test-fixtures/chain_id.json"
	commit               responseKey = ".test-fixtures/commit.json"
	contract             responseKey = ".test-fixtures/contract.json"
	connections          responseKey = ".test-fixtures/connections.json"
	constants            responseKey = ".test-fixtures/constants.json"
	contractEntrypoints  responseKey = ".test-fixtures/entrypoints.json"
	counter              responseKey = ".test-fixtures/counter.json"
	currentLevel         responseKey = ".test-fixtures/current_level.json"
	cycle                responseKey = ".test-fixtures/cycle.json"
	delegate             responseKey = ".test-fixtures/delegate.json"
	delegatedcontracts   responseKey = ".test-fixtures/delegated_contracts.json"
	endorsingrights      responseKey = ".test-fixtures/endorsing_rights.json"
	frozenbalance        responseKey = ".test-fixtures/frozen_balance.json"
	frozenbalanceByCycle responseKey = ".test-fixtures/frozen_balance_by_cycle.json"
	header               responseKey = ".test-fixtures/header.json"
	headerShell          responseKey = ".test-fixtures/header_shell.json"
	operationhashes      responseKey = ".test-fixtures/operation_hashes.json"
	parseOperations      responseKey = ".test-fixtures/parse_operations.json"
	preapplyOperations   responseKey = ".test-fixtures/preapply_operations.json"
	proposals            responseKey = ".test-fixtures/proposals.json"
	protocolData         responseKey = ".test-fixtures/protocol_data.json"
	rpcerrors            responseKey = ".test-fixtures/rpc_errors.json"
	version              responseKey = ".test-fixtures/version.json"
	voteListings         responseKey = ".test-fixtures/vote_listings.json"
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
		var out rpc.BakingRights
		json.Unmarshal(f, &out)
		return &out
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
	case bootstrap:
		f := readResponse(key)
		var out rpc.Bootstrap
		json.Unmarshal(f, &out)
		return out
	case chainid:
		f := readResponse(key)
		var out string
		json.Unmarshal(f, &out)
		return out
	case commit:
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
		var out rpc.EndorsingRights
		json.Unmarshal(f, &out)
		return &out
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
	case operationhashes:
		f := readResponse(key)
		var out [][]string
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
	case version:
		f := readResponse(key)
		var out rpc.Version
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
	regBlocks                       = regexp.MustCompile(`\/chains\/main\/blocks`)
	regBigMap                       = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/big_maps\/[0-9]+\/[A-z0-9]+`)
	regBoostrap                     = regexp.MustCompile(`\/monitor\/bootstrapped`)
	regChainID                      = regexp.MustCompile(`\/chains\/main\/chain_id`)
	regCommit                       = regexp.MustCompile(`\/monitor\/commit_hash`)
	regContracts                    = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/contracts`)
	regContract                     = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/contracts\/[A-z0-9]+`)
	regConnections                  = regexp.MustCompile(`\/network\/connections`)
	regCompletePrefix               = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/helpers/complete/[A-z0-9]+`)
	regConstants                    = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/constants`)
	regContractDelegate             = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/contracts\/[A-z0-9]+\/delegate`)
	regContractEntrypoints          = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/contracts\/[A-z0-9]+\/entrypoints`)
	regContractEntrypoint           = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/contracts\/[A-z0-9]+\/entrypoints\/[A-z]+`)
	regContractManagerKey           = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/contracts\/[A-z0-9]+\/manager_key`)
	regContractScript               = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/contracts\/[A-z0-9]+\/script`)
	regContractCounter              = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/contracts\/[A-z0-9]+\/counter`)
	regCurrentPeriodKind            = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/votes\/current_period_kind`)
	regCurrentLevel                 = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/helpers\/current_level`)
	regCurrentProposal              = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/votes\/current_proposal`)
	regCurrentQuorum                = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/votes\/current_quorum`)
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

	regEndorsingRights = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/helpers\/endorsing_rights`)
	regEndorsingPower  = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/endorsing_power`)
	regFrozenBalance   = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/raw\/json\/contracts\/index\/[A-z0-9]+\/frozen_balance\/[0-9]+`)
	//	regForgeOperationWithRPC   = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/helpers\/forge\/operations`)
	regHash                    = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/hash`)
	regHeader                  = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/header`)
	regHeaderRaw               = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/header\/raw`)
	regHeaderShell             = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/header\/shell`)
	regHeaderProtocolData      = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/header\/protocol_data`)
	regHeaderProtocolDataRaw   = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/header\/protocol_data/raw`)
	regInjectionBlock          = regexp.MustCompile(`\/injection\/block`)
	regInjectionOperation      = regexp.MustCompile(`\/injection\/operation`)
	regNonces                  = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/nonces\/[0-9]+`)
	regOperationHashes         = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/operation_hashes`)
	regPreapplyOperations      = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/helpers\/preapply\/operations`)
	regProposals               = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/votes\/proposals`)
	regRunOperation            = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/helpers\/scripts\/run_operation`)
	regRawBytes                = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/raw\/bytes`)
	regSeed                    = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/seed`)
	regDelegateStakingBalance  = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/delegates\/[A-z0-9]+\/staking_balance`)
	regContractStorage         = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/contracts\/[A-z0-9]+\/storage`)
	regUnforgeOperationWithRPC = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/helpers\/parse\/operations`)
	regVersions                = regexp.MustCompile(`\/network\/version`)
	regVoteListings            = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/votes\/listings`)
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

func activeChainsHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regActiveChains.MatchString(r.URL.String()) {
			w.Write(resp)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func bakingRightsHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regBakingRights.MatchString(r.URL.String()) {
			w.Write(resp)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func ballotListHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regBallotList.MatchString(r.URL.String()) {
			w.Write(resp)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func ballotsHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regBallots.MatchString(r.URL.String()) {
			w.Write(resp)
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

func blocksHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regBlocks.MatchString(r.URL.String()) {
			w.Write(resp)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func bootstrapHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regBoostrap.MatchString(r.URL.String()) {
			w.Write(resp)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func chainIDHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regChainID.MatchString(r.URL.String()) {
			w.Write(resp)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func commitHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regCommit.MatchString(r.URL.String()) {
			w.Write(resp)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func connectionsHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regConnections.MatchString(r.URL.String()) {
			w.Write(resp)
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

func currentPeriodKindHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regCurrentPeriodKind.MatchString(r.URL.String()) {
			w.Write(resp)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func currentProposalHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regCurrentProposal.MatchString(r.URL.String()) {
			w.Write(resp)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func currentQuorumHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regCurrentQuorum.MatchString(r.URL.String()) {
			w.Write(resp)
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

func delegateGracePeriodHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regDelegateGracePeriod.MatchString(r.URL.String()) {
			w.Write(resp)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func delegateVotingPowerHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regDelegateVotingPower.MatchString(r.URL.String()) {
			w.Write(resp)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func endorsingRightsHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regEndorsingRights.MatchString(r.URL.String()) {
			w.Write(resp)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// func forgeOperationWithRPCMock(resp []byte, next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if regForgeOperationWithRPC.MatchString(r.URL.String()) {
// 			w.Write(resp)
// 			return
// 		}

// 		next.ServeHTTP(w, r)
// 	})
// }

func injectionBlockHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regInjectionBlock.MatchString(r.URL.String()) {
			w.Write(resp)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func injectionOperationHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regInjectionOperation.MatchString(r.URL.String()) {
			w.Write(resp)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func operationHashesHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regOperationHashes.MatchString(r.URL.String()) {
			w.Write(resp)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func preapplyOperationsHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regPreapplyOperations.MatchString(r.URL.String()) {
			w.Write(resp)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func proposalsHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regProposals.MatchString(r.URL.String()) {
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

func unforgeOperationWithRPCMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regUnforgeOperationWithRPC.MatchString(r.URL.String()) {
			w.Write(resp)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func versionsHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regVersions.MatchString(r.URL.String()) {
			w.Write(resp)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func voteListingsHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regVoteListings.MatchString(r.URL.String()) {
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
