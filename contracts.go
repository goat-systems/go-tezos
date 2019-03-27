package gotezos

// GetContractStorage gets the contract storage for a contract
func (gt *GoTezos) GetContractStorage(contract string) ([]byte, error) {
	get := "/chains/main/blocks/head/context/contracts/" + contract + "/storage"
	resp, err := gt.GetResponse(get, "{}")
	if err != nil {
		return resp.Bytes, err
	}
	return resp.Bytes, nil
}
