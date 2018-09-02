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
Description: Sends tezos from your specified wallet to another account.
Param amount (float64): Amount of tezos to be sent
Param toAddress (string): The address you are sending tezos to.
Param alias (string): The named alias assigned to your wallet you are sending out of.

****WARNING****
If not using the ledger there is nothing stopping this from actually sending Tezos.
With the ledger you have to physically confirm the transaction, without the ledger you don't.

BE CAREFUL WHEN CALLING THIS FUNCTION!!!!!
****WARNING****
*/
// func SendTezos(amount float64, toAddress string, alias string) error {
// 	strAmount := strconv.FormatFloat(amount, 'f', -1, 64)
// 	_, err := TezosDo("transfer", strAmount, "from", alias, "to", toAddress)
// 	if err != nil {
// 		return errors.New("Could not send " + strAmount + " XTZ from " + alias + " to " + toAddress + ": tezosDo(args ...string) failed: " + err.Error())
// 	}
// 	return nil
// }

// /*
// Description: Sends tezos from your specified wallet to another account, but makes you confirm the transaction.
// Param amount (float64): Amount of tezos to be sent
// Param toAddress (string): The address you are sending tezos to.
// Param alias (string): The named alias assigned to your wallet you are sending out of.
// */
// func SafeSendTezos(amount float64, toAddress string, alias string) error {
// 	strAmount := strconv.FormatFloat(amount, 'f', -1, 64)

// 	confirmStatement := "Send " + strAmount + " XTZ from " + alias + " to " + toAddress + "?"
// 	confirmation := askForConfirmation(confirmStatement)

// 	if confirmation {
// 		_, err := TezosDo("transfer", strAmount, "from", alias, "from", toAddress)
// 		if err != nil {
// 			return errors.New("Could not send " + strAmount + " XTZ from " + alias + " to " + toAddress + ": tezosDo(args ...string) failed: " + err.Error())
// 		}
// 	} else {
// 		return errors.New("Cancelled: Send " + strAmount + " XTZ from " + alias + " to " + toAddress)
// 	}
// 	return nil
// }

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
func TezosRPCPost(arg string) ([]byte, error) {
	get := RPCURL + arg
	resp, err := resty.R().Get(get)
	if err != nil {
		return resp.Body(), err
	}

	return resp.Body(), nil
}
