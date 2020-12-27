package rpc

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"

	validator "github.com/go-playground/validator/v10"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

/*
BakingRights represents the baking rights RPC on the tezos network.

Path:
	../<block_id>/helpers/baking_rights?(level=<block_level>)*&(cycle=<block_cycle>)*&(delegate=<pkh>)*&[max_priority=<int>]&[all]

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-helpers-baking-rights
*/
type BakingRights []struct {
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
func (c *Client) BakingRights(input BakingRightsInput) (*resty.Response, *BakingRights, error) {
	err := validator.New().Struct(input)
	if err != nil {
		return nil, &BakingRights{}, errors.Wrap(err, "failed to get baking rights: invalid input")
	}

	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/helpers/baking_rights", c.chain, input.BlockID.ID()), input.contructRPCOptions()...)
	if err != nil {
		return resp, &BakingRights{}, errors.Wrapf(err, "failed to get baking rights")
	}

	var bakingRights BakingRights
	err = json.Unmarshal(resp.Body(), &bakingRights)
	if err != nil {
		return resp, &BakingRights{}, errors.Wrapf(err, "failed to get baking rights: failed to parse json")
	}

	return resp, &bakingRights, nil
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
	// The hash of block (height) of which you want to make the query.
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
type EndorsingRights []struct {
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
func (c *Client) EndorsingRights(input EndorsingRightsInput) (*resty.Response, *EndorsingRights, error) {
	err := validator.New().Struct(input)
	if err != nil {
		return nil, &EndorsingRights{}, errors.Wrap(err, "failed to get endorsing rightsL invalid input")
	}

	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/helpers/endorsing_rights", c.chain, input.BlockID.ID()), input.contructRPCOptions()...)
	if err != nil {
		return resp, &EndorsingRights{}, errors.Wrap(err, "failed to get endorsing rights")
	}

	var endorsingRights EndorsingRights
	err = json.Unmarshal(resp.Body(), &endorsingRights)
	if err != nil {
		return resp, &EndorsingRights{}, errors.Wrapf(err, "failed to get endorsing rights: failed to parse json")
	}

	return resp, &endorsingRights, nil
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
	BlockID  BlockID  `validate:"required"`
	Branch   string   `validate:"required"`
	Contents Contents `validate:"required"`
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

	resp, err := c.post(fmt.Sprintf("/chains/%s/blocks/%s/helpers/forge/operations", c.chain, input.BlockID.ID()), v)
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

	resp, block, err := c.Block(input.BlockID)
	if err != nil {
		return resp, operation, errors.Wrap(err, "failed to forge operation: unable to get blockhash for BlockID")
	}

	_, operations, err := rpc.UnforgeOperation(UnforgeOperationInput{
		Blockhash: block.Hash,
		Operations: []UnforgeOperation{
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
