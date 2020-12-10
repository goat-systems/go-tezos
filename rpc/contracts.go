package rpc

import (
	"fmt"

	validator "github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

/*
ScriptExpression is a string that will eventually be forged into a script_expression

Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-context-big-maps-big-map-id-script-expr
*/
type ScriptExpression string

/*
BigMapInput is the input for the goTezos.BigMap function.

Function:
	func (t *GoTezos) BigMap(input BigMapInput) ([]byte, error) {}
*/
type BigMapInput struct {
	Cycle            int
	Blockhash        string
	BigMapID         int              `validate:"required"`
	ScriptExpression ScriptExpression `validate:"required"`
}

func (b *BigMapInput) validate() error {
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
ContractStorageInput is the input for the client.ContractStorage() function.

Function:
	func (c *Client) ContractStorage(input ContractStorageInput) ([]byte, error)  {}
*/
type ContractStorageInput struct {
	Blockhash string `validate:"required"`
	Contract  string `validate:"required"`
}

/*
ContractStorage gets access the data of the contract.

Path:
	../<block_id>/context/contracts/<contract_id>/storage (GET)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-context-contracts-contract-id-storage
*/
func (c *Client) ContractStorage(input ContractStorageInput) ([]byte, error) {
	err := validator.New().Struct(input)
	if err != nil {
		return []byte{}, errors.Wrap(err, "invalid input")
	}

	query := fmt.Sprintf("/chains/%s/blocks/%s/context/contracts/%s/storage", c.chain, input.Blockhash, input.Contract)
	resp, err := c.get(query)
	if err != nil {
		return []byte{}, errors.Wrap(err, "could not get storage '%s'")
	}

	return resp, nil
}

/*
BigMap reads data from a big_map.

Path:
 	../<block_id>/context/big_maps/<big_map_id>/<script_expr> (GET)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-context-big-maps-big-map-id-script-expr
*/
func (c *Client) BigMap(input BigMapInput) ([]byte, error) {
	err := input.validate()
	if err != nil {
		return []byte{}, errors.Wrapf(err, "could not get big map '%d' at cycle '%d'", input.BigMapID, input.Cycle)
	}

	input.Blockhash, err = c.extractBlockHash(input.Cycle, input.Blockhash)
	if err != nil {
		return []byte{}, errors.Wrapf(err, "could not get big map '%d' at cycle '%d'", input.BigMapID, input.Cycle)
	}

	query := fmt.Sprintf("/chains/%s/blocks/%s/context/big_maps/%d/%s", c.chain, input.Blockhash, input.BigMapID, input.ScriptExpression)
	resp, err := c.get(query)
	if err != nil {
		return []byte{}, errors.Wrap(err, "could not get big_map storage '%s'")
	}

	return resp, nil
}
