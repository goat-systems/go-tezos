package gotezos

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

/*
BigInt Wrapper
Description: BigInt wraps go's big.Int.
*/
type BigInt struct {
	big.Int
}

/*
UnmarshalJSON Function
Description: Implements the json.Marshaler interface for BigInt

Parameters:
	b:
		The byte representation of a BigInt.
*/
func (i *BigInt) UnmarshalJSON(b []byte) error {
	var val string
	err := json.Unmarshal(b, &val)
	if err != nil {
		return err
	}

	i.SetString(val, 10)

	return nil
}

/*
MarshalJSON Function
Description: Implements the json.Marshaler interface for BigInt
*/
func (i *BigInt) MarshalJSON() ([]byte, error) {
	return i.MarshalText()

}

/*
Block Resp
RPC: /chains/<chain_id>/blocks/<block_id> (<dyn>)
Link: https://tezos.gitlab.io/api/rpc.html#get-block-id-context-contracts-contract-id-balance
*/
type Block struct {
	Protocol   string         `json:"protocol"`
	ChainID    string         `json:"chain_id"`
	Hash       string         `json:"hash"`
	Header     Header         `json:"header"`
	Metadata   Metadata       `json:"metadata"`
	Operations [][]Operations `json:"operations"`
}

/*
Header <block>
RPC: /chains/<chain_id>/blocks/<block_id> (<dyn>)
Link: https://tezos.gitlab.io/api/rpc.html#get-block-id-context-contracts-contract-id-balance
*/
type Header struct {
	Level            int       `json:"level"`
	Proto            int       `json:"proto"`
	Predecessor      string    `json:"Predecessor"`
	Timestamp        time.Time `json:"timestamp"`
	ValidationPass   int       `json:"validation_pass"`
	OperationsHash   string    `json:"operations_hash"`
	Fitness          []string  `json:"fitness"`
	Context          string    `json:"context"`
	Priority         int       `json:"priority"`
	ProofOfWorkNonce string    `json:"proof_of_work_nonce"`
	Signature        string    `json:"signature"`
}

/*
Metadata <block>
RPC: /chains/<chain_id>/blocks/<block_id> (<dyn>)
Link: https://tezos.gitlab.io/api/rpc.html#get-block-id-context-contracts-contract-id-balance
*/
type Metadata struct {
	Protocol               string                   `json:"protocol"`
	NextProtocol           string                   `json:"next_protocol"`
	TestChainStatus        TestChainStatus          `json:"test_chain_status"`
	MaxOperationsTTL       int                      `json:"max_operations_ttl"`
	MaxOperationDataLength int                      `json:"max_operation_data_length"`
	MaxBlockHeaderLength   int                      `json:"max_block_header_length"`
	MaxOperationListLength []MaxOperationListLength `json:"max_operation_list_length"`
	Baker                  string                   `json:"baker"`
	Level                  Level                    `json:"level"`
	VotingPeriodKind       string                   `json:"voting_period_kind"`
	NonceHash              interface{}              `json:"nonce_hash"`
	ConsumedGas            string                   `json:"consumed_gas"`
	Deactivated            []string                 `json:"deactivated"`
	BalanceUpdates         []BalanceUpdates         `json:"balance_updates"`
}

/*
TestChainStatus <block>
RPC: /chains/<chain_id>/blocks/<block_id> (<dyn>)
Link: https://tezos.gitlab.io/api/rpc.html#get-block-id-context-contracts-contract-id-balance
*/
type TestChainStatus struct {
	Status string `json:"status"`
}

/*
MaxOperationListLength <block>
RPC: /chains/<chain_id>/blocks/<block_id> (<dyn>)
Link: https://tezos.gitlab.io/api/rpc.html#get-block-id-context-contracts-contract-id-balance
*/
type MaxOperationListLength struct {
	MaxSize int `json:"max_size"`
	MaxOp   int `json:"max_op,omitempty"`
}

/*
Level <block>
RPC: /chains/<chain_id>/blocks/<block_id> (<dyn>)
Link: https://tezos.gitlab.io/api/rpc.html#get-block-id-context-contracts-contract-id-balance
*/
type Level struct {
	Level                int  `json:"level"`
	LevelPosition        int  `json:"level_position"`
	Cycle                int  `json:"cycle"`
	CyclePosition        int  `json:"cycle_position"`
	VotingPeriod         int  `json:"voting_period"`
	VotingPeriodPosition int  `json:"voting_period_position"`
	ExpectedCommitment   bool `json:"expected_commitment"`
}

/*
BalanceUpdates <block>
RPC: /chains/<chain_id>/blocks/<block_id> (<dyn>)
Link: https://tezos.gitlab.io/api/rpc.html#get-block-id-context-contracts-contract-id-balance
*/
type BalanceUpdates struct {
	Kind     string `json:"kind"`
	Contract string `json:"contract,omitempty"`
	Change   BigInt `json:"change"`
	Category string `json:"category,omitempty"`
	Delegate string `json:"delegate,omitempty"`
	Cycle    int    `json:"cycle,omitempty"`
	Level    int    `json:"level,omitempty"`
}

/*
OperationResult <block>
RPC: /chains/<chain_id>/blocks/<block_id> (<dyn>)
Link: https://tezos.gitlab.io/api/rpc.html#get-block-id-context-contracts-contract-id-balance
*/
type OperationResult struct {
	Status      string  `json:"status"`
	ConsumedGas BigInt  `json:"consumed_gas,omitempty"`
	Errors      []Error `json:"errors,omitempty"`
}

/*
Operations <block>
RPC: /chains/<chain_id>/blocks/<block_id> (<dyn>)
Link: https://tezos.gitlab.io/api/rpc.html#get-block-id-context-contracts-contract-id-balance
*/
type Operations struct {
	Protocol  string     `json:"protocol"`
	ChainID   string     `json:"chain_id"`
	Hash      string     `json:"hash"`
	Branch    string     `json:"branch"`
	Contents  []Contents `json:"contents"`
	Signature string     `json:"signature"`
}

/*
Contents <block>
RPC: /chains/<chain_id>/blocks/<block_id> (<dyn>)
Link: https://tezos.gitlab.io/api/rpc.html#get-block-id-context-contracts-contract-id-balance
*/
type Contents struct {
	Kind             string            `json:"kind,omitempty"`
	Source           string            `json:"source,omitempty"`
	Fee              BigInt            `json:"fee,omitempty"`
	Counter          BigInt            `json:"counter,omitempty"`
	GasLimit         BigInt            `json:"gas_limit,omitempty"`
	StorageLimit     BigInt            `json:"storage_limit,omitempty"`
	Amount           BigInt            `json:"amount,omitempty"`
	Destination      string            `json:"destination,omitempty"`
	Delegate         string            `json:"delegate,omitempty"`
	Phk              string            `json:"phk,omitempty"`
	Secret           string            `json:"secret,omitempty"`
	Level            int               `json:"level,omitempty"`
	ManagerPublicKey string            `json:"managerPubkey,omitempty"`
	Balance          BigInt            `json:"balance,omitempty"`
	Period           int               `json:"period,omitempty"`
	Proposal         string            `json:"proposal,omitempty"`
	Proposals        []string          `json:"proposals,omitempty"`
	Ballot           string            `json:"ballot,omitempty"`
	Metadata         *ContentsMetadata `json:"metadata,omitempty"`
}

/*
ContentsMetadata <block>
RPC: /chains/<chain_id>/blocks/<block_id> (<dyn>)
Link: https://tezos.gitlab.io/api/rpc.html#get-block-id-context-contracts-contract-id-balance
*/
type ContentsMetadata struct {
	BalanceUpdates  []BalanceUpdates `json:"balance_updates"`
	OperationResult *OperationResult `json:"operation_result,omitempty"`
	Slots           []int            `json:"slots"`
}

/*
Error <block>
RPC: /chains/<chain_id>/blocks/<block_id> (<dyn>)
Link: https://tezos.gitlab.io/api/rpc.html#get-block-id-context-contracts-contract-id-balance
*/
type Error struct {
	Kind string `json:"kind"`
	ID   string `json:"id"`
}

/*
Blocks RPC
Path: /chains/<chain_id>/blocks (GET)
Link: https://tezos.gitlab.io/api/rpc.html#get-chains-chain-id-blocks
Description:  Lists known heads of the blockchain sorted with decreasing fitness.
Optional arguments allows to returns the list of predecessors for known heads or
the list of predecessors for a given list of blocks.

Parameters:
	opts:
		length = <int> : The requested number of predecessors to returns (per requested head).
		head = <block_hash> : An empty argument requests blocks from the current heads. A non empty list allow to request specific fragment of the chain.
		min_date = <date> : When `min_date` is provided, heads with a timestamp before `min_date` are filtered out
*/

/*
Head RPC
Path: /chains/<chain_id>/blocks/head (GET)
Link: https://tezos.gitlab.io/api/rpc.html#get-chains-chain-id-blocks
Description: All the information about the head block.
*/
func (t *GoTezos) Head() (Block, error) {
	resp, err := t.get("/chains/main/blocks/head")
	if err != nil {
		return Block{}, errors.Wrapf(err, "could not get head block")
	}

	var block Block
	err = json.Unmarshal(resp, &block)
	if err != nil {
		return block, errors.Wrapf(err, "could not get head block")
	}

	return block, nil
}

/*
Block RPC
Path: /chains/<chain_id>/blocks/<block_id> (GET)
Link: https://tezos.gitlab.io/api/rpc.html#get-chains-chain-id-blocks
Description:  All the information about block.

Parameters:
	id:
		hash = <string> : The block hash.
		level = <int> : The block level.
*/
func (t *GoTezos) Block(id interface{}) (Block, error) {
	blockID, err := idToString(id)
	if err != nil {
		return Block{}, errors.Wrapf(err, "could not get block '%s'", blockID)
	}

	resp, err := t.get(fmt.Sprintf("/chains/main/blocks/%s", blockID))
	if err != nil {
		return Block{}, errors.Wrapf(err, "could not get block '%s'", blockID)
	}

	var block Block
	err = json.Unmarshal(resp, &block)
	if err != nil {
		return block, errors.Wrapf(err, "could not get block '%s'", blockID)
	}

	return block, nil
}

/*
OperationHashes RPC
Path: ../<block_id>/operation_hashes (GET)
Link: https://tezos.gitlab.io/api/rpc.html#get-block-id-context-contracts-contract-id-balance
Description: The hashes of all the operations included in the block.

Parameters:
	blockhash:
		The hash of block (height) of which you want to make the query.
*/
func (t *GoTezos) OperationHashes(blockhash string) ([]string, error) {
	resp, err := t.get(fmt.Sprintf("/chains/main/blocks/%s/operation_hashes", blockhash))
	if err != nil {
		return []string{}, errors.Wrapf(err, "could not get operation hashes")
	}

	var operations []string
	err = json.Unmarshal(resp, &operations)
	if err != nil {
		return []string{}, errors.Wrapf(err, "could not unmarshal operation hashes")
	}

	return operations, nil
}

func idToString(id interface{}) (string, error) {
	switch v := id.(type) {
	case int:
		return strconv.Itoa(v), nil
	case string:
		return v, nil
	default:
		return "", errors.Errorf("id must be block level (int) or block hash (string)")
	}
}
