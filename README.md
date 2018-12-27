# goTezos: A Tezos Go Library

The purpose of this library is to allow developers to build go driven applications for Tezos of the Tezos RPC. This library is a work in progress, and not complete. 

More robust documentation will come soon.

## Installation

Install pkg_config (debian example below):
```
sudo apt-get install pkg_config
```

Install [libsoidum](https://libsodium.gitbook.io/doc/installation)

Get goTezos 
```
go get github.com/DefinitelyNotAGoat/goTezos
```

## goTezos Documentation

[GoDoc](https://godoc.org/github.com/DefinitelyNotAGoat/goTezos)

The goTezos Library requires you to set the RPC URL for a node to query. 


Usage:

```
package main

import (
	"fmt"
	goTezos "github.com/DefinitelyNotAGoat/go-tezos"
)

func main() {
	gt := goTezos.NewGoTezos()
	gt.AddNewClient(goTezos.NewTezosRPCClient("localhost",":8732"))

	block,_ := gt.GetBlockAtLevel(1000)
	fmt.Println(block.Hash)
}
```

I will create a wiki shortly describing the functions available.

## Contributers

* [**DefinitelyNotAGoat**](https://github.com/DefinitelyNotAGoat)
* [**Magic_Gum**](https://github.com/fkbenjamin)
* [**Johann**](https://github.com/tulpenhaendler)
* [**utdrmac**](https://github.com/utdrmac)

See the list of [contributors](https://github.com/DefinitelyNotAGoat/goTezos/graphs/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
