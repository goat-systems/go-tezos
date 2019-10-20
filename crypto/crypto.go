package crypto

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type Prefix []byte

var (
	// For (de)constructing addresses
	Prefix_tz1       Prefix = []byte{6, 161, 159}
	Prefix_tz2       Prefix = []byte{6, 161, 151}
	Prefix_tz3       Prefix = []byte{6, 161, 164}
	Prefix_kt        Prefix = []byte{2, 90, 121}
	Prefix_edsk      Prefix = []byte{43, 246, 78, 7}
	Prefix_edsk2     Prefix = []byte{13, 15, 58, 7}
	Prefix_edpk      Prefix = []byte{13, 15, 37, 217}
	Prefix_edesk     Prefix = []byte{7, 90, 60, 179, 41}
	Prefix_edsig     Prefix = []byte{9, 245, 205, 134, 18}
	Prefix_branch    Prefix = []byte{1, 52}
	Prefix_watermark Prefix = []byte{3}
)

//B58cencode encodes a byte array into base58 with prefix
func B58cencode(payload []byte, prefix Prefix) string {
	n := make([]byte, (len(prefix) + len(payload)))
	for k := range prefix {
		n[k] = prefix[k]
	}
	for l := range payload {
		n[l+len(prefix)] = payload[l]
	}
	b58c := Encode(n)
	return b58c
}

func B58cdecode(payload string, prefix []byte) []byte {
	b58c, _ := Decode(payload)
	return b58c[len(prefix):]
}

//Helper Functions to round float64
func RoundPlus(f float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return round(f*shift) / shift
}

func round(f float64) float64 {
	return math.Floor(f + .5)
}

// checkAndRemovePrefixToHex(base58CheckEncodedPayload: string, tezosPrefix: Uint8Array): string {
//     const prefixHex = Buffer.from(tezosPrefix).toString('hex')
//     const payload = bs58check.decode(base58CheckEncodedPayload).toString('hex')
//     if (payload.startsWith(prefixHex)) {
//       return payload.substring(tezosPrefix.length * 2)
//     } else {
//       throw new Error('payload did not match prefix: ' + prefixHex)
//     }
//   }

func RemovePrefixToHex(s string, prefix Prefix) (string, error) {
	p := string(prefix)
	branch58, err := Decode(s) //base58 TODO clean up naming
	if err != nil {
		return "", errors.Wrap(err, "could not forge operation")
	}
	branch58str := string(branch58)
	if strings.HasPrefix(branch58str, p) {
		return branch58str[len(prefix)*2:], nil
	}

	return "", errors.New("could not forge operation, missing prefix")
}

func BigNumberToZarith(bigNum big.Int) (string, error) {
	bitString := fmt.Sprintf("%b", &bigNum)
	for len(bitString)%7 != 0 {
		bitString = "0" + bitString
	}

	var resHex strings.Builder
	for i := len(bitString); i > 0; i -= 7 {
		section := bitString[i-7:]
		if i == 7 {
			section = "0" + section
		} else {
			section = "1" + section
		}

		i, err := strconv.ParseInt(section, 2, 64)
		if err != nil {
			return "", errors.Wrap(err, "could not forge operation")
		}
		hexSection := fmt.Sprintf("%x", i)
		if len(hexSection)%2 == 0 {
			hexSection = "0" + hexSection
		}

		resHex.WriteString(hexSection)
	}

	return resHex.String(), nil
}

// bigNumberToZarith(inputNumber: BigNumber) {
//     let bitString: string = inputNumber.toString(2)
//     while (bitString.length % 7 !== 0) {
//       bitString = '0' + bitString // fill up with leading '0'
//     }

//     let resultHexString = ''
//     // because it's little endian we start from behind...
//     for (let i = bitString.length; i > 0; i -= 7) {
//       let bitStringSection: string = bitString.substring(i - 7, i)
//       if (i === 7) {
//         // the last byte will show it's the last with a leading '0'
//         bitStringSection = '0' + bitStringSection
//       } else {
//         // the others will show more will come with a leading '1'
//         bitStringSection = '1' + bitStringSection
//       }
//       let hexStringSection = parseInt(bitStringSection, 2).toString(16)

//       if (hexStringSection.length % 2) {
//         hexStringSection = '0' + hexStringSection
//       }

//       resultHexString += hexStringSection
//     }
//     return resultHexString
//   }
const alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

func Encode(dataBytes []byte) string {

	// Performing SHA256 twice
	sha256hash := sha256.New()
	sha256hash.Write(dataBytes)
	middleHash := sha256hash.Sum(nil)
	sha256hash = sha256.New()
	sha256hash.Write(middleHash)
	hash := sha256hash.Sum(nil)

	checksum := hash[:4]
	dataBytes = append(dataBytes, checksum...)

	// For all the "00" versions or any prepended zeros as base58 removes them
	zeroCount := 0
	for i := 0; i < len(dataBytes); i++ {
		if dataBytes[i] == 0 {
			zeroCount++
		} else {
			break
		}
	}

	// Performing base58 encoding
	encoded := b58encode(dataBytes)

	for i := 0; i < zeroCount; i++ {
		encoded = "1" + encoded
	}

	return encoded
}

func Decode(encoded string) ([]byte, error) {
	zeroCount := 0
	for i := 0; i < len(encoded); i++ {
		if encoded[i] == 49 {
			zeroCount++
		} else {
			break
		}
	}

	dataBytes, err := b58decode(encoded)
	if err != nil {
		return []byte{}, err
	}
	data, checksum := dataBytes[:len(dataBytes)-4], dataBytes[len(dataBytes)-4:]

	for i := 0; i < zeroCount; i++ {
		data = append([]byte{0}, data...)
	}

	// Performing SHA256 twice to validate checksum
	sha256hash := sha256.New()
	sha256hash.Write(data)
	middleHash := sha256hash.Sum(nil)
	sha256hash = sha256.New()
	sha256hash.Write(middleHash)
	hash := sha256hash.Sum(nil)

	if !reflect.DeepEqual(checksum, hash[:4]) {
		return []byte{}, errors.New("data and checksum don't match")
	}

	return data, nil
}

func b58encode(data []byte) string {
	var encoded string
	decimalData := new(big.Int)
	decimalData.SetBytes(data)
	divisor, zero := big.NewInt(58), big.NewInt(0)

	for decimalData.Cmp(zero) > 0 {
		mod := new(big.Int)
		decimalData.DivMod(decimalData, divisor, mod)
		encoded = string(alphabet[mod.Int64()]) + encoded
	}

	return encoded
}

func b58decode(data string) ([]byte, error) {
	decimalData := new(big.Int)
	alphabetBytes := []byte(alphabet)
	multiplier := big.NewInt(58)

	for _, value := range data {
		pos := bytes.IndexByte(alphabetBytes, byte(value))
		if pos == -1 {
			return nil, errors.New("character not found in alphabet")
		}
		decimalData.Mul(decimalData, multiplier)
		decimalData.Add(decimalData, big.NewInt(int64(pos)))
	}

	return decimalData.Bytes(), nil
}
