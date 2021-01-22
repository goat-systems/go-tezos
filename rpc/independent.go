package rpc

import (
	"encoding/json"
	"fmt"
	"time"

	validator "github.com/go-playground/validator/v10"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

/*
InjectionOperationInput is the input for the InjectionOperation function.

RPC:
	https://tezos.gitlab.io/shell/rpc.html#post-injection-operation
*/
type InjectionOperationInput struct {
	// The operation string.
	Operation string `validate:"required"`

	// If ?async is true, the function returns immediately.
	Async bool

	// Specify the ChainID.
	ChainID string
}

func (i *InjectionOperationInput) contructRPCOptions() []rpcOptions {
	var opts []rpcOptions
	if i.Async {
		opts = append(opts, rpcOptions{
			"async",
			"true",
		})
	}

	if i.ChainID != "" {
		opts = append(opts, rpcOptions{
			"chain_id",
			i.ChainID,
		})
	}
	return opts
}

/*
InjectionOperation injects an operation in node and broadcast it. Returns the ID of the operation.
The `signedOperationContents` should be constructed using a contextual RPCs from the latest block
and signed by the client. By default, the RPC will wait for the operation to be (pre-)validated
before answering. See RPCs under /blocks/prevalidation for more details on the prevalidation context.
If ?async is true, the function returns immediately. Otherwise, the operation will be validated before
the result is returned. An optional ?chain parameter can be used to specify whether to inject on the
test chain or the main chain.

Path:
	/injection/operation (POST)

RPC:
	https://tezos.gitlab.io/shell/rpc.html#post-injection-operation
*/
func (c *Client) InjectionOperation(input InjectionOperationInput) (*resty.Response, string, error) {
	err := validator.New().Struct(input)
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to inject operation: invalid input")
	}

	v, err := json.Marshal(input.Operation)
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to inject operation")
	}
	resp, err := c.post("/injection/operation", v, input.contructRPCOptions()...)
	if err != nil {
		return resp, "", errors.Wrap(err, "failed to inject operation")
	}

	var opstring string
	err = json.Unmarshal(resp.Body(), &opstring)
	if err != nil {
		return resp, "", errors.Wrap(err, "failed to inject operation: failed to parse json")
	}

	return resp, opstring, nil
}

/*
InjectionBlockInput is the input for the InjectionBlock function.

RPC:
	https://tezos.gitlab.io/shell/rpc.html#post-injection-block
*/
type InjectionBlockInput struct {

	// Block header signature
	SignedBlock string `validate:"required"`

	// Operations to include in the block.
	// This is not the same as operations found in mempool and also not like preapply result
	Operations [][]interface{} `validate:"required"`

	// If ?async is true, the function returns immediately.
	Async bool

	// If ?force is true, it will be injected even on non strictly increasing fitness.
	Force bool

	// Specify the ChainID.
	ChainID string
}

func (i *InjectionBlockInput) contructRPCOptions() []rpcOptions {
	var opts []rpcOptions
	if i.Async {
		opts = append(opts, rpcOptions{
			"async",
			"true",
		})
	}

	if i.Force {
		opts = append(opts, rpcOptions{
			"force",
			"true",
		})
	}

	if i.ChainID != "" {
		opts = append(opts, rpcOptions{
			"chain_id",
			i.ChainID,
		})
	}
	return opts
}

/*
InjectionBlock inject a block in the node and broadcast it. The `operations`
embedded in `blockHeader` might be pre-validated using a contextual RPCs
from the latest block (e.g. '/blocks/head/context/preapply'). Returns the
ID of the block. By default, the RPC will wait for the block to be validated
before answering. If ?async is true, the function returns immediately. Otherwise,
the block will be validated before the result is returned. If ?force is true, it
will be injected even on non strictly increasing fitness. An optional ?chain parameter
can be used to specify whether to inject on the test chain or the main chain.

Path:
	/injection/operation (POST)

RPC:
	https://tezos.gitlab.io/shell/rpc.html#post-injection-block
*/
func (c *Client) InjectionBlock(input InjectionBlockInput) (*resty.Response, error) {
	err := validator.New().Struct(input)
	if err != nil {
		return nil, errors.Wrap(err, "failed to inject block: invalid input")
	}

	// Create an anonymous struct containing the data required by RPC
	newBlock := struct {
		SignedBlock string  `json:"data"`
		Ops [][]interface{} `json:"operations"`
	}{
		input.SignedBlock,
		input.Operations,
	}

	v, err := json.Marshal(newBlock)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal new block")
	}

	resp, err := c.post("/injection/block", v, input.contructRPCOptions()...)
	if err != nil {
		return resp, errors.Wrap(err, "failed to inject block")
	}

	return resp, nil
}

/*
Connections is the network connections of a tezos node.

Path:
	/network/connections (GET)

RPC:
	https://tezos.gitlab.io/shell/rpc.html#get-network-connections
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
Connections lists the running P2P connection.

Path:
	/network/connections (GET)

RPC:
	https://tezos.gitlab.io/shell/rpc.html#get-network-connections
*/
func (c *Client) Connections() (*resty.Response, Connections, error) {
	resp, err := c.get("/network/connections")
	if err != nil {
		return resp, Connections{}, errors.Wrapf(err, "failed to get network connections")
	}

	var connections Connections
	err = json.Unmarshal(resp.Body(), &connections)
	if err != nil {
		return resp, Connections{}, errors.Wrapf(err, "failed to get network connections: failed to parse json")
	}

	return resp, connections, nil
}

/*
ActiveChains is the active chains on the tezos network.

RPC:
	https://tezos.gitlab.io/shell/rpc.html#get-monitor-active-chains
*/
type ActiveChains []struct {
	ChainID        string    `json:"chain_id"`
	TestProtocol   string    `json:"test_protocol"`
	ExpirationDate time.Time `json:"expiration_date"`
	Stopping       string    `json:"stopping"`
}

/*
ActiveChains monitor every chain creation and destruction. Currently active chains will be given as first elements.

Path:
	/monitor/active_chains (GET)

RPC:
	https://tezos.gitlab.io/shell/rpc.html#get-monitor-active-chains
*/
func (c *Client) ActiveChains() (*resty.Response, ActiveChains, error) {
	resp, err := c.get("/monitor/active_chains")
	if err != nil {
		return nil, ActiveChains{}, errors.Wrap(err, "failed to get active chains")
	}

	var activeChains ActiveChains
	err = json.Unmarshal(resp.Body(), &activeChains)
	if err != nil {
		return resp, ActiveChains{}, errors.Wrap(err, "failed to get active chains: failed to parse json")
	}

	return resp, activeChains, nil
}

/*
MempoolInput is the input for the goTezos.Mempool function.
Function:
	func (c *Client) Mempool(input *MempoolInput) (Mempool, error) {}
*/
type MempoolInput struct {
	// Mempool filters
	Applied       bool
	BranchDelayed bool
	Refused       bool
	BranchRefused bool
}

/*
Mempool represents the contents of the Tezos mempool.
RPC:
    /chains/<chain_id>/mempool/pending_operations (GET)
*/
type Mempool struct {
	Applied       []Operations    `json:"applied"`
	Refused       []OperationsAlt `json:"refused"`
	BranchRefused []OperationsAlt `json:"branch_refused"`
	BranchDelayed []OperationsAlt `json:"branch_delayed"`
	Unprocessed   []OperationsAlt `json:"unprocessed"`
}

/*
Mempool fetches the current contents of main the chain mempool.
Path:
    /chains/<chain_id>/mempool/pending_operations (GET)
Parameters:
    None
*/
func (c *Client) Mempool(input MempoolInput) (*resty.Response, *Mempool, error) {
	resp, err := c.get(fmt.Sprintf("/chains/%s/mempool/pending_operations", c.chain), input.constructRPCOptions()...)
	if err != nil {
		return resp, &Mempool{}, errors.Wrap(err, "failed to fetch mempool contents")
	}

	var mempool Mempool
	err = json.Unmarshal(resp.Body(), &mempool)
	if err != nil {
		return resp, &mempool, errors.Wrap(err, "failed to unmarshal mempool contents")
	}

	return resp, &mempool, nil
}

func (m *MempoolInput) constructRPCOptions() []rpcOptions {
	var opts []rpcOptions
	if m.Applied {
		opts = append(opts, rpcOptions{
			"applied",
			"true",
		})
	} else {
		opts = append(opts, rpcOptions{
			"applied",
			"false",
		})
	}

	if m.BranchDelayed {
		opts = append(opts, rpcOptions{
			"branch_delayed",
			"true",
		})
	} else {
		opts = append(opts, rpcOptions{
			"branch_delayed",
			"false",
		})
	}

	if m.Refused {
		opts = append(opts, rpcOptions{
			"refused",
			"true",
		})
	} else {
		opts = append(opts, rpcOptions{
			"refused",
			"false",
		})
	}

	if m.BranchRefused {
		opts = append(opts, rpcOptions{
			"branch_refused",
			"true",
		})
	} else {
		opts = append(opts, rpcOptions{
			"branch_refused",
			"false",
		})
	}

	return opts
}
