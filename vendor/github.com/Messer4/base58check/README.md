# base58check
[![Build Status](https://travis-ci.org/Messer4/base58check.svg?branch=master)](https://travis-ci.org/Messer4/base58check)
[![GoDoc](https://godoc.org/github.com/Messer4/base58check?status.svg)](https://godoc.org/github.com/Messer4/base58check)
[![Go Report Card](https://goreportcard.com/badge/github.com/Messer4/base58check)](https://goreportcard.com/report/github.com/Messer4/base58check)

This package in Go provides functions to encode and decode in `base58check`, a specific `base58` encoding format for encoding Bitcoin addresses.

Functions:
```go
func Encode([]byte) (string, error) // takes the version and data as hexadecimal strings and returns the encoded string
func Decode(string) ([]byte, error) // takes the encoded string and returns the decoded version prepended hexadecimal string
```

### Installation

```bash
go get github.com/Messer4/base58check
```

### Usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/Messer4/base58check"
)

func main() {
	encoded, err := base58check.Encode("80", "44D00F6EB2E5491CD7AB7E7185D81B67A23C4980F62B2ED0914D32B7EB1C5581")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(encoded) // 5JLbJxi9koHHvyFEAERHLYwG7VxYATnf8YdA9fiC6kXMghkYXpk

	decoded, err := base58check.Decode("1mayif3H2JDC62S4N3rLNtBNRAiUUP99k")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(decoded) // 00086eaa677895f92d4a6c5ef740c168932b5e3f44
}

```

### References

+ [Base58Check encoding](https://en.bitcoin.it/wiki/Base58Check_encoding)
