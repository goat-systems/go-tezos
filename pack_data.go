package gotezos

// reference ConseilJS/src/chain/tezos/TezosMessageUtils.ts

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
)

func PackUint64(num uint64) (string, error) {
	raw := uint2hex(num)
	data, err := hex.DecodeString(raw)
	if err != nil {
		return "", fmt.Errorf("PackUint64: decode raw, %s", err)
	}
	for index, d := range data[1:] {
		data[index+1] = d ^ 0x80
	}
	length := len(data)
	tempData := make([]byte, length)
	for index, d := range data {
		tempData[length-index-1] = d
	}
	result := hex.EncodeToString(tempData)

	return "0500" + result, nil
}

func PackInt64(num int64) string {
	if num == 0 {
		return "0500" + "00"
	}
	value := big.NewInt(num)
	n := new(big.Int).Abs(value)
	l := n.BitLen()
	arr := make([]uint64, 0)
	v := n
	big0x3f := big.NewInt(0x3f)
	big0x7f := big.NewInt(0x7f)
	big0x40 := big.NewInt(0x40)
	big0x80 := big.NewInt(0x80)
	big0 := big.NewInt(0)
	for i := 0; i < l; i += 7 {
		b := new(big.Int)
		if i == 0 {
			b = new(big.Int).And(v, big0x3f)
			v.Rsh(v, 6)
		} else {
			b = new(big.Int).And(v, big0x7f)
			v.Rsh(v, 7)
		}
		if value.Cmp(big0) < 0 && i == 0 {
			b = b.Or(b, big0x40)
		}
		if i+7 < l {
			b = b.Or(b, big0x80)
		}
		arr = append(arr, b.Uint64())
	}
	if l%7 == 0 {
		arr[len(arr)-1] = arr[len(arr)-1] | 0x80
		arr = append(arr, 1)
	}
	result := ""
	for _, a := range arr {
		result = result + twoCharAtTail(strconv.FormatUint(a, 16))
	}
	return "0500" + result
}

func uint2hex(num uint64) string {
	result := ""
	if num < 128 {
		result = strconv.FormatUint(num, 16)
		if len(result) == 1 {
			result = "0" + result
		}
	} else if num > 2147483648 {
		r := new(big.Int).SetUint64(num)
		zero := big.NewInt(0)
		big127 := big.NewInt(127)
		for r.Cmp(zero) > 0 {
			temp := new(big.Int).And(r, big127)
			tempStr := twoCharAtTail(hex.EncodeToString(temp.Bytes()))
			result = tempStr + result
			r.Rsh(r, 7)
		}
	} else {
		r := num
		for r > 0 {
			temp := r & 127
			tempStr := twoCharAtTail(strconv.FormatUint(temp, 16))
			result = tempStr + result
			r = r >> 7
		}
	}
	return result
}

func twoCharAtTail(str string) string {
	if len(str) == 0 {
		str = "00"
	} else if len(str) == 1 {
		str = "0" + str
	} else {
		str = str[len(str)-2:]
	}
	return str
}

func PackAddress(address string) (string, error) {
	data, err := decode(address)
	if err != nil {
		return "", fmt.Errorf("decode %s, %s", address, err)
	}
	// remove prefix
	data = data[3:]
	source := hex.EncodeToString(data)
	if len(source) > 40 {
		return "", fmt.Errorf("invalid source length %d", len(source))
	}

	switch address[:3] {
	case "tz1":
		source = "0000" + source
	case "tz2":
		source = "0001" + source
	case "tz3":
		source = "0002" + source
	case "KT1":
		source = "01" + source + "00"
	default:
		return "", fmt.Errorf("invalid prefix")
	}
	return "050a" + genDataLengthSection(len(source)/2) + source, nil
}

func PackStr(str string) string {
	valueSection := hex.EncodeToString([]byte(str))
	return "0501" + genDataLengthSection(len(str)) + valueSection
}

func genDataLengthSection(length int) string {
	lengthSection := "0000000" + strconv.FormatInt(int64(length), 16)
	lengthSection = lengthSection[len(lengthSection)-8:]
	return lengthSection
}
