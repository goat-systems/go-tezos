package gotezos

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

type responseKey string

func (r *responseKey) String() string {
	return string(*r)
}

const (
	activechains       responseKey = ".test-fixtures/active_chains.json"
	bakingrights       responseKey = ".test-fixtures/baking_rights.json"
	balance            responseKey = ".test-fixtures/balance.json"
	block              responseKey = ".test-fixtures/block.json"
	blocks             responseKey = ".test-fixtures/blocks.json"
	bootstrap          responseKey = ".test-fixtures/bootstrap.json"
	chainid            responseKey = ".test-fixtures/chain_id.json"
	checkpoint         responseKey = ".test-fixtures/checkpoint.json"
	commit             responseKey = ".test-fixtures/commit.json"
	connections        responseKey = ".test-fixtures/connections.json"
	constants          responseKey = ".test-fixtures/constants.json"
	counter            responseKey = ".test-fixtures/counter.json"
	cycle              responseKey = ".test-fixtures/cycle.json"
	delegate           responseKey = ".test-fixtures/delegate.json"
	delegatedcontracts responseKey = ".test-fixtures/delegated_contracts.json"
	endorsingrights    responseKey = ".test-fixtures/endorsing_rights.json"
	frozenbalance      responseKey = ".test-fixtures/frozen_balance.json"
	invalidblock       responseKey = ".test-fixtures/invalid_block.json"
	invalidblocks      responseKey = ".test-fixtures/invalid_blocks.json"
	operationhashes    responseKey = ".test-fixtures/operation_hashes.json"
	rpcerrors          responseKey = ".test-fixtures/rpc_errors.json"
	version            responseKey = ".test-fixtures/version.json"
)

func readResponse(key responseKey) []byte {
	f, _ := ioutil.ReadFile(key.String())
	return f
}

func getResponse(key responseKey) interface{} {
	switch key {
	case activechains:
		f := readResponse(key)
		var out ActiveChains
		json.Unmarshal(f, &out)
		return out
	case bakingrights:
		f := readResponse(key)
		var out BakingRights
		json.Unmarshal(f, &out)
		return &out
	case balance:
		f := readResponse(key)
		var out Int
		json.Unmarshal(f, &out)
		return &out
	case block:
		f := readResponse(key)
		var out Block
		json.Unmarshal(f, &out)
		return &out
	case blocks:
		f := readResponse(key)
		var out [][]string
		json.Unmarshal(f, &out)
		return out
	case bootstrap:
		f := readResponse(key)
		var out Bootstrap
		json.Unmarshal(f, &out)
		return out
	case chainid:
		f := readResponse(key)
		var out string
		json.Unmarshal(f, &out)
		return out
	case checkpoint:
		f := readResponse(key)
		var out Checkpoint
		json.Unmarshal(f, &out)
		return out
	case commit:
		f := readResponse(key)
		var out string
		json.Unmarshal(f, &out)
		return out
	case connections:
		f := readResponse(key)
		var out Connections
		json.Unmarshal(f, &out)
		return out
	case constants:
		f := readResponse(key)
		var out Constants
		json.Unmarshal(f, &out)
		return out
	case counter:
		f := readResponse(key)
		var out int
		json.Unmarshal(f, &out)
		return out
	case cycle:
		f := readResponse(key)
		var out Cycle
		json.Unmarshal(f, &out)
		return out
	case delegate:
		f := readResponse(key)
		var out Delegate
		json.Unmarshal(f, &out)
		return out
	case delegatedcontracts:
		f := readResponse(key)
		var out []*string
		json.Unmarshal(f, &out)
		return out
	case endorsingrights:
		f := readResponse(key)
		var out EndorsingRights
		json.Unmarshal(f, &out)
		return &out
	case frozenbalance:
		f := readResponse(key)
		var out FrozenBalance
		json.Unmarshal(f, &out)
		return out
	case invalidblock:
		f := readResponse(key)
		var out InvalidBlock
		json.Unmarshal(f, &out)
		return out
	case invalidblocks:
		f := readResponse(key)
		var out []InvalidBlock
		json.Unmarshal(f, &out)
		return out
	case operationhashes:
		f := readResponse(key)
		var out [][]string
		json.Unmarshal(f, &out)
		return out
	case rpcerrors:
		f := readResponse(key)
		var out RPCErrors
		json.Unmarshal(f, &out)
		return out
	case version:
		f := readResponse(key)
		var out Version
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
	regActiveChains            = regexp.MustCompile(`\/monitor\/active_chains`)
	regBakingRights            = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/helpers\/baking_rights`)
	regBalance                 = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/contracts\/[A-z0-9]+\/balance`)
	regBlock                   = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+`)
	regBlocks                  = regexp.MustCompile(`\/chains\/main\/blocks`)
	regBoostrap                = regexp.MustCompile(`\/monitor\/bootstrapped`)
	regChainID                 = regexp.MustCompile(`\/chains\/main\/chain_id`)
	regCheckpoint              = regexp.MustCompile(`\/chains\/main\/checkpoint`)
	regCommit                  = regexp.MustCompile(`\/monitor\/commit_hash`)
	regConnections             = regexp.MustCompile(`\/network\/connections`)
	regConstants               = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/constants`)
	regCounter                 = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/contracts\/[A-z0-9]+\/counter`)
	regCycle                   = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/raw\/json\/cycle\/[0-9]+`)
	regDelegate                = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/delegates\/[A-z0-9]+`)
	regDelegates               = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/delegates`)
	regDelegatedContracts      = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/delegates\/[A-z0-9]+\/delegated_contracts`)
	regEndorsingRights         = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/helpers\/endorsing_rights`)
	regFrozenBalance           = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/raw\/json\/contracts\/index\/[A-z0-9]+\/frozen_balance\/[0-9]+`)
	regForgeOperationWithRPC   = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/helpers\/forge\/operations`)
	regInjectionBlock          = regexp.MustCompile(`\/injection\/block`)
	regInjectionOperation      = regexp.MustCompile(`\/injection\/operation`)
	regInvalidBlocks           = regexp.MustCompile(`\/chains\/main\/invalid_blocks`)
	regOperationHashes         = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/operation_hashes`)
	regParseOperations         = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/helpers\/parse\/operations`)
	regPreapplyOperations      = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/helpers\/preapply\/operations`)
	regStakingBalance          = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/delegates\/[A-z0-9]+\/staking_balance`)
	regStorage                 = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/context\/contracts\/[A-z0-9]+\/storage`)
	regUnforgeOperationWithRPC = regexp.MustCompile(`\/chains\/main\/blocks\/[A-z0-9]+\/helpers\/parse\/operations`)
	regVersions                = regexp.MustCompile(`\/network\/version`)
)

// blankHandler handles the end of a http test handler chain
var blankHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

// ----------------------------------------- //
// Mock Handlers
// The below handlers are to simulate the Tezos RPC server for unit testing.
// ----------------------------------------- //

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

func balanceHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regBalance.MatchString(r.URL.String()) {
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

func checkpointHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regCheckpoint.MatchString(r.URL.String()) {
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

func counterHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regCounter.MatchString(r.URL.String()) {
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

func delegateHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regDelegate.MatchString(r.URL.String()) {
			w.Write(resp)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func delegatesHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regDelegates.MatchString(r.URL.String()) {
			w.Write(resp)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func delegationsHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regDelegatedContracts.MatchString(r.URL.String()) {
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

func frozenBalanceHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regFrozenBalance.MatchString(r.URL.String()) {
			w.Write(resp)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func forgeOperationWithRPCMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regForgeOperationWithRPC.MatchString(r.URL.String()) {
			fmt.Println(resp)
			w.Write(resp)
			return
		}

		fmt.Println("Here")
		next.ServeHTTP(w, r)
	})
}

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

func invalidBlocksHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regInvalidBlocks.MatchString(r.URL.String()) {
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

func parseOperationsHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regParseOperations.MatchString(r.URL.String()) {
			w.Write(resp)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func preapplyOperationsHandlerMock(preapplyResp, blockResp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regPreapplyOperations.MatchString(r.URL.String()) {
			w.Write(preapplyResp)
			return
		}

		if regBlock.MatchString(r.URL.String()) {
			w.Write(blockResp)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func stakingBalanceHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regStakingBalance.MatchString(r.URL.String()) {
			w.Write(resp)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func storageHandlerMock(resp []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regStorage.MatchString(r.URL.String()) {
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

func checkErr(t *testing.T, wantErr bool, errContains string, err error) {
	if wantErr {
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errContains)
	} else {
		assert.Nil(t, err)
	}
}
