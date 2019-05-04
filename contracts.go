package gotezos

// ContractService is a struct wrapper for contract functions
type ContractService struct {
	gt *GoTezos
}

// returns a new newContractService
func (gt *GoTezos) newContractService() *ContractService {
	return &ContractService{gt: gt}
}

// GetStorage gets the contract storage for a contract
func (s *ContractService) GetStorage(contract string) ([]byte, error) {
	query := "/chains/main/blocks/head/context/contracts/" + contract + "/storage"
	resp, err := s.gt.Get(query, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
