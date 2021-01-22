package rpc

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

// Kind is a contents kind
type Kind string

const (
	// ENDORSEMENT kind
	ENDORSEMENT Kind = "endorsement"
	// SEEDNONCEREVELATION kind
	SEEDNONCEREVELATION Kind = "seed_nonce_revelation"
	// DOUBLEENDORSEMENTEVIDENCE kind
	DOUBLEENDORSEMENTEVIDENCE Kind = "double_endorsement_evidence"
	// DOUBLEBAKINGEVIDENCE kind
	DOUBLEBAKINGEVIDENCE Kind = "Double_baking_evidence"
	// ACTIVATEACCOUNT kind
	ACTIVATEACCOUNT Kind = "activate_account"
	// PROPOSALS kind
	PROPOSALS Kind = "proposals"
	// BALLOT kind
	BALLOT Kind = "ballot"
	// REVEAL kind
	REVEAL Kind = "reveal"
	// TRANSACTION kind
	TRANSACTION Kind = "transaction"
	// ORIGINATION kind
	ORIGINATION Kind = "origination"
	// DELEGATION kind
	DELEGATION Kind = "delegation"
)

// BigMapDiffAction is an Action in a BigMapDiff
type BigMapDiffAction string

const (
	// UPDATE is a big_map_diff action
	UPDATE BigMapDiffAction = "update"
	// REMOVE is a big_map_diff action
	REMOVE BigMapDiffAction = "remove"
	// COPY is a big_map_diff action
	COPY BigMapDiffAction = "copy"
	// ALLOC is a big_map_diff action
	ALLOC BigMapDiffAction = "alloc"
)

/*
Block represents a Tezos block.

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id
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
	https://tezos.gitlab.io/008/rpc.html#get-block-id
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
	https://tezos.gitlab.io/008/rpc.html#get-block-id
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
	https://tezos.gitlab.io/008/rpc.html#get-block-id
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
	https://tezos.gitlab.io/008/rpc.html#get-block-id
*/
type MaxOperationListLength struct {
	MaxSize int `json:"max_size"`
	MaxOp   int `json:"max_op,omitempty"`
}

/*
Level represents the level in a Tezos block

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id
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
	https://tezos.gitlab.io/008/rpc.html#get-block-id
*/
type BalanceUpdates struct {
	Kind     string `json:"kind"`
	Contract string `json:"contract,omitempty"`
	Change   string `json:"change"`
	Category string `json:"category,omitempty"`
	Delegate string `json:"delegate,omitempty"`
	Cycle    int    `json:"cycle,omitempty"`
	Level    int    `json:"level,omitempty"`
}

// ResultError are errors reported by OperationResults
type ResultError struct {
	Kind           string           `json:"kind"`
	ID             string           `json:"id,omitempty"`
	With           *json.RawMessage `json:"with,omitempty"`
	Msg            string           `json:"msg,omitempty"`
	Location       int              `json:"location,omitempty"`
	ContractHandle string           `json:"contract_handle,omitempty"`
	ContractCode   *json.RawMessage `json:"contract_code,omitempty"`
}

/*
OperationResult represents the operation result in a Tezos block

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id
*/
type OperationResult struct {
	Status                       string           `json:"status"`
	Storage                      *json.RawMessage `json:"storage"`
	BigMapDiff                   BigMapDiffs      `json:"big_map_diff"`
	BalanceUpdates               []BalanceUpdates `json:"balance_updates"`
	OriginatedContracts          []string         `json:"originated_contracts"`
	ConsumedGas                  string           `json:"consumed_gas,omitempty"`
	StorageSize                  string           `json:"storage_size,omitempty"`
	AllocatedDestinationContract bool             `json:"allocated_destination_contract,omitempty"`
	Errors                       []ResultError    `json:"errors,omitempty"`
}

/*
Operations represents the operations in a Tezos block

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id
*/
type Operations struct {
	Protocol  string   `json:"protocol,omitempty"`
	ChainID   string   `json:"chain_id,omitempty"`
	Hash      string   `json:"hash,omitempty"`
	Branch    string   `json:"branch"`
	Contents  Contents `json:"contents"`
	Signature string   `json:"signature,omitempty"`
}

/*
OperationsAlt represents a JSON array containing an opHash at index 0, and
an Operation object at index 1. The RPC does not properly objectify Refused,
BranchRefused, BranchDelayed, and Unprocessed sections of the mempool,
so we must parse them manually.
Code hints used from github.com/blockwatch-cc/tzindex/rpc/mempool.go
*/
type OperationsAlt Operations

func (o *OperationsAlt) UnmarshalJSON(buf []byte) error {
	return unmarshalNamedJSONArray(buf, &o.Hash, (*Operations)(o))
}

func unmarshalNamedJSONArray(data []byte, v ...interface{}) error {
	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if len(raw) < len(v) {
		return fmt.Errorf("JSON array is too short, expected %d, got %d", len(v), len(raw))
	}

	for i, vv := range v {
		if err := json.Unmarshal(raw[i], vv); err != nil {
			return err
		}
	}

	return nil
}

/*
OrganizedContents represents the contents in Tezos operations orginized by kind.

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id
*/
type OrganizedContents struct {
	Endorsements              []Endorsement
	SeedNonceRevelations      []SeedNonceRevelation
	DoubleEndorsementEvidence []DoubleEndorsementEvidence
	DoubleBakingEvidence      []DoubleBakingEvidence
	AccountActivations        []AccountActivation
	Proposals                 []Proposal
	Ballots                   []Ballot
	Reveals                   []Reveal
	Transactions              []Transaction
	Originations              []Origination
	Delegations               []Delegation
}

// ToContents converts OrganizedContents into Contents
func (o *OrganizedContents) ToContents() Contents {
	var contents Contents
	for _, endorsement := range o.Endorsements {
		contents = append(contents, endorsement.ToContent())
	}

	for _, seedNonceRevelation := range o.SeedNonceRevelations {
		contents = append(contents, seedNonceRevelation.ToContent())
	}

	for _, doubleEndorsementEvidence := range o.DoubleEndorsementEvidence {
		contents = append(contents, doubleEndorsementEvidence.ToContent())
	}

	for _, doubleBakingEvidence := range o.DoubleBakingEvidence {
		contents = append(contents, doubleBakingEvidence.ToContent())
	}

	for _, accountActivation := range o.AccountActivations {
		contents = append(contents, accountActivation.ToContent())
	}

	for _, proposal := range o.Proposals {
		contents = append(contents, proposal.ToContent())
	}

	for _, ballot := range o.Ballots {
		contents = append(contents, ballot.ToContent())
	}

	for _, reveal := range o.AccountActivations {
		contents = append(contents, reveal.ToContent())
	}

	for _, transaction := range o.Transactions {
		contents = append(contents, transaction.ToContent())
	}

	for _, origination := range o.Originations {
		contents = append(contents, origination.ToContent())
	}

	for _, delegation := range o.Delegations {
		contents = append(contents, delegation.ToContent())
	}
	return contents
}

/*
Contents represents the contents in Tezos operations.

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id
*/
type Contents []Content

// Content is an element of Contents
type Content struct {
	Kind          Kind                `json:"kind,omitempty"`
	Level         int                 `json:"level,omitempty"`
	Nonce         string              `json:"nonce,omitempty"`
	Op1           *InlinedEndorsement `json:"Op1,omitempty"`
	Op2           *InlinedEndorsement `json:"Op2,omitempty"`
	Pkh           string              `json:"pkh,omitempty"`
	Secret        string              `json:"secret,omitempty"`
	Bh1           *BlockHeader        `json:"bh1,omitempty"`
	Bh2           *BlockHeader        `json:"bh2,omitempty"`
	Source        string              `json:"source,omitempty"`
	Period        int                 `json:"period,omitempty"`
	Proposals     []string            `json:"proposals,omitempty"`
	Proposal      string              `json:"proposal,omitempty"`
	Ballot        string              `json:"ballot,omitempty"`
	Fee           string              `json:"fee,omitempty"`
	Counter       string              `json:"counter,omitempty"`
	GasLimit      string              `json:"gas_limit,omitempty"`
	StorageLimit  string              `json:"storage_limit,omitempty"`
	PublicKey     string              `json:"public_key,omitempty"`
	ManagerPubkey string              `json:"managerPubKey,omitempty"`
	Amount        string              `json:"amount,omitempty"`
	Destination   string              `json:"destination,omitempty"`
	Balance       string              `json:"balance,omitempty"`
	Delegate      string              `json:"delegate,omitempty"`
	Script        Script              `json:"script,omitempty"`
	Parameters    *Parameters         `json:"parameters,omitempty"`
	Metadata      *ContentsMetadata   `json:"metadata,omitempty"`
}

// MarshalJSON implements json.Marshaler in order to correctly marshal contents based of kind
func (c *Content) MarshalJSON() ([]byte, error) {
	if c.Kind == ENDORSEMENT {
		return json.Marshal(c.ToEndorsement())
	} else if c.Kind == SEEDNONCEREVELATION {
		return json.Marshal(c.ToSeedNonceRevelations())
	} else if c.Kind == DOUBLEENDORSEMENTEVIDENCE {
		return json.Marshal(c.ToDoubleEndorsementEvidence())
	} else if c.Kind == DOUBLEBAKINGEVIDENCE {
		return json.Marshal(c.ToDoubleBakingEvidence())
	} else if c.Kind == ACTIVATEACCOUNT {
		return json.Marshal(c.ToAccountActivation())
	} else if c.Kind == PROPOSALS {
		return json.Marshal(c.ToProposal())
	} else if c.Kind == BALLOT {
		return json.Marshal(c.ToBallot())
	} else if c.Kind == REVEAL {
		return json.Marshal(c.ToReveal())
	} else if c.Kind == TRANSACTION {
		return json.Marshal(c.ToTransaction())
	} else if c.Kind == ORIGINATION {
		return json.Marshal(c.ToOrigination())
	} else if c.Kind == DELEGATION {
		return json.Marshal(c.ToDelegation())
	}

	return nil, errors.New("failed to find content kind to marshal into")
}

// Organize converts contents into OrganizedContents where contents are organized by Kind
func (c Contents) Organize() OrganizedContents {
	var organizeContents OrganizedContents
	for _, content := range c {
		if content.Kind == ENDORSEMENT {
			organizeContents.Endorsements = append(organizeContents.Endorsements, content.ToEndorsement())
		} else if content.Kind == SEEDNONCEREVELATION {
			organizeContents.SeedNonceRevelations = append(organizeContents.SeedNonceRevelations, content.ToSeedNonceRevelations())
		} else if content.Kind == DOUBLEENDORSEMENTEVIDENCE {
			organizeContents.DoubleEndorsementEvidence = append(organizeContents.DoubleEndorsementEvidence, content.ToDoubleEndorsementEvidence())
		} else if content.Kind == DOUBLEBAKINGEVIDENCE {
			organizeContents.DoubleBakingEvidence = append(organizeContents.DoubleBakingEvidence, content.ToDoubleBakingEvidence())
		} else if content.Kind == ACTIVATEACCOUNT {
			organizeContents.AccountActivations = append(organizeContents.AccountActivations, content.ToAccountActivation())
		} else if content.Kind == PROPOSALS {
			organizeContents.Proposals = append(organizeContents.Proposals, content.ToProposal())
		} else if content.Kind == BALLOT {
			organizeContents.Ballots = append(organizeContents.Ballots, content.ToBallot())
		} else if content.Kind == REVEAL {
			organizeContents.Reveals = append(organizeContents.Reveals, content.ToReveal())
		} else if content.Kind == TRANSACTION {
			organizeContents.Transactions = append(organizeContents.Transactions, content.ToTransaction())
		} else if content.Kind == ORIGINATION {
			organizeContents.Originations = append(organizeContents.Originations, content.ToOrigination())
		} else if content.Kind == DELEGATION {
			organizeContents.Delegations = append(organizeContents.Delegations, content.ToDelegation())
		}
	}

	return organizeContents
}

/*
Parameters represents parameters in Tezos operations.

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id
*/
type Parameters struct {
	Entrypoint string           `json:"entrypoint"`
	Value      *json.RawMessage `json:"value"`
}

/*
ContentsMetadata represents metadata in contents in Tezos operations.

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id
*/
type ContentsMetadata struct {
	BalanceUpdates          []BalanceUpdates           `json:"balance_updates,omitempty"`
	Delegate                string                     `json:"delegate,omitempty"`
	Slots                   []int                      `json:"slots,omitempty"`
	OperationResults        *OperationResults          `json:"operation_result,omitempty"`
	InternalOperationResult []InternalOperationResults `json:"internal_operation_results,omitempty"`
}

/*
OperationResults represents the operation_results in Tezos operations.

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id
*/
type OperationResults struct {
	Status                       string           `json:"status"`
	BigMapDiff                   BigMapDiffs      `json:"big_map_diff,omitempty"`
	BalanceUpdates               []BalanceUpdates `json:"balance_updates,omitempty"`
	OriginatedContracts          []string         `json:"originated_contracts,omitempty"`
	ConsumedGas                  string           `json:"consumed_gas,omitempty"`
	StorageSize                  string           `json:"storage_size,omitempty"`
	PaidStorageSizeDiff          string           `json:"paid_storage_size_diff,omitempty"`
	Errors                       []ResultError    `json:"errors,omitempty"`
	Storage                      *json.RawMessage `json:"storage,omitempty"`
	AllocatedDestinationContract bool             `json:"allocated_destination_contract,omitempty"`
}

func (o *OperationResults) toOperationResultsReveal() OperationResultReveal {
	return OperationResultReveal{
		Status:      o.Status,
		ConsumedGas: o.ConsumedGas,
		Errors:      o.Errors,
	}
}

func (o *OperationResults) toOperationResultsTransfer() OperationResultTransfer {
	var (
		storage    *json.RawMessage
		bigMapDiff BigMapDiffs
	)

	if o.Storage != nil {
		storage = o.Storage
	}

	if o.BigMapDiff != nil {
		bigMapDiff = o.BigMapDiff
	}

	return OperationResultTransfer{
		Status:                       o.Status,
		Storage:                      storage,
		BigMapDiff:                   bigMapDiff,
		BalanceUpdates:               o.BalanceUpdates,
		OriginatedContracts:          o.OriginatedContracts,
		ConsumedGas:                  o.ConsumedGas,
		StorageSize:                  o.StorageSize,
		PaidStorageSizeDiff:          o.PaidStorageSizeDiff,
		AllocatedDestinationContract: o.AllocatedDestinationContract,
		Errors:                       o.Errors,
	}
}

func (o *OperationResults) toOperationResultsOrigination() OperationResultOrigination {
	var bigMapDiff BigMapDiffs

	if o.BigMapDiff != nil {
		bigMapDiff = o.BigMapDiff
	}

	return OperationResultOrigination{
		Status:              o.Status,
		BigMapDiff:          bigMapDiff,
		BalanceUpdates:      o.BalanceUpdates,
		OriginatedContracts: o.OriginatedContracts,
		ConsumedGas:         o.ConsumedGas,
		StorageSize:         o.StorageSize,
		PaidStorageSizeDiff: o.PaidStorageSizeDiff,
		Errors:              o.Errors,
	}
}

func (o *OperationResults) toOperationResultsDelegation() OperationResultDelegation {
	return OperationResultDelegation{
		Status:      o.Status,
		ConsumedGas: o.ConsumedGas,
		Errors:      o.Errors,
	}
}

// ToEndorsement converts Content to Endorsement.
func (c *Content) ToEndorsement() Endorsement {
	var metadata *EndorsementMetadata

	if c.Metadata != nil {
		metadata = &EndorsementMetadata{
			BalanceUpdates: c.Metadata.BalanceUpdates,
			Delegate:       c.Metadata.Delegate,
			Slots:          c.Metadata.Slots,
		}
	}

	return Endorsement{
		Kind:     c.Kind,
		Level:    c.Level,
		Metadata: metadata,
	}
}

// ToSeedNonceRevelations converts Content to SeedNonceRevelations.
func (c *Content) ToSeedNonceRevelations() SeedNonceRevelation {
	var metadata *SeedNonceRevelationMetadata

	if c.Metadata != nil {
		metadata = &SeedNonceRevelationMetadata{
			BalanceUpdates: c.Metadata.BalanceUpdates,
		}
	}

	return SeedNonceRevelation{
		Kind:     c.Kind,
		Level:    c.Level,
		Nonce:    c.Nonce,
		Metadata: metadata,
	}
}

// ToDoubleEndorsementEvidence converts Content to DoubleEndorsementEvidence.
func (c *Content) ToDoubleEndorsementEvidence() DoubleEndorsementEvidence {
	var (
		metadata *DoubleEndorsementEvidenceMetadata
		op1      *InlinedEndorsement
		op2      *InlinedEndorsement
	)

	if c.Op1 != nil {
		op1 = c.Op1
	}

	if c.Op2 != nil {
		op2 = c.Op2
	}

	if c.Metadata != nil {
		metadata = &DoubleEndorsementEvidenceMetadata{
			BalanceUpdates: c.Metadata.BalanceUpdates,
		}
	}

	return DoubleEndorsementEvidence{
		Kind:     c.Kind,
		Op1:      op1,
		Op2:      op2,
		Metadata: metadata,
	}
}

// ToDoubleBakingEvidence converts Content to DoubleBakingEvidence.
func (c *Content) ToDoubleBakingEvidence() DoubleBakingEvidence {
	var (
		metadata *DoubleBakingEvidenceMetadata
		bh1      *BlockHeader
		bh2      *BlockHeader
	)

	if c.Bh1 != nil {
		bh1 = c.Bh1
	}

	if c.Bh2 != nil {
		bh2 = c.Bh2
	}

	if c.Metadata != nil {
		metadata = &DoubleBakingEvidenceMetadata{
			BalanceUpdates: c.Metadata.BalanceUpdates,
		}
	}

	return DoubleBakingEvidence{
		Kind:     c.Kind,
		Bh1:      bh1,
		Bh2:      bh2,
		Metadata: metadata,
	}
}

// ToAccountActivation converts Content to AccountActivation.
func (c *Content) ToAccountActivation() AccountActivation {
	var metadata *AccountActivationMetadata

	if c.Metadata != nil {
		metadata = &AccountActivationMetadata{
			BalanceUpdates: c.Metadata.BalanceUpdates,
		}
	}

	return AccountActivation{
		Kind:     c.Kind,
		Pkh:      c.Pkh,
		Secret:   c.Secret,
		Metadata: metadata,
	}
}

// ToProposal converts Content to Proposal.
func (c *Content) ToProposal() Proposal {
	return Proposal{
		Kind:      c.Kind,
		Source:    c.Source,
		Period:    c.Period,
		Proposals: c.Proposals,
	}
}

// ToBallot converts Content to Ballot.
func (c *Content) ToBallot() Ballot {
	return Ballot{
		Kind:     c.Kind,
		Source:   c.Source,
		Period:   c.Period,
		Proposal: c.Proposal,
		Ballot:   c.Ballot,
	}
}

// ToReveal converts Content to Reveal.
func (c *Content) ToReveal() Reveal {
	var metadata *RevealMetadata

	if c.Metadata != nil {
		metadata = &RevealMetadata{
			BalanceUpdates:           c.Metadata.BalanceUpdates,
			OperationResult:          c.Metadata.OperationResults.toOperationResultsReveal(),
			InternalOperationResults: c.Metadata.InternalOperationResult,
		}
	}

	return Reveal{
		Kind:         c.Kind,
		Source:       c.Source,
		Fee:          c.Fee,
		Counter:      c.Counter,
		GasLimit:     c.GasLimit,
		StorageLimit: c.StorageLimit,
		PublicKey:    c.PublicKey,
		Metadata:     metadata,
	}
}

// ToTransaction converts Content to Transaction.
func (c *Content) ToTransaction() Transaction {
	var (
		metadata   *TransactionMetadata
		parameters *Parameters
	)

	if c.Metadata != nil {
		metadata = &TransactionMetadata{
			BalanceUpdates:           c.Metadata.BalanceUpdates,
			OperationResult:          c.Metadata.OperationResults.toOperationResultsTransfer(),
			InternalOperationResults: c.Metadata.InternalOperationResult,
		}
	}

	if c.Parameters != nil {
		parameters = &Parameters{
			Entrypoint: c.Parameters.Entrypoint,
			Value:      c.Parameters.Value,
		}
	}

	return Transaction{
		Kind:         c.Kind,
		Source:       c.Source,
		Fee:          c.Fee,
		Counter:      c.Counter,
		GasLimit:     c.GasLimit,
		StorageLimit: c.StorageLimit,
		Amount:       c.Amount,
		Destination:  c.Destination,
		Parameters:   parameters,
		Metadata:     metadata,
	}
}

// ToOrigination converts Content to Origination.
func (c *Content) ToOrigination() Origination {
	var metadata *OriginationMetadata

	if c.Metadata != nil {
		metadata = &OriginationMetadata{
			BalanceUpdates:           c.Metadata.BalanceUpdates,
			OperationResults:         c.Metadata.OperationResults.toOperationResultsOrigination(),
			InternalOperationResults: c.Metadata.InternalOperationResult,
		}
	}

	return Origination{
		Kind:          c.Kind,
		Source:        c.Source,
		Fee:           c.Fee,
		Counter:       c.Counter,
		GasLimit:      c.GasLimit,
		StorageLimit:  c.StorageLimit,
		Balance:       c.Balance,
		Delegate:      c.Delegate,
		Script:        c.Script,
		ManagerPubkey: c.ManagerPubkey,
		Metadata:      metadata,
	}
}

// ToDelegation converts Content to Origination.
func (c *Content) ToDelegation() Delegation {
	var metadata *DelegationMetadata

	if c.Metadata != nil {
		metadata = &DelegationMetadata{
			BalanceUpdates:           c.Metadata.BalanceUpdates,
			OperationResults:         c.Metadata.OperationResults.toOperationResultsDelegation(),
			InternalOperationResults: c.Metadata.InternalOperationResult,
		}
	}

	return Delegation{
		Kind:         c.Kind,
		Source:       c.Source,
		Fee:          c.Fee,
		Counter:      c.Counter,
		GasLimit:     c.GasLimit,
		StorageLimit: c.StorageLimit,
		Delegate:     c.Delegate,
		Metadata:     metadata,
	}
}

//MarshalJSON satisfies the json.MarshalJSON interface for contents
func (o *OrganizedContents) MarshalJSON() ([]byte, error) {
	contents := o.ToContents()
	return json.Marshal(&contents)
}

/*
Endorsement represents an endorsement in the $operation.alpha.operation_contents_and_result in the tezos block schema

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id
*/
type Endorsement struct {
	Kind     Kind                 `json:"kind"`
	Level    int                  `json:"level"`
	Metadata *EndorsementMetadata `json:"metadata,omitempty"`
}

/*
EndorsementMetadata represents the metadata of an endorsement in the $operation.alpha.operation_contents_and_result in the tezos block schema

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id
*/
type EndorsementMetadata struct {
	BalanceUpdates []BalanceUpdates `json:"balance_updates"`
	Delegate       string           `json:"delegate"`
	Slots          []int            `json:"slots"`
}

// ToContent converts Endorsement to Content
func (e *Endorsement) ToContent() Content {
	var metadata *ContentsMetadata

	if e.Metadata != nil {
		metadata = &ContentsMetadata{
			BalanceUpdates: e.Metadata.BalanceUpdates,
			Delegate:       e.Metadata.Delegate,
			Slots:          e.Metadata.Slots,
		}
	}

	return Content{
		Kind:     e.Kind,
		Level:    e.Level,
		Metadata: metadata,
	}
}

/*
SeedNonceRevelation represents an Seed_nonce_revelation in the $operation.alpha.operation_contents_and_result in the tezos block schema

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id
*/
type SeedNonceRevelation struct {
	Kind     Kind                         `json:"kind"`
	Level    int                          `json:"level"`
	Nonce    string                       `json:"nonce"`
	Metadata *SeedNonceRevelationMetadata `json:"metadata,omitempty"`
}

// ToContent converts a SeedNonceRevelation to Content
func (s *SeedNonceRevelation) ToContent() Content {
	var metadata *ContentsMetadata

	if s.Metadata != nil {
		metadata = &ContentsMetadata{
			BalanceUpdates: s.Metadata.BalanceUpdates,
		}
	}

	return Content{
		Kind:     s.Kind,
		Level:    s.Level,
		Nonce:    s.Nonce,
		Metadata: metadata,
	}
}

/*
SeedNonceRevelationMetadata represents the metadata for Seed_nonce_revelation in the $operation.alpha.operation_contents_and_result in the tezos block schema

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id
*/
type SeedNonceRevelationMetadata struct {
	BalanceUpdates []BalanceUpdates `json:"balance_updates"`
}

/*
DoubleEndorsementEvidence represents an Double_endorsement_evidence in the $operation.alpha.operation_contents_and_result in the tezos block schema

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id
*/
type DoubleEndorsementEvidence struct {
	Kind     Kind                               `json:"kind"`
	Op1      *InlinedEndorsement                `json:"Op1"`
	Op2      *InlinedEndorsement                `json:"Op2"`
	Metadata *DoubleEndorsementEvidenceMetadata `json:"metadata,omitempty"`
}

/*
DoubleEndorsementEvidenceMetadata represents the metadata for Double_endorsement_evidence in the $operation.alpha.operation_contents_and_result in the tezos block schema

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id
*/
type DoubleEndorsementEvidenceMetadata struct {
	BalanceUpdates []BalanceUpdates `json:"balance_updates"`
}

/*
InlinedEndorsement represents $inlined.endorsement in the tezos block schema

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id
*/
type InlinedEndorsement struct {
	Branch     string                        `json:"branch"`
	Operations *InlinedEndorsementOperations `json:"operations"`
	Signature  string                        `json:"signature"`
}

/*
InlinedEndorsementOperations represents operations in $inlined.endorsement in the tezos block schema

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id
*/
type InlinedEndorsementOperations struct {
	Kind  string `json:"kind"`
	Level int    `json:"level"`
}

// ToContent converts a DoubleEndorsementEvidence to Content
func (d *DoubleEndorsementEvidence) ToContent() Content {
	var (
		metadata *ContentsMetadata
		op1      *InlinedEndorsement
		op2      *InlinedEndorsement
	)

	if d.Op1 != nil {
		op1 = d.Op1
	}

	if d.Op2 != nil {
		op2 = d.Op2
	}

	if d.Metadata != nil {
		metadata = &ContentsMetadata{
			BalanceUpdates: d.Metadata.BalanceUpdates,
		}
	}

	return Content{
		Kind:     d.Kind,
		Op1:      op1,
		Op2:      op2,
		Metadata: metadata,
	}
}

/*
DoubleBakingEvidence represents an Double_baking_evidence in the $operation.alpha.operation_contents_and_result in the tezos block schema

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id
*/
type DoubleBakingEvidence struct {
	Kind     Kind                          `json:"kind"`
	Bh1      *BlockHeader                  `json:"bh1"`
	Bh2      *BlockHeader                  `json:"bh2"`
	Metadata *DoubleBakingEvidenceMetadata `json:"metadata,omitempty"`
}

/*
DoubleBakingEvidenceMetadata represents the metadata of Double_baking_evidence in the $operation.alpha.operation_contents_and_result in the tezos block schema

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id
*/
type DoubleBakingEvidenceMetadata struct {
	BalanceUpdates []BalanceUpdates `json:"balance_updates"`
}

/*
BlockHeader represents $block_header.alpha.full_header in the tezos block schema

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id
*/
type BlockHeader struct {
	Level            int       `json:"level"`
	Proto            int       `json:"proto"`
	Predecessor      string    `json:"predecessor"`
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

// ToContent converts a DoubleBakingEvidence to Content
func (d *DoubleBakingEvidence) ToContent() Content {
	var (
		metadata *ContentsMetadata
		bh1      *BlockHeader
		bh2      *BlockHeader
	)

	if d.Bh1 != nil {
		bh1 = d.Bh1
	}

	if d.Bh2 != nil {
		bh2 = d.Bh2
	}

	if d.Metadata != nil {
		metadata = &ContentsMetadata{
			BalanceUpdates: d.Metadata.BalanceUpdates,
		}
	}

	return Content{
		Kind:     d.Kind,
		Bh1:      bh1,
		Bh2:      bh2,
		Metadata: metadata,
	}
}

/*
AccountActivation represents an Activate_account in the $operation.alpha.operation_contents_and_result in the tezos block schema

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id
*/
type AccountActivation struct {
	Kind     Kind                       `json:"kind"`
	Pkh      string                     `json:"pkh"`
	Secret   string                     `json:"secret"`
	Metadata *AccountActivationMetadata `json:"metadata,omitempty"`
}

/*
AccountActivationMetadata represents the metadata for Activate_account in the $operation.alpha.operation_contents_and_result in the tezos block schema

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id
*/
type AccountActivationMetadata struct {
	BalanceUpdates []BalanceUpdates `json:"balance_updates"`
}

// ToContent converts a AccountActivation to Content
func (a *AccountActivation) ToContent() Content {
	var metadata *ContentsMetadata
	if a.Metadata != nil {
		metadata = &ContentsMetadata{
			BalanceUpdates: a.Metadata.BalanceUpdates,
		}
	}

	return Content{
		Kind:     a.Kind,
		Pkh:      a.Pkh,
		Secret:   a.Secret,
		Metadata: metadata,
	}
}

/*
Proposal represents a Proposal in the $operation.alpha.operation_contents_and_result in the tezos block schema

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id
*/
type Proposal struct {
	Kind      Kind     `json:"kind"`
	Source    string   `json:"source"`
	Period    int      `json:"period"`
	Proposals []string `json:"proposals"`
}

// ToContent converts a Proposal to Content
func (p *Proposal) ToContent() Content {
	return Content{
		Kind:      p.Kind,
		Source:    p.Source,
		Period:    p.Period,
		Proposals: p.Proposals,
	}
}

/*
Ballot represents a Ballot in the $operation.alpha.operation_contents_and_result in the tezos block schema

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id
*/
type Ballot struct {
	Kind     Kind   `json:"kind"`
	Source   string `json:"source"`
	Period   int    `json:"period"`
	Proposal string `json:"proposal"`
	Ballot   string `json:"ballot"`
}

// ToContent converts a Ballot to Content
func (b *Ballot) ToContent() Content {
	return Content{
		Kind:     b.Kind,
		Source:   b.Source,
		Period:   b.Period,
		Proposal: b.Proposal,
		Ballot:   b.Ballot,
	}
}

/*
Reveal represents a Reveal in the $operation.alpha.operation_contents_and_result in the tezos block schema

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id
*/
type Reveal struct {
	Kind         Kind            `json:"kind"`
	Source       string          `json:"source" validate:"required"`
	Fee          string          `json:"fee" validate:"required"`
	Counter      string          `json:"counter" validate:"required"`
	GasLimit     string          `json:"gas_limit" validate:"required"`
	StorageLimit string          `json:"storage_limit"`
	PublicKey    string          `json:"public_key" validate:"required"`
	Metadata     *RevealMetadata `json:"metadata,omitempty"`
}

/*
RevealMetadata represents the metadata for Reveal in the $operation.alpha.operation_contents_and_result in the tezos block schema

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id
*/
type RevealMetadata struct {
	BalanceUpdates           []BalanceUpdates           `json:"balance_updates"`
	OperationResult          OperationResultReveal      `json:"operation_result"`
	InternalOperationResults []InternalOperationResults `json:"internal_operation_result,omitempty"`
}

// ToContent converts a Reveal to Content
func (r *Reveal) ToContent() Content {
	var metadata *ContentsMetadata

	if r.Metadata != nil {
		metadata = &ContentsMetadata{
			BalanceUpdates: r.Metadata.BalanceUpdates,
			OperationResults: &OperationResults{
				Status:      r.Metadata.OperationResult.Status,
				ConsumedGas: r.Metadata.OperationResult.ConsumedGas,
				Errors:      r.Metadata.OperationResult.Errors,
			},
			InternalOperationResult: r.Metadata.InternalOperationResults,
		}
	}

	return Content{
		Kind:         r.Kind,
		Source:       r.Source,
		Fee:          r.Fee,
		Counter:      r.Counter,
		GasLimit:     r.GasLimit,
		StorageLimit: r.StorageLimit,
		PublicKey:    r.PublicKey,
		Metadata:     metadata,
	}
}

/*
Transaction represents a Transaction in the $operation.alpha.operation_contents_and_result in the tezos block schema

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id
*/
type Transaction struct {
	Kind         Kind                 `json:"kind"`
	Source       string               `json:"source" validate:"required"`
	Fee          string               `json:"fee" validate:"required"`
	Counter      string               `json:"counter" validate:"required"`
	GasLimit     string               `json:"gas_limit" validate:"required"`
	StorageLimit string               `json:"storage_limit"`
	Amount       string               `json:"amount"`
	Destination  string               `json:"destination" validate:"required"`
	Parameters   *Parameters          `json:"parameters,omitempty"`
	Metadata     *TransactionMetadata `json:"metadata,omitempty"`
}

/*
TransactionMetadata represents the metadata of Transaction in the $operation.alpha.operation_contents_and_result in the tezos block schema

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id
*/
type TransactionMetadata struct {
	BalanceUpdates           []BalanceUpdates           `json:"balance_updates"`
	OperationResult          OperationResultTransfer    `json:"operation_result"`
	InternalOperationResults []InternalOperationResults `json:"internal_operation_results,omitempty"`
}

// ToContent converts a Transaction to Content
func (t *Transaction) ToContent() Content {
	var (
		parameters *Parameters
		metadata   *ContentsMetadata
	)

	if t.Parameters != nil {
		parameters = &Parameters{
			Entrypoint: t.Parameters.Entrypoint,
			Value:      t.Parameters.Value,
		}
	}

	if t.Metadata != nil {
		metadata = &ContentsMetadata{
			BalanceUpdates: t.Metadata.BalanceUpdates,
			OperationResults: &OperationResults{
				Status:                       t.Metadata.OperationResult.Status,
				Storage:                      t.Metadata.OperationResult.Storage,
				BigMapDiff:                   t.Metadata.OperationResult.BigMapDiff,
				BalanceUpdates:               t.Metadata.OperationResult.BalanceUpdates,
				OriginatedContracts:          t.Metadata.OperationResult.OriginatedContracts,
				ConsumedGas:                  t.Metadata.OperationResult.ConsumedGas,
				StorageSize:                  t.Metadata.OperationResult.StorageSize,
				PaidStorageSizeDiff:          t.Metadata.OperationResult.PaidStorageSizeDiff,
				AllocatedDestinationContract: t.Metadata.OperationResult.AllocatedDestinationContract,
				Errors:                       t.Metadata.OperationResult.Errors,
			},
			InternalOperationResult: t.Metadata.InternalOperationResults,
		}
	}

	return Content{
		Kind:         t.Kind,
		Source:       t.Source,
		Fee:          t.Fee,
		Counter:      t.Counter,
		GasLimit:     t.GasLimit,
		StorageLimit: t.StorageLimit,
		Amount:       t.Amount,
		Destination:  t.Destination,
		Parameters:   parameters,
		Metadata:     metadata,
	}
}

/*
Origination represents a Origination in the $operation.alpha.operation_contents_and_result in the tezos block schema

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id
*/
type Origination struct {
	Kind          Kind                 `json:"kind"`
	Source        string               `json:"source" validate:"required"`
	Fee           string               `json:"fee" validate:"required"`
	Counter       string               `json:"counter" validate:"required"`
	GasLimit      string               `json:"gas_limit" validate:"required"`
	StorageLimit  string               `json:"storage_limit" validate:"required"`
	Balance       string               `json:"balance"`
	Delegate      string               `json:"delegate,omitempty"`
	Script        Script               `json:"script" validate:"required"`
	ManagerPubkey string               `json:"managerPubkey,omitempty"`
	Metadata      *OriginationMetadata `json:"metadata,omitempty"`
}

/*
Script represents the script in an Origination in the $operation.alpha.operation_contents_and_result -> $scripted.contracts in the tezos block schema

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id
*/
type Script struct {
	Code    *json.RawMessage `json:"code,omitempty"`
	Storage *json.RawMessage `json:"storage,omitempty"`
}

/*
OriginationMetadata represents the metadata of Origination in the $operation.alpha.operation_contents_and_result in the tezos block schema

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id
*/
type OriginationMetadata struct {
	BalanceUpdates           []BalanceUpdates           `json:"balance_updates"`
	OperationResults         OperationResultOrigination `json:"operation_result"`
	InternalOperationResults []InternalOperationResults `json:"internal_operation_results,omitempty"`
}

// ToContent converts a Origination to Content
func (o *Origination) ToContent() Content {
	var metadata *ContentsMetadata

	if o.Metadata != nil {
		metadata = &ContentsMetadata{
			BalanceUpdates: o.Metadata.BalanceUpdates,
			OperationResults: &OperationResults{
				Status:              o.Metadata.OperationResults.Status,
				BigMapDiff:          o.Metadata.OperationResults.BigMapDiff,
				BalanceUpdates:      o.Metadata.OperationResults.BalanceUpdates,
				OriginatedContracts: o.Metadata.OperationResults.OriginatedContracts,
				ConsumedGas:         o.Metadata.OperationResults.ConsumedGas,
				StorageSize:         o.Metadata.OperationResults.StorageSize,
				PaidStorageSizeDiff: o.Metadata.OperationResults.PaidStorageSizeDiff,
				Errors:              o.Metadata.OperationResults.Errors,
			},
			InternalOperationResult: o.Metadata.InternalOperationResults,
		}
	}

	return Content{
		Kind:          o.Kind,
		Source:        o.Source,
		Fee:           o.Fee,
		Counter:       o.Counter,
		GasLimit:      o.GasLimit,
		StorageLimit:  o.StorageLimit,
		Balance:       o.Balance,
		Delegate:      o.Delegate,
		Script:        o.Script,
		ManagerPubkey: o.ManagerPubkey,
		Metadata:      metadata,
	}
}

/*
Delegation represents a Delegation in the $operation.alpha.operation_contents_and_result in the tezos block schema

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id
*/
type Delegation struct {
	Kind         Kind                `json:"kind"`
	Source       string              `json:"source" validate:"required"`
	Fee          string              `json:"fee" validate:"required"`
	Counter      string              `json:"counter" validate:"required"`
	GasLimit     string              `json:"gas_limit" validate:"required"`
	StorageLimit string              `json:"storage_limit" validate:"required"`
	Delegate     string              `json:"delegate,omitempty"`
	Metadata     *DelegationMetadata `json:"metadata,omitempty"`
}

/*
DelegationMetadata represents the metadata Delegation in the $operation.alpha.operation_contents_and_result in the tezos block schema

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id
*/
type DelegationMetadata struct {
	BalanceUpdates           []BalanceUpdates           `json:"balance_updates"`
	OperationResults         OperationResultDelegation  `json:"operation_result"`
	InternalOperationResults []InternalOperationResults `json:"internal_operation_results,omitempty"`
}

// ToContent converts a Delegation to Content
func (d *Delegation) ToContent() Content {
	var metadata *ContentsMetadata
	if d.Metadata != nil {
		metadata = &ContentsMetadata{
			BalanceUpdates: d.Metadata.BalanceUpdates,
			OperationResults: &OperationResults{
				Status:      d.Metadata.OperationResults.Status,
				ConsumedGas: d.Metadata.OperationResults.ConsumedGas,
				Errors:      d.Metadata.OperationResults.Errors,
			},
			InternalOperationResult: d.Metadata.InternalOperationResults,
		}
	}

	return Content{
		Kind:         d.Kind,
		Source:       d.Source,
		Fee:          d.Fee,
		Counter:      d.Counter,
		GasLimit:     d.GasLimit,
		StorageLimit: d.StorageLimit,
		Delegate:     d.Delegate,
		Metadata:     metadata,
	}
}

/*
InternalOperationResults represents an InternalOperationResults in the $operation.alpha.internal_operation_result in the tezos block schema

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id
*/
type InternalOperationResults struct {
	Kind        string            `json:"kind"`
	Source      string            `json:"source"`
	Nonce       int               `json:"nonce"`
	Amount      string            `json:"amount,omitempty"`
	PublicKey   string            `json:"public_key,omitempty"`
	Destination string            `json:"destination,omitempty"`
	Balance     string            `json:"balance,omitempty"`
	Delegate    string            `json:"delegate,omitempty"`
	Script      ScriptedContracts `json:"script,omitempty"`
	Parameters  struct {
		Entrypoint string           `json:"entrypoint"`
		Value      *json.RawMessage `json:"value"`
	} `json:"paramaters,omitempty"`
	Result OperationResult `json:"result"`
}

/*
OperationResultReveal represents an OperationResultReveal in the $operation.alpha.operation_result.reveal in the tezos block schema

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id
*/
type OperationResultReveal struct {
	Status      string        `json:"status"`
	ConsumedGas string        `json:"consumed_gas,omitempty"`
	Errors      []ResultError `json:"rpc_error,omitempty"`
}

/*
OperationResultTransfer represents $operation.alpha.operation_result.transaction in the tezos block schema

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id
*/
type OperationResultTransfer struct {
	Status                       string           `json:"status"`
	Storage                      *json.RawMessage `json:"storage,omitempty"`
	BigMapDiff                   BigMapDiffs      `json:"big_map_diff,omitempty"`
	BalanceUpdates               []BalanceUpdates `json:"balance_updates,omitempty"`
	OriginatedContracts          []string         `json:"originated_contracts,omitempty"`
	ConsumedGas                  string           `json:"consumed_gas,omitempty"`
	StorageSize                  string           `json:"storage_size,omitempty"`
	PaidStorageSizeDiff          string           `json:"paid_storage_size_diff,omitempty"`
	AllocatedDestinationContract bool             `json:"allocated_destination_contract,omitempty"`
	Errors                       []ResultError    `json:"errors,omitempty"`
}

/*
OperationResultOrigination represents $operation.alpha.operation_result.origination in the tezos block schema

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id
*/
type OperationResultOrigination struct {
	Status              string           `json:"status"`
	BigMapDiff          BigMapDiffs      `json:"big_map_diff,omitempty"`
	BalanceUpdates      []BalanceUpdates `json:"balance_updates,omitempty"`
	OriginatedContracts []string         `json:"originated_contracts,omitempty"`
	ConsumedGas         string           `json:"consumed_gas,omitempty"`
	StorageSize         string           `json:"storage_size,omitempty"`
	PaidStorageSizeDiff string           `json:"paid_storage_size_diff,omitempty"`
	Errors              []ResultError    `json:"errors,omitempty"`
}

/*
OperationResultDelegation represents $operation.alpha.operation_result.delegation in the tezos block schema

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id
*/
type OperationResultDelegation struct {
	Status      string        `json:"status"`
	ConsumedGas string        `json:"consumed_gas,omitempty"`
	Errors      []ResultError `json:"errors,omitempty"`
}

// OrganizedBigMapDiff represents a BigMapDiffs organized by kind.
type OrganizedBigMapDiff struct {
	Updates  []BigMapDiffUpdate
	Removals []BigMapDiffRemove
	Copies   []BigMapDiffCopy
	Allocs   []BigMapDiffAlloc
}

// ToBigMapDiffs converts OrganizedBigMapDiff to BigMapDiffs
func (o *OrganizedBigMapDiff) ToBigMapDiffs() BigMapDiffs {
	var bigMapDiffs BigMapDiffs
	for _, update := range o.Updates {
		bigMapDiffs = append(bigMapDiffs, update.toBigMapDiff())
	}

	for _, removal := range o.Removals {
		bigMapDiffs = append(bigMapDiffs, removal.toBigMapDiff())
	}

	for _, copy := range o.Copies {
		bigMapDiffs = append(bigMapDiffs, copy.toBigMapDiff())
	}

	for _, alloc := range o.Allocs {
		bigMapDiffs = append(bigMapDiffs, alloc.toBigMapDiff())
	}

	return bigMapDiffs
}

/*
BigMapDiffs represents $contract.big_map_diff in the tezos block schema

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id
*/
type BigMapDiffs []BigMapDiff

// BigMapDiff is an element of BigMapDiffs
type BigMapDiff struct {
	Action            BigMapDiffAction `json:"action,omitempty"`
	BigMap            string           `json:"big_map,omitempty"`
	KeyHash           string           `json:"key_hash,omitempty"`
	Key               *json.RawMessage `json:"key,omitempty"`
	Value             *json.RawMessage `json:"value,omitempty"`
	SourceBigMap      string           `json:"source_big_map,omitempty"`
	DestinationBigMap string           `json:"destination_big_map,omitempty"`
	KeyType           *json.RawMessage `json:"key_type,omitempty"`
	ValueType         *json.RawMessage `json:"value_type,omitempty"`
}

// ToUpdate converts BigMapDiff to BigMapDiffUpdate
func (b BigMapDiff) ToUpdate() BigMapDiffUpdate {
	return BigMapDiffUpdate{
		Action:  b.Action,
		BigMap:  b.BigMap,
		KeyHash: b.KeyHash,
		Key:     b.Key,
		Value:   b.Value,
	}
}

// ToRemove converts BigMapDiff to BigMapDiffRemove
func (b BigMapDiff) ToRemove() BigMapDiffRemove {
	return BigMapDiffRemove{
		Action: b.Action,
		BigMap: b.BigMap,
	}
}

// ToCopy converts BigMapDiff to BigMapDiffCopy
func (b BigMapDiff) ToCopy() BigMapDiffCopy {
	return BigMapDiffCopy{
		Action:            b.Action,
		SourceBigMap:      b.SourceBigMap,
		DestinationBigMap: b.DestinationBigMap,
	}
}

// ToAlloc converts BigMapDiff to BigMapDiffAlloc
func (b BigMapDiff) ToAlloc() BigMapDiffAlloc {
	return BigMapDiffAlloc{
		Action:    b.Action,
		BigMap:    b.BigMap,
		KeyType:   b.KeyType,
		ValueType: b.ValueType,
	}
}

// Organize converts BigMapDiffs into OrganizedBigMapDiff
func (b BigMapDiffs) Organize() OrganizedBigMapDiff {
	var organizedBigMapDiff OrganizedBigMapDiff
	for _, bigMapDiff := range b {
		if bigMapDiff.Action == UPDATE {
			organizedBigMapDiff.Updates = append(organizedBigMapDiff.Updates, bigMapDiff.ToUpdate())
		} else if bigMapDiff.Action == REMOVE {
			organizedBigMapDiff.Removals = append(organizedBigMapDiff.Removals, bigMapDiff.ToRemove())
		} else if bigMapDiff.Action == COPY {
			organizedBigMapDiff.Copies = append(organizedBigMapDiff.Copies, bigMapDiff.ToCopy())
		} else if bigMapDiff.Action == ALLOC {
			organizedBigMapDiff.Allocs = append(organizedBigMapDiff.Allocs, bigMapDiff.ToAlloc())
		}
	}

	return organizedBigMapDiff
}

/*
BigMapDiffUpdate represents $contract.big_map_diff in the tezos block schema

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id
*/
type BigMapDiffUpdate struct {
	Action  BigMapDiffAction `json:"action"`
	BigMap  string           `json:"big_map,omitempty"`
	KeyHash string           `json:"key_hash,omitempty"`
	Key     *json.RawMessage `json:"key"`
	Value   *json.RawMessage `json:"value,omitempty"`
}

func (b *BigMapDiffUpdate) toBigMapDiff() BigMapDiff {
	return BigMapDiff{
		Action:  b.Action,
		BigMap:  b.BigMap,
		KeyHash: b.KeyHash,
		Key:     b.Key,
		Value:   b.Value,
	}
}

/*
BigMapDiffRemove represents $contract.big_map_diff in the tezos block schema

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id
*/
type BigMapDiffRemove struct {
	Action BigMapDiffAction `json:"action"`
	BigMap string           `json:"big_map"`
}

func (b *BigMapDiffRemove) toBigMapDiff() BigMapDiff {
	return BigMapDiff{
		Action: b.Action,
		BigMap: b.BigMap,
	}
}

/*
BigMapDiffCopy represents $contract.big_map_diff in the tezos block schema

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id
*/
type BigMapDiffCopy struct {
	Action            BigMapDiffAction `json:"action"`
	SourceBigMap      string           `json:"source_big_map"`
	DestinationBigMap string           `json:"destination_big_map"`
}

func (b *BigMapDiffCopy) toBigMapDiff() BigMapDiff {
	return BigMapDiff{
		Action:            b.Action,
		SourceBigMap:      b.SourceBigMap,
		DestinationBigMap: b.DestinationBigMap,
	}
}

/*
BigMapDiffAlloc represents $contract.big_map_diff in the tezos block schema

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id
*/
type BigMapDiffAlloc struct {
	Action    BigMapDiffAction `json:"action"`
	BigMap    string           `json:"big_map"`
	KeyType   *json.RawMessage `json:"key_type"`
	ValueType *json.RawMessage `json:"value_type"`
}

func (b *BigMapDiffAlloc) toBigMapDiff() BigMapDiff {
	return BigMapDiff{
		Action:    b.Action,
		BigMap:    b.BigMap,
		KeyType:   b.KeyType,
		ValueType: b.ValueType,
	}
}

/*
ScriptedContracts represents $scripted.contracts in the tezos block schema

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id
*/
type ScriptedContracts struct {
	Code    *json.RawMessage `json:"code"`
	Storage *json.RawMessage `json:"storage"`
}

// BlockID represents an ID for a Block
type BlockID interface {
	ID() string
}

// BlockIDHead is the BlockID for the head block
type BlockIDHead struct{}

// ID satisfies the BlockID interface
func (b *BlockIDHead) ID() string {
	return "head"
}

// BlockIDLevel is the BlockID for a specific level
type BlockIDLevel int

// ID satisfies the BlockID interface
func (b *BlockIDLevel) ID() string {
	return strconv.Itoa(int(*b))
}

// BlockIDHash is the BlockID for a specific hash
type BlockIDHash string

// ID satisfies the BlockID interface
func (b *BlockIDHash) ID() string {
	return string(*b)
}

// BlockIDHeadPredecessor is a BlockID equivilent to head~<diff_level>
type BlockIDHeadPredecessor int

// ID satisfies the BlockID interface
func (b *BlockIDHeadPredecessor) ID() string {
	return fmt.Sprintf("head~%d", *b)
}

// BlockIDPredecessor is a BlockID equivilent to hash~<diff_level>
type BlockIDPredecessor struct {
	Hash      string
	DiffLevel int
}

// ID satisfies the BlockID interface
func (b *BlockIDPredecessor) ID() string {
	return fmt.Sprintf("%s~%d", b.Hash, b.DiffLevel)
}

/*
Block gets all the information about a specific block

Path
	/chains/<chain_id>/blocks/<block_id> (GET)
RPC
	https://tezos.gitlab.io/008/rpc.html#get-block-id
*/
func (c *Client) Block(blockID BlockID) (*resty.Response, *Block, error) {
	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s", c.chain, blockID.ID()))
	if err != nil {
		return resp, &Block{}, errors.Wrapf(err, "failed to get block '%s'", blockID.ID())
	}

	var block Block
	err = json.Unmarshal(resp.Body(), &block)
	if err != nil {
		return resp, &block, errors.Wrapf(err, "failed to get block '%s': failed to parse json", blockID.ID())
	}

	return resp, &block, nil
}

/*
EndorsingPowerInput is the input for the EndorsingPower function

Path
	 ../<block_id>/endorsing_power (POST)
RPC
	https://tezos.gitlab.io/008/rpc.html#post-block-id-endorsing-power
*/
type EndorsingPowerInput struct {
	// The block of which you want to make the query.
	BlockID BlockID
	// The cycle to get the balance at. If not provided Blockhash is required.
	Cycle int
	// The Operation you wish to get the endorsing power for
	EndorsingPower EndorsingPower
}

/*
EndorsingPower the body of the operation for endorsing power

Path
	 ../<block_id>/endorsing_power (POST)
RPC
	https://tezos.gitlab.io/008/rpc.html#post-block-id-endorsing-power
*/
type EndorsingPower struct {
	EndorsementOperation EndorsingOperation `json:"endorsement_operation"`
	ChainID              string             `json:"chain_id"`
}

/*
EndorsingOperation the body of the operation for endorsing power

Path
	 ../<block_id>/endorsing_power (POST)
RPC
	https://tezos.gitlab.io/008/rpc.html#post-block-id-endorsing-power
*/
type EndorsingOperation struct {
	Branch    string    `json:"branch"`
	Contents  []Content `json:"contents"`
	Signature string    `json:"signature"`
}

/*
EndorsingPower gets the endorsing power of an endorsement, that is, the number of slots that the endorser has

Path
	 ../<block_id>/endorsing_power (POST)
RPC
	https://tezos.gitlab.io/008/rpc.html#post-block-id-endorsing-power
*/
func (c *Client) EndorsingPower(input EndorsingPowerInput) (*resty.Response, int, error) {
	resp, blockID, err := c.processContextRequest(input, input.Cycle, input.BlockID)
	if err != nil {
		return resp, 0, errors.Wrap(err, "failed to get endorsing power")
	}

	resp, err = c.post(fmt.Sprintf("/chains/%s/blocks/%s/endorsing_power", c.chain, blockID.ID()), input.EndorsingPower)
	if err != nil {
		return resp, 0, errors.Wrap(err, "failed to get endorsing power")
	}

	var endorsingPower int
	err = json.Unmarshal(resp.Body(), &endorsingPower)
	if err != nil {
		return resp, 0, errors.Wrap(err, "failed to get endorsing power: failed to parse json")
	}

	return resp, endorsingPower, nil
}

/*
Hash gets the block's hash, its unique identifier.

Path
	  ../<block_id>/hash (GET)
RPC
	https://tezos.gitlab.io/008/rpc.html#get-block-id-hash
*/
func (c *Client) Hash(blockID BlockID) (*resty.Response, string, error) {
	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/hash", c.chain, blockID.ID()))
	if err != nil {
		return resp, "", errors.Wrapf(err, "failed to get block '%s' hash", blockID.ID())
	}

	var hash string
	err = json.Unmarshal(resp.Body(), &hash)
	if err != nil {
		return resp, "", errors.Wrapf(err, "failed to get block '%s' hash: failed to parse json", blockID.ID())
	}

	return resp, hash, nil
}

/*
Header gets the whole block header.

Path
	../<block_id>/header (GET)
RPC
	https://tezos.gitlab.io/008/rpc.html#get-block-id-header
*/
func (c *Client) Header(blockID BlockID) (*resty.Response, Header, error) {
	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/header", c.chain, blockID.ID()))
	if err != nil {
		return resp, Header{}, errors.Wrapf(err, "failed to get block '%s' header", blockID.ID())
	}

	var header Header
	err = json.Unmarshal(resp.Body(), &header)
	if err != nil {
		return resp, Header{}, errors.Wrapf(err, "failed to get block '%s' header: failed to parse json", blockID.ID())
	}

	return resp, header, nil
}

/*
HeaderRaw gets the whole block header (unparsed).

Path
	../<block_id>/header/raw (GET)
RPC
	https://tezos.gitlab.io/008/rpc.html#get-block-id-header-raw
*/
func (c *Client) HeaderRaw(blockID BlockID) (*resty.Response, string, error) {
	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/header/raw", c.chain, blockID.ID()))
	if err != nil {
		return resp, "", errors.Wrapf(err, "failed to get block '%s' raw header", blockID.ID())
	}

	var header string
	err = json.Unmarshal(resp.Body(), &header)
	if err != nil {
		return resp, "", errors.Wrapf(err, "failed to get block '%s' raw header: failed to parse json", blockID.ID())
	}

	return resp, header, nil
}

/*
HeaderShell is the shell-specific fragment of the block header.

Path
	../<block_id>/header/shell (GET)
RPC
	https://tezos.gitlab.io/008/rpc.html#get-block-id-header-shell
*/
type HeaderShell struct {
	Level          int       `json:"level"`
	Proto          int       `json:"proto"`
	Predecessor    string    `json:"predecessor"`
	Timestamp      time.Time `json:"timestamp"`
	ValidationPass int       `json:"validation_pass"`
	OperationsHash string    `json:"operations_hash"`
	Fitness        []string  `json:"fitness"`
	Context        string    `json:"context"`
}

/*
HeaderShell gets the shell-specific fragment of the block header.

Path
	../<block_id>/header/shell (GET)
RPC
	https://tezos.gitlab.io/008/rpc.html#get-block-id-header-shell
*/
func (c *Client) HeaderShell(blockID BlockID) (*resty.Response, HeaderShell, error) {
	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/header/shell", c.chain, blockID.ID()))
	if err != nil {
		return resp, HeaderShell{}, errors.Wrapf(err, "failed to get block '%s' header shell", blockID.ID())
	}

	var headerShell HeaderShell
	err = json.Unmarshal(resp.Body(), &headerShell)
	if err != nil {
		return resp, HeaderShell{}, errors.Wrapf(err, "failed to get block '%s' header shell: failed to parse json", blockID.ID())
	}

	return resp, headerShell, nil
}

/*
ProtocolData is the version-specific fragment of the block header.

Path
	../<block_id>/header/protocol_data (GET)
RPC
	https://tezos.gitlab.io/008/rpc.html#get-block-id-header-protocol-data
*/
type ProtocolData struct {
	Protocol         string `json:"protocol"`
	Priority         int    `json:"priority"`
	ProofOfWorkNonce string `json:"proof_of_work_nonce"`
	Signature        string `json:"signature"`
}

/*
HeaderProtocolData gets the version-specific fragment of the block header.

Path
	../<block_id>/header/protocol_data (GET)
RPC
	https://tezos.gitlab.io/008/rpc.html#get-block-id-header-protocol-data
*/
func (c *Client) HeaderProtocolData(blockID BlockID) (*resty.Response, ProtocolData, error) {
	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/header/protocol_data", c.chain, blockID.ID()))
	if err != nil {
		return resp, ProtocolData{}, errors.Wrapf(err, "failed to get block '%s' protocol data", blockID.ID())
	}

	var protocolData ProtocolData
	err = json.Unmarshal(resp.Body(), &protocolData)
	if err != nil {
		return resp, ProtocolData{}, errors.Wrapf(err, "failed to get block '%s' protocol data: failed to parse json", blockID.ID())
	}

	return resp, protocolData, nil
}

/*
HeaderProtocolDataRaw gets the version-specific fragment of the block header (unparsed).

Path
	../<block_id>/header/protocol_data/raw (GET)
RPC
	https://tezos.gitlab.io/008/rpc.html#get-block-id-header-protocol-data-raw
*/
func (c *Client) HeaderProtocolDataRaw(blockID BlockID) (*resty.Response, string, error) {
	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/header/protocol_data/raw", c.chain, blockID.ID()))
	if err != nil {
		return resp, "", errors.Wrapf(err, "failed to get block '%s' raw protocol data", blockID.ID())
	}

	var protocolData string
	err = json.Unmarshal(resp.Body(), &protocolData)
	if err != nil {
		return resp, "", errors.Wrapf(err, "failed to get block '%s' raw protocol data: failed to parse json", blockID.ID())
	}

	return resp, protocolData, nil
}

/*
LiveBlocks lists the ancestors of the given block which, if referred to as
the branch in an operation header, are recent enough for that operation to
be included in the current block.

Path
	../<block_id>/live_blocks (GET)
RPC
	https://tezos.gitlab.io/008/rpc.html#get-block-id-live-blocks
*/
func (c *Client) LiveBlocks(blockID BlockID) (*resty.Response, []string, error) {
	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/live_blocks", c.chain, blockID.ID()))
	if err != nil {
		return resp, []string{}, errors.Wrapf(err, "failed to get live blocks at '%s'", blockID.ID())
	}

	var liveBlocks []string
	err = json.Unmarshal(resp.Body(), &liveBlocks)
	if err != nil {
		return resp, []string{}, errors.Wrapf(err, "failed to get live blocks at '%s': failed to parse json", blockID.ID())
	}

	return resp, liveBlocks, nil
}

/*
Metadata returns all the metadata associated to the block.

Path
	../<block_id>/metadata (GET)
RPC
	https://tezos.gitlab.io/008/rpc.html#get-block-id-metadata
*/
func (c *Client) Metadata(blockID BlockID) (*resty.Response, Metadata, error) {
	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/metadata", c.chain, blockID.ID()))
	if err != nil {
		return resp, Metadata{}, errors.Wrapf(err, "failed to get block '%s' metadata", blockID.ID())
	}

	var metadata Metadata
	err = json.Unmarshal(resp.Body(), &metadata)
	if err != nil {
		return resp, Metadata{}, errors.Wrapf(err, "failed to get block '%s' metadata: failed to parse json", blockID.ID())
	}

	return resp, metadata, nil
}

/*
MetadataHash returns the Hash of the metadata associated to the block. This is only set on blocks starting from environment V1.

Path
	../<block_id>/metadata_hash (GET)
RPC
	https://tezos.gitlab.io/008/rpc.html#get-block-id-metadata-hash
*/
func (c *Client) MetadataHash(blockID BlockID) (*resty.Response, string, error) {
	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/metadata_hash", c.chain, blockID.ID()))
	if err != nil {
		return resp, "", errors.Wrapf(err, "failed to get block '%s' metadata hash", blockID.ID())
	}

	var metadataHash string
	err = json.Unmarshal(resp.Body(), &metadataHash)
	if err != nil {
		return resp, "", errors.Wrapf(err, "failed to get block '%s' metadata hash: failed to parse json", blockID.ID())
	}

	return resp, metadataHash, nil
}

/*
MinimalValidTimeInput is the input for the MinimalValidTime function

RPC
	https://tezos.gitlab.io/008/rpc.html#get-block-id-minimal-valid-time
*/
type MinimalValidTimeInput struct {
	// The block of which you want to make the query.
	BlockID        BlockID
	Priority       int
	EndorsingPower int
}

func (i *MinimalValidTimeInput) contructRPCOptions() []rpcOptions {
	var opts []rpcOptions
	opts = append(opts, rpcOptions{
		"priority",
		strconv.Itoa(i.Priority),
	})

	// Endorsing power
	opts = append(opts, rpcOptions{
		"endorsing_power",
		strconv.Itoa(i.EndorsingPower),
	})
	return opts
}


/*
MinimalValidTime returns the minimal valid time for a block given a priority and an endorsing power.

Path
	../<block_id>/minimal_valid_time?[priority=<int>]&[endorsing_power=<int>] (GET)
RPC
	https://tezos.gitlab.io/008/rpc.html#get-block-id-minimal-valid-time
*/
func (c *Client) MinimalValidTime(input MinimalValidTimeInput) (*resty.Response, time.Time, error) {
	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/minimal_valid_time", c.chain, input.BlockID.ID()), input.contructRPCOptions()...)
	if err != nil {
		return resp, time.Time{}, errors.Wrapf(err, "failed to get minimal valid time at '%s'", input.BlockID.ID())
	}

	var minimalValidTime time.Time
	err = json.Unmarshal(resp.Body(), &minimalValidTime)
	if err != nil {
		return resp, time.Time{}, errors.Wrapf(err, "failed to get minimal valid time at '%s': failed to parse json", input.BlockID.ID())
	}

	return resp, minimalValidTime, nil
}

/*
OperationHashesInput is the input to the OperationHashes function

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-operation-hashes
	https://tezos.gitlab.io/008/rpc.html#get-block-id-operation-hashes-list-offset
	https://tezos.gitlab.io/008/rpc.html#get-block-id-operation-hashes-list-offset-operation-offset
*/
type OperationHashesInput struct {
	// The block of which you want to make the query.
	BlockID         BlockID
	ListOffset      string
	OperationOffset string
}

func (o *OperationHashesInput) path(chain string) string {
	if o.ListOffset != "" && o.OperationOffset != "" {
		return fmt.Sprintf("/chains/%s/blocks/%s/operation_hashes/%s/%s", chain, o.BlockID.ID(), o.ListOffset, o.OperationOffset)
	}

	if o.ListOffset != "" && o.OperationOffset == "" {
		return fmt.Sprintf("/chains/%s/blocks/%s/operation_hashes/%s", chain, o.BlockID.ID(), o.ListOffset)
	}

	return fmt.Sprintf("/chains/%s/blocks/%s/operation_hashes", chain, o.BlockID.ID())
}

/*
OperationHashes is the operations hashes in the OperationHashes function

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-operation-hashes
	https://tezos.gitlab.io/008/rpc.html#get-block-id-operation-hashes-list-offset
	https://tezos.gitlab.io/008/rpc.html#get-block-id-operation-hashes-list-offset-operation-offset
*/
type OperationHashes []string

// UnmarshalJSON satisfies json.Marsheler
func (o *OperationHashes) UnmarshalJSON(b []byte) error {
	var flatOps []string
	var operations [][]string
	if err := json.Unmarshal(b, &operations); err != nil {
		var operations []string
		if err = json.Unmarshal(b, &operations); err != nil {
			var operation string
			if err = json.Unmarshal(b, &operation); err != nil {
				return err
			}
			flatOps = append(flatOps, operation)
		} else {
			flatOps = append(flatOps, operations...)
		}
	} else {
		for _, x := range operations {
			flatOps = append(flatOps, x...)
		}
	}

	*o = flatOps
	return nil
}

/*
OperationHashes returns the hashes of operations included in a block

Path:
	 ../<block_id>/operation_hashes (GET)
	../<block_id>/operation_hashes/<list_offset> (GET)
	../<block_id>/operation_hashes/<list_offset>/<operation_offset> (GET)

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-operation-hashes
	https://tezos.gitlab.io/008/rpc.html#get-block-id-operation-hashes-list-offset
	https://tezos.gitlab.io/008/rpc.html#get-block-id-operation-hashes-list-offset-operation-offset
*/
func (c *Client) OperationHashes(input OperationHashesInput) (*resty.Response, OperationHashes, error) {
	resp, err := c.get(input.path(c.chain))
	if err != nil {
		return nil, []string{}, errors.Wrapf(err, "failed to get block '%s' operation hashes", input.BlockID.ID())
	}

	var operationHashes OperationHashes
	err = json.Unmarshal(resp.Body(), &operationHashes)
	if err != nil {
		return resp, []string{}, errors.Wrapf(err, "failed to get block '%s' operation hashes: failed to parse json", input.BlockID.ID())
	}

	return resp, operationHashes, nil
}

/*
OperationMetadataHashesInput is the operations metadata hashes in the OperationMetadataHashes function

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-operation-hashes
	https://tezos.gitlab.io/008/rpc.html#get-block-id-operation-hashes-list-offset
	https://tezos.gitlab.io/008/rpc.html#get-block-id-operation-hashes-list-offset-operation-offset
*/
type OperationMetadataHashesInput struct {
	// The block of which you want to make the query.
	BlockID         BlockID
	ListOffset      string
	OperationOffset string
}

func (o *OperationMetadataHashesInput) path(chain string) string {
	if o.ListOffset != "" && o.OperationOffset != "" {
		return fmt.Sprintf("/chains/%s/blocks/%s/operation_metadata_hashes/%s/%s", chain, o.BlockID.ID(), o.ListOffset, o.OperationOffset)
	}

	if o.ListOffset != "" && o.OperationOffset == "" {
		return fmt.Sprintf("/chains/%s/blocks/%s/operation_metadata_hashes/%s", chain, o.BlockID.ID(), o.ListOffset)
	}

	return fmt.Sprintf("/chains/%s/blocks/%s/operation_metadata_hashes", chain, o.BlockID.ID())
}

/*
OperationMetadataHashes is the operations hashes in the OperationMetadataHashes function

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-operation-hashes
	https://tezos.gitlab.io/008/rpc.html#get-block-id-operation-hashes-list-offset
	https://tezos.gitlab.io/008/rpc.html#get-block-id-operation-hashes-list-offset-operation-offset
*/
type OperationMetadataHashes []string

// UnmarshalJSON satisfies json.Marsheler
func (o *OperationMetadataHashes) UnmarshalJSON(b []byte) error {
	var flatOps []string
	var operations [][]string
	if err := json.Unmarshal(b, &operations); err != nil {
		var operations []string
		if err = json.Unmarshal(b, &operations); err != nil {
			var operation string
			if err = json.Unmarshal(b, &operation); err != nil {
				return err
			}
			flatOps = append(flatOps, operation)
		} else {
			flatOps = append(flatOps, operations...)
		}
	} else {
		for _, x := range operations {
			flatOps = append(flatOps, x...)
		}
	}

	*o = flatOps
	return nil
}

/*
OperationMetadataHashes returns the hashes of all the operation metadata included in the block.
This is only set on blocks starting from environment V1.

Path:
	 ../<block_id>/operation_metadata_hashes (GET)
	../<block_id>/operation_metadata_hashes/<list_offset> (GET)
	../<block_id>/operation_metadata_hashes/<list_offset>/<operation_offset> (GET)

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-operation-metadata-hashes
	https://tezos.gitlab.io/008/rpc.html#get-block-id-operation-metadata-hashes-list-offset
	https://tezos.gitlab.io/008/rpc.html#get-block-id-operation-metadata-hashes-list-offset-operation-offset
*/
func (c *Client) OperationMetadataHashes(input OperationMetadataHashesInput) (*resty.Response, OperationMetadataHashes, error) {
	resp, err := c.get(input.path(c.chain))
	if err != nil {
		return nil, []string{}, errors.Wrapf(err, "failed to get block '%s' operation metadata hashes", input.BlockID.ID())
	}

	var operationMetadataHashes OperationMetadataHashes
	err = json.Unmarshal(resp.Body(), &operationMetadataHashes)
	if err != nil {
		return resp, []string{}, errors.Wrapf(err, "failed to get block '%s' operation metadata hashes: failed to parse json", input.BlockID.ID())
	}

	return resp, operationMetadataHashes, nil
}

/*
OperationsInput is the input for the Operations function

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-operations
	https://tezos.gitlab.io/008/rpc.html#get-block-id-operations-list-offset
	https://tezos.gitlab.io/008/rpc.html#get-block-id-operations-list-offset-operation-offset
*/
type OperationsInput struct {
	// The block of which you want to make the query.
	BlockID         BlockID
	ListOffset      string
	OperationOffset string
}

func (o *OperationsInput) path(chain string) string {
	if o.ListOffset != "" && o.OperationOffset != "" {
		return fmt.Sprintf("/chains/%s/blocks/%s/operations/%s/%s", chain, o.BlockID.ID(), o.ListOffset, o.OperationOffset)
	}

	if o.ListOffset != "" && o.OperationOffset == "" {
		return fmt.Sprintf("/chains/%s/blocks/%s/operations/%s", chain, o.BlockID.ID(), o.ListOffset)
	}

	return fmt.Sprintf("/chains/%s/blocks/%s/operations", chain, o.BlockID.ID())
}

/*
FlattenedOperations is Opperations expressed in a single slice

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-operations
	https://tezos.gitlab.io/008/rpc.html#get-block-id-operations-list-offset
	https://tezos.gitlab.io/008/rpc.html#get-block-id-operations-list-offset-operation-offset
*/
type FlattenedOperations []Operations

// UnmarshalJSON satisfies json.Marsheler
func (f *FlattenedOperations) UnmarshalJSON(b []byte) error {
	var flatOps []Operations
	var operations [][]Operations
	if err := json.Unmarshal(b, &operations); err != nil {
		var operations []Operations
		if err = json.Unmarshal(b, &operations); err != nil {
			var operation Operations
			if err = json.Unmarshal(b, &operation); err != nil {
				return err
			}
			flatOps = append(flatOps, operation)
		} else {
			flatOps = append(flatOps, operations...)
		}
	} else {
		for _, x := range operations {
			flatOps = append(flatOps, x...)
		}
	}

	*f = flatOps
	return nil
}

/*
Operations gets the operations included in a block

Path:
	 ../<block_id>/operations (GET)
	../<block_id>/operations/<list_offset> (GET)
	../<block_id>/operations/<list_offset>/<operation_offset> (GET)

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-operations
	https://tezos.gitlab.io/008/rpc.html#get-block-id-operations-list-offset
	https://tezos.gitlab.io/008/rpc.html#get-block-id-operations-list-offset-operation-offset
*/
func (c *Client) Operations(input OperationsInput) (*resty.Response, FlattenedOperations, error) {
	resp, err := c.get(input.path(c.chain))
	if err != nil {
		return nil, FlattenedOperations{}, errors.Wrapf(err, "failed to get block '%s' operations", input.BlockID.ID())
	}

	var operations FlattenedOperations
	err = json.Unmarshal(resp.Body(), &operations)
	if err != nil {
		return resp, FlattenedOperations{}, errors.Wrapf(err, "failed to get block '%s' operations: failed to parse json", input.BlockID.ID())
	}

	return resp, operations, nil
}

/*
OperationsMetadataHash returns the root hash of the operations metadata from the block.
This is only set on blocks starting from environment V1.

Path:
	../<block_id>/operations_metadata_hash (GET)

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-operations-metadata-hash
*/
func (c *Client) OperationsMetadataHash(blockID BlockID) (*resty.Response, string, error) {
	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/operations_metadata_hash", c.chain, blockID.ID()))
	if err != nil {
		return nil, "", errors.Wrapf(err, "failed to get block '%s' operations metadata hash", blockID.ID())
	}

	var metadataHash string
	err = json.Unmarshal(resp.Body(), &metadataHash)
	if err != nil {
		return resp, "", errors.Wrapf(err, "failed to get block '%s' operations metadata hash: failed to parse json", blockID.ID())
	}

	return resp, metadataHash, nil
}

/*
Protocols is the current and next protocol.

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-protocols
*/
type Protocols struct {
	Protocol     string `json:"protocol"`
	NextProtocol string `json:"next_protocol"`
}

/*
Protocols returns the current and next protocol.

Path:
	../<block_id>/protocols (GET)

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-protocols
*/
func (c *Client) Protocols(blockID BlockID) (*resty.Response, Protocols, error) {
	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/protocols", c.chain, blockID.ID()))
	if err != nil {
		return nil, Protocols{}, errors.Wrapf(err, "failed to get block '%s' protocols", blockID.ID())
	}

	var protocols Protocols
	err = json.Unmarshal(resp.Body(), &protocols)
	if err != nil {
		return resp, Protocols{}, errors.Wrapf(err, "failed to get block '%s' protocols: failed to parse json", blockID.ID())
	}

	return resp, protocols, nil
}

/*
RequiredEndorsementsInput is the input for RequiredEndorsements functions.

Path:
	../<block_id>/required_endorsements?[block_delay=<int64>] (GET)

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-required-endorsements
*/
type RequiredEndorsementsInput struct {
	// The block of which you want to make the query.
	BlockID    BlockID
	BlockDelay int64
}

func (r *RequiredEndorsementsInput) constructRPCOptions() []rpcOptions {
	var options []rpcOptions
	if r.BlockDelay != 0 {
		options = append(options, rpcOptions{
			"block_delay",
			fmt.Sprintf("%d", r.BlockDelay),
		})
	}

	return options
}

/*
RequiredEndorsements returns the minimum number of endorsements for a block to be valid, given a delay of the block's timestamp with respect to the minimum time to bake at the block's priority

Path:
	../<block_id>/required_endorsements?[block_delay=<int64>] (GET)

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-required-endorsements
*/
func (c *Client) RequiredEndorsements(input RequiredEndorsementsInput) (*resty.Response, int, error) {
	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/required_endorsements", c.chain, input.BlockID.ID()), input.constructRPCOptions()...)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "failed to get block '%s' required endorsements", input.BlockID.ID())
	}

	var endrosements int
	err = json.Unmarshal(resp.Body(), &endrosements)
	if err != nil {
		return resp, 0, errors.Wrapf(err, "failed to get block '%s' required endorsements: failed to parse json", input.BlockID.ID())
	}

	return resp, endrosements, nil
}
