package gotezos

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

// FrozenBalance is representation of frozen balance on the Tezos network
type FrozenBalance struct {
	Deposits string `json:"deposits"`
	Fees     string `json:"fees"`
	Rewards  string `json:"rewards"`
}

// Delegate is representation of a delegate on the Tezos Network
type Delegate struct {
	Balance              string                 `json:"balance"`
	FrozenBalance        string                 `json:"frozen_balance"`
	FrozenBalanceByCycle []FrozenBalanceByCycle `json:"frozen_balance_by_cycle"`
	StakingBalance       string                 `json:"staking_balance"`
	DelegateContracts    []string               `json:"delegated_contracts"`
	DelegatedBalance     string                 `json:"delegated_balance"`
	Deactivated          bool                   `json:"deactivated"`
	GracePeriod          int                    `json:"grace_period"`
}

// FrozenBalanceByCycle a representation of frozen balance by cycle on the Tezos network
type FrozenBalanceByCycle struct {
	Cycle   int    `json:"cycle"`
	Deposit string `json:"deposit"`
	Fees    string `json:"fees"`
	Rewards string `json:"rewards"`
}

// BakingRights a representation of baking rights on the Tezos network
type BakingRights []struct {
	Level         int       `json:"level"`
	Delegate      string    `json:"delegate"`
	Priority      int       `json:"priority"`
	EstimatedTime time.Time `json:"estimated_time"`
}

// EndorsingRights is a representation of endorsing rights on the Tezos network
type EndorsingRights []struct {
	Level         int       `json:"level"`
	Delegate      string    `json:"delegate"`
	Slots         []int     `json:"slots"`
	EstimatedTime time.Time `json:"estimated_time"`
}

/*
DelegatedContracts RPC
Path: ../<block_id>/context/delegates/<pkh>/delegated_contracts (GET)
Link: https://tezos.gitlab.io/api/rpc.html#get-block-id-context-delegates-pkh-delegated-contracts
Description: Returns the list of contracts that delegate to a given delegate.

Parameters:
	blockhash:
		The hash of block (height) of which you want to make the query.
	delegate:
		The tz(1-3) address of the delegate.
*/
func (t *GoTezos) DelegatedContracts(blockhash, delegate string) ([]string, error) {
	resp, err := t.get(fmt.Sprintf("/chains/main/blocks/%s/context/delegates/%s/delegated_contracts", blockhash, delegate))
	if err != nil {
		return []string{}, errors.Wrapf(err, "could not get delegations for '%s'", delegate)
	}

	var list []string
	err = json.Unmarshal(resp, &list)
	if err != nil {
		return []string{}, errors.Wrapf(err, "could not unmarshal delegations for '%s'", delegate)
	}

	return list, nil
}

/*
DelegatedContractsAtCycle RPC
Path: ../<block_id>/context/delegates/<pkh>/delegated_contracts (GET)
Link: https://tezos.gitlab.io/api/rpc.html#get-block-id-context-delegates-pkh-delegated-contracts
Description: Returns the list of contracts that delegate to a given delegate at a specific cycle or snapshot.

Parameters:
	cycle:
		The cycle of which you want to make the query.
	delegate:
		The tz(1-3) address of the delegate.
*/
func (t *GoTezos) DelegatedContractsAtCycle(cycle int, delegate string) ([]string, error) {
	snapshot, err := t.Cycle(cycle)
	if err != nil {
		return []string{}, errors.Wrapf(err, "could not get delegations for '%s' at cycle '%d'", delegate, cycle)
	}

	delegations, err := t.DelegatedContracts(snapshot.BlockHash, delegate)
	if err != nil {
		return []string{}, errors.Wrapf(err, "could not get delegations at cycle '%d'", cycle)
	}

	return delegations, nil
}

/*
FrozenBalance RPC
Path: ../<block_id>/context/delegates/<pkh>/frozen_balance (GET)
Link: https://tezos.gitlab.io/api/rpc.html#get-block-id-context-delegates-pkh-frozen-balance
Description: Returns the total frozen balances of a given delegate, this includes the frozen deposits, rewards and fees.

Parameters:
	cycle:
		The cycle of which you want to make the query.
	delegate:
		The tz(1-3) address of the delegate.
*/
func (t *GoTezos) FrozenBalance(cycle int, delegate string) (FrozenBalance, error) {
	snapshot, err := t.Cycle(cycle)
	if err != nil {
		return FrozenBalance{}, errors.Wrapf(err, "could not get frozen balance at cycle '%d' for delegate '%s'", cycle, delegate)
	}

	resp, err := t.get(fmt.Sprintf("/chains/main/blocks/%s/context/raw/json/contracts/index/%s/frozen_balance/%d/", snapshot.BlockHash, delegate, cycle))
	if err != nil {
		return FrozenBalance{}, errors.Wrapf(err, "could not get frozen balance at cycle '%d' for delegate '%s'", cycle, delegate)
	}

	var frozenBalance FrozenBalance
	err = json.Unmarshal(resp, &frozenBalance)
	if err != nil {
		return frozenBalance, errors.Wrapf(err, "could not unmarshal frozen balance at cycle '%d' for delegate '%s'", cycle, delegate)
	}

	return frozenBalance, nil
}

/*
Delegate RPC
Path: ../<block_id>/context/delegates/<pkh> (GET)
Link: https://tezos.gitlab.io/api/rpc.html#get-block-id-context-delegates-pkh
Description: Everything about a delegate.

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
StakingBalance RPC
Path: ../<block_id>/context/delegates/<pkh> (GET)
Link: https://tezos.gitlab.io/api/rpc.html#get-block-id-context-delegates-pkh-staking-balance
Description: Everything about a delegate.

Parameters:
	blockhash:
		The hash of block (height) of which you want to make the query.
	delegate:
		The tz(1-3) address of the delegate.
*/
func (t *GoTezos) StakingBalance(blockhash, delegate string) (string, error) {
	resp, err := t.get(fmt.Sprintf("/chains/main/blocks/%s/context/delegates/%s/staking_balance", blockhash, delegate))
	if err != nil {
		return "", errors.Wrapf(err, "could not get staking balance for '%s'", delegate)
	}

	var balance string
	err = json.Unmarshal(resp, &balance)
	if err != nil {
		return balance, errors.Wrapf(err, "could not unmarshal staking balance for '%s'", delegate)
	}

	return balance, nil
}

/*
StakingBalanceAtCycle RPC
Path: ../<block_id>/context/delegates/<pkh> (GET)
Link: https://tezos.gitlab.io/api/rpc.html#get-block-id-context-delegates-pkh-staking-balance
Description: Everything about a delegate.

Parameters:
	cycle:
		The cycle of which you want to make the query.
	delegate:
		The tz(1-3) address of the delegate.
*/
func (t *GoTezos) StakingBalanceAtCycle(cycle int, delegate string) (string, error) {
	snapshot, err := t.Cycle(cycle)
	if err != nil {
		return "", errors.Wrapf(err, "could not get staking balance for '%s' at cycle '%d'", delegate, cycle)
	}

	balance, err := t.StakingBalance(snapshot.BlockHash, delegate)
	if err != nil {
		return balance, errors.Wrapf(err, "could not get staking balance for '%s' at cycle '%d'", delegate, cycle)
	}

	return balance, nil
}

/*
BakingRightsInput -
Description: The input for the baking_rights rpc query.
Function: func (t *GoTezos) BakingRights(input *BakingRightsInput) (*BakingRights, error) {}
*/
type BakingRightsInput struct {
	// The block level of which you want to make the query.
	Level *int

	// The cycle of which you want to make the query.
	Cycle *int

	// The delegate public key hash of which you want to make the query.
	Delegate *string

	// The max priotity of which you want to make the query.
	MaxPriority *int

	// The hash of block (height) of which you want to make the query.
	// Required.
	BlockHash *string `validate:"required"`
}

/*
BakingRights RPC
Path: ../<block_id>/helpers/baking_rights (GET)
Link: https://tezos.gitlab.io/api/rpc.html#get-block-id-context-delegates-pkh-staking-balance
Description: Retrieves the list of delegates allowed to bake a block. By default, it gives the best baking priorities
for bakers that have at least one opportunity below the 64th priority for the next block.
Parameters `level` and `cycle` can be used to specify the (valid) level(s) in the past or
future at which the baking rights have to be returned. Parameter `delegate` can be used to
restrict the results to the given delegates. If parameter `all` is set, all the baking
opportunities for each baker at each level are returned, instead of just the first one.
Returns the list of baking slots. Also returns the minimal timestamps that correspond
to these slots. The timestamps are omitted for levels in the past, and are only estimates
for levels later that the next block, based on the hypothesis that all predecessor blocks
were baked at the first priority.

Parameters:
	blockhash:
		The hash of block (height) of which you want to make the query.
	BakingRightsInput:
		Modifies the BakingRights RPC query by passing optional URL parameters. If <nil> no URL parameters will be used.

*/
func (t *GoTezos) BakingRights(input *BakingRightsInput) (*BakingRights, error) {
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

func (b *BakingRightsInput) contructRPCOptions() []RPCOptions {
	var opts []RPCOptions
	if b.Cycle != nil {
		opts = append(opts, RPCOptions{
			"cycle",
			strconv.Itoa(*b.Cycle),
		})
	}

	if b.Delegate != nil {
		opts = append(opts, RPCOptions{
			"delegate",
			*b.Delegate,
		})
	}

	if b.Level != nil {
		opts = append(opts, RPCOptions{
			"level",
			strconv.Itoa(*b.Level),
		})
	}

	if b.MaxPriority != nil {
		opts = append(opts, RPCOptions{
			"max_priority",
			strconv.Itoa(*b.MaxPriority),
		})
	}

	return opts
}

// EndorsingRights gets the endorsing rights for a specific cycle
func (t *GoTezos) EndorsingRights(cycle int) (EndorsingRights, error) {
	snapshot, err := t.Cycle(cycle)
	if err != nil {
		return EndorsingRights{}, errors.Wrapf(err, "could not get endorsing rights for cycle %d", cycle)
	}

	resp, err := t.get(
		fmt.Sprintf("/chains/main/blocks/%s/helpers/baking_rights", snapshot.BlockHash),
		RPCOptions{
			Key:   "cycle",
			Value: strconv.Itoa(cycle),
		},
	)
	if err != nil {
		return EndorsingRights{}, errors.Wrapf(err, "could not get endorsing rights for cycle '%d'", cycle)
	}

	var endorsingRights EndorsingRights
	err = json.Unmarshal(resp, &endorsingRights)
	if err != nil {
		return endorsingRights, errors.Wrapf(err, "could not unmarshal endorsing rights for cycle '%d'", cycle)
	}

	return endorsingRights, nil
}

// EndorsingRightsForDelegate gets the endorsing rights for a specific cycle
func (t *GoTezos) EndorsingRightsForDelegate(cycle int, delegate string) (EndorsingRights, error) {
	snapshot, err := t.Cycle(cycle)
	if err != nil {
		return EndorsingRights{}, errors.Wrapf(err, "could not get endorsing rights for delegate %s at cycle %d", delegate, cycle)
	}

	resp, err := t.get(
		fmt.Sprintf("/chains/main/blocks/%s/helpers/endorsing_rights", snapshot.BlockHash),
		RPCOptions{
			Key:   "cycle",
			Value: strconv.Itoa(cycle),
		},
		RPCOptions{
			Key:   "delegate",
			Value: delegate,
		},
	)
	if err != nil {
		return EndorsingRights{}, errors.Wrapf(err, "could not get endorsing rights for delegate '%s'", delegate)
	}

	var endorsingRights EndorsingRights
	err = json.Unmarshal(resp, &endorsingRights)
	if err != nil {
		return endorsingRights, errors.Wrapf(err, "could not unmarshal endorsing rights for delegate '%s'", delegate)
	}

	return endorsingRights, nil
}

// Delegates gets a list of all delegates at a blockhash
func (t *GoTezos) Delegates(blockhash string) ([]string, error) {
	resp, err := t.get(fmt.Sprintf("/chains/main/blocks/%s/context/delegates", blockhash))
	if err != nil {
		return []string{}, errors.Wrap(err, "could not get delegates")
	}

	var list []string
	err = json.Unmarshal(resp, &list)
	if err != nil {
		return list, errors.Wrap(err, "could not unmarshal delegates")
	}

	return list, nil
}
