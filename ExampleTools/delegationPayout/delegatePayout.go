package main
/*
Author: DefinitelyNotAGoat/MagicAglet
Version: 0.0.1
Description: A tool to be used for delegation services to calculate their payouts for each contract
License: MIT
*/

import (
  "fmt"
  "flag"
  "bufio"
  "os"
  "regexp"
  "strconv"
  "github.com/DefinitelyNotAGoat/goTezos"
  "gopkg.in/cheggaaa/pb.v1"
)

var (
    bar = pb.StartNew(5).Prefix("Taco Mission") //A fun loading bar to let you know the progress
)

func main() {
  var cycleRange [2]int

  delegateAddr := flag.String("delegateaddr", "nil", "The address to your delegation service (eg. tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc)") //The tz1 address of the delegation service the program is used for.
  cycle := flag.Int("cycle", -1, "The cycle you are querying for.") //If only querying one cycle, use this arg.
  cycles := flag.String("cycles", "nil", "A range of cycles to compute delegated contracts. Example: 10-20") //If querying a range of cycles use this arg
  fee := flag.Float64("fee", .05, "Your dynamic fee percentage noted as a decimal.") //The fee your service charges.
  //alias := flag.String("alias", "nil", "The alias to your baking address on your node. Needed for to send XTZ for payouts.") //The alias to your delegate wallet, used for sending out payments
  report := flag.Bool("report", true, "Generates a list of all the Delegated Contracts for the request cycle(s).") //Generate a report of a cycle or cycles
  payout := flag.Bool("payout", false, "Pays each of your delegated contracts their share less your percentage fee.")//Pay your contracts
  flag.Parse()


  var delegatedContracts []goTezos.DelegatedContract //Our delegated contracts in a cycle or cycles

  _, err := goTezos.GetBalanceFor(*delegateAddr) //A dirty trick to check if delegate address is real
  if (err != nil){
    fmt.Println("Invalid Delegator Address " + *delegateAddr)
    fmt.Println("func main() failed: " + err.Error())
    os.Exit(1)
  }


  if (*cycle != -1){
    delegatedContracts = singleCycleOp(*cycle, *delegateAddr, *fee) //perform operations over a single cycle
  } else if (*cycles != "nil"){
    cycleRange = parseCyclesInput(*cycles)
    delegatedContracts = multiCycleOp(cycleRange[0], cycleRange[1],*delegateAddr, *fee) //perform operations over multiple cycles
  } else{
    fmt.Println("No cycle(s) provided. Exiting...")
    os.Exit(1)
  }

  if (*report == true){ //If the program was ran to get a report only
    generateReport(delegatedContracts)
    bar.Increment()
  } else if (*report == false && *payout == true){ //If the program was ran to payout only
    //goTezos.PayoutDelegatedContracts(delegatedClients, *alias)
    bar.Increment()
    fmt.Println("PayoutDelegatedContracts Function is temporary disabled for safety!")
  } else if (*report == true && *payout == true){ //If the program was ran to generate a report and payout
    generateReport(delegatedContracts)
    //goTezos.PayoutDelegatedContracts(delegatedClients, *alias)
    bar.Increment()
    fmt.Println("PayoutDelegatedContracts Function is temporary disabled for safety!")
  }
  bar.FinishPrint("Tacos Have Been Made On This Day!")
}

/*
Description: This function will take a cycle, get all delegated contracts for a delegate in that cycle.
             Calculate the percentage share of that cycle for each contract, calculate the fees for each contract,
             and Return the all contracts with the information above.
Param cycle (int): Cycle to query for.
Param delegateAddr (string): The delegate address we are querying
Param fee (float64): The fee for the delegate
Returns ([]DelegatedClient): A list of all delegated contracts and the needed info
*/
func singleCycleOp(cycle int, delegateAddr string, fee float64) []goTezos.DelegatedContract{
  var delegatedClients []goTezos.DelegatedContract
  contracts, err := goTezos.GetDelegatedContractsForCycle(cycle, delegateAddr)
  if (err != nil){
    fmt.Println("func singleCycleOp(cycle, delegateAddr, fee) failed: " + err.Error())
    os.Exit(-1)
  }
  for _, delegatedClientAddr := range contracts {
    delegatedClients = append(delegatedClients, goTezos.DelegatedContract{Address:delegatedClientAddr, Delegate:false, TotalPayout:0})
  }
  delegatedClients = append(delegatedClients, goTezos.DelegatedContract{Address:delegateAddr, Delegate:true, TotalPayout:0}) //Need to keep track of your own baking rewards, to avoid accidentally including them in your fee system.
  bar.Increment()
  delegatedClients = goTezos.SortDelegateContracts(delegatedClients) //Put the oldest contract at the begining of the array
  bar.Increment()
  delegatedClients, err = goTezos.CalculateAllContractsForCycle(delegatedClients, cycle, fee, false, delegateAddr)
  delegatedClients = goTezos.CalculateDelegateNetPayout(delegatedClients)
  if (err != nil){
    fmt.Println("func singleCycleOp(cycle, delegateAddr, fee) failed: " + err.Error())
    os.Exit(-1)
  }
  bar.Increment()
  delegatedClients = goTezos.CalculateAllTotalPayout(delegatedClients)
  bar.Increment()

  return delegatedClients
}

/*
Description: This function will take a range of cycles, get all delegated contracts for a delegate in those cycles.
             Calculate the percentage share of the cycle(s) for each contract, calculate the fees for each contract,
             and Return the all contracts with the information above.
Param cycleStart (int): Start of cycle range.
Param cycleEnd (int): End of cycle range.
Param delegateAddr (string): The delegate address we are querying
Param fee (float64): The fee for the delegate
Returns ([]DelegatedClient): A list of all delegated contracts and the needed info
*/
func multiCycleOp(cycleStart int, cycleEnd int, delegateAddr string, fee float64) []goTezos.DelegatedContract{
  var delegatedClients []goTezos.DelegatedContract
  contracts, err := goTezos.GetAllDelegatedContracts(delegateAddr)
  if (err != nil){
    fmt.Println("func multiCycleOp(cycleStart, cycleEnd, delegateAddr, fee) failed: " + err.Error())
    os.Exit(-1)
  }
  for _, delegatedClientAddr := range contracts {
    delegatedClients = append(delegatedClients, goTezos.DelegatedContract{Address:delegatedClientAddr, Delegate:false, TotalPayout:0})
  }
  delegatedClients = append(delegatedClients, goTezos.DelegatedContract{Address:delegateAddr, Delegate:true, TotalPayout:0}) //Need to keep track of your own baking rewards, to avoid accidentally including them in your fee system.
  bar.Increment()
  delegatedClients = goTezos.SortDelegateContracts(delegatedClients) //Put the oldest contract at the begining of the array
  bar.Increment()
  delegatedClients, err = goTezos.CalculateAllContractsForCycles(delegatedClients, cycleStart, cycleEnd, fee, false, delegateAddr)
  delegatedClients = goTezos.CalculateDelegateNetPayout(delegatedClients)
  if (err != nil){
    fmt.Println("func multiCycleOp(cycleStart, cycleEnd, delegateAddr, fee) failed: " + err.Error())
    os.Exit(-1)
  }
  bar.Increment()
  delegatedClients = goTezos.CalculateAllTotalPayout(delegatedClients)
  bar.Increment()

  return delegatedClients
}

/*
Description: Generates a JSON report of all delegated contracts in a cycle(s), and writes them to report.json
Param delegatedClients ([]DelegatedClient): List of delegated contracts
*/
func generateReport(delegatedClients []goTezos.DelegatedContract){
  report := goTezos.PrettyReport(delegatedClients)
  write(report)
}

/*
Description: Basic File IO specifically for writing to a file
Param report (string): a string to be written to the file
*/
func write(report string){
  f, err := os.Create("./report.json")
  if err != nil {
    panic("func write(report) failed: " + err.Error())
  }

  defer f.Close()

  writer := bufio.NewWriter(f)
  writer.WriteString(report)
  writer.Flush()
}

/*
Description: Parses the command line option for cycles into two integers to be used as range
Param cycles (string): the string command line argument *cycles
Returns ([]int): An array of size 2, index 0 = first cycle in range, index 1 = last cycle in range
*/
func parseCyclesInput(cycles string) [2]int{
  reCycles := regexp.MustCompile(`([0-9]+)-([0-9]+)`)
  arrayCycles := reCycles.FindStringSubmatch(cycles)
  var cycleRange [2]int
  cycleRange[0], _ = strconv.Atoi(arrayCycles[1])
  cycleRange[1], _ = strconv.Atoi(arrayCycles[2])

  return cycleRange
}
