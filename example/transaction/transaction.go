package main

import (
	"fmt"
	"math/big"
	"os"
	"strconv"

	"github.com/completium/go-tezos/v4/forge"
	"github.com/completium/go-tezos/v4/keys"
	"github.com/completium/go-tezos/v4/rpc"
)

func main() {
	key, err := keys.FromEncryptedSecret("edesk...", "password")
	if err != nil {
		fmt.Printf("failed to import keys: %s\n", err.Error())
		os.Exit(1)
	}

	client, err := rpc.New("http://tezos_rpc_addr")
	if err != nil {
		fmt.Printf("failed to initialize rpc client: %s\n", err.Error())
		os.Exit(1)
	}

	resp, counter, err := client.ContractCounter(rpc.ContractCounterInput{
		BlockID:    &rpc.BlockIDHead{},
		ContractID: key.PubKey.GetAddress(),
	})
	if err != nil {
		fmt.Printf("failed to get (%s) counter: %s\n", resp.Status(), err.Error())
		os.Exit(1)
	}
	counter++

	big.NewInt(0).SetString("10000000000000000000000000000", 10)

	transaction := rpc.Transaction{
		Source:      key.PubKey.GetPublicKey(),
		Fee:         "2941",
		GasLimit:    "26283",
		Counter:     strconv.Itoa(counter),
		Amount:      "0",
		Destination: "<some_dest>",
	}

	resp, head, err := client.Block(&rpc.BlockIDHead{})
	if err != nil {
		fmt.Printf("failed to get (%s) head block: %s\n", resp.Status(), err.Error())
		os.Exit(1)
	}

	op, err := forge.Encode(head.Hash, transaction.ToContent())
	if err != nil {
		fmt.Printf("failed to forge transaction: %s\n", err.Error())
		os.Exit(1)
	}

	signature, err := key.SignGeneric(op)
	if err != nil {
		fmt.Printf("failed to sign operation: %s\n", err.Error())
		os.Exit(1)
	}

	resp, ophash, err := client.InjectionOperation(rpc.InjectionOperationInput{
		Operation: signature.AppendToHex(op),
	})
	if err != nil {
		fmt.Printf("failed to inject (%s): %s\n", resp.Status(), err.Error())
		os.Exit(1)
	}

	fmt.Println(ophash)
}
