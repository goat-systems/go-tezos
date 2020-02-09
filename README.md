[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://godoc.org/github.com/goat-systems/go-tezos/v2)
# A Tezos Go Library

Go Tezos is a GoLang driven library for your Tezos node. 

## Installation

Get goTezos 
```
go get github.com/DefinitelyNotAGoat/go-tezos/v2
```

### Getting A Block

```
package main

import (
	"fmt"
	goTezos "github.com/DefinitelyNotAGoat/go-tezos/v2"
)

func main() {
	gt, err := goTezos.New("http://127.0.0.1:8732")
	if err != nil {
		fmt.Printf("could not connect to network: %v", err)
	}

	block, err := gt.Block(1000)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(block)
}
```

### Getting a Cycle
```
	cycle, err := gt.Cycle(50)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cycle)
```

### More Documentation
See [github pages](https://definitelynotagoat.github.io/go-tezos/v2/)

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
