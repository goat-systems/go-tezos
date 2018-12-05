package main

import (
	"fmt"
	"github.com/DefinitelyNotAGoat/go-tezos"
)

//Used to show how to use added features to Create Batch Payments. The signed operations are not injeced but rather returned as an array.
func main() {
	gt := goTezos.NewGoTezos()
	gt.AddNewClient(goTezos.NewTezosRPCClient("localhost",":8732"))

	block,_ := gt.GetBlockAtLevel(1000)
	fmt.Println(block.Hash)
}
