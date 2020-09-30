package gotezos

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	validator "github.com/go-playground/validator/v10"
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
	Status string `json:"status"`
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
	Data      string     `json:"data,omitempty"`
	Errors    []Error    `json:"errors,omitempty"`
}

/*
Contents represents the contents in a Tezos operations

RPC:
	/chains/<chain_id>/blocks/<block_id> (<dyn>)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-context-contracts-contract-id-balance
*/
type Contents struct {
	Kind             string            `json:"kind,omitempty"`
	Source           string            `json:"source,omitempty"`
	Fee              *Int              `json:"fee,omitempty"`
	Counter          *Int              `json:"counter,omitempty"`
	GasLimit         *Int              `json:"gas_limit,omitempty"`
	StorageLimit     *Int              `json:"storage_limit,omitempty"`
	Amount           *Int              `json:"amount,omitempty"`
	Destination      string            `json:"destination,omitempty"`
	Delegate         string            `json:"delegate,omitempty"`
	Phk              string            `json:"phk,omitempty"`
	Secret           string            `json:"secret,omitempty"`
	Level            int               `json:"level,omitempty"`
	ManagerPublicKey string            `json:"managerPubkey,omitempty"`
	Balance          *Int              `json:"balance,omitempty"`
	Period           int               `json:"period,omitempty"`
	Proposal         string            `json:"proposal,omitempty"`
	Proposals        []string          `json:"proposals,omitempty"`
	Ballot           string            `json:"ballot,omitempty"`
	Metadata         *ContentsMetadata `json:"metadata,omitempty"`
	Nonce            string            `json:"nonce,omitempty"`  // Revealing nonce for block
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
type InternalOperationResults struct {
	Kind        string           `json:"kind"`
	Source      string           `json:"source"`
	Nonce       uint64           `json:"nonce"`
	Amount      string           `json:"amount"`
	Destination string           `json:"destination"`
	Result      *OperationResult `json:"result"`
}

/*
Error respresents an error for operation results

RPC:
	/chains/<chain_id>/blocks/<block_id> (<dyn>)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-context-contracts-contract-id-balance
*/
type Error struct {
	Kind     string `json:"kind"`
	ID       string `json:"id"`
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
ProtocolData is information required in any newly forged Blocks

*/
type ProtocolData struct {
	Protocol      string `json:"protocol"`
	Priority      int    `json:"priority"`
	POWNonce      string `json:"proof_of_work_nonce"`
	Signature     string `json:"signature"`
	SeedNonceHash string `json:"seed_nonce_hash,omitempty"`
}

/*
PreapplyBlockOperationInput is the input for the goTezos.PreapplyBlockOperation function.

Function:
    func (t *GoTezos) PreapplyBlockOperation(input PreapplyBlockOperationInput) ([]interface{}, error) {}
*/
type PreapplyBlockOperationInput struct {
	// Header data for the new block
	Protocoldata ProtocolData `validate:"required"`

	// Operations (txn, endorsements, etc) to be contained in the next block
	Operations [][]Operations `validate:"required"`

	// Sort operations on return
	Sort bool `validate:"required"`

	// Earliest timestamp of next block
	Timestamp time.Time `validate:"required"`
}

func (p *PreapplyBlockOperationInput) constructRPCOptions() []rpcOptions {
    var opts []rpcOptions
    if p.Sort {
        opts = append(opts, rpcOptions{
            "sort",
            "true",
        })
    }

    if !p.Timestamp.IsZero() {
        opts = append(opts, rpcOptions{
            "timestamp",
            strconv.FormatInt(p.Timestamp.Unix(), 10),
        })
    }

    return opts
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

/*
PreapplyBlock returns a succesful baked block

Path:
    /chains/<chain_id>/blocks/head/helpers/preapply/block

Parameters:

    sort:
        Sorts the resulting operations
    timestamp:
        ??
*/
type ShellHeader struct {
	Level            int       `json:"level"`
	Proto            int       `json:"proto"`
	Predecessor      string    `json:"predecessor"`
	Timestamp        time.Time `json:"timestamp"`
	ValidationPass   int       `json:"validation_pass"`
	OperationsHash   string    `json:"operations_hash"`
	Fitness          []string  `json:"fitness"`
	Context          string    `json:"context"`
	ProtocolData     string    `json:"protocol_data"`
}

type PreApplyOperationsAlt PreApplyOperations

func (o *PreApplyOperationsAlt) UnmarshalJSON(buf []byte) error {
	return unmarshalNamedJSONArray(buf, &o.Hash, (*PreApplyOperations)(o))
}

type ShellOperations struct {
	Applied       []PreApplyOperations    `json:"applied"`
	Refused       []PreApplyOperationsAlt `json:"refused"`
	BranchRefused []PreApplyOperationsAlt `json:"branch_refused"`
	BranchDelayed []PreApplyOperationsAlt `json:"branch_delayed"`
	Unprocessed   []PreApplyOperationsAlt `json:"unprocessed"`
}

type PreApplyOperations struct {
	Hash   string `json:"hash"`
	Branch string `json:"branch"`
	Data   string `json:"data"`
}

type PreapplyResult struct {
	Shellheader ShellHeader `json:"shell_header"`
	Operations []ShellOperations `json:"operations"`
}

func (t *GoTezos) PreapplyBlockOperation(input PreapplyBlockOperationInput) (PreapplyResult, error) {
	var preapplyRes PreapplyResult
	err  := validator.New().Struct(input)
	if err != nil {
		return preapplyRes, errors.Wrap(err, "invalid input")
	}

	// construct the block
	bh := struct {
		Pd ProtocolData `json:"protocol_data"`
		Ops [][]Operations `json:"operations"`
	}{
		input.Protocoldata,
		input.Operations,
	}

	v, err := json.Marshal(bh)
	if err != nil {
		return preapplyRes, errors.Wrap(err, "failed to marshal new block")
	}

	resp, err := t.post("/chains/main/blocks/head/helpers/preapply/block", v, input.constructRPCOptions()...)
	if err != nil {
		return preapplyRes, errors.Wrap(err, "failed to preapply new block")
	}

	// Parse the response from preapply
	err = json.Unmarshal(resp, &preapplyRes)
	if err != nil {
		return preapplyRes, errors.Wrap(err, "failed to unmarshal preapply result of new block")
	}

	return preapplyRes, nil
}

type ForgedBlockHeader struct {
	BlockHeader string `json:"block"`
}

func (t *GoTezos) ForgeBlockHeader(sh ShellHeader) (ForgedBlockHeader, error) {

	var fBlockheader ForgedBlockHeader

	v, err := json.Marshal(sh)
	if err != nil {
		return fBlockheader, errors.Wrap(err, "failed to marshal shell header for forge")
	}

	resp, err := t.post("/chains/main/blocks/head/helpers/forge_block_header", v)
	if err != nil {
		return fBlockheader, errors.Wrap(err, "failed to forge new block")
	}
	
	// Parse the response from preapply
	err = json.Unmarshal(resp, &fBlockheader)
	if err != nil {
		return fBlockheader, errors.Wrap(err, "failed to unmarshal forged block header result")
	}

	return fBlockheader, nil
}

type EndorsingPowerInput struct {
	Endorsement  Operations `json:"endorsement_operation"`
	ChainID      string     `json:"chain_id"`
}

func (t *GoTezos) GetEndorsingPower(operation EndorsingPowerInput) (int, error) {

	// Remove Protocol value
	operation.Endorsement.Protocol = ""

	// Marshal
	v, err := json.Marshal(operation)
	if err != nil {
		return 0, errors.Wrap(err, "failed to marshal operation for endorsing power")
	}

	resp, err := t.post(fmt.Sprintf("/chains/%s/blocks/head/endorsing_power", operation.ChainID), v)
	if err != nil {
		return 0, errors.Wrap(err, "failed to fetch endorsing power")
	}

	// RPC contains newline; strip in order to parse to int
	power, err := strconv.Atoi(strings.TrimSuffix(string(resp), "\n"))
	if err != nil {
		return 0, errors.Wrap(err, "failed to decode endorsing power")
	}

	return power, nil
}

func (t *GoTezos) MinimalValidTime(endorsingPower, priority int, chainID string) (time.Time, error) {

	var minTimestamp time.Time

	var opts []rpcOptions
	opts = append(opts, rpcOptions{
		"priority", strconv.Itoa(priority),
	})
	opts = append(opts, rpcOptions{
		"endorsing_power", strconv.Itoa(endorsingPower),
	})

	resp, err := t.get(fmt.Sprintf("/chains/%s/blocks/head/minimal_valid_time", chainID), opts...)
	if err != nil {
		return minTimestamp, errors.Wrap(err, "failed to fetch minimal_valid_time")
	}

	minTimestamp, err = time.Parse(time.RFC3339, stripQuote(string(resp)))
	if err != nil {
		return minTimestamp, errors.Wrap(err, "failed to parse minimal valid timestamp")
	}

	return minTimestamp, nil
}

func stripQuote(s string) string {
	m := strings.TrimSpace(s)
	if len(m) > 0 && m[0] == '"' {
		m = m[1:]
	}
	if len(m) > 0 && m[len(m)-1] == '"' {
		m = m[:len(m)-1]
	}
	return m
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
