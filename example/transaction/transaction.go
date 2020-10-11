package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"strconv"

	"github.com/goat-systems/go-tezos/v3/forge"
	"github.com/goat-systems/go-tezos/v3/keys"
	"github.com/goat-systems/go-tezos/v3/rpc"
)

func main() {
	key, err := keys.NewKey(keys.NewKeyInput{
		Esk:      "edesk...",
		Password: "password",
	})
	if err != nil {
		fmt.Printf("failed to import keys: %s\n", err.Error())
		os.Exit(1)
	}

	client, err := rpc.New("http://tezos_rpc_addr")
	if err != nil {
		fmt.Printf("failed to initialize rpc client: %s\n", err.Error())
		os.Exit(1)
	}

	head, err := client.Head()
	if err != nil {
		fmt.Printf("failed to get head block: %s\n", err.Error())
		os.Exit(1)
	}

	counter, err := client.Counter(rpc.CounterInput{
		Blockhash: head.Hash,
		Address:   key.PubKey.GetPublicKeyHash(),
	})
	if err != nil {
		fmt.Printf("failed to get counter: %s\n", err.Error())
		os.Exit(1)
	}
	counter++

	transaction := rpc.Transaction{
		Source:      key.PubKey.GetPublicKey(),
		Fee:         "2941",
		GasLimit:    "26283",
		Counter:     strconv.Itoa(counter),
		Amount:      "0",
		Destination: "<some_dest>",
	}

	op, err := forge.Encode(head.Hash, transaction.ToContent())
	if err != nil {
		fmt.Printf("failed to forge transaction: %s\n", err.Error())
		os.Exit(1)
	}

	signature, err := key.Sign(keys.SignInput{
		Message: op,
	})
	if err != nil {
		fmt.Printf("failed to sign operation: %s\n", err.Error())
		os.Exit(1)
	}

	ophash, err := client.InjectionOperation(rpc.InjectionOperationInput{
		Operation: fmt.Sprintf("%s%s", op, hex.EncodeToString(signature.Bytes)),
	})
	if err != nil {
		fmt.Printf("failed to inject: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Println(ophash)
}
