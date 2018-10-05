package goTezos

/*
Author: DefinitelyNotAGoat/MagicAglet
Version: 0.0.1
Description: The Tezos API written in GO, for easy development.
License: MIT
*/

import (
	"gopkg.in/resty.v1"
)

var RPCURL string

func SetRPCURL(url string) {
	RPCURL = url
}

/*
Description: A function that executes an rpc get arg
Param args ([]string): Arguments to be executed
Returns (string): Returns the output of the executed command as a string
*/
func TezosRPCGet(arg string) ([]byte, error) {
	get := RPCURL + arg
	resp, err := resty.R().Get(get)
	if err != nil {
		return resp.Body(), err
	}

	return resp.Body(), nil
}

/*
Description: A function that executes an rpc get arg
Param args ([]string): Arguments to be executed
Returns (string): Returns the output of the executed command as a string
*/
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
