package rpc

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
Cycle represents the cycle RPC.

RPC:
	../blocks/<block_id>/context/raw/json/cycle/<cycle_number> (GET)
*/
type Cycle struct {
	LastRoll     []string `json:"last_roll,omitempty"`
	Nonces       []string `json:"nonces,omitempty"`
	RandomSeed   string   `json:"random_seed"`
	RollSnapshot int      `json:"roll_snapshot"`
	BlockHash    string   `json:"-"`
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
func (c *Client) Version() (Version, error) {
	resp, err := c.get("/network/version")
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
Connections lists the running P2P connection.

Path:
	/network/connections (GET)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-network-connections
*/
func (c *Client) Connections() (Connections, error) {
	resp, err := c.get("/network/connections")
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
func (c *Client) Bootstrap() (Bootstrap, error) {
	resp, err := c.get("/monitor/bootstrapped")
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
func (c *Client) Commit() (string, error) {
	resp, err := c.get("/monitor/commit_hash")
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
func (c *Client) Cycle(cycle int) (Cycle, error) {
	head, err := c.Head()
	if err != nil {
		return Cycle{}, errors.Wrapf(err, "could not get cycle '%d'", cycle)
	}

	if cycle > head.Metadata.Level.Cycle+c.networkConstants.PreservedCycles-1 {
		return Cycle{}, errors.Errorf("could not get cycle '%d': request is in the future", cycle)
	}

	var cyc Cycle
	if cycle < head.Metadata.Level.Cycle {
		block, err := c.Block(cycle*c.networkConstants.BlocksPerCycle + 1)
		if err != nil {
			return Cycle{}, errors.Wrapf(err, "could not get cycle '%d'", cycle)
		}

		cyc, err = c.getCycleAtHash(block.Hash, cycle)
		if err != nil {
			return Cycle{}, errors.Wrapf(err, "could not get cycle '%d'", cycle)
		}

	} else {
		var err error
		cyc, err = c.getCycleAtHash(head.Hash, cycle)
		if err != nil {
			return Cycle{}, errors.Wrapf(err, "could not get cycle '%d'", cycle)
		}
	}

	level := ((cycle - c.networkConstants.PreservedCycles - 2) * c.networkConstants.BlocksPerCycle) + (cyc.RollSnapshot+1)*c.networkConstants.BlocksPerRollSnapshot
	if level < 1 {
		level = 1
	}

	block, err := c.Block(level)
	if err != nil {
		return cyc, errors.Wrapf(err, "could not get cycle '%d'", cycle)
	}

	cyc.BlockHash = block.Hash
	return cyc, nil
}

func (c *Client) getCycleAtHash(blockhash string, cycle int) (Cycle, error) {
	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/context/raw/json/cycle/%d", c.chain, blockhash, cycle))
	if err != nil {
		return Cycle{}, errors.Wrapf(err, "could not get cycle at hash '%s'", blockhash)
	}

	var cyc Cycle
	err = json.Unmarshal(resp, &cyc)
	if err != nil {
		return cyc, errors.Wrapf(err, "could not unmarshal at cycle hash '%s'", blockhash)
	}

	return cyc, nil
}

/*
ActiveChains monitor every chain creation and destruction. Currently active chains will be given as first elements.

Path:
	/monitor/active_chains (GET)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-monitor-active-chains
*/
func (c *Client) ActiveChains() (ActiveChains, error) {
	resp, err := c.get("/monitor/active_chains")
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
