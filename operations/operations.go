package operations

import (
	"encoding/hex"
	"encoding/json"
	"strconv"

	"golang.org/x/crypto/blake2b"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ed25519"

	"github.com/DefinitelyNotAGoat/go-tezos/account"
	"github.com/DefinitelyNotAGoat/go-tezos/block"
	tzc "github.com/DefinitelyNotAGoat/go-tezos/client"
	"github.com/DefinitelyNotAGoat/go-tezos/crypto"
	"github.com/DefinitelyNotAGoat/go-tezos/delegate"
)

var (
	// maxBatchSize tells how many Transactions per batch are allowed.
	maxBatchSize = 200
)

// OperationService is a struct wrapper for operation related functions
type OperationService struct {
	blockService block.TezosBlockService
	tzclient     tzc.TezosClient
}

// Conts is helper structure to build out the contents of a a transfer operation to post to the Tezos RPC
type Conts struct {
	Contents []block.Contents `json:"contents"`
	Branch   string           `json:"branch"`
}

// Transfer a complete transfer request
type Transfer struct {
	Conts
	Protocol  string `json:"protocol"`
	Signature string `json:"signature"`
}

// NewOperationService returns a New Operation Service
func NewOperationService(blockService block.TezosBlockService, tzclient tzc.TezosClient) *OperationService {
	return &OperationService{
		blockService: blockService,
		tzclient:     tzclient,
	}
}

// CreateBatchPayment forges batch payments and returns them ready to inject to a Tezos RPC. PaymentFee must be expressed in mutez and the max batch size allowed is 200.
func (o *OperationService) CreateBatchPayment(payments []delegate.Payment, wallet account.Wallet, paymentFee int, gasLimit int, batchSize int) ([]string, error) {

	if batchSize > maxBatchSize {
		batchSize = maxBatchSize
	}

	var operationSignatures []string

	// Get current branch head
	blockHead, err := o.blockService.GetHead()
	if err != nil {
		return operationSignatures, errors.Wrap(err, "could not create batch payment")
	}

	// Get the counter for the payment address and increment it
	counter, err := o.getAddressCounter(wallet.Address)
	if err != nil {
		return operationSignatures, errors.Wrap(err, "could not create batch payment")
	}
	counter++

	// Split our slice of []Payment into batches
	batches := o.splitPaymentIntoBatches(payments, batchSize)
	operationSignatures = make([]string, len(batches))

	for k := range batches {

		// Convert (ie: forge) each 'Payment' into an actual Tezos transfer operation
		operationBytes, operationContents, newCounter, err := o.forgeOperationBytes(blockHead.Hash, counter, wallet, batches[k], paymentFee, gasLimit)
		if err != nil {
			return operationSignatures, errors.Wrap(err, "could not create batch payment")
		}
		counter = newCounter

		// Sign gt batch of operations with the secret key; return that signature
		edsig, err := o.signOperationBytes(operationBytes, wallet)
		if err != nil {
			return operationSignatures, errors.Wrap(err, "could not create batch payment")
		}

		// Extract and decode the bytes of the signature
		decodedSignature, err := o.decodeSignature(edsig)
		if err != nil {
			return operationSignatures, errors.Wrap(err, "could not create batch payment")
		}

		decodedSignature = decodedSignature[10:(len(decodedSignature))]

		// The signed bytes of gt batch
		fullOperation := operationBytes + decodedSignature

		// We can validate gt batch against the node for any errors
		err = o.preApplyOperations(operationContents, edsig, blockHead)
		if err != nil {
			return operationSignatures, errors.Wrap(err, "could not create batch payment")
		}
		// Add the signature (raw operation bytes & signature of operations) of gt batch of transfers to the returning slice
		// gt will be used to POST to /injection/operation
		operationSignatures[k] = fullOperation

	}

	return operationSignatures, nil
}

//Sign previously forged Operation bytes using secret key of wallet
func (o *OperationService) signOperationBytes(operationBytes string, wallet account.Wallet) (string, error) {

	opBytes, err := hex.DecodeString(operationBytes)
	if err != nil {
		return "", errors.Wrap(err, "could not sign operation bytes")
	}
	opBytes = append(crypto.Prefix_watermark, opBytes...)

	// Generic hash of 32 bytes
	genericHash, err := blake2b.New(32, []byte{})
	if err != nil {
		return "", errors.Wrap(err, "could not sign operation bytes")
	}

	// Write operation bytes to hash
	i, err := genericHash.Write(opBytes)

	if err != nil {
		return "", errors.Wrap(err, "could not sign operation bytes")
	}
	if i != len(opBytes) {
		return "", errors.Errorf("could not sign operation, generic hash length %d does not match bytes length %d", i, len(opBytes))
	}

	finalHash := genericHash.Sum([]byte{})

	// Sign the finalized generic hash of operations and b58 encode
	sig := ed25519.Sign(wallet.Kp.PrivKey, finalHash)
	//sig := sodium.Bytes(finalHash).SignDetached(wallet.Kp.PrivKey)
	edsig := crypto.B58cencode(sig, crypto.Prefix_edsig)

	return edsig, nil
}

func (o *OperationService) forgeOperationBytes(branchHash string, counter int, wallet account.Wallet, batch []delegate.Payment, paymentFee int, gaslimit int) (string, Conts, int, error) {

	var contents Conts
	var combinedOps []block.Contents

	//left here to display how to reveal a new wallet (needs funds to be revealed!)
	/**
	  combinedOps = append(combinedOps, StructContents{Kind: "reveal", PublicKey: wallet.pk , Source: wallet.address, Fee: "0", GasLimit: "127", StorageLimit: "0", Counter: strCounter})
	  counter++
	**/

	for k := range batch {

		if batch[k].Amount > 0 {

			operation := block.Contents{
				Kind:         "transaction",
				Source:       wallet.Address,
				Fee:          strconv.Itoa(paymentFee),
				GasLimit:     strconv.Itoa(gaslimit),
				StorageLimit: "0",
				Amount:       strconv.FormatFloat(crypto.RoundPlus(batch[k].Amount, 0), 'f', -1, 64),
				Destination:  batch[k].Address,
				Counter:      strconv.Itoa(counter),
			}
			combinedOps = append(combinedOps, operation)
			counter++
		}
	}
	contents.Contents = combinedOps
	contents.Branch = branchHash

	var opBytes string

	forge := "/chains/main/blocks/head/helpers/forge/operations"
	output, err := o.tzclient.Post(forge, contents.string())
	if err != nil {
		return "", contents, counter, errors.Wrapf(err, "could not forge operation '%s' with contents '%s'", forge, contents.string())
	}

	err = json.Unmarshal(output, &opBytes)
	if err != nil {
		return "", contents, counter, errors.Wrapf(err, "could not forge operation '%s' with contents '%s'", forge, contents.string())
	}

	return opBytes, contents, counter, nil
}

// Pre-apply an operation, or batch of operations, to a Tezos node to ensure correctness
func (o *OperationService) preApplyOperations(paymentOperations Conts, signature string, blockHead block.Block) error {

	// Create a full transfer request
	var transfer Transfer
	transfer.Signature = signature
	transfer.Contents = paymentOperations.Contents
	transfer.Branch = blockHead.Hash
	transfer.Protocol = blockHead.Protocol

	// RPC says outer element must be JSON array
	var transfers = []Transfer{transfer}

	// Convert object to JSON string
	transfersOp, err := json.Marshal(transfers)
	if err != nil {
		return errors.Wrap(err, "could not preapply operations, could not marshal into json")
	}

	// POST the JSON to the RPC
	query := "/chains/main/blocks/head/helpers/preapply/operations"
	_, err = o.tzclient.Post(query, string(transfersOp))
	if err != nil {
		return errors.Wrapf(err, "could not preapply operations '%s' with contents '%s'", query, string(transfersOp))
	}

	return nil
}

// InjectOperation injects an signed operation string and returns the response
func (o *OperationService) InjectOperation(op string) ([]byte, error) {
	post := "/injection/operation"
	jsonBytes, err := json.Marshal(op)
	if err != nil {
		return nil, errors.Wrapf(err, "could not inject operation '%s' with contents '%s'", post, string(jsonBytes))
	}
	resp, err := o.tzclient.Post(post, string(jsonBytes))
	if err != nil {
		return resp, errors.Wrapf(err, "could not inject operation '%s' with contents '%s'", post, string(jsonBytes))
	}
	return resp, nil
}

//Getting the Counter of an address from the RPC
func (o *OperationService) getAddressCounter(address string) (int, error) {
	rpc := "/chains/main/blocks/head/context/contracts/" + address + "/counter"
	resp, err := o.tzclient.Get(rpc, nil)
	if err != nil {
		return 0, errors.Wrapf(err, "could not get address counter '%s'", rpc)
	}
	rtnStr, err := unmarshalString(resp)
	if err != nil {
		return 0, errors.Wrapf(err, "could not get address counter '%s'", rpc)
	}
	counter, err := strconv.Atoi(rtnStr)
	if err != nil {
		return 0, errors.Wrapf(err, "could not get address counter '%s'", rpc)
	}
	return counter, nil
}

func (o *OperationService) splitPaymentIntoBatches(rewards []delegate.Payment, batchSize int) [][]delegate.Payment {
	var batches [][]delegate.Payment
	for i := 0; i < len(rewards); i += batchSize {
		end := i + batchSize
		if end > len(rewards) {
			end = len(rewards)
		}
		batches = append(batches, rewards[i:end])
	}
	return batches
}

// GetBlockOperationHashes returns list of operations in block at specific level
func (o *OperationService) GetBlockOperationHashes(id interface{}) ([]string, error) {

	var operations []string
	block, err := o.blockService.Get(id)
	if err != nil {
		return operations, errors.Wrap(err, "could not get operation hashes")
	}

	query := "/chains/main/blocks/" + block.Hash + "/operation_hashes"
	resp, err := o.tzclient.Get(query, nil)
	if err != nil {
		return operations, errors.Wrapf(err, "could not get operation hashes '%s'", query)
	}

	operations, err = unmarshalMultiStrJSON(resp)
	if err != nil {
		return operations, errors.Wrapf(err, "could not get operation hashes '%s'", query)
	}

	return operations, nil
}

//Helper function to return the decoded signature
func (o *OperationService) decodeSignature(sig string) (string, error) {
	decBytes, err := crypto.Decode(sig)
	if err != nil {
		return "", errors.Wrap(err, "could not decode signature")
	}
	return hex.EncodeToString(decBytes), nil
}

func (c Conts) string() string {
	res, _ := json.Marshal(c)
	return string(res)
}

// unmarshalMultiStrJSON unmarhsels bytes into OperationHashes
func unmarshalMultiStrJSON(v []byte) ([]string, error) {
	dops := [][]string{}
	ops := []string{}

	err := json.Unmarshal(v, &dops)
	if err != nil {
		return ops, errors.Wrap(err, "could not unmarshal bytes into []string")
	}

	for _, i := range dops {
		ops = append(ops, i...)
	}
	return ops, nil
}

// unmarshalString unmarshals the bytes received as a parameter, into the type string.
func unmarshalString(v []byte) (string, error) {
	var str string
	err := json.Unmarshal(v, &str)
	if err != nil {
		return str, errors.Wrap(err, "could not unmarshal bytes to string")
	}
	return str, nil
}
