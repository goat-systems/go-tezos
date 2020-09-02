package gotezos

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

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
	Change   int64  `json:"change,string"`
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
	ConsumedGas         int64            `json:"consumed_gas,string,omitempty"`
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
	Protocol  string   `json:"protocol,omitempty"`
	ChainID   string   `json:"chain_id,omitempty"`
	Hash      string   `json:"hash,omitempty"`
	Branch    string   `json:"branch"`
	Contents  Contents `json:"contents"`
	Signature string   `json:"signature,omitempty"`
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
	Transactions              []Transaction
	Originations              []Origination
	Delegations               []Delegation
}

// ContentsHelper used for unmarshaling and marshaling json block contents
type ContentsHelper struct {
	Kind          string                    `json:"kind,omitempty"`
	Level         int                       `json:"level,omitempty"`
	Nonce         string                    `json:"nonce,omitempty"`
	Op1           *InlinedEndorsement       `json:"Op1,omitempty"`
	Op2           *InlinedEndorsement       `json:"Op2,omitempty"`
	Pkh           string                    `json:"pkh,omitempty"`
	Secret        string                    `json:"secret,omitempty"`
	Bh1           *BlockHeader              `json:"bh1,omitempty"`
	Bh2           *BlockHeader              `json:"bh2,omitempty"`
	Source        string                    `json:"source,omitempty"`
	Period        int                       `json:"period,omitempty"`
	Proposals     []string                  `json:"proposals,omitempty"`
	Proposal      string                    `json:"proposal,omitempty"`
	Ballot        string                    `json:"ballot,omitempty"`
	Fee           int64                     `json:"fee,string,omitempty"`
	Counter       int                       `json:"counter,string,omitempty"`
	GasLimit      int64                     `json:"gas_limit,string,omitempty"`
	StorageLimit  int64                     `json:"storage_limit,string,omitempty"`
	PublicKey     string                    `json:"public_key,omitempty"`
	ManagerPubkey string                    `json:"managerPubKey,omitempty"`
	Amount        int64                     `json:"amount,string,omitempty"`
	Destination   string                    `json:"destination,omitempty"`
	Balance       int64                     `json:"balance,string,omitempty"`
	Delegate      string                    `json:"delegate,omitempty"`
	Script        Script                    `json:"script,omitempty"`
	Parameters    *ContentsHelperParameters `json:"parameters,omitempty"`
	Metadata      *ContentsHelperMetadata   `json:"metadata,omitempty"`
}

// ContentsHelperParameters used for unmarshaling and marshaling json block contents
type ContentsHelperParameters struct {
	Entrypoint string               `json:"entrypoint"`
	Value      *MichelineExpression `json:"value"`
}

// ContentsHelperMetadata used for unmarshaling and marshaling json block contents
type ContentsHelperMetadata struct {
	BalanceUpdates          []BalanceUpdates           `json:"balance_updates,omitempty"`
	Delegate                string                     `json:"delegate,omitempty"`
	Slots                   []int                      `json:"slots,omitempty"`
	OperationResults        *OperationResultsHelper    `json:"operation_result,omitempty"`
	InternalOperationResult []InternalOperationResults `json:"internal_operation_result,omitempty"`
}

// OperationResultsHelper is a helper to unmarhsal and marshal OperationResults data
type OperationResultsHelper struct {
	Status                       string               `json:"status"`
	BigMapDiff                   *BigMapDiff          `json:"big_map_diff,omitempty"`
	BalanceUpdates               []BalanceUpdates     `json:"balance_updates,omitempty"`
	OriginatedContracts          []string             `json:"originated_contracts,omitempty"`
	ConsumedGas                  int64                `json:"consumed_gas,string,omitempty"`
	StorageSize                  int64                `json:"storage_size,string,omitempty"`
	PaidStorageSizeDiff          int64                `json:"paid_storage_size_diff,string,omitempty"`
	Errors                       []RPCError           `json:"errors,omitempty"`
	Storage                      *MichelineExpression `json:"storage,omitempty"`
	AllocatedDestinationContract bool                 `json:"allocated_destination_contract,omitempty"`
}

func (o *OperationResultsHelper) toOperationResultsReveal() OperationResultReveal {
	return OperationResultReveal{
		Status:      o.Status,
		ConsumedGas: o.ConsumedGas,
		Errors:      o.Errors,
	}
}

func (o *OperationResultsHelper) toOperationResultsTransfer() OperationResultTransfer {
	var (
		storage    *MichelineExpression
		bigMapDiff *BigMapDiff
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

func (o *OperationResultsHelper) toOperationResultsOrigination() OperationResultOrigination {
	var bigMapDiff *BigMapDiff

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

func (o *OperationResultsHelper) toOperationResultsDelegation() OperationResultDelegation {
	return OperationResultDelegation{
		Status:      o.Status,
		ConsumedGas: o.ConsumedGas,
		Errors:      o.Errors,
	}
}

// ToEndorsement converts ContentsHelper to an endorsement.
func (c *ContentsHelper) toEndorsement() Endorsement {
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

// ToSeedNonceRevelations converts ContentsHelper to an SeedNonceRevelations.
func (c *ContentsHelper) toSeedNonceRevelations() SeedNonceRevelation {
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

// ToDoubleEndorsementEvidence converts ContentsHelper to an DoubleEndorsementEvidence.
func (c *ContentsHelper) toDoubleEndorsementEvidence() DoubleEndorsementEvidence {
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

// ToDoubleBakingEvidence converts ContentsHelper to an DoubleBakingEvidence.
func (c *ContentsHelper) toDoubleBakingEvidence() DoubleBakingEvidence {
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

// ToAccountActivation converts ContentsHelper to an AccountActivation.
func (c *ContentsHelper) toAccountActivation() AccountActivation {
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

// ToProposal converts ContentsHelper to an Proposal.
func (c *ContentsHelper) toProposal() Proposal {
	return Proposal{
		Kind:      c.Kind,
		Source:    c.Source,
		Period:    c.Period,
		Proposals: c.Proposals,
	}
}

// ToBallot converts ContentsHelper to an Proposal.
func (c *ContentsHelper) toBallot() Ballot {
	return Ballot{
		Kind:     c.Kind,
		Source:   c.Source,
		Period:   c.Period,
		Proposal: c.Proposal,
		Ballot:   c.Ballot,
	}
}

// ToReveal converts ContentsHelper to a Reveal.
func (c *ContentsHelper) toReveal() Reveal {
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

// ToTransaction converts ContentsHelper to a Transaction.
func (c *ContentsHelper) toTransaction() Transaction {
	var (
		metadata   *TransactionMetadata
		parameters *TransactionParameters
	)

	if c.Metadata != nil {
		metadata = &TransactionMetadata{
			BalanceUpdates:           c.Metadata.BalanceUpdates,
			OperationResult:          c.Metadata.OperationResults.toOperationResultsTransfer(),
			InternalOperationResults: c.Metadata.InternalOperationResult,
		}
	}

	if c.Parameters != nil {
		parameters = &TransactionParameters{
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

// ToOrigination converts ContentsHelper to a Origination.
func (c *ContentsHelper) toOrigination() Origination {
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

// ToDelegation converts ContentsHelper to a Origination.
func (c *ContentsHelper) toDelegation() Delegation {
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

//UnmarshalJSON satisfies the json.Unmarshal interface for contents
func (c *Contents) UnmarshalJSON(v []byte) error {
	var contentsHelper []ContentsHelper
	err := json.Unmarshal(v, &contentsHelper)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal contents into ContentsHelper")
	}

	for _, content := range contentsHelper {
		switch content.Kind {
		case "endorsement":
			c.Endorsements = append(c.Endorsements, content.toEndorsement())
		case "seed_nonce_revelation":
			c.SeedNonceRevelations = append(c.SeedNonceRevelations, content.toSeedNonceRevelations())
		case "double_endorsement_evidence":
			c.DoubleEndorsementEvidence = append(c.DoubleEndorsementEvidence, content.toDoubleEndorsementEvidence())
		case "double_baking_evidence":
			c.DoubleBakingEvidence = append(c.DoubleBakingEvidence, content.toDoubleBakingEvidence())
		case "activate_account":
			c.AccountActivations = append(c.AccountActivations, content.toAccountActivation())
		case "proposals":
			c.Proposals = append(c.Proposals, content.toProposal())
		case "ballot":
			c.Ballots = append(c.Ballots, content.toBallot())
		case "reveal":
			c.Reveals = append(c.Reveals, content.toReveal())
		case "transaction":
			c.Transactions = append(c.Transactions, content.toTransaction())
		case "origination":
			c.Originations = append(c.Originations, content.toOrigination())
		case "delegation":
			c.Delegations = append(c.Delegations, content.toDelegation())
		default:
			return errors.New("failed to map contents to valid operation")
		}
	}

	return nil
}

//MarshalJSON satisfies the json.MarshalJSON interface for contents
func (c *Contents) MarshalJSON() ([]byte, error) {
	var contentsHelper []ContentsHelper

	for _, endorsement := range c.Endorsements {
		contentsHelper = append(contentsHelper, endorsement.toContentsHelper())
	}

	for _, seedNonceRevelation := range c.SeedNonceRevelations {
		contentsHelper = append(contentsHelper, seedNonceRevelation.toContentsHelper())
	}

	for _, doubleEndorsementEvidence := range c.DoubleEndorsementEvidence {
		contentsHelper = append(contentsHelper, doubleEndorsementEvidence.toContentsHelper())
	}

	for _, doubleBakingEvidence := range c.DoubleBakingEvidence {
		contentsHelper = append(contentsHelper, doubleBakingEvidence.toContentsHelper())
	}

	for _, accountActivation := range c.AccountActivations {
		contentsHelper = append(contentsHelper, accountActivation.toContentsHelper())
	}

	for _, proposal := range c.Proposals {
		contentsHelper = append(contentsHelper, proposal.toContentsHelper())
	}

	for _, ballot := range c.Ballots {
		contentsHelper = append(contentsHelper, ballot.toContentsHelper())
	}

	for _, reveal := range c.AccountActivations {
		contentsHelper = append(contentsHelper, reveal.toContentsHelper())
	}

	for _, transaction := range c.Transactions {
		contentsHelper = append(contentsHelper, transaction.toContentsHelper())
	}

	for _, origination := range c.Originations {
		contentsHelper = append(contentsHelper, origination.toContentsHelper())
	}

	for _, delegation := range c.Delegations {
		contentsHelper = append(contentsHelper, delegation.toContentsHelper())
	}

	return json.Marshal(&contentsHelper)
}

/*
Endorsement represents an endorsement in the $operation.alpha.operation_contents_and_result in the tezos block schema
See: tezos-client RPC format GET /chains/main/blocks/head
*/
type Endorsement struct {
	Kind     string               `json:"kind"`
	Level    int                  `json:"level"`
	Metadata *EndorsementMetadata `json:"metadata"`
}

/*
EndorsementMetadata represents the metadata of an endorsement in the $operation.alpha.operation_contents_and_result in the tezos block schema
See: tezos-client RPC format GET /chains/main/blocks/head
*/
type EndorsementMetadata struct {
	BalanceUpdates []BalanceUpdates `json:"balance_updates"`
	Delegate       string           `json:"delegate"`
	Slots          []int            `json:"slots"`
}

func (e *Endorsement) toContentsHelper() ContentsHelper {
	var metadata *ContentsHelperMetadata

	if e.Metadata != nil {
		metadata = &ContentsHelperMetadata{
			BalanceUpdates: e.Metadata.BalanceUpdates,
			Delegate:       e.Metadata.Delegate,
			Slots:          e.Metadata.Slots,
		}
	}

	return ContentsHelper{
		Kind:     e.Kind,
		Level:    e.Level,
		Metadata: metadata,
	}
}

/*
SeedNonceRevelation represents an Seed_nonce_revelation in the $operation.alpha.operation_contents_and_result in the tezos block schema
See: tezos-client RPC format GET /chains/main/blocks/head
*/
type SeedNonceRevelation struct {
	Kind     string                       `json:"kind"`
	Level    int                          `json:"level"`
	Nonce    string                       `json:"nonce"`
	Metadata *SeedNonceRevelationMetadata `json:"metadata"`
}

func (s *SeedNonceRevelation) toContentsHelper() ContentsHelper {
	var metadata *ContentsHelperMetadata

	if s.Metadata != nil {
		metadata = &ContentsHelperMetadata{
			BalanceUpdates: s.Metadata.BalanceUpdates,
		}
	}

	return ContentsHelper{
		Kind:     s.Kind,
		Level:    s.Level,
		Nonce:    s.Nonce,
		Metadata: metadata,
	}
}

/*
SeedNonceRevelationMetadata represents the metadata for Seed_nonce_revelation in the $operation.alpha.operation_contents_and_result in the tezos block schema
See: tezos-client RPC format GET /chains/main/blocks/head
*/
type SeedNonceRevelationMetadata struct {
	BalanceUpdates []BalanceUpdates `json:"balance_updates"`
}

/*
DoubleEndorsementEvidence represents an Double_endorsement_evidence in the $operation.alpha.operation_contents_and_result in the tezos block schema
See: tezos-client RPC format GET /chains/main/blocks/head
*/
type DoubleEndorsementEvidence struct {
	Kind     string                             `json:"kind"`
	Op1      *InlinedEndorsement                `json:"Op1"`
	Op2      *InlinedEndorsement                `json:"Op2"`
	Metadata *DoubleEndorsementEvidenceMetadata `json:"metadata"`
}

/*
DoubleEndorsementEvidenceMetadata represents the metadata for Double_endorsement_evidence in the $operation.alpha.operation_contents_and_result in the tezos block schema
See: tezos-client RPC format GET /chains/main/blocks/head
*/
type DoubleEndorsementEvidenceMetadata struct {
	BalanceUpdates []BalanceUpdates `json:"balance_updates"`
}

/*
InlinedEndorsement represents $inlined.endorsement in the tezos block schema
See: tezos-client RPC format GET /chains/main/blocks/head
*/
type InlinedEndorsement struct {
	Branch     string                        `json:"branch"`
	Operations *InlinedEndorsementOperations `json:"operations"`
	Signature  string                        `json:"signature"`
}

/*
InlinedEndorsementOperations represents operations in $inlined.endorsement in the tezos block schema
See: tezos-client RPC format GET /chains/main/blocks/head
*/
type InlinedEndorsementOperations struct {
	Kind  string `json:"kind"`
	Level int    `json:"level"`
}

func (d *DoubleEndorsementEvidence) toContentsHelper() ContentsHelper {
	var (
		metadata *ContentsHelperMetadata
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
		metadata = &ContentsHelperMetadata{
			BalanceUpdates: d.Metadata.BalanceUpdates,
		}
	}

	return ContentsHelper{
		Kind:     d.Kind,
		Op1:      op1,
		Op2:      op2,
		Metadata: metadata,
	}
}

/*
DoubleBakingEvidence represents an Double_baking_evidence in the $operation.alpha.operation_contents_and_result in the tezos block schema
See: tezos-client RPC format GET /chains/main/blocks/head
*/
type DoubleBakingEvidence struct {
	Kind     string                        `json:"kind"`
	Bh1      *BlockHeader                  `json:"bh1"`
	Bh2      *BlockHeader                  `json:"bh2"`
	Metadata *DoubleBakingEvidenceMetadata `json:"metadata"`
}

/*
DoubleBakingEvidenceMetadata represents the metadata of Double_baking_evidence in the $operation.alpha.operation_contents_and_result in the tezos block schema
See: tezos-client RPC format GET /chains/main/blocks/head
*/
type DoubleBakingEvidenceMetadata struct {
	BalanceUpdates []BalanceUpdates `json:"balance_updates"`
}

/*
BlockHeader represents $block_header.alpha.full_header in the tezos block schema
See: tezos-client RPC format GET /chains/main/blocks/head
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

func (d *DoubleBakingEvidence) toContentsHelper() ContentsHelper {
	var (
		metadata *ContentsHelperMetadata
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
		metadata = &ContentsHelperMetadata{
			BalanceUpdates: d.Metadata.BalanceUpdates,
		}
	}

	return ContentsHelper{
		Kind:     d.Kind,
		Bh1:      bh1,
		Bh2:      bh2,
		Metadata: metadata,
	}
}

/*
AccountActivation represents an Activate_account in the $operation.alpha.operation_contents_and_result in the tezos block schema
See: tezos-client RPC format GET /chains/main/blocks/head
*/
type AccountActivation struct {
	Kind     string                     `json:"kind"`
	Pkh      string                     `json:"pkh"`
	Secret   string                     `json:"secret"`
	Metadata *AccountActivationMetadata `json:"metadata"`
}

/*
AccountActivationMetadata represents the metadata for Activate_account in the $operation.alpha.operation_contents_and_result in the tezos block schema
See: tezos-client RPC format GET /chains/main/blocks/head
*/
type AccountActivationMetadata struct {
	BalanceUpdates []BalanceUpdates `json:"balance_updates"`
}

func (a *AccountActivation) toContentsHelper() ContentsHelper {
	var metadata *ContentsHelperMetadata
	if a.Metadata != nil {
		metadata = &ContentsHelperMetadata{
			BalanceUpdates: a.Metadata.BalanceUpdates,
		}
	}

	return ContentsHelper{
		Kind:     a.Kind,
		Pkh:      a.Pkh,
		Secret:   a.Secret,
		Metadata: metadata,
	}
}

/*
Proposal represents a Proposal in the $operation.alpha.operation_contents_and_result in the tezos block schema
See: tezos-client RPC format GET /chains/main/blocks/head
*/
type Proposal struct {
	Kind      string   `json:"kind"`
	Source    string   `json:"source"`
	Period    int      `json:"period"`
	Proposals []string `json:"proposals"`
}

func (p *Proposal) toContentsHelper() ContentsHelper {
	return ContentsHelper{
		Kind:      p.Kind,
		Source:    p.Source,
		Period:    p.Period,
		Proposals: p.Proposals,
	}
}

/*
Ballot represents a Ballot in the $operation.alpha.operation_contents_and_result in the tezos block schema
See: tezos-client RPC format GET /chains/main/blocks/head
*/
type Ballot struct {
	Kind     string `json:"kind"`
	Source   string `json:"source"`
	Period   int    `json:"period"`
	Proposal string `json:"proposal"`
	Ballot   string `json:"ballot"`
}

func (b *Ballot) toContentsHelper() ContentsHelper {
	return ContentsHelper{
		Kind:     b.Kind,
		Source:   b.Source,
		Period:   b.Period,
		Proposal: b.Proposal,
		Ballot:   b.Ballot,
	}
}

/*
Reveal represents a Reveal in the $operation.alpha.operation_contents_and_result in the tezos block schema
See: tezos-client RPC format GET /chains/main/blocks/head
*/
type Reveal struct {
	Kind         string          `json:"kind"`
	Source       string          `json:"source" validate:"required"`
	Fee          int64           `json:"fee,string" validate:"required"`
	Counter      int             `json:"counter,string" validate:"required"`
	GasLimit     int64           `json:"gas_limit,string" validate:"required"`
	StorageLimit int64           `json:"storage_limit,string" validate:"required"`
	PublicKey    string          `json:"public_key" validate:"required"`
	Metadata     *RevealMetadata `json:"metadata"`
}

/*
RevealMetadata represents the metadata for Reveal in the $operation.alpha.operation_contents_and_result in the tezos block schema
See: tezos-client RPC format GET /chains/main/blocks/head
*/
type RevealMetadata struct {
	BalanceUpdates           []BalanceUpdates           `json:"balance_updates"`
	OperationResult          OperationResultReveal      `json:"operation_result"`
	InternalOperationResults []InternalOperationResults `json:"internal_operation_result,omitempty"`
}

func (r *Reveal) toContentsHelper() ContentsHelper {
	var metadata *ContentsHelperMetadata

	if r.Metadata != nil {
		metadata = &ContentsHelperMetadata{
			BalanceUpdates: r.Metadata.BalanceUpdates,
			OperationResults: &OperationResultsHelper{
				Status:      r.Metadata.OperationResult.Status,
				ConsumedGas: r.Metadata.OperationResult.ConsumedGas,
				Errors:      r.Metadata.OperationResult.Errors,
			},
			InternalOperationResult: r.Metadata.InternalOperationResults,
		}
	}

	return ContentsHelper{
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
See: tezos-client RPC format GET /chains/main/blocks/head
*/
type Transaction struct {
	Kind         string                 `json:"kind"`
	Source       string                 `json:"source" validate:"required"`
	Fee          int64                  `json:"fee,string" validate:"required"`
	Counter      int                    `json:"counter,string" validate:"required"`
	GasLimit     int64                  `json:"gas_limit,string" validate:"required"`
	StorageLimit int64                  `json:"storage_limit,string"`
	Amount       int64                  `json:"amount,string"`
	Destination  string                 `json:"destination" validate:"required"`
	Parameters   *TransactionParameters `json:"parameters,omitempty"`
	Metadata     *TransactionMetadata   `json:"metadata"`
}

/*
TransactionParameters represents the parameters of a Transaction in the $operation.alpha.operation_contents_and_result in the tezos block schema
See: tezos-client RPC format GET /chains/main/blocks/head
*/
type TransactionParameters struct {
	Entrypoint string               `json:"entrypoint"`
	Value      *MichelineExpression `json:"value"`
}

/*
TransactionMetadata represents the metadata of Transaction in the $operation.alpha.operation_contents_and_result in the tezos block schema
See: tezos-client RPC format GET /chains/main/blocks/head
*/
type TransactionMetadata struct {
	BalanceUpdates           []BalanceUpdates           `json:"balance_updates"`
	OperationResult          OperationResultTransfer    `json:"operation_result"`
	InternalOperationResults []InternalOperationResults `json:"internal_operation_results,omitempty"`
}

func (t *Transaction) toContentsHelper() ContentsHelper {
	var (
		parameters *ContentsHelperParameters
		metadata   *ContentsHelperMetadata
	)

	if t.Parameters != nil {
		parameters = &ContentsHelperParameters{
			Entrypoint: t.Parameters.Entrypoint,
			Value:      t.Parameters.Value,
		}
	}

	if t.Metadata != nil {
		metadata = &ContentsHelperMetadata{
			BalanceUpdates: t.Metadata.BalanceUpdates,
			OperationResults: &OperationResultsHelper{
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

	return ContentsHelper{
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
See: tezos-client RPC format GET /chains/main/blocks/head
*/
type Origination struct {
	Kind          string               `json:"kind"`
	Source        string               `json:"source" validate:"required"`
	Fee           int64                `json:"fee,string" validate:"required"`
	Counter       int                  `json:"counter,string" validate:"required"`
	GasLimit      int64                `json:"gas_limit,string" validate:"required"`
	StorageLimit  int64                `json:"storage_limit,string" validate:"required"`
	Balance       int64                `json:"balance,string"`
	Delegate      string               `json:"delegate,omitempty"`
	Script        Script               `json:"script" validate:"required"`
	ManagerPubkey string               `json:"managerPubkey,omitempty"`
	Metadata      *OriginationMetadata `json:"metadata"`
}

/*
Script represents the script in an Origination in the $operation.alpha.operation_contents_and_result -> $scripted.contracts in the tezos block schema
See: tezos-client RPC format GET /chains/main/blocks/head
*/
type Script struct {
	Code    *MichelineExpression `json:"code,omitempty"`
	Storage *MichelineExpression `json:"storage,omitempty"`
}

/*
OriginationMetadata represents the metadata of Origination in the $operation.alpha.operation_contents_and_result in the tezos block schema
See: tezos-client RPC format GET /chains/main/blocks/head
*/
type OriginationMetadata struct {
	BalanceUpdates           []BalanceUpdates           `json:"balance_updates"`
	OperationResults         OperationResultOrigination `json:"operation_result"`
	InternalOperationResults []InternalOperationResults `json:"internal_operation_results,omitempty"`
}

func (o *Origination) toContentsHelper() ContentsHelper {
	var metadata *ContentsHelperMetadata

	if o.Metadata != nil {
		metadata = &ContentsHelperMetadata{
			BalanceUpdates: o.Metadata.BalanceUpdates,
			OperationResults: &OperationResultsHelper{
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

	return ContentsHelper{
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
See: tezos-client RPC format GET /chains/main/blocks/head
*/
type Delegation struct {
	Kind         string              `json:"kind"`
	Source       string              `json:"source" validate:"required"`
	Fee          int64               `json:"fee,string" validate:"required"`
	Counter      int                 `json:"counter,string" validate:"required"`
	GasLimit     int64               `json:"gas_limit,string" validate:"required"`
	StorageLimit int64               `json:"storage_limit,string" validate:"required"`
	Delegate     string              `json:"delegate,omitempty"`
	Metadata     *DelegationMetadata `json:"metadata"`
}

/*
DelegationMetadata represents the metadata Delegation in the $operation.alpha.operation_contents_and_result in the tezos block schema
See: tezos-client RPC format GET /chains/main/blocks/head
*/
type DelegationMetadata struct {
	BalanceUpdates           []BalanceUpdates           `json:"balance_updates"`
	OperationResults         OperationResultDelegation  `json:"operation_result"`
	InternalOperationResults []InternalOperationResults `json:"internal_operation_results,omitempty"`
}

func (d *Delegation) toContentsHelper() ContentsHelper {
	var metadata *ContentsHelperMetadata
	if d.Metadata != nil {
		metadata = &ContentsHelperMetadata{
			BalanceUpdates: d.Metadata.BalanceUpdates,
			OperationResults: &OperationResultsHelper{
				Status:      d.Metadata.OperationResults.Status,
				ConsumedGas: d.Metadata.OperationResults.ConsumedGas,
				Errors:      d.Metadata.OperationResults.Errors,
			},
			InternalOperationResult: d.Metadata.InternalOperationResults,
		}
	}

	return ContentsHelper{
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
See: tezos-client RPC format GET /chains/main/blocks/head
*/
type InternalOperationResults struct {
	Kind        string            `json:"kind"`
	Source      string            `json:"source"`
	Nonce       int               `json:"nonce"`
	Amount      int64             `json:"amount,string,omitempty"`
	PublicKey   string            `json:"public_key,omitempty"`
	Destination string            `json:"destination,omitempty"`
	Balance     int64             `json:"balance,string,omitempty"`
	Delegate    string            `json:"delegate,omitempty"`
	Script      ScriptedContracts `json:"script,omitempty"`
	Parameters  struct {
		Entrypoint string               `json:"entrypoint"`
		Value      *MichelineExpression `json:"value"`
	} `json:"paramaters,omitempty"`
	Result interface{} `json:"result"` //TODO This could be other things
}

/*
OperationResultReveal represents an OperationResultReveal in the $operation.alpha.operation_result.reveal in the tezos block schema
See: tezos-client RPC format GET /chains/main/blocks/head
*/
type OperationResultReveal struct {
	Status      string     `json:"status"`
	ConsumedGas int64      `json:"consumed_gas,omitempty,string"`
	Errors      []RPCError `json:"rpc_error,omitempty"`
}

/*
OperationResultTransfer represents $operation.alpha.operation_result.transaction in the tezos block schema
See: tezos-client RPC format GET /chains/main/blocks/head
*/
type OperationResultTransfer struct {
	Status                       string               `json:"status"`
	Storage                      *MichelineExpression `json:"storage,omitempty"`
	BigMapDiff                   *BigMapDiff          `json:"big_map_diff,omitempty"`
	BalanceUpdates               []BalanceUpdates     `json:"balance_updates,omitempty"`
	OriginatedContracts          []string             `json:"originated_contracts,omitempty"`
	ConsumedGas                  int64                `json:"consumed_gas,string,omitempty"`
	StorageSize                  int64                `json:"storage_size,string,omitempty"`
	PaidStorageSizeDiff          int64                `json:"paid_storage_size_diff,string,omitempty"`
	AllocatedDestinationContract bool                 `json:"allocated_destination_contract,omitempty"`
	Errors                       []RPCError           `json:"errors,omitempty"`
}

/*
OperationResultOrigination represents $operation.alpha.operation_result.origination in the tezos block schema
See: tezos-client RPC format GET /chains/main/blocks/head
*/
type OperationResultOrigination struct {
	Status              string           `json:"status"`
	BigMapDiff          *BigMapDiff      `json:"big_map_diff,omitempty"`
	BalanceUpdates      []BalanceUpdates `json:"balance_updates,omitempty"`
	OriginatedContracts []string         `json:"originated_contracts,omitempty"`
	ConsumedGas         int64            `json:"consumed_gas,string,omitempty"`
	StorageSize         int64            `json:"storage_size,string,omitempty"`
	PaidStorageSizeDiff int64            `json:"paid_storage_size_diff,string,omitempty"`
	Errors              []RPCError       `json:"errors,omitempty"`
}

/*
OperationResultDelegation represents $operation.alpha.operation_result.delegation in the tezos block schema
See: tezos-client RPC format GET /chains/main/blocks/head
*/
type OperationResultDelegation struct {
	Status      string     `json:"status"`
	ConsumedGas int64      `json:"consumed_gas,string,omitempty"`
	Errors      []RPCError `json:"errors,omitempty"`
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

// BigMapDiffHelper is a helper for unmarshaling and marshaling BigMapDiff
type BigMapDiffHelper struct {
	Action            string               `json:"action,omitempty"`
	BigMap            int                  `json:"big_map,string,omitempty"`
	KeyHash           string               `json:"key_hash,omitempty"`
	Key               *MichelineExpression `json:"key,omitempty"`
	Value             *MichelineExpression `json:"value,omitempty"`
	SourceBigMap      int                  `json:"source_big_map,string,omitempty"`
	DestinationBigMap int                  `json:"destination_big_map,string,omitempty"`
	KeyType           *MichelineExpression `json:"key_type,omitempty"`
	ValueType         *MichelineExpression `json:"value_type,omitempty"`
}

/*
UnmarshalJSON implements the json.UnmarshalJSON interface for BigMapDiff
*/
func (b *BigMapDiff) UnmarshalJSON(v []byte) error {
	var bigMapDiffHelpers []BigMapDiffHelper
	if err := json.Unmarshal(v, &bigMapDiffHelpers); err != nil {
		return errors.Wrap(err, "failed to unmarshal BigMapDiff")
	}

	var bigMapDiffUpdate []BigMapDiffUpdate
	if err := json.Unmarshal(v, &bigMapDiffUpdate); err != nil {
		return errors.Wrap(err, "failed to unmarshal BigMapDiff")
	}

	for _, bigMapDiffHelper := range bigMapDiffHelpers {
		if bigMapDiffHelper.Action == "update" {

			bigMapDiffUpdate := BigMapDiffUpdate{
				bigMapDiffHelper.Action,
				bigMapDiffHelper.BigMap,
				bigMapDiffHelper.KeyHash,
				bigMapDiffHelper.Key,
				bigMapDiffHelper.Value,
			}

			b.Updates = append(b.Updates, bigMapDiffUpdate)
		} else if bigMapDiffHelper.Action == "remove" {
			bigMapDiffRemove := BigMapDiffRemove{
				bigMapDiffHelper.Action,
				bigMapDiffHelper.BigMap,
			}

			b.Removals = append(b.Removals, bigMapDiffRemove)
		} else if bigMapDiffHelper.Action == "copy" {
			bigMapDiffCopy := BigMapDiffCopy{
				bigMapDiffHelper.Action,
				bigMapDiffHelper.SourceBigMap,
				bigMapDiffHelper.DestinationBigMap,
			}

			b.Copies = append(b.Copies, bigMapDiffCopy)
		} else if bigMapDiffHelper.Action == "alloc" {
			bigMapDiffAlloc := BigMapDiffAlloc{
				bigMapDiffHelper.Action,
				bigMapDiffHelper.BigMap,
				bigMapDiffHelper.KeyType,
				bigMapDiffHelper.ValueType,
			}

			b.Alloc = append(b.Alloc, bigMapDiffAlloc)
		}
	}

	return nil
}

/*
MarshalJSON implements the json.Marshaler interface for BigMapDiff
*/
func (b *BigMapDiff) MarshalJSON() ([]byte, error) {
	var bigMapDiffHelpers []BigMapDiffHelper
	for _, update := range b.Updates {
		bigMapDiffHelpers = append(bigMapDiffHelpers, BigMapDiffHelper{
			Action:  update.Action,
			BigMap:  update.BigMap,
			KeyHash: update.KeyHash,
			Key:     update.Key,
			Value:   update.Value,
		})
	}

	for _, remove := range b.Removals {
		bigMapDiffHelpers = append(bigMapDiffHelpers, BigMapDiffHelper{
			Action: remove.Action,
			BigMap: remove.BigMap,
		})
	}

	for _, copy := range b.Copies {
		bigMapDiffHelpers = append(bigMapDiffHelpers, BigMapDiffHelper{
			Action:            copy.Action,
			SourceBigMap:      copy.SourceBigMap,
			DestinationBigMap: copy.DestinationBigMap,
		})
	}

	for _, alloc := range b.Alloc {
		bigMapDiffHelpers = append(bigMapDiffHelpers, BigMapDiffHelper{
			Action:       alloc.Action,
			SourceBigMap: alloc.BigMap,
			KeyType:      alloc.KeyType,
			ValueType:    alloc.ValueType,
		})
	}

	v, err := json.Marshal(&bigMapDiffHelpers)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal BigMapDiff")
	}

	return v, nil
}

/*
BigMapDiffUpdate represents $contract.big_map_diff in the tezos block schema
See: tezos-client RPC format GET /chains/main/blocks/head
*/
type BigMapDiffUpdate struct {
	Action  string               `json:"action"`
	BigMap  int                  `json:"big_map,string,omitempty"`
	KeyHash string               `json:"key_hash,omitempty"`
	Key     *MichelineExpression `json:"key"`
	Value   *MichelineExpression `json:"value,omitempty"`
}

/*
BigMapDiffRemove represents $contract.big_map_diff in the tezos block schema
See: tezos-client RPC format GET /chains/main/blocks/head
*/
type BigMapDiffRemove struct {
	Action string `json:"action"`
	BigMap int    `json:"big_map,string"`
}

/*
BigMapDiffCopy represents $contract.big_map_diff in the tezos block schema
See: tezos-client RPC format GET /chains/main/blocks/head
*/
type BigMapDiffCopy struct {
	Action            string `json:"action"`
	SourceBigMap      int    `json:"source_big_map,string"`
	DestinationBigMap int    `json:"destination_big_map,string"`
}

/*
BigMapDiffAlloc represents $contract.big_map_diff in the tezos block schema
See: tezos-client RPC format GET /chains/main/blocks/head
*/
type BigMapDiffAlloc struct {
	Action    string               `json:"action"`
	BigMap    int                  `json:"big_map,string"`
	KeyType   *MichelineExpression `json:"key_type"`
	ValueType *MichelineExpression `json:"value_type"`
}

/*
ScriptedContracts represents $scripted.contracts in the tezos block schema
See: tezos-client RPC format GET /chains/main/blocks/head
*/
type ScriptedContracts struct {
	Code    MichelineExpression `json:"code"`
	Storage MichelineExpression `json:"storage"`
}

/*
Micheline represents $micheline.michelson_v1.expression in the tezos block schema
See: tezos-client RPC format GET /chains/main/blocks/head
*/
type Micheline struct {
	Int    string         `json:"int,omitempty"`
	String string         `json:"string,omitempty"`
	Bytes  string         `json:"bytes,omitempty"`
	Prim   string         `json:"prim,omitempty"`
	Args   *MichelineArgs `json:"args,omitempty"`
	Annots []string       `json:"annots,omitempty"`
}

/*
MichelineArgs represents the Args field for $micheline.michelson_v1.expression  in the
tezos block schema. Because the JSON for $micheline.michelson_v1.expression.args is dynamic
(can be slice or multi-slice) one of these fields will always be nil.
*/
type MichelineArgs struct {
	Array      []Micheline
	MultiArray [][]Micheline
}

/*
MichelineExpression represents $micheline.michelson_v1.expression in the tezos block schema.
Because the JSON for $micheline.michelson_v1.expression is dynamic (can be object or slice)
one of these fields will always be nil.
*/
type MichelineExpression struct {
	Object *Micheline
	Array  []Micheline
}

// UnmarshalJSON satisfies json.Marshaler
func (m *MichelineExpression) UnmarshalJSON(data []byte) error {
	var errs []error

	var object Micheline
	if err := json.Unmarshal(data, &object); err == nil {
		m.Object = &object
	} else {
		errs = append(errs, err)
	}

	var array []Micheline
	if err := json.Unmarshal(data, &array); err == nil {
		m.Array = array
	} else {
		errs = append(errs, err)
	}

	if m.Array == nil && m.Object == nil {
		return errors.Wrap(errs[1], errs[0].Error())
	}

	return nil
}

// MarshalJSON satisfies json.Marshaler
func (m *MichelineExpression) MarshalJSON() ([]byte, error) {
	if m.Object != nil {
		return json.Marshal(m.Object)
	}

	if m.Array != nil {
		return json.Marshal(m.Array)
	}

	return []byte{}, errors.New("nil values in MichelineExpression")
}

// UnmarshalJSON satisfies json.Marshaler
func (m *MichelineArgs) UnmarshalJSON(data []byte) error {
	var errs []error

	var array []Micheline
	if err := json.Unmarshal(data, &array); err == nil {
		m.Array = array
	} else {
		errs = append(errs, err)
	}

	var multiArray [][]Micheline
	if err := json.Unmarshal(data, &multiArray); err == nil {
		m.MultiArray = multiArray
	} else {
		errs = append(errs, err)
	}

	if m.Array == nil && m.MultiArray == nil {
		return errors.Wrap(errs[1], errs[0].Error())
	}

	return nil
}

// MarshalJSON satisfies json.Marshaler
func (m *MichelineArgs) MarshalJSON() ([]byte, error) {
	if m.Array != nil {
		return json.Marshal(m.Array)
	}

	if m.MultiArray != nil {
		return json.Marshal(m.MultiArray)
	}

	return []byte{}, errors.New("nil values in MichelineExpression")
}

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
