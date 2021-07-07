[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://godoc.org/github.com/goat-systems/go-tezos/v4)
# A Tezos Go Library

Go Tezos is a GoLang driven library for your Tezos node. This library has received a grant from the Tezos Foundation to ensure it's continuous development through 2020.

## Installation

Get GoTezos 
```
go get github.com/goat-systems/go-tezos/v4
```

### Getting A Block

```go
package main

import (
	"fmt"
	"os"
	"github.com/goat-systems/go-tezos/v4/rpc"
)

func main() {
	client, err := rpc.New("http://127.0.0.1:8732")
	if err != nil {
		fmt.Printf("failed tp connect to network: %v", err)
	}

	resp, block, err := client.Block(&rpc.BlockIDHead{})
	if err != nil {
		fmt.Printf("failed to get (%s) head block: %s\n", resp.Status(), err.Error())
		os.Exit(1)
	}
	fmt.Println(block)
}
```

### Getting a Cycle
```
	resp, cycle, err := client.Cycle(50)
	if err != nil {
		fmt.Printf("failed to get (%s) cycle: %s\n", resp.Status(), err.Error())
		os.Exit(1)
	}
	fmt.Println(cycle)
```

### More Examples
You can find more examples by looking through the unit tests and integration tests in each package. [Here](example/transaction/transaction.go) is an example on
how to forge and inject an operation. 

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
