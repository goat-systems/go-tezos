package rpc

import (
	"encoding/json"
	"regexp"
	"strconv"

	validator "github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

/*
GetFA12BalanceInput is the input for the goTezos.GetFA12Balance function.

Function:
	func (c *Client) GetFA12Balance(input GetFA12BalanceInput) (int, error) {}
*/
type GetFA12BalanceInput struct {
	// Blockhash is the block height at which to make the query. Can leave blank if using Cycle.
	Blockhash string
	// Cycle is the cycle in which to make the query. Can leave blank if using Blockhash.
	Cycle int
	// ChainID is the Chain ID of the chain you want to query
	ChainID string `validate:"required"`
	// Source to form the contents with. The operation is not forged or injected so it is possible for XTZ to be spent.
	Source string `validate:"required"`
	// FA12Contract address of the FA1.2 Contract you wish to query.
	FA12Contract string `validate:"required"`
	// OwnerAddress is the address to get the balance for in the FA1.2 contract
	OwnerAddress string `validate:"required"`
	// If true the function will use an intermediate contract deployed on Carthagenet, default mainnet.
	Testnet bool
	// If provided this will be the contract view address used to query the FA1.2 contract
	ContractViewAddress string
}

func (g *GetFA12BalanceInput) validate() error {
	if g.Blockhash == "" && g.Cycle == 0 {
		return errors.New("invalid input: missing key cycle or blockhash")
	} else if g.Blockhash != "" && g.Cycle != 0 {
		return errors.New("invalid input: cannot have both cycle and blockhash")
	}

	err := validator.New().Struct(g)
	if err != nil {
		return errors.Wrap(err, "invalid input")
	}

	return nil
}

/*
GetFA12SupplyInput is the input for the goTezos.GetFA12Supply function.

Function:
	func (c *Client) GetFA12Supply(input GetFA12SupplyInput) (int, error) {}
*/
type GetFA12SupplyInput struct {
	// Blockhash is the block height at which to make the query. Can leave blank if using Cycle.
	Blockhash string
	// Cycle is the cycle in which to make the query. Can leave blank if using Blockhash.
	Cycle int
	// ChainID is the Chain ID of the chain you want to query
	ChainID string `validate:"required"`
	// Source to form the contents with. The operation is not forged or injected so it is possible for XTZ to be spent.
	Source string `validate:"required"`
	// FA12Contract address of the FA1.2 Contract you wish to query.
	FA12Contract string `validate:"required"`
	// If true the function will use an intermediate contract deployed on Carthagenet, default mainnet.
	Testnet bool
	// If provided this will be the contract view address used to query the FA1.2 contract
	ContractViewAddress string
}

func (g *GetFA12SupplyInput) validate() error {
	if g.Blockhash == "" && g.Cycle == 0 {
		return errors.New("invalid input: missing key cycle or blockhash")
	} else if g.Blockhash != "" && g.Cycle != 0 {
		return errors.New("invalid input: cannot have both cycle and blockhash")
	}

	err := validator.New().Struct(g)
	if err != nil {
		return errors.Wrap(err, "invalid input")
	}

	return nil
}

/*
GetFA12AllowanceInput is the input for the goTezos.GetFA12Allowance function.

Function:
	func (c *Client) GetFA12Allowance(input GetFA12AllowanceInput) (int, error) {}
*/
type GetFA12AllowanceInput struct {
	// Blockhash is the block height at which to make the query. Can leave blank if using Cycle.
	Blockhash string
	// Cycle is the cycle in which to make the query. Can leave blank if using Blockhash.
	Cycle int
	// ChainID is the Chain ID of the chain you want to query
	ChainID string `validate:"required"`
	// Source to form the contents with. The operation is not forged or injected so it is possible for XTZ to be spent.
	Source string `validate:"required"`
	// FA12Contract address of the FA1.2 Contract you wish to query.
	FA12Contract string `validate:"required"`
	// OwnerAddress is the address to get the balance for in the FA1.2 contract
	OwnerAddress string `validate:"required"`
	// SpenderAddress is the address to check an allowance for on behalf of an owner
	SpenderAddress string `validate:"required"`
	// If true the function will use an intermediate contract deployed on Carthagenet, default mainnet.
	Testnet bool
	// If provided this will be the contract view address used to query the FA1.2 contract
	ContractViewAddress string
}

func (g *GetFA12AllowanceInput) validate() error {
	if g.Blockhash == "" && g.Cycle == 0 {
		return errors.New("invalid input: missing key cycle or blockhash")
	} else if g.Blockhash != "" && g.Cycle != 0 {
		return errors.New("invalid input: cannot have both cycle and blockhash")
	}

	err := validator.New().Struct(g)
	if err != nil {
		return errors.Wrap(err, "invalid input")
	}

	return nil
}

var (
	regexTokenAddress        = regexp.MustCompile(`\$token_address`)
	regexOwnerAddress        = regexp.MustCompile(`\$owner`)
	regexSpenderAddress      = regexp.MustCompile(`\$spender`)
	regexContractViewAddress = regexp.MustCompile(`\$contract_view_address`)
	regexBalance             = regexp.MustCompile(`"int":"([0-9]+)"`)
)

const (
	defaultBalanceArgs     = `[{"prim":"PUSH","args":[{"prim":"mutez"},{"int":"0"}]},{"prim":"NONE","args":[{"prim":"key_hash"}]},{"prim":"CREATE_CONTRACT","args":[[{"prim":"parameter","args":[{"prim":"nat"}]},{"prim":"storage","args":[{"prim":"unit"}]},{"prim":"code","args":[[{"prim":"FAILWITH"}]]}]]},{"prim":"DIP","args":[[{"prim":"DIP","args":[[{"prim":"LAMBDA","args":[{"prim":"pair","args":[{"prim":"address"},{"prim":"unit"}]},{"prim":"pair","args":[{"prim":"list","args":[{"prim":"operation"}]},{"prim":"unit"}]},[{"prim":"CAR"},{"prim":"CONTRACT","args":[{"prim":"nat"}]},{"prim":"IF_NONE","args":[[{"prim":"PUSH","args":[{"prim":"string"},{"string":"a"}]},{"prim":"FAILWITH"}],[]]},{"prim":"PUSH","args":[{"prim":"address"},{"string":"$owner"}]},{"prim":"PAIR"},{"prim":"DIP","args":[[{"prim":"PUSH","args":[{"prim":"address"},{"string":"$token_address"}]},{"prim":"CONTRACT","args":[{"prim":"pair","args":[{"prim":"address"},{"prim":"contract","args":[{"prim":"nat"}]}]}],"annots":["%getBalance"]},{"prim":"IF_NONE","args":[[{"prim":"PUSH","args":[{"prim":"string"},{"string":"b"}]},{"prim":"FAILWITH"}],[]]},{"prim":"PUSH","args":[{"prim":"mutez"},{"int":"0"}]}]]},{"prim":"TRANSFER_TOKENS"},{"prim":"DIP","args":[[{"prim":"NIL","args":[{"prim":"operation"}]}]]},{"prim":"CONS"},{"prim":"DIP","args":[[{"prim":"UNIT"}]]},{"prim":"PAIR"}]]}]]},{"prim":"APPLY"},{"prim":"DIP","args":[[{"prim":"PUSH","args":[{"prim":"address"},{"string":"$contract_view_address"}]},{"prim":"CONTRACT","args":[{"prim":"lambda","args":[{"prim":"unit"},{"prim":"pair","args":[{"prim":"list","args":[{"prim":"operation"}]},{"prim":"unit"}]}]}]},{"prim":"IF_NONE","args":[[{"prim":"PUSH","args":[{"prim":"string"},{"string":"c"}]},{"prim":"FAILWITH"}],[]]},{"prim":"PUSH","args":[{"prim":"mutez"},{"int":"0"}]}]]},{"prim":"TRANSFER_TOKENS"},{"prim":"DIP","args":[[{"prim":"NIL","args":[{"prim":"operation"}]}]]},{"prim":"CONS"}]]},{"prim":"CONS"},{"prim":"DIP","args":[[{"prim":"UNIT"}]]},{"prim":"PAIR"}]`
	defaultTotalSupplyArgs = `[{"prim":"PUSH","args":[{"prim":"mutez"},{"int":"0"}]},{"prim":"NONE","args":[{"prim":"key_hash"}]},{"prim":"CREATE_CONTRACT","args":[[{"prim":"parameter","args":[{"prim":"nat"}]},{"prim":"storage","args":[{"prim":"unit"}]},{"prim":"code","args":[[{"prim":"FAILWITH"}]]}]]},{"prim":"DIP","args":[[{"prim":"DIP","args":[[{"prim":"LAMBDA","args":[{"prim":"pair","args":[{"prim":"address"},{"prim":"unit"}]},{"prim":"pair","args":[{"prim":"list","args":[{"prim":"operation"}]},{"prim":"unit"}]},[{"prim":"CAR"},{"prim":"CONTRACT","args":[{"prim":"nat"}]},{"prim":"IF_NONE","args":[[{"prim":"PUSH","args":[{"prim":"string"},{"string":"a"}]},{"prim":"FAILWITH"}],[]]},{"prim":"PUSH","args":[{"prim":"unit"},{"prim":"Unit"}]},{"prim":"PAIR"},{"prim":"DIP","args":[[{"prim":"PUSH","args":[{"prim":"address"},{"string":"$token_address"}]},{"prim":"CONTRACT","args":[{"prim":"pair","args":[{"prim":"unit"},{"prim":"contract","args":[{"prim":"nat"}]}]}],"annots":["%getTotalSupply"]},{"prim":"IF_NONE","args":[[{"prim":"PUSH","args":[{"prim":"string"},{"string":"b"}]},{"prim":"FAILWITH"}],[]]},{"prim":"PUSH","args":[{"prim":"mutez"},{"int":"0"}]}]]},{"prim":"TRANSFER_TOKENS"},{"prim":"DIP","args":[[{"prim":"NIL","args":[{"prim":"operation"}]}]]},{"prim":"CONS"},{"prim":"DIP","args":[[{"prim":"UNIT"}]]},{"prim":"PAIR"}]]}]]},{"prim":"APPLY"},{"prim":"DIP","args":[[{"prim":"PUSH","args":[{"prim":"address"},{"string":"$contract_view_address"}]},{"prim":"CONTRACT","args":[{"prim":"lambda","args":[{"prim":"unit"},{"prim":"pair","args":[{"prim":"list","args":[{"prim":"operation"}]},{"prim":"unit"}]}]}]},{"prim":"IF_NONE","args":[[{"prim":"PUSH","args":[{"prim":"string"},{"string":"c"}]},{"prim":"FAILWITH"}],[]]},{"prim":"PUSH","args":[{"prim":"mutez"},{"int":"0"}]}]]},{"prim":"TRANSFER_TOKENS"},{"prim":"DIP","args":[[{"prim":"NIL","args":[{"prim":"operation"}]}]]},{"prim":"CONS"}]]},{"prim":"CONS"},{"prim":"DIP","args":[[{"prim":"UNIT"}]]},{"prim":"PAIR"}]`
	defaultAllowanceArgs   = `[{"prim":"PUSH","args":[{"prim":"mutez"},{"int":"0"}]},{"prim":"NONE","args":[{"prim":"key_hash"}]},{"prim":"CREATE_CONTRACT","args":[[{"prim":"parameter","args":[{"prim":"nat"}]},{"prim":"storage","args":[{"prim":"unit"}]},{"prim":"code","args":[[{"prim":"FAILWITH"}]]}]]},{"prim":"DIP","args":[[{"prim":"DIP","args":[[{"prim":"LAMBDA","args":[{"prim":"pair","args":[{"prim":"address"},{"prim":"unit"}]},{"prim":"pair","args":[{"prim":"list","args":[{"prim":"operation"}]},{"prim":"unit"}]},[{"prim":"CAR"},{"prim":"CONTRACT","args":[{"prim":"nat"}]},{"prim":"IF_NONE","args":[[{"prim":"PUSH","args":[{"prim":"string"},{"string":"a"}]},{"prim":"FAILWITH"}],[]]},{"prim":"PUSH","args":[{"prim":"address"},{"string":"$spender"}]},{"prim":"PUSH","args":[{"prim":"address"},{"string":"$owner"}]},{"prim":"PAIR"},{"prim":"PAIR"},{"prim":"DIP","args":[[{"prim":"PUSH","args":[{"prim":"address"},{"string":"$token_address"}]},{"prim":"CONTRACT","args":[{"prim":"pair","args":[{"prim":"pair","args":[{"prim":"address"},{"prim":"address"}]},{"prim":"contract","args":[{"prim":"nat"}]}]}],"annots":["%getAllowance"]},{"prim":"IF_NONE","args":[[{"prim":"PUSH","args":[{"prim":"string"},{"string":"b"}]},{"prim":"FAILWITH"}],[]]},{"prim":"PUSH","args":[{"prim":"mutez"},{"int":"0"}]}]]},{"prim":"TRANSFER_TOKENS"},{"prim":"DIP","args":[[{"prim":"NIL","args":[{"prim":"operation"}]}]]},{"prim":"CONS"},{"prim":"DIP","args":[[{"prim":"UNIT"}]]},{"prim":"PAIR"}]]}]]},{"prim":"APPLY"},{"prim":"DIP","args":[[{"prim":"PUSH","args":[{"prim":"address"},{"string":"$contract_view_address"}]},{"prim":"CONTRACT","args":[{"prim":"lambda","args":[{"prim":"unit"},{"prim":"pair","args":[{"prim":"list","args":[{"prim":"operation"}]},{"prim":"unit"}]}]}]},{"prim":"IF_NONE","args":[[{"prim":"PUSH","args":[{"prim":"string"},{"string":"c"}]},{"prim":"FAILWITH"}],[]]},{"prim":"PUSH","args":[{"prim":"mutez"},{"int":"0"}]}]]},{"prim":"TRANSFER_TOKENS"},{"prim":"DIP","args":[[{"prim":"NIL","args":[{"prim":"operation"}]}]]},{"prim":"CONS"}]]},{"prim":"CONS"},{"prim":"DIP","args":[[{"prim":"UNIT"}]]},{"prim":"PAIR"}]`
)

func newBalanceArgs(contractViewAddress, tokenAddress, ownerAddress string) []byte {
	args := regexOwnerAddress.ReplaceAllString(defaultBalanceArgs, ownerAddress)
	args = regexTokenAddress.ReplaceAllString(args, tokenAddress)
	args = regexContractViewAddress.ReplaceAllString(args, contractViewAddress)

	return []byte(args)
}

func newTotalSupplyArgs(contractViewAddress, tokenAddress string) []byte {
	args := regexTokenAddress.ReplaceAllString(defaultTotalSupplyArgs, tokenAddress)
	args = regexContractViewAddress.ReplaceAllString(args, contractViewAddress)

	return []byte(args)
}

func newAllowanceArgs(contractViewAddress, tokenAddress, ownerAddress, spenderAddress string) []byte {
	args := regexOwnerAddress.ReplaceAllString(defaultAllowanceArgs, ownerAddress)
	args = regexTokenAddress.ReplaceAllString(args, tokenAddress)
	args = regexContractViewAddress.ReplaceAllString(args, contractViewAddress)
	args = regexSpenderAddress.ReplaceAllString(args, spenderAddress)

	return []byte(args)
}

func parseBalance(operation Operations) (string, error) {
	if len(operation.Contents) > 0 {
		transaction := operation.Contents[0].ToTransaction()
		if transaction.Metadata != nil {
			if len(transaction.Metadata.InternalOperationResults) >= 4 {
				internalOperationResult := transaction.Metadata.InternalOperationResults[3]
				if len(internalOperationResult.Result.Errors) >= 2 {
					operationErr := internalOperationResult.Result.Errors[1]
					val := string([]byte(*operationErr.With))
					ints := regexBalance.FindStringSubmatch(val)
					if len(ints) == 2 {
						return ints[1], nil
					}
				}
			}

		}
	}

	return "0", errors.New("failed to parse balance from response")
}

func parseSupply(operation Operations) (string, error) {
	if len(operation.Contents) > 0 {
		transaction := operation.Contents[0].ToTransaction()
		if transaction.Metadata != nil {
			if len(transaction.Metadata.InternalOperationResults) >= 4 {
				internalOperationResult := transaction.Metadata.InternalOperationResults[3]
				if len(internalOperationResult.Result.Errors) >= 2 {
					operationErr := internalOperationResult.Result.Errors[1]
					val := string([]byte(*operationErr.With))
					ints := regexBalance.FindStringSubmatch(val)
					if len(ints) == 2 {
						return ints[1], nil
					}
				}
			}

		}
	}

	return "0", errors.New("failed to parse supply from response")
}

func parseAllowance(operation Operations) (string, error) {
	if len(operation.Contents) > 0 {
		transaction := operation.Contents[0].ToTransaction()
		if transaction.Metadata != nil {
			if len(transaction.Metadata.InternalOperationResults) >= 4 {
				internalOperationResult := transaction.Metadata.InternalOperationResults[3]
				if len(internalOperationResult.Result.Errors) >= 2 {
					operationErr := internalOperationResult.Result.Errors[1]
					val := string([]byte(*operationErr.With))
					ints := regexBalance.FindStringSubmatch(val)
					if len(ints) == 2 {
						return ints[1], nil
					}
				}
			}

		}
	}

	return "0", errors.New("failed to parse allowance from response")
}

/*
GetFA12Balance is a helper function to get the balance of a participant in an FA1.2 contracts.
There isn't really a good way to get the balance naturally because the FA1.2 contract entrypoints
are meant to be called from another contract. As a result of this this function will run an operation
that calls an intermediary contract which calls the FA1.2 contract and parses the result.

See: https://gitlab.com/camlcase-dev/dexter-integration/-/blob/master/call_fa1.2_view_entrypoints.md


*/
func (c *Client) GetFA12Balance(input GetFA12BalanceInput) (string, error) {
	err := input.validate()
	if err != nil {
		return "0", errors.Wrapf(err, "could not get fa1.2 balance for '%s' in contract '%s'", input.OwnerAddress, input.FA12Contract)
	}

	if input.Cycle != 0 {
		snapshot, err := c.Cycle(input.Cycle)
		if err != nil {
			return "0", errors.Wrapf(err, "could not get fa1.2 balance for '%s' in contract '%s'", input.OwnerAddress, input.FA12Contract)
		}

		input.Blockhash = snapshot.BlockHash
	}

	counter, err := c.Counter(input.Blockhash, input.Source)
	if err != nil {
		return "0", errors.Wrapf(err, "could not get fa1.2 balance for '%s' in contract '%s'", input.OwnerAddress, input.FA12Contract)
	}
	counter++

	if input.ContractViewAddress == "" {
		if !input.Testnet {
			return "0", errors.Wrapf(errors.New("mainnet not supported yet"), "could not get fa1.2 balance for '%s' in contract '%s'", input.OwnerAddress, input.FA12Contract)
		}
		input.ContractViewAddress = "KT1Njyz94x2pNJGh5uMhKj24VB9JsGCdkySN"
	}

	parameters := json.RawMessage(newBalanceArgs(input.ContractViewAddress, input.FA12Contract, input.OwnerAddress))
	contents := Contents{
		{
			Kind:         TRANSACTION,
			Source:       input.Source,
			Destination:  input.ContractViewAddress,
			Fee:          "0",
			GasLimit:     "1040000",
			StorageLimit: "60000",
			Amount:       "0",
			Counter:      strconv.Itoa(counter),
			Parameters: &Parameters{
				Entrypoint: "default",
				Value:      &parameters,
			},
		},
	}

	operation, err := c.RunOperation(RunOperationInput{
		Blockhash: input.Blockhash,
		Operation: RunOperation{
			Operation: Operations{
				Branch:    input.Blockhash,
				Contents:  contents,
				Signature: "edsigtXomBKi5CTRf5cjATJWSyaRvhfYNHqSUGrn4SdbYRcGwQrUGjzEfQDTuqHhuA8b2d8NarZjz8TRf65WkpQmo423BtomS8Q", // no validation on sig for this func
			},
			ChainID: input.ChainID,
		},
	})
	if err != nil {
		return "0", errors.Wrapf(err, "could not get fa1.2 balance for '%s' in contract '%s'", input.OwnerAddress, input.FA12Contract)
	}

	balance, err := parseBalance(operation)
	if err != nil {
		return "0", errors.Wrapf(err, "could not get fa1.2 balance for '%s' in contract '%s'", input.OwnerAddress, input.FA12Contract)
	}

	return balance, nil
}

/*
GetFA12Supply is a helper function to get the total supply of an FA1.2 contract.
There isn't really a good way to get the supply naturally because the FA1.2 contract entrypoints
are meant to be called from another contract. As a result of this this function will run an operation
that calls an intermediary contract which calls the FA1.2 contract and parses the result.

See: https://gitlab.com/camlcase-dev/dexter-integration/-/blob/master/call_fa1.2_view_entrypoints.md


*/
func (c *Client) GetFA12Supply(input GetFA12SupplyInput) (string, error) {
	err := validator.New().Struct(input)
	if err != nil {
		return "0", errors.Wrap(err, "invalid input")
	}

	if input.Cycle != 0 {
		snapshot, err := c.Cycle(input.Cycle)
		if err != nil {
			return "0", errors.Wrapf(err, "could not get fa1.2 supply for contract '%s'", input.FA12Contract)
		}

		input.Blockhash = snapshot.BlockHash
	}

	counter, err := c.Counter(input.Blockhash, input.Source)
	if err != nil {
		return "0", err
	}
	counter++

	if input.ContractViewAddress == "" {
		if !input.Testnet {
			return "0", errors.New("mainnet not supported yet")
		}
		input.ContractViewAddress = "KT1Njyz94x2pNJGh5uMhKj24VB9JsGCdkySN"
	}

	parameters := json.RawMessage(newTotalSupplyArgs(input.ContractViewAddress, input.FA12Contract))
	contents := Contents{
		{
			Kind:         TRANSACTION,
			Source:       input.Source,
			Destination:  input.ContractViewAddress,
			Fee:          "0",
			GasLimit:     "1040000",
			StorageLimit: "60000",
			Amount:       "0",
			Counter:      strconv.Itoa(counter),
			Parameters: &Parameters{
				Entrypoint: "default",
				Value:      &parameters,
			},
		},
	}

	operation, err := c.RunOperation(RunOperationInput{
		Blockhash: input.Blockhash,
		Operation: RunOperation{
			Operation: Operations{
				Branch:    input.Blockhash,
				Contents:  contents,
				Signature: "edsigtXomBKi5CTRf5cjATJWSyaRvhfYNHqSUGrn4SdbYRcGwQrUGjzEfQDTuqHhuA8b2d8NarZjz8TRf65WkpQmo423BtomS8Q",
			},
			ChainID: input.ChainID,
		},
	})
	if err != nil {
		return "0", err
	}

	return parseSupply(operation)
}

/*
GetFA12Allowance is a helper function to get the allowance of an FA1.2 contract.
There isn't really a good way to get the allowance naturally because the FA1.2 contract entrypoints
are meant to be called from another contract. As a result of this this function will run an operation
that calls an intermediary contract which calls the FA1.2 contract and parses the result.

See: https://gitlab.com/camlcase-dev/dexter-integration/-/blob/master/call_fa1.2_view_entrypoints.md


*/
func (c *Client) GetFA12Allowance(input GetFA12AllowanceInput) (string, error) {
	err := validator.New().Struct(input)
	if err != nil {
		return "0", errors.Wrap(err, "invalid input")
	}

	if input.Cycle != 0 {
		snapshot, err := c.Cycle(input.Cycle)
		if err != nil {
			return "0", errors.Wrapf(err, "could not get fa1.2 supply for contract '%s'", input.FA12Contract)
		}

		input.Blockhash = snapshot.BlockHash
	}

	counter, err := c.Counter(input.Blockhash, input.Source)
	if err != nil {
		return "0", err
	}
	counter++

	if input.ContractViewAddress == "" {
		if !input.Testnet {
			return "0", errors.New("mainnet not supported yet")
		}
		input.ContractViewAddress = "KT1Njyz94x2pNJGh5uMhKj24VB9JsGCdkySN"
	}

	parameters := json.RawMessage(newAllowanceArgs(input.ContractViewAddress, input.FA12Contract, input.OwnerAddress, input.SpenderAddress))
	contents := Contents{
		{
			Kind:         TRANSACTION,
			Source:       input.Source,
			Destination:  input.ContractViewAddress,
			Fee:          "0",
			GasLimit:     "1040000",
			StorageLimit: "60000",
			Amount:       "0",
			Counter:      strconv.Itoa(counter),
			Parameters: &Parameters{
				Entrypoint: "default",
				Value:      &parameters,
			},
		},
	}

	operation, err := c.RunOperation(RunOperationInput{
		Blockhash: input.Blockhash,
		Operation: RunOperation{
			Operation: Operations{
				Branch:    input.Blockhash,
				Contents:  contents,
				Signature: "edsigtXomBKi5CTRf5cjATJWSyaRvhfYNHqSUGrn4SdbYRcGwQrUGjzEfQDTuqHhuA8b2d8NarZjz8TRf65WkpQmo423BtomS8Q",
			},
			ChainID: input.ChainID,
		},
	})
	if err != nil {
		return "0", err
	}

	return parseAllowance(operation)
}
