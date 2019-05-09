package gotezos

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"strconv"

	"golang.org/x/crypto/blake2b"

	"github.com/Messer4/base58check"
	"golang.org/x/crypto/ed25519"
)

var (
	// How many Transactions per batch are injected. I recommend 100. Now 30 for easier testing
	batchSize = 100

	// For (de)constructing addresses
	tz1   = []byte{6, 161, 159}
	edsk  = []byte{43, 246, 78, 7}
	edsk2 = []byte{13, 15, 58, 7}
	edpk  = []byte{13, 15, 37, 217}
	edesk = []byte{7, 90, 60, 179, 41}
)

// OperationService is a struct wrapper for operation related functions
type OperationService struct {
	gt *GoTezos
}

// Conts is helper structure to build out the contents of a a transfer operation to post to the Tezos RPC
type Conts struct {
	Contents []StructContents `json:"contents"`
	Branch   string           `json:"branch"`
}

// Transfer a complete transfer request
type Transfer struct {
	Conts
	Protocol  string `json:"protocol"`
	Signature string `json:"signature"`
}

// NewOperationService returns a New Operation Service
func (gt *GoTezos) newOperationService() *OperationService {
	return &OperationService{gt: gt}
}

// CreateBatchPayment forges batch payments and returns them ready to inject to a Tezos RPC. PaymentFee must be expressed in mutez.
func (o *OperationService) CreateBatchPayment(payments []Payment, wallet Wallet, paymentFee int, gaslimit int) ([]string, error) {

	var operationSignatures []string

	// Get current branch head
	blockHead, err := o.gt.Block.GetHead()
	if err != nil {
		return operationSignatures, fmt.Errorf("could not create batch payment: %v", err)
	}

	// Get the counter for the payment address and increment it
	counter, err := o.getAddressCounter(wallet.Address)
	if err != nil {
		return operationSignatures, fmt.Errorf("could not create batch payment: %v", err)
	}
	counter++

	// Split our slice of []Payment into batches
	batches := o.splitPaymentIntoBatches(payments)
	operationSignatures = make([]string, len(batches))

	for k := range batches {

		// Convert (ie: forge) each 'Payment' into an actual Tezos transfer operation
		operationBytes, operationContents, newCounter, err := o.forgeOperationBytes(blockHead.Hash, counter, wallet, batches[k], paymentFee, gaslimit)
		if err != nil {
			return operationSignatures, fmt.Errorf("could not create batch payment: %v", err)
		}
		counter = newCounter

		// Sign gt batch of operations with the secret key; return that signature
		edsig, err := o.signOperationBytes(operationBytes, wallet)
		if err != nil {
			return operationSignatures, fmt.Errorf("could not create batch payment: %v", err)
		}

		// Extract and decode the bytes of the signature
		decodedSignature := o.decodeSignature(edsig)
		decodedSignature = decodedSignature[10:(len(decodedSignature))]

		// The signed bytes of gt batch
		fullOperation := operationBytes + decodedSignature

		// We can validate gt batch against the node for any errors
		if err := o.preApplyOperations(operationContents, edsig, blockHead); err != nil {
			return operationSignatures, fmt.Errorf("could not create batch payment: failed to Pre-Apply: %v", err)
		}
		// Add the signature (raw operation bytes & signature of operations) of gt batch of transfers to the returnning slice
		// gt will be used to POST to /injection/operation
		operationSignatures[k] = fullOperation

	}

	return operationSignatures, nil
}

//Sign previously forged Operation bytes using secret key of wallet
func (o *OperationService) signOperationBytes(operationBytes string, wallet Wallet) (string, error) {

	//Prefixes
	edsigByte := []byte{9, 245, 205, 134, 18}
	watermark := []byte{3}

	opBytes, err := hex.DecodeString(operationBytes)
	if err != nil {
		return "", fmt.Errorf("could not sign operation: %v", err)
	}
	opBytes = append(watermark, opBytes...)

	// Generic hash of 32 bytes
	genericHash, err := blake2b.New(32, []byte{})

	// Write operation bytes to hash
	i, err := genericHash.Write(opBytes)
	if i != len(opBytes) || err != nil {
		return "", fmt.Errorf("could not sign operation: unable to write operations to generic hash")
	}
	finalHash := genericHash.Sum([]byte{})

	// Sign the finalized generic hash of operations and b58 encode
	sig := ed25519.Sign(wallet.Kp.PrivKey, finalHash)
	//sig := sodium.Bytes(finalHash).SignDetached(wallet.Kp.PrivKey)
	edsig := b58cencode(sig, edsigByte)

	return edsig, nil
}

func (o *OperationService) forgeOperationBytes(branchHash string, counter int, wallet Wallet, batch []Payment, paymentFee int, gaslimit int) (string, Conts, int, error) {

	var contents Conts
	var combinedOps []StructContents

	//left here to display how to reveal a new wallet (needs funds to be revealed!)
	/**
	  combinedOps = append(combinedOps, StructContents{Kind: "reveal", PublicKey: wallet.pk , Source: wallet.address, Fee: "0", GasLimit: "127", StorageLimit: "0", Counter: strCounter})
	  counter++
	**/

	for k := range batch {

		if batch[k].Amount > 0 {

			operation := StructContents{
				Kind:         "transaction",
				Source:       wallet.Address,
				Fee:          strconv.Itoa(paymentFee),
				GasLimit:     strconv.Itoa(gaslimit),
				StorageLimit: "0",
				Amount:       strconv.FormatFloat(roundPlus(batch[k].Amount, 0), 'f', -1, 64),
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
	output, err := o.gt.Post(forge, contents.string())
	if err != nil {
		return "", contents, counter, fmt.Errorf("could not forge operation: %v", err)
	}

	err = json.Unmarshal(output, &opBytes)
	if err != nil {
		return "", contents, counter, fmt.Errorf("could not forge operation: %v", err)
	}

	return opBytes, contents, counter, nil
}

// Pre-apply an operation, or batch of operations, to a Tezos node to ensure correctness
func (o *OperationService) preApplyOperations(paymentOperations Conts, signature string, blockHead Block) error {

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
		return err
	}

	// POST the JSON to the RPC
	_, err = o.gt.Post("/chains/main/blocks/head/helpers/preapply/operations", string(transfersOp))
	if err != nil {
		return fmt.Errorf("could not preapply operation: %v", err)
	}

	return nil
}

// InjectOperation injects an signed operation string and returns the response
func (o *OperationService) InjectOperation(op string) ([]byte, error) {
	post := "/injection/operation"
	jsonBytes, err := json.Marshal(op)
	if err != nil {
		return nil, fmt.Errorf("could not inject operation: %v", err)
	}
	resp, err := o.gt.Post(post, string(jsonBytes))
	if err != nil {
		return resp, fmt.Errorf("could not inject operation: %v", err)
	}
	return resp, nil
}

//Getting the Counter of an address from the RPC
func (o *OperationService) getAddressCounter(address string) (int, error) {
	rpc := "/chains/main/blocks/head/context/contracts/" + address + "/counter"
	resp, err := o.gt.Get(rpc, nil)
	if err != nil {
		return 0, fmt.Errorf("could not get address counter: %v", err)
	}
	rtnStr, err := unmarshalString(resp)
	if err != nil {
		return 0, fmt.Errorf("could not get address counter: %v", err)
	}
	counter, err := strconv.Atoi(rtnStr)
	if err != nil {
		return 0, fmt.Errorf("could not get address counter: %v", err)
	}
	return counter, nil
}

func (o *OperationService) splitPaymentIntoBatches(rewards []Payment) [][]Payment {
	var batches [][]Payment
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
	block, err := o.gt.Block.Get(id)
	if err != nil {
		return operations, fmt.Errorf("could not get operation hashes: %v", err)
	}

	query := "/chains/main/blocks/" + block.Hash + "/operation_hashes"
	resp, err := o.gt.Get(query, nil)
	if err != nil {
		return operations, fmt.Errorf("could not get operation hashes: %v", err)
	}

	operations, err = unmarshalMultiStrJSON(resp)
	if err != nil {
		return operations, fmt.Errorf("could not get operation hashes: %v", err)
	}

	return operations, nil
}

//Helper function to return the decoded signature
func (o *OperationService) decodeSignature(sig string) string {
	decBytes, err := base58check.Decode(sig)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	return hex.EncodeToString(decBytes)
}

//Helper Function to get the right format for wallet.
func b58cencode(payload []byte, prefix []byte) string {
	n := make([]byte, (len(prefix) + len(payload)))
	for k := range prefix {
		n[k] = prefix[k]
	}
	for l := range payload {
		n[l+len(prefix)] = payload[l]
	}
	b58c := base58check.Encode(n)
	return b58c
}

func b58cdecode(payload string, prefix []byte) []byte {
	b58c, _ := base58check.Decode(payload)
	return b58c[len(prefix):]
}

//Helper Functions to round float64
func roundPlus(f float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return round(f*shift) / shift
}

func round(f float64) float64 {
	return math.Floor(f + .5)
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
		return ops, err
	}

	for _, i := range dops {
		for _, j := range i {
			ops = append(ops, j)
		}
	}
	return ops, nil
}
