package gotezos

import (
	"bytes"
	"encoding/hex"
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

// import (
// 	"encoding/hex"
// 	"encoding/json"
// 	"fmt"
// 	"strconv"

// 	"golang.org/x/crypto/blake2b"

// 	"github.com/DefinitelyNotAGoat/go-tezos/account"
// 	"github.com/DefinitelyNotAGoat/go-tezos/block"
// 	"github.com/DefinitelyNotAGoat/go-tezos/crypto"
// 	"github.com/DefinitelyNotAGoat/go-tezos/v2/delegate"
// 	"github.com/pkg/errors"
// 	"golang.org/x/crypto/ed25519"
// )

// var (
// 	// maxBatchSize tells how many Transactions per batch are allowed.
// 	maxBatchSize = 200
// )

// // Conts is helper structure to build out the contents of a a transfer operation to post to the Tezos RPC
// type Conts struct {
// 	Contents []Contents `json:"contents"`
// 	Branch   string     `json:"branch"`
// }

// // Transfer a complete transfer request
// type Transfer struct {
// 	Conts
// 	Protocol  string `json:"protocol"`
// 	Signature string `json:"signature"`
// }

// // ForgeOperations returns a batch of operations in string representation.
// func (w *Wallet) ForgeOperations(fee, gaslimit int, operations ...Contents) ([]string, error) {
// 	var operationSignatures []string

// 	// Get current branch head
// 	block, err := t.HeadBlock()
// 	if err != nil {
// 		return operationSignatures, errors.Wrap(err, "could not create batch payment")
// 	}

// 	// Get the counter for the payment address and increment it
// 	counter, err := t.getAddressCounter(wallet.Address)
// 	if err != nil {
// 		return operationSignatures, errors.Wrap(err, "could not create batch payment")
// 	}
// 	counter++

// 	// Split our slice of []Payment into batches
// 	batches := t.splitPaymentIntoBatches(payments, batchSize)
// 	operationSignatures = make([]string, len(batches))

// 	for k := range batches {

// 		// Convert (ie: forge) each 'Payment' into an actual Tezos transfer operation
// 		operationBytes, operationContents, newCounter, err := t.forgeOperationBytes(blockHead.Hash, counter, wallet, batches[k], paymentFee, gasLimit)
// 		if err != nil {
// 			return operationSignatures, errors.Wrap(err, "could not create batch payment")
// 		}
// 		counter = newCounter

// 		// Sign gt batch of operations with the secret key; return that signature
// 		edsig, err := t.signOperationBytes(operationBytes, wallet)
// 		if err != nil {
// 			return operationSignatures, errors.Wrap(err, "could not create batch payment")
// 		}

// 		// Extract and decode the bytes of the signature
// 		decodedSignature, err := t.decodeSignature(edsig)
// 		if err != nil {
// 			return operationSignatures, errors.Wrap(err, "could not create batch payment")
// 		}

// 		decodedSignature = decodedSignature[10:(len(decodedSignature))]

// 		// The signed bytes of gt batch
// 		fullOperation := operationBytes + decodedSignature

// 		// We can validate gt batch against the node for any errors
// 		err = t.preApplyOperations(operationContents, edsig, blockHead)
// 		if err != nil {
// 			return operationSignatures, errors.Wrap(err, "could not create batch payment")
// 		}
// 		// Add the signature (raw operation bytes & signature of operations) of gt batch of transfers to the returning slice
// 		// gt will be used to POST to /injection/operation
// 		operationSignatures[k] = fullOperation

// 	}

// 	return operationSignatures, nil
// }

// //Sign previously forged Operation bytes using secret key of wallet
// func signOperationBytes(operationBytes string, wallet account.Wallet) (string, error) {

// 	opBytes, err := hex.DecodeString(operationBytes)
// 	if err != nil {
// 		return "", errors.Wrap(err, "could not sign operation bytes")
// 	}
// 	opBytes = append(crypto.Prefix_watermark, opBytes...)

// 	// Generic hash of 32 bytes
// 	genericHash, err := blake2b.New(32, []byte{})
// 	if err != nil {
// 		return "", errors.Wrap(err, "could not sign operation bytes")
// 	}

// 	// Write operation bytes to hash
// 	i, err := genericHash.Write(opBytes)

// 	if err != nil {
// 		return "", errors.Wrap(err, "could not sign operation bytes")
// 	}
// 	if i != len(opBytes) {
// 		return "", errors.Errorf("could not sign operation, generic hash length %d does not match bytes length %d", i, len(opBytes))
// 	}

// 	finalHash := genericHash.Sum([]byte{})

// 	// Sign the finalized generic hash of operations and b58 encode
// 	sig := ed25519.Sign(wallet.Kp.PrivKey, finalHash)
// 	//sig := sodium.Bytes(finalHash).SignDetached(wallet.Kp.PrivKey)
// 	edsig := crypto.B58cencode(sig, crypto.Prefix_edsig)

// 	return edsig, nil
// }

// func (t *GoTezos) forgeOperationBytes(branchHash string, counter int, wallet account.Wallet, batch []delegate.Payment, paymentFee int, gaslimit int) (string, Conts, int, error) {

// 	var contents Conts
// 	var combinedOps []block.Contents

// 	//left here to display how to reveal a new wallet (needs funds to be revealed!)
// 	/**
// 	  combinedOps = append(combinedOps, StructContents{Kind: "reveal", PublicKey: wallet.pk , Source: wallet.address, Fee: "0", GasLimit: "127", StorageLimit: "0", Counter: strCounter})
// 	  counter++
// 	**/

// 	for k := range batch {

// 		if batch[k].Amount > 0 {

// 			operation := block.Contents{
// 				Kind:         "transaction",
// 				Source:       wallet.Address,
// 				Fee:          strconv.Itoa(paymentFee),
// 				GasLimit:     strconv.Itoa(gaslimit),
// 				StorageLimit: "0",
// 				Amount:       strconv.FormatFloat(crypto.RoundPlus(batch[k].Amount, 0), 'f', -1, 64),
// 				Destination:  batch[k].Address,
// 				Counter:      strconv.Itoa(counter),
// 			}
// 			combinedOps = append(combinedOps, operation)
// 			counter++
// 		}
// 	}
// 	contents.Contents = combinedOps
// 	contents.Branch = branchHash

// 	var opBytes string

// 	forge := "/chains/main/blocks/head/helpers/forge/operations"
// 	output, err := o.tzclient.Post(forge, contents.string())
// 	if err != nil {
// 		return "", contents, counter, errors.Wrapf(err, "could not forge operation '%s' with contents '%s'", forge, contents.string())
// 	}

// 	err = json.Unmarshal(output, &opBytes)
// 	if err != nil {
// 		return "", contents, counter, errors.Wrapf(err, "could not forge operation '%s' with contents '%s'", forge, contents.string())
// 	}

// 	return opBytes, contents, counter, nil
// }

// // PreapplyOperations, or batch of operations, to a Tezos node to ensure correctness
// func (t *GoTezos) PreapplyOperations(paymentOperations Conts, signature string, blockHead block.Block) error {

// 	// Create a full transfer request
// 	var transfer Transfer
// 	transfer.Signature = signature
// 	transfer.Contents = paymentOperations.Contents
// 	transfer.Branch = blockHead.Hash
// 	transfer.Protocol = blockHead.Protocol

// 	// RPC says outer element must be JSON array
// 	var transfers = []Transfer{transfer}

// 	// Convert object to JSON string
// 	transfersOp, err := json.Marshal(transfers)
// 	if err != nil {
// 		return errors.Wrap(err, "could not preapply operations, could not marshal into json")
// 	}

// 	// POST the JSON to the RPC
// 	query := "/chains/main/blocks/head/helpers/preapply/operations"
// 	_, err = o.tzclient.Post(query, string(transfersOp))
// 	if err != nil {
// 		return errors.Wrapf(err, "could not preapply operations '%s' with contents '%s'", query, string(transfersOp))
// 	}

// 	return nil
// }

// // InjectOperation injects an signed operation string and returns the response
// func (t *GoTezos) InjectOperation(operation string) ([]byte, error) {
// 	jsonBytes, err := json.Marshal(operation)
// 	if err != nil {
// 		return []byte{}, errors.Wrapf(err, "could not inject operation with contents '%s'", string(jsonBytes))
// 	}
// 	resp, err := t.post("/injection/operation", jsonBytes)
// 	if err != nil {
// 		return resp, errors.Wrapf(err, "could not inject operation '%s' with contents '%s'", string(jsonBytes))
// 	}
// 	return resp, nil
// }

// //Counter returns the counter of the current account
// func (t *GoTezos) Counter(blockhash, address string) (int, error) {
// 	resp, err := t.get(fmt.Sprintf("/chains/main/blocks/%s/context/contracts/%s/counter", blockhash, address))
// 	if err != nil {
// 		return 0, errors.Wrapf(err, "could not get account counter")
// 	}
// 	rtnStr, err := unmarshalString(resp)
// 	if err != nil {
// 		return 0, errors.Wrapf(err, "could not get account counter")
// 	}
// 	counter, err := strconv.Atoi(rtnStr)
// 	if err != nil {
// 		return 0, errors.Wrapf(err, "could not get account counter")
// 	}
// 	return counter, nil
// }

// func splitPaymentIntoBatches(rewards []delegate.Payment, batchSize int) [][]delegate.Payment {
// 	var batches [][]delegate.Payment
// 	for i := 0; i < len(rewards); i += batchSize {
// 		end := i + batchSize
// 		if end > len(rewards) {
// 			end = len(rewards)
// 		}
// 		batches = append(batches, rewards[i:end])
// 	}
// 	return batches
// }

// //decodeSignature decodes and returns a signature
// func decodeSignature(sig string) (string, error) {
// 	decBytes, err := decode(sig)
// 	if err != nil {
// 		return "", errors.Wrap(err, "could not decode signature")
// 	}
// 	return hex.EncodeToString(decBytes), nil
// }

// func (c Conts) string() string {
// 	res, _ := json.Marshal(c)
// 	return string(res)
// }

// // unmarshalString unmarshals the bytes received as a parameter, into the type string.
// func unmarshalString(v []byte) (string, error) {
// 	var str string
// 	err := json.Unmarshal(v, &str)
// 	if err != nil {
// 		return str, errors.Wrap(err, "could not unmarshal bytes to string")
// 	}
// 	return str, nil
// }

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
			contents, err := t.unforgeRevealOperation(hexString)
			if err != nil {
				return branch, contents, errors.New("failed to unforge operation")
			}
		case "6c":
		case "6d":
		case "6e":
		default:
			return branch, contents, fmt.Errorf("failed to unforge operation: transaction operation unkown %s", result)
		}
	}

	return branch, contents, nil
}

func (t *GoTezos) unforgeRevealOperation(hexString string) (Contents, error) {
	result, rest := split(hexString, 42)
	source, err := parseTzAddress(result)
	if err != nil {
		return Contents{}, errors.Wrap(err, "failed to unforge reveal operation")
	}

	var contents Contents

	zEndIndex, err := findZarithEndIndex(rest)
	if err != nil {
		return Contents{}, errors.Wrap(err, "failed to unforge reveal operation")
	}
	result, rest = split(rest, zEndIndex)
	zBigNum, err := zarithToBigNumber(result)
	if err != nil {
		return Contents{}, errors.Wrap(err, "failed to unforge reveal operation")
	}
	contents.Fee = zBigNum

	zEndIndex, err = findZarithEndIndex(rest)
	if err != nil {
		return Contents{}, errors.Wrap(err, "failed to unforge reveal operation")
	}
	result, rest = split(rest, zEndIndex)
	zBigNum, err = zarithToBigNumber(result)
	if err != nil {
		return Contents{}, errors.Wrap(err, "failed to unforge reveal operation")
	}
	contents.Counter = zBigNum

	zEndIndex, err = findZarithEndIndex(rest)
	if err != nil {
		return Contents{}, errors.Wrap(err, "failed to unforge reveal operation")
	}
	result, rest = split(rest, zEndIndex)
	zBigNum, err = zarithToBigNumber(result)
	if err != nil {
		return Contents{}, errors.Wrap(err, "failed to unforge reveal operation")
	}
	contents.GasLimit = zBigNum

	zEndIndex, err = findZarithEndIndex(rest)
	if err != nil {
		return Contents{}, errors.Wrap(err, "failed to unforge reveal operation")
	}
	result, rest = split(rest, zEndIndex)
	zBigNum, err = zarithToBigNumber(result)
	if err != nil {
		return Contents{}, errors.Wrap(err, "failed to unforge reveal operation")
	}
	contents.StorageLimit = zBigNum

	zEndIndex, err = findZarithEndIndex(rest)
	if err != nil {
		return Contents{}, errors.Wrap(err, "failed to unforge reveal operation")
	}
	result, rest = split(rest, zEndIndex)
	phk, err := parsePublicKey(result)
	if err != nil {
		return Contents{}, errors.Wrap(err, "failed to unforge reveal operation")
	}
	contents.Phk = phk

	return contents, nil
}

// public unforgeRevealOperation(hexString: string): { tezosRevealOperation: TezosRevealOperation; rest: string } {
//     let { result, rest }: { result: string; rest: string } = this.splitAndReturnRest(hexString, 42)
//     const source: string = this.parseTzAddress(result)

//       // fee, counter, gas_limit, storage_limit
//     ;({ result, rest } = this.splitAndReturnRest(rest, this.findZarithEndIndex(rest)))
//     const fee: BigNumber = this.zarithToBigNumber(result)
//     ;({ result, rest } = this.splitAndReturnRest(rest, this.findZarithEndIndex(rest)))
//     const counter: BigNumber = this.zarithToBigNumber(result)
//     ;({ result, rest } = this.splitAndReturnRest(rest, this.findZarithEndIndex(rest)))
//     const gasLimit: BigNumber = this.zarithToBigNumber(result)
//     ;({ result, rest } = this.splitAndReturnRest(rest, this.findZarithEndIndex(rest)))
//     const storageLimit: BigNumber = this.zarithToBigNumber(result)
//     ;({ result, rest } = this.splitAndReturnRest(rest, 66))
//     const publicKey: string = this.parsePublicKey(result)

//     return {
//       tezosRevealOperation: {
//         kind: TezosOperationType.REVEAL,
//         fee: fee.toFixed(),
//         gas_limit: gasLimit.toFixed(),
//         storage_limit: storageLimit.toFixed(),
//         counter: counter.toFixed(),
//         public_key: publicKey,
//         source
//       },
//       rest
//     }
//   }

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
