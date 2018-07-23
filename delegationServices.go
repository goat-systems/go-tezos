package goTezos
/*
Author: DefinitelyNotAGoat/MagicAglet
Version: 0.0.1
Description: This file contains specific functions for delegation services
License: MIT
*/

import (
  //"math/rand"
  //"time"
  "strconv"
  "errors"
  "math"
  "fmt"
)

/*
Description: Calculates the percentage share of a specific cycle for all delegated contracts on a range of cycles.
Param delegatedContracts ([]DelegatedClient): A list of all the delegated contracts
Param cycleStart (int): The first cycle we are calculating
Param cycleEnd (int): The last cycle we are calculating
Returns delegatedContracts ([]DelegatedContract): A list of all the delegated contracts
*/
func CalculateAllContractsForCycles(delegatedContracts []DelegatedContract, cycleStart int, cycleEnd int, rate float64, spillage bool, delegateAddr string) ([]DelegatedContract, error){
  var err error

  for cycleStart <= cycleEnd {
    fmt.Println(cycleStart)
    delegatedContracts, err = CalculateAllContractsForCycle(delegatedContracts, cycleStart, rate, spillage, delegateAddr)
    if (err != nil){
      return delegatedContracts, errors.New("Could not calculate all commitments for cycles " + strconv.Itoa(cycleStart) + "-" +  strconv.Itoa(cycleEnd) + ":CalculateAllCommitmentsForCycle(delegatedContracts []DelegatedContract, cycle int, rate float64) failed: " + err.Error())
    }
    cycleStart = cycleStart + 1
  }
   return delegatedContracts, nil
}


/*
Description: Calculates the percentage share of a specific cycle for all delegated contracts
Param delegatedContracts ([]DelegatedContract): A list of all the delegated contracts
Param cycle (int): The cycle we are calculating
Param rate (float64): Fee rate of the delegate
Param spillage (bool): If the delegate wants to hard cap the payouts at complete rolls
Returns delegatedContracts ([]DelegatedContract): A list of all the delegated contracts
*/
func CalculateAllContractsForCycle(delegatedContracts []DelegatedContract, cycle int, rate float64, spillage bool, delegateAddr string) ([]DelegatedContract, error) {
  var err error
  var balance float64
  delegationsForCycle, _ := GetDelegatedContractsForCycle(cycle, delegateAddr)

  for index, delegation := range delegatedContracts{
    balance, err = GetAccountBalanceAtSnapshot(delegation.Address, cycle)
    if (err != nil){
      return delegatedContracts, errors.New("Could not calculate all commitments for cycle " + strconv.Itoa(cycle) + ":GetAccountBalanceAtSnapshot(tezosAddr string, cycle int) failed: " + err.Error())
    }
    if (isDelegationInGroup(delegatedContracts[index].Address, delegationsForCycle) && delegatedContracts[index].Address != delegateAddr){
      delegatedContracts[index].Contracts = append(delegatedContracts[index].Contracts, Contract{Cycle:cycle, Amount:balance})
    } else{
      delegatedContracts[index].Contracts = append(delegatedContracts[index].Contracts, Contract{Cycle:cycle, Amount:0})
    }
    fmt.Println(delegatedContracts[index].Contracts)
  }

  delegatedContracts, err = CalculatePercentageSharesForCycle(delegatedContracts, cycle, rate, spillage, delegateAddr)
  if (err != nil){
    return delegatedContracts, errors.New("func CalculateAllContractsForCycle(delegatedContracts []DelegatedContract, cycle int, rate float64, spillage bool, delegateAddr string) failed: " + err.Error())
  }
  return delegatedContracts, nil
}

func isDelegationInGroup(phk string, group []string) bool{
  for _, address := range group{
    if (address == phk){
      return true
    }
  }
  return false
}

/*
Description: Calculates the share percentage of a cycle over a list of delegated contracts


*/
func CalculatePercentageSharesForCycle(delegatedContracts []DelegatedContract, cycle int, rate float64, spillage bool, delegateAddr string) ([]DelegatedContract, error){
  var stakingBalance float64
  //var balance float64
  var err error

  spillAlert := false

  stakingBalance, err = GetDelegateStakingBalance(delegateAddr, cycle)
  if (err != nil){
    return delegatedContracts, errors.New("func CalculateRollSpillage(delegatedContracts []DelegatedContract, delegateAddr string) failed: " + err.Error())
  }

  mod := math.Mod(stakingBalance, 10000)
  sum := stakingBalance - mod
  balanceCheck := stakingBalance - mod

  for index, delegation := range  delegatedContracts{
    counter := 0
    for i, _ := range delegation.Contracts {
      if (delegatedContracts[index].Contracts[i].Cycle == cycle){
        break
      }
      counter = counter + 1
    }
    balanceCheck = balanceCheck - delegatedContracts[index].Contracts[counter].Amount
    //fmt.Println(stakingBalance)
    if (spillAlert){
      delegatedContracts[index].Contracts[counter].SharePercentage = 0
      delegatedContracts[index].Contracts[counter].RollInclusion = 0
    } else if (balanceCheck < 0 && spillage){
      spillAlert = true
      delegatedContracts[index].Contracts[counter].SharePercentage = (delegatedContracts[index].Contracts[counter].Amount + stakingBalance) / sum
      delegatedContracts[index].Contracts[counter].RollInclusion = delegatedContracts[index].Contracts[counter].Amount + stakingBalance
    } else{
      delegatedContracts[index].Contracts[counter].SharePercentage = delegatedContracts[index].Contracts[counter].Amount / stakingBalance
      delegatedContracts[index].Contracts[counter].RollInclusion = delegatedContracts[index].Contracts[counter].Amount
    }
    delegatedContracts[index].Contracts[counter] = CalculatePayoutForContract(delegatedContracts[index].Contracts[counter], rate, delegatedContracts[index].Delegate)
    delegatedContracts[index].Fee = delegatedContracts[index].Fee + delegatedContracts[index].Contracts[counter].Fee
  }

  return delegatedContracts, nil
}

/*
Description: Retrieves the list of addresses delegated to a delegate
Param SnapShot: A SnapShot object describing the desired snap shot.
Param delegateAddr: A string that represents a delegators tz address.
Returns []string: An array of contracts delegated to the delegator during the snap shot
*/
func GetDelegatedContractsForCycle(cycle int, delegateAddr string) ([]string, error){
  var rtnString []string
  snapShot, err := GetSnapShot(cycle)
  // fmt.Println(snapShot)
  if (err != nil){
    return rtnString, errors.New("Could not get delegated contracts for cycle " + strconv.Itoa(cycle) + ": GetSnapShot(cycle int) failed: " + err.Error())
  }
  hash, err:= GetBlockLevelHash(snapShot.AssociatedBlock)
  if (err != nil){
    return rtnString, errors.New("Could not get delegated contracts for cycle " + strconv.Itoa(cycle) + ": GetBlockLevelHash(level int) failed: " + err.Error())
  }
  // fmt.Println(hash)
  getDelegatedContracts := "/chains/main/blocks/" + hash + "/context/delegates/" + delegateAddr + "/delegated_contracts"

  s, err := TezosRPCGet(getDelegatedContracts)
  if (err != nil){
    return rtnString, errors.New("Could not get delegated contracts for cycle " + strconv.Itoa(cycle) + ": TezosRPCGet(arg string) failed: " + err.Error())
  }

  DelegatedContracts := reDelegatedContracts.FindAllStringSubmatch(s, -1)
  if (DelegatedContracts == nil){
    return rtnString, errors.New("Could not get delegated contracts for cycle " + strconv.Itoa(cycle) + ": You have no contracts.")
  }
  rtnString = addressesToArray(DelegatedContracts)
  return rtnString, nil
}

/*
Description: Gets a list of all of the delegated contacts to a delegator
Param delegateAddr (string): string representation of the address of a delegator
Returns ([]string): An array of addresses (delegated contracts) that are delegated to the delegator
*/
func GetAllDelegatedContracts(delegateAddr string) ([]string, error){
  var rtnString []string
  delegatedContractsCmd := "/chains/main/blocks/head/context/delegates/" + delegateAddr + "/delegated_contracts"
  s, err := TezosRPCGet(delegatedContractsCmd)
  if (err != nil){
    return rtnString, errors.New("Could not get delegated contracts: TezosRPCGet(arg string) failed: " + err.Error())
  }

  DelegatedContracts := reDelegatedContracts.FindAllStringSubmatch(s, -1) //TODO Error checking
  if (DelegatedContracts == nil){
    return rtnString, errors.New("Could not get all delegated contracts: Regex failed")
  }
  rtnString = addressesToArray(DelegatedContracts)
  return rtnString, nil
}

/*
Description: Takes a commitment, and calculates the GrossPayout, NetPayout, and Fee.
Param commitment (Commitment): The commitment we are doing the operation on.
Param rate (float64): The delegation percentage fee written as decimal.
Param totalNodeRewards: Total rewards for the cyle the commitment represents. //TODO Make function to get total rewards for delegate in cycle
Param delegate (bool): Is this the delegate
Returns (Commitment): Returns a commitment with the calculations made
Note: This function assumes Commitment.SharePercentage is already calculated.
*/
func CalculatePayoutForContract(contract Contract, rate float64, delegate bool) Contract{
  ////-------------JUST FOR TESTING -------------////
  totalNodeRewards := 378 //Amount of rewards for my delegation in cycle 11
 ////--------------END TESTING ------------------////

  grossRewards := contract.SharePercentage * float64(totalNodeRewards)
  contract.GrossPayout = grossRewards
  fee := rate * grossRewards
  contract.Fee = fee
  var netRewards float64
  if (delegate){
    netRewards = grossRewards
    contract.NetPayout = netRewards
    contract.Fee = 0
  } else {
    netRewards = grossRewards - fee
    contract.NetPayout = contract.NetPayout + netRewards
  }

  return contract
}

/*
Description: Calculates all the fees for each contract and adds them to the delegates net payout


*/
func CalculateDelegateNetPayout(delegatedContracts []DelegatedContract) []DelegatedContract{
  var delegateIndex int

  for index, delegate := range delegatedContracts{
    if (delegate.Delegate){
      delegateIndex = index
    }
  }

  for _, delegate := range delegatedContracts{
    if (!delegate.Delegate){
      delegatedContracts[delegateIndex].TotalPayout = delegatedContracts[delegateIndex].TotalPayout + delegate.Fee
    }
  }
  return delegatedContracts
}

/*
Description: A function to Payout rewards for all contracts in delegatedContracts
Param delegatedContracts ([]DelegatedClient): List of all contracts to be paid out
Param alias (string): The alias name to your known delegation wallet on your node
****WARNING****
If not using the ledger there is nothing stopping this from actually sending Tezos.
With the ledger you have to physically confirm the transaction, without the ledger you don't.
BE CAREFUL WHEN CALLING THIS FUNCTION!!!!!
****WARNING****
*/
func PayoutDelegatedContracts(delegatedContracts []DelegatedContract, alias string) error{
  for _, delegatedContract := range delegatedContracts {
    err := SendTezos(delegatedContract.TotalPayout, delegatedContract.Address, alias)
    if (err != nil){
      return errors.New("Could not Payout Delegated Contracts: SendTezos(amount float64, toAddress string, alias string) failed: " + err.Error())
    }
  }
  return nil
}

/*
Description: Calculates the total payout in all commitments for a delegated contract
Param delegatedContracts (DelegatedClient): the delegated contract to calulate over
Returns (DelegatedClient): return the contract with the Total Payout
*/
func CalculateTotalPayout(delegatedContract DelegatedContract) DelegatedContract{
  for _, contract := range delegatedContract.Contracts{
    delegatedContract.TotalPayout = delegatedContract.TotalPayout + contract.NetPayout
  }
  return delegatedContract
}

/*
Description: payout in all commitments for a delegated contract for all contracts
Param delegatedContracts (DelegatedClient): the delegated contracts to calulate over
Returns (DelegatedClient): return the contract with the Total Payout for all contracts
*/
func CalculateAllTotalPayout(delegatedContracts []DelegatedContract) []DelegatedContract{
  for index, delegatedContract := range delegatedContracts{
    delegatedContracts[index] = CalculateTotalPayout(delegatedContract)
  }

  return delegatedContracts
}

/*
Description: A test function that loops through the commitments of each delegated contract for a specific cycle,
             then it computes the share value of each one. The output should be = 1. With my tests it was, so you
             can really just ignore this.
Param cycle (int): The cycle number to be queryed
Param delegatedContracts ([]DelegatedClient): the group of delegated DelegatedContracts
Returns (float64): The sum of all shares
*/
func CheckPercentageSumForCycle(cycle int, delegatedContracts []DelegatedContract) float64{
  var sum float64
  sum = 0
  for x := 0; x < len(delegatedContracts); x++{
    counter := 0
    for y := 0; y < len(delegatedContracts[x].Contracts); y++{
      if (delegatedContracts[x].Contracts[y].Cycle == cycle){
        break
      }
      counter = counter + 1
    }

    sum = sum + delegatedContracts[x].Contracts[counter].SharePercentage
  }
  return sum
}

/*
Description: A function to account for incomplete rolls, and the payouts associated with that
TODO: In Progress
*/
func CalculateRollSpillage(delegatedContracts []DelegatedContract, delegateAddr string, cycle int) ([]DelegatedContract, error) {
  stakingBalance, err := GetDelegateStakingBalance(delegateAddr, cycle)
  if (err != nil){
    return delegatedContracts, errors.New("func CalculateRollSpillage(delegatedContracts []DelegatedContract, delegateAddr string) failed: " + err.Error())
  }

  mod := math.Mod(stakingBalance, 10000)
  sum := mod * 10000

  for index, delegatedContract := range delegatedContracts{
    for i, contract := range delegatedContract.Contracts{
      if (contract.Cycle == cycle){
        stakingBalance = stakingBalance - contract.Amount
        if (stakingBalance < 0){
          delegatedContracts[index].Contracts[i].SharePercentage = (contract.Amount - stakingBalance) / sum
        }
      }
    }
  }

  return delegatedContracts, nil
}

/*
Description: Reverse the order of an array of DelegatedClient.
             Used when fisrt retreiving contracts because the
             Tezos RPC API returns the newest contract first.
Param delegatedContracts ([]DelegatedClient) Delegated
*/
func SortDelegateContracts(delegatedContracts []DelegatedContract) []DelegatedContract{
   for i, j := 0, len(delegatedContracts)-1; i < j; i, j = i+1, j-1 {
       delegatedContracts[i], delegatedContracts[j] = delegatedContracts[j], delegatedContracts[i]
   }
   return delegatedContracts
}
