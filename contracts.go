package gotezos

import (
	"fmt"
	"github.com/pkg/errors"
)

// ContractStorage gets the contract storage for a contract
func (t *GoTezos) ContractStorage(blockhash string, KT1 string) ([]byte, error) {
	query := fmt.Sprintf("/chains/main/blocks/%s/context/contracts/%s/storage", blockhash, KT1)
	resp, err := t.get(query)
	if err != nil {
		return resp, errors.Wrap(err, "could not get storage '%s'")
	}
	return resp, nil
}
