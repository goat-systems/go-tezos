package rpc

import (
	"encoding/json"
	"fmt"

	validator "github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

/*
BalanceInput is the input for the gotezos.Balance function.

Function:
	func (c *Client) Balance(input BalanceInput) (string, error) {}
*/
type BalanceInput struct {
	// The block level of which you want to make the query. If not provided Cycle is required.
	Blockhash string
	// The cycle to get the balance at. If not provided Blockhash is required.
	Cycle int
	// The delegate that you want to make the query.
	Address string `validate:"required"`
}

func (b *BalanceInput) validate() error {
	if b.Blockhash == "" && b.Cycle == 0 {
		return errors.New("invalid input: missing key cycle or blockhash")
	} else if b.Blockhash != "" && b.Cycle != 0 {
		return errors.New("invalid input: cannot have both cycle and blockhash")
	}

	err := validator.New().Struct(b)
	if err != nil {
		return errors.Wrap(err, "invalid input")
	}

	return nil
}

/*
Balance gives access to the balance of a contract.

Path:
	../<block_id>/context/contracts/<contract_id>/balance (GET)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-context-contracts-contract-id-balance
*/
func (c *Client) Balance(input BalanceInput) (string, error) {
	if err := input.validate(); err != nil {
		return "", errors.Wrapf(err, "could not get balance for '%s'", input.Address)
	}

	var resp []byte
	if input.Cycle != 0 {
		snapshot, err := c.Cycle(input.Cycle)
		if err != nil {
			return "", errors.Wrapf(err, "could not get balance for '%s' at cycle '%d'", input.Address, input.Cycle)
		}

		input.Blockhash = snapshot.BlockHash
	}

	query := fmt.Sprintf("/chains/%s/blocks/%s/context/contracts/%s/balance", c.chain, input.Blockhash, input.Address)
	resp, err := c.get(query)
	if err != nil {
		return "0", errors.Wrap(err, "failed to get balance")
	}

	var balance string
	if err = json.Unmarshal(resp, &balance); err != nil {
		return "0", errors.Wrap(err, "failed to get balance")
	}

	return balance, nil
}
