package gotezos

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// FrozenBalanceRewards is a FrozenBalanceRewards query returned by the Tezos RPC API.
type FrozenBalanceRewards struct {
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

// Delegations retrieves a list of all currently delegated contracts for a delegate.
func (t *GoTezos) Delegations(blockhash, delegate string) ([]string, error) {
	resp, err := t.get(fmt.Sprintf("/chains/main/blocks/%s/context/delegates/%s/delegated_contracts", blockhash, delegate))
	if err != nil {
		return []string{}, errors.Wrapf(err, "could not get delegations for '%s'", delegate)
	}

	var list []string
	err = json.Unmarshal(resp, &list)
	if err != nil {
		return list, errors.Wrapf(err, "could not unmarshal delegations for '%s'", delegate)
	}

	return list, nil
}

// DelegationsAtCycle retrieves a list of all currently delegated contracts for a delegate at a specific cycle.
func (t *GoTezos) DelegationsAtCycle(cycle int, delegate string) ([]string, error) {
	snapshot, err := t.Cycle(cycle)
	if err != nil {
		return []string{}, errors.Wrapf(err, "could not get delegations for '%s' at cycle '%d'", delegate, cycle)
	}

	delegations, err := t.Delegations(snapshot.BlockHash, delegate)
	if err != nil {
		return delegations, errors.Wrapf(err, "could not get delegations for '%s' at cycle '%d'", delegate, cycle)
	}

	return delegations, nil
}

// FrozenBalance gets the rewards earned by a delegate for a specific cycle.
func (t *GoTezos) FrozenBalance(cycle int, delegate string) (FrozenBalanceRewards, error) {
	snapshot, err := t.Cycle(cycle)
	if err != nil {
		return FrozenBalanceRewards{}, errors.Wrapf(err, "could not get frozen balance at cycle '%d' for delegate '%s'", cycle, delegate)
	}

	resp, err := t.get(fmt.Sprintf("/chains/main/blocks/%s/context/raw/json/contracts/index/%s/frozen_balance/%d/", snapshot.BlockHash, delegate, cycle))
	if err != nil {
		return FrozenBalanceRewards{}, errors.Wrapf(err, "could not get frozen balance at cycle '%d' for delegate '%s'", cycle, delegate)
	}

	var frozenBalance FrozenBalanceRewards
	err = json.Unmarshal(resp, &frozenBalance)
	if err != nil {
		return frozenBalance, errors.Wrapf(err, "could not unmarshal frozen balance at cycle '%d' for delegate '%s'", cycle, delegate)
	}

	return frozenBalance, nil
}

// Delegate retrieves information about a delegate at the head block
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

// StakingBalance gets the staking balance of a delegate at a specific block hash
func (t *GoTezos) StakingBalance(headhash, delegate string) (string, error) {
	resp, err := t.get(fmt.Sprintf("/chains/main/blocks/%s/context/delegates/%s/staking_balance", headhash, delegate))
	if err != nil {
		return "", errors.Wrap(err, "could not get staking balance")
	}

	var balance string
	err = json.Unmarshal(resp, &balance)
	if err != nil {
		return balance, errors.Wrap(err, "could not unmarshal staking balance")
	}

	return balance, nil
}

// StakingBalanceAtCycle gets the staking balance of a delegate at a specific cycle
func (t *GoTezos) StakingBalanceAtCycle(cycle int, delegate string) (string, error) {
	snapshot, err := t.Cycle(cycle)
	if err != nil {
		return "", errors.Wrapf(err, "could not get staking balance for %s at cycle %d", delegate, cycle)
	}

	balance, err := t.StakingBalance(snapshot.BlockHash, delegate)
	if err != nil {
		return balance, errors.Wrapf(err, "could not get staking balance for %s at cycle %d", delegate, cycle)
	}

	return balance, nil
}

// BakingRights gets the baking rights at a specific blockhash with the given priority.
func (t *GoTezos) BakingRights(blockhash string, priority int) (BakingRights, error) {
	resp, err := t.get(
		fmt.Sprintf("/chains/main/blocks/%s/helpers/baking_rights", blockhash),
		params{
			key:   "max_priority",
			value: strconv.Itoa(priority),
		},
	)
	if err != nil {
		return BakingRights{}, errors.Wrapf(err, "could not get baking rights")
	}

	var bakingRights BakingRights
	err = json.Unmarshal(resp, &bakingRights)
	if err != nil {
		return bakingRights, errors.Wrapf(err, "could not unmarshal baking rights")
	}

	return bakingRights, nil
}

// BakingRightsAtCycle gets the baking rights at a specific cycle
func (t *GoTezos) BakingRightsAtCycle(cycle, priority int) (BakingRights, error) {
	snapshot, err := t.Cycle(cycle)
	if err != nil {
		return BakingRights{}, err
	}

	resp, err := t.get(
		fmt.Sprintf("/chains/main/blocks/%s/helpers/baking_rights", snapshot.BlockHash),
		params{
			key:   "cycle",
			value: strconv.Itoa(cycle),
		},
		params{
			key:   "max_priority",
			value: strconv.Itoa(priority),
		},
	)
	if err != nil {
		return BakingRights{}, errors.Wrap(err, "could not get baking rights")
	}

	var bakingRights BakingRights
	err = json.Unmarshal(resp, &bakingRights)
	if err != nil {
		return bakingRights, errors.Wrap(err, "could not unmarshal baking rights")
	}

	return bakingRights, nil
}

// BakingRightsForDelegate gets the baking rights for a delegate at a specific cycle with a certain priority level
func (t *GoTezos) BakingRightsForDelegate(cycle int, delegate string, priority int) (BakingRights, error) {
	snapshot, err := t.Cycle(cycle)
	if err != nil {
		return BakingRights{}, errors.Wrapf(err, "could not get baking rights for delegate %s at cycle %d", delegate, cycle)
	}

	resp, err := t.get(
		fmt.Sprintf("/chains/main/blocks/%s/helpers/baking_rights", snapshot.BlockHash),
		params{
			key:   "cycle",
			value: strconv.Itoa(cycle),
		},
		params{
			key:   "max_priority",
			value: strconv.Itoa(priority),
		},
	)
	if err != nil {
		return BakingRights{}, errors.Wrapf(err, "could not get baking rights for delegate '%s'", delegate)
	}

	var bakingRights BakingRights
	err = json.Unmarshal(resp, &bakingRights)
	if err != nil {
		return bakingRights, errors.Wrapf(err, "could not unmarshal baking rights for delegate '%s'", delegate)
	}

	return bakingRights, nil
}

// EndorsingRights gets the endorsing rights for a specific cycle
func (t *GoTezos) EndorsingRights(cycle int) (EndorsingRights, error) {
	snapshot, err := t.Cycle(cycle)
	if err != nil {
		return EndorsingRights{}, errors.Wrapf(err, "could not get endorsing rights for cycle %d", cycle)
	}

	resp, err := t.get(
		fmt.Sprintf("/chains/main/blocks/%s/helpers/baking_rights", snapshot.BlockHash),
		params{
			key:   "cycle",
			value: strconv.Itoa(cycle),
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
		params{
			key:   "cycle",
			value: strconv.Itoa(cycle),
		},
		params{
			key:   "delegate",
			value: delegate,
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
