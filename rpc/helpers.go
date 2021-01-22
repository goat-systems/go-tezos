package rpc

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"

	validator "github.com/go-playground/validator/v10"
	"github.com/go-resty/resty/v2"
	"github.com/goat-systems/go-tezos/v4/internal/crypto"
	"github.com/pkg/errors"
)

/*
BakingRights represents the baking rights RPC on the tezos network.

Path:
	../<block_id>/helpers/baking_rights?(level=<block_level>)*&(cycle=<block_cycle>)*&(delegate=<pkh>)*&[max_priority=<int>]&[all]

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-helpers-baking-rights
*/
type BakingRights struct {
	Level         int       `json:"level"`
	Delegate      string    `json:"delegate"`
	Priority      int       `json:"priority"`
	EstimatedTime time.Time `json:"estimated_time"`
}

/*
BakingRightsInput is the input for the BakingRights function.

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-helpers-baking-rights
*/
type BakingRightsInput struct {
	// The hash of block (height) of which you want to make the query.
	BlockID BlockID `validate:"required"`
	// The cycle of which you want to make the query.
	Cycle int
	// The block level of which you want to make the query.
	Level int
	// The delegate public key hash of which you want to make the query.
	Delegate string
	// The max priotity of which you want to make the query.
	MaxPriority int
	// All baking rights
	All bool
}

func (b *BakingRightsInput) contructRPCOptions() []rpcOptions {
	var opts []rpcOptions
	if b.Cycle != 0 {
		opts = append(opts, rpcOptions{
			"cycle",
			strconv.Itoa(b.Cycle),
		})
	}

	if b.Delegate != "" {
		opts = append(opts, rpcOptions{
			"delegate",
			b.Delegate,
		})
	}

	if b.Level != 0 {
		opts = append(opts, rpcOptions{
			"level",
			strconv.Itoa(b.Level),
		})
	}

	if b.MaxPriority != 0 {
		opts = append(opts, rpcOptions{
			"max_priority",
			strconv.Itoa(b.MaxPriority),
		})
	}

	if b.All {
		opts = append(opts, rpcOptions{
			"all",
			"True",
		})
	}

	return opts
}

/*
BakingRights retrieves the list of delegates allowed to bake a block. By default, it gives the best baking
priorities for bakers that have at least one opportunity below the 64th priority for the next block. Parameters
`level` and `cycle` can be used to specify the (valid) level(s) in the past or future at which the baking
rights have to be returned. When asked for (a) whole cycle(s), baking opportunities are given by default up to
the priority 8. Parameter `delegate` can be used to restrict the results to the given delegates. If parameter
`all` is set, all the baking opportunities for each baker at each level are returned, instead of just the first
one. Returns the list of baking slots. Also returns the minimal timestamps that correspond to these slots. The
timestamps are omitted for levels in the past, and are only estimates for levels later that the next block, based
on the hypothesis that all predecessor blocks were baked at the first priority.

Path:
	../<block_id>/helpers/baking_rights?(level=<block_level>)*&(cycle=<block_cycle>)*&(delegate=<pkh>)*&[max_priority=<int>]&[all] (GET)

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-helpers-baking-rights
*/
func (c *Client) BakingRights(input BakingRightsInput) (*resty.Response, []BakingRights, error) {
	err := validator.New().Struct(input)
	if err != nil {
		return nil, []BakingRights{}, errors.Wrap(err, "failed to get baking rights: invalid input")
	}

	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/helpers/baking_rights", c.chain, input.BlockID.ID()), input.contructRPCOptions()...)
	if err != nil {
		return resp, []BakingRights{}, errors.Wrapf(err, "failed to get baking rights")
	}

	var bakingRights []BakingRights
	err = json.Unmarshal(resp.Body(), &bakingRights)
	if err != nil {
		return resp, []BakingRights{}, errors.Wrapf(err, "failed to get baking rights: failed to parse json")
	}

	return resp, bakingRights, nil
}

/*
CompletePrefixInput is the input for the Prefix function.

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-helpers-complete-prefix
*/
type CompletePrefixInput struct {
	// The hash of block (height) of which you want to make the query.
	BlockID BlockID `validate:"required"`
	// The prefix of you wish to complete
	Prefix string
}

/*
CompletePrefix tries to complete a prefix of a Base58Check-encoded data.
This RPC is actually able to complete hashes of block, operations, public_keys and contracts.

Path:
	../<block_id>/helpers/complete/<prefix>
RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-helpers-complete-prefix
*/
func (c *Client) CompletePrefix(input CompletePrefixInput) (*resty.Response, []string, error) {
	err := validator.New().Struct(input)
	if err != nil {
		return nil, []string{}, errors.Wrap(err, "failed to complete prefix: invalid input")
	}

	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/helpers/complete/%s", c.chain, input.BlockID.ID(), input.Prefix))
	if err != nil {
		return resp, []string{}, errors.Wrapf(err, "failed to complete prefix")
	}

	var completion []string
	err = json.Unmarshal(resp.Body(), &completion)
	if err != nil {
		return resp, []string{}, errors.Wrapf(err, "failed to complete prefix: failed to parse json")
	}

	return resp, completion, nil
}

/*
CurrentLevelInput is the input to the CurrentLevel function.

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-helpers-current-level
*/
type CurrentLevelInput struct {
	// The hash of block (height) of which you want to make the query.
	BlockID BlockID `validate:"required"`
	// The next block if `offset` is 1.
	Offset int32
}

func (c *CurrentLevelInput) constructRPCOptions() []rpcOptions {
	var opts []rpcOptions
	if c.Offset != 0 {
		opts = append(opts, rpcOptions{
			"offset",
			fmt.Sprintf("%d", c.Offset),
		})
	}

	return opts
}

/*
CurrentLevel is the the level of the interrogated block, or the one of a block located
`offset` blocks after in the chain (or before when negative). For instance, the next block
if `offset` is 1.

Path:
	../<block_id>/helpers/current_level?[offset=<int32>]

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-helpers-current-level
*/
type CurrentLevel struct {
	Level                int  `json:"level"`
	LevelPosition        int  `json:"level_position"`
	Cycle                int  `json:"cycle"`
	CyclePosition        int  `json:"cycle_position"`
	VotingPeriod         int  `json:"voting_period"`
	VotingPeriodPosition int  `json:"voting_period_position"`
	ExpectedCommitment   bool `json:"expected_commitment"`
}

/*
CurrentLevel returns the level of the interrogated block, or the one of a block located
`offset` blocks after in the chain (or before when negative). For instance, the next block
if `offset` is 1.

Path:
	../<block_id>/helpers/current_level?[offset=<int32>]
RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-helpers-current-level
*/
func (c *Client) CurrentLevel(input CurrentLevelInput) (*resty.Response, CurrentLevel, error) {
	err := validator.New().Struct(input)
	if err != nil {
		return nil, CurrentLevel{}, errors.Wrap(err, "failed to get current level: invalid input")
	}

	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/helpers/current_level", c.chain, input.BlockID.ID()), input.constructRPCOptions()...)
	if err != nil {
		return resp, CurrentLevel{}, errors.Wrapf(err, "failed to get current level")
	}

	var currentLevel CurrentLevel
	err = json.Unmarshal(resp.Body(), &currentLevel)
	if err != nil {
		return resp, CurrentLevel{}, errors.Wrapf(err, "failed to get current level: failed to parse json")
	}

	return resp, currentLevel, nil
}

/*
EndorsingRightsInput is the input for the EndorsingRights function.

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-helpers-endorsing-rights
*/
type EndorsingRightsInput struct {
	// The block (height) of which you want to make the query.
	BlockID BlockID `validate:"required"`
	// The block level of which you want to make the query.
	Level int
	// The cycle of which you want to make the query.
	Cycle int
	// The delegate public key hash of which you want to make the query.
	Delegate string
}

/*
EndorsingRights represents the endorsing rights RPC on the tezos network.

Path:
	../<block_id>/helpers/endorsing_rights?(level=<block_level>)*&(cycle=<block_cycle>)*&(delegate=<pkh>)* (GET)

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-helpers-endorsing-rights
*/
type EndorsingRights struct {
	Level         int       `json:"level"`
	Delegate      string    `json:"delegate"`
	Slots         []int     `json:"slots"`
	EstimatedTime time.Time `json:"estimated_time"`
}

/*
EndorsingRights retrieves the delegates allowed to endorse a block. By default,
it gives the endorsement slots for delegates that have at least one in the
next block. Parameters `level` and `cycle` can be used to specify the (valid)
level(s) in the past or future at which the endorsement rights have to be returned.
Parameter `delegate` can be used to restrict the results to the given delegates.
Returns the list of endorsement slots. Also returns the minimal timestamps that
correspond to these slots. The timestamps are omitted for levels in the past, and
are only estimates for levels later that the next block, based on the hypothesis
that all predecessor blocks were baked at the first priority.

Path:
	../<block_id>/helpers/endorsing_rights?(level=<block_level>)*&(cycle=<block_cycle>)*&(delegate=<pkh>)* (GET)

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-helpers-endorsing-rights
*/
func (c *Client) EndorsingRights(input EndorsingRightsInput) (*resty.Response, []EndorsingRights, error) {
	err := validator.New().Struct(input)
	if err != nil {
		return nil, []EndorsingRights{}, errors.Wrap(err, "failed to get endorsing rightsL invalid input")
	}

	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/helpers/endorsing_rights", c.chain, input.BlockID.ID()), input.contructRPCOptions()...)
	if err != nil {
		return resp, []EndorsingRights{}, errors.Wrap(err, "failed to get endorsing rights")
	}

	var endorsingRights []EndorsingRights
	err = json.Unmarshal(resp.Body(), &endorsingRights)
	if err != nil {
		return resp, []EndorsingRights{}, errors.Wrapf(err, "failed to get endorsing rights: failed to parse json")
	}

	return resp, endorsingRights, nil
}

func (b *EndorsingRightsInput) contructRPCOptions() []rpcOptions {
	var opts []rpcOptions
	if b.Cycle != 0 {
		opts = append(opts, rpcOptions{
			"cycle",
			strconv.Itoa(b.Cycle),
		})
	}

	if b.Delegate != "" {
		opts = append(opts, rpcOptions{
			"delegate",
			b.Delegate,
		})
	}

	if b.Level != 0 {
		opts = append(opts, rpcOptions{
			"level",
			strconv.Itoa(b.Level),
		})
	}

	return opts
}

/*
ForgeOperationsInput is the input for the function ForgeOperation

RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-helpers-forge-operations

*/
type ForgeOperationsInput struct {
	// The hash of block (height) of which you want to make the query.
	BlockIDHash BlockIDHash `validate:"required"`
	Branch      string      `validate:"required"`
	Contents    Contents    `validate:"required"`
	// Using the RPC to forge an operation is dangerous, you can mitigate this
	// danger by passing a different host to CheckRPCAddr which will unforge the
	// operation and compare the results to filter something malicious.
	// OR just use the go-tezos/forge package for forging locally.
	CheckRPCAddr string
}

/*
ForgeOperations forges an operation.

Path:
	../<block_id>/helpers/forge/operations (POST)

RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-helpers-forge-operations

*/
func (c *Client) ForgeOperations(input ForgeOperationsInput) (*resty.Response, string, error) {
	err := validator.New().Struct(input)
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to forge operation: invalid input")
	}

	op := Operations{
		Branch:   input.Branch,
		Contents: input.Contents,
	}

	v, err := json.Marshal(op)
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to forge operation")
	}

	resp, err := c.post(fmt.Sprintf("/chains/%s/blocks/%s/helpers/forge/operations", c.chain, input.BlockIDHash.ID()), v)
	if err != nil {
		return resp, "", errors.Wrap(err, "failed to forge operation")
	}

	var operation string
	err = json.Unmarshal(resp.Body(), &operation)
	if err != nil {
		return resp, "", errors.Wrap(err, "failed to forge operation")
	}

	_, opstr, err := stripBranchFromForgedOperation(operation, false)
	if err != nil {
		return resp, operation, errors.Wrap(err, "failed to forge operation: unable to verify rpc returned a valid contents")
	}

	var rpc *Client
	if input.CheckRPCAddr != "" {
		rpc, err = New(input.CheckRPCAddr)
		if err != nil {
			return resp, operation, errors.Wrap(err, "failed to forge operation: unable to verify rpc returned a valid contents with alternative node")
		}
	} else {
		rpc = c
	}

	resp, operations, err := rpc.ParseOperations(ParseOperationsInput{
		BlockID: &input.BlockIDHash,
		Operations: []ParseOperationsBody{
			{
				Data:   fmt.Sprintf("%s00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", opstr),
				Branch: input.Branch,
			},
		},
		CheckSignature: false,
	})
	if err != nil {
		return resp, operation, errors.Wrap(err, "failed to forge operation: unable to verify rpc returned a valid contents")
	}

	for _, op := range operations {
		ok := reflect.DeepEqual(op.Contents, input.Contents)
		if !ok {
			return resp, operation, errors.New("failed to forge operation: alert rpc returned invalid contents")
		}
	}

	return resp, operation, nil
}

/*
ForgeBlockHeaderInput is the input for the function ForgeBlockHeader

RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-helpers-forge-block-header

*/
type ForgeBlockHeaderInput struct {
	// The block (height) of which you want to make the query.
	BlockID BlockID `validate:"required"`
	// The block header you wish to forge
	BlockHeader ForgeBlockHeaderBody `validate:"required"`
}

/*
ForgeBlockHeaderBody is the block header to forge

RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-helpers-forge-block-header

*/
type ForgeBlockHeaderBody struct {
	Level          int       `json:"level"`
	Proto          int       `json:"proto"`
	Predecessor    string    `json:"predecessor"`
	Timestamp      time.Time `json:"timestamp"`
	ValidationPass int       `json:"validation_pass"`
	OperationsHash string    `json:"operations_hash"`
	Fitness        []string  `json:"fitness"`
	Context        string    `json:"context"`
	ProtocolData   string    `json:"protocol_data"`
}

/*
ForgeBlockHeader is the block header received from forging

RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-helpers-forge-block-header
*/
type ForgeBlockHeader struct {
	Block string `json:"block"`
}

/*
ForgeBlockHeader is the block header received from forging

Path:
	../<block_id>/helpers/forge_block_header (POST)

RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-helpers-forge-block-header
*/
func (c *Client) ForgeBlockHeader(input ForgeBlockHeaderInput) (*resty.Response, ForgeBlockHeader, error) {
	err := validator.New().Struct(input)
	if err != nil {
		return nil, ForgeBlockHeader{}, errors.Wrap(err, "failed to forge block header: invalid input")
	}

	resp, err := c.post(fmt.Sprintf("/chains/%s/blocks/%s/helpers/forge_block_header", c.chain, input.BlockID.ID()), input.BlockHeader)
	if err != nil {
		return resp, ForgeBlockHeader{}, errors.Wrap(err, "failed to forge block header")
	}

	var blockHeader ForgeBlockHeader
	err = json.Unmarshal(resp.Body(), &blockHeader)
	if err != nil {
		return resp, ForgeBlockHeader{}, errors.Wrap(err, "failed to forge block header: failed to parse json")
	}

	return resp, blockHeader, nil
}

/*
LevelsInCurrentCycleInput is the input for the LevelsInCurrentCycle function

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-helpers-levels-in-current-cycle
*/
type LevelsInCurrentCycleInput struct {
	// The block (height) of which you want to make the query.
	BlockID BlockID `validate:"required"`
	Offset  int32
}

/*
LevelsInCurrentCycle is the levels of a cycle

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-helpers-levels-in-current-cycle
*/
type LevelsInCurrentCycle struct {
	First int `json:"first"`
	Last  int `json:"last"`
}

/*
LevelsInCurrentCycle is the levels of a cycle

Path:
	../<block_id>/helpers/levels_in_current_cycle?[offset=<int32>] (GET)

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-helpers-levels-in-current-cycle
*/
func (c *Client) LevelsInCurrentCycle(input LevelsInCurrentCycleInput) (*resty.Response, LevelsInCurrentCycle, error) {
	err := validator.New().Struct(input)
	if err != nil {
		return nil, LevelsInCurrentCycle{}, errors.Wrap(err, "failed to get levels in current cycle: invalid input")
	}

	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/helpers/levels_in_current_cycle", c.chain, input.BlockID.ID()))
	if err != nil {
		return resp, LevelsInCurrentCycle{}, errors.Wrap(err, "failed to get levels in current cycle")
	}

	var levelsInCurrentCycle LevelsInCurrentCycle
	err = json.Unmarshal(resp.Body(), &levelsInCurrentCycle)
	if err != nil {
		return resp, LevelsInCurrentCycle{}, errors.Wrap(err, "failed to get levels in current cycle: failed to parse json")
	}

	return resp, levelsInCurrentCycle, nil
}

/*
ParseBlockInput is the input for the function ParseBlock function

RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-helpers-parse-block
*/
type ParseBlockInput struct {
	// The block (height) of which you want to make the query.
	BlockID BlockID `validate:"required"`
	// The block header you wish to forge
	BlockHeader ForgeBlockHeaderBody `validate:"required"`
}

/*
BlockHeaderSignedContents is signed header contents returend from parsing a block

RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-helpers-parse-block
*/
type BlockHeaderSignedContents struct {
	Priority         int    `json:"priority"`
	ProofOfWorkNonce string `json:"proof_of_work_nonce"`
	SeedNonceHash    string `json:"seed_nonce_hash"`
	Signature        string `json:"signature"`
}

/*
ParseBlock is signed header contents returend from parsing a block

Path:
	../<block_id>/helpers/parse/block (POST)
RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-helpers-parse-block
*/
func (c *Client) ParseBlock(input ParseBlockInput) (*resty.Response, BlockHeaderSignedContents, error) {
	err := validator.New().Struct(input)
	if err != nil {
		return nil, BlockHeaderSignedContents{}, errors.Wrap(err, "failed to parse block: invalid input")
	}

	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/helpers/parse/block", c.chain, input.BlockID.ID()))
	if err != nil {
		return resp, BlockHeaderSignedContents{}, errors.Wrap(err, "failed to parse block")
	}

	var blockHeaderSignedContents BlockHeaderSignedContents
	err = json.Unmarshal(resp.Body(), &blockHeaderSignedContents)
	if err != nil {
		return resp, BlockHeaderSignedContents{}, errors.Wrap(err, "failed to parse block: failed to parse json")
	}

	return resp, blockHeaderSignedContents, nil
}

/*
ParseOperationsInput is the input for the ParseOperations function

RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-helpers-parse-operations
*/
type ParseOperationsInput struct {
	// The block (height) of which you want to make the query.
	BlockID BlockID `validate:"required"`
	// The operations to parse
	Operations []ParseOperationsBody `validate:"required"`
	// Whether to check the signature or not
	CheckSignature bool
}

/*
ParseOperationsBody is the operations you wish to parse

RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-helpers-parse-operations
*/
type ParseOperationsBody struct {
	Branch string `json:"branch"`
	Data   string `json:"data"`
}

/*
ParseOperations parses encoded operations to a slice of Operations

Path:
	../<block_id>/helpers/parse/operations
RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-helpers-parse-operations
*/
func (c *Client) ParseOperations(input ParseOperationsInput) (*resty.Response, []Operations, error) {
	err := validator.New().Struct(input)
	if err != nil {
		return nil, []Operations{}, errors.Wrap(err, "failed to parse operations: invalid input")
	}

	operations := struct {
		Operations     []ParseOperationsBody `json:"operations"`
		CheckSignature bool                  `json:"check_signature"`
	}{
		input.Operations,
		input.CheckSignature,
	}

	resp, err := c.post(fmt.Sprintf("/chains/%s/blocks/%s/helpers/parse/operations", c.chain, input.BlockID.ID()), operations)
	if err != nil {
		return resp, []Operations{}, errors.Wrap(err, "failed to parse operations")
	}

	var ops []Operations
	err = json.Unmarshal(resp.Body(), &ops)
	if err != nil {
		return resp, []Operations{}, errors.Wrap(err, "failed to parse operations: failed to parse json")
	}

	return resp, ops, nil
}

/*
PreapplyBlockInput is the input for the PreapplyBlock function

RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-helpers-preapply-block
*/
type PreapplyBlockInput struct {
	// The block (height) of which you want to make the query.
	BlockID BlockID `validate:"required"`
	// The block to preapply
	Block     PreapplyBlockBody `validate:"required"`
	Sort      bool
	Timestamp *time.Time
}

func (p *PreapplyBlockInput) constructRPCOptions() []rpcOptions {
	var options []rpcOptions
	if p.Sort {
		options = append(options, rpcOptions{
			"sort",
			"true",
		})
	}

	if p.Timestamp != nil {
		options = append(options, rpcOptions{
			"timestamp",
			strconv.FormatInt(p.Timestamp.Unix(), 10),
		})
	}

	return options
}

/*
PreapplyBlockBody is the block to preapply in the PreapplyBlock function

RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-helpers-preapply-block
*/
type PreapplyBlockBody struct {
	ProtocolData PreapplyBlockProtocolData `json:"protocol_data"`
	Operations   [][]Operations            `json:"operations"`
}

/*
PreapplyBlockProtocolData is the protocol data of the block to preapply in the PreapplyBlock function

RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-helpers-preapply-block
*/
type PreapplyBlockProtocolData struct {
	Protocol         string `json:"protocol"`
	Priority         int    `json:"priority"`
	ProofOfWorkNonce string `json:"proof_of_work_nonce"`
	SeedNonceHash    string `json:"seed_nonce_hash,omitempty"`
	Signature        string `json:"signature"`
}

/*
PreappliedBlock is the preapplied block returned by the PreapplyBlock function

RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-helpers-preapply-block
*/
type PreappliedBlock struct {
	ShellHeader HeaderShell                 `json:"shell_header"`
	Operations  []PreappliedBlockOperations `json:"operations"`
}

/*
PreappliedBlockOperations is the preapplied block operations returned by the PreapplyBlock function

RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-helpers-preapply-block
*/
type PreappliedBlockOperations struct {
	Applied       []PreappliedBlockOperationsStatus `json:"applied"`
	Refused       []PreappliedBlockOperationsStatus `json:"refused"`
	BranchRefused []PreappliedBlockOperationsStatus `json:"branch_refused"`
	BranchDelayed []PreappliedBlockOperationsStatus `json:"branch_delayed"`
}

/*
PreappliedBlockOperationsStatus is the preapplied block operation status returned by the PreapplyBlock function

RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-helpers-preapply-block
*/
type PreappliedBlockOperationsStatus struct {
	Hash   string      `json:"hash"`
	Branch string      `json:"branch"`
	Data   string      `json:"data"`
	Error  ResultError `json:"error,omitempty"`
}

/*
PreapplyBlock simulates the validation of a block that would contain
the given operations and return the resulting fitness and context hash.

Path:
	../<block_id>/helpers/preapply/block?[sort]&[timestamp=<date>]
RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-helpers-preapply-block
*/
func (c *Client) PreapplyBlock(input PreapplyBlockInput) (*resty.Response, PreappliedBlock, error) {
	err := validator.New().Struct(input)
	if err != nil {
		return nil, PreappliedBlock{}, errors.Wrap(err, "failed to preapply block: invalid input")
	}

	resp, err := c.post(fmt.Sprintf("/chains/%s/blocks/%s/helpers/preapply/block", c.chain, input.BlockID.ID()), input.Block, input.constructRPCOptions()...)
	if err != nil {
		return resp, PreappliedBlock{}, errors.Wrap(err, "failed to preapply block")
	}

	var preappliedBlock PreappliedBlock
	err = json.Unmarshal(resp.Body(), &preappliedBlock)
	if err != nil {
		return resp, PreappliedBlock{}, errors.Wrap(err, "failed to preapply block: failed to parse json")
	}

	return resp, preappliedBlock, nil
}

/*
PreapplyOperationsInput is the input for the PreapplyOperations function

RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-helpers-preapply-operations
*/
type PreapplyOperationsInput struct {
	// The block (height) of which you want to make the query.
	BlockID BlockID `validate:"required"`
	// The operations to parse
	Operations []Operations `validate:"required"`
}

/*
PreapplyOperations simulates the validation of an operation.

Path:
	../<block_id>/helpers/preapply/operations (POST)

RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-helpers-preapply-operations
*/
func (c *Client) PreapplyOperations(input PreapplyOperationsInput) (*resty.Response, []Operations, error) {
	err := validator.New().Struct(input)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to preapply operations: invalid input")
	}

	resp, err := c.post(fmt.Sprintf("/chains/%s/blocks/%s/helpers/preapply/operations", c.chain, input.BlockID.ID()), input.Operations)
	if err != nil {
		return resp, nil, errors.Wrap(err, "failed to preapply operations")
	}

	var operations []Operations
	err = json.Unmarshal(resp.Body(), &operations)
	if err != nil {
		return resp, nil, errors.Wrap(err, "failed to preapply operations: failed to parse json")
	}

	return resp, operations, nil
}

/*
EntrypointInput is the input for the Entrypoint function

RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-helpers-scripts-entrypoint
*/
type EntrypointInput struct {
	// The block (height) of which you want to make the query.
	BlockID BlockID `validate:"required"`
	// The entrypoint to get the type of
	Entrypoint EntrypointBody `validate:"required"`
}

/*
EntrypointBody is the entrypoint body for the Entrypoint function

RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-helpers-scripts-entrypoint
*/
type EntrypointBody struct {
	Script     *json.RawMessage `json:"script"`
	Entrypoint string           `json:"entrypoint,omitempty"`
}

/*
Entrypoint is the return value for the Entrypoint function and contains the entrypoint type

RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-helpers-scripts-entrypoint
*/
type Entrypoint struct {
	EntrypointType *json.RawMessage `json:"entrypoint_type"`
}

/*
Entrypoint returns the type of the given entrypoint.

Path:
	../<block_id>/helpers/scripts/entrypoint (POST)

RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-helpers-scripts-entrypoint
*/
func (c *Client) Entrypoint(input EntrypointInput) (*resty.Response, Entrypoint, error) {
	err := validator.New().Struct(input)
	if err != nil {
		return nil, Entrypoint{}, errors.Wrap(err, "failed to get entrypoint type: invalid input")
	}

	resp, err := c.post(fmt.Sprintf("/chains/%s/blocks/%s/helpers/scripts/entrypoint", c.chain, input.BlockID.ID()), input.Entrypoint)
	if err != nil {
		return resp, Entrypoint{}, errors.Wrap(err, "failed to get entrypoint type")
	}

	var entrypoint Entrypoint
	err = json.Unmarshal(resp.Body(), &entrypoint)
	if err != nil {
		return resp, Entrypoint{}, errors.Wrap(err, "failed to get entrypoint type: failed to parse json")
	}

	return resp, entrypoint, nil
}

/*
EntrypointsInput is the input for the Entrypoints function

RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-helpers-scripts-entrypoints
*/
type EntrypointsInput struct {
	// The block (height) of which you want to make the query.
	BlockID BlockID `validate:"required"`
	// The script to get the entrypoints for
	Entrypoints EntrypointsBody `validate:"required"`
}

/*
EntrypointsBody is the entrypoints body for the Entrypoints function

RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-helpers-scripts-entrypoints
*/
type EntrypointsBody struct {
	Script *json.RawMessage `json:"script"`
}

/*
Entrypoints is the return value for the Entrypoints function and contains the entrypoints for a script

RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-helpers-scripts-entrypoints
*/
type Entrypoints struct {
	Unreachable           []UnreachableEntrypoints `json:"unreachable,omitempty"`
	EntrypointsFromScript *json.RawMessage         `json:"entrypoints"`
}

/*
UnreachableEntrypoints is the unreachable entrypoints in theEntrypoints function

RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-helpers-scripts-entrypoints
*/
type UnreachableEntrypoints struct {
	Path []*json.RawMessage `json:"path"`
}

/*
Entrypoints returns the list of entrypoints of the given script

Path:
	../<block_id>/helpers/scripts/entrypoints (POST)

RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-helpers-scripts-entrypoints
*/
func (c *Client) Entrypoints(input EntrypointsInput) (*resty.Response, Entrypoints, error) {
	err := validator.New().Struct(input)
	if err != nil {
		return nil, Entrypoints{}, errors.Wrap(err, "failed to get entrypoints: invalid input")
	}

	resp, err := c.post(fmt.Sprintf("/chains/%s/blocks/%s/helpers/scripts/entrypoints", c.chain, input.BlockID.ID()), input.Entrypoints)
	if err != nil {
		return resp, Entrypoints{}, errors.Wrap(err, "failed to get entrypoints")
	}

	var entrypoints Entrypoints
	err = json.Unmarshal(resp.Body(), &entrypoints)
	if err != nil {
		return resp, Entrypoints{}, errors.Wrap(err, "failed to get entrypoints: failed to parse json")
	}

	return resp, entrypoints, nil
}

/*
PackDataInput is the input for the PackData function

RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-helpers-scripts-pack-data
*/
type PackDataInput struct {
	// The block (height) of which you want to make the query.
	BlockID BlockID `validate:"required"`
	// The data to pack
	Data PackDataBody `validate:"required"`
}

/*
PackDataBody is the data to pack for the PackData function

RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-helpers-scripts-pack-data
*/
type PackDataBody struct {
	Data *json.RawMessage `json:"data"`
	Type *json.RawMessage `json:"type"`
	Gas  string           `json:"gas"`
}

/*
PackedData is the packed data for the PackData function

RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-helpers-scripts-pack-data
*/
type PackedData struct {
	Packed string `json:"packed"`
	Gas    string `json:"gas"`
}

/*
PackData computes the serialized version of some data expression using the same algorithm as script instruction PACK

Path:
	../<block_id>/helpers/scripts/pack_data (POST)

RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-helpers-scripts-pack-data
*/
func (c *Client) PackData(input PackDataInput) (*resty.Response, PackedData, error) {
	err := validator.New().Struct(input)
	if err != nil {
		return nil, PackedData{}, errors.Wrap(err, "failed to pack data: invalid input")
	}

	resp, err := c.post(fmt.Sprintf("/chains/%s/blocks/%s/helpers/scripts/pack_data", c.chain, input.BlockID.ID()), input.Data)
	if err != nil {
		return resp, PackedData{}, errors.Wrap(err, "failed to pack data")
	}

	var packedData PackedData
	err = json.Unmarshal(resp.Body(), &packedData)
	if err != nil {
		return resp, PackedData{}, errors.Wrap(err, "failed to pack data: failed to parse json")
	}

	return resp, packedData, nil
}

/*
RunCodeInput is the input for the RunCode function

RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-helpers-scripts-run-code
*/
type RunCodeInput struct {
	// The block (height) of which you want to make the query.
	BlockID BlockID `validate:"required"`
	// The code to run
	Code RunCodeBody `validate:"required"`
}

/*
RunCodeBody is the body of the RunCode RPC

RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-helpers-scripts-run-code
*/
type RunCodeBody struct {
	Script     *json.RawMessage `json:"script"`
	Storage    *json.RawMessage `json:"storage"`
	Input      *json.RawMessage `json:"input"`
	Amount     string           `json:"amount"`
	Balance    string           `json:"balance"`
	ChainID    string           `json:"chain_id"`
	Source     string           `json:"source,omitempty"`
	Payer      string           `json:"payer,omitempty"`
	Gas        string           `json:"gas,omitempty"`
	Entrypoint string           `json:"entrypoint,omitempty"`
}

/*
RanCode is the response to running code with the RunCode function

RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-helpers-scripts-run-code
*/
type RanCode struct {
	Storage     *json.RawMessage `json:"storage"`
	Operations  []Operations     `json:"operations"`
	BigMapDiffs []BigMapDiff     `json:"big_map_diff,omitempty"`
}

/*
RunCode runs a piece of code in the current context

Path:
	../<block_id>/helpers/scripts/run_code (POST)

RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-helpers-scripts-run-code
*/
func (c *Client) RunCode(input RunCodeInput) (*resty.Response, RanCode, error) {
	err := validator.New().Struct(input)
	if err != nil {
		return nil, RanCode{}, errors.Wrap(err, "failed to run code: invalid input")
	}

	resp, err := c.post(fmt.Sprintf("/chains/%s/blocks/%s/helpers/scripts/run_code", c.chain, input.BlockID.ID()), input.Code)
	if err != nil {
		return resp, RanCode{}, errors.Wrap(err, "failed to run code")
	}

	var rancode RanCode
	err = json.Unmarshal(resp.Body(), &rancode)
	if err != nil {
		return resp, RanCode{}, errors.Wrap(err, "failed to run code: failed to parse json")
	}

	return resp, rancode, nil
}

/*
RunOperationInput is the input for the RunOperation function.

RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-helpers-scripts-run-operation
*/
type RunOperationInput struct {
	// The block (height) of which you want to make the query.
	BlockID BlockID `validate:"required"`
	// The operation to run
	Operation RunOperation `json:"operation" validate:"required"`
}

/*
RunOperation is the operation to run in the RunOperation function.

RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-helpers-scripts-run-operation
*/
type RunOperation struct {
	Operation Operations `json:"operation" validate:"required"`
	ChainID   string     `json:"chain_id" validate:"required"`
}

/*
RunOperation will run an operation without signature checks.

Path:
	../<block_id>/helpers/scripts/run_operation (POST)

RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-helpers-scripts-run-operation
*/
func (c *Client) RunOperation(input RunOperationInput) (*resty.Response, Operations, error) {
	err := validator.New().Struct(input)
	if err != nil {
		return nil, Operations{}, errors.Wrap(err, "failed to run operation: invalid input")
	}

	resp, err := c.post(fmt.Sprintf("/chains/%s/blocks/%s/helpers/scripts/run_operation", c.chain, input.BlockID.ID()), input.Operation)
	if err != nil {
		return resp, input.Operation.Operation, errors.Wrapf(err, "failed to run operation")
	}

	var op Operations
	err = json.Unmarshal(resp.Body(), &op)
	if err != nil {
		return resp, input.Operation.Operation, errors.Wrap(err, "failed to run operation: failed to parse json")
	}

	return resp, op, nil
}

/*
TraceCodeInput is the input for TraceCode function

RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-helpers-scripts-trace-code
*/
type TraceCodeInput struct {
	// The block (height) of which you want to make the query.
	BlockID BlockID `validate:"required"`
	// The code to trace
	Code RunCodeBody
}

/*
TracedCode is traced code returned from the TraceCode function

RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-helpers-scripts-trace-code
*/
type TracedCode struct {
	RanCode
	Trace Trace `json:"trace"`
}

/*
Trace is a trace in traced code returned from the TraceCode function

RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-helpers-scripts-trace-code
*/
type Trace struct {
	Location int     `json:"location"`
	Gas      string  `json:"gas"`
	Stack    []Stack `json:"stack"`
}

/*
Stack is a stack in a trace in traced code returned from the TraceCode function

RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-helpers-scripts-trace-code
*/
type Stack struct {
	Item  *json.RawMessage `json:"item"`
	Annot string           `json:"annot,omitempty"`
}

/*
TraceCode runs a piece of code in the current context, keeping a trace

Path:
	../<block_id>/helpers/scripts/trace_code (POST)

RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-helpers-scripts-trace-code
*/
func (c *Client) TraceCode(input TraceCodeInput) (*resty.Response, TracedCode, error) {
	err := validator.New().Struct(input)
	if err != nil {
		return nil, TracedCode{}, errors.Wrap(err, "failed to trace code: invalid input")
	}

	resp, err := c.post(fmt.Sprintf("/chains/%s/blocks/%s/helpers/scripts/trace_code", c.chain, input.BlockID.ID()), input.Code)
	if err != nil {
		return resp, TracedCode{}, errors.Wrapf(err, "failed to trace code")
	}

	var tracedCode TracedCode
	err = json.Unmarshal(resp.Body(), &tracedCode)
	if err != nil {
		return resp, TracedCode{}, errors.Wrap(err, "failed to trace code: failed to parse json")
	}

	return resp, tracedCode, nil
}

/*
TypeCheckcodeInput is the input for the TypecheckCode functions

RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-helpers-scripts-typecheck-code
*/
type TypeCheckcodeInput struct {
	// The block (height) of which you want to make the query.
	BlockID BlockID `validate:"required"`
	// The code to type check
	Code TypecheckCodeBody `validate:"required"`
}

/*
TypecheckCodeBody is body for the input for the TypecheckCode functions

RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-helpers-scripts-typecheck-code
*/
type TypecheckCodeBody struct {
	Program *json.RawMessage `json:"program"`
	Gas     string           `json:"gas"`
	Legacy  bool             `json:"legacy,omitempty"`
}

/*
TypecheckedCode is typechecked code returned by the TypecheckCode function

RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-helpers-scripts-typecheck-code
*/
type TypecheckedCode struct {
	TypeMap []struct {
		Location    int                `json:"location"`
		StackBefore []*json.RawMessage `json:"stack_before"`
		StackAfter  []*json.RawMessage `json:"stack_after"`
	} `json:"type_map"`
	Gas string `json:"gas"`
}

/*
TypecheckCode typechecks a piece of code in the current context

Path:
	../<block_id>/helpers/scripts/typecheck_code (POST)

RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-helpers-scripts-typecheck-code
*/
func (c *Client) TypecheckCode(input TypeCheckcodeInput) (*resty.Response, TypecheckedCode, error) {
	err := validator.New().Struct(input)
	if err != nil {
		return nil, TypecheckedCode{}, errors.Wrap(err, "failed to typecheck code: invalid input")
	}

	resp, err := c.post(fmt.Sprintf("/chains/%s/blocks/%s/helpers/scripts/typecheck_code", c.chain, input.BlockID.ID()), input.Code)
	if err != nil {
		return resp, TypecheckedCode{}, errors.Wrapf(err, "failed to typecheck code")
	}

	var typecheckCode TypecheckedCode
	err = json.Unmarshal(resp.Body(), &typecheckCode)
	if err != nil {
		return resp, TypecheckedCode{}, errors.Wrap(err, "failed to typecheck code: failed to parse json")
	}

	return resp, typecheckCode, nil
}

/*
TypecheckDataInput is the input for the TypecheckData functions

RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-helpers-scripts-typecheck-data
*/
type TypecheckDataInput struct {
	// The block (height) of which you want to make the query.
	BlockID BlockID `validate:"required"`
	// The code to type check
	Data TypecheckDataBody `validate:"required"`
}

/*
TypecheckDataBody is body for the input for the TypecheckData functions

RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-helpers-scripts-typecheck-data
*/
type TypecheckDataBody struct {
	Data   *json.RawMessage `json:"data"`
	Type   *json.RawMessage `json:"type"`
	Gas    string           `json:"gas"`
	Legacy bool             `json:"legacy,omitempty"`
}

/*
TypecheckedData is body for the input for the TypecheckData functions

RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-helpers-scripts-typecheck-data
*/
type TypecheckedData struct {
	Gas string `json:"gas"`
}

/*
TypecheckData checks that some data expression is well formed and of a given type in the current context

Path:
	../<block_id>/helpers/scripts/typecheck_data (POST)

RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-helpers-scripts-typecheck-data
*/
func (c *Client) TypecheckData(input TypecheckDataInput) (*resty.Response, TypecheckedData, error) {
	err := validator.New().Struct(input)
	if err != nil {
		return nil, TypecheckedData{}, errors.Wrap(err, "failed to typecheck data: invalid input")
	}

	resp, err := c.post(fmt.Sprintf("/chains/%s/blocks/%s/helpers/scripts/typecheck_data", c.chain, input.BlockID.ID()), input.Data)
	if err != nil {
		return resp, TypecheckedData{}, errors.Wrapf(err, "failed to typecheck data")
	}

	var typecheckData TypecheckedData
	err = json.Unmarshal(resp.Body(), &typecheckData)
	if err != nil {
		return resp, TypecheckedData{}, errors.Wrap(err, "failed to typecheck data: failed to parse json")
	}

	return resp, typecheckData, nil
}

func stripBranchFromForgedOperation(operation string, signed bool) (string, string, error) {
	if signed && len(operation) <= 128 {
		return "", operation, errors.New("failed to unforge branch from operation")
	}

	if signed {
		operation = operation[:len(operation)-128]
	}

	var result, rest string
	if len(operation) < 64 {
		result = operation
	} else {
		result = operation[:64]
		rest = operation[64:]
	}

	resultByts, err := hex.DecodeString(result)
	if err != nil {
		return "", operation, errors.New("failed to unforge branch from operation")
	}

	branch := crypto.B58cencode(resultByts, []byte{1, 52})
	if err != nil {
		return branch, rest, errors.Wrap(err, "failed to unforge branch from operation")
	}

	return branch, rest, nil
}
