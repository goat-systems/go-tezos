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
    bar = pb.StartNew(5).Prefix("Taco Mission")
)

func main() {
  var cycleRange [2]int

  delegateAddr := flag.String("delegateaddr", "nil", "The address to your delegation service (eg. tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc)") //TODO Validate Input
  cycle := flag.Int("cycle", -1, "The cycle you are querying for.")
  cycles := flag.String("cycles", "nil", "A range of cycles to compute delegated contracts. Example: 10-20")
  fee := flag.Float64("fee", .05, "Your dynamic fee percentage noted as a decimal.")
  //alias := flag.String("alias", "nil", "The alias to your baking address on your node. Needed for to send XTZ for payouts.")
  report := flag.Bool("report", true, "Generates a list of all the Delegated Contracts for the request cycle(s).")
  payout := flag.Bool("payout", false, "Pays each of your delegated contracts their share less your percentage fee.")
  flag.Parse()


  var delegatedClients []goTezos.DelegatedClient

  _, err := goTezos.GetBalanceFor(*delegateAddr) //A dirty trick to check if the address is real
  if (err != nil){
    fmt.Println("Invalid Delegator Address " + *delegateAddr + ": GetBalanceFor(*delegateAddr) failed: " + err.Error())
    os.Exit(1)
  }

  if (*cycle != -1){
    delegatedClients = singleCycleOp(*cycle, *delegateAddr, *fee)
  } else if (*cycles != "nil"){
    cycleRange = parseCyclesInput(*cycles)
    delegatedClients = multiCycleOp(cycleRange[0], cycleRange[1],*delegateAddr, *fee)
  } else{
    fmt.Println("No cycle(s) provided. Exiting...")
    os.Exit(1)
  }

  if (*report == true){
    generateReport(delegatedClients)
    bar.Increment()
  } else if (*report == false && *payout == true){
    //goTezos.PayoutDelegatedContracts(delegatedClients, *alias)
    bar.Increment()
    fmt.Println("PayoutDelegatedContracts Function is temporary disabled for safety!")
  } else if (*report == true && *payout == true){
    generateReport(delegatedClients)
    //goTezos.PayoutDelegatedContracts(delegatedClients, *alias)
    bar.Increment()
    fmt.Println("PayoutDelegatedContracts Function is temporary disabled for safety!")
  }
  bar.FinishPrint("Tacos Have Been Made On This Day!")
}

func singleCycleOp(cycle int, delegateAddr string, fee float64, ) []goTezos.DelegatedClient{
  var delegatedClients []goTezos.DelegatedClient
  contracts, err := goTezos.GetDelegatedContractsForCycle(cycle, delegateAddr)
  if (err != nil){
    fmt.Println(err)
    os.Exit(-1)
  }
  for _, delegatedClientAddr := range contracts {
    delegatedClients = append(delegatedClients, goTezos.DelegatedClient{Address:delegatedClientAddr, Delegator:false, TotalPayout:0})
  }
  delegatedClients = append(delegatedClients, goTezos.DelegatedClient{Address:delegateAddr, Delegator:true, TotalPayout:0}) //Need to keep track of your own baking rewards, to avoid accidentally including them in your fee system.
  bar.Increment()
  delegatedClients = goTezos.SortDelegateContracts(delegatedClients) //Put the oldest contract at the begining of the array
  bar.Increment()
  delegatedClients, err = goTezos.CalculateAllCommitmentsForCycle(delegatedClients, cycle, fee)
  if (err != nil){
    fmt.Println(err)
    os.Exit(-1)
  }
  bar.Increment()
  delegatedClients = goTezos.CalculateAllTotalPayout(delegatedClients)
  bar.Increment()

  return delegatedClients
}

func multiCycleOp(cycleStart int, cycleEnd int, delegateAddr string, fee float64) []goTezos.DelegatedClient{
  var delegatedClients []goTezos.DelegatedClient
  contracts, err := goTezos.GetAllDelegatedContracts(delegateAddr)
  if (err != nil){
    fmt.Println(err)
    os.Exit(-1)
  }
  for _, delegatedClientAddr := range contracts {
    delegatedClients = append(delegatedClients, goTezos.DelegatedClient{Address:delegatedClientAddr, Delegator:false, TotalPayout:0})
  }
  delegatedClients = append(delegatedClients, goTezos.DelegatedClient{Address:delegateAddr, Delegator:true, TotalPayout:0}) //Need to keep track of your own baking rewards, to avoid accidentally including them in your fee system.
  bar.Increment()
  delegatedClients = goTezos.SortDelegateContracts(delegatedClients) //Put the oldest contract at the begining of the array
  bar.Increment()
  delegatedClients, err = goTezos.CalculateAllCommitmentsForCycles(delegatedClients, cycleStart, cycleEnd, fee)
  if (err != nil){
    fmt.Println(err)
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
func generateReport(delegatedClients []goTezos.DelegatedClient){
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
    panic(err)
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
