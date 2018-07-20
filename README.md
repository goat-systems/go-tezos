# goTezos: A Tezos Go Library

The purpose of this library is to allow developers to build go driven applications for Tezos. I have built a tool called delegatePayout
as an example use of this library. You can find this tool in the ExampleTools directory.

If you would like to send me some coffee money:
```
tz1hyaA2mLUQLqQo3TVk6cQHXDc7xoKcBSbN
```

If you would like to delegate to me to show your support (5% dynamic fee):
```
tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc
```


More robust documentation will come soon.

## Installation
```
go get github.com/DefinitelyNotAGoat/goTezos
```

To use the library import it into your go application:
```
import "github.com/DefinitelyNotAGoat/goTezos"
```


## goTezos Documentation
The goTezos Library requires the use of an env variable called TEZOSPATH.


Example:

```
export TEZOSPATH=/home/tezosuser/tezos
```

I will create a wiki shortly describing the functions available.


## goTezos ExampleTools
As of now there is only one tool. I hope to develop more tools, as I continue development on goTezos. The one tool available is delegatePayout. If you haven't guessed, it's a tool to generate the books for a delegation service on any given cycle or a range of cycles. The tool is still in development, but you can test it's functionality.

See the [README.md](https://github.com/DefinitelyNotAGoat/goTezos/blob/master/ExampleTools/delegationPayout/README.md) for more information.

## Authors

* **DefinitelyNotAGoat**

See also the list of [contributors](https://github.com/DefinitelyNotAGoat/goTezos/graphs/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
