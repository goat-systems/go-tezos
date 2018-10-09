//Package goTezos exposes the Tezos RPC API in goLang.
package goTezos

import (
	"gopkg.in/resty.v1"
)

var RPCURL string

//Set the RPC API URL pointing to a Tezos Node. 
//This function is required for goTezos to perform queries. 
func SetRPCURL(url string) {
	RPCURL = url
}

//Take an RPC url according to the Tezos documentation and 
//performs a GET request. 
func TezosRPCGet(arg string) ([]byte, error) {
	get := RPCURL + arg
	resp, err := resty.R().Get(get)
	if err != nil {
		return resp.Body(), err
	}

	return resp.Body(), nil
}

//Take an RPC url according to the Tezos documentation, and a Conts type,
//which is the contents of a Tezos operation. 
func TezosRPCPost(arg string, post Conts) ([]byte, error) {
	url := RPCURL + arg
	resp, err := resty.R().
		SetHeader("Content-Type", "application/json").
		SetBody(post).
		Post(url)
	if err != nil {
		return resp.Body(), err
	}

	return resp.Body(), nil
}
