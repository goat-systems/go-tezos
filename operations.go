package gotezos

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"strings"

	validator "github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

const (
	// TRANSACTION is a kind of operation
	TRANSACTION = "transaction"
	// REVEAL is a kind of operation
	REVEAL = "reveal"
	// ORIGINATION is a kind of operation
	ORIGINATION = "origination"
	// DELEGATION is a kind of operation
	DELEGATION = "delegation"
	// ENDORSEMENT is a kind of operation
	ENDORSEMENT = "endorsement"
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
	Blockhash string   `validate:"required"`
	Protocol  string   `validate:"required"`
	Signature string   `validate:"required"`
	Contents  Contents `validate:"required"`
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
func (t *GoTezos) PreapplyOperations(input PreapplyOperationsInput) ([]Operations, error) {
	err := validator.New().Struct(input)
	if err != nil {
		return nil, errors.Wrap(err, "invalid input")
	}

	ops := []Operations{
		{
			Protocol:  input.Protocol,
			Branch:    input.Blockhash,
			Contents:  input.Contents,
			Signature: input.Signature,
		},
	}

	op, err := json.Marshal(ops)
	if err != nil {
		return nil, errors.Wrap(err, "failed to preapply operation")
	}

	resp, err := t.post(fmt.Sprintf("/chains/main/blocks/%s/helpers/preapply/operations", input.Blockhash), op)
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
func (t *GoTezos) InjectionOperation(input InjectionOperationInput) (string, error) {
	err := validator.New().Struct(input)
	if err != nil {
		return "", errors.Wrap(err, "invalid input")
	}

	v, err := json.Marshal(input.Operation)
	if err != nil {
		return "", errors.Wrap(err, "failed to inject operation")
	}
	resp, err := t.post("/injection/operation", v, input.contructRPCOptions()...)
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
ForgeOperationWithRPC will forge an operation with the tezos RPC. For
security purposes ForgeOperationWithRPC will preapply an operation to
verify the node forged the operation with the requested contents.

If you would rather not use a node at all, GoTezos supports local forging
operations REVEAL, TRANSFER, ORIGINATION, and DELEGATION.

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
func (t *GoTezos) ForgeOperationWithRPC(input ForgeOperationWithRPCInput) (string, error) {
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

	resp, err := t.post(fmt.Sprintf("/chains/main/blocks/%s/helpers/forge/operations", input.Blockhash), v)
	if err != nil {
		return "", errors.Wrap(err, "failed to forge operation")
	}

	var operation string
	err = json.Unmarshal(resp, &operation)
	if err != nil {
		return "", errors.Wrap(err, "failed to forge operation")
	}

	_, opstr, err := StripBranchFromForgedOperation(operation, false)
	if err != nil {
		return operation, errors.Wrap(err, "failed to forge operation: unable to verify rpc returned a valid contents")
	}

	var gt *GoTezos
	if input.CheckRPCAddr != "" {
		gt, err = New(input.CheckRPCAddr)
		if err != nil {
			return operation, errors.Wrap(err, "failed to forge operation: unable to verify rpc returned a valid contents with alternative node")
		}
	} else {
		gt = t
	}

	operations, err := gt.UnforgeOperationWithRPC(UnforgeOperationWithRPCInput{
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
func (t *GoTezos) UnforgeOperationWithRPC(input UnforgeOperationWithRPCInput) ([]Operations, error) {
	err := validator.New().Struct(input)
	if err != nil {
		return []Operations{}, errors.Wrap(err, "invalid input")
	}

	v, err := json.Marshal(input)
	if err != nil {
		return []Operations{}, errors.Wrap(err, "failed to unforge forge operations with RPC")
	}

	resp, err := t.post(fmt.Sprintf("/chains/main/blocks/%s/helpers/parse/operations", input.Blockhash), v)
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
func (t *GoTezos) InjectionBlock(input InjectionBlockInput) ([]byte, error) {
	err := validator.New().Struct(input)
	if err != nil {
		return []byte{}, errors.Wrap(err, "invalid input")
	}

	v, err := json.Marshal(*input.Block)
	if err != nil {
		return []byte{}, errors.Wrap(err, "failed to inject block")
	}
	resp, err := t.post("/injection/block", v, input.contructRPCOptions()...)
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
func (t *GoTezos) Counter(blockhash, pkh string) (int, error) {
	resp, err := t.get(fmt.Sprintf("/chains/main/blocks/%s/context/contracts/%s/counter", blockhash, pkh))
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
ForgeOperation forges an operation locally. GoTezos does not use the RPC or a trusted source to forge operations.
Current supported operations include transfer, reveal, delegation, and origination.

Parameters:

	branch:
		The branch to forge the operation on.

	contents:
		The operation contents to be formed.
*/
func ForgeOperation(branch string, contents Contents) (string, error) {
	if len(contents.DoubleEndorsementEvidence) > 0 || len(contents.AccountActivations) > 0 || len(contents.DoubleBakingEvidence) > 0 ||
		len(contents.Ballots) > 0 || len(contents.Proposals) > 0 || len(contents.SeedNonceRevelations) > 0 {
		return "", errors.New("unsupported kind: supported kinds for local forging are transaction, origination, reveal, and delegation")
	}

	cleanBranch, err := cleanBranch(branch)
	if err != nil {
		return "", errors.Wrap(err, "failed to forge operation")
	}

	var sb strings.Builder
	sb.WriteString(cleanBranch)

	for _, transaction := range contents.Transactions {
		forge, err := transaction.Forge("")
		if err != nil {
			return "", errors.Wrap(err, "failed to forge operation")
		}
		sb.WriteString(forge)
	}

	for _, reveal := range contents.Reveals {
		forge, err := reveal.Forge("")
		if err != nil {
			return "", errors.Wrap(err, "failed to forge operation")
		}
		sb.WriteString(forge)
	}

	for _, origination := range contents.Originations {
		forge, err := origination.Forge("")
		if err != nil {
			return "", errors.Wrap(err, "failed to forge operation")
		}
		sb.WriteString(forge)
	}

	for _, delegation := range contents.Delegations {
		forge, err := delegation.Forge("")
		if err != nil {
			return "", errors.Wrap(err, "failed to forge operation")
		}
		sb.WriteString(forge)
	}

	return sb.String(), nil
}

func cleanBranch(branch string) (string, error) {
	cleanBranch, err := removeHexPrefix(branch, branchprefix)
	if err != nil {
		return "", errors.Wrap(err, "failed to clean branch")
	}

	if len(cleanBranch) != 64 {
		return "", fmt.Errorf("failed to clean branch: operation branch invalid length %d", len(cleanBranch))
	}

	return cleanBranch, nil
}

/*
Forge forges a transaction operation locally. GoTezos does not use the RPC or a trusted source to forge operations.
Current supported operations include transfer, reveal, delegation, and origination. Forging a transaction locally is limited
and does not support code (TODO).

Parameters:
	branch:
		The branch to forge the operation on.
*/
func (t *Transaction) Forge(branch string) (string, error) {
	err := validator.New().Struct(t)
	if err != nil {
		return "", errors.Wrap(err, "invalid input")
	}

	var sb strings.Builder
	if branch != "" {
		cleanBranch, err := cleanBranch(branch)
		if err != nil {
			return "", errors.Wrap(err, "failed to forge transaction(s)")
		}
		sb.WriteString(cleanBranch)
	}

	commonFields, err := forgeCommonFields(t.toContentsHelper())
	if err != nil {
		return "", errors.Wrap(err, "failed to forge transaction")
	}
	sb.WriteString("6c")
	sb.WriteString(commonFields)
	sb.WriteString(bigNumberToZarith(*t.Amount))

	var cleanDestination string
	if strings.HasPrefix(strings.ToLower(t.Destination), "kt") {
		dest, err := removeHexPrefix(t.Destination, ktprefix)
		if err != nil {
			return "", errors.Wrap(err, "failed to forge transaction: provided destination is not a valid KT1 address")
		}
		cleanDestination = fmt.Sprintf("%s%s%s", "01", dest, "00")
	} else {
		cleanDestination, err = removeHexPrefix(t.Destination, tz1prefix)
		if err != nil {
			return "", errors.Wrap(err, "failed to forge transaction: provided destination is not a valid tz1 address")
		}
	}

	if len(cleanDestination) > 44 {
		return "", errors.New("failed to forge transaction: provided destination is of invalid length")
	}

	for len(cleanDestination) != 44 {
		cleanDestination = fmt.Sprintf("0%s", cleanDestination)
	}

	// TODO account for code
	sb.WriteString(cleanDestination)
	sb.WriteString("00")

	return sb.String(), nil
}

/*
Forge forges a reveal operation locally. GoTezos does not use the RPC or a trusted source to forge operations.
Current supported operations include transfer, reveal, delegation, and origination.

Parameters:
	branch:
		The branch to forge the operation on.
*/
func (r *Reveal) Forge(branch string) (string, error) {
	err := validator.New().Struct(r)
	if err != nil {
		return "", errors.Wrap(err, "invalid input")
	}

	var sb strings.Builder
	if branch != "" {
		cleanBranch, err := cleanBranch(branch)
		if err != nil {
			return "", errors.Wrap(err, "failed to forge transaction(s)")
		}
		sb.WriteString(cleanBranch)
	}

	sb.WriteString("6b")
	common, err := forgeCommonFields(r.toContentsHelper())
	if err != nil {
		return "", errors.Wrap(err, "failed to forge reveal operation")
	}
	sb.WriteString(common)

	cleanPubKey, err := removeHexPrefix(r.PublicKey, edpkprefix)
	if err != nil {
		return "", errors.Wrap(err, "failed to forge reveal operation")
	}

	if len(cleanPubKey) == 32 {
		errors.Wrap(err, "failed to forge reveal operation: public key is invalid")
	}

	sb.WriteString(fmt.Sprintf("00%s", cleanPubKey))

	return sb.String(), nil
}

/*
Forge forges a origination operation(s) locally. GoTezos does not use the RPC or a trusted source to forge operations.
Current supported operations include transfer, reveal, delegation, and origination.

Parameters:
	branch:
		The branch to forge the operation on.
*/
func (o *Origination) Forge(branch string) (string, error) {
	err := validator.New().Struct(o)
	if err != nil {
		return "", errors.Wrap(err, "invalid input")
	}

	var sb strings.Builder
	if branch != "" {
		cleanBranch, err := cleanBranch(branch)
		if err != nil {
			return "", errors.Wrap(err, "failed to forge transaction(s)")
		}
		sb.WriteString(cleanBranch)
	}
	sb.WriteString("6d")

	common, err := forgeCommonFields(o.toContentsHelper())
	if err != nil {
		return "", errors.Wrap(err, "failed to forge origination operation")
	}

	sb.WriteString(common)
	sb.WriteString(bigNumberToZarith(*o.Balance))

	source, err := removeHexPrefix(o.Source, tz1prefix)
	if err != nil {
		return "", errors.Wrap(err, "failed to forge origination operation")
	}

	if len(source) > 42 {
		return "", errors.Wrap(err, "failed to forge origination operation: source is invalid")
	}

	for len(source) != 42 {
		source = fmt.Sprintf("0%s", source)
	}

	if o.Delegate != "" {

		dest, err := removeHexPrefix(o.Delegate, tz1prefix)
		if err != nil {
			return "", errors.Wrap(err, "failed to forge origination operation")
		}

		if len(source) > 42 {
			return "", errors.Wrap(err, "failed to forge origination operation: source is invalid")
		}

		for len(dest) != 42 {
			dest = fmt.Sprintf("0%s", dest)
		}

		sb.WriteString("ff")
		sb.WriteString(dest)
	} else {
		sb.WriteString("00")
	}

	sb.WriteString("000000c602000000c105000764085e036c055f036d0000000325646f046c000000082564656661756c740501035d050202000000950200000012020000000d03210316051f02000000020317072e020000006a0743036a00000313020000001e020000000403190325072c020000000002000000090200000004034f0327020000000b051f02000000020321034c031e03540348020000001e020000000403190325072c020000000002000000090200000004034f0327034f0326034202000000080320053d036d0342")
	sb.WriteString("0000001a")
	sb.WriteString("0a")
	sb.WriteString("00000015")
	sb.WriteString(source)

	return sb.String(), nil
}

/*
Forge forges a delegation operation locally. GoTezos does not use the RPC or a trusted source to forge operations.
Current supported operations include transfer, reveal, delegation, and origination.

Parameters:
	branch:
		The branch to forge the operation on.
*/
func (d *Delegation) Forge(branch string) (string, error) {
	err := validator.New().Struct(d)
	if err != nil {
		return "", errors.Wrap(err, "invalid input")
	}

	var sb strings.Builder
	if branch != "" {
		cleanBranch, err := cleanBranch(branch)
		if err != nil {
			return "", errors.Wrap(err, "failed to forge transaction(s)")
		}
		sb.WriteString(cleanBranch)
	}
	sb.WriteString("6e")

	common, err := forgeCommonFields(d.toContentsHelper())
	if err != nil {
		return "", errors.Wrap(err, "failed to forge delegation operation")
	}
	sb.WriteString(common)

	var dest string
	if d.Delegate != "" {
		sb.WriteString("ff")

		if strings.HasPrefix(strings.ToLower(d.Delegate), "tz1") {
			dest, err = removeHexPrefix(d.Delegate, tz1prefix)
			if err != nil {
				return "", errors.Wrap(err, "failed to forge delegation operation")
			}
		} else if strings.HasPrefix(strings.ToLower(d.Delegate), "kt1") {
			dest, err = removeHexPrefix(d.Delegate, ktprefix)
			if err != nil {
				return "", errors.Wrap(err, "failed to forge delegation operation")
			}
		}

		if len(dest) > 42 {
			return "", errors.Wrap(err, "failed to forge delegation operation: dest is invalid")
		}

		for len(dest) != 42 {
			dest = fmt.Sprintf("0%s", dest)
		}

		sb.WriteString(dest)
	} else {
		sb.WriteString("00")
	}

	return sb.String(), nil
}

func forgeCommonFields(contents ContentsHelper) (string, error) {
	source, err := removeHexPrefix(contents.Source, tz1prefix)
	if err != nil {
		return "", errors.New("failed to remove tz1 from source prefix")
	}

	if len(source) > 42 {
		return "", fmt.Errorf("invalid source length %d", len(source))
	}

	for len(source) != 42 {
		source = fmt.Sprintf("0%s", source)
	}

	var sb strings.Builder
	sb.WriteString(source)
	sb.WriteString(bigNumberToZarith(*contents.Fee))
	sb.WriteString(bigNumberToZarith(*NewInt(contents.Counter)))
	sb.WriteString(bigNumberToZarith(*contents.GasLimit))
	sb.WriteString(bigNumberToZarith(*contents.StorageLimit))

	return sb.String(), nil
}

/*
UnforgeOperation takes a forged/encoded tezos operation and decodes it by returning the
operations branch, and contents.

Parameters:

	operation:
		The hex string encoded operation.

	signed:
		The ?true Unforge will decode a signed operation.
*/
func UnforgeOperation(operation string, signed bool) (string, Contents, error) {
	if signed && len(operation) <= 128 {
		return "", Contents{}, errors.New("failed to unforge operation: not a valid signed transaction")
	}

	if signed {
		operation = operation[:len(operation)-128]
	}

	result, rest := splitAndReturnRest(operation, 64)
	branch, err := prefixAndBase58Encode(result, branchprefix)
	if err != nil {
		return branch, Contents{}, errors.Wrap(err, "failed to unforge operation")
	}

	var contents Contents
	for len(rest) > 0 {
		result, rest = splitAndReturnRest(rest, 2)
		if result == "00" || len(result) < 2 {
			break
		}

		switch result {
		case "6b":
			reveal, rst, err := unforgeRevealOperation(rest)
			if err != nil {
				return branch, contents, errors.Wrap(err, "failed to unforge operation")
			}
			rest = rst
			contents.Reveals = append(contents.Reveals, reveal)
		case "6c":
			transaction, rst, err := unforgeTransactionOperation(rest)
			if err != nil {
				return branch, contents, errors.Wrap(err, "failed to unforge operation")
			}
			rest = rst
			contents.Transactions = append(contents.Transactions, transaction)
		case "6d":
			origination, rst, err := unforgeOriginationOperation(rest)
			if err != nil {
				return branch, contents, errors.Wrap(err, "failed to unforge operation")
			}
			rest = rst
			contents.Originations = append(contents.Originations, origination)
		case "6e":
			delegation, rst, err := unforgeDelegationOperation(rest)
			if err != nil {
				return branch, contents, errors.Wrap(err, "failed to unforge operation")
			}
			rest = rst
			contents.Delegations = append(contents.Delegations, delegation)
		default:
			return branch, contents, fmt.Errorf("failed to unforge operation: transaction operation unkown %s", result)
		}
	}

	return branch, contents, nil
}

func unforgeRevealOperation(hexString string) (Reveal, string, error) {
	result, rest := splitAndReturnRest(hexString, 42)
	source, err := parseTzAddress(result)
	if err != nil {
		return Reveal{}, rest, errors.Wrap(err, "failed to unforge reveal operation")
	}

	var reveal Reveal
	reveal.Kind = REVEAL
	reveal.Source = source

	zEndIndex, err := findZarithEndIndex(rest)
	if err != nil {
		return reveal, rest, errors.Wrap(err, "failed to unforge reveal operation")
	}
	result, rest = splitAndReturnRest(rest, zEndIndex)
	zBigNum, err := zarithToBigNumber(result)
	if err != nil {
		return reveal, rest, errors.Wrap(err, "failed to unforge reveal operation")
	}
	reveal.Fee = zBigNum

	zEndIndex, err = findZarithEndIndex(rest)
	if err != nil {
		return reveal, rest, errors.Wrap(err, "failed to unforge reveal operation")
	}
	result, rest = splitAndReturnRest(rest, zEndIndex)
	zBigNum, err = zarithToBigNumber(result)
	if err != nil {
		return reveal, rest, errors.Wrap(err, "failed to unforge reveal operation")
	}
	reveal.Counter = int(zBigNum.Big.Int64())

	zEndIndex, err = findZarithEndIndex(rest)
	if err != nil {
		return reveal, rest, errors.Wrap(err, "failed to unforge reveal operation")
	}
	result, rest = splitAndReturnRest(rest, zEndIndex)
	zBigNum, err = zarithToBigNumber(result)
	if err != nil {
		return reveal, rest, errors.Wrap(err, "failed to unforge reveal operation")
	}
	reveal.GasLimit = zBigNum

	zEndIndex, err = findZarithEndIndex(rest)
	if err != nil {
		return reveal, rest, errors.Wrap(err, "failed to unforge reveal operation")
	}
	result, rest = splitAndReturnRest(rest, zEndIndex)
	zBigNum, err = zarithToBigNumber(result)
	if err != nil {
		return reveal, rest, errors.Wrap(err, "failed to unforge reveal operation")
	}
	reveal.StorageLimit = zBigNum

	result, rest = splitAndReturnRest(rest, 66)
	phk, err := parsePublicKey(result)
	if err != nil {
		return reveal, rest, errors.Wrap(err, "failed to unforge reveal operation")
	}
	reveal.PublicKey = phk

	return reveal, rest, nil
}

func unforgeTransactionOperation(hexString string) (Transaction, string, error) {
	result, rest := splitAndReturnRest(hexString, 42)
	source, err := parseTzAddress(result)
	if err != nil {
		return Transaction{}, rest, errors.Wrap(err, "failed to unforge transaction operation")
	}

	var transaction Transaction
	transaction.Source = source
	transaction.Kind = TRANSACTION

	zarithEndIndex, err := findZarithEndIndex(rest)
	if err != nil {
		return transaction, "", errors.Wrap(err, "failed to unforge transaction operation")
	}

	result, rest = splitAndReturnRest(rest, zarithEndIndex)
	zBigNum, err := zarithToBigNumber(result)
	if err != nil {
		return transaction, "", errors.Wrap(err, "failed to unforge transaction operation")
	}
	transaction.Fee = zBigNum

	zarithEndIndex, err = findZarithEndIndex(rest)
	if err != nil {
		return transaction, "", errors.Wrap(err, "failed to unforge transaction operation")
	}
	result, rest = splitAndReturnRest(rest, zarithEndIndex)
	zBigNum, err = zarithToBigNumber(result)
	if err != nil {
		return transaction, "", errors.Wrap(err, "failed to unforge transaction operation")
	}
	transaction.Counter = int(zBigNum.Big.Int64())

	zarithEndIndex, err = findZarithEndIndex(rest)
	if err != nil {
		return transaction, "", errors.Wrap(err, "failed to unforge transaction operation")
	}
	result, rest = splitAndReturnRest(rest, zarithEndIndex)
	zBigNum, err = zarithToBigNumber(result)
	if err != nil {
		return transaction, "", errors.Wrap(err, "failed to unforge transaction operation")
	}
	transaction.GasLimit = zBigNum

	zarithEndIndex, err = findZarithEndIndex(rest)
	if err != nil {
		return transaction, "", errors.Wrap(err, "failed to unforge transaction operation")
	}
	result, rest = splitAndReturnRest(rest, zarithEndIndex)
	zBigNum, err = zarithToBigNumber(result)
	if err != nil {
		return transaction, "", errors.Wrap(err, "failed to unforge transaction operation")
	}
	transaction.StorageLimit = zBigNum

	zarithEndIndex, err = findZarithEndIndex(rest)
	if err != nil {
		return transaction, "", errors.Wrap(err, "failed to unforge transaction operation")
	}
	result, rest = splitAndReturnRest(rest, zarithEndIndex)
	zBigNum, err = zarithToBigNumber(result)
	if err != nil {
		return transaction, "", errors.Wrap(err, "failed to unforge transaction operation")
	}
	transaction.Amount = zBigNum

	result, rest = splitAndReturnRest(rest, 44)
	address, err := parseAddress(result)
	if err != nil {
		return transaction, "", errors.Wrap(err, "failed to unforge transaction operation")
	}
	transaction.Destination = address

	// TODO Handle Contracts
	// hasParameters, err := checkBoolean(result)
	// if err != nil {
	// 	return Contents{}, "", errors.Wrap(err, "failed to unforge transaction operation: could not check for parameters")
	// }

	// Temporary: Trim 00
	if len(rest) > 2 {
		rest = rest[2:]
	}

	transaction.Kind = TRANSACTION

	return transaction, rest, nil
}

func unforgeOriginationOperation(hexString string) (Origination, string, error) {
	result, rest := splitAndReturnRest(hexString, 42)
	source, err := parseTzAddress(result)
	if err != nil {
		return Origination{}, rest, errors.Wrap(err, "failed to unforge origination operation")
	}

	origination := Origination{
		Source: source,
		Kind:   ORIGINATION,
	}

	zEndIndex, err := findZarithEndIndex(rest)
	if err != nil {
		return origination, "", errors.Wrap(err, "failed to unforge origination operation")
	}
	result, rest = splitAndReturnRest(rest, zEndIndex)
	zBigNum, err := zarithToBigNumber(result)
	if err != nil {
		return origination, "", errors.Wrap(err, "failed to unforge origination operation")
	}
	origination.Fee = zBigNum

	zEndIndex, err = findZarithEndIndex(rest)
	if err != nil {
		return origination, "", errors.Wrap(err, "failed to unforge origination operation")
	}
	result, rest = splitAndReturnRest(rest, zEndIndex)
	zBigNum, err = zarithToBigNumber(result)
	if err != nil {
		return origination, "", errors.Wrap(err, "failed to unforge origination operation")
	}
	origination.Counter = int(zBigNum.Big.Int64())

	zEndIndex, err = findZarithEndIndex(rest)
	if err != nil {
		return origination, "", errors.Wrap(err, "failed to unforge origination operation")
	}
	result, rest = splitAndReturnRest(rest, zEndIndex)
	zBigNum, err = zarithToBigNumber(result)
	if err != nil {
		return origination, "", errors.Wrap(err, "failed to unforge origination operation")
	}
	origination.GasLimit = zBigNum

	zEndIndex, err = findZarithEndIndex(rest)
	if err != nil {
		return origination, "", errors.Wrap(err, "failed to unforge origination operation")
	}
	result, rest = splitAndReturnRest(rest, zEndIndex)
	zBigNum, err = zarithToBigNumber(result)
	if err != nil {
		return origination, "", errors.Wrap(err, "failed to unforge origination operation")
	}
	origination.StorageLimit = zBigNum

	zEndIndex, err = findZarithEndIndex(rest)
	if err != nil {
		return origination, "", errors.Wrap(err, "failed to unforge origination operation")
	}
	result, rest = splitAndReturnRest(rest, zEndIndex)
	zBigNum, err = zarithToBigNumber(result)
	if err != nil {
		return origination, "", errors.Wrap(err, "failed to unforge origination operation")
	}
	origination.Balance = zBigNum

	result, rest = splitAndReturnRest(rest, 2)
	hasDelegate, err := checkBoolean(result)
	if err != nil {
		return origination, "", errors.Wrap(err, "failed to unforge origination operation")
	}

	var delegate string
	if hasDelegate {
		result, rest = splitAndReturnRest(rest, 42)
		delegate, err = parseAddress(fmt.Sprintf("00%s", result))
		if err != nil {
			return origination, "", errors.Wrap(err, "failed to unforge origination operation")
		}
	}
	origination.Delegate = delegate

	// TODO: decode script

	return origination, rest, nil
}

func unforgeDelegationOperation(hexString string) (Delegation, string, error) {
	result, rest := splitAndReturnRest(hexString, 42)
	source, err := parseAddress(result)
	if err != nil {
		return Delegation{}, rest, errors.Wrap(err, "failed to unforge delegation operation")
	}

	delegation := Delegation{
		Source: source,
		Kind:   DELEGATION,
	}

	zEndIndex, err := findZarithEndIndex(rest)
	if err != nil {
		return delegation, "", errors.Wrap(err, "failed to unforge delegation operation")
	}
	result, rest = splitAndReturnRest(rest, zEndIndex)
	zBigNum, err := zarithToBigNumber(result)
	if err != nil {
		return delegation, "", errors.Wrap(err, "failed to unforge delegation operation")
	}
	delegation.Fee = zBigNum

	zEndIndex, err = findZarithEndIndex(rest)
	if err != nil {
		return delegation, "", errors.Wrap(err, "failed to unforge delegation operation")
	}
	result, rest = splitAndReturnRest(rest, zEndIndex)
	zBigNum, err = zarithToBigNumber(result)
	if err != nil {
		return delegation, "", errors.Wrap(err, "failed to unforge delegation operation")
	}
	delegation.Counter = int(zBigNum.Big.Int64())

	zEndIndex, err = findZarithEndIndex(rest)
	if err != nil {
		return delegation, "", errors.Wrap(err, "failed to unforge delegation operation")
	}
	result, rest = splitAndReturnRest(rest, zEndIndex)
	zBigNum, err = zarithToBigNumber(result)
	if err != nil {
		return delegation, "", errors.Wrap(err, "failed to unforge delegation operation")
	}
	delegation.GasLimit = zBigNum

	zEndIndex, err = findZarithEndIndex(rest)
	if err != nil {
		return delegation, "", errors.Wrap(err, "failed to unforge delegation operation")
	}
	result, rest = splitAndReturnRest(rest, zEndIndex)
	zBigNum, err = zarithToBigNumber(result)
	if err != nil {
		return delegation, "", errors.Wrap(err, "failed to unforge delegation operation")
	}
	delegation.StorageLimit = zBigNum

	var delegate string
	if len(rest) == 42 {
		result, rest = splitAndReturnRest(fmt.Sprintf("01%s", rest[2:]), 42)
		delegate, err = parseAddress(result)
		if err != nil {
			return delegation, "", errors.Wrap(err, "failed to unforge delegation operation")
		}
	} else if len(rest) > 42 {
		result, rest = splitAndReturnRest(fmt.Sprintf("00%s", rest[2:]), 44)
		delegate, err = parseAddress(result)
		if err != nil {
			return delegation, "", errors.Wrap(err, "failed to unforge delegation operation")
		}
	} else if len(rest) == 2 && rest == "00" {
		rest = ""
	}
	delegation.Delegate = delegate

	return delegation, rest, nil
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
func StripBranchFromForgedOperation(operation string, signed bool) (string, string, error) {
	if signed && len(operation) <= 128 {
		return "", operation, errors.New("failed to unforge branch from operation")
	}

	if signed {
		operation = operation[:len(operation)-128]
	}

	result, rest := splitAndReturnRest(operation, 64)
	branch, err := prefixAndBase58Encode(result, branchprefix)
	if err != nil {
		return branch, rest, errors.Wrap(err, "failed to unforge branch from operation")
	}

	return branch, rest, nil
}

// type parameters struct {
// 	amount      BigInt
// 	destination string
// 	rest        string
// }

// func UnforgeParameters(hexString string) (parameters, error) {
// 	result, rest := split(hexString, 2)
// 	i := &big.Int{}
// 	i.SetString(result, 16)
// 	result, rest = split(rest, 40)
// 	result, rest = split(rest, 42)
// 	destination, err := parseTzAddress(result)
// 	if err != nil {
// 		return parameters{}, errors.Wrap(err, "failed to parse destination address from parameters")
// 	}
// 	result, rest = split(rest, 12)
// 	i = i.Mul(i, big.NewInt(2))
// 	i = i.Sub(i, big.NewInt(106))
// 	result, rest = split(rest, int(i.Int64()))

// 	amount := new(big.Int)
// 	amount.SetString(result[2:], 16)
// 	result, rest = split(rest, 12)

// 	return parameters{
// 		amount:      BigInt{*amount},
// 		destination: destination,
// 		rest:        rest,
// 	}, nil
// }

func checkBoolean(hexString string) (bool, error) {
	if hexString == "ff" {
		return true, nil
	} else if hexString == "00" {
		return false, nil
	}
	return false, errors.New("boolean value is invalid")
}

func parseAddress(rawHexAddress string) (string, error) {
	result, rest := splitAndReturnRest(rawHexAddress, 2)
	if strings.HasPrefix(rawHexAddress, "0000") {
		rawHexAddress = rawHexAddress[2:]
	}
	if result == "00" {
		return parseTzAddress(rawHexAddress)
	} else if result == "01" {
		encode, err := prefixAndBase58Encode(rest[:len(rest)-2], ktprefix)
		if err != nil {
			errors.Wrap(err, "address format not supported")
		}
		return encode, nil
	}

	return "", errors.New("address format not supported")
}

func parseTzAddress(rawHexAddress string) (string, error) {
	result, rest := splitAndReturnRest(rawHexAddress, 2)
	if result == "00" {
		encode, err := prefixAndBase58Encode(rest, tz1prefix)
		if err != nil {
			errors.Wrap(err, "address format not supported")
		}
		return encode, nil
	}

	return "", errors.New("address format not supported")
}

func parsePublicKey(rawHexPublicKey string) (string, error) {
	result, rest := splitAndReturnRest(rawHexPublicKey, 2)
	if result == "00" {
		encode, err := prefixAndBase58Encode(rest, edpkprefix)
		if err != nil {
			errors.Wrap(err, "failed to base58 encode public key")
		}
		return encode, nil
	}

	return "", errors.New("public key format not supported")
}

func findZarithEndIndex(hexString string) (int, error) {
	for i := 0; i < len(hexString); i += 2 {
		byteSection := hexString[i : i+2]
		byteInt, err := strconv.ParseUint(byteSection, 16, 64)
		if err != nil {
			return 0, errors.New("failed to find Zarith end index")
		}

		if len(strconv.FormatInt(int64(byteInt), 2)) != 8 {
			return i + 2, nil
		}
	}

	return 0, errors.New("provided hex string is not Zarith encoded")
}

func zarithToBigNumber(hexString string) (*Int, error) {
	var bitString string
	for i := 0; i < len(hexString); i += 2 {
		byteSection := hexString[i : i+2]
		intSection, err := strconv.ParseInt(byteSection, 16, 64)
		if err != nil {
			return NewInt(0), errors.New("failed to find Zarith end index")
		}

		bitSection := fmt.Sprintf("00000000%s", strconv.FormatInt(intSection, 2))
		bitSection = bitSection[len(bitSection)-7:]
		bitString = fmt.Sprintf("%s%s", bitSection, bitString)
	}

	n := new(big.Int)
	n, ok := n.SetString(bitString, 2)
	if !ok {
		return NewInt(0), errors.New("failed to find Zarith end index")
	}

	b := Int{n}
	return &b, nil
}

func prefixAndBase58Encode(hexPayload string, prefix prefix) (string, error) {
	v, err := hex.DecodeString(fmt.Sprintf("%s%s", hex.EncodeToString(prefix), hexPayload))
	if err != nil {
		return "", errors.Wrap(err, "failed to encode to base58")
	}
	return encode(v), nil
}

func splitAndReturnRest(payload string, length int) (string, string) {
	if len(payload) < length {
		return payload, ""
	}

	return payload[:length], payload[length:]
}

func bigNumberToZarith(num Int) string {
	bitString := fmt.Sprintf("%b", num.Big.Int64())
	for len(bitString)%7 != 0 {
		bitString = fmt.Sprintf("0%s", bitString)
	}

	var resultHexString string
	for i := len(bitString); i > 0; i -= 7 {
		bitStringSection := bitString[i-7 : i]

		if i == 7 {
			bitStringSection = fmt.Sprintf("0%s", bitStringSection)
		} else {
			bitStringSection = fmt.Sprintf("1%s", bitStringSection)
		}

		x, _ := strconv.ParseInt(bitStringSection, 2, 64)
		hexStringSection := strconv.FormatInt(x, 16)

		if len(hexStringSection)%2 != 0 {
			hexStringSection = fmt.Sprintf("0%s", hexStringSection)
		}

		resultHexString = fmt.Sprintf("%s%s", resultHexString, hexStringSection)
	}

	return resultHexString
}

func removeHexPrefix(base58CheckEncodedPayload string, prefix prefix) (string, error) {
	strPrefix := hex.EncodeToString([]byte(prefix))
	base58CheckEncodedPayloadBytes, err := decode(base58CheckEncodedPayload)
	if err != nil {
		return "", fmt.Errorf("failed to decode payload: %s", base58CheckEncodedPayload)
	}
	base58CheckEncodedPayload = hex.EncodeToString(base58CheckEncodedPayloadBytes)

	if strings.HasPrefix(base58CheckEncodedPayload, strPrefix) {
		return base58CheckEncodedPayload[len(prefix)*2:], nil
	}

	return "", fmt.Errorf("payload did not match prefix: %s", strPrefix)
}
