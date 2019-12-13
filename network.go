package gotezos

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

// Version represents the network version returned by the Tezos network.
type Version struct {
	Name    string `json:"name"`
	Major   int    `json:"major"`
	Minor   int    `json:"minor"`
	Network string // Human readable network name
}

// Versions is an array of Version
type Versions []Version

// Constants represents the network constants returned by the Tezos network.
type Constants struct {
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

// Connections represents network connections
type Connections []struct {
	Incoming bool   `json:"incoming"`
	PeerID   string `json:"peer_id"`
	IDPoint  struct {
		Addr string `json:"addr"`
		Port int    `json:"port"`
	} `json:"id_point"`
	RemoteSocketPort int `json:"remote_socket_port"`
	Versions         []struct {
		Name  string `json:"name"`
		Major int    `json:"major"`
		Minor int    `json:"minor"`
	} `json:"versions"`
	Private       bool `json:"private"`
	LocalMetadata struct {
		DisableMempool bool `json:"disable_mempool"`
		PrivateNode    bool `json:"private_node"`
	} `json:"local_metadata"`
	RemoteMetadata struct {
		DisableMempool bool `json:"disable_mempool"`
		PrivateNode    bool `json:"private_node"`
	} `json:"remote_metadata"`
}

// Versions gets the network versions.
func (t *GoTezos) Versions() (Versions, error) {
	resp, err := t.get("/network/versions")
	if err != nil {
		return []Version{}, errors.Wrap(err, "could not get network versions")
	}

	var versions Versions
	err = json.Unmarshal(resp, &versions)
	if err != nil {
		return versions, errors.Wrap(err, "could not get network versions")
	}

	return versions, nil
}

// Constants gets the network constants
func (t *GoTezos) Constants(blockhash string) (Constants, error) {
	resp, err := t.get(fmt.Sprintf("/chains/main/blocks/%s/context/constants", blockhash))
	if err != nil {
		return Constants{}, errors.Wrapf(err, "could not get network constants")
	}

	var constants Constants
	err = json.Unmarshal(resp, &constants)
	if err != nil {
		return constants, errors.Wrapf(err, "could not unmarshal network constants")
	}

	return constants, nil
}

// ChainID gets the id of the chain with the most fitness
func (t *GoTezos) ChainID() (string, error) {
	resp, err := t.get("/chains/main/chain_id")
	if err != nil {
		return "", errors.Wrapf(err, "could not get chain ID")
	}

	var chainID string
	err = json.Unmarshal(resp, &chainID)
	if err != nil {
		return chainID, errors.Wrapf(err, "could unmarshal chain ID")
	}

	return chainID, nil
}

// NetworkConnections gets the network connections
func (t *GoTezos) NetworkConnections() (Connections, error) {
	resp, err := t.get("/network/connections")
	if err != nil {
		return Connections{}, errors.Wrapf(err, "could not get network connections")
	}

	var connections Connections
	err = json.Unmarshal(resp, &connections)
	if err != nil {
		return connections, errors.Wrapf(err, "could not unmarshal network connections")
	}

	return connections, nil
}
