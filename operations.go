package gotezos

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type kind string

var (
	TRANSACTION kind = "transaction"
	REVEAL      kind = "reveal"
	ORIGINATION kind = "origination"
	DELEGATION  kind = "delegation"
)

// PreapplyOperations pre-applies an operation
func (t *GoTezos) PreapplyOperations(headhash string, contents []Contents, signature string) ([]byte, error) {
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

	resp, err := t.post("/chains/main/blocks/head/helpers/preapply/operations", op)
	if err != nil {
		return resp, errors.Wrap(err, "failed to preapply operation")
	}

	return resp, nil
}

// InjectOperation injects an signed operation string and returns the response
func (t *GoTezos) InjectOperation(operation string) ([]byte, error) {
	v, err := json.Marshal(operation)
	if err != nil {
		return []byte{}, errors.Wrapf(err, "failed to inject operation '%s'", operation)
	}
	resp, err := t.post("/injection/operation", v)
	if err != nil {
		return resp, errors.Wrapf(err, "failed to inject operation '%s'", operation)
	}
	return resp, nil
}

//Counter returns the counter of the current account
func (t *GoTezos) Counter(blockhash, address string) (int, error) {
	resp, err := t.get(fmt.Sprintf("/chains/main/blocks/%s/context/contracts/%s/counter", blockhash, address))
	if err != nil {
		return 0, errors.Wrapf(err, "failed to get counter")
	}
	var strCounter string
	err = json.Unmarshal(resp, &strCounter)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to get counter")
	}

	counter, err := strconv.Atoi(strCounter)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to get counter")
	}
	return counter, nil
}

// ForgeOperation will return a forged operation based of the content(s) passed.
// Operations Supported: transaction, reveal, origination, delegation
func (t *GoTezos) ForgeOperation(branch string, contents ...Contents) (string, error) {
	cleanBranch, err := removeHexPrefix(branch, prefix_branch)
	if err != nil {
		return "", errors.Wrap(err, "failed to forge operation")
	}

	if len(cleanBranch) != 64 {
		return "", fmt.Errorf("failed to forge operation: operation branch invalid length %d", len(cleanBranch))
	}

	var sb strings.Builder
	sb.WriteString(cleanBranch)

	for _, c := range contents {
		switch c.Kind {
		case string(TRANSACTION):
			forge, err := t.forgeTransactionOperation(c)
			if err != nil {
				return "", errors.Wrap(err, "failed to forge operation")
			}
			sb.WriteString(forge)
		case string(REVEAL):
			forge, err := t.forgeRevealOperation(c)
			if err != nil {
				return "", errors.Wrap(err, "failed to forge operation")
			}
			sb.WriteString(forge)
		case string(ORIGINATION):
			forge, err := t.forgeOriginationOperation(c)
			if err != nil {
				return "", errors.Wrap(err, "failed to forge operation")
			}
			sb.WriteString(forge)
		case string(DELEGATION):
			forge, err := t.forgeDelegationOperation(c)
			if err != nil {
				return "", errors.Wrap(err, "failed to forge operation")
			}
			sb.WriteString(forge)
		default:
			return "", fmt.Errorf("failed to forge operation: unsupported kind %s", c.Kind)
		}
	}

	return sb.String(), nil
}

func (t *GoTezos) forgeTransactionOperation(contents Contents) (string, error) {
	commonFields, err := t.forgeCommonFields(contents)
	if err != nil {
		return "", errors.Wrap(err, "could not forge transaction")
	}
	var sb strings.Builder
	sb.WriteString("6c")
	sb.WriteString(commonFields)
	sb.WriteString(bigNumberToZarith(contents.Amount))

	var cleanDestination string
	if strings.HasPrefix(strings.ToLower(contents.Destination), "kt") {
		dest, err := removeHexPrefix(contents.Destination, prefix_kt)
		if err != nil {
			return "", errors.Wrap(err, "could not forge transaction: provided destination is invalid")
		}
		cleanDestination = fmt.Sprintf("%s%s%s", "01", dest, "00")
	} else {
		cleanDestination, err = removeHexPrefix(contents.Destination, prefix_tz1)
		if err != nil {
			return "", errors.Wrap(err, "could not forge transaction: provided destination is invalid")
		}
	}

	if len(cleanDestination) > 44 {
		return "", errors.New("could not forge transaction: provided destination is invalid")
	}

	for len(cleanDestination) != 44 {
		cleanDestination = fmt.Sprintf("0%s", cleanDestination)
	}

	// TODO account for code
	sb.WriteString(cleanDestination)
	sb.WriteString("00")

	return sb.String(), nil
}

func (t *GoTezos) forgeRevealOperation(contents Contents) (string, error) {
	var sb strings.Builder
	sb.WriteString("6b")
	common, err := t.forgeCommonFields(contents)
	if err != nil {
		return "", errors.Wrap(err, "failed to forge reveal operation")
	}
	sb.WriteString(common)

	cleanPubKey, err := removeHexPrefix(contents.Phk, prefix_edpk)
	if err != nil {
		return "", errors.Wrap(err, "failed to forge reveal operation")
	}

	if len(cleanPubKey) == 32 {
		errors.Wrap(err, "failed to forge reveal operation: public key is invalid")
	}

	sb.WriteString(fmt.Sprintf("00%s", cleanPubKey))

	return sb.String(), nil
}

func (t *GoTezos) forgeOriginationOperation(contents Contents) (string, error) {
	var sb strings.Builder
	sb.WriteString("6d")

	common, err := t.forgeCommonFields(contents)
	if err != nil {
		return "", errors.Wrap(err, "failed to forge origination operation")
	}

	sb.WriteString(common)
	sb.WriteString(bigNumberToZarith(contents.Balance))

	source, err := removeHexPrefix(contents.Source, prefix_tz1)
	if err != nil {
		return "", errors.Wrap(err, "failed to forge origination operation")
	}
	sb.WriteString(source)

	if len(source) > 42 {
		return "", errors.Wrap(err, "failed to forge origination operation: source is invalid")
	}

	for len(source) != 42 {
		source = fmt.Sprintf("0%s", source)
	}

	if contents.Delegate != "" {
		dest, err := removeHexPrefix(contents.Destination, prefix_tz1)
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

func (t *GoTezos) forgeDelegationOperation(contents Contents) (string, error) {
	var sb strings.Builder
	sb.WriteString("6e")

	common, err := t.forgeCommonFields(contents)
	if err != nil {
		return "", errors.Wrap(err, "failed to forge delegation operation")
	}
	sb.WriteString(common)

	var dest string
	if contents.Delegate != "" {
		sb.WriteString("ff")

		if strings.HasPrefix(strings.ToLower(contents.Delegate), "tz1") {
			dest, err = removeHexPrefix(contents.Delegate, prefix_tz1)
			if err != nil {
				return "", errors.Wrap(err, "failed to forge delegation operation")
			}
		} else if strings.HasPrefix(strings.ToLower(contents.Delegate), "kt1") {
			dest, err = removeHexPrefix(contents.Delegate, prefix_kt)
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

func (t *GoTezos) forgeCommonFields(contents Contents) (string, error) {
	source, err := removeHexPrefix(contents.Source, prefix_tz1)
	if err != nil {
		return "", errors.New("failed to remove tz1 from source prefix")
	}

	lensrc := len(source)
	if lensrc > 42 {
		return "", fmt.Errorf("invalid source length %d", lensrc)
	}

	for lensrc != 42 {
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

// UnforgeOperation will unforge an operation by returning its branch and contents.
// Operations Supported: transaction, reveal, origination, delegation
func (t *GoTezos) UnforgeOperation(hexString string, signed bool) (string, []Contents, error) {
	if signed && len(hexString) <= 128 {
		return "", []Contents{}, errors.New("failed to unforge operation: not a valid signed transaction")
	}

	if signed {
		hexString = hexString[:len(hexString)-128]
	}

	result, rest := split(hexString, 64)
	branch := prefixAndBase58Encode(result, prefix_branch)

	var contents []Contents
	for len(rest) > 0 {
		result, rest = split(rest, 2)
		switch result {
		case "6b":
			c, r, err := t.unforgeRevealOperation(hexString)
			if err != nil {
				return branch, contents, errors.New("failed to unforge operation")
			}
			rest = r
			contents = append(contents, c)
		case "6c":
			c, r, err := t.unforgeTransactionOperation(hexString)
			if err != nil {
				return branch, contents, errors.New("failed to unforge operation")
			}
			rest = r
			contents = append(contents, c)
		case "6d":
			c, r, err := t.unforgeOriginationOperation(hexString)
			if err != nil {
				return branch, contents, errors.New("failed to unforge operation")
			}
			rest = r
			contents = append(contents, c)
		case "6e":
			c, r, err := t.unforgeDelegationOperation(hexString)
			if err != nil {
				return branch, contents, errors.New("failed to unforge operation")
			}
			rest = r
			contents = append(contents, c)
		default:
			return branch, contents, fmt.Errorf("failed to unforge operation: transaction operation unkown %s", result)
		}
	}

	return branch, contents, nil
}

func (t *GoTezos) unforgeRevealOperation(hexString string) (Contents, string, error) {
	result, rest := split(hexString, 42)
	source, err := parseTzAddress(result)
	if err != nil {
		return Contents{}, rest, errors.Wrap(err, "failed to unforge reveal operation")
	}

	var contents Contents
	contents.Source = source

	zEndIndex, err := findZarithEndIndex(rest)
	if err != nil {
		return Contents{}, rest, errors.Wrap(err, "failed to unforge reveal operation")
	}
	result, rest = split(rest, zEndIndex)
	zBigNum, err := zarithToBigNumber(result)
	if err != nil {
		return Contents{}, rest, errors.Wrap(err, "failed to unforge reveal operation")
	}
	contents.Fee = zBigNum

	zEndIndex, err = findZarithEndIndex(rest)
	if err != nil {
		return Contents{}, rest, errors.Wrap(err, "failed to unforge reveal operation")
	}
	result, rest = split(rest, zEndIndex)
	zBigNum, err = zarithToBigNumber(result)
	if err != nil {
		return Contents{}, rest, errors.Wrap(err, "failed to unforge reveal operation")
	}
	contents.Counter = zBigNum

	zEndIndex, err = findZarithEndIndex(rest)
	if err != nil {
		return Contents{}, rest, errors.Wrap(err, "failed to unforge reveal operation")
	}
	result, rest = split(rest, zEndIndex)
	zBigNum, err = zarithToBigNumber(result)
	if err != nil {
		return Contents{}, rest, errors.Wrap(err, "failed to unforge reveal operation")
	}
	contents.GasLimit = zBigNum

	zEndIndex, err = findZarithEndIndex(rest)
	if err != nil {
		return Contents{}, rest, errors.Wrap(err, "failed to unforge reveal operation")
	}
	result, rest = split(rest, zEndIndex)
	zBigNum, err = zarithToBigNumber(result)
	if err != nil {
		return Contents{}, rest, errors.Wrap(err, "failed to unforge reveal operation")
	}
	contents.StorageLimit = zBigNum

	zEndIndex, err = findZarithEndIndex(rest)
	if err != nil {
		return Contents{}, rest, errors.Wrap(err, "failed to unforge reveal operation")
	}
	result, rest = split(rest, zEndIndex)
	phk, err := parsePublicKey(result)
	if err != nil {
		return Contents{}, rest, errors.Wrap(err, "failed to unforge reveal operation")
	}
	contents.Phk = phk

	return contents, rest, nil
}

func (t *GoTezos) unforgeTransactionOperation(hexString string) (Contents, string, error) {
	result, rest := split(hexString, 42)
	source, err := parseTzAddress(result)
	if err != nil {
		return Contents{}, rest, errors.Wrap(err, "failed to unforge transaction operation")
	}

	var contents Contents
	contents.Source = source

	zEndIndex, err := findZarithEndIndex(rest)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge transaction operation")
	}
	result, rest = split(rest, zEndIndex)
	zBigNum, err := zarithToBigNumber(result)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge transaction operation")
	}
	contents.Fee = zBigNum

	zEndIndex, err = findZarithEndIndex(rest)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge transaction operation")
	}
	result, rest = split(rest, zEndIndex)
	zBigNum, err = zarithToBigNumber(result)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge transaction operation")
	}
	contents.Counter = zBigNum

	zEndIndex, err = findZarithEndIndex(rest)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge transaction operation")
	}
	result, rest = split(rest, zEndIndex)
	zBigNum, err = zarithToBigNumber(result)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge transaction operation")
	}
	contents.GasLimit = zBigNum

	zEndIndex, err = findZarithEndIndex(rest)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge transaction operation")
	}
	result, rest = split(rest, zEndIndex)
	zBigNum, err = zarithToBigNumber(result)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge transaction operation")
	}
	contents.StorageLimit = zBigNum

	zEndIndex, err = findZarithEndIndex(rest)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge transaction operation")
	}
	result, rest = split(rest, zEndIndex)
	zBigNum, err = zarithToBigNumber(result)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge transaction operation")
	}
	contents.Amount = zBigNum

	zEndIndex, err = findZarithEndIndex(rest)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge transaction operation")
	}
	result, rest = split(rest, zEndIndex)
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

	return contents, rest, nil
}

func (t *GoTezos) unforgeOriginationOperation(hexString string) (Contents, string, error) {
	result, rest := split(hexString, 42)
	source, err := parseTzAddress(result)
	if err != nil {
		return Contents{}, rest, errors.Wrap(err, "failed to unforge origination operation")
	}

	contents := Contents{
		Source: source,
	}

	zEndIndex, err := findZarithEndIndex(rest)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge origination operation")
	}
	result, rest = split(rest, zEndIndex)
	zBigNum, err := zarithToBigNumber(result)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge origination operation")
	}
	contents.Fee = zBigNum

	zEndIndex, err = findZarithEndIndex(rest)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge origination operation")
	}
	result, rest = split(rest, zEndIndex)
	zBigNum, err = zarithToBigNumber(result)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge origination operation")
	}
	contents.Counter = zBigNum

	zEndIndex, err = findZarithEndIndex(rest)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge origination operation")
	}
	result, rest = split(rest, zEndIndex)
	zBigNum, err = zarithToBigNumber(result)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge origination operation")
	}
	contents.GasLimit = zBigNum

	zEndIndex, err = findZarithEndIndex(rest)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge origination operation")
	}
	result, rest = split(rest, zEndIndex)
	zBigNum, err = zarithToBigNumber(result)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge origination operation")
	}
	contents.StorageLimit = zBigNum

	zEndIndex, err = findZarithEndIndex(rest)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge origination operation")
	}
	result, rest = split(rest, zEndIndex)
	zBigNum, err = zarithToBigNumber(result)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge origination operation")
	}
	contents.Balance = zBigNum

	hasDelegate, err := checkBoolean(result)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge origination operation")
	}

	var delegate string
	if hasDelegate {
		result, rest = split(rest, 42)
		delegate, err = parseAddress(fmt.Sprintf("00%s", result))
	}
	contents.Delegate = delegate

	// TODO: decode script

	return contents, rest, nil
}

func (t *GoTezos) unforgeDelegationOperation(hexString string) (Contents, string, error) {
	result, rest := split(hexString, 42)
	source, err := parseTzAddress(result)
	if err != nil {
		return Contents{}, rest, errors.Wrap(err, "failed to unforge origination operation")
	}

	contents := Contents{
		Source: source,
	}

	zEndIndex, err := findZarithEndIndex(rest)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge origination operation")
	}
	result, rest = split(rest, zEndIndex)
	zBigNum, err := zarithToBigNumber(result)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge origination operation")
	}
	contents.Fee = zBigNum

	zEndIndex, err = findZarithEndIndex(rest)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge origination operation")
	}
	result, rest = split(rest, zEndIndex)
	zBigNum, err = zarithToBigNumber(result)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge origination operation")
	}
	contents.Counter = zBigNum

	zEndIndex, err = findZarithEndIndex(rest)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge origination operation")
	}
	result, rest = split(rest, zEndIndex)
	zBigNum, err = zarithToBigNumber(result)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge origination operation")
	}
	contents.GasLimit = zBigNum

	zEndIndex, err = findZarithEndIndex(rest)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge origination operation")
	}
	result, rest = split(rest, zEndIndex)
	zBigNum, err = zarithToBigNumber(result)
	if err != nil {
		return Contents{}, "", errors.Wrap(err, "failed to unforge origination operation")
	}
	contents.StorageLimit = zBigNum

	var delegate string
	if len(rest) == 42 {
		result, rest = split(fmt.Sprintf("01%s", rest[2:]), 42)
		delegate, err = parseAddress(result)
		if err != nil {
			return Contents{}, "", errors.Wrap(err, "failed to unforge origination operation")
		}
	} else if len(rest) > 42 {
		result, rest = split(fmt.Sprintf("00%s", rest[2:]), 44)
		delegate, err = parseAddress(result)
		if err != nil {
			return Contents{}, "", errors.Wrap(err, "failed to unforge origination operation")
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
	result, rest := split(rawHexAddress, 2)
	if result == "00" {
		return parseTzAddress(rest)
	} else if result == "01" {
		return prefixAndBase58Encode(rest[:len(rest)-2], prefix_kt), nil
	}

	return "", errors.New("address format not supported")
}

func parseTzAddress(rawHexAddress string) (string, error) {
	result, rest := split(rawHexAddress, 2)
	if result == "00" {
		return prefixAndBase58Encode(rest, prefix_tz1), nil
	}

	return "", errors.New("address format not supported")
}

func parsePublicKey(rawHexPublicKey string) (string, error) {
	result, rest := split(rawHexPublicKey, 2)
	if result == "00" {
		return prefixAndBase58Encode(rest, prefix_edpk), nil
	}

	return "", errors.New("public key format not supported")
}

func findZarithEndIndex(hexString string) (int, error) {
	for i := 0; i < len(hexString); i += 2 {
		byteSection := hexString[i:2]
		byteInt, err := strconv.ParseInt(byteSection, 16, 64)
		if err != nil {
			return 0, errors.New("failed to find Zarith end index")
		}

		if len(strconv.FormatInt(byteInt, 2)) != 8 {
			return i + 2, nil
		}
	}

	return 0, errors.New("provided hex string is not Zarith encoded")
}

func zarithToBigNumber(hexString string) (BigInt, error) {
	var bitString string
	for i := 0; i < len(hexString); i += 2 {
		byteSection := hexString[i:2]
		intSection, err := strconv.ParseInt(byteSection, 16, 64)
		if err != nil {
			return BigInt{}, errors.New("failed to find Zarith end index")
		}
		bitSection := fmt.Sprintf("00000000%d", strconv.FormatInt(intSection, 2)[:])
		bitString = fmt.Sprintf("%s%s", bitSection, bitString)
	}

	n := new(big.Int)
	n, ok := n.SetString(bitString, 2)
	if !ok {
		return BigInt{}, errors.New("failed to find Zarith end index")
	}

	b := BigInt{*n}
	return b, nil
}

func prefixAndBase58Encode(hexStringPayload string, prefix prefix) string {
	prefixHex := hex.EncodeToString(prefix)

	return b58encode(bytes.NewBufferString(fmt.Sprintf("%s%s", prefixHex, hexStringPayload)).Bytes())
}

func split(payload string, length int) (string, string) {
	res := payload[:length]
	rest := payload[length : len(payload)-length]

	return res, rest
}

func bigNumberToZarith(num BigInt) string {
	bitString := fmt.Sprintf("%b", num.Int64())
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

		if len(hexStringSection)%2 == 0 {
			hexStringSection = fmt.Sprintf("0%s", hexStringSection)
		}

		resultHexString = fmt.Sprintf("%s%s", resultHexString, hexStringSection)
	}

	return resultHexString
}

func removeHexPrefix(payload string, prefix prefix) (string, error) {
	strPrefix := hex.EncodeToString([]byte(prefix))
	bytePayload, err := b58decode(payload)
	if err != nil {
		return "", err
	}

	strPayload := hex.EncodeToString(bytePayload)
	if strings.HasPrefix(strPayload, strPrefix) {
		if len(strPrefix) < len(strPrefix)*2 {
			return "", errors.New("invalid payload")
		}
		return strPayload[:len(strPrefix)*2], nil
	}

	return "", fmt.Errorf("payload did not match prefix: %s", strPrefix)
}
