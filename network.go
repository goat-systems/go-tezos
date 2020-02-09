package gotezos

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

/*
Version Result
RPC: /network/version (GET)
Link: https://tezos.gitlab.io/api/rpc.html#get-network-version
*/
type Version struct {
	ChainName            string `json:"chain_name"`
	DistributedDbVersion int    `json:"distributed_db_version"`
	P2PVersion           int    `json:"p2p_version"`
}

/*
Constants Result
RPC: ../<block_id>/context/constants (GET)
Link: https://tezos.gitlab.io/api/rpc.html#get-block-id-context-constants
*/
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

// Cycle is a Snapshot returned by the Tezos RPC API.
type Cycle struct {
	RandomSeed   string `json:"random_seed"`
	RollSnapshot int    `json:"roll_snapshot"`
	BlockHash    string `json:"-"`
}

/*
Connections Result
RPC: /network/connections (GET)
Link: https://tezos.gitlab.io/api/rpc.html#get-network-connections
*/
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

/*
Bootstrap Result
RPC: /monitor/bootstrapped (GET)
Link: https://tezos.gitlab.io/api/rpc.html#get-monitor-bootstrapped
*/
type Bootstrap struct {
	Block     string    `json:"block"`
	Timestamp time.Time `json:"timestamp"`
}

/*
Version RPC
Path: /network/version (GET)
Link: https://tezos.gitlab.io/api/rpc.html#get-network-version
Description: Supported network layer version.
*/
func (t *GoTezos) Version() (*Version, error) {
	resp, err := t.get("/network/version")
	if err != nil {
		return &Version{}, errors.Wrap(err, "could not get network version")
	}

	var version Version
	err = json.Unmarshal(resp, &version)
	if err != nil {
		return &Version{}, errors.Wrap(err, "could not unmarshal network version")
	}

	return &version, nil
}

/*
Version RPC
Path: ../<block_id>/context/constants (GET)
Link: https://tezos.gitlab.io/api/rpc.html#get-block-id-context-constants
Description: All constants.
*/
func (t *GoTezos) Constants(blockhash string) (*Constants, error) {
	resp, err := t.get(fmt.Sprintf("/chains/main/blocks/%s/context/constants", blockhash))
	if err != nil {
		return &Constants{}, errors.Wrapf(err, "could not get network constants")
	}

	var constants Constants
	err = json.Unmarshal(resp, &constants)
	if err != nil {
		return &constants, errors.Wrapf(err, "could not unmarshal network constants")
	}

	return &constants, nil
}

/*
Connections RPC
Path: /network/connections (GET)
Link: https://tezos.gitlab.io/api/rpc.html#get-network-connections
Description: List the running P2P connection.
*/
func (t *GoTezos) Connections() (*Connections, error) {
	resp, err := t.get("/network/connections")
	if err != nil {
		return &Connections{}, errors.Wrapf(err, "could not get network connections")
	}

	var connections Connections
	err = json.Unmarshal(resp, &connections)
	if err != nil {
		return &Connections{}, errors.Wrapf(err, "could not unmarshal network connections")
	}

	return &connections, nil
}

/*
Bootstrap RPC
Path: /monitor/bootstrapped (GET)
Link: https://tezos.gitlab.io/api/rpc.html#get-monitor-bootstrapped
Description: Wait for the node to have synchronized its chain with a few peers (configured by the node's administrator),
streaming head updates that happen during the bootstrapping process, and closing the stream at the end. If the node was
already bootstrapped, returns the current head immediately.
*/
func (t *GoTezos) Bootstrap() (*Bootstrap, error) {
	resp, err := t.get("/monitor/bootstrapped")
	if err != nil {
		return &Bootstrap{}, errors.Wrap(err, "could not get bootstrap")
	}

	var bootstrap Bootstrap
	err = json.Unmarshal(resp, &bootstrap)
	if err != nil {
		return &bootstrap, errors.Wrap(err, "could not unmarshal bootstrap")
	}

	return &bootstrap, nil
}

/*
Commit RPC
Path: /monitor/commit_hash (GET)
Link: https://tezos.gitlab.io/api/rpc.html#get-monitor-commit-hash
Description: Get information on the build of the node.
*/
func (t *GoTezos) Commit() (*string, error) {
	resp, err := t.get("/monitor/commit_hash")
	if err != nil {
		return nil, errors.Wrap(err, "could not get commit hash")
	}

	var commit string
	err = json.Unmarshal(resp, &commit)
	if err != nil {
		return nil, errors.Wrap(err, "could unmarshal commit")
	}

	return &commit, nil
}

/*
Cycle RPC
Path: ../context/raw/json/cycle/%d" (GET)
Link: https://tezos.gitlab.io/api/rpc.html#get-block-id-context-raw-bytes
Description: Gets information about a tezos snapshot or cycle.
*/
func (t *GoTezos) Cycle(cycle int) (*Cycle, error) {
	head, err := t.Head()
	if err != nil {
		return &Cycle{}, errors.Wrapf(err, "could not get cycle '%d'", cycle)
	}

	if cycle > head.Metadata.Level.Cycle+t.networkConstants.PreservedCycles-1 {
		return &Cycle{}, errors.Errorf("could not get cycle '%d': request is in the future", cycle)
	}

	var c Cycle
	if cycle < head.Metadata.Level.Cycle {
		block, err := t.Block(cycle*t.networkConstants.BlocksPerCycle + 1)
		if err != nil {
			return &Cycle{}, errors.Wrapf(err, "could not get cycle '%d'", cycle)
		}
		c, err = t.getCycleAtHash(block.Hash, cycle)
		if err != nil {
			return &Cycle{}, errors.Wrapf(err, "could not get cycle '%d'", cycle)
		}

	} else {
		var err error
		c, err = t.getCycleAtHash(head.Hash, cycle)
		if err != nil {
			return &Cycle{}, errors.Wrapf(err, "could not get cycle '%d'", cycle)
		}
	}

	level := ((cycle - t.networkConstants.PreservedCycles - 2) * t.networkConstants.BlocksPerCycle) + (c.RollSnapshot+1)*t.networkConstants.BlocksPerRollSnapshot
	if level < 1 {
		level = 1
	}

	block, err := t.Block(level)
	if err != nil {
		return &c, errors.Wrapf(err, "could not get cycle '%d'", cycle)
	}

	c.BlockHash = block.Hash
	return &c, nil
}

func (t *GoTezos) getCycleAtHash(blockhash string, cycle int) (Cycle, error) {
	resp, err := t.get(fmt.Sprintf("/chains/main/blocks/%s/context/raw/json/cycle/%d", blockhash, cycle))
	if err != nil {
		return Cycle{}, errors.Wrapf(err, "could not get cycle at hash '%s'", blockhash)
	}

	var c Cycle
	err = json.Unmarshal(resp, &c)
	if err != nil {
		return c, errors.Wrapf(err, "could not unmarshal at cycle hash '%s'", blockhash)
	}

	return c, nil
}
