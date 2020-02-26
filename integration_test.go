// +build integration

package gotezos

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Head_Integration(t *testing.T) {
	gt, err := New(getMainnetIntegration())
	assert.Nil(t, err)
	head, err := gt.Head()
	assert.Nil(t, err)
	assert.NotNil(t, head)
}

func Test_Block_Integration(t *testing.T) {
	gt, err := New(getMainnetIntegration())
	assert.Nil(t, err)
	head, err := gt.Block(100000)
	assert.Nil(t, err)
	assert.NotNil(t, head)
}

func Test_Operation_Hashes_Integration(t *testing.T) {
	gt, err := New(getMainnetIntegration())
	assert.Nil(t, err)

	head, err := gt.Block(100000)
	assert.Nil(t, err)

	hashes, err := gt.OperationHashes(head.Hash)
	assert.Nil(t, err)
	assert.NotNil(t, hashes)
}

func Test_Balance_Integration(t *testing.T) {
	gt, err := New(getMainnetIntegration())
	assert.Nil(t, err)

	head, err := gt.Block(100000)
	assert.Nil(t, err)

	hashes, err := gt.Balance(head.Hash, mockAddressTz1)
	assert.Nil(t, err)
	assert.NotNil(t, hashes)
}

func Test_Blocks_Integration(t *testing.T) {
	gt, err := New(getMainnetIntegration())
	assert.Nil(t, err)

	hashes, err := gt.Blocks(&BlocksInput{})
	assert.Nil(t, err)
	assert.NotNil(t, hashes)
}

func Test_ChainID_Integration(t *testing.T) {
	gt, err := New(getMainnetIntegration())
	assert.Nil(t, err)

	chainid, err := gt.ChainID()
	assert.Nil(t, err)
	assert.NotNil(t, chainid)
}

func Test_Checkpoint_Integration(t *testing.T) {
	gt, err := New(getMainnetIntegration())
	assert.Nil(t, err)

	checkpoint, err := gt.Checkpoint()
	assert.Nil(t, err)
	assert.NotNil(t, checkpoint)
}

func Test_Invalid_Blocks_Integration(t *testing.T) {
	gt, err := New(getMainnetIntegration())
	assert.Nil(t, err)

	checkpoint, err := gt.InvalidBlocks()
	assert.Nil(t, err)
	assert.NotNil(t, checkpoint)
}

func getMainnetIntegration() string {
	return os.Getenv("TEZOS_MAINNET_URL")
}

// func getTestnetIntegration() string {
// 	return os.Getenv("TEZOS_TESTNET_URL")
// }
