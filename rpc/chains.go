package rpc

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

/*
Checkpoint represents a Tezos checkpoint.

RPC:
	/chains/<chain_id>/checkpoint (GET)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-chains-chain-id-checkpoint
*/
type Checkpoint struct {
	Block struct {
		Level          int       `json:"level"`
		Proto          int       `json:"proto"`
		Predecessor    string    `json:"predecessor"`
		Timestamp      time.Time `json:"timestamp"`
		ValidationPass int       `json:"validation_pass"`
		OperationsHash string    `json:"operations_hash"`
		Fitness        []string  `json:"fitness"`
		Context        string    `json:"context"`
		ProtocolData   string    `json:"protocol_data"`
	} `json:"block"`
	SavePoint   int    `json:"save_point"`
	Caboose     int    `json:"caboose"`
	HistoryMode string `json:"history_mode"`
}

/*
InvalidBlock represents Tezos invalid blocks.

RPC:
	/chains/<chain_id>/invalid_blocks (GET)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-chains-chain-id-invalid-blocks
*/
type InvalidBlock struct {
	Block  string    `json:"block"`
	Level  int       `json:"level"`
	Errors RPCErrors `json:"errors"`
}

/*
BlocksInput is the input for the goTezos.Blocks function.

Function:
	func (t *GoTezos) EndorsingRights(input *EndorsingRightsInput) (*EndorsingRights, error) {}
*/
type BlocksInput struct {
	//length is the requested number of predecessors to returns (per requested head).
	Length int
	//An empty argument requests blocks from the current heads. A non empty list allow to request specific fragment of the chain.
	Head string
	// When `min_date` is provided, heads with a timestamp before `min_date` are filtered out
	MinDate *time.Time
}

/*
Blocks lists known heads of the blockchain sorted with decreasing fitness.
Optional arguments allows to returns the list of predecessors for known heads or
the list of predecessors for a given list of blocks.

Path:
	/chains/<chain_id>/blocks (GET)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-chains-chain-id-blocks

Parameters:

	input:
		Modifies the Blocks RPC query by passing optional URL parameters.
*/
func (c *Client) Blocks(input BlocksInput) ([][]string, error) {
	resp, err := c.get("/chains/main/blocks", input.contructRPCOptions()...)
	if err != nil {
		return [][]string{}, errors.Wrap(err, "failed to get blocks")
	}

	var blocks [][]string
	err = json.Unmarshal(resp, &blocks)
	if err != nil {
		return [][]string{}, errors.Wrap(err, "failed to unmarshal blocks")
	}

	return blocks, nil
}

func (b *BlocksInput) contructRPCOptions() []rpcOptions {
	var opts []rpcOptions
	if b.Length > 0 {
		opts = append(opts, rpcOptions{
			"length",
			strconv.Itoa(b.Length),
		})
	}

	if b.Head != "" {
		opts = append(opts, rpcOptions{
			"head",
			b.Head,
		})
	}

	if b.MinDate != nil {
		opts = append(opts, rpcOptions{
			"min_date",
			strconv.Itoa(int(b.MinDate.Unix())),
		})
	}

	return opts
}

/*
ChainID gets the chain unique identifier.

Path:
	/chains/<chain_id>/chain_id (GET)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-chains-chain-id-chain-id
*/
func (c *Client) ChainID() (string, error) {
	resp, err := c.get("/chains/main/chain_id")
	if err != nil {
		return "", errors.Wrapf(err, "failed to get chain id")
	}

	var chainID string
	err = json.Unmarshal(resp, &chainID)
	if err != nil {
		return "", errors.Wrapf(err, "failed to unmarshal chain id")
	}

	return chainID, nil
}

/*
Checkpoint gets the current checkpoint for this chain.

Path:
	/chains/<chain_id>/checkpoint (GET)RPC

Link:
	https://tezos.gitlab.io/api/rpc.html#get-chains-chain-id-checkpoint
*/
func (c *Client) Checkpoint() (Checkpoint, error) {
	resp, err := c.get("/chains/main/checkpoint")
	if err != nil {
		return Checkpoint{}, errors.Wrap(err, "failed to get checkpoint")
	}

	var checkpoint Checkpoint
	err = json.Unmarshal(resp, &c)
	if err != nil {
		return checkpoint, errors.Wrap(err, "failed to unmarshal checkpoint")
	}

	return checkpoint, nil
}

/*
InvalidBlocks lists blocks that have been declared invalid
along with the errors that led to them being declared invalid.

Path:
	/chains/<chain_id>/invalid_blocks (GET)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-chains-chain-id-invalid-blocks
*/
func (c *Client) InvalidBlocks() ([]InvalidBlock, error) {
	resp, err := c.get("/chains/main/invalid_blocks")
	if err != nil {
		return []InvalidBlock{}, errors.Wrap(err, "failed to get invalid blocks")
	}

	var blocks []InvalidBlock
	err = json.Unmarshal(resp, &blocks)
	if err != nil {
		return []InvalidBlock{}, errors.Wrap(err, "failed to unmarshal invalid blocks")
	}

	return blocks, nil
}

/*
InvalidBlock gets the errors that appears during the block (in)validation.

Path:
	/chains/<chain_id>/invalid_blocks/<block_hash> (GET)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-chains-chain-id-invalid-blocks-block-hash
*/
func (c *Client) InvalidBlock(blockHash string) (InvalidBlock, error) {
	resp, err := c.get(fmt.Sprintf("/chains/main/invalid_blocks/%s", blockHash))
	if err != nil {
		return InvalidBlock{}, errors.Wrap(err, "failed to get invalid blocks")
	}

	var block InvalidBlock
	err = json.Unmarshal(resp, &block)
	if err != nil {
		return InvalidBlock{}, errors.Wrap(err, "failed to unmarshal invalid blocks")
	}

	return block, nil
}

/*
DeleteInvalidBlock remove an invalid block for the tezos storage.

Path:
	/chains/<chain_id>/invalid_blocks/<block_hash> (DELETE)

Link:
	https://tezos.gitlab.io/api/rpc.html#delete-chains-chain-id-invalid-blocks-block-hash
*/
func (c *Client) DeleteInvalidBlock(blockHash string) error {
	_, err := c.delete(fmt.Sprintf("/chains/main/invalid_blocks/%s", blockHash))
	if err != nil {
		return errors.Wrap(err, "failed to delete invalid blocks")
	}

	return nil
}
