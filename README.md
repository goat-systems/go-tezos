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
	client    *client
	Constants NetworkConstants
	Block     *BlockService
	SnapShot  *SnapShotService
	Cycle     *CycleService
	Account   *AccountService
	Delegate  *DelegateService
	Network   *NetworkService
	Operation *OperationService
	Contract  *ContractService
}
```
You can see GoTezos is a wrapper for an http client, and services such as `block`,  `SnapShot`, `Cycle`, `Account`, `Delegate`, `Network`, `Operation`, and `Contract`.
Each service has it's own set of functions. You can see examples of using the `Block` and `SnapShot` service below.


### Getting A Block

```
package main

import (
	"fmt"
	goTezos "github.com/DefinitelyNotAGoat/go-tezos"
)

func main() {
	gt, err := NewGoTezos("http://127.0.0.1:8732")
	if err != nil {
		t.Errorf("could not connect to network: %v", err)
	}

	block, err := gt.Block.Get(1000)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(block)
}
```

### Getting a Snap Shot For A Cycle
```
	snapshot, err := gt.SnapShot.Get(50)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(snapshot)
```

### More Documentation
See [github pages](https://definitelynotagoat.github.io/go-tezos/)

## Contributers

### Note

Because this project is gaining some traction I made the decision to squash all commits into a single commit, so that we can have a clean and well organized
commit history going forward. By doing that, the commit history doesn't reflect the contributions of some very helpful people apart of Tezos community and the go-tezos project. Please take a look at the pull request history to see their individual contributions. 


### Special Thank You

I want to make sure the following people are recognized and give a special thank you to some of the original contributers to go-tezos:  

* [**BrianBland**](https://github.com/BrianBland)
* [**utdrmac**](https://github.com/utdrmac)
* [**Magic_Gum**](https://github.com/fkbenjamin)
* [**Johann**](https://github.com/tulpenhaendler)
* [**leopoldjoy**](https://github.com/leopoldjoy)

## Pull Requests
Go Tezos is a relatively large project and has the potential to be larger, because of that it's important to maintain quality code and PR's. Please review the pull request guide lines [here](PULL_REQUEST_GUIDE.md).

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
