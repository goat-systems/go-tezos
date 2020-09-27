[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://godoc.org/github.com/goat-systems/go-tezos/v2)
# A Tezos Go Library

Go Tezos is a GoLang driven library for your Tezos node. This library has received a grant from the Tezos Foundation to ensure it's continuous development through 2020. 

## Installation

Get goTezos 
```
go get github.com/goat-systems/go-tezos/v3
```

### Getting A Block

```
package main

import (
	"fmt"
	goTezos "github.com/goat-systems/go-tezos/v3/rpc"
)

func main() {
	rpc, err := client.New("http://127.0.0.1:8732")
	if err != nil {
		fmt.Printf("could not connect to network: %v", err)
	}

	block, err := rpc.Block(1000)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(block)
}
```

### Getting a Cycle
```
	cycle, err := rpc.Cycle(50)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cycle)
```

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
