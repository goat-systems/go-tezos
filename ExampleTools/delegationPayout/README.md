# delegationPayout: A Tool for Delegation Services

This is a tool for Tezos delegation services developed by me, DefinitelyNotAGoat. The purpose of this tool is to keep the books of your delegation service. This tool can give you a report for all delegated contracts in a cycle(s), and calculate what you owe them based off your dynamic percentage fee. Automatic payouts currently under testing and development.

If you would like to send me some coffee money:
```
tz1hyaA2mLUQLqQo3TVk6cQHXDc7xoKcBSbN
```

If you would like to delegate to me to show your support (5% dynamic fee):
```
tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc
```

## Installation
```
go get gopkg.in/cheggaaa/pb.v1
go build delegationPayout.go
```

You will also need to export the path to your `tezos-client`. Example:

```
export TEZOSPATH=/home/tezosuser/tezos
```

## Usage
In the below examples we will be using the delegationPayout tool to generate reports for a large delegation service, for now its Tezos.Community(Tezos.Community).

You can find all program options by running `./delegatePayout -help`:
```
Usage of ./delegatePayout:
  -cycle int
    	The cycle you are querying for. (default -1)
  -cycles string
    	A range of cycles to compute delegated contracts. Example: 10-20 (default "nil")
  -delegateaddr string
    	The address to your delegation service (eg. tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc) (default "nil")
  -fee float
    	Your dynamic fee percentage noted as a decimal. (default 0.05)
  -payout
    	Pays each of your delegated contracts their share less your percentage fee.
  -report
    	Generates a list of all the Delegated Contracts for the request cycle(s). (default true)
```

Let's say we want to generate a report for all the people delegated to Tezos.Community(Tezos.Community) for cycle 7. We would then run:
```
./delegatePayout -delegateaddr=tz1TDSmoZXwVevLTEvKCTHWpomG76oC9S2fJ -cycle=7 -report=true -fee=0.05
```
This will generate a JSON report called `./report.json`. Look below for an example of the report.
```
[
  {
    "Address": "tz1TDSmoZXwVevLTEvKCTHWpomG76oC9S2fJ",
    "Commitments": [
      {
        "Cycle": 7,
        "Amount": 3000.73,
        "SharePercentage": 0.00801117027950095,
        "GrossPayout": 822.627020150555,
        "NetPayout": 781.4956691430273,
        "Fee": 41.131351007527755
      }
    ],
    "Delegator": true,
    "TotalPayout": 781.4956691430273
  },
  {
    "Address": "KT18goog2JjfTYVBxG3h3L7ERAskpJpCNsaK",
    "Commitments": [
      {
        "Cycle": 7,
        "Amount": 5020.218,
        "SharePercentage": 0.013402679094158987,
        "GrossPayout": 1376.2541027837156,
        "NetPayout": 1307.4413976445298,
        "Fee": 68.81270513918578
      }
    ],
    "Delegator": false,
    "TotalPayout": 1307.4413976445298
  },
  {
    "Address": "KT19Dh3Xd52vMXdnSZAkAwJPJV67G2kGFQxf",
    "Commitments": [
      {
        "Cycle": 7,
        "Amount": 2148.06,
        "SharePercentage": 0.005734762684608348,
        "GrossPayout": 588.8741062690083,
        "NetPayout": 559.4304009555578,
        "Fee": 29.443705313450415
      }
    ],
    "Delegator": false,
    "TotalPayout": 559.4304009555578
  }
]
```

Let's say we want to generate a report for all the people delegated for cycle 7-10. We would then run:
```
./delegatePayout -delegateaddr=tz1TDSmoZXwVevLTEvKCTHWpomG76oC9S2fJ -cycles=7-10 -report=true -fee=0.05
```

It will produce a similar report as we saw below, but with multiple commitments in a contract:



### Note 1:
When you run the tool, there will be a loading bar. The loading bar does hang on 2/5 because it is computing a lot of data. This is especially true for doing calculations over multiple cycles. I believe the slowness is caused by the RPC client, and not the program itself. 

```
Taco Mission 2 / 5 [====------]  40.00%
```

When the report is finished you will see this:

```
Taco Mission 5 / 5 [==========]  100.00%
Tacos Have Been Made On This Day!
```

### Note 2:
This tool does not yet have the ability to calculate the total rewards received by a delegate in a cycle. I will be implementing this shortly. But to test the functionality of the calculations I generate a random total reward amount between 70,000 - 105,000 (XTZ). See the function below to see what I mean:


```
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
```


## Authors

* **DefinitelyNotAGoat**

See also the list of [contributors](https://github.com/DefinitelyNotAGoat/goTezos/graphs/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
