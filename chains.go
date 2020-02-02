package gotezos

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
)

/*
Checkpoint Result
RPC: /chains/<chain_id>/checkpoint (GET)
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
Blocks RPC
Path: /chains/<chain_id>/blocks (GET)
Description:  Lists known heads of the blockchain sorted with decreasing fitness.
Optional arguments allows to returns the list of predecessors for known heads or
the list of predecessors for a given list of blocks.

Options:
	length = <int> : The requested number of predecessors to returns (per requested head).
    head = <block_hash> : An empty argument requests blocks from the current heads. A non empty list allow to request specific fragment of the chain.
    min_date = <date> : When `min_date` is provided, heads with a timestamp before `min_date` are filtered out
*/
func (t *GoTezos) Blocks(opts ...RPCOptions) ([][]string, error) {
	resp, err := t.get("/chains/main/blocks", opts...)
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

/*
ChainID RPC
Path: /chains/<chain_id>/chain_id (GET)
Description: The chain unique identifier.
*/
func (t *GoTezos) ChainID() (string, error) {
	resp, err := t.get("/chains/main/chain_id")
	if err != nil {
		return "", errors.Wrapf(err, "failed to get chain id")
	}

	var chainID string
	err = json.Unmarshal(resp, &chainID)
	if err != nil {
		return chainID, errors.Wrapf(err, "failed to unmarshal chain id")
	}

	return chainID, nil
}

/*
Checkpoint RPC
Path: /chains/<chain_id>/checkpoint (GET)
Description:  The current checkpoint for this chain.
*/
func (t *GoTezos) Checkpoint() (Checkpoint, error) {
	resp, err := t.get("/chains/main/checkpoint")
	if err != nil {
		return Checkpoint{}, errors.Wrap(err, "failed to get checkpoint")
	}

	var c Checkpoint
	err = json.Unmarshal(resp, &c)
	if err != nil {
		return c, errors.Wrap(err, "failed to unmarshal checkpoint")
	}

	return c, nil
}

/*
InvalidBlocks RPC
Path: /chains/<chain_id>/invalid_blocks (GET)
Description: Lists blocks that have been declared invalid
along with the errors that led to them being declared invalid.
*/
func (t *GoTezos) InvalidBlocks() ([]string, error) {
	resp, err := t.get("/chains/main/invalid_blocks")
	if err != nil {
		return []string{}, errors.Wrap(err, "failed to get invalid blocks")
	}

	var blocks []string
	err = json.Unmarshal(resp, &blocks)
	if err != nil {
		return []string{}, errors.Wrap(err, "failed to unmarshal invalid blocks")
	}

	return blocks, nil
}

// TODO: /chains/<chain_id>/invalid_blocks/<block_hash> (GET) (DELETE)
