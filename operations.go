package gotezos

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"strconv"

	"github.com/Messer4/base58check"
)

var (
	// How many Transactions per batch are injected. I recommend 100. Now 30 for easier testing
	batchSize = 30

	// For (de)constructing addresses
	tz1   = []byte{6, 161, 159}
	edsk  = []byte{43, 246, 78, 7}
	edsk2 = []byte{13, 15, 58, 7}
	edpk  = []byte{13, 15, 37, 217}
	edesk = []byte{7, 90, 60, 179, 41}
)

// Pre-apply an operation, or batch of operations, to a Tezos node to ensure correctness
func (gt *GoTezos) preApplyOperations(paymentOperations Conts, signature string, blockHead Block) error {

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

	if gt.debug {
		fmt.Println("\n== preApplyOperations Submit:", string(transfersOp))
	}

	// POST the JSON to the RPC
	preApplyResp, err := gt.PostResponse("/chains/main/blocks/head/helpers/preapply/operations", string(transfersOp))
	if err != nil {
		return err
	}

	if gt.debug {
		fmt.Println("\n== preApplyOperations Result:", string(preApplyResp.Bytes))
	}

	return nil
}

func (gt *GoTezos) InjectOperation(ops []string) ([][]byte, error) {
	post := "/injection/operation"
	responses := [][]byte{}
	for _, op := range ops {
		jsonBytes, err := json.Marshal(op)
		if err != nil {
			return responses, err
		}
		resp, err := gt.PostResponse(post, string(jsonBytes))
		responses = append(responses, resp.Bytes)
		if err != nil {
			return responses, err
		}
	}
	return responses, nil
}

//Getting the Counter of an address from the RPC
func (gt *GoTezos) getAddressCounter(address string) (int, error) {
	rpc := "/chains/main/blocks/head/context/contracts/" + address + "/counter"
	resp, err := gt.GetResponse(rpc, "{}")
	if err != nil {
		return 0, err
	}
	rtnStr, err := unmarshalString(resp.Bytes)
	if err != nil {
		return 0, err
	}
	counter, err := strconv.Atoi(rtnStr)
	return counter, err
}

func (gt *GoTezos) splitPaymentIntoBatches(rewards []Payment) [][]Payment {
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

//Helper function to return the decoded signature
func (gt *GoTezos) decodeSignature(sig string) string {
	decBytes, err := base58check.Decode(sig)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	return hex.EncodeToString(decBytes)
}

//Helper Function to get the right format for wallet.
func (gt *GoTezos) b58cencode(payload []byte, prefix []byte) string {
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

func (gt *GoTezos) b58cdecode(payload string, prefix []byte) []byte {
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
