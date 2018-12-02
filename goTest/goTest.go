package main

import (
  "fmt"
  "github.com/fkbenjamin/go-tezos"
)

//Please use another rpc for real money tx. This is just for testing purpose!
var url = "https://rpc.tezrpc.me:443"

//Used to show how to use added features to Create Batch Payments. The signed operations are not injeced but rather returned as an array.
func main() {
  goTezos.SetRPCURL(url);

  var testPay goTezos.Payment
  testPay.Address = "tz1fHyywxNwfEgxzCj95hCGXPdTTPtj2C5BA"
  testPay.Amount = 1

  var payments []goTezos.Payment
  for i := 0; i < 100; i++ {
    payments = append(payments, testPay)
  }
  dec_sigs := goTezos.CreateBatchPayment(payments)
  fmt.Println(dec_sigs)
}
