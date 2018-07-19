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
go get github.com/DefinitelyNotAGoat/goTezos
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

Let's say we want to generate a report for all the people delegated to [Tezos.Community](https://www.tezos.community/) for cycle 7. We would then run:
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

It will produce a similar report as we saw above, but with multiple commitments in a contract:

```
[
  {
    "Address": "KT1XnbBbxyHgyrkpVYDzCQdspyvJMGxUUwHQ",
    "Commitments": [
      {
        "Cycle": 7,
        "Amount": 0,
        "SharePercentage": 0,
        "GrossPayout": 0,
        "NetPayout": 0,
        "Fee": 0
      },
      {
        "Cycle": 8,
        "Amount": 10396.31,
        "SharePercentage": 0.0020119177581536915,
        "GrossPayout": 141.2507100466962,
        "NetPayout": 134.1881745443614,
        "Fee": 7.062535502334811
      },
      {
        "Cycle": 9,
        "Amount": 10396.31,
        "SharePercentage": 0.0019680386659993486,
        "GrossPayout": 176.62163008011154,
        "NetPayout": 167.79054857610598,
        "Fee": 8.831081504005578
      },
      {
        "Cycle": 10,
        "Amount": 10396.31,
        "SharePercentage": 0.0012310976368299208,
        "GrossPayout": 121.93283434218267,
        "NetPayout": 115.83619262507354,
        "Fee": 6.096641717109134
      }
    ],
    "Delegator": false,
    "TotalPayout": 417.81491574554093
  },
  {
    "Address": "KT1XoEvvXkDKkF7wDbP32GQWACXwJviUNJAd",
    "Commitments": [
      {
        "Cycle": 7,
        "Amount": 0,
        "SharePercentage": 0,
        "GrossPayout": 0,
        "NetPayout": 0,
        "Fee": 0
      },
      {
        "Cycle": 8,
        "Amount": 1798,
        "SharePercentage": 0.0003479530842347273,
        "GrossPayout": 24.4287421848675,
        "NetPayout": 23.207305075624124,
        "Fee": 1.2214371092433751
      },
      {
        "Cycle": 9,
        "Amount": 1798,
        "SharePercentage": 0.00034036437173062643,
        "GrossPayout": 30.546000540965068,
        "NetPayout": 29.018700513916816,
        "Fee": 1.5273000270482535
      },
      {
        "Cycle": 10,
        "Amount": 1798,
        "SharePercentage": 0.0002129133847509547,
        "GrossPayout": 21.087793279273555,
        "NetPayout": 20.033403615309876,
        "Fee": 1.0543896639636778
      }
    ],
    "Delegator": false,
    "TotalPayout": 72.25940920485081
  },
  {
    "Address": "KT1XtS1RTVhwWHNV9Nb1npRC18TpAu2Jg9rX",
    "Commitments": [
      {
        "Cycle": 7,
        "Amount": 0,
        "SharePercentage": 0,
        "GrossPayout": 0,
        "NetPayout": 0,
        "Fee": 0
      },
      {
        "Cycle": 8,
        "Amount": 5999.742999,
        "SharePercentage": 0.0011610840273179992,
        "GrossPayout": 81.51622630591477,
        "NetPayout": 77.44041499061903,
        "Fee": 4.075811315295739
      },
      {
        "Cycle": 9,
        "Amount": 5999.742999,
        "SharePercentage": 0.0011357612660733369,
        "GrossPayout": 101.92889482375162,
        "NetPayout": 96.83245008256404,
        "Fee": 5.0964447411875815
      },
      {
        "Cycle": 10,
        "Amount": 5999.742999,
        "SharePercentage": 0.0007104702945233224,
        "GrossPayout": 70.36781985076794,
        "NetPayout": 66.84942885822954,
        "Fee": 3.5183909925383974
      }
    ],
    "Delegator": false,
    "TotalPayout": 241.12229393141263
  }
  ]

```


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
