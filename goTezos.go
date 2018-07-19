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
  "regexp"
  "strconv"
)

var (
  tezosPath = ""
)

/*
Description: This library needs the TEZOSPATH enviroment variable to function
*/
func init() {
  tezosPath = os.Getenv("TEZOSPATH") + "tezos-client"
  if (tezosPath == ""){
    fmt.Println("Error: goTezos needs the enviroment variable TEZOSPATH. Please export it.")
    os.Exit(1)
  }
}

/*
Description: Gets the snapshot number for a certain cycle and returns the block level
Param cycle (int): Takes a cycle number as an integer to query for that cycles snapshot
Returns struct SnapShot: A SnapShot Structure defined above.
*/
func GetSnapShot(cycle int) SnapShot{
  var snapShot SnapShot
  snapShot.Cycle = cycle

  snapshotStr := "/chains/main/blocks/head/context/raw/json/cycle/" + strconv.Itoa(cycle)

  s := TezosRPCGet(snapshotStr)

  regRandomSeed := reGetRandomSeed.FindStringSubmatch(s)

  if (len(regRandomSeed) <1){
    snapShot.Decided = false
  } else {
    regRollSnapShot := reGetRollSnapShot.FindStringSubmatch(s)
    number, _ := strconv.Atoi(regRollSnapShot[1])
    snapShot.Number = number
    snapShot.Decided = true
    snapShot.AssociatedBlock =((cycle - 7) * 4096) + (number + 1) * 256
  }

  return snapShot
}

/*
Description: Will retreive the current block level as an integer
Returns (int): Returns integer representation of block level
*/
func GetBlockLevelHead() int{
  s := TezosRPCGet("chains/main/blocks/head")

  regHeadLevelResult := reGetBlockLevelHead.FindStringSubmatch(s)
  headlevel, _ := strconv.Atoi(regHeadLevelResult[1]) //TODO Error Checking

  return headlevel
}

/*
Description: Takes a block level, and returns the hash for that specific level
Param level (int): An integer representation of the block level to query
Returns (string): A string representation of the hash for the block level queried.
*/
func GetBlockLevelHash(level int) string{
  diff := GetBlockLevelHead() - level
  diffStr := strconv.Itoa(diff)
  getBlockByLevel := "chains/main/blocks/head~" + diffStr

  s := TezosRPCGet(getBlockByLevel)

  hash := reGetHash.FindStringSubmatch(s) //TODO Error check the regex

  return hash[1]
}

/*
// TODO Need to finish and test
Description: Returns the balance to a specific tezos address
Param tezosAddr (string): Takes a string representation of the address querying
Returns (float64): Returns a float64 representation of the balance for the account
*/
func GetBalanceFor(tezosAddr string) float64{

  s := TezosDo("get", "balance", "for", tezosAddr)

  regGetBalance := reGetBalance.FindStringSubmatch(s) //TODO Regex error checking
  floatBalance, _ := strconv.ParseFloat(regGetBalance[1], 64) //TODO error checking

  return floatBalance
}

/*
Description: Will get the balance of an account at a specific snapshot
Param tezosAddr (string): Takes a string representation of the address querying
Param cycle (int): The cycle we are getting the snapshot for
Returns (float64): Returns a float64 representation of the balance for the account
*/
func GetBalanceAtSnapShotFor(tezosAddr string, cycle int) float64{
  snapShot := GetSnapShot(cycle)
  hash := GetBlockLevelHash(snapShot.AssociatedBlock)

  balanceCmdStr := "chains/main/blocks/" + hash + "/context/contracts/" + tezosAddr + "/balance"

  s := TezosRPCGet(balanceCmdStr)

  regGetBalance := reGetBalance.FindStringSubmatch(s) //TODO Regex error checking

  var returnBalance float64

  if (len(regGetBalance) < 1){
    returnBalance = 0
  } else{
    floatBalance, _ := strconv.ParseFloat(regGetBalance[1], 64) //TODO error checking
    returnBalance = floatBalance
  }

  return returnBalance / 1000000
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
func SendTezos(amount float64, toAddress string, alias string){
  strAmount := strconv.FormatFloat(amount, 'f', -1, 64)
  TezosDo("transfer", strAmount, "from", alias, "from", toAddress)

}

/*
Description: Sends tezos from your specified wallet to another account, but makes you confirm the transaction.
Param amount (float64): Amount of tezos to be sent
Param toAddress (string): The address you are sending tezos to.
Param alias (string): The named alias assigned to your wallet you are sending out of.
*/
func SafeSendTezos(amount float64, toAddress string, alias string){
  strAmount := strconv.FormatFloat(amount, 'f', -1, 64)

  confirmStatement := "Send " + strAmount + " XTZ from " + alias + " to " + toAddress + "?"
  confirmation := askForConfirmation(confirmation)

  if confirmation{
    TezosDo("transfer", strAmount, "from", alias, "from", toAddress)
  } else {
    fmt.Println("Cancelled: Send " + strAmount + " XTZ from " + alias + " to " + toAddress)
  }
}


/*
Description: Will list the known addresses to your node and parse them into a multi-array.
Returns ([]KnownAddress): A structure containing the known address
*/
func ListKownAddresses() []KnownAddress{
  s := tezosDo("list", "known", "addresses")
  parseKownAddresses := reListKownAddresses.FindAllStringSubmatch(s, -1)

  var knownAddresses []KnownAddress
  for _, address := range parseKownAddresses{
    knownAddresses = append(knownAddresses, KnownAddress{Address:address[1],Alias:address[0],Sk:address[2]})
  }

  return knownAddresses
}

/*
Description: A function that executes a command to the tezos-client
Param args ([]string): Arguments to be executed
Returns (string): Returns the output of the executed command as a string
*/
func TezosDo(args ...string) string{
  out, err := exec.Command(tezosPath, args...).Output()
	if err != nil {
		fmt.Println(err)
	}

  return string(out[:])
}

/*
Description: A function that executes an rpc get arg
Param args ([]string): Arguments to be executed
Returns (string): Returns the output of the executed command as a string
*/
func TezosRPCGet(arg string) string{
  return TezosDo("rpc", "get", arg)
}
