package main

/*
Author: DefinitelyNotAGoat
Description: This is a program to test all functions in the goTezos library.
*/

import (
	"fmt"
	"github.com/DefinitelyNotAGoat/goTezos"
	"io/ioutil"
	"json"
	"os"
)

/*
Description: A structure to store the contents of config.json
Param delegateAddr (string): The address of the delegate to use for the unit tests
Param contractAddr (string): A contact address to use (KT1)
Param cycle (int): The cycle number you would like to unit test converts
Param cycles (string): For functions that need multiple cycles
Param fee (float64): The fee for delegation functions
*/
type Config struct {
	delegateAddr      string
	contractAddr      string
	cycle             int
	cycles            string
	fee               float64
	TestSendFunctions bool
}

func main() {
	file, e := ioutil.ReadFile("./config.json")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	var conf Config
	json.Unmarshal(file, &conf)

	//**********Testing goTezos.go**********

	//Testing func GetSnapShot(cycle int) (SnapShot, error)
	snapShot, err := goTezos.GetSnapShot(cycle)
	errorCheck("func GetSnapShot(cycle int) (SnapShot, error) failed", err)
	fmt.Println("func GetSnapShot(cycle int) (SnapShot, error) success! Output: ")
	fmt.Println(snapShot + "\n")

	//Testing GetBlockLevelHead() (int, error)
	level, err := goTezos.GetBlockLevelHead()
	errorCheck("func GetBlockLevelHead() (int, error) failed", err)
	fmt.Println("func GetBlockLevelHead() (int, error) success! Output: ")
	fmt.Println(level + "\n")

	//Testing GetBlockLevelHash(level int) (string, error)
	hash, err := goTezos.GetBlockLevelHash(snapShot.AssociatedBlock)
	errorCheck("func GetBlockLevelHash(level int) (string, error) failed", err)
	fmt.Println("func GetBlockLevelHash(level int) (string, error) success! Output: ")
	fmt.Println(hash + "\n")

	//Testing GetBalanceFor(tezosAddr string) (float64, error)
	balance, err := GetBalanceFor(delegateAddr)(float64, error)
	errorCheck("func GetBalanceFor(tezosAddr string) (float64, error) failed", err)
	fmt.Println("func GetBalanceFor(tezosAddr string) (float64, error) success! Output: ")
	fmt.Println(balance + "\n")

	//Testing func GetBalanceAtSnapShotFor(tezosAddr string, cycle int) (float64, error)
	balanceAtSnapShot, err := GetBalanceAtSnapShotFor(contractAddr, cycle)
	errorCheck("func GetBalanceAtSnapShotFor(tezosAddr string, cycle int) (float64, error) failed", err)
	fmt.Println("func GetBalanceAtSnapShotFor(tezosAddr string, cycle int) (float64, error) success! Output: ")
	fmt.Println(balanceAtSnapShot + "\n")

	if TestSendFunctions {
    //Testing SendTezos(amount float64, toAddress string, alias string) error
    err := SendTezos("0.5", contractAddr, alias)
    errorCheck("func SendTezos(amount float64, toAddress string, alias string) error failed", err)
    fmt.Println("func SendTezos(amount float64, toAddress string, alias string) error success!")

    err := SafeSendTezos("0.5", contractAddr, alias)
    errorCheck("func SendTezos(amount float64, toAddress string, alias string) error failed", err)
    fmt.Println("func SendTezos(amount float64, toAddress string, alias string) error success!")
	}

  //Testing ListKownAddresses() ([]KnownAddress, error)

}

func errorCheck(msg string, err error) {
	if err != nil {
		fmt.Println(msg + ": " + err.Error())
	}
}
