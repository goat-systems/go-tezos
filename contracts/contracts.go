package contracts

import (
	tzc "github.com/DefinitelyNotAGoat/go-tezos/client"
	"github.com/pkg/errors"
)

// ContractService is a struct wrapper for contract functions
type ContractService struct {
	tzclient tzc.TezosClient
}

// NewContractService returns a new ContractService
func NewContractService(tzclient tzc.TezosClient) *ContractService {
	return &ContractService{tzclient: tzclient}
}

// GetStorage gets the contract storage for a contract
func (s *ContractService) GetStorage(contract string) ([]byte, error) {
	query := "/chains/main/blocks/head/context/contracts/" + contract + "/storage"
	resp, err := s.tzclient.Get(query, nil)
	if err != nil {
		return resp, errors.Wrap(err, "could not get storage '%s'")
	}
	return resp, nil
}
