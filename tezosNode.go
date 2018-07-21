package goTezos

/*
Author: DefinitelyNotAGoat/MagicAglet
Version: 0.0.1
Description: The Tezos API written in GO, for easy development.
License: MIT
*/

import (
  "fmt"
  "os"
  "os/exec"
  "strconv"
  "errors"
)

var TezosPath string

/*
Description: This library needs the TEZOSPATH enviroment variable to function
*/
func init() {
  var ok bool
  TezosPath, ok = os.LookupEnv("TEZOSPATH")
  if !ok {
	   fmt.Println("Error: Could not retrieve TEZOSPATH")
	   os.Exit(1)
  }
  TezosPath = TezosPath + "tezos-client"
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
func SendTezos(amount float64, toAddress string, alias string) error{
  strAmount := strconv.FormatFloat(amount, 'f', -1, 64)
  _, err := TezosDo("transfer", strAmount, "from", alias, "to", toAddress)
  if (err != nil){
    return errors.New("Could not send " + strAmount + " XTZ from " + alias + " to " + toAddress + ": tezosDo(args ...string) failed: " + err.Error())
  }
  return nil
}

/*
Description: Sends tezos from your specified wallet to another account, but makes you confirm the transaction.
Param amount (float64): Amount of tezos to be sent
Param toAddress (string): The address you are sending tezos to.
Param alias (string): The named alias assigned to your wallet you are sending out of.
*/
func SafeSendTezos(amount float64, toAddress string, alias string) error{
  strAmount := strconv.FormatFloat(amount, 'f', -1, 64)

  confirmStatement := "Send " + strAmount + " XTZ from " + alias + " to " + toAddress + "?"
  confirmation := askForConfirmation(confirmStatement)

  if confirmation{
    _, err := TezosDo("transfer", strAmount, "from", alias, "from", toAddress)
    if (err != nil){
      return errors.New("Could not send " + strAmount + " XTZ from " + alias + " to " + toAddress + ": tezosDo(args ...string) failed: " + err.Error())
    }
  } else {
    return errors.New("Cancelled: Send " + strAmount + " XTZ from " + alias + " to " + toAddress)
  }
  return nil
}


/*
Description: Will list the known addresses to your node and parse them into a multi-array.
Returns ([]KnownAddress): A structure containing the known address
*/
func ListKownAddresses() ([]KnownAddress, error){
  var knownAddresses []KnownAddress

  s, err := TezosDo("list", "known", "addresses")
  if (err != nil){
    return knownAddresses, errors.New("Could not list known addresses: tezosDo(args ...string) failed: " + err.Error())
  }

  parseKownAddresses := reListKownAddresses.FindAllStringSubmatch(s, -1)
  if (parseKownAddresses == nil){
    return knownAddresses, errors.New("Could not parse known addresses")
  }

  for _, address := range parseKownAddresses{
    knownAddresses = append(knownAddresses, KnownAddress{Address:address[1],Alias:address[0],Sk:address[2]})
  }

  return knownAddresses, nil
}

/*
Description: A function that executes a command to the tezos-client
Param args ([]string): Arguments to be executed
Returns (string): Returns the output of the executed command as a string
*/
func TezosDo(args ...string) (string, error){
  out, err := exec.Command(TezosPath, args...).Output()
  if err != nil {
    return "", err
  }

  return string(out[:]), nil
}

/*
Description: A function that executes an rpc get arg
Param args ([]string): Arguments to be executed
Returns (string): Returns the output of the executed command as a string
*/
func TezosRPCGet(arg string) (string, error){
  output, err := TezosDo("rpc", "get", arg)
  if (err != nil){
    return output, errors.New("Could not rpc get " + arg + " : tezosDo(args ...string) failed: " + err.Error())
  }
  return output, nil
}
