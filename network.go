package gotezos

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

// NetworkService is wrapper representing network functions
type NetworkService struct {
	gt *GoTezos
}

// NetworkVersion represents the network version returned by the Tezos network.
type NetworkVersion struct {
	Name    string `json:"name"`
	Major   int    `json:"major"`
	Minor   int    `json:"minor"`
	Network string // Human readable network name
}

// NetworkVersions is an array of NetworkVersion
type NetworkVersions []NetworkVersion

// NetworkConstants represents the network constants returned by the Tezos network.
type NetworkConstants struct {
	ProofOfWorkNonceSize         int      `json:"proof_of_work_nonce_size"`
	NonceLength                  int      `json:"nonce_length"`
	MaxRevelationsPerBlock       int      `json:"max_revelations_per_block"`
	MaxOperationDataLength       int      `json:"max_operation_data_length"`
	MaxProposalsPerDelegate      int      `json:"max_proposals_per_delegate"`
	PreservedCycles              int      `json:"preserved_cycles"`
	BlocksPerCycle               int      `json:"blocks_per_cycle"`
	BlocksPerCommitment          int      `json:"blocks_per_commitment"`
	BlocksPerRollSnapshot        int      `json:"blocks_per_roll_snapshot"`
	BlocksPerVotingPeriod        int      `json:"blocks_per_voting_period"`
	TimeBetweenBlocks            []string `json:"time_between_blocks"`
	EndorsersPerBlock            int      `json:"endorsers_per_block"`
	HardGasLimitPerOperation     string   `json:"hard_gas_limit_per_operation"`
	HardGasLimitPerBlock         string   `json:"hard_gas_limit_per_block"`
	ProofOfWorkThreshold         string   `json:"proof_of_work_threshold"`
	TokensPerRoll                string   `json:"tokens_per_roll"`
	MichelsonMaximumTypeSize     int      `json:"michelson_maximum_type_size"`
	SeedNonceRevelationTip       string   `json:"seed_nonce_revelation_tip"`
	OriginationSize              int      `json:"origination_size"`
	BlockSecurityDeposit         string   `json:"block_security_deposit"`
	EndorsementSecurityDeposit   string   `json:"endorsement_security_deposit"`
	BlockReward                  string   `json:"block_reward"`
	EndorsementReward            string   `json:"endorsement_reward"`
	CostPerByte                  string   `json:"cost_per_byte"`
	HardStorageLimitPerOperation string   `json:"hard_storage_limit_per_operation"`
}

// NewNetworkService returns a new NetworkService
func (gt *GoTezos) newNetworkService() *NetworkService {
	return &NetworkService{gt: gt}
}

// GetVersions gets the network versions of Tezos network the client is using.
func (n *NetworkService) GetVersions() ([]NetworkVersion, error) {

	networkVersions := make([]NetworkVersion, 0)

	resp, err := n.gt.Get("/network/versions", nil)
	if err != nil {
		return networkVersions, err
	}

	var nvs NetworkVersions
	nvs, err = nvs.unmarshalJSON(resp)
	if err != nil {
		return networkVersions, err
	}

	// Extract just the network name and append to returning slice
	// 'range' operates on a copy of the struct so cannot update-in-place
	for _, v := range nvs {

		parts := strings.Split(v.Name, "_")
		if len(parts) == 3 {
			v.Network = parts[1]
		}

		networkVersions = append(networkVersions, v)
	}

	return networkVersions, nil
}

// GetConstants gets the network constants for the Tezos network the client is using.
func (n *NetworkService) GetConstants() (NetworkConstants, error) {
	networkConstants := NetworkConstants{}
	resp, err := n.gt.Get("/chains/main/blocks/head/context/constants", nil)
	if err != nil {
		return networkConstants, err
	}
	networkConstants, err = networkConstants.unmarshalJSON(resp)
	if err != nil {
		return networkConstants, err
	}

	return networkConstants, nil
}

// GetChainID gets the id of the chain with the most fitness
func (n *NetworkService) GetChainID() (string, error) {
	query := "/chains/main/chain_id"
	resp, err := n.gt.Get(query, nil)
	if err != nil {
		return "", fmt.Errorf("could not get chain ID: %v", err)
	}

	chainID, err := unmarshalString(resp)
	if err != nil {
		return "", fmt.Errorf("could not get chain ID: %v", err)
	}

	return chainID, nil
}

// UnmarshalJSON unmarshals the bytes received as a parameter, into the type NetworkVersion.
func (nvs *NetworkVersions) unmarshalJSON(v []byte) (NetworkVersions, error) {
	networkVersions := NetworkVersions{}
	err := json.Unmarshal(v, &networkVersions)
	if err != nil {
		return networkVersions, err
	}
	return networkVersions, nil
}

// UnmarshalJSON unmarshals bytes received as a parameter, into the type NetworkConstants.
func (nc *NetworkConstants) unmarshalJSON(v []byte) (NetworkConstants, error) {
	networkConstants := NetworkConstants{}
	err := json.Unmarshal(v, &networkConstants)
	if err != nil {
		log.Println("Could not get unMarshal bytes into NetworkConstants: " + err.Error())
		return networkConstants, err
	}
	return networkConstants, nil
}
