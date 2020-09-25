package rpc

import (
	"encoding/json"
	"fmt"
	"strconv"

	validator "github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

/*
BalanceInput is the input for the gotezos.Balance function.

Function:
	func (t *GoTezos) Balance(blockhash, address string) (int, error) {}
*/
type BalanceInput struct {
	// The block level of which you want to make the query.
	Blockhash string
	// The delegate that you want to make the query.
	Address string `validate:"required"`
	// The cycle to get the balance at (optional).
	Cycle int
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

Parameters:

	blockhash:
		The hash of block (height) of which you want to make the query.

	address:
		Any tezos public address.
*/
func (c *Client) Balance(input BalanceInput) (int, error) {
	if err := input.validate(); err != nil {
		return 0, errors.Wrapf(err, "could not get balance for '%s'", input.Address)
	}

	var resp []byte
	if input.Cycle != 0 {
		snapshot, err := c.Cycle(input.Cycle)
		if err != nil {
			return 0, errors.Wrapf(err, "could not get balance for '%s' at cycle '%d'", input.Address, input.Cycle)
		}

		input.Blockhash = snapshot.BlockHash
	}

	query := fmt.Sprintf("/chains/%s/blocks/%s/context/contracts/%s/balance", c.chain, input.Blockhash, input.Address)
	resp, err := c.get(query)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get balance")
	}

	var balanceStr string
	if err = json.Unmarshal(resp, &balanceStr); err != nil {
		return 0, errors.Wrap(err, "failed to get balance")
	}

	balance, err := strconv.Atoi(balanceStr)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get balance")
	}

	return balance, nil
}
