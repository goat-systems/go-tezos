package gotezos

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"time"

	validator "github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

/*
FrozenBalance represents the frozen balance RPC on the tezos network.

RPC:
	../<block_id>/context/delegates/<pkh>/frozen_balance (GET)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-context-delegates-pkh-frozen-balance
*/
type FrozenBalance struct {
	Deposits *Int `json:"deposits"`
	Fees     *Int `json:"fees"`
	Rewards  *Int `json:"rewards"`
}

/*
Delegate represents the frozen delegate RPC on the tezos network.

RPC:
	../<block_id>/context/delegates/<pkh> (GET)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-context-delegates-pkh
*/
type Delegate struct {
	Balance              string `json:"balance"`
	FrozenBalance        string `json:"frozen_balance"`
	FrozenBalanceByCycle []struct {
		Cycle   int  `json:"cycle"`
		Deposit *Int `json:"deposit"`
		Fees    *Int `json:"fees"`
		Rewards *Int `json:"rewards"`
	} `json:"frozen_balance_by_cycle"`
	StakingBalance    string   `json:"staking_balance"`
	DelegateContracts []string `json:"delegated_contracts"`
	DelegatedBalance  string   `json:"delegated_balance"`
	Deactivated       bool     `json:"deactivated"`
	GracePeriod       int      `json:"grace_period"`
}

/*
BakingRights represents the baking rights RPC on the tezos network.

RPC:
	../<block_id>/helpers/baking_rights (GET)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-helpers-baking-rights
*/
type BakingRights []BakingRight

type BakingRight struct {
	Level         int       `json:"level"`
	Delegate      string    `json:"delegate"`
	Priority      int       `json:"priority"`
	EstimatedTime time.Time `json:"estimated_time"`
}

/*
EndorsingRights represents the endorsing rights RPC on the tezos network.

RPC:
	../<block_id>/helpers/baking_rights (GET)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-helpers-baking-rights
*/
type EndorsingRights []EndorsingRight

type EndorsingRight struct {
	Level         int       `json:"level"`
	Delegate      string    `json:"delegate"`
	Slots         []int     `json:"slots"`
	EstimatedTime time.Time `json:"estimated_time"`
}

/*
BakingRightsInput is the input for the goTezos.BakingRights function.

Function:
	func (t *GoTezos) BakingRights(input *BakingRightsInput) (*BakingRights, error) {}
*/
type BakingRightsInput struct {
	// The block level of which you want to make the query.
	Level int

	// The cycle of which you want to make the query.
	Cycle int

	// The delegate public key hash of which you want to make the query.
	Delegate string

	// The max priotity of which you want to make the query.
	MaxPriority int

	// The hash of block (height) of which you want to make the query.
	// Required.
	BlockHash string `validate:"required"`
}

/*
EndorsingRightsInput is the input for the goTezos.EndorsingRights function.

Function:
	func (t *GoTezos) EndorsingRights(input *EndorsingRightsInput) (*EndorsingRights, error) {}
*/
type EndorsingRightsInput struct {
	// The block level of which you want to make the query.
	Level int

	// The cycle of which you want to make the query.
	Cycle int

	// The delegate public key hash of which you want to make the query.
	Delegate string

	// The hash of block (height) of which you want to make the query.
	// Required.
	BlockHash string `validate:"required"`
}

/*
DelegatesInput is the input for the goTezos.Delegates function.

Function:
	func (t *GoTezos) Delegates(blockhash string) ([]string, error) {}
*/
type DelegatesInput struct {
	// The block level of which you want to make the query.
	active bool

	// The cycle of which you want to make the query.
	inactive bool

	// The hash of block (height) of which you want to make the query.
	// Required.
	BlockHash string `validate:"required"`
}

/*
DelegatedContracts Returns the list of contracts that delegate to a given delegate.

Path:
	../<block_id>/context/delegates/<pkh>/delegated_contracts (GET)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-context-delegates-pkh-delegated-contracts

Parameters:

	blockhash:
		The hash of block (height) of which you want to make the query.

	delegate:
		The tz(1-3) address of the delegate.
*/
func (t *GoTezos) DelegatedContracts(blockhash, delegate string) ([]*string, error) {
	resp, err := t.get(fmt.Sprintf("/chains/main/blocks/%s/context/delegates/%s/delegated_contracts", blockhash, delegate))
	if err != nil {
		return []*string{}, errors.Wrapf(err, "could not get delegations for '%s'", delegate)
	}

	var list []*string
	err = json.Unmarshal(resp, &list)
	if err != nil {
		return []*string{}, errors.Wrapf(err, "could not unmarshal delegations for '%s'", delegate)
	}

	return list, nil
}

/*
DelegatedContractsAtCycle returns the list of contracts that delegate to a given delegate at a specific cycle or snapshot.

Path:
	../<block_id>/context/delegates/<pkh>/delegated_contracts (GET)
Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-context-delegates-pkh-delegated-contracts

Parameters:

	cycle:
		The cycle of which you want to make the query.

	delegate:
		The tz(1-3) address of the delegate.
*/
func (t *GoTezos) DelegatedContractsAtCycle(cycle int, delegate string) ([]*string, error) {
	snapshot, err := t.Cycle(cycle)
	if err != nil {
		return []*string{}, errors.Wrapf(err, "could not get delegations for '%s' at cycle '%d'", delegate, cycle)
	}

	delegations, err := t.DelegatedContracts(snapshot.BlockHash, delegate)
	if err != nil {
		return []*string{}, errors.Wrapf(err, "could not get delegations at cycle '%d'", cycle)
	}

	return delegations, nil
}

/*
FrozenBalance returns the total frozen balances of a given delegate, this includes the frozen deposits, rewards and fees.

Path:
	../<block_id>/context/delegates/<pkh>/frozen_balance (GET)
Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-context-delegates-pkh-frozen-balance

Parameters:

	cycle:
		The cycle of which you want to make the query.

	delegate:
		The tz(1-3) address of the delegate.
*/
func (t *GoTezos) FrozenBalance(cycle int, delegate string) (FrozenBalance, error) {
	level := (cycle+1)*(t.NetworkConstants.BlocksPerCycle) + 1

	head, err := t.Block(level)
	if err != nil {
		return FrozenBalance{}, errors.Wrapf(err, "failed to get frozen balance at cycle '%d' for delegate '%s'", cycle, delegate)
	}

	resp, err := t.get(fmt.Sprintf("/chains/main/blocks/%s/context/raw/json/contracts/index/%s/frozen_balance/%d/", head.Hash, delegate, cycle))
	if err != nil {
		return FrozenBalance{}, errors.Wrapf(err, "failed to get frozen balance at cycle '%d' for delegate '%s'", cycle, delegate)
	}

	var frozenBalance FrozenBalance
	err = json.Unmarshal(resp, &frozenBalance)
	if err != nil {
		return frozenBalance, errors.Wrapf(err, "failed to unmarshal frozen balance at cycle '%d' for delegate '%s'", cycle, delegate)
	}

	return frozenBalance, nil
}

/*
Delegate gets everything about a delegate.

Path:
	../<block_id>/context/delegates/<pkh> (GET)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-context-delegates-pkh

Parameters:

	cycle:
		The cycle of which you want to make the query.

	delegate:
		The tz(1-3) address of the delegate.
*/
func (t *GoTezos) Delegate(blockhash, delegate string) (Delegate, error) {
	resp, err := t.get(fmt.Sprintf("/chains/main/blocks/%s/context/delegates/%s", blockhash, delegate))
	if err != nil {
		return Delegate{}, errors.Wrapf(err, "could not get delegate '%s'", delegate)
	}

	var d Delegate
	err = json.Unmarshal(resp, &d)
	if err != nil {
		return d, errors.Wrapf(err, "could not unmarshal delegate '%s'", delegate)
	}

	return d, nil
}

/*
StakingBalance returns the total amount of tokens delegated to a given delegate. This includes the balances of all the contracts
that delegate to it, but also the balance of the delegate itself and its frozen fees and deposits. The rewards do not count in
the delegated balance until they are unfrozen.

Path:
	../<block_id>/context/delegates/<pkh>/staking_balance (GET)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-context-delegates-pkh-staking-balance

Parameters:

	blockhash:
		The hash of block (height) of which you want to make the query.

	delegate:
		The tz(1-3) address of the delegate.
*/
func (t *GoTezos) StakingBalance(blockhash, delegate string) (*big.Int, error) {
	resp, err := t.get(fmt.Sprintf("/chains/main/blocks/%s/context/delegates/%s/staking_balance", blockhash, delegate))
	if err != nil {
		return big.NewInt(0), errors.Wrapf(err, "could not get staking balance for '%s'", delegate)
	}

	balance, err := newInt(resp)
	if err != nil {
		return big.NewInt(0), errors.Wrapf(err, "could not unmarshal staking balance for '%s'", delegate)
	}

	return balance.Big, nil
}

/*
StakingBalanceAtCycle returns the total amount of tokens delegated to a given delegate. This includes the balances of all the contracts
that delegate to it, but also the balance of the delegate itself and its frozen fees and deposits. The rewards do not count in
the delegated balance until they are unfrozen.

Path:
	../<block_id>/context/delegates/<pkh>/staking_balance (GET)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-context-delegates-pkh-staking-balance

Parameters:

	cycle:
		The cycle of which you want to make the query.

	delegate:
		The tz(1-3) address of the delegate.
*/
func (t *GoTezos) StakingBalanceAtCycle(cycle int, delegate string) (*big.Int, error) {
	snapshot, err := t.Cycle(cycle)
	if err != nil {
		return big.NewInt(0), errors.Wrapf(err, "could not get staking balance for '%s' at cycle '%d'", delegate, cycle)
	}

	balance, err := t.StakingBalance(snapshot.BlockHash, delegate)
	if err != nil {
		return big.NewInt(0), errors.Wrapf(err, "could not get staking balance for '%s' at cycle '%d'", delegate, cycle)
	}

	return balance, nil
}

/*
BakingRights retrieves the list of delegates allowed to bake a block. By default, it gives the best baking priorities
for bakers that have at least one opportunity below the 64th priority for the next block.
Parameters `level` and `cycle` can be used to specify the (valid) level(s) in the past or
future at which the baking rights have to be returned. Parameter `delegate` can be used to
restrict the results to the given delegates. If parameter `all` is set, all the baking
opportunities for each baker at each level are returned, instead of just the first one.
Returns the list of baking slots. Also returns the minimal timestamps that correspond
to these slots. The timestamps are omitted for levels in the past, and are only estimates
for levels later that the next block, based on the hypothesis that all predecessor blocks
were baked at the first priority.

Path:
	../<block_id>/helpers/baking_rights (GET)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-helpers-baking-rights

Parameters:

	BakingRightsInput:
		Modifies the BakingRights RPC query by passing optional URL parameters. BlockHash is required.

*/
func (t *GoTezos) BakingRights(input BakingRightsInput) (*BakingRights, error) {
	err := validator.New().Struct(input)
	if err != nil {
		return &BakingRights{}, errors.Wrap(err, "invalid input")
	}

	resp, err := t.get(fmt.Sprintf("/chains/main/blocks/%s/helpers/baking_rights", input.BlockHash), input.contructRPCOptions()...)
	if err != nil {
		return &BakingRights{}, errors.Wrapf(err, "could not get baking rights")
	}

	var bakingRights BakingRights
	err = json.Unmarshal(resp, &bakingRights)
	if err != nil {
		return &BakingRights{}, errors.Wrapf(err, "could not unmarshal baking rights")
	}

	return &bakingRights, nil
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

	return opts
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
	../<block_id>/helpers/endorsing_rights (GET)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-helpers-endorsing-rights

Parameters:
	EndorsingRightsInput:
		Modifies the EndorsingRights RPC query by passing optional URL parameters. BlockHash is required.

*/
func (t *GoTezos) EndorsingRights(input EndorsingRightsInput) (*EndorsingRights, error) {
	err := validator.New().Struct(input)
	if err != nil {
		return &EndorsingRights{}, errors.Wrap(err, "invalid input")
	}

	resp, err := t.get(fmt.Sprintf("/chains/main/blocks/%s/helpers/endorsing_rights", input.BlockHash), input.contructRPCOptions()...)
	if err != nil {
		return &EndorsingRights{}, errors.Wrap(err, "could not get endorsing rights")
	}

	var endorsingRights EndorsingRights
	err = json.Unmarshal(resp, &endorsingRights)
	if err != nil {
		return &EndorsingRights{}, errors.Wrapf(err, "could not unmarshal endorsing rights")
	}

	return &endorsingRights, nil
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
Delegates lists all registered delegates.

Path:
	../<block_id>/context/delegates (GET)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-context-delegates

Parameters:

	cycle:
		The cycle of which you want to make the query.

	delegate:
		The tz(1-3) address of the delegate.
*/
func (t *GoTezos) Delegates(input DelegatesInput) ([]*string, error) {
	err := validator.New().Struct(input)
	if err != nil {
		return []*string{}, errors.Wrap(err, "invalid input")
	}

	resp, err := t.get(fmt.Sprintf("/chains/main/blocks/%s/context/delegates", input.BlockHash), input.contructRPCOptions()...)
	if err != nil {
		return []*string{}, errors.Wrap(err, "could not get delegates")
	}

	var list []*string
	err = json.Unmarshal(resp, &list)
	if err != nil {
		return []*string{}, errors.Wrap(err, "could not unmarshal delegates")
	}

	return list, nil
}

func (d *DelegatesInput) contructRPCOptions() []rpcOptions {
	var opts []rpcOptions
	if d.active {
		opts = append(opts, rpcOptions{
			"active",
			"true",
		})
	} else {
		opts = append(opts, rpcOptions{
			"active",
			"false",
		})
	}

	if d.inactive {
		opts = append(opts, rpcOptions{
			"inactive",
			"true",
		})
	} else {
		opts = append(opts, rpcOptions{
			"inactive",
			"false",
		})
	}

	return opts
}
