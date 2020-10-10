package rpc

import (
	"encoding/json"
	"fmt"
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
	Deposits int `json:"deposits"`
	Fees     int `json:"fees"`
	Rewards  int `json:"rewards"`
}

// UnmarshalJSON satisfies json.Marshaler
func (f *FrozenBalance) UnmarshalJSON(data []byte) error {
	type FrozenBalanceHelper struct {
		Deposits string `json:"deposits"`
		Fees     string `json:"fees"`
		Rewards  string `json:"rewards"`
	}

	var frozenBalanceHelper FrozenBalanceHelper
	if err := json.Unmarshal(data, &frozenBalanceHelper); err != nil {
		return err
	}

	deposits, err := strconv.Atoi(frozenBalanceHelper.Deposits)
	if err != nil {
		return err
	}
	f.Deposits = deposits

	fees, err := strconv.Atoi(frozenBalanceHelper.Fees)
	if err != nil {
		return err
	}
	f.Fees = fees

	rewards, err := strconv.Atoi(frozenBalanceHelper.Rewards)
	if err != nil {
		return err
	}
	f.Rewards = rewards

	return nil
}

// MarshalJSON satisfies json.Marshaler
func (f *FrozenBalance) MarshalJSON() ([]byte, error) {
	frozenBalance := struct {
		Deposits string `json:"deposits"`
		Fees     string `json:"fees"`
		Rewards  string `json:"rewards"`
	}{
		strconv.Itoa(f.Deposits),
		strconv.Itoa(f.Fees),
		strconv.Itoa(f.Rewards),
	}

	return json.Marshal(frozenBalance)
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
		Cycle   int `json:"cycle"`
		Deposit int `json:"deposit,string"`
		Fees    int `json:"fees,string"`
		Rewards int `json:"rewards,string"`
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
type BakingRights []struct {
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
	https://tezos.gitlab.io/api/rpc.html#get-block-id-helpers-endorsing-rights
*/
type EndorsingRights []struct {
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
	// The hash of block (height) of which you want to make the query.
	BlockHash string `validate:"required"`

	// The block level of which you want to make the query.
	Level int

	// The cycle of which you want to make the query.
	Cycle int

	// The delegate public key hash of which you want to make the query.
	Delegate string

	// The max priotity of which you want to make the query.
	MaxPriority int
}

/*
EndorsingRightsInput is the input for the goTezos.EndorsingRights function.

Function:
	func (t *GoTezos) EndorsingRights(input *EndorsingRightsInput) (*EndorsingRights, error) {}
*/
type EndorsingRightsInput struct {
	// The hash of block (height) of which you want to make the query.
	BlockHash string `validate:"required"`

	// The block level of which you want to make the query.
	Level int

	// The cycle of which you want to make the query.
	Cycle int

	// The delegate public key hash of which you want to make the query.
	Delegate string
}

/*
DelegatesInput is the input for the goTezos.Delegates function.

Function:
	func (t *GoTezos) Delegates(input DelegatesInput) ([]string, error) {}
*/
type DelegatesInput struct {
	// The block level of which you want to make the query. If empty Cycle is required.
	Blockhash string
	// The cycle to get the balance at. If empty Blockhash is required.
	Cycle int
	// The block level of which you want to make the query.
	active bool
	// The cycle of which you want to make the query.
	inactive bool
}

func (d *DelegatesInput) validate() error {
	if d.Blockhash == "" && d.Cycle == 0 {
		return errors.New("invalid input: missing key cycle or blockhash")
	} else if d.Blockhash != "" && d.Cycle != 0 {
		return errors.New("invalid input: cannot have both cycle and blockhash")
	}

	err := validator.New().Struct(d)
	if err != nil {
		return errors.Wrap(err, "invalid input")
	}

	return nil
}

/*
StakingBalanceInput is the input for the goTezos.StakingBalance function.

Function:
	func (t *GoTezos) StakingBalance(blockhash, delegate string) (int, error) {}
*/
type StakingBalanceInput struct {
	// The block level of which you want to make the query. If empty Cycle is required.
	Blockhash string
	// The cycle to get the balance at. If empty Blockhash is required.
	Cycle int
	// The delegate that you want to make the query.
	Delegate string `validate:"required"`
}

func (s *StakingBalanceInput) validate() error {
	if s.Blockhash == "" && s.Cycle == 0 {
		return errors.New("invalid input: missing key cycle or blockhash")
	} else if s.Blockhash != "" && s.Cycle != 0 {
		return errors.New("invalid input: cannot have both cycle and blockhash")
	}

	err := validator.New().Struct(s)
	if err != nil {
		return errors.Wrap(err, "invalid input")
	}

	return nil
}

/*
DelegateInput is the input for the client.Delegate() function.

Function:
	func (t *GoTezos) DelegatedContracts(input DelegatedContractsInput) ([]string, error)  {}
*/
type DelegateInput struct {
	// The block level of which you want to make the query. If empty Cycle is required.
	Blockhash string
	// The cycle to get the balance at. If empty Blockhash is required.
	Cycle int
	// The delegate that you want to make the query.
	Delegate string `validate:"required"`
}

func (s *DelegateInput) validate() error {
	if s.Blockhash == "" && s.Cycle == 0 {
		return errors.New("invalid input: missing key cycle or blockhash")
	} else if s.Blockhash != "" && s.Cycle != 0 {
		return errors.New("invalid input: cannot have both cycle and blockhash")
	}

	err := validator.New().Struct(s)
	if err != nil {
		return errors.Wrap(err, "invalid input")
	}

	return nil
}

/*
FrozenBalanceInput is the input for the client.Delegate() function.

Function:
	func (t *GoTezos) DelegatedContracts(input DelegatedContractsInput) ([]string, error)  {}
*/
type FrozenBalanceInput struct {
	// The cycle to get the balance at.
	Cycle int `validate:"required"`
	// The delegate that you want to make the query.
	Delegate string `validate:"required"`
}

/*
DelegatedContractsInput is the input for the goTezos.DelegatedContractsInput function.

Function:
	func (t *GoTezos) DelegatedContracts(input DelegatedContractsInput) ([]string, error)  {}
*/
type DelegatedContractsInput struct {
	// The block level of which you want to make the query. If empty Cycle is required.
	Blockhash string
	// The cycle to get the balance at. If empty Blockhash is required.
	Cycle int
	// The delegate that you want to make the query.
	Delegate string `validate:"required"`
}

func (s *DelegatedContractsInput) validate() error {
	if s.Blockhash == "" && s.Cycle == 0 {
		return errors.New("invalid input: missing key cycle or blockhash")
	} else if s.Blockhash != "" && s.Cycle != 0 {
		return errors.New("invalid input: cannot have both cycle and blockhash")
	}

	err := validator.New().Struct(s)
	if err != nil {
		return errors.Wrap(err, "invalid input")
	}

	return nil
}

/*
DelegatedContracts Returns the list of contracts that delegate to a given delegate.

Path:
	../<block_id>/context/delegates/<pkh>/delegated_contracts (GET)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-context-delegates-pkh-delegated-contracts
*/
func (c *Client) DelegatedContracts(input DelegatedContractsInput) ([]string, error) {
	err := input.validate()
	if err != nil {
		return []string{}, errors.Wrapf(err, "could not get delegations for delegate '%s'", input.Delegate)
	}

	input.Blockhash, err = c.extractBlockHash(input.Cycle, input.Blockhash)
	if err != nil {
		return []string{}, errors.Wrapf(err, "could not get delegations for delegate '%s'", input.Delegate)
	}

	var resp []byte
	resp, err = c.get(fmt.Sprintf("/chains/%s/blocks/%s/context/delegates/%s/delegated_contracts", c.chain, input.Blockhash, input.Delegate))
	if err != nil {
		return []string{}, errors.Wrapf(err, "could not get delegations for delegate '%s'", input.Delegate)
	}

	var list []string
	err = json.Unmarshal(resp, &list)
	if err != nil {
		return []string{}, errors.Wrapf(err, "could not unmarshal delegations for delegate '%s'", input.Delegate)
	}

	return list, nil
}

/*
FrozenBalance returns the total frozen balances of a given delegate, this includes the frozen deposits, rewards and fees.

Path:
	../<block_id>/context/delegates/<pkh>/frozen_balance (GET)
Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-context-delegates-pkh-frozen-balance
*/
func (c *Client) FrozenBalance(input FrozenBalanceInput) (FrozenBalance, error) {
	err := validator.New().Struct(input)
	if err != nil {
		return FrozenBalance{}, errors.Wrap(err, "invalid input")
	}

	level := (input.Cycle+1)*(c.NetworkConstants.BlocksPerCycle) + 1
	head, err := c.Block(level)
	if err != nil {
		return FrozenBalance{}, errors.Wrapf(err, "failed to get frozen balance at cycle '%d' for delegate '%s'", input.Cycle, input.Delegate)
	}

	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/context/raw/json/contracts/index/%s/frozen_balance/%d", c.chain, head.Hash, input.Delegate, input.Cycle))
	if err != nil {
		return FrozenBalance{}, errors.Wrapf(err, "failed to get frozen balance for delegate '%s'", input.Delegate)
	}

	var frozenBalance FrozenBalance
	err = json.Unmarshal(resp, &frozenBalance)
	if err != nil {
		return frozenBalance, errors.Wrapf(err, "failed to unmarshal frozen balance for delegate '%s'", input.Delegate)
	}

	return frozenBalance, nil
}

/*
Delegate gets everything about a delegate.

Path:
	../<block_id>/context/delegates/<pkh> (GET)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-context-delegates-pkh
*/
func (c *Client) Delegate(input DelegateInput) (Delegate, error) {
	err := input.validate()
	if err != nil {
		return Delegate{}, errors.Wrapf(err, "could not get delegate '%s'", input.Delegate)
	}

	input.Blockhash, err = c.extractBlockHash(input.Cycle, input.Blockhash)
	if err != nil {
		return Delegate{}, errors.Wrapf(err, "could not get delegate '%s'", input.Delegate)
	}

	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/context/delegates/%s", c.chain, input.Blockhash, input.Delegate))
	if err != nil {
		return Delegate{}, errors.Wrapf(err, "could not get delegate '%s'", input.Delegate)
	}

	var d Delegate
	err = json.Unmarshal(resp, &d)
	if err != nil {
		return d, errors.Wrapf(err, "could not unmarshal delegate '%s'", input.Delegate)
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

	StakingBalanceInput
		Modifies the StakingBalance RPC query by passing optional parameters. Delegate and (Cycle or Blockhash) is required.
*/
func (c *Client) StakingBalance(input StakingBalanceInput) (int, error) {
	err := input.validate()
	if err != nil {
		return 0, errors.Wrapf(err, "could not get staking balance for '%s'", input.Delegate)
	}

	input.Blockhash, err = c.extractBlockHash(input.Cycle, input.Blockhash)
	if err != nil {
		return 0, errors.Wrapf(err, "could not get staking balance for '%s'", input.Delegate)
	}

	var resp []byte
	resp, err = c.get(fmt.Sprintf("/chains/%s/blocks/%s/context/delegates/%s/staking_balance", c.chain, input.Blockhash, input.Delegate))
	if err != nil {
		return 0, errors.Wrapf(err, "could not get staking balance for '%s'", input.Delegate)
	}

	var stakingBalanceStr string
	err = json.Unmarshal(resp, &stakingBalanceStr)
	if err != nil {
		return 0, errors.Wrapf(err, "could not unmarshal staking balance for '%s'", input.Delegate)
	}

	return strconv.Atoi(stakingBalanceStr)
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
*/
func (c *Client) BakingRights(input BakingRightsInput) (*BakingRights, error) {
	err := validator.New().Struct(input)
	if err != nil {
		return &BakingRights{}, errors.Wrap(err, "invalid input")
	}

	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/helpers/baking_rights", c.chain, input.BlockHash), input.contructRPCOptions()...)
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
*/
func (c *Client) EndorsingRights(input EndorsingRightsInput) (*EndorsingRights, error) {
	err := validator.New().Struct(input)
	if err != nil {
		return &EndorsingRights{}, errors.Wrap(err, "invalid input")
	}

	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/helpers/endorsing_rights", c.chain, input.BlockHash), input.contructRPCOptions()...)
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
*/
func (c *Client) Delegates(input DelegatesInput) ([]string, error) {
	err := input.validate()
	if err != nil {
		return []string{}, errors.Wrap(err, "could not get delegates")
	}

	input.Blockhash, err = c.extractBlockHash(input.Cycle, input.Blockhash)
	if err != nil {
		return []string{}, errors.Wrap(err, "could not get delegates")
	}

	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/context/delegates", c.chain, input.Blockhash), input.contructRPCOptions()...)
	if err != nil {
		return []string{}, errors.Wrap(err, "could not get delegates")
	}

	var list []string
	err = json.Unmarshal(resp, &list)
	if err != nil {
		return []string{}, errors.Wrap(err, "could not unmarshal delegates")
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

func (c *Client) extractBlockHash(cycle int, blockhash string) (string, error) {
	if cycle != 0 {
		snapshot, err := c.Cycle(cycle)
		if err != nil {
			return "", errors.Wrapf(err, "failed to get cycle: %d", cycle)
		}

		return snapshot.BlockHash, nil
	}

	return blockhash, nil
}
