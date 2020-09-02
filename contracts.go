package gotezos

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

/*
ContractStorage gets access the data of the contract.

Path:
	../<block_id>/context/contracts/<contract_id>/storage (GET)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-context-contracts-contract-id-storage

Parameters:

	blockhash:
		The hash of block (height) of which you want to make the query.

	KT1:
		The contract address.
*/
func (t *GoTezos) ContractStorage(blockhash string, KT1 string) (MichelineExpression, error) {
	query := fmt.Sprintf("/chains/main/blocks/%s/context/contracts/%s/storage", blockhash, KT1)
	resp, err := t.get(query)
	if err != nil {
		return MichelineExpression{}, errors.Wrap(err, "could not get storage '%s'")
	}

	var micheline MichelineExpression
	err = json.Unmarshal(resp, &micheline)
	if err != nil {
		return micheline, errors.Wrapf(err, "failed to get storage for contract '%s'", KT1)
	}

	return micheline, nil
}
