package gotezos

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"
)

// BlockService is a struct wrapper for all block functions
type BlockService struct {
	gt *GoTezos
}

// Block is a block returned by the Tezos RPC API.
type Block struct {
	Protocol   string               `json:"protocol"`
	ChainID    string               `json:"chain_id"`
	Hash       string               `json:"hash"`
	Header     StructHeader         `json:"header"`
	Metadata   StructMetadata       `json:"metadata"`
	Operations [][]StructOperations `json:"operations"`
}

// StructHeader is a header in a block returned by the Tezos RPC API.
type StructHeader struct {
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

// StructMetadata is the Metadata in a block returned by the Tezos RPC API.
type StructMetadata struct {
	Protocol               string                         `json:"protocol"`
	NextProtocol           string                         `json:"next_protocol"`
	TestChainStatus        StructTestChainStatus          `json:"test_chain_status"`
	MaxOperationsTTL       int                            `json:"max_operations_ttl"`
	MaxOperationDataLength int                            `json:"max_operation_data_length"`
	MaxBlockHeaderLength   int                            `json:"max_block_header_length"`
	MaxOperationListLength []StructMaxOperationListLength `json:"max_operation_list_length"`
	Baker                  string                         `json:"baker"`
	Level                  StructLevel                    `json:"level"`
	VotingPeriodKind       string                         `json:"voting_period_kind"`
	NonceHash              interface{}                    `json:"nonce_hash"`
	ConsumedGas            string                         `json:"consumed_gas"`
	Deactivated            []string                       `json:"deactivated"`
	BalanceUpdates         []StructBalanceUpdates         `json:"balance_updates"`
}

// StructTestChainStatus is the TestChainStatus found in the Metadata of a block returned by the Tezos RPC API.
type StructTestChainStatus struct {
	Status string `json:"status"`
}

// StructMaxOperationListLength is the MaxOperationListLength found in the Metadata of a block returned by the Tezos RPC API.
type StructMaxOperationListLength struct {
	MaxSize int `json:"max_size"`
	MaxOp   int `json:"max_op,omitempty"`
}

// StructLevel the Level found in the Metadata of a block returned by the Tezos RPC API.
type StructLevel struct {
	Level                int  `json:"level"`
	LevelPosition        int  `json:"level_position"`
	Cycle                int  `json:"cycle"`
	CyclePosition        int  `json:"cycle_position"`
	VotingPeriod         int  `json:"voting_period"`
	VotingPeriodPosition int  `json:"voting_period_position"`
	ExpectedCommitment   bool `json:"expected_commitment"`
}

// StructBalanceUpdates is the BalanceUpdates found in the Metadata of a block returned by the Tezos RPC API.
type StructBalanceUpdates struct {
	Kind     string `json:"kind"`
	Contract string `json:"contract,omitempty"`
	Change   string `json:"change"`
	Category string `json:"category,omitempty"`
	Delegate string `json:"delegate,omitempty"`
	Level    int    `json:"level,omitempty"`
}

// StructOperations is the Operations found in a block returned by the Tezos RPC API.
type StructOperations struct {
	Protocol  string           `json:"protocol"`
	ChainID   string           `json:"chain_id"`
	Hash      string           `json:"hash"`
	Branch    string           `json:"branch"`
	Contents  []StructContents `json:"contents"`
	Signature string           `json:"signature"`
}

// StructContents is the Contents found in a operation of a block returned by the Tezos RPC API.
type StructContents struct {
	Kind             string           `json:"kind"`
	Source           string           `json:"source"`
	Fee              string           `json:"fee"`
	Counter          string           `json:"counter"`
	GasLimit         string           `json:"gas_limit"`
	StorageLimit     string           `json:"storage_limit"`
	Amount           string           `json:"amount"`
	Destination      string           `json:"destination"`
	Delegate         string           `json:"delegate"`
	Phk              string           `json:"phk"`
	Secret           string           `json:"secret"`
	Level            int              `json:"level"`
	ManagerPublicKey string           `json:"managerPubkey"`
	Balance          string           `json:"balance"`
	Metadata         ContentsMetadata `json:"metadata"`
}

// ContentsMetadata is the Metadata found in the Contents in a operation of a block returned by the Tezos RPC API.
type ContentsMetadata struct {
	BalanceUpdates []StructBalanceUpdates `json:"balance_updates"`
	Slots          []int                  `json:"slots"`
}

// NewBlockService creates a new BlockService
func (gt *GoTezos) newBlockService() *BlockService {
	return &BlockService{gt: gt}
}

// GetHead returns the head block
func (b *BlockService) GetHead() (Block, error) {
	var block Block
	query := "/chains/main/blocks/head"
	resp, err := b.gt.Get(query, nil)
	if err != nil {
		return block, err
	}

	block, err = block.unmarshalJSON(resp)
	if err != nil {
		return block, err
	}

	return block, nil
}

// Get returns a Block at a specific level or hash
func (b *BlockService) Get(id interface{}) (Block, error) {
	var block Block

	// Get current head block
	head, err := b.GetHead()
	if err != nil {
		return block, fmt.Errorf("cannot get block %v: %v", id, err)
	}

	query := "/chains/main/blocks/"
	switch v := id.(type) {
	case int:
		diff := strconv.Itoa(head.Header.Level - v)
		query = query + head.Hash + "~" + diff
	case string:
		query = query + v
	default:
		return block, fmt.Errorf("invalid block id type, must be string or int")
	}

	resp, err := b.gt.Get(query, nil)
	if err != nil {
		return block, err
	}

	block, err = block.unmarshalJSON(resp)
	if err != nil {
		return block, err
	}

	return block, nil
}

// UnmarshalJSON unmarshals the bytes received as a parameter, into the type Block.
func (b *Block) unmarshalJSON(v []byte) (Block, error) {
	block := Block{}
	err := json.Unmarshal(v, &block)
	if err != nil {
		log.Println("Could not get unMarshal bytes into block: " + err.Error())
		return block, err
	}
	return block, nil
}
