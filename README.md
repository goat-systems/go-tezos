[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://godoc.org/github.com/DefinitelyNotAGoat/go-tezos)
# A Tezos Go Library

Go Tezos is a GoLang driven library for your Tezos node. 

## Installation

Get goTezos 
```
go get github.com/DefinitelyNotAGoat/go-tezos
```

## Quick Start 
Go Tezos is split into multiple services underneath to help organize it's functionality and also makes the library easier to maintain. 

To understand how Go Tezos works, take a look at the GoTezos Structure: 
```
type GoTezos struct {
	Client    tzc.TezosClient
	Constants network.Constants
	Block     block.TezosBlockService
	Snapshot  snapshot.TezosSnapshotService
	Cycle     cycle.TezosCycleService
	Account   account.TezosAccountService
	Delegate  delegate.TezosDelegateService
	Network   network.TezosNetworkService
	Operation operations.TezosOperationsService
	Contract  contracts.TezosContractsService
	Node      node.TezosNodeService
}
```
You can see GoTezos is a wrapper for several services such as `block`,  `Snapshot`, `Cycle`, `Account`, `Delegate`, `Network`, `Operation`, `Node`, and `Contract`.
Each service has it's own set of functions. You can see examples of using the `Block` and `SnapShot` service below.


### Getting A Block

```
package main

import (
	"fmt"
	goTezos "github.com/DefinitelyNotAGoat/go-tezos"
)

func main() {
	gt, err := goTezos.NewGoTezos("http://127.0.0.1:8732")
	if err != nil {
		fmt.Printf("could not connect to network: %v", err)
	}

	block, err := gt.Block.Get(1000)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(block)
}
```

### Getting a Snapshot For A Cycle
```
	snapshot, err := gt.Snapshot.Get(50)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(snapshot)
```

### More Documentation
See [github pages](https://definitelynotagoat.github.io/go-tezos/)

## Contributing

### The Makefile
The makefile is there as a helper to run quality code checks. To run vet and staticchecks please run: 
```
make checks
```

## Contributers: A Special Thank You

* [**BrianBland**](https://github.com/BrianBland)
* [**utdrmac**](https://github.com/utdrmac)
* [**Magic_Gum**](https://github.com/fkbenjamin)
* [**Johann**](https://github.com/tulpenhaendler)
* [**leopoldjoy**](https://github.com/leopoldjoy)
* [**RomarQ**](https://github.com/RomarQ)
* [**surzm**](https://github.com/surzm)
* [**fredcy**](https://github.com/fredcy)

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
