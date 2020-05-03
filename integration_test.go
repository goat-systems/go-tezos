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

	hashes, err := gt.Blocks(BlocksInput{})
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

	invalidBlocks, err := gt.InvalidBlocks()
	assert.Nil(t, err)
	assert.NotNil(t, invalidBlocks)
}

// func Test_UnforgeOperationWithRPC(t *testing.T) {
// 	gt, err := New(getMainnetIntegration())
// 	assert.Nil(t, err)

// 	head, err := gt.Block(100000)
// 	assert.Nil(t, err)

// 	_, err = gt.UnforgeOperationWithRPC(head.Hash, UnforgeOperationWithRPCInput{
// 		Operations: []string{"a732d3520eeaa3de98d78e5e5cb6c85f72204fd46feb9f76853841d4a701add36c0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e00b960000008ba0cb2fad622697145cf1665124096d25bc31e006c0008ba0cb2fad622697145cf1665124096d25bc31ed3e7bd1008d3bb0300b1a803018b88e99e66c1c2587f87118449f781cb7d44c9c40000"},
// 	})
// 	assert.Nil(t, err)
// }

func getMainnetIntegration() string {
	return os.Getenv("TEZOS_MAINNET_URL")
}
