package gotezos

// import (
// 	"encoding/hex"
// 	"encoding/json"
// 	"fmt"
// 	"strconv"

// 	"golang.org/x/crypto/blake2b"

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
// func (t *GoTezos) ForgeOperations(wallet Wallet, fee, gaslimit int, operations ...Contents) ([]string, error) {
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
