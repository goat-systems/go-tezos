package goTezos

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strconv"
)

func GetSignatureForOp(tezosClientPath, opBytes, walletAlias string) (string, error) {
	opBytes = "0x03" + opBytes
	cmd := exec.Command(tezosClientPath, "sign", "bytes", opBytes, "for", walletAlias)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		defer stdin.Close()

		io.WriteString(stdin, "abcd1234")
	}()
	out1, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", out1)
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Result: " + out.String())
	return out.String(), nil
	// out, err := exec.Command(tezosClientPath, "sign", "bytes", opBytes, "for", "alias").Output()
	// if err != nil {
	// 	fmt.Println("Err in exec")
	// 	return "", err
	// }
	// rtnStr, err := unMarshelString(out)
	// if err != nil {
	// 	return "", err
	// }
	// return rtnStr, nil
}

func ForgeMultiTransferOpertion(delegatedContracts []DelegatedContract, source string) (string, error) {
	var contents Conts
	var transOps []TransOp
	strCounter, err := GetCounterForDelegate(source)
	if err != nil {
		return "", err
	}
	intCounter, err := strconv.Atoi(strCounter)
	if err != nil {
		return "", err
	}

	for _, contract := range delegatedContracts {
		if contract.Address != source {
			intCounter++
			counter := strconv.Itoa(intCounter)
			pay := strconv.FormatFloat(contract.TotalPayout, 'f', 6, 64)
			i, err := strconv.ParseFloat(pay, 64)
			if err != nil {
				return "", errors.New("Could not get parse amount to payout: " + err.Error())
			}
			i = i * 1000000
			if i != 0 {
				transOps = append(transOps, TransOp{Kind: "transaction", Source: source, Fee: "0", GasLimit: "100", StorageLimit: "0", Amount: strconv.Itoa(int(i)), Destination: contract.Address, Counter: counter})
			}

		}
	}
	contents.Contents = transOps
	contents.Branch = "BKs9YjNzRbhwtiqjkHFXGZ4jkCDBRRdMmm5DKZHQ6o8jH7imzZU"

	forge := "/chains/main/blocks/head/helpers/forge/operations"
	fmt.Println(PrettyReport(contents))

	output, err := TezosRPCPost(forge, contents)
	if err != nil {
		return "", err
	}
	opBytes, err := unMarshelString(output)
	if err != nil {
		return "", err
	}
	return opBytes, nil
}

func GetCounterForDelegate(phk string) (string, error) {
	//8732/chains/main/blocks/head/context/contracts/tz1TP6gRyCSfgxauWwsbXi6s59MmvejcBAZa/counter
	rpc := "/chains/main/blocks/head/context/contracts/" + phk + "/counter"
	resp, err := TezosRPCGet(rpc)
	if err != nil {
		return "", err
	}
	rtnStr, err := unMarshelString(resp)
	if err != nil {
		return "", err
	}
	return rtnStr, nil
}
