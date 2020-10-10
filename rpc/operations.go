package rpc

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	validator "github.com/go-playground/validator/v10"
	"github.com/utdrmac/go-tezos/v3/crypto"
	"github.com/pkg/errors"
)

/*
InjectionOperationInput is the input for the goTezos.InjectionOperation function.

Function:
	func (c *Client) InjectionOperation(input InjectionOperationInput) ([]byte, error) {}
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
	func (c *Client) InjectionBlock(input InjectionBlockInput) ([]byte, error) {}
*/
type InjectionBlockInput struct {

	// Block header signature
	SignedBytes string `validate:"required"`

	// Operations included in the block
	// This is not the same as operations found in mempool
	// and also not like preapply result
	Operations [][]interface{} `validate:"required"`

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

// RunOperation is a sub structure of RunOperationInput
type RunOperation struct {
	Operation Operations `json:"operation" validate:"required"`
	ChainID   string     `json:"chain_id" validate:"required"`
}

/*
UnforgeOperationInput is the input for the goTezos.UnforgeOperationWithRPC function.

Function:
	func (c *Client) UnforgeOperationWithRPC(blockhash string, operation string, checkSignature bool) (Operations, error) {}
*/
type UnforgeOperationInput struct {
	Blockhash      string             `validate:"required"`
	Operations     []UnforgeOperation `json:"operations" validate:"required"`
	CheckSignature bool               `json:"check_signature"`
}

// UnforgeOperation is a sub structure of UnforgeOperationWithRPCInput
type UnforgeOperation struct {
	Data   string `json:"data" validate:"required"`
	Branch string `json:"branch" validate:"required"`
}

/*
ForgeOperationInput is the input for the client.ForgeOperation function.

Function:
	func (c *Client) ForgeOperation(input ForgeOperationInput) (string, error) {}
*/
type ForgeOperationInput struct {
	Blockhash    string   `validate:"required"`
	Branch       string   `validate:"required"`
	Contents     Contents `validate:"required"`
	CheckRPCAddr string
}

/*
CounterInput is the input for the client.Counter function.

Function:
	func (c *Client) Counter(input CounterInput) (int, error) {}
*/
type CounterInput struct {
	Blockhash string `validate:"required"`
	Address   string `validate:"required"`
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
*/
func (c *Client) ForgeOperation(input ForgeOperationInput) (string, error) {
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

	operations, err := rpc.UnforgeOperation(UnforgeOperationInput{
		Blockhash: input.Blockhash,
		Operations: []UnforgeOperation{
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
UnforgeOperation will unforge an operation with the tezos RPC.

If you would rather not use a node at all, GoTezos supports local unforging
operations REVEAL, TRANSFER, ORIGINATION, and DELEGATION.

Path:
	../<block_id>/helpers/parse/operations (POST)

Link:
	https://tezos.gitlab.io/api/rpc.html#post-block-id-helpers-parse-operations
*/
func (c *Client) UnforgeOperation(input UnforgeOperationInput) ([]Operations, error) {
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
	/injection/block (POST)

Link:
	https/tezos.gitlab.io/api/rpc.html#post-injection-operation
*/
func (c *Client) InjectionBlock(input InjectionBlockInput) ([]byte, error) {
	err := validator.New().Struct(input)
	if err != nil {
		return []byte{}, errors.Wrap(err, "invalid input")
	}

	block := struct {
		SignedBytes string  `json:"data"`
		Ops [][]interface{} `json:"operations"`
	}{
		input.SignedBytes,
		input.Operations,
	}

	v, err := json.Marshal(block)
	if err != nil {
		return []byte{}, errors.Wrap(err, "failed to marshal new block")
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
*/
func (c *Client) Counter(input CounterInput) (int, error) {
	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/context/contracts/%s/counter", c.chain, input.Blockhash, input.Address))
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

/*
ShellHeader comprises the "shell" of each block
*/
type ShellHeader struct {
	Level            int       `json:"level"`
	Proto            int       `json:"proto"`
	Predecessor      string    `json:"predecessor"`
	Timestamp        time.Time `json:"timestamp"`
	ValidationPass   int       `json:"validation_pass"`
	OperationsHash   string    `json:"operations_hash"`
	Fitness          []string  `json:"fitness"`
	Context          string    `json:"context"`
	ProtocolData     string    `json:"protocol_data"`
}

/*
PreApplyOperations and PreApplyOperationsAlt are responses from PreApplying a forged block for baking
*/
type PreApplyOperationsAlt PreApplyOperations

func (o *PreApplyOperationsAlt) UnmarshalJSON(buf []byte) error {
	return unmarshalNamedJSONArray(buf, &o.Hash, (*PreApplyOperations)(o))
}

type ShellOperations struct {
	Applied       []PreApplyOperations    `json:"applied"`
	Refused       []PreApplyOperationsAlt `json:"refused"`
	BranchRefused []PreApplyOperationsAlt `json:"branch_refused"`
	BranchDelayed []PreApplyOperationsAlt `json:"branch_delayed"`
	Unprocessed   []PreApplyOperationsAlt `json:"unprocessed"`
}

type PreApplyOperations struct {
	Hash   string `json:"hash"`
	Branch string `json:"branch"`
	Data   string `json:"data"`
}

type PreapplyResult struct {
	Shellheader ShellHeader `json:"shell_header"`
	Operations []ShellOperations `json:"operations"`
}

/*
ProtocolData is information required in any newly forged Blocks
*/
type ProtocolData struct {
	Protocol      string `json:"protocol"`
	Priority      int    `json:"priority"`
	POWNonce      string `json:"proof_of_work_nonce"`
	Signature     string `json:"signature"`
	SeedNonceHash string `json:"seed_nonce_hash,omitempty"`
}

/*
PreapplyBlockOperationInput is the input for the goTezos.PreapplyBlockOperation function.
Function:
    func (t *GoTezos) PreapplyBlockOperation(input PreapplyBlockOperationInput) ([]interface{}, error) {}
*/
type PreapplyBlockOperationInput struct {
	// Header data for the new block
	Protocoldata ProtocolData `validate:"required"`

	// Operations (txn, endorsements, etc) to be contained in the next block
	Operations [][]Operations `validate:"required"`

	// Sort operations on return
	Sort bool `validate:"required"`

	// Earliest timestamp of next block
	Timestamp time.Time `validate:"required"`
}

func (p *PreapplyBlockOperationInput) constructRPCOptions() []rpcOptions {
    var opts []rpcOptions
    if p.Sort {
        opts = append(opts, rpcOptions{
            "sort",
            "true",
        })
    }

    if !p.Timestamp.IsZero() {
        opts = append(opts, rpcOptions{
            "timestamp",
            strconv.FormatInt(p.Timestamp.Unix(), 10),
        })
    }

    return opts
}

func (c *Client) PreapplyBlockOperation(input PreapplyBlockOperationInput) (PreapplyResult, error) {
	var preapplyRes PreapplyResult
	err  := validator.New().Struct(input)
	if err != nil {
		return preapplyRes, errors.Wrap(err, "invalid input")
	}

	// construct the block
	bh := struct {
		Pd ProtocolData `json:"protocol_data"`
		Ops [][]Operations `json:"operations"`
	}{
		input.Protocoldata,
		input.Operations,
	}

	v, err := json.Marshal(bh)
	if err != nil {
		return preapplyRes, errors.Wrap(err, "failed to marshal new block")
	}

	resp, err := c.post(fmt.Sprintf("/chains/%s/blocks/head/helpers/preapply/block", c.chain), v, input.constructRPCOptions()...)
	if err != nil {
		return preapplyRes, errors.Wrap(err, "failed to preapply new block")
	}

	// Parse the response from preapply
	err = json.Unmarshal(resp, &preapplyRes)
	if err != nil {
		return preapplyRes, errors.Wrap(err, "failed to unmarshal preapply result of new block")
	}

	return preapplyRes, nil
}

type ForgedBlockHeader struct {
	BlockHeader string `json:"block"`
}

func (c *Client) ForgeBlockHeader(sh ShellHeader) (ForgedBlockHeader, error) {

	var fBlockheader ForgedBlockHeader

	v, err := json.Marshal(sh)
	if err != nil {
		return fBlockheader, errors.Wrap(err, "failed to marshal shell header for forge")
	}

	resp, err := c.post("/chains/main/blocks/head/helpers/forge_block_header", v)
	if err != nil {
		return fBlockheader, errors.Wrap(err, "failed to forge new block")
	}
	
	// Parse the response from preapply
	err = json.Unmarshal(resp, &fBlockheader)
	if err != nil {
		return fBlockheader, errors.Wrap(err, "failed to unmarshal forged block header result")
	}

	return fBlockheader, nil
}

type EndorsingPowerInput struct {
	Endorsement  Operations `json:"endorsement_operation"`
	ChainID      string     `json:"chain_id"`
}

func (c *Client) GetEndorsingPower(operation EndorsingPowerInput) (int, error) {

	// Remove Protocol value
	operation.Endorsement.Protocol = ""

	// Marshal
	v, err := json.Marshal(operation)
	if err != nil {
		return 0, errors.Wrap(err, "failed to marshal operation for endorsing power")
	}

	resp, err := c.post(fmt.Sprintf("/chains/%s/blocks/head/endorsing_power", operation.ChainID), v)
	if err != nil {
		return 0, errors.Wrap(err, "failed to fetch endorsing power")
	}

	// RPC contains newline; strip in order to parse to int
	power, err := strconv.Atoi(strings.TrimSuffix(string(resp), "\n"))
	if err != nil {
		return 0, errors.Wrap(err, "failed to decode endorsing power")
	}

	return power, nil
}

func (c *Client) MinimalValidTime(endorsingPower, priority int, chainID string) (time.Time, error) {

	var minTimestamp time.Time

	var opts []rpcOptions
	opts = append(opts, rpcOptions{
		"priority", strconv.Itoa(priority),
	})
	opts = append(opts, rpcOptions{
		"endorsing_power", strconv.Itoa(endorsingPower),
	})

	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/head/minimal_valid_time", chainID), opts...)
	if err != nil {
		return minTimestamp, errors.Wrap(err, "failed to fetch minimal_valid_time")
	}

	minTimestamp, err = time.Parse(time.RFC3339, stripQuote(string(resp)))
	if err != nil {
		return minTimestamp, errors.Wrap(err, "failed to parse minimal valid timestamp")
	}

	return minTimestamp, nil
}
