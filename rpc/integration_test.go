// +build integration

package rpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// var (
// 	skipNonExposed = false
// 	mainnetURL     string
// )

// func init() {
// 	mainnetURL = os.Getenv("GOTEZOS_MAINNET")
// 	if mainnetURL == "" {
// 		mainnetURL = "https://mainnet-tezos.giganode.io"
// 		skipNonExposed = true
// 	}
// }

// func Test_Balance_Integration(t *testing.T) {
// 	gt, err := New(mainnetURL)
// 	assert.Nil(t, err)

// 	head, err := gt.Head()
// 	assert.Nil(t, err)

// 	_, err = gt.Balance(head.Hash, "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc")
// 	assert.Nil(t, err)
// }

// func Test_Head_Integration(t *testing.T) {
// 	gt, err := New(mainnetURL)
// 	assert.Nil(t, err)

// 	_, err = gt.Head()
// 	assert.Nil(t, err)
// }

// func Test_Block_Integration(t *testing.T) {
// 	gt, err := New(mainnetURL)
// 	assert.Nil(t, err)

// 	min := 7
// 	max := 1000000

// 	var randomBlocks []int
// 	for i := 0; i < 300; i++ {
// 		randomBlocks = append(randomBlocks, rand.Intn(max-min)+min)
// 	}

// 	for _, block := range randomBlocks {
// 		fmt.Printf("block: %d\n", block)
// 		_, err := gt.Block(block)
// 		assert.Nil(t, err, fmt.Sprintf("Failed to get block: %d", block))
// 	}
// }

// func Test_OperationHashes_Integration(t *testing.T) {
// 	gt, err := New(mainnetURL)
// 	assert.Nil(t, err)

// 	head, err := gt.Head()
// 	assert.Nil(t, err)

// 	_, err = gt.OperationHashes(head.Hash)
// 	assert.Nil(t, err)
// }

// func Test_BallotList_Integration(t *testing.T) {
// 	gt, err := New(mainnetURL)
// 	assert.Nil(t, err)

// 	head, err := gt.Block(647830)
// 	assert.Nil(t, err)

// 	_, err = gt.BallotList(head.Hash)
// 	assert.Nil(t, err)
// }

// func Test_Ballots_Integration(t *testing.T) {
// 	gt, err := New(mainnetURL)
// 	assert.Nil(t, err)

// 	head, err := gt.Block(647830)
// 	assert.Nil(t, err)

// 	_, err = gt.Ballots(head.Hash)
// 	assert.Nil(t, err)
// }

// func Test_CurrentPeriodKind_Integration(t *testing.T) {
// 	gt, err := New(mainnetURL)
// 	assert.Nil(t, err)

// 	head, err := gt.Block(647830)
// 	assert.Nil(t, err)

// 	_, err = gt.CurrentPeriodKind(head.Hash)
// 	assert.Nil(t, err)
// }

// func Test_CurrentProposal_Integration(t *testing.T) {
// 	gt, err := New(mainnetURL)
// 	assert.Nil(t, err)

// 	head, err := gt.Block(647830)
// 	assert.Nil(t, err)

// 	_, err = gt.CurrentProposal(head.Hash)
// 	assert.Nil(t, err)
// }

// func Test_VoteListings_Integration(t *testing.T) {
// 	gt, err := New(mainnetURL)
// 	assert.Nil(t, err)

// 	head, err := gt.Block(647830)
// 	assert.Nil(t, err)

// 	_, err = gt.VoteListings(head.Hash)
// 	assert.Nil(t, err)
// }

// func Test_Proposals_Integration(t *testing.T) {
// 	gt, err := New(mainnetURL)
// 	assert.Nil(t, err)

// 	head, err := gt.Block(550000)
// 	assert.Nil(t, err)

// 	_, err = gt.Proposals(head.Hash)
// 	assert.Nil(t, err)
// }

// func Test_Blocks_Integration(t *testing.T) {
// 	gt, err := New(mainnetURL)
// 	assert.Nil(t, err)

// 	blocks, err := gt.Blocks(BlocksInput{
// 		Length: 100,
// 	})
// 	assert.Nil(t, err)

// 	length := 0
// 	for _, b := range blocks {
// 		length = len(b)
// 	}
// 	assert.Equal(t, 100, length)
// }

// func Test_ChainID_Integration(t *testing.T) {
// 	gt, err := New(mainnetURL)
// 	assert.Nil(t, err)

// 	_, err = gt.ChainID()
// 	assert.Nil(t, err)
// }

// func Test_Checkpoint_Integration(t *testing.T) {
// 	gt, err := New(mainnetURL)
// 	assert.Nil(t, err)

// 	_, err = gt.Checkpoint()
// 	assert.Nil(t, err)
// }

// func Test_InvalidBlocks_Integration(t *testing.T) {
// 	gt, err := New(mainnetURL)
// 	assert.Nil(t, err)

// 	_, err = gt.InvalidBlocks()
// 	assert.Nil(t, err)
// }

// func Test_DelegatedContracts_Integration(t *testing.T) {
// 	gt, err := New(mainnetURL)
// 	assert.Nil(t, err)

// 	head, err := gt.Head()
// 	assert.Nil(t, err)

// 	_, err = gt.DelegatedContracts(head.Hash, "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc")
// 	assert.Nil(t, err)
// }

// func Test_DelegatedContractsAtCycle_Integration(t *testing.T) {
// 	gt, err := New(mainnetURL)
// 	assert.Nil(t, err)

// 	_, err = gt.DelegatedContractsAtCycle(100, "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc")
// 	assert.Nil(t, err)
// }

// func Test_FrozenBalance_Integration(t *testing.T) {
// 	gt, err := New(mainnetURL)
// 	assert.Nil(t, err)

// 	_, err = gt.FrozenBalance(100, "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc")
// 	assert.Nil(t, err)
// }

// func Test_Delegate_Integration(t *testing.T) {
// 	gt, err := New(mainnetURL)
// 	assert.Nil(t, err)

// 	head, err := gt.Head()
// 	assert.Nil(t, err)

// 	_, err = gt.Delegate(head.Hash, "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc")
// 	assert.Nil(t, err)
// }

// func Test_StakingBalance_Integration(t *testing.T) {
// 	gt, err := New(mainnetURL)
// 	assert.Nil(t, err)

// 	head, err := gt.Head()
// 	assert.Nil(t, err)

// 	_, err = gt.StakingBalance(head.Hash, "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc")
// 	assert.Nil(t, err)
// }

// func Test_StakingBalanceAtCycle_Integration(t *testing.T) {
// 	gt, err := New(mainnetURL)
// 	assert.Nil(t, err)

// 	_, err = gt.StakingBalanceAtCycle(100, "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc")
// 	assert.Nil(t, err)
// }

// func Test_BakingRights_Integration(t *testing.T) {
// 	gt, err := New(mainnetURL)
// 	assert.Nil(t, err)

// 	head, err := gt.Head()
// 	assert.Nil(t, err)

// 	_, err = gt.BakingRights(BakingRightsInput{
// 		BlockHash: head.Hash,
// 	})
// 	assert.Nil(t, err)
// }

// func Test_EndorsingRights_Integration(t *testing.T) {
// 	gt, err := New(mainnetURL)
// 	assert.Nil(t, err)

// 	head, err := gt.Head()
// 	assert.Nil(t, err)

// 	_, err = gt.EndorsingRights(EndorsingRightsInput{
// 		BlockHash: head.Hash,
// 	})
// 	assert.Nil(t, err)
// }

// func Test_Delegates_Integration(t *testing.T) {
// 	gt, err := New(mainnetURL)
// 	assert.Nil(t, err)

// 	head, err := gt.Head()
// 	assert.Nil(t, err)

// 	_, err = gt.Delegates(DelegatesInput{
// 		BlockHash: head.Hash,
// 	})
// 	assert.Nil(t, err)
// }

// func Test_Version_Integration(t *testing.T) {
// 	gt, err := New(mainnetURL)
// 	assert.Nil(t, err)

// 	_, err = gt.Version()
// 	assert.Nil(t, err)
// }

// func Test_Constants_Integration(t *testing.T) {
// 	gt, err := New(mainnetURL)
// 	assert.Nil(t, err)

// 	head, err := gt.Head()
// 	assert.Nil(t, err)

// 	_, err = gt.Constants(head.Hash)
// 	assert.Nil(t, err)
// }

// func Test_Connections_Integration(t *testing.T) {
// 	if !skipNonExposed {
// 		gt, err := New(mainnetURL)
// 		assert.Nil(t, err)

// 		_, err = gt.Connections()
// 		assert.Nil(t, err)
// 	}
// }

// func Test_Bootstrap_Integration(t *testing.T) {
// 	gt, err := New(mainnetURL)
// 	assert.Nil(t, err)

// 	_, err = gt.Bootstrap()
// 	assert.Nil(t, err)
// }

// func Test_Commit_Integration(t *testing.T) {
// 	gt, err := New(mainnetURL)
// 	assert.Nil(t, err)

// 	_, err = gt.Commit()
// 	assert.Nil(t, err)
// }

// func Test_Cycle_Integration(t *testing.T) {
// 	gt, err := New(mainnetURL)
// 	assert.Nil(t, err)

// 	_, err = gt.Cycle(100)
// 	assert.Nil(t, err)
// }

func Test_Operations(t *testing.T) {
	gt, err := New("https://testnet-tezos.giganode.io")
	assert.Nil(t, err)

	head, err := gt.Head()
	assert.Nil(t, err)

	accountActivation := AccountActivation{
		Kind:   ACTIVATE_ACCOUNT,
		Pkh:    "tz1Yx9DZpJh2hBttzqNxApr7reEz5pi8mjfb",
		Secret: "4b28cca3859e9cd9803f5dca84372154e819df12",
	}

	operation, err := ForgeOperation(head.Hash, &accountActivation)
	assert.Nil(t, err)

	_, err = gt.InjectionOperation(InjectionOperationInput{
		Operation: operation,
	})
	assert.Nil(t, err)
}

// func Test_ActiveChains_Integration(t *testing.T) {
// 	gt, err := New(mainnetURL)
// 	assert.Nil(t, err)

// 	_, err = gt.ActiveChains()
// 	assert.Nil(t, err)
// }
