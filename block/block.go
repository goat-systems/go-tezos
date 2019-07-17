package block

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/pkg/errors"

	tzc "github.com/DefinitelyNotAGoat/go-tezos/client"
)

// BlockService is a struct wrapper for all block functions
type BlockService struct {
	tzclient tzc.TezosClient
}

// Block is a block returned by the Tezos RPC API.
type Block struct {
	Protocol   string         `json:"protocol"`
	ChainID    string         `json:"chain_id"`
	Hash       string         `json:"hash"`
	Header     Header         `json:"header"`
	Metadata   Metadata       `json:"metadata"`
	Operations [][]Operations `json:"operations"`
}

// Header is a header in a block returned by the Tezos RPC API.
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

// Metadata is the Metadata in a block returned by the Tezos RPC API.
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

// TestChainStatus is the TestChainStatus found in the Metadata of a block returned by the Tezos RPC API.
type TestChainStatus struct {
	Status string `json:"status"`
}

// MaxOperationListLength is the MaxOperationListLength found in the Metadata of a block returned by the Tezos RPC API.
type MaxOperationListLength struct {
	MaxSize int `json:"max_size"`
	MaxOp   int `json:"max_op,omitempty"`
}

// Level the Level found in the Metadata of a block returned by the Tezos RPC API.
type Level struct {
	Level                int  `json:"level"`
	LevelPosition        int  `json:"level_position"`
	Cycle                int  `json:"cycle"`
	CyclePosition        int  `json:"cycle_position"`
	VotingPeriod         int  `json:"voting_period"`
	VotingPeriodPosition int  `json:"voting_period_position"`
	ExpectedCommitment   bool `json:"expected_commitment"`
}

// BalanceUpdates is the BalanceUpdates found in the Metadata of a block returned by the Tezos RPC API.
type BalanceUpdates struct {
	Kind     string `json:"kind"`
	Contract string `json:"contract,omitempty"`
	Change   string `json:"change"`
	Category string `json:"category,omitempty"`
	Delegate string `json:"delegate,omitempty"`
	Cycle    int    `json:"cycle,omitempty"`
	Level    int    `json:"level,omitempty"`
}

// OperationResult is the OperationResult found in metadata of block returned by the Tezos RPC API.
type OperationResult struct {
	Status      string  `json:"status"`
	ConsumedGas string  `json:"consumed_gas,omitempty"`
	Errors      []Error `json:"errors,omitempty"`
}

// Operations is the Operations found in a block returned by the Tezos RPC API.
type Operations struct {
	Protocol  string     `json:"protocol"`
	ChainID   string     `json:"chain_id"`
	Hash      string     `json:"hash"`
	Branch    string     `json:"branch"`
	Contents  []Contents `json:"contents"`
	Signature string     `json:"signature"`
}

// Contents is the Contents found in a operation of a block returned by the Tezos RPC API.
type Contents struct {
	Kind             string            `json:"kind,omitempty"`
	Source           string            `json:"source,omitempty"`
	Fee              string            `json:"fee,omitempty"`
	Counter          string            `json:"counter,omitempty"`
	GasLimit         string            `json:"gas_limit,omitempty"`
	StorageLimit     string            `json:"storage_limit,omitempty"`
	Amount           string            `json:"amount,omitempty"`
	Destination      string            `json:"destination,omitempty"`
	Delegate         string            `json:"delegate,omitempty"`
	Phk              string            `json:"phk,omitempty"`
	Secret           string            `json:"secret,omitempty"`
	Level            int               `json:"level,omitempty"`
	ManagerPublicKey string            `json:"managerPubkey,omitempty"`
	Balance          string            `json:"balance,omitempty"`
	Period           int               `json:"period,omitempty"`
	Proposal         string            `json:"proposal,omitempty"`
	Proposals        []string          `json:"proposals,omitempty"`
	Ballot           string            `json:"ballot,omitempty"`
	Metadata         *ContentsMetadata `json:"metadata,omitempty"`
}

// ContentsMetadata is the Metadata found in the Contents in a operation of a block returned by the Tezos RPC API.
type ContentsMetadata struct {
	BalanceUpdates  []BalanceUpdates `json:"balance_updates"`
	OperationResult *OperationResult `json:"operation_result,omitempty"`
	Slots           []int            `json:"slots"`
}

// Error is the Error found in the OperationResult in a metadata of operation of a block returned by the Tezos RPC API.
type Error struct {
	Kind string `json:"kind"`
	ID   string `json:"id"`
}

// NewBlockService creates a new BlockService
func NewBlockService(tzclient tzc.TezosClient) *BlockService {
	return &BlockService{tzclient: tzclient}
}

// GetHead returns the head block
func (b *BlockService) GetHead() (Block, error) {
	var block Block
	query := "/chains/main/blocks/head"
	resp, err := b.tzclient.Get(query, nil)
	if err != nil {
		return block, errors.Wrapf(err, "could not get head block '%s'", query)
	}

	block, err = block.unmarshalJSON(resp)
	if err != nil {
		return block, errors.Wrapf(err, "could not get head block '%s'", query)
	}

	return block, nil
}

// Get returns a Block at a specific level or hash
func (b *BlockService) Get(id interface{}) (Block, error) {
	var block Block

	query := "/chains/main/blocks/"

	blockID, err := b.IDToString(id)
	if err != nil {
		return block, err
	}

	query += blockID

	resp, err := b.tzclient.Get(query, nil)
	if err != nil {
		return block, errors.Wrap(err, "could not get block '%s'")
	}

	block, err = block.unmarshalJSON(resp)
	if err != nil {
		return block, errors.Wrap(err, "could not get block '%s'")
	}

	return block, nil
}

// IDToString returns a queryable block reference for a specific level or hash
func (b *BlockService) IDToString(id interface{}) (string, error) {
	switch v := id.(type) {
	case int:
		return strconv.Itoa(v), nil
	case string:
		return v, nil
	default:
		return "", errors.Errorf("invalid block id type, must be string or int")
	}
}

// UnmarshalJSON unmarshals the bytes received as a parameter, into the type Block.
func (b *Block) unmarshalJSON(v []byte) (Block, error) {
	block := Block{}
	err := json.Unmarshal(v, &block)
	if err != nil {
		return block, errors.Wrap(err, "could not unmarshal bytes to Block")
	}
	return block, nil
}
