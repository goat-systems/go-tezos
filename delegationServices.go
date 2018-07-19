package goTezos
/*
Author: DefinitelyNotAGoat/MagicAglet
Version: 0.0.1
Description: A tool to be used for delegation services to calculate their payouts for each contract
License: MIT
*/

import (
  "encoding/json"
  "math/rand"
  "time"
  "regexp"
)

//Regex Used
var (
  reDelegatedContracts = regexp.MustCompile(`"([A-Z0-9a-z]+)"`)
)

/*
Description: A way to repesent each Delegated Contact, and their share for each cycle
Address: A string value representing the delegated contracts address
Commitments: An array of a structure that holds the amount commited for a cycle, and the percentage share
Delegator: Is this contract the delegator?
*/
type DelegatedClient struct {
    Address string //Public Key Hash
    Commitments []Commitment //Percentage of total delegation for profit share for each cycle participated
    Delegator bool //If this client is yourself or not.
    TotalPayout float64
}

/*
Description: A representation of the amount commited in a cycle, and the percentage share for that amount.
Cycle: The cycle number
Amount: XTZ value of the amount commited in the cycle
SharePercentage: The percentage value of the amount to all commitments made in that cycle
Payout: Amount of rewards to be paid out for the commitment
Timestamp: A timestamp to show when the commitment was made
*/
type Commitment struct {
  Cycle int
  Amount float64
  SharePercentage float64
  GrossPayout float64
  NetPayout float64
  Fee float64
}

/*
Description: Calculates the percentage share of a specific cycle for all delegated contracts on a range of cycles.
Param delegatedClients ([]DelegatedClient): A list of all the delegated contracts
Param cycleStart (int): The first cycle we are calculating
Param cycleEnd (int): The last cycle we are calculating
Returns delegatedClients ([]DelegatedClient): A list of all the delegated contracts
*/
func CalculateAllCommitmentsForCycles(delegatedClients []DelegatedClient, cycleStart int, cycleEnd int, rate float64) []DelegatedClient{
  for cycleStart <= cycleEnd {
    delegatedClients = CalculateAllCommitmentsForCycle(delegatedClients, cycleStart, rate)
    cycleStart = cycleStart + 1
  }
   return delegatedClients
}

/*
Description: Calculates the percentage share of a specific cycle for all delegated contracts
Param delegatedClients ([]DelegatedClient): A list of all the delegated contracts
Param cycle (int): The cycle we are calculating
Returns delegatedClients ([]DelegatedClient): A list of all the delegated contracts
*/
func CalculateAllCommitmentsForCycle(delegatedClients []DelegatedClient, cycle int, rate float64) []DelegatedClient{
  var sum float64
  sum = 0
  for i := 0; i < len(delegatedClients); i++{
    balance := GetBalanceAtSnapShotFor(delegatedClients[i].Address, cycle)
    sum = sum + balance
    delegatedClients[i].Commitments = append(delegatedClients[i].Commitments, Commitment{Cycle:cycle, Amount:balance})
  }

  for x := 0; x < len(delegatedClients); x++{
    counter := 0
    for y := 0; y < len(delegatedClients[x].Commitments); y++{
      if (delegatedClients[x].Commitments[y].Cycle == cycle){
        break
      }
      counter = counter + 1
    }
    delegatedClients[x].Commitments[counter].SharePercentage = delegatedClients[x].Commitments[counter].Amount / sum
    delegatedClients[x].Commitments[counter] = CalculatePayoutForCommitment(delegatedClients[x].Commitments[counter], rate)
  }
  return delegatedClients
}

/*
Description: Retrieves the list of addresses delegated to a delegate
Param SnapShot: A SnapShot object describing the desired snap shot.
Param delegateAddr: A string that represents a delegators tz address.
Returns []string: An array of contracts delegated to the delegator during the snap shot
*/
func GetDelegatedContractsForCycle(cycle int, delegateAddr string) []string{
  snapShot := GetSnapShot(cycle)
  hash := GetBlockLevelHash(snapShot.AssociatedBlock)
  getDelegatedContracts := "/chains/main/blocks/" + hash + "/context/delegates/" + delegateAddr + "/delegated_contracts"

  s := TezosRPCGet(getDelegatedContracts)

  DelegatedContracts := reDelegatedContracts.FindAllStringSubmatch(s, -1)

  return addressesToArray(DelegatedContracts)
}

/*
Description: Gets a list of all of the delegated contacts to a delegator
Param delegateAddr (string): string representation of the address of a delegator
Returns ([]string): An array of addresses (delegated contracts) that are delegated to the delegator
*/
func GetAllDelegatedContracts(delegateAddr string) []string{
  delegatedContractsCmd := "/chains/main/blocks/head/context/delegates/" + delegateAddr + "/delegated_contracts"
  s := TezosRPCGet(delegatedContractsCmd)

  DelegatedContracts := reDelegatedContracts.FindAllStringSubmatch(s, -1) //TODO Error checking

  return addressesToArray(DelegatedContracts)
}

/*
Description: Takes a commitment, and calculates the GrossPayout, NetPayout, and Fee.
Param commitment (Commitment): The commitment we are doing the operation on.
Param rate (float64): The delegation percentage fee written as decimal.
Param totalNodeRewards: Total rewards for the cyle the commitment represents. //TODO Make function to get total rewards for delegate in cycle
Returns (Commitment): Returns a commitment with the calculations made
Note: This function assumes Commitment.SharePercentage is already calculated.
*/
func CalculatePayoutForCommitment(commitment Commitment, rate float64) Commitment{
  ////-------------JUST FOR TESTING -------------////
  rand.Seed(time.Now().Unix())
  totalNodeRewards := rand.Intn(105000 - 70000) + 70000
 ////--------------END TESTING ------------------////

  grossRewards := commitment.SharePercentage * float64(totalNodeRewards)
  commitment.GrossPayout = grossRewards
  fee := rate * grossRewards
  netRewards := grossRewards - fee
  commitment.NetPayout = netRewards
  commitment.Fee = fee

  return commitment
}

/*
Description: A function to Payout rewards for all contracts in delegatedClients
Param delegatedClients ([]DelegatedClient): List of all contracts to be paid out
Param alias (string): The alias name to your known delegation wallet on your node


****WARNING****
If not using the ledger there is nothing stopping this from actually sending Tezos.
With the ledger you have to physically confirm the transaction, without the ledger you don't.

BE CAREFUL WHEN CALLING THIS FUNCTION!!!!!
****WARNING****
*/
func PayoutDelegatedContracts(delegatedClients []DelegatedClient, alias string){
  for _, delegatedClient := range delegatedClients {
    SendTezos(delegatedClient.TotalPayout, delegatedClient.Address, alias)
  }
}

/*
Description: Calculates the total payout in all commitments for a delegated contract
Param delegatedClients (DelegatedClient): the delegated contract to calulate over
Returns (DelegatedClient): return the contract with the Total Payout
*/
func CalculateTotalPayout(delegatedClient DelegatedClient) DelegatedClient{
  for _, commitment := range delegatedClient.Commitments{
    delegatedClient.TotalPayout = delegatedClient.TotalPayout + commitment.NetPayout
  }
  return delegatedClient
}

/*
Description: payout in all commitments for a delegated contract for all contracts
Param delegatedClients (DelegatedClient): the delegated contracts to calulate over
Returns (DelegatedClient): return the contract with the Total Payout for all contracts
*/
func CalculateAllTotalPayout(delegatedClients []DelegatedClient) []DelegatedClient{
  for index, delegatedClient := range delegatedClients{
    delegatedClients[index] = CalculateTotalPayout(delegatedClient)
  }

  return delegatedClients
}

/*
Description: Takes a multi-dimmensional array of addresses from a regex parse, and converts them into a single index(able) array.
Param matches ([][]string): All the addresses found and parsed by regex (ex. DelegatedContracts := reDelegatedContracts.FindAllStringSubmatch(s, -1) returns a multi dimmensional array)
Returns ([]string): Returns an index(able) string array of the matches input.
*/
func addressesToArray(matches [][]string) []string{
  var addresses []string
  for _, x := range matches {
    addresses = append(addresses, x[1])
  }

  return addresses
}

/*
Description: Takes an  array of interface (struct in our case), jsonifies it, and allows a much neater print.
Param v (interface{}): Array of an interface
*/
func PrettyReport(v interface{}) string {
  b, err := json.MarshalIndent(v, "", "  ")
  if err == nil {
    return string(b)
  }
  return ""
}

/*
Description: A test function that loops through the commitments of each delegated contract for a specific cycle,
             then it computes the share value of each one. The output should be = 1. With my tests it was, so you
             can really just ignore this.
Param cycle (int): The cycle number to be queryed
Param delegatedClients ([]DelegatedClient): the group of delegated DelegatedContracts
Returns (float64): The sum of all shares
*/
func CheckPercentageSumForCycle(cycle int, delegatedClients []DelegatedClient) float64{
  var sum float64
  sum = 0
  for x := 0; x < len(delegatedClients); x++{
    counter := 0
    for y := 0; y < len(delegatedClients[x].Commitments); y++{
      if (delegatedClients[x].Commitments[y].Cycle == cycle){
        break
      }
      counter = counter + 1
    }

    sum = sum + delegatedClients[x].Commitments[counter].SharePercentage
  }
  return sum
}

/*
Description: Reverse the order of an array of DelegatedClient.
             Used when fisrt retreiving contracts because the
             Tezos RPC API returns the newest contract first.
Param delegatedClients ([]DelegatedClient) Delegated

*/
func SortDelegateContracts(delegatedClients []DelegatedClient) []DelegatedClient{
   for i, j := 0, len(delegatedClients)-1; i < j; i, j = i+1, j-1 {
       delegatedClients[i], delegatedClients[j] = delegatedClients[j], delegatedClients[i]
   }
   return delegatedClients
}
