package gotezos

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

const (
	// TRANSACTIONOP is a kind of operation
	TRANSACTIONOP = "transaction"
	// REVEALOP is a kind of operation
	REVEALOP = "reveal"
	// ORIGINATIONOP is a kind of operation
	ORIGINATIONOP = "origination"
	// DELEGATIONOP is a kind of operation
	DELEGATIONOP = "delegation"
)

/*
InjectionOperationInput is the input for the goTezos.InjectionOperation function.

Function:
	func (t *GoTezos) InjectionOperation(input *InjectionOperationInput) (*[]byte, error) {}
*/
type InjectionOperationInput struct {
	// The operation string.
	Operation *string `validate:"required"`

	// If ?async is true, the function returns immediately.
	Async bool

	// Specify the ChainID.
	ChainID *string
}

/*
InjectionBlockInput is the input for the goTezos.InjectionBlock function.

Function:
	func (t *GoTezos) InjectionBlock(input *InjectionBlockInput) (**[]byte, error) {}
*/
type InjectionBlockInput struct {
	// Block to inject
	Block *Block `validate:"required"`

	// If ?async is true, the function returns immediately.
	Async bool

	// If ?force is true, it will be injected even on non strictly increasing fitness.
	Force bool

	// Specify the ChainID.
	ChainID *string
}

/*
PreapplyOperations simulates the validation of an operation.

Path:
	../<block_id>/helpers/preapply/operations (POST)

Link:
	https://tezos.gitlab.io/api/rpc.html#post-block-id-helpers-preapply-operations

Parameters:

	blockhash:
		The hash of block (height) of which you want to make the query.

	contents:
		The contents of the of the operation.

	signature:
		The operation signature.
*/
func (t *GoTezos) PreapplyOperations(blockhash string, contents []Contents, signature string) (*[]byte, error) {
	head, err := t.Head()
	if err != nil {
		return nil, errors.Wrap(err, "failed to preapply operation")
	}

	operations := []Operations{
		Operations{
			Protocol:  head.Protocol,
			Branch:    head.Hash,
			Contents:  contents,
			Signature: signature,
		},
	}

	op, err := json.Marshal(operations)
	if err != nil {
		return nil, errors.Wrap(err, "failed to preapply operation")
	}

	resp, err := t.post(fmt.Sprintf("/chains/main/blocks/%s/helpers/preapply/operations", blockhash), op)
	if err != nil {
		return &resp, errors.Wrap(err, "failed to preapply operation")
	}

	return &resp, nil
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
func (t *GoTezos) InjectionOperation(input *InjectionOperationInput) (*[]byte, error) {
	err := validator.New().Struct(input)
	if err != nil {
		return &[]byte{}, errors.Wrap(err, "invalid input")
	}

	v, err := json.Marshal(*input.Operation)
	if err != nil {
		return &[]byte{}, errors.Wrap(err, "failed to inject operation")
	}
	resp, err := t.post("/injection/operation", v, input.contructRPCOptions()...)
	if err != nil {
		return &resp, errors.Wrap(err, "failed to inject operation")
	}
	return &resp, nil
}

func (i *InjectionOperationInput) contructRPCOptions() []rpcOptions {
	var opts []rpcOptions
	if i.Async {
		opts = append(opts, rpcOptions{
			"async",
			"true",
		})
	}

	if i.ChainID != nil {
		opts = append(opts, rpcOptions{
			"chain_id",
			*i.ChainID,
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

Link:
	https/tezos.gitlab.io/api/rpc.html#post-injection-operation

Parameters:

	input:
		Modifies the InjectionBlock RPC query by passing optional URL parameters. Block is required.
*/
func (t *GoTezos) InjectionBlock(input *InjectionBlockInput) (*[]byte, error) {
	err := validator.New().Struct(input)
	if err != nil {
		return &[]byte{}, errors.Wrap(err, "invalid input")
	}

	v, err := json.Marshal(*input.Block)
	if err != nil {
		return &[]byte{}, errors.Wrap(err, "failed to inject block")
	}
	resp, err := t.post("/injection/block", v, input.contructRPCOptions()...)
	if err != nil {
		return &resp, errors.Wrap(err, "failed to inject block")
	}
	return &resp, nil
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

	if i.ChainID != nil {
		opts = append(opts, rpcOptions{
			"chain_id",
			*i.ChainID,
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
func (t *GoTezos) Counter(blockhash, pkh string) (*int, error) {
	resp, err := t.get(fmt.Sprintf("/chains/main/blocks/%s/context/contracts/%s/counter", blockhash, pkh))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get counter")
	}
	var strCounter string
	err = json.Unmarshal(resp, &strCounter)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal counter")
	}

	counter, err := strconv.Atoi(strCounter)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get counter")
	}
	return &counter, nil
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
func ForgeOperation(branch string, contents ...Contents) (*string, error) {
	cleanBranch, err := removeHexPrefix(branch, branchprefix)
	if err != nil {
		return nil, errors.Wrap(err, "failed to forge operation")
	}

	if len(cleanBranch) != 64 {
		return nil, fmt.Errorf("failed to forge operation: operation branch invalid length %d", len(cleanBranch))
	}

	var sb strings.Builder
	sb.WriteString(cleanBranch)

	for _, c := range contents {
		switch c.Kind {
		case TRANSACTIONOP:
			forge, err := forgeTransactionOperation(c)
			if err != nil {
				return nil, errors.Wrap(err, "failed to forge operation")
			}
			sb.WriteString(forge)
		case REVEALOP:
			forge, err := forgeRevealOperation(c)
			if err != nil {
				return nil, errors.Wrap(err, "failed to forge operation")
			}
			sb.WriteString(forge)
		case ORIGINATIONOP:
			forge, err := forgeOriginationOperation(c)
			if err != nil {
				return nil, errors.Wrap(err, "failed to forge operation")
			}
			sb.WriteString(forge)
		case DELEGATIONOP:
			forge, err := forgeDelegationOperation(c)
			if err != nil {
				return nil, errors.Wrap(err, "failed to forge operation")
			}
			sb.WriteString(forge)
		default:
			return nil, fmt.Errorf("failed to forge operation: unsupported kind %s", c.Kind)
		}
	}
	operation := sb.String()

	return &operation, nil
}

func forgeTransactionOperation(contents Contents) (string, error) {
	commonFields, err := forgeCommonFields(contents)
	if err != nil {
		return "", errors.Wrap(err, "could not forge transaction")
	}
	var sb strings.Builder
	sb.WriteString("6c")
	sb.WriteString(commonFields)
	sb.WriteString(bigNumberToZarith(contents.Amount))

	var cleanDestination string
	if strings.HasPrefix(strings.ToLower(contents.Destination), "kt") {
		dest, err := removeHexPrefix(contents.Destination, ktprefix)
		if err != nil {
			return "", errors.Wrap(err, "could not forge transaction: provided destination is not a valid KT1 address")
		}
		cleanDestination = fmt.Sprintf("%s%s%s", "01", dest, "00")
	} else {
		cleanDestination, err = removeHexPrefix(contents.Destination, tz1prefix)
		if err != nil {
			return "", errors.Wrap(err, "could not forge transaction: provided destination is not a valid tz1 address")
		}
	}

	if len(cleanDestination) > 44 {
		return "", errors.New("could not forge transaction: provided destination is of invalid length")
	}

	for len(cleanDestination) != 44 {
		cleanDestination = fmt.Sprintf("0%s", cleanDestination)
	}

	// TODO account for code
	sb.WriteString(cleanDestination)
	sb.WriteString("00")

	return sb.String(), nil
}

func forgeRevealOperation(contents Contents) (string, error) {
	var sb strings.Builder
	sb.WriteString("6b")
	common, err := forgeCommonFields(contents)
	if err != nil {
		return "", errors.Wrap(err, "failed to forge reveal operation")
	}
	sb.WriteString(common)

	cleanPubKey, err := removeHexPrefix(contents.Phk, edpkprefix)
	if err != nil {
		return "", errors.Wrap(err, "failed to forge reveal operation")
	}

	if len(cleanPubKey) == 32 {
		errors.Wrap(err, "failed to forge reveal operation: public key is invalid")
	}

	sb.WriteString(fmt.Sprintf("00%s", cleanPubKey))

	return sb.String(), nil
}

func forgeOriginationOperation(contents Contents) (string, error) {
	var sb strings.Builder
	sb.WriteString("6d")

	common, err := forgeCommonFields(contents)
	if err != nil {
		return "", errors.Wrap(err, "failed to forge origination operation")
	}

	sb.WriteString(common)
	sb.WriteString(bigNumberToZarith(contents.Balance))

	source, err := removeHexPrefix(contents.Source, tz1prefix)
	if err != nil {
		return "", errors.Wrap(err, "failed to forge origination operation")
	}

	if len(source) > 42 {
		return "", errors.Wrap(err, "failed to forge origination operation: source is invalid")
	}

	for len(source) != 42 {
		source = fmt.Sprintf("0%s", source)
	}

	if contents.Delegate != "" {

		dest, err := removeHexPrefix(contents.Delegate, tz1prefix)
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

func forgeDelegationOperation(contents Contents) (string, error) {
	var sb strings.Builder
	sb.WriteString("6e")

	common, err := forgeCommonFields(contents)
	if err != nil {
		return "", errors.Wrap(err, "failed to forge delegation operation")
	}
	sb.WriteString(common)

	var dest string
	if contents.Delegate != "" {
		sb.WriteString("ff")

		if strings.HasPrefix(strings.ToLower(contents.Delegate), "tz1") {
			dest, err = removeHexPrefix(contents.Delegate, tz1prefix)
			if err != nil {
				return "", errors.Wrap(err, "failed to forge delegation operation")
			}
		} else if strings.HasPrefix(strings.ToLower(contents.Delegate), "kt1") {
			dest, err = removeHexPrefix(contents.Delegate, ktprefix)
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

func forgeCommonFields(contents Contents) (string, error) {
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
	sb.WriteString(bigNumberToZarith(contents.Fee))
	sb.WriteString(bigNumberToZarith(contents.Counter))
	sb.WriteString(bigNumberToZarith(contents.GasLimit))
	sb.WriteString(bigNumberToZarith(contents.StorageLimit))

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
func UnforgeOperation(operation string, signed bool) (*string, *[]Contents, error) {
	if signed && len(operation) <= 128 {
		return nil, &[]Contents{}, errors.New("failed to unforge operation: not a valid signed transaction")
	}

	if signed {
		operation = operation[:len(operation)-128]
	}

	result, rest := splitAndReturnRest(operation, 64)
	branch, err := prefixAndBase58Encode(result, branchprefix)
	if err != nil {
		return &branch, &[]Contents{}, errors.Wrap(err, "failed to unforge operation")
	}

	var contents []Contents
	for len(rest) > 0 {
		result, rest = splitAndReturnRest(rest, 2)
		if result == "00" || len(result) < 2 {
			break
		}

		switch result {
		case "6b":
			c, r, err := unforgeRevealOperation(rest)
			if err != nil {
				return &branch, &contents, errors.Wrap(err, "failed to unforge operation")
			}
			rest = r
			contents = append(contents, c)
		case "6c":
			c, r, err := unforgeTransactionOperation(rest)
			if err != nil {
				return &branch, &contents, errors.Wrap(err, "failed to unforge operation")
			}
			rest = r
			contents = append(contents, c)
		case "6d":
			c, r, err := unforgeOriginationOperation(rest)
			if err != nil {
				return &branch, &contents, errors.Wrap(err, "failed to unforge operation")
			}
			rest = r
			contents = append(contents, c)
		case "6e":
			c, r, err := unforgeDelegationOperation(rest)
			if err != nil {
				return &branch, &contents, errors.Wrap(err, "failed to unforge operation")
			}
			rest = r
			contents = append(contents, c)
		default:
			return &branch, &contents, fmt.Errorf("failed to unforge operation: transaction operation unkown %s", result)
		}
	}

	return &branch, &contents, nil
}

func unforgeRevealOperation(hexString string) (Contents, string, error) {
	result, rest := splitAndReturnRest(hexString, 42)
	source, err := parseTzAddress(result)
	if err != nil {
		return Contents{}, rest, errors.Wrap(err, "failed to unforge reveal operation")
	}

	var contents Contents
	contents.Kind = REVEALOP
	contents.Source = source

	zEndIndex, err := findZarithEndIndex(rest)
	if err != nil {
		return Contents{}, rest, errors.Wrap(err, "failed to unforge reveal operation")
	}
	result, rest = splitAndReturnRest(rest, zEndIndex)
	zBigNum, err := zarithToBigNumber(result)
	if err != nil {
		return Contents{}, rest, errors.Wrap(err, "failed to unforge reveal operation")
	}
	contents.Fee = zBigNum

	zEndIndex, err = findZarithEndIndex(rest)
	if err != nil {
		return Contents{}, rest, errors.Wrap(err, "failed to unforge reveal operation")
	}
	result, rest = splitAndReturnRest(rest, zEndIndex)
	zBigNum, err = zarithToBigNumber(result)
	if err != nil {
		return Contents{}, rest, errors.Wrap(err, "failed to unforge reveal operation")
	}
	contents.Counter = zBigNum

	zEndIndex, err = findZarithEndIndex(rest)
	if err != nil {
		return Contents{}, rest, errors.Wrap(err, "failed to unforge reveal operation")
	}
	result, rest = splitAndReturnRest(rest, zEndIndex)
	zBigNum, err = zarithToBigNumber(result)
	if err != nil {
		return Contents{}, rest, errors.Wrap(err, "failed to unforge reveal operation")
	}
	contents.GasLimit = zBigNum

	zEndIndex, err = findZarithEndIndex(rest)
	if err != nil {
		return Contents{}, rest, errors.Wrap(err, "failed to unforge reveal operation")
	}
	result, rest = splitAndReturnRest(rest, zEndIndex)
	zBigNum, err = zarithToBigNumber(result)
	if err != nil {
		return Contents{}, rest, errors.Wrap(err, "failed to unforge reveal operation")
	}
	contents.StorageLimit = zBigNum

	result, rest = splitAndReturnRest(rest, 66)
	phk, err := parsePublicKey(result)
	if err != nil {
		return Contents{}, rest, errors.Wrap(err, "failed to unforge reveal operation")
	}
	contents.Phk = phk

	return contents, rest, nil
}

func unforgeTransactionOperation(hexString string) (Contents, string, error) {
	result, rest := splitAndReturnRest(hexString, 42)
	source, err := parseTzAddress(result)
	if err != nil {
		return Contents{}, rest, errors.Wrap(err, "failed to unforge transaction operation")
	}

	var contents Contents
	contents.Source = source
	contents.Kind = TRANSACTIONOP

	zarithEndIndex, err := findZarithEndIndex(rest)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge transaction operation")
	}

	result, rest = splitAndReturnRest(rest, zarithEndIndex)
	zBigNum, err := zarithToBigNumber(result)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge transaction operation")
	}
	contents.Fee = zBigNum

	zarithEndIndex, err = findZarithEndIndex(rest)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge transaction operation")
	}
	result, rest = splitAndReturnRest(rest, zarithEndIndex)
	zBigNum, err = zarithToBigNumber(result)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge transaction operation")
	}
	contents.Counter = zBigNum

	zarithEndIndex, err = findZarithEndIndex(rest)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge transaction operation")
	}
	result, rest = splitAndReturnRest(rest, zarithEndIndex)
	zBigNum, err = zarithToBigNumber(result)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge transaction operation")
	}
	contents.GasLimit = zBigNum

	zarithEndIndex, err = findZarithEndIndex(rest)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge transaction operation")
	}
	result, rest = splitAndReturnRest(rest, zarithEndIndex)
	zBigNum, err = zarithToBigNumber(result)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge transaction operation")
	}
	contents.StorageLimit = zBigNum

	zarithEndIndex, err = findZarithEndIndex(rest)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge transaction operation")
	}
	result, rest = splitAndReturnRest(rest, zarithEndIndex)
	zBigNum, err = zarithToBigNumber(result)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge transaction operation")
	}
	contents.Amount = zBigNum

	result, rest = splitAndReturnRest(rest, 44)
	address, err := parseAddress(result)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge transaction operation")
	}
	contents.Destination = address

	// TODO Handle Contracts
	// hasParameters, err := checkBoolean(result)
	// if err != nil {
	// 	return Contents{}, "", errors.Wrap(err, "failed to unforge transaction operation: could not check for parameters")
	// }

	// Temporary: Trim 00
	if len(rest) > 2 {
		rest = rest[2:]
	}

	contents.Kind = TRANSACTIONOP

	return contents, rest, nil
}

func unforgeOriginationOperation(hexString string) (Contents, string, error) {
	result, rest := splitAndReturnRest(hexString, 42)
	source, err := parseTzAddress(result)
	if err != nil {
		return Contents{}, rest, errors.Wrap(err, "failed to unforge origination operation")
	}

	contents := Contents{
		Source: source,
		Kind:   ORIGINATIONOP,
	}

	zEndIndex, err := findZarithEndIndex(rest)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge origination operation")
	}
	result, rest = splitAndReturnRest(rest, zEndIndex)
	zBigNum, err := zarithToBigNumber(result)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge origination operation")
	}
	contents.Fee = zBigNum

	zEndIndex, err = findZarithEndIndex(rest)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge origination operation")
	}
	result, rest = splitAndReturnRest(rest, zEndIndex)
	zBigNum, err = zarithToBigNumber(result)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge origination operation")
	}
	contents.Counter = zBigNum

	zEndIndex, err = findZarithEndIndex(rest)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge origination operation")
	}
	result, rest = splitAndReturnRest(rest, zEndIndex)
	zBigNum, err = zarithToBigNumber(result)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge origination operation")
	}
	contents.GasLimit = zBigNum

	zEndIndex, err = findZarithEndIndex(rest)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge origination operation")
	}
	result, rest = splitAndReturnRest(rest, zEndIndex)
	zBigNum, err = zarithToBigNumber(result)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge origination operation")
	}
	contents.StorageLimit = zBigNum

	zEndIndex, err = findZarithEndIndex(rest)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge origination operation")
	}
	result, rest = splitAndReturnRest(rest, zEndIndex)
	zBigNum, err = zarithToBigNumber(result)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge origination operation")
	}
	contents.Balance = zBigNum

	result, rest = splitAndReturnRest(rest, 2)
	hasDelegate, err := checkBoolean(result)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge origination operation")
	}

	var delegate string
	if hasDelegate {
		result, rest = splitAndReturnRest(rest, 42)
		delegate, err = parseAddress(fmt.Sprintf("00%s", result))
		if err != nil {
			return Contents{}, "", errors.Wrap(err, "failed to unforge origination operation")
		}
	}
	contents.Delegate = delegate

	// TODO: decode script

	return contents, rest, nil
}

func unforgeDelegationOperation(hexString string) (Contents, string, error) {
	result, rest := splitAndReturnRest(hexString, 42)
	source, err := parseAddress(result)
	if err != nil {
		return Contents{}, rest, errors.Wrap(err, "failed to unforge delegation operation")
	}

	contents := Contents{
		Source: source,
		Kind:   DELEGATIONOP,
	}

	zEndIndex, err := findZarithEndIndex(rest)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge delegation operation")
	}
	result, rest = splitAndReturnRest(rest, zEndIndex)
	zBigNum, err := zarithToBigNumber(result)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge delegation operation")
	}
	contents.Fee = zBigNum

	zEndIndex, err = findZarithEndIndex(rest)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge delegation operation")
	}
	result, rest = splitAndReturnRest(rest, zEndIndex)
	zBigNum, err = zarithToBigNumber(result)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge delegation operation")
	}
	contents.Counter = zBigNum

	zEndIndex, err = findZarithEndIndex(rest)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge delegation operation")
	}
	result, rest = splitAndReturnRest(rest, zEndIndex)
	zBigNum, err = zarithToBigNumber(result)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge delegation operation")
	}
	contents.GasLimit = zBigNum

	zEndIndex, err = findZarithEndIndex(rest)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge delegation operation")
	}
	result, rest = splitAndReturnRest(rest, zEndIndex)
	zBigNum, err = zarithToBigNumber(result)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge delegation operation")
	}
	contents.StorageLimit = zBigNum

	var delegate string
	if len(rest) == 42 {
		result, rest = splitAndReturnRest(fmt.Sprintf("01%s", rest[2:]), 42)
		delegate, err = parseAddress(result)
		if err != nil {
			return Contents{}, "", errors.Wrap(err, "failed to unforge delegation operation")
		}
	} else if len(rest) > 42 {
		result, rest = splitAndReturnRest(fmt.Sprintf("00%s", rest[2:]), 44)
		delegate, err = parseAddress(result)
		if err != nil {
			return Contents{}, "", errors.Wrap(err, "failed to unforge delegation operation")
		}
	} else if len(rest) == 2 && rest == "00" {
		rest = ""
	}
	contents.Delegate = delegate

	return contents, rest, nil
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

func zarithToBigNumber(hexString string) (Int, error) {
	var bitString string
	for i := 0; i < len(hexString); i += 2 {
		byteSection := hexString[i : i+2]
		intSection, err := strconv.ParseInt(byteSection, 16, 64)
		if err != nil {
			return Int{}, errors.New("failed to find Zarith end index")
		}

		bitSection := fmt.Sprintf("00000000%s", strconv.FormatInt(intSection, 2))
		bitSection = bitSection[len(bitSection)-7:]
		bitString = fmt.Sprintf("%s%s", bitSection, bitString)
	}

	n := new(big.Int)
	n, ok := n.SetString(bitString, 2)
	if !ok {
		return Int{}, errors.New("failed to find Zarith end index")
	}

	b := Int{n}
	return b, nil
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
