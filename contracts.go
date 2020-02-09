package gotezos

import (
	"fmt"

	"github.com/pkg/errors"
)

/*
ContractStorage RPC
Path: ../<block_id>/context/contracts/<contract_id>/storage (GET)
Link: https://tezos.gitlab.io/api/rpc.html#get-block-id-context-contracts-contract-id-storage
Description: Access the data of the contract.

Parameters:
	blockhash:
		The hash of block (height) of which you want to make the query.
	KT1:
		The contract address.
*/
func (t *GoTezos) ContractStorage(blockhash string, KT1 string) (*[]byte, error) {
	query := fmt.Sprintf("/chains/main/blocks/%s/context/contracts/%s/storage", blockhash, KT1)
	resp, err := t.get(query)
	if err != nil {
		return &resp, errors.Wrap(err, "could not get storage '%s'")
	}
	return &resp, nil
}
