// +build integration

package rpc

import (
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	skipNonExposed = false
	mainnetURL     string
)

func init() {
	mainnetURL = os.Getenv("GOTEZOS_MAINNET")
	if mainnetURL == "" {
		mainnetURL = "https://mainnet-tezos.giganode.io"
		skipNonExposed = true
	}
}

func Test_Balance_Integration(t *testing.T) {
	rpc, err := New(mainnetURL)
	assert.Nil(t, err)

	head, err := rpc.Head()
	assert.Nil(t, err)

	_, err = rpc.Balance(BalanceInput{
		Blockhash: head.Hash,
		Address:   "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
	})
	assert.Nil(t, err)
}

func Test_Head_Integration(t *testing.T) {
	rpc, err := New(mainnetURL)
	assert.Nil(t, err)

	_, err = rpc.Head()
	assert.Nil(t, err)
}

func Test_Block_Integration(t *testing.T) {
	rpc, err := New(mainnetURL)
	assert.Nil(t, err)

	min := 7
	max := 1000000

	var randomBlocks []int
	for i := 0; i < 50; i++ {
		randomBlocks = append(randomBlocks, rand.Intn(max-min)+min)
	}

	for _, block := range randomBlocks {
		_, err := rpc.Block(block)
		assert.Nil(t, err, fmt.Sprintf("Failed to get block: %d", block))
	}
}

func Test_OperationHashes_Integration(t *testing.T) {
	rpc, err := New(mainnetURL)
	assert.Nil(t, err)

	head, err := rpc.Head()
	assert.Nil(t, err)

	_, err = rpc.OperationHashes(head.Hash)
	assert.Nil(t, err)
}

func Test_BallotList_Integration(t *testing.T) {
	rpc, err := New(mainnetURL)
	assert.Nil(t, err)

	head, err := rpc.Block(647830)
	assert.Nil(t, err)

	_, err = rpc.BallotList(head.Hash)
	assert.Nil(t, err)
}

func Test_Ballots_Integration(t *testing.T) {
	rpc, err := New(mainnetURL)
	assert.Nil(t, err)

	head, err := rpc.Block(647830)
	assert.Nil(t, err)

	_, err = rpc.Ballots(head.Hash)
	assert.Nil(t, err)
}

func Test_CurrentPeriodKind_Integration(t *testing.T) {
	rpc, err := New(mainnetURL)
	assert.Nil(t, err)

	head, err := rpc.Block(647830)
	assert.Nil(t, err)

	_, err = rpc.CurrentPeriodKind(head.Hash)
	assert.Nil(t, err)
}

func Test_CurrentProposal_Integration(t *testing.T) {
	rpc, err := New(mainnetURL)
	assert.Nil(t, err)

	head, err := rpc.Block(647830)
	assert.Nil(t, err)

	_, err = rpc.CurrentProposal(head.Hash)
	assert.Nil(t, err)
}

func Test_VoteListings_Integration(t *testing.T) {
	rpc, err := New(mainnetURL)
	assert.Nil(t, err)

	head, err := rpc.Block(647830)
	assert.Nil(t, err)

	_, err = rpc.VoteListings(head.Hash)
	assert.Nil(t, err)
}

func Test_Proposals_Integration(t *testing.T) {
	rpc, err := New(mainnetURL)
	assert.Nil(t, err)

	head, err := rpc.Block(550000)
	assert.Nil(t, err)

	_, err = rpc.Proposals(head.Hash)
	assert.Nil(t, err)
}

func Test_Blocks_Integration(t *testing.T) {
	rpc, err := New(mainnetURL)
	assert.Nil(t, err)

	blocks, err := rpc.Blocks(BlocksInput{
		Length: 100,
	})
	assert.Nil(t, err)

	length := 0
	for _, b := range blocks {
		length = len(b)
	}
	assert.Equal(t, 100, length)
}

func Test_ChainID_Integration(t *testing.T) {
	rpc, err := New(mainnetURL)
	assert.Nil(t, err)

	_, err = rpc.ChainID()
	assert.Nil(t, err)
}

func Test_Checkpoint_Integration(t *testing.T) {
	rpc, err := New(mainnetURL)
	assert.Nil(t, err)

	_, err = rpc.Checkpoint()
	assert.Nil(t, err)
}

func Test_InvalidBlocks_Integration(t *testing.T) {
	rpc, err := New(mainnetURL)
	assert.Nil(t, err)

	_, err = rpc.InvalidBlocks()
	assert.Nil(t, err)
}

func Test_DelegatedContracts_Integration(t *testing.T) {
	rpc, err := New(mainnetURL)
	assert.Nil(t, err)

	head, err := rpc.Head()
	assert.Nil(t, err)

	_, err = rpc.DelegatedContracts(DelegatedContractsInput{
		Blockhash: head.Hash,
		Delegate:  "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
	})
	assert.Nil(t, err)

	_, err = rpc.DelegatedContracts(DelegatedContractsInput{
		Cycle:    100,
		Delegate: "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
	})
	assert.Nil(t, err)
}

func Test_FrozenBalance_Integration(t *testing.T) {
	rpc, err := New(mainnetURL)
	assert.Nil(t, err)

	_, err = rpc.FrozenBalance(FrozenBalanceInput{
		Cycle:    100,
		Delegate: "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
	})
	assert.Nil(t, err)
}

func Test_Delegate_Integration(t *testing.T) {
	rpc, err := New(mainnetURL)
	assert.Nil(t, err)

	head, err := rpc.Head()
	assert.Nil(t, err)

	_, err = rpc.Delegate(DelegateInput{
		Blockhash: head.Hash,
		Delegate:  "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
	})
	assert.Nil(t, err)

	_, err = rpc.Delegate(DelegateInput{
		Cycle:    100,
		Delegate: "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
	})
	assert.Nil(t, err)
}

func Test_StakingBalance_Integration(t *testing.T) {
	rpc, err := New(mainnetURL)
	assert.Nil(t, err)

	head, err := rpc.Head()
	assert.Nil(t, err)

	_, err = rpc.StakingBalance(StakingBalanceInput{
		Blockhash: head.Hash,
		Delegate:  "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
	})
	assert.Nil(t, err)

	_, err = rpc.StakingBalance(StakingBalanceInput{
		Cycle:    100,
		Delegate: "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
	})
	assert.Nil(t, err)
}

func Test_BakingRights_Integration(t *testing.T) {
	rpc, err := New(mainnetURL)
	assert.Nil(t, err)

	head, err := rpc.Head()
	assert.Nil(t, err)

	_, err = rpc.BakingRights(BakingRightsInput{
		BlockHash: head.Hash,
	})
	assert.Nil(t, err)
}

func Test_EndorsingRights_Integration(t *testing.T) {
	rpc, err := New(mainnetURL)
	assert.Nil(t, err)

	head, err := rpc.Head()
	assert.Nil(t, err)

	_, err = rpc.EndorsingRights(EndorsingRightsInput{
		BlockHash: head.Hash,
	})
	assert.Nil(t, err)
}

func Test_Delegates_Integration(t *testing.T) {
	rpc, err := New(mainnetURL)
	assert.Nil(t, err)

	head, err := rpc.Head()
	assert.Nil(t, err)

	_, err = rpc.Delegates(DelegatesInput{
		Blockhash: head.Hash,
	})
	assert.Nil(t, err)

	_, err = rpc.Delegates(DelegatesInput{
		Cycle: 100,
	})
	assert.Nil(t, err)
}

func Test_Version_Integration(t *testing.T) {
	rpc, err := New(mainnetURL)
	assert.Nil(t, err)

	_, err = rpc.Version()
	assert.Nil(t, err)
}

func Test_Constants_Integration(t *testing.T) {
	rpc, err := New(mainnetURL)
	assert.Nil(t, err)

	head, err := rpc.Head()
	assert.Nil(t, err)

	_, err = rpc.Constants(head.Hash)
	assert.Nil(t, err)
}

func Test_Connections_Integration(t *testing.T) {
	if !skipNonExposed {
		rpc, err := New(mainnetURL)
		assert.Nil(t, err)

		_, err = rpc.Connections()
		assert.Nil(t, err)
	}
}

func Test_Bootstrap_Integration(t *testing.T) {
	rpc, err := New(mainnetURL)
	assert.Nil(t, err)

	_, err = rpc.Bootstrap()
	assert.Nil(t, err)
}

func Test_Commit_Integration(t *testing.T) {
	rpc, err := New(mainnetURL)
	assert.Nil(t, err)

	_, err = rpc.Commit()
	assert.Nil(t, err)
}

func Test_Cycle_Integration(t *testing.T) {
	rpc, err := New(mainnetURL)
	assert.Nil(t, err)

	_, err = rpc.Cycle(100)
	assert.Nil(t, err)
}
