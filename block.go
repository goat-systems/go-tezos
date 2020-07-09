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
Int Wrapper
Description: Int wraps go's big.Int.
*/
type Int struct {
	Big *big.Int
}

// NewInt returns a pointer GoTezos's wrapper Int
func NewInt(i int) *Int {
	return &Int{Big: big.NewInt(int64(i))}
}

func newInt(bigintstring []byte) (*Int, error) {
	i := &Int{}
	err := i.UnmarshalJSON(bigintstring)
	return i, err
}

/*
UnmarshalJSON implements the json.Marshaler interface for BigInt

Parameters:

	b:
		The byte representation of a BigInt.
*/
func (i *Int) UnmarshalJSON(b []byte) error {
	var val string
	err := json.Unmarshal(b, &val)
	if err != nil {
		return err
	}
	i.Big = big.NewInt(0)
	i.Big.SetString(val, 10)

	return nil
}

/*
MarshalJSON implements the json.Marshaler interface for BigInt
*/
func (i *Int) MarshalJSON() ([]byte, error) {
	val, err := i.Big.MarshalText()
	if err != nil {
		return nil, err
	}

	return []byte(fmt.Sprintf("\"%s\"", val)), nil
}

/*
Block represents a Tezos block.

RPC:
	/chains/<chain_id>/blocks/<block_id> (<dyn>)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-context-contracts-contract-id-balance
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
Header represents the header in a Tezos block

RPC:
	/chains/<chain_id>/blocks/<block_id> (<dyn>)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-context-contracts-contract-id-balance
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
	SeedNonceHash    string    `json:"seed_nonce_hash"`
	Signature        string    `json:"signature"`
}

/*
Metadata represents the metadata in a Tezos block

RPC:
	/chains/<chain_id>/blocks/<block_id> (<dyn>)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-context-contracts-contract-id-balance
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
TestChainStatus represents the testchainstatus in a Tezos block

RPC:
	/chains/<chain_id>/blocks/<block_id> (<dyn>)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-context-contracts-contract-id-balance
*/
type TestChainStatus struct {
	Status     string    `json:"status"`
	Protocol   string    `json:"protocol"`
	ChainID    string    `json:"chain_id"`
	Genesis    string    `json:"genesis"`
	Expiration time.Time `json:"expiration"`
}

/*
MaxOperationListLength represents the maxoperationlistlength in a Tezos block

RPC:
	/chains/<chain_id>/blocks/<block_id> (<dyn>)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-context-contracts-contract-id-balance
*/
type MaxOperationListLength struct {
	MaxSize int `json:"max_size"`
	MaxOp   int `json:"max_op,omitempty"`
}

/*
Level represents the level in a Tezos block

RPC:
	/chains/<chain_id>/blocks/<block_id> (<dyn>)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-context-contracts-contract-id-balance
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
BalanceUpdates represents the balance updates in a Tezos block

RPC:
	/chains/<chain_id>/blocks/<block_id> (<dyn>)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-context-contracts-contract-id-balance
*/
type BalanceUpdates struct {
	Kind     string `json:"kind"`
	Contract string `json:"contract,omitempty"`
	Change   *Int   `json:"change"`
	Category string `json:"category,omitempty"`
	Delegate string `json:"delegate,omitempty"`
	Cycle    int    `json:"cycle,omitempty"`
	Level    int    `json:"level,omitempty"`
}

/*
OperationResult represents the operation result in a Tezos block

RPC:
	/chains/<chain_id>/blocks/<block_id> (<dyn>)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-context-contracts-contract-id-balance
*/
type OperationResult struct {
	BalanceUpdates      []BalanceUpdates `json:"balance_updates"`
	OriginatedContracts []string         `json:"originated_contracts"`
	Status              string           `json:"status"`
	ConsumedGas         *Int             `json:"consumed_gas,omitempty"`
	Errors              []Error          `json:"errors,omitempty"`
}

/*
Operations represents the operations in a Tezos block

RPC:
	/chains/<chain_id>/blocks/<block_id> (<dyn>)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-context-contracts-contract-id-balance
*/
type Operations struct {
	Protocol  string     `json:"protocol,omitempty"`
	ChainID   string     `json:"chain_id,omitempty"`
	Hash      string     `json:"hash,omitempty"`
	Branch    string     `json:"branch"`
	Contents  []Contents `json:"contents"`
	Signature string     `json:"signature,omitempty"`
}

/*
Contents represents the contents in a Tezos operations

RPC:
	/chains/<chain_id>/blocks/<block_id> (<dyn>)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-context-contracts-contract-id-balance
*/
type Contents struct {
	Endorsements              []Endorsement
	SeedNonceRevelations      []SeedNonceRevelation
	DoubleEndorsementEvidence []DoubleEndorsementEvidence
	DoubleBakingEvidence      []DoubleBakingEvidence
	AccountActivations        []AccountActivation
	Proposals                 []Proposal
	Ballots                   []Ballot
	Reveals                   []Reveal
}

type Endorsement struct {
	Kind     string `json:"kind"`
	Level    int    `json:"level"`
	Metadata *struct {
		BalanceUpdates []BalanceUpdates `json:"balance_updates"`
		Delegate       string           `json:"delegate"`
		Slots          []int            `json:"slots"`
	} `json:"metadata"`
}

type SeedNonceRevelation struct {
	Kind     string `json:"kind"`
	Level    int    `json:"level"`
	Nonce    string `json:"nonce"`
	Metadata *struct {
		BalanceUpdates []BalanceUpdates `json:"balance_updates"`
	} `json:"metadata"`
}

type DoubleEndorsementEvidence struct {
	Kind     string             `json:"kind"`
	Op1      InlinedEndorsement `json:"Op1"`
	Op2      InlinedEndorsement `json:"Op2"`
	Metadata *struct {
		BalanceUpdates []BalanceUpdates `json:"balance_updates"`
	} `json:"metadata"`
}

type InlinedEndorsement struct {
	Branch     string `json:"branch"`
	Operations struct {
		Kind  string `json:"kind"`
		Level int    `json:"level"`
	} `json:"operations"`
	Signature string `json:"signature"`
}

type DoubleBakingEvidence struct {
	Kind     string      `json:"kind"`
	Bh1      BlockHeader `json:"bh1"`
	Bh2      BlockHeader `json:"bh2"`
	Metadata *struct {
		BalanceUpdates []BalanceUpdates `json:"balance_updates"`
	} `json:"metadata"`
}

type BlockHeader struct {
	Level            int       `json:"level"`
	Proto            int       `json:"proto"`
	Predecessor      string    `json:"predecessor"`
	Timestamp        time.Time `json:"timestamp"`
	ValidationPass   int       `json:"validation_pass"`
	OperationsHash   string    `json:"operations_hash"`
	Fitness          string    `json:"fitness"`
	Context          string    `json:"context"`
	Priority         int       `json:"priority"`
	ProofOfWorkNonce string    `json:"proof_of_work_nonce"`
	SeedNonceHash    string    `json:"seed_nonce_hash"`
	Signature        string    `json:"signature"`
}

type AccountActivation struct {
	Kind     string `json:"kind"`
	Pkh      string `json:"pkh"`
	Secret   string `json:"secret"`
	Metadata *struct {
		BalanceUpdates []BalanceUpdates `json:"balance_updates"`
	} `json:"metadata"`
}

type Proposal struct {
	Kind      string   `json:"kind"`
	Source    string   `json:"source"`
	Period    int      `json:"period"`
	Proposals []string `json:"proposals"`
}

type Ballot struct {
	Kind     string `json:"kind"`
	Source   string `json:"source"`
	Period   int    `json:"period"`
	Proposal string `json:"proposal"`
	Ballot   string `json:"ballot"`
}

type Reveal struct {
	Kind         string `json:"kind"`
	Source       string `json:"source"`
	Fee          Int    `json:"fee"`
	Counter      int    `json:"counter"`
	GasLimit     Int    `json:"gas_limit"`
	StorageLimit Int    `json:"storage_limit"`
	PublicKey    string `json:"public_key"`
	Metadata     *struct {
		BalanceUpdates          []BalanceUpdates         `json:"balance_updates"`
		OperationResult         OperationResultReveal    `json:"operation_result"`
		InternalOperationResult InternalOperationResults `json:"internal_operation_result"`
	} `json:"metadata"`
}

type OperationResultReveal struct {
	Status      string     `json:"status"`
	ConsumedGas Int        `json:"consumed_gas"`
	Errors      []RPCError `json:"rpc_error"`
}

/*
OperationResultTransfer represents $operation.alpha.operation_result.transaction in the tezos block schema
See: tezos-client RPC format GET /chains/main/blocks/head
*/
type OperationResultTransfer struct {
	OperationResultTransferApplied     *OperationResultTransferApplied
	OperationResultTransferFailed      *OperationResultTransferFailed
	OperationResultTransferSkipped     *OperationResultTransferSkipped
	OperationResultTransferBacktracked *OperationResultTransferBacktracked
}

/*
UnmarshalJSON implements the json.UnmarshalJSON interface for OperationResultTransfer
*/
func (b *OperationResultTransfer) UnmarshalJSON(v []byte) error {
	m := map[string]interface{}{}

	status, ok := m["status"]
	if !ok {
		return errors.New("failed to unmarshal OperationResultTransfer")
	}

	if status == "applied" {
		var operationResultTransferApplied OperationResultTransferApplied
		if err := json.Unmarshal(v, &operationResultTransferApplied); err != nil {
			return errors.Wrap(err, "failed to unmarshal OperationResultTransfer")
		}

		b.OperationResultTransferApplied = &operationResultTransferApplied
		return nil
	} else if status == "failed" {
		var operationResultTransferFailed OperationResultTransferFailed
		if err := json.Unmarshal(v, &operationResultTransferFailed); err != nil {
			return errors.Wrap(err, "failed to unmarshal OperationResultTransfer")
		}

		b.OperationResultTransferFailed = &operationResultTransferFailed
		return nil
	} else if status == "skipped" {
		var operationResultTransferSkipped OperationResultTransferSkipped
		if err := json.Unmarshal(v, &operationResultTransferSkipped); err != nil {
			return errors.Wrap(err, "failed to unmarshal OperationResultTransfer")
		}

		b.OperationResultTransferSkipped = &operationResultTransferSkipped
		return nil
	} else if status == "backtracked" {
		var operationResultTransferBacktracked OperationResultTransferBacktracked
		if err := json.Unmarshal(v, &operationResultTransferBacktracked); err != nil {
			return errors.Wrap(err, "failed to unmarshal OperationResultTransfer")
		}

		b.OperationResultTransferBacktracked = &operationResultTransferBacktracked
		return nil
	}

	return errors.New("failed to unmarshal OperationResultTransfer")
}

/*
MarshalJSON implements the json.Marshaler interface for OperationResultTransfer
*/
func (b *OperationResultTransfer) MarshalJSON() ([]byte, error) {
	if b.OperationResultTransferApplied != nil {
		v, err := json.Marshal(b.OperationResultTransferApplied)
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal OperationResultTransfer")
		}

		return v, nil
	} else if b.OperationResultTransferFailed != nil {
		v, err := json.Marshal(b.OperationResultTransferFailed)
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal OperationResultTransfer")
		}

		return v, nil
	} else if b.OperationResultTransferSkipped != nil {
		v, err := json.Marshal(b.OperationResultTransferSkipped)
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal OperationResultTransfer")
		}

		return v, nil
	} else if b.OperationResultTransferBacktracked != nil {
		v, err := json.Marshal(b.OperationResultTransferBacktracked)
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal OperationResultTransfer")
		}

		return v, nil
	}

	return nil, nil
}

/*
OperationResultTransferApplied represents $operation.alpha.operation_result.transaction in the tezos block schema
See: tezos-client RPC format GET /chains/main/blocks/head
*/
type OperationResultTransferApplied struct {
	Status                       string                          `json:"status"`
	Storage                      *MichelineMichelsonV1Expression `json:"storage,omitempty"`
	BigMapDiff                   *BigMapDiff                     `json:"big_map_diff,omitempty"`
	BalanceUpdates               *BalanceUpdates                 `json:"balance_updates,omitempty"`
	OriginatedContracts          []string                        `json:"originated_contracts,omitempty"`
	ConsumedGas                  *Int                            `json:"consumed_gas,omitempty"`
	StorageSize                  *Int                            `json:"storage_size,omitempty"`
	PaidStorageSizeDiff          *Int                            `json:"paid_storage_size_diff,omitempty"`
	AllocatedDestinationContract *bool                           `json:"allocated_destination_contract,omitempty"`
}

/*
OperationResultTransferFailed represents $operation.alpha.operation_result.transaction in the tezos block schema
See: tezos-client RPC format GET /chains/main/blocks/head
*/
type OperationResultTransferFailed struct {
	Status string     `json:"status"`
	Errors []RPCError `json:"errors"`
}

/*
OperationResultTransferSkipped represents $operation.alpha.operation_result.transaction in the tezos block schema
See: tezos-client RPC format GET /chains/main/blocks/head
*/
type OperationResultTransferSkipped struct {
	Status string `json:"status"`
}

/*
OperationResultTransferBacktracked represents $operation.alpha.operation_result.transaction in the tezos block schema
See: tezos-client RPC format GET /chains/main/blocks/head
*/
type OperationResultTransferBacktracked struct {
	Status                       string                          `json:"status"`
	Errors                       []RPCError                      `json:"errors,omitempty"`
	Storage                      *MichelineMichelsonV1Expression `json:"storage,omitempty"`
	BigMapDiff                   *BigMapDiff                     `json:"big_map_diff,omitempty"`
	BalanceUpdates               *BalanceUpdates                 `json:"balance_updates,omitempty"`
	OriginatedContracts          []string                        `json:"originated_contracts,omitempty"`
	ConsumedGas                  *Int                            `json:"consumed_gas,omitempty"`
	StorageSize                  *Int                            `json:"storage_size,omitempty"`
	PaidStorageSizeDiff          *Int                            `json:"paid_storage_size_diff,omitempty"`
	AllocatedDestinationContract *bool                           `json:"allocated_destination_contract,omitempty"`
}

/*
BigMapDiff represents $contract.big_map_diff in the tezos block schema
See: tezos-client RPC format GET /chains/main/blocks/head
*/
type BigMapDiff struct {
	Updates  []BigMapDiffUpdate
	Removals []BigMapDiffRemove
	Copies   []BigMapDiffCopy
	Alloc    []BigMapDiffAlloc
}

/*
UnmarshalJSON implements the json.UnmarshalJSON interface for BigMapDiff
*/
func (b *BigMapDiff) UnmarshalJSON(v []byte) error {
	data := [][]byte{}
	if err := json.Unmarshal(v, &data); err != nil {
		return errors.Wrap(err, "failed to unmarshal BigMapDiff")
	}

	var bigMapDiff BigMapDiff

	for _, d := range data {
		m := map[string]interface{}{}
		action, ok := m["action"]
		if !ok {
			return errors.New("failed to unmarshal BigMapDiff")
		}

		if action == "update" {
			var update BigMapDiffUpdate
			if err := json.Unmarshal(d, &update); err != nil {
				return errors.Wrap(err, "failed to unmarshal BigMapDiff")
			}

			bigMapDiff.Updates = append(bigMapDiff.Updates, update)
		} else if action == "remove" {
			var remove BigMapDiffRemove
			if err := json.Unmarshal(d, &remove); err != nil {
				return errors.Wrap(err, "failed to unmarshal BigMapDiff")
			}

			bigMapDiff.Removals = append(bigMapDiff.Removals, remove)
		} else if action == "copy" {
			var copy BigMapDiffCopy
			if err := json.Unmarshal(d, &copy); err != nil {
				return errors.Wrap(err, "failed to unmarshal BigMapDiff")
			}

			bigMapDiff.Copies = append(bigMapDiff.Copies, copy)
		} else if action == "alloc" {
			var alloc BigMapDiffAlloc
			if err := json.Unmarshal(d, &alloc); err != nil {
				return errors.Wrap(err, "failed to unmarshal BigMapDiff")
			}

			bigMapDiff.Alloc = append(bigMapDiff.Alloc, alloc)
		}
	}

	b = &bigMapDiff

	return nil
}

/*
MarshalJSON implements the json.Marshaler interface for BigMapDiff
*/
func (b *BigMapDiff) MarshalJSON() ([]byte, error) {
	var bigMapDiff []interface{}
	for _, update := range b.Updates {
		bigMapDiff = append(bigMapDiff, update)
	}

	for _, remove := range b.Removals {
		bigMapDiff = append(bigMapDiff, remove)
	}

	for _, copy := range b.Copies {
		bigMapDiff = append(bigMapDiff, copy)
	}

	for _, alloc := range b.Alloc {
		bigMapDiff = append(bigMapDiff, alloc)
	}

	bigMapDiffBytes, err := json.Marshal(bigMapDiff)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal BigMapDiff")
	}

	return bigMapDiffBytes, nil
}

/*
BigMapDiffUpdate represents $contract.big_map_diff in the tezos block schema
See: tezos-client RPC format GET /chains/main/blocks/head
*/
type BigMapDiffUpdate struct {
	Action  string                          `json:"action"`
	BigMap  Int                             `json:"big_map"`
	KeyHash string                          `json:"key_hash"`
	Key     MichelineMichelsonV1Expression  `json:"key"`
	Value   *MichelineMichelsonV1Expression `json:"value,omitempty"`
}

/*
BigMapDiffRemove represents $contract.big_map_diff in the tezos block schema
See: tezos-client RPC format GET /chains/main/blocks/head
*/
type BigMapDiffRemove struct {
	Action string `json:"action"`
	BigMap Int    `json:"big_map"`
}

/*
BigMapDiffCopy represents $contract.big_map_diff in the tezos block schema
See: tezos-client RPC format GET /chains/main/blocks/head
*/
type BigMapDiffCopy struct {
	Action            string `json:"action"`
	SourceBigMap      Int    `json:"source_big_map"`
	DestinationBigMap Int    `json:"destination_big_map"`
}

/*
BigMapDiffAlloc represents $contract.big_map_diff in the tezos block schema
See: tezos-client RPC format GET /chains/main/blocks/head
*/
type BigMapDiffAlloc struct {
	Action    string                         `json:"action"`
	BigMap    Int                            `json:"big_map"`
	KeyType   MichelineMichelsonV1Expression `json:"key_type"`
	ValueType MichelineMichelsonV1Expression `json:"value_type"`
}

type InternalOperationResults struct {
	Reveals      []InternalOperationResultReveal
	Transactions []InternalOperationResultTransaction
}

type InternalOperationResultReveal struct {
	Kind      string                `json:"kind"`
	Source    string                `json:"source"`
	Nonce     int                   `json:"nonce"`
	PublicKey string                `json:"public_key"`
	Result    OperationResultReveal `json:"result"`
}

type InternalOperationResultTransaction struct {
	Kind        string `json:"kind"`
	Source      string `json:"source"`
	Nonce       int    `json:"nonce"`
	Amount      Int    `json:"amount"`
	Destination string `json:"destination"`
	Parameters  struct {
		Entrypoint string                         `json:"entrypoint"`
		Value      MichelineMichelsonV1Expression `json:"value"`
	} `json:"paramaters"`
}

type MichelineMichelsonV1Expression struct {
	Int                            *int
	String                         *string
	Bytes                          []byte
	MichelineMichelsonV1Expression []MichelineMichelsonV1Expression
	GenericPrimitive               GenericPrimitive
}

type GenericPrimitive struct {
	Prim   string
	Args   []MichelineMichelsonV1Expression
	Annots []string
}

func (c *Contents) equal(contents Contents) (bool, error) {
	x, err := json.Marshal(c)
	if err != nil {
		return false, errors.New("failed to compare")
	}

	y, err := json.Marshal(contents)
	if err != nil {
		return false, errors.New("failed to compare")
	}

	if string(x) == string(y) {
		return true, nil
	}

	return false, nil
}

/*
ContentsMetadata represents the contents metadata in a Tezos operations

RPC:
	/chains/<chain_id>/blocks/<block_id> (<dyn>)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-context-contracts-contract-id-balance
*/
type ContentsMetadata struct {
	BalanceUpdates           []BalanceUpdates            `json:"balance_updates"`
	OperationResult          *OperationResult            `json:"operation_result,omitempty"`
	Slots                    []int                       `json:"slots"`
	InternalOperationResults []*InternalOperationResults `json:"internal_operation_results,omitempty"`
}

/*
InternalOperationResults represents a field in contents metadata in a Tezos operations

RPC:
	/chains/<chain_id>/blocks/<block_id> (<dyn>)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-context-contracts-contract-id-balance
*/
// type InternalOperationResults struct {
// 	Kind        string           `json:"kind"`
// 	Source      string           `json:"source"`
// 	Nonce       uint64           `json:"nonce"`
// 	Amount      string           `json:"amount"`
// 	Destination string           `json:"destination"`
// 	Result      *OperationResult `json:"result"`
// }

/*
Error respresents an error for operation results

RPC:
	/chains/<chain_id>/blocks/<block_id> (<dyn>)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-context-contracts-contract-id-balance
*/
type Error struct {
	Kind string `json:"kind"`
	ID   string `json:"id"`
}

/*
BallotList represents a list of casted ballots in a block.

Path:
	../<block_id>/votes/ballot_list (GET)
Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-votes-ballot-list
*/
type BallotList []struct {
	PublicKeyHash string `json:"pkh"`
	Ballot        string `json:"ballot"`
}

/*
Ballots represents a ballot total.

Path:
	../<block_id>/votes/ballots (GET)
Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-votes-ballots
*/
type Ballots struct {
	Yay  int `json:"yay"`
	Nay  int `json:"nay"`
	Pass int `json:"pass"`
}

/*
Listings represents a list of delegates with their voting weight, in number of rolls.

Path:
	../<block_id>/votes/listings (GET)
Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-votes-listings
*/
type Listings []struct {
	PublicKeyHash string `json:"pkh"`
	Rolls         int    `json:"rolls"`
}

/*
Proposals represents a list of proposals with number of supporters.

Path:
	../<block_id>/votes/proposals (GET)
Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-votes-proposals
*/
type Proposals []struct {
	Hash       string
	Supporters int
}

/*
UnmarshalJSON implements the json.Marshaler interface for Proposals

Parameters:

	b:
		The byte representation of a Proposals.
*/
func (p *Proposals) UnmarshalJSON(b []byte) error {
	var out [][]interface{}
	if err := json.Unmarshal(b, &out); err != nil {
		return err
	}

	var proposals Proposals
	for _, x := range out {
		if len(x) != 2 {
			return errors.New("unexpected bytes")
		}

		hash := fmt.Sprintf("%v", x[0])
		supportersStr := fmt.Sprintf("%v", x[1])
		supporters, err := strconv.Atoi(supportersStr)
		if err != nil {
			return errors.New("unexpected bytes")
		}

		proposals = append(proposals, struct {
			Hash       string
			Supporters int
		}{
			Hash:       hash,
			Supporters: supporters,
		})
	}

	p = &proposals
	return nil
}

/*
Head gets all the information about the head block.

Path:
	/chains/<chain_id>/blocks/head (GET)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-chains-chain-id-blocks
*/
func (t *GoTezos) Head() (*Block, error) {
	resp, err := t.get("/chains/main/blocks/head")
	if err != nil {
		return &Block{}, errors.Wrapf(err, "could not get head block")
	}

	var block Block
	err = json.Unmarshal(resp, &block)
	if err != nil {
		return &block, errors.Wrapf(err, "could not get head block")
	}

	return &block, nil
}

/*
Block gets all the information about block.RPC

Path
	/chains/<chain_id>/blocks/<block_id> (GET)
Link
	https://tezos.gitlab.io/api/rpc.html#get-chains-chain-id-blocks

Parameters:

	id:
		hash = <string> : The block hash.
		level = <int> : The block level.
*/
func (t *GoTezos) Block(id interface{}) (*Block, error) {
	blockID, err := idToString(id)
	if err != nil {
		return &Block{}, errors.Wrapf(err, "could not get block '%s'", blockID)
	}

	resp, err := t.get(fmt.Sprintf("/chains/main/blocks/%s", blockID))
	if err != nil {
		return &Block{}, errors.Wrapf(err, "could not get block '%s'", blockID)
	}

	var block Block
	err = json.Unmarshal(resp, &block)
	if err != nil {
		return &block, errors.Wrapf(err, "could not get block '%s'", blockID)
	}

	return &block, nil
}

/*
OperationHashes is the hashes of all the operations included in the block.

Path:
	../<block_id>/operation_hashes (GET)
Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-context-contracts-contract-id-balance

Parameters:

	blockhash:
		The hash of block (height) of which you want to make the query.
*/
func (t *GoTezos) OperationHashes(blockhash string) ([][]string, error) {
	resp, err := t.get(fmt.Sprintf("/chains/main/blocks/%s/operation_hashes", blockhash))
	if err != nil {
		return [][]string{}, errors.Wrapf(err, "could not get operation hashes")
	}

	var operations [][]string
	err = json.Unmarshal(resp, &operations)
	if err != nil {
		return [][]string{}, errors.Wrapf(err, "could not unmarshal operation hashes")
	}

	return operations, nil
}

/*
BallotList returns ballots casted so far during a voting period.

Path:
	../<block_id>/votes/ballot_list (GET)
Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-votes-ballot-list

Parameters:

	blockhash:
		The hash of block (height) of which you want to make the query.
*/
func (t *GoTezos) BallotList(blockhash string) (BallotList, error) {
	resp, err := t.get(fmt.Sprintf("/chains/main/blocks/%s/votes/ballot_list", blockhash))
	if err != nil {
		return BallotList{}, errors.Wrapf(err, "failed to get ballot list")
	}

	var ballotList BallotList
	err = json.Unmarshal(resp, &ballotList)
	if err != nil {
		return BallotList{}, errors.Wrapf(err, "failed to unmarshal ballot list")
	}

	return ballotList, nil
}

/*
Ballots returns sum of ballots casted so far during a voting period.

Path:
	../<block_id>/votes/ballots (GET)
Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-votes-ballots

Parameters:

	blockhash:
		The hash of block (height) of which you want to make the query.
*/
func (t *GoTezos) Ballots(blockhash string) (Ballots, error) {
	resp, err := t.get(fmt.Sprintf("/chains/main/blocks/%s/votes/ballots", blockhash))
	if err != nil {
		return Ballots{}, errors.Wrapf(err, "failed to get ballots")
	}

	var ballots Ballots
	err = json.Unmarshal(resp, &ballots)
	if err != nil {
		return Ballots{}, errors.Wrapf(err, "failed to unmarshal ballots")
	}

	return ballots, nil
}

/*
CurrentPeriodKind returns the current period kind.

Path:
	../<block_id>/votes/current_period_kind (GET)
Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-votes-current-period-kind

Parameters:

	blockhash:
		The hash of block (height) of which you want to make the query.
*/
func (t *GoTezos) CurrentPeriodKind(blockhash string) (string, error) {
	resp, err := t.get(fmt.Sprintf("/chains/main/blocks/%s/votes/current_period_kind", blockhash))
	if err != nil {
		return "", errors.Wrapf(err, "failed to get current period kind")
	}

	var currentPeriodKind string
	err = json.Unmarshal(resp, &currentPeriodKind)
	if err != nil {
		return "", errors.Wrapf(err, "failed to unmarshal current period kind")
	}

	return currentPeriodKind, nil
}

/*
CurrentProposal returns the current proposal under evaluation.

Path:
	../<block_id>/votes/current_proposal (GET)
Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-votes-current-proposal

Parameters:

	blockhash:
		The hash of block (height) of which you want to make the query.
*/
func (t *GoTezos) CurrentProposal(blockhash string) (string, error) {
	resp, err := t.get(fmt.Sprintf("/chains/main/blocks/%s/votes/current_proposal", blockhash))
	if err != nil {
		return "", errors.Wrapf(err, "failed to get current proposal")
	}

	var currentProposal string
	err = json.Unmarshal(resp, &currentProposal)
	if err != nil {
		return "", errors.Wrapf(err, "failed to unmarshal current proposal")
	}

	return currentProposal, nil
}

/*
CurrentQuorum returns the current expected quorum.

Path:
	../<block_id>/votes/current_proposal (GET)
Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-votes-current-quorum

Parameters:

	blockhash:
		The hash of block (height) of which you want to make the query.
*/
func (t *GoTezos) CurrentQuorum(blockhash string) (int, error) {
	resp, err := t.get(fmt.Sprintf("/chains/main/blocks/%s/votes/current_quorum", blockhash))
	if err != nil {
		return 0, errors.Wrapf(err, "failed to get current quorum")
	}

	var currentQuorum int
	err = json.Unmarshal(resp, &currentQuorum)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to unmarshal current quorum")
	}

	return currentQuorum, nil
}

/*
VoteListings returns a list of delegates with their voting weight, in number of rolls.

Path:
	../<block_id>/votes/listings (GET)
Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-votes-listings

Parameters:

	blockhash:
		The hash of block (height) of which you want to make the query.
*/
func (t *GoTezos) VoteListings(blockhash string) (Listings, error) {
	resp, err := t.get(fmt.Sprintf("/chains/main/blocks/%s/votes/listings", blockhash))
	if err != nil {
		return Listings{}, errors.Wrapf(err, "failed to get listings")
	}

	var listings Listings
	err = json.Unmarshal(resp, &listings)
	if err != nil {
		return Listings{}, errors.Wrapf(err, "failed to unmarshal listings")
	}

	return listings, nil
}

/*
Proposals returns a list of proposals with number of supporters.

Path:
	../<block_id>/votes/proposals (GET)
Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-votes-proposals

Parameters:

	blockhash:
		The hash of block (height) of which you want to make the query.
*/
func (t *GoTezos) Proposals(blockhash string) (Proposals, error) {
	resp, err := t.get(fmt.Sprintf("/chains/main/blocks/%s/votes/proposals", blockhash))
	if err != nil {
		return Proposals{}, errors.Wrapf(err, "failed to get proposals")
	}

	var proposals Proposals
	err = json.Unmarshal(resp, &proposals)
	if err != nil {
		return Proposals{}, errors.Wrapf(err, "failed to unmarshal proposals")
	}

	return proposals, nil
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
