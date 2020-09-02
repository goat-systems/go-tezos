package gotezos

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

/*
Version represents the Version RPC.

RPC:
	/network/version (GET)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-network-version
*/
type Version struct {
	ChainName            string `json:"chain_name"`
	DistributedDbVersion int    `json:"distributed_db_version"`
	P2PVersion           int    `json:"p2p_version"`
}

/*
Constants represents the constants RPC.

RPC:
	../<block_id>/context/constants (GET)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-context-constants
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
	HardGasLimitPerOperation     int      `json:"hard_gas_limit_per_operation,string"`
	HardGasLimitPerBlock         int      `json:"hard_gas_limit_per_block,string"`
	ProofOfWorkThreshold         string   `json:"proof_of_work_threshold"`
	TokensPerRoll                string   `json:"tokens_per_roll"`
	MichelsonMaximumTypeSize     int      `json:"michelson_maximum_type_size"`
	SeedNonceRevelationTip       string   `json:"seed_nonce_revelation_tip"`
	OriginationSize              int      `json:"origination_size"`
	BlockSecurityDeposit         int      `json:"block_security_deposit,string"`
	EndorsementSecurityDeposit   int      `json:"endorsement_security_deposit,string"`
	BlockReward                  []int    `json:"block_reward"`
	EndorsementReward            []int    `json:"endorsement_reward"`
	CostPerByte                  int      `json:"cost_per_byte,string"`
	HardStorageLimitPerOperation int      `json:"hard_storage_limit_per_operation,string"`
}

/*
Cycle represents the cycle RPC.

RPC:
	../blocks/<block_id>/context/raw/json/cycle/<cycle_number> (GET)
*/
type Cycle struct {
	RandomSeed   string `json:"random_seed"`
	RollSnapshot int    `json:"roll_snapshot"`
	BlockHash    string `json:"-"`
}

/*
Connections represents the connections RPC.

RPC:
	/network/connections (GET)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-network-connections
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
Bootstrap represents the bootstrap RPC.

RPC:
	/monitor/bootstrapped (GET)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-monitor-bootstrapped
*/
type Bootstrap struct {
	Block     string    `json:"block"`
	Timestamp time.Time `json:"timestamp"`
}

/*
ActiveChains represents the active chains RPC.

RPC:
	/monitor/active_chains (GET)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-monitor-active-chains
*/
type ActiveChains []struct {
	ChainID        string    `json:"chain_id"`
	TestProtocol   string    `json:"test_protocol"`
	ExpirationDate time.Time `json:"expiration_date"`
	Stopping       string    `json:"stopping"`
}

/*
Version gets supported network layer version.

Path:
	/network/version (GET)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-network-version
*/
func (t *GoTezos) Version() (Version, error) {
	resp, err := t.get("/network/version")
	if err != nil {
		return Version{}, errors.Wrap(err, "could not get network version")
	}

	var version Version
	err = json.Unmarshal(resp, &version)
	if err != nil {
		return Version{}, errors.Wrap(err, "could not unmarshal network version")
	}

	return version, nil
}

/*
Constants gets all constants.

Path:
	../<block_id>/context/constants (GET)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-context-constants
*/
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

/*
Connections lists the running P2P connection.

Path:
	/network/connections (GET)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-network-connections
*/
func (t *GoTezos) Connections() (Connections, error) {
	resp, err := t.get("/network/connections")
	if err != nil {
		return Connections{}, errors.Wrapf(err, "could not get network connections")
	}

	var connections Connections
	err = json.Unmarshal(resp, &connections)
	if err != nil {
		return Connections{}, errors.Wrapf(err, "could not unmarshal network connections")
	}

	return connections, nil
}

/*
Bootstrap waits for the node to have synchronized its chain with a few peers (configured by the node's administrator),
streaming head updates that happen during the bootstrapping process, and closing the stream at the end. If the node was
already bootstrapped, returns the current head immediately.

Path:
	/monitor/bootstrapped (GET)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-monitor-bootstrapped
*/
func (t *GoTezos) Bootstrap() (Bootstrap, error) {
	resp, err := t.get("/monitor/bootstrapped")
	if err != nil {
		return Bootstrap{}, errors.Wrap(err, "could not get bootstrap")
	}

	var bootstrap Bootstrap
	err = json.Unmarshal(resp, &bootstrap)
	if err != nil {
		return bootstrap, errors.Wrap(err, "could not unmarshal bootstrap")
	}

	return bootstrap, nil
}

/*
Commit gets information on the build of the node.

Path:
	/monitor/commit_hash (GET)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-monitor-commit-hash
*/
func (t *GoTezos) Commit() (string, error) {
	resp, err := t.get("/monitor/commit_hash")
	if err != nil {
		return "", errors.Wrap(err, "could not get commit hash")
	}

	var commit string
	err = json.Unmarshal(resp, &commit)
	if err != nil {
		return "", errors.Wrap(err, "could unmarshal commit")
	}

	return commit, nil
}

/*
Cycle gets information about a tezos snapshot or cycle.

Path:
	../context/raw/json/cycle/%d" (GET)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-context-raw-bytes
*/
func (t *GoTezos) Cycle(cycle int) (Cycle, error) {
	head, err := t.Head()
	if err != nil {
		return Cycle{}, errors.Wrapf(err, "could not get cycle '%d'", cycle)
	}

	if cycle > head.Metadata.Level.Cycle+t.networkConstants.PreservedCycles-1 {
		return Cycle{}, errors.Errorf("could not get cycle '%d': request is in the future", cycle)
	}

	var c Cycle
	if cycle < head.Metadata.Level.Cycle {
		block, err := t.Block(cycle*t.networkConstants.BlocksPerCycle + 1)
		if err != nil {
			return Cycle{}, errors.Wrapf(err, "could not get cycle '%d'", cycle)
		}
		c, err = t.getCycleAtHash(block.Hash, cycle)
		if err != nil {
			return Cycle{}, errors.Wrapf(err, "could not get cycle '%d'", cycle)
		}

	} else {
		var err error
		c, err = t.getCycleAtHash(head.Hash, cycle)
		if err != nil {
			return Cycle{}, errors.Wrapf(err, "could not get cycle '%d'", cycle)
		}
	}

	level := ((cycle - t.networkConstants.PreservedCycles - 2) * t.networkConstants.BlocksPerCycle) + (c.RollSnapshot+1)*t.networkConstants.BlocksPerRollSnapshot
	if level < 1 {
		level = 1
	}

	block, err := t.Block(level)
	if err != nil {
		return c, errors.Wrapf(err, "could not get cycle '%d'", cycle)
	}

	c.BlockHash = block.Hash
	return c, nil
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

/*
ActiveChains monitor every chain creation and destruction. Currently active chains will be given as first elements.

Path:
	/monitor/active_chains (GET)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-monitor-active-chains
*/
func (t *GoTezos) ActiveChains() (ActiveChains, error) {
	resp, err := t.get("/monitor/active_chains")
	if err != nil {
		return nil, errors.Wrap(err, "failed to get active chains")
	}

	var activeChains ActiveChains
	err = json.Unmarshal(resp, &activeChains)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal active chains")
	}

	return activeChains, nil
}
