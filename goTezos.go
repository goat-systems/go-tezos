package goTezos

/*
Author: DefinitelyNotAGoat/MagicAglet
Version: 0.0.1
Description: The Tezos API written in GO, for easy development.
License: MIT
*/

import (
  "strconv"
  "errors"
)

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
func GetAccountBalanceAtSnapshot(tezosAddr string, cycle int) (float64, error){
  snapShot, err := GetSnapShot(cycle)
  if (err != nil){
    return 0, errors.New("Could not get balance at snapshot for " +  tezosAddr + ": GetSnapShot(cycle int) failed: " + err.Error())
  }

  hash, err := GetBlockLevelHash(snapShot.AssociatedBlock)
  if (err != nil){
    return 0, errors.New("Could not get hash for block " +  strconv.Itoa(snapShot.AssociatedBlock) + ": GetBlockLevelHead() failed: " + err.Error())
  }

  balanceCmdStr := "/chains/main/blocks/" + hash + "/context/contracts/" + tezosAddr + "/balance"

  s, err := TezosRPCGet(balanceCmdStr)
  if (err != nil){
    return 0, errors.New("Could not get balance at snapshot for " +  tezosAddr + ": TezosRPCGet(arg string) failed: " + err.Error())
  }

  var returnBalance float64
  var regGetBalance []string

  if (s == "No service found at this URL"){
    returnBalance = 0
  } else{
    regGetBalance = reGetBalance.FindStringSubmatch(s)
    if (regGetBalance == nil){
      return 0, errors.New("Could not parse balance for " + s)
    }
  }

  if (len(regGetBalance) < 1){
    returnBalance = 0
  } else{
    floatBalance, _ := strconv.ParseFloat(regGetBalance[1], 64) //TODO error checking
    returnBalance = floatBalance
  }

  return returnBalance / 1000000, nil
}

/*
Description: Will get the staking balance of a delegate
Param delegateAddr (string): Takes a string representation of the address querying
Returns (float64): Returns a float64 representation of the balance for the account
*/
func GetDelegateStakingBalance(delegateAddr string, cycle int) (float64, error){
  var snapShot SnapShot
  var err error
  var hash string
  var s string

  snapShot, err = GetSnapShot(cycle)
  if (err != nil){
    return 0, errors.New("GetDelegateStakingBalance(delegateAddr string, cycle int) failed: " + err.Error())
  }

  hash, err = GetBlockLevelHash(snapShot.AssociatedBlock)
  if (err != nil){
    return 0, errors.New("GetDelegateStakingBalance(delegateAddr string, cycle int) failed: " + err.Error())
  }

  rpcCall := "/chains/main/blocks/" + hash + "/context/delegates/" + delegateAddr + "/staking_balance"

  s, err = TezosRPCGet(rpcCall)
  if (err != nil){
    return 0, errors.New("GetDelegateStakingBalance(delegateAddr string) failed: " + err.Error())
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
