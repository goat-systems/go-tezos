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

var (
  tezosPath = ""
)

/*
Description: This library needs the TEZOSPATH enviroment variable to function
*/
func init() {
  tezosPath, ok := os.LookupEnv("TEZOSPATH")
  if !ok {
	   fmt.Println("Error: Could not retrieve TEZOSPATH")
	   os.Exit(1)
  }
  tezosPath += "tezos-client"
  fmt.Println(tezosPath)
}

/*
Description: Gets the snapshot number for a certain cycle and returns the block level
Param cycle (int): Takes a cycle number as an integer to query for that cycles snapshot
Returns struct SnapShot: A SnapShot Structure defined above.
*/
func GetSnapShot(cycle int) (SnapShot, error){
  var snapShot SnapShot
  snapShot.Cycle = cycle
  strCycle := strconv.Itoa(cycle)

  snapshotStr := "/chains/main/blocks/head/context/raw/json/cycle/" + strCycle

  s, err := TezosRPCGet(snapshotStr)
  if (err != nil){
    return snapShot, errors.New("Could not get snapshot for cycle " + strCycle + ": TezosRPCGet(arg string) failed: " + err.Error())
  }

  regRandomSeed := reGetRandomSeed.FindStringSubmatch(s)
  if (regRandomSeed == nil){
    return snapShot, errors.New("No random seed, could not get snapshot for cycle " + strCycle)
  }


  regRollSnapShot := reGetRollSnapShot.FindStringSubmatch(s)
  if (regRollSnapShot == nil){
    return snapShot, errors.New("Could not get snapshot for cycle " + strCycle)
  }
  number, _ := strconv.Atoi(regRollSnapShot[1])
  snapShot.Number = number
  snapShot.AssociatedBlock =((cycle - 7) * 4096) + (number + 1) * 256

  return snapShot, nil
}

/*
Description: Will retreive the current block level as an integer
Returns (int): Returns integer representation of block level
*/
func GetBlockLevelHead() (int, error){
  s, err := TezosRPCGet("chains/main/blocks/head")
  if (err != nil){
    return 0, errors.New("Could not get block level for head: TezosRPCGet(arg string) failed: " + err.Error())
  }


  regHeadLevelResult := reGetBlockLevelHead.FindStringSubmatch(s)
  if (regHeadLevelResult == nil){
    return 0, errors.New("Could not get block level for head")
  }
  headlevel, _ := strconv.Atoi(regHeadLevelResult[1]) //TODO Error Checking

  return headlevel, nil
}

/*
Description: Takes a block level, and returns the hash for that specific level
Param level (int): An integer representation of the block level to query
Returns (string): A string representation of the hash for the block level queried.
*/
func GetBlockLevelHash(level int) (string, error){
  head, err := GetBlockLevelHead()
  if (err != nil){
    return "", errors.New("Could not get hash for block " +  strconv.Itoa(level) + ": GetBlockLevelHead() failed: " + err.Error())
  }
  diff :=  head - level

  diffStr := strconv.Itoa(diff)
  getBlockByLevel := "chains/main/blocks/head~" + diffStr

  s, err := TezosRPCGet(getBlockByLevel)
  if (err != nil){
    return "", errors.New("Could not get hash for block " +  strconv.Itoa(level) + ": TezosRPCGet(arg string) failed: " + err.Error())
  }

  hash := reGetHash.FindStringSubmatch(s) //TODO Error check the regex
  if (hash == nil){
    return "", errors.New("Could not get hash for block " + strconv.Itoa(level))
  }

  return hash[1], nil
}

/*
Description: Returns the balance to a specific tezos address
Param tezosAddr (string): Takes a string representation of the address querying
Returns (float64): Returns a float64 representation of the balance for the account
*/
func GetBalanceFor(tezosAddr string) (float64, error){

  s, err := TezosDo("get", "balance", "for", tezosAddr)
  if (err != nil){
    return 0, errors.New("Could not get balance for " + tezosAddr + ": tezosDo(args ...string) failed: " + err.Error())
  }

  regGetBalance := reGetBalance.FindStringSubmatch(s) //TODO Regex error checking
  if (regGetBalance == nil){
    return 0, errors.New("Could not get balance for " + tezosAddr)
  }
  floatBalance, _ := strconv.ParseFloat(regGetBalance[1], 64) //TODO error checking

  return floatBalance, nil
}

/*
Description: Will get the balance of an account at a specific snapshot
Param tezosAddr (string): Takes a string representation of the address querying
Param cycle (int): The cycle we are getting the snapshot for
Returns (float64): Returns a float64 representation of the balance for the account
*/
func GetBalanceAtSnapShotFor(tezosAddr string, cycle int) (float64, error){
  snapShot, err := GetSnapShot(cycle)
  if (err != nil){
    return 0, errors.New("Could not get balance at snapshot for " +  tezosAddr + ": GetSnapShot(cycle int) failed: " + err.Error())
  }

  hash, err := GetBlockLevelHash(snapShot.AssociatedBlock)
  if (err != nil){
    return 0, errors.New("Could not get hash for block " +  strconv.Itoa(snapShot.AssociatedBlock) + ": GetBlockLevelHead() failed: " + err.Error())
  }

  balanceCmdStr := "chains/main/blocks/" + hash + "/context/contracts/" + tezosAddr + "/balance"

  s, err := TezosRPCGet(balanceCmdStr)
  if (err != nil){
    return 0, errors.New("Could not get balance at snapshot for " +  tezosAddr + ": TezosRPCGet(arg string) failed: " + err.Error())
  }

  regGetBalance := reGetBalance.FindStringSubmatch(s)
  if (regGetBalance == nil){
    return 0, errors.New("Could not parse balance for " + s)
  }

  var returnBalance float64

  if (len(regGetBalance) < 1){
    returnBalance = 0
  } else{
    floatBalance, _ := strconv.ParseFloat(regGetBalance[1], 64) //TODO error checking
    returnBalance = floatBalance
  }

  return returnBalance / 1000000, nil
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
  _, err := TezosDo("transfer", strAmount, "from", alias, "from", toAddress)
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
  fmt.Println(args)
  out, err := exec.Command(tezosPath, args...).Output()
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
