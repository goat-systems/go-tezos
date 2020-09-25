package rpc

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	validator "github.com/go-playground/validator/v10"
	"github.com/goat-systems/go-tezos/v3/internal/crypto"
	"github.com/pkg/errors"
)

/*
InjectionOperationInput is the input for the goTezos.InjectionOperation function.

Function:
	func (t *GoTezos) InjectionOperation(input InjectionOperationInput) ([]byte, error) {}
*/
type InjectionOperationInput struct {
	// The operation string.
	Operation string `validate:"required"`

	// If ?async is true, the function returns immediately.
	Async bool

	// Specify the ChainID.
	ChainID string
}

/*
InjectionBlockInput is the input for the goTezos.InjectionBlock function.

Function:
	func (t *GoTezos) InjectionBlock(input InjectionBlockInput) ([]byte, error) {}
*/
type InjectionBlockInput struct {
	// Block to inject
	Block *Block `validate:"required"`

	// If ?async is true, the function returns immediately.
	Async bool

	// If ?force is true, it will be injected even on non strictly increasing fitness.
	Force bool

	// Specify the ChainID.
	ChainID string
}

/*
RunOperationInput is the input for the rpc.RunOperation function.

Function:
	func (c *Client) RunOperation(input RunOperationInput) (Operations, error)
*/
type RunOperationInput struct {
	Blockhash string       `validate:"required"`
	Operation RunOperation `json:"operation" validate:"required"`
}

// RunOperation -
type RunOperation struct {
	Operation Operations `json:"operation" validate:"required"`
	ChainID   string     `json:"chain_id" validate:"required"`
}

/*
UnforgeOperationWithRPCInput is the input for the goTezos.UnforgeOperationWithRPC function.

Function:
	func (t *GoTezos) UnforgeOperationWithRPC(blockhash string, operation string, checkSignature bool) (Operations, error) {}
*/
type UnforgeOperationWithRPCInput struct {
	Blockhash      string                             `validate:"required"`
	Operations     []UnforgeOperationWithRPCOperation `json:"operations" validate:"required"`
	CheckSignature bool                               `json:"check_signature"`
}

// UnforgeOperationWithRPCOperation -
type UnforgeOperationWithRPCOperation struct {
	Data   string `json:"data" validate:"required"`
	Branch string `json:"branch" validate:"required"`
}

/*
ForgeOperationWithRPCInput is the input for the goTezos.ForgeOperationWithRPC function.

Fields:

	Blockhash:
		The hash of block (height) of which you want to make the query.

	Contents:
		The contents of the of the operation.

	Branch:
		The branch of the operation to be forged.

	CheckRPCAddr:
		Overides the GoTezos client with a new one pointing to a different address. This allows the user to validate the forge against different nodes for security.

Function:
	func (t *GoTezos) ForgeOperationWithRPC(blockhash, branch string, contents ...Contents) (string, error) {}
*/
type ForgeOperationWithRPCInput struct {
	Blockhash    string   `validate:"required"`
	Branch       string   `validate:"required"`
	Contents     Contents `validate:"required"`
	CheckRPCAddr string
}

/*
PreapplyOperationsInput is the input for the PreapplyOperations.

Function:
	func PreapplyOperations(input PreapplyOperationsInput) ([]byte, error) {}
*/
type PreapplyOperationsInput struct {
	Blockhash  string `validate:"required"`
	Operations []Operations
}

/*
PreapplyOperations simulates the validation of an operation.

Path:
	../<block_id>/helpers/preapply/operations (POST)

Link:
	https://tezos.gitlab.io/api/rpc.html#post-block-id-helpers-preapply-operations

Parameters:

	input:
		PreapplyOperationsInput contains the blockhash, protocol, signature, and operation contents needed to fufill this RPC.
*/
func (c *Client) PreapplyOperations(input PreapplyOperationsInput) ([]Operations, error) {
	err := validator.New().Struct(input)
	if err != nil {
		return nil, errors.Wrap(err, "invalid input")
	}

	op, err := json.Marshal(input.Operations)
	if err != nil {
		return nil, errors.Wrap(err, "failed to preapply operation")
	}

	resp, err := c.post(fmt.Sprintf("/chains/%s/blocks/%s/helpers/preapply/operations", c.chain, input.Blockhash), op)
	if err != nil {
		return nil, errors.Wrap(err, "failed to preapply operation")
	}

	var operations []Operations
	err = json.Unmarshal(resp, &operations)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal operations")
	}

	return operations, nil
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

Link:
	https/tezos.gitlab.io/api/rpc.html#post-injection-operation

Parameters:

	input:
		Modifies the InjectionOperation RPC query by passing optional URL parameters. Operation is required.
*/
func (c *Client) InjectionOperation(input InjectionOperationInput) (string, error) {
	err := validator.New().Struct(input)
	if err != nil {
		return "", errors.Wrap(err, "invalid input")
	}

	v, err := json.Marshal(input.Operation)
	if err != nil {
		return "", errors.Wrap(err, "failed to inject operation")
	}
	resp, err := c.post("/injection/operation", v, input.contructRPCOptions()...)
	if err != nil {
		return "", errors.Wrap(err, "failed to inject operation")
	}

	var opstring string
	err = json.Unmarshal(resp, &opstring)
	if err != nil {
		return "", errors.Wrap(err, "failed to unmarshal operation")
	}

	return opstring, nil
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
ForgeOperation will forge an operation with the tezos RPC. For
security purposes ForgeOperationWithRPC will preapply an operation to
verify the node forged the operation with the requested contents.

NOTE:
	* Is is recommended that you forge locally with the go-tezos/v3/forge package instead. This eliminates the risk for a blind signature attack.
	* Forging with the RPC also unforges with the RPC and compares the expected contents for some security.

Path:
	../<block_id>/helpers/forge/operations (POST)

Link:
	https://tezos.gitlab.io/api/rpc.html#post-block-id-helpers-forge-operations

Parameters:

	blockhash:
		The hash of block (height) of which you want to make the query.

	branch:
		The branch of the operation.

	contents:
		The contents of the of the operation.
*/
func (c *Client) ForgeOperation(input ForgeOperationWithRPCInput) (string, error) {
	err := validator.New().Struct(input)
	if err != nil {
		return "", errors.Wrap(err, "invalid input")
	}

	op := Operations{
		Branch:   input.Branch,
		Contents: input.Contents,
	}

	v, err := json.Marshal(op)
	if err != nil {
		return "", errors.Wrap(err, "failed to forge operation")
	}

	resp, err := c.post(fmt.Sprintf("/chains/%s/blocks/%s/helpers/forge/operations", c.chain, input.Blockhash), v)
	if err != nil {
		return "", errors.Wrap(err, "failed to forge operation")
	}

	var operation string
	err = json.Unmarshal(resp, &operation)
	if err != nil {
		return "", errors.Wrap(err, "failed to forge operation")
	}

	_, opstr, err := stripBranchFromForgedOperation(operation, false)
	if err != nil {
		return operation, errors.Wrap(err, "failed to forge operation: unable to verify rpc returned a valid contents")
	}

	var rpc *Client
	if input.CheckRPCAddr != "" {
		rpc, err = New(input.CheckRPCAddr)
		if err != nil {
			return operation, errors.Wrap(err, "failed to forge operation: unable to verify rpc returned a valid contents with alternative node")
		}
	} else {
		rpc = c
	}

	operations, err := rpc.UnforgeOperation(UnforgeOperationWithRPCInput{
		Blockhash: input.Blockhash,
		Operations: []UnforgeOperationWithRPCOperation{
			{
				Data:   fmt.Sprintf("%s00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", opstr),
				Branch: input.Branch,
			},
		},
		CheckSignature: false,
	})
	if err != nil {
		return operation, errors.Wrap(err, "failed to forge operation: unable to verify rpc returned a valid contents")
	}

	for _, op := range operations {
		ok := reflect.DeepEqual(op.Contents, input.Contents)
		if !ok {
			return operation, errors.New("failed to forge operation: alert rpc returned invalid contents")
		}
	}

	return operation, nil
}

/*
UnforgeOperationWithRPC will unforge an operation with the tezos RPC.

If you would rather not use a node at all, GoTezos supports local unforging
operations REVEAL, TRANSFER, ORIGINATION, and DELEGATION.

Path:
	../<block_id>/helpers/parse/operations (POST)

Link:
	https://tezos.gitlab.io/api/rpc.html#post-block-id-helpers-parse-operations

Parameters:

	blockhash:
		The hash of block (height) of which you want to make the query.

	input:
		Contains the operations and the option to verify the operations signatures.
*/
func (c *Client) UnforgeOperation(input UnforgeOperationWithRPCInput) ([]Operations, error) {
	err := validator.New().Struct(input)
	if err != nil {
		return []Operations{}, errors.Wrap(err, "invalid input")
	}

	v, err := json.Marshal(input)
	if err != nil {
		return []Operations{}, errors.Wrap(err, "failed to unforge forge operations with RPC")
	}

	resp, err := c.post(fmt.Sprintf("/chains/%s/blocks/%s/helpers/parse/operations", c.chain, input.Blockhash), v)
	if err != nil {
		return []Operations{}, errors.Wrap(err, "failed to unforge forge operations with RPC")
	}

	var operations []Operations
	err = json.Unmarshal(resp, &operations)
	if err != nil {
		return []Operations{}, errors.Wrap(err, "failed to unforge forge operations with RPC")
	}

	return operations, nil
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

Link:
	https/tezos.gitlab.io/api/rpc.html#post-injection-operation

Parameters:

	input:
		Modifies the InjectionBlock RPC query by passing optional URL parameters. Block is required.
*/
func (c *Client) InjectionBlock(input InjectionBlockInput) ([]byte, error) {
	err := validator.New().Struct(input)
	if err != nil {
		return []byte{}, errors.Wrap(err, "invalid input")
	}

	v, err := json.Marshal(*input.Block)
	if err != nil {
		return []byte{}, errors.Wrap(err, "failed to inject block")
	}
	resp, err := c.post("/injection/block", v, input.contructRPCOptions()...)
	if err != nil {
		return resp, errors.Wrap(err, "failed to inject block")
	}
	return resp, nil
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
Counter access the counter of a contract, if any.

Path:
	../<block_id>/context/contracts/<contract_id>/counter (GET)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-context-contracts-contract-id-counter

Parameters:

	blockhash:
		The hash of block (height) of which you want to make the query.

	pkh:
		The pkh (address) of the contract for the query.
*/
func (c *Client) Counter(blockhash, pkh string) (int, error) {
	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/context/contracts/%s/counter", c.chain, blockhash, pkh))
	if err != nil {
		return 0, errors.Wrapf(err, "failed to get counter")
	}
	var strCounter string
	err = json.Unmarshal(resp, &strCounter)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to unmarshal counter")
	}

	counter, err := strconv.Atoi(strCounter)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to get counter")
	}
	return counter, nil
}

/*
RunOperation will run an operation without signature checks.

Path:
	../<block_id>/helpers/scripts/run_operation (POST)

Link:
	https://tezos.gitlab.io/api/rpc.html#post-block-id-helpers-scripts-run-operation

Parameters:

	input:
		Contains the block hash at which to make the query, and all required parameters.
*/
func (c *Client) RunOperation(input RunOperationInput) (Operations, error) {
	err := validator.New().Struct(input)
	if err != nil {
		return Operations{}, errors.Wrap(err, "invalid input")
	}

	v, err := json.Marshal(&input.Operation)
	if err != nil {
		return input.Operation.Operation, errors.Wrap(err, "failed to marshal operation")
	}

	resp, err := c.post(fmt.Sprintf("/chains/%s/blocks/%s/helpers/scripts/run_operation", c.chain, input.Blockhash), v)
	if err != nil {
		return input.Operation.Operation, errors.Wrapf(err, "failed to run_operation")
	}

	var op Operations
	err = json.Unmarshal(resp, &op)
	if err != nil {
		return input.Operation.Operation, errors.Wrap(err, "failed to unmarshal operation")
	}

	return op, nil
}

/*
StripBranchFromForgedOperation will strip the branch off an operation and resturn it with the
rest of the operation string minus the signature if signed.
Parameters:
	operation:
		The operation string.
	signed:
		Whether or not the operation is signed.
*/
func stripBranchFromForgedOperation(operation string, signed bool) (string, string, error) {
	if signed && len(operation) <= 128 {
		return "", operation, errors.New("failed to unforge branch from operation")
	}

	if signed {
		operation = operation[:len(operation)-128]
	}

	var result, rest string
	if len(operation) < 64 {
		result = operation
	} else {
		result = operation[:64]
		rest = operation[64:]
	}

	resultByts, err := hex.DecodeString(result)
	if err != nil {
		return "", operation, errors.New("failed to unforge branch from operation")
	}

	branch := crypto.B58cencode(resultByts, []byte{1, 52})
	if err != nil {
		return branch, rest, errors.Wrap(err, "failed to unforge branch from operation")
	}

	return branch, rest, nil
}
