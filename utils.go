package goTezos
/*
Author: DefinitelyNotAGoat/MagicAglet
Version: 0.0.1
Description: Functions not necessarily related to the Tezos API, but used in goTezos
License: MIT
*/

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

/*
Description: Asks user a yes or no question and returns true or false.
Param s (string): Question to be asked
Returns (bool): true or false
*/
func askForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
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
