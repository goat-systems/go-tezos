package gotezos

import (
	"bytes"
	"fmt"
	"math"
	"math/big"
	"strconv"

	"github.com/btcsuite/btcutil/base58"
	"github.com/pkg/errors"
)

var (
	BranchPrefix    []byte = []byte{1, 52}
	ProposalPrefix  []byte = []byte{2, 170}
	SigPrefix       []byte = []byte{4, 130, 43}
	OperationPrefix []byte = []byte{29, 159, 109}
	ContextPrefix   []byte = []byte{79, 179}
)

func operationTags(kind string) int64 {
	tags := map[string]int64{
		"endorsement":                 0,
		"proposals":                   5,
		"ballot":                      6,
		"seed_nonce_revelation":       1,
		"double_endorsement_evidence": 2,
		"double_baking_evidence":      3,
		"activate_account":            4,
		"reveal":                      107,
		"transaction":                 108,
		"origination":                 109,
		"delegation":                  110,
	}

	return tags[kind]
}

func primTags(prim string) byte {
	tags := map[string]byte{
		"parameter":        0x00,
		"storage":          0x01,
		"code":             0x02,
		"False":            0x03,
		"Elt":              0x04,
		"Left":             0x05,
		"None":             0x06,
		"Pair":             0x07,
		"Right":            0x08,
		"Some":             0x09,
		"True":             0x0A,
		"Unit":             0x0B,
		"PACK":             0x0C,
		"UNPACK":           0x0D,
		"BLAKE2B":          0x0E,
		"SHA256":           0x0F,
		"SHA512":           0x10,
		"ABS":              0x11,
		"ADD":              0x12,
		"AMOUNT":           0x13,
		"AND":              0x14,
		"BALANCE":          0x15,
		"CAR":              0x16,
		"CDR":              0x17,
		"CHECK_SIGNATURE":  0x18,
		"COMPARE":          0x19,
		"CONCAT":           0x1A,
		"CONS":             0x1B,
		"CREATE_ACCOUNT":   0x1C,
		"CREATE_CONTRACT":  0x1D,
		"IMPLICIT_ACCOUNT": 0x1E,
		"DIP":              0x1F,
		"DROP":             0x20,
		"DUP":              0x21,
		"EDIV":             0x22,
		"EMPTY_MAP":        0x23,
		"EMPTY_SET":        0x24,
		"EQ":               0x25,
		"EXEC":             0x26,
		"FAILWITH":         0x27,
		"GE":               0x28,
		"GET":              0x29,
		"GT":               0x2A,
		"HASH_KEY":         0x2B,
		"IF":               0x2C,
		"IF_CONS":          0x2D,
		"IF_LEFT":          0x2E,
		"IF_NONE":          0x2F,
		"INT":              0x30,
		"LAMBDA":           0x31,
		"LE":               0x32,
		"LEFT":             0x33,
		"LOOP":             0x34,
		"LSL":              0x35,
		"LSR":              0x36,
		"LT":               0x37,
		"MAP":              0x38,
		"MEM":              0x39,
		"MUL":              0x3A,
		"NEG":              0x3B,
		"NEQ":              0x3C,
		"NIL":              0x3D,
		"NONE":             0x3E,
		"NOT":              0x3F,
		"NOW":              0x40,
		"OR":               0x41,
		"PAIR":             0x42,
		"PUSH":             0x43,
		"RIGHT":            0x44,
		"SIZE":             0x45,
		"SOME":             0x46,
		"SOURCE":           0x47,
		"SENDER":           0x48,
		"SELF":             0x49,
		"STEPS_TO_QUOTA":   0x4A,
		"SUB":              0x4B,
		"SWAP":             0x4C,
		"TRANSFER_TOKENS":  0x4D,
		"SET_DELEGATE":     0x4E,
		"UNIT":             0x4F,
		"UPDATE":           0x50,
		"XOR":              0x51,
		"ITER":             0x52,
		"LOOP_LEFT":        0x53,
		"ADDRESS":          0x54,
		"CONTRACT":         0x55,
		"ISNAT":            0x56,
		"CAST":             0x57,
		"RENAME":           0x58,
		"bool":             0x59,
		"contract":         0x5A,
		"int":              0x5B,
		"key":              0x5C,
		"key_hash":         0x5D,
		"lambda":           0x5E,
		"list":             0x5F,
		"map":              0x60,
		"big_map":          0x61,
		"nat":              0x62,
		"option":           0x63,
		"or":               0x64,
		"pair":             0x65,
		"set":              0x66,
		"signature":        0x67,
		"string":           0x68,
		"bytes":            0x69,
		"mutez":            0x6A,
		"timestamp":        0x6B,
		"unit":             0x6C,
		"operation":        0x6D,
		"address":          0x6E,
		"SLICE":            0x6F,
		"DIG":              0x70,
		"DUG":              0x71,
		"EMPTY_BIG_MAP":    0x72,
		"APPLY":            0x73,
		"chain_id":         0x74,
		"CHAIN_ID":         0x75,
	}

	return tags[prim]
}

func (t *Transaction) Forge_Prototype() ([]byte, error) {

	result := bytes.NewBuffer([]byte{})

	kind, err := forgeNat(operationTags(t.Kind))
	if err != nil {
		return []byte{}, errors.Wrap(err, "failed to forge kind")
	}

	source, err := forgeSource(t.Source)
	if err != nil {
		return []byte{}, errors.Wrap(err, "failed to forge source")
	}

	fee, err := forgeNat(t.Fee.Big.Int64())
	if err != nil {
		return []byte{}, errors.Wrap(err, "failed to forge fee")
	}

	counter, err := forgeNat(int64(t.Counter))
	if err != nil {
		return []byte{}, errors.Wrap(err, "failed to forge counter")
	}

	gasLimit, err := forgeNat(t.GasLimit.Big.Int64())
	if err != nil {
		return []byte{}, errors.Wrap(err, "failed to forge gas_limit")
	}

	storageLimit, err := forgeNat(t.StorageLimit.Big.Int64())
	if err != nil {
		return []byte{}, errors.Wrap(err, "failed to forge storage_limit")
	}

	amount, err := forgeNat(t.Amount.Big.Int64())
	if err != nil {
		return []byte{}, errors.Wrap(err, "failed to forge amount")
	}

	destination, err := forgeAddress(t.Destination)
	if err != nil {
		return []byte{}, errors.Wrap(err, "failed to forge destination")
	}

	result.Write(kind)
	result.Write(source)
	result.Write(fee)
	result.Write(counter)
	result.Write(gasLimit)
	result.Write(storageLimit)
	result.Write(amount)
	result.Write(destination)

	if t.Parameters != nil {
		result.Write(forgeBool(true))
		result.Write(forgeEntrypoint(t.Parameters.Entrypoint))

		micheline, err := forgeMicheline(&t.Parameters.Value)
		if err != nil {
			return []byte{}, errors.Wrap(err, "failed to forge parameters")
		}
		result.Write(forgeArray(micheline, 4))
	} else {
		result.Write(forgeBool(false))
	}

	return result.Bytes(), nil
}

func forgeNat(value int64) ([]byte, error) {
	if value < 0 {
		return nil, fmt.Errorf("nat value (%d) cannot be negative", value)
	}

	buf := bytes.NewBuffer([]byte{})
	more := true

	for more {
		b := byte(value & 0x7f)
		value >>= 7
		if value > 0 {
			b |= 0x80
		} else {
			more = false
		}

		buf.WriteByte(b)
	}

	return buf.Bytes(), nil
}

func forgeSource(source string) ([]byte, error) {
	var prefix string
	if len(source) != 36 {
		return []byte{}, fmt.Errorf("invalid length (%d!=36) source address", len(source))
	}
	prefix = source[0:3]
	buf := base58.Decode(source)[3:]

	switch prefix {
	case "tz1":
		buf = append([]byte{0}, buf...)
	case "tz2":
		buf = append([]byte{1}, buf...)
	case "tz3":
		buf = append([]byte{2}, buf...)
	default:
		return []byte{}, fmt.Errorf("invalid source prefix '%s'", prefix)
	}

	return buf, nil
}

func forgeAddress(address string) ([]byte, error) {
	var prefix string
	if len(address) != 36 {
		return []byte{}, fmt.Errorf("invalid length (%d!=36) source address", len(address))
	}
	prefix = address[0:3]
	buf := base58.Decode(address)[3:]

	switch prefix {
	case "tz1":
		buf = append([]byte{0, 0}, buf...)
	case "tz2":
		buf = append([]byte{0, 1}, buf...)
	case "tz3":
		buf = append([]byte{0, 2}, buf...)
	case "KT1":
		buf = append([]byte{1}, buf...)
		buf = append(buf, byte(0))
	default:
		return []byte{}, fmt.Errorf("invalid address prefix '%s'", prefix)
	}

	return buf, nil
}

func forgeBool(value bool) []byte {
	if true {
		return []byte{255}
	}

	return []byte{0}
}

func forgeEntrypoint(value string) []byte {
	buf := bytes.NewBuffer([]byte{})

	entrypointTags := map[string]byte{
		"default":         0,
		"root":            1,
		"do":              2,
		"set_delegate":    3,
		"remove_delegate": 4,
	}

	if val, ok := entrypointTags[value]; ok == true {
		buf.WriteByte(val)
	} else {
		buf.WriteByte(byte(255))
		buf.Write(forgeArray(bytes.NewBufferString(value).Bytes(), 1))
	}

	return buf.Bytes()
}

func forgeArray(value []byte, l int) []byte {
	buf := bytes.NewBuffer(reverseBytes([]byte(strconv.Itoa(len(value)))[0:l]))
	buf.Write(value)
	return buf.Bytes()
}

// static byte[] ForgeInt(int value)
//         {
//             var binary = Convert.ToString(Math.Abs(value), 2);

//             var pad = 6;
//             if ((binary.Length - 6) % 7 == 0)
//                 pad = binary.Length;
//             else if (binary.Length > 6)
//                 pad = binary.Length + 7 - (binary.Length - 6) % 7;

//             binary = binary.PadLeft(pad, '0');

//             var septets = new List<string>();

//             for (var i = 0; i <= pad / 7; i++)
//                 septets.Add(binary.Substring(7 * i, Math.Min(7, pad - 7 * i)));

//             septets.Reverse();

//             septets[0] = (value >= 0 ? "0" : "1") + septets[0];

//             var res = new byte[]{};

//             for (var i = 0; i < septets.Count; i++)
//             {
//                 var prefix = i == septets.Count - 1 ? "0" : "1";
//                 res = res.Concat(new []{Convert.ToByte(prefix + septets[i], 2)});
//             }

//             return res;
//         }

func forgeInt(value int) []byte {
	binary := strconv.FormatInt(int64(value), 2)
	lenBin := len(binary)

	pad := 6
	if (lenBin-6)%7 == 0 {
		pad = lenBin
	} else if lenBin > 6 {
		pad = lenBin + 7 - (lenBin-6)%7
	}

	binary = fmt.Sprintf("0%s", binary)
	septets := []string{}

	for i := 0; i <= pad/7; i++ {
		septets = append(septets, binary[7*i:int(math.Min(7, float64(pad-7*i)))])
	}

	septets = reverseStrings(septets)
	if value >= 0 {
		septets[0] = fmt.Sprintf("0%s", septets)
	} else {
		septets[0] = fmt.Sprintf("1%s", septets)
	}

	buf := bytes.NewBuffer([]byte{})
	for i := 0; i < len(septets); i++ {
		prefix := "1"
		if i == len(septets)-1 {
			prefix = "0"
		}
		n := new(big.Int)
		n.SetString(prefix+septets[i], 2)
		buf.Write(n.Bytes())
	}

	return buf.Bytes()
}

func reverseStrings(s []string) []string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}

	return s
}

func reverseBytes(s []byte) []byte {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}

	return s
}

func forgeMicheline(micheline *MichelineMichelsonV1Expression) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	lenTags := []map[bool]byte{
		{
			false: 3,
			true:  4,
		},
		{
			false: 5,
			true:  6,
		},
		{
			false: 7,
			true:  8,
		},
		{
			false: 9,
			true:  9,
		},
	}

	if micheline.MichelineMichelsonV1Expression != nil {
		buf.WriteByte(0x02)
		// TODO buf.Write(forgeArray())
	} else if micheline.Prim != "" {

	} else if micheline.Bytes != "" {
		buf.WriteByte(0x0A)

	} else if micheline.Int != "" {
		buf.WriteByte(0x00)
		i, err := strconv.Atoi(micheline.Int)
		if err != nil {
			return []byte{}, errors.New("failed to forge \"int\"")
		}

		buf.Write(forgeInt(i))
	} else if micheline.String != "" {
		buf.WriteByte(0x01)

	}

	return buf.Bytes(), nil
}

// static IEnumerable<byte> ForgeMicheline(JToken data)
// {
// 	var res = new List<byte>();

// 	#region Tags
// 	#endregion

// 	switch (data)
// 	{
// 		case JArray _:
// 			res.Add(0x02);
// 			res.AddRange(ForgeArray(data.Select(ForgeMicheline).SelectMany(x => x).ToArray()).ToList());
// //                    Console.WriteLine($"JArray {Hex.Convert(res.ToArray())}");
// 			break;
// 		case JObject _ when data["prim"] != null:
// 		{
// 			var argsLen = data["args"]?.Count() ?? 0;
// 			var annotsLen = data["annots"]?.Count() ?? 0;

// 			res.Add(lenTags[argsLen][annotsLen > 0]);
// 			res.Add(primTags[data["prim"].ToString()]);
// //                    Console.WriteLine($"Args {Hex.Convert(res.ToArray())}");

// 			if (argsLen > 0)
// 			{
// 				var args = data["args"].Select(ForgeMicheline).SelectMany(x => x);
// 				if (argsLen < 3)
// 				{
// 					res.AddRange(args.ToList());
// //                            Console.WriteLine($"argsLen > 0 {Hex.Convert(res.ToArray())}");
// 				}
// 				else
// 				{
// 					res.AddRange(ForgeArray(args.ToArray()));
// //                            Console.WriteLine($"argsLen <= 0 {Hex.Convert(res.ToArray())}");
// 				}
// 			}

// 			if (annotsLen > 0)
// 			{
// 				res.AddRange(ForgeArray(Encoding.UTF8.GetBytes(string.Join(" ", data["annots"]))));
// //                        Console.WriteLine($"annotsLen > 0 {Hex.Convert(res.ToArray())}");
// 			}

// 			else if (argsLen == 3)
// 				res.AddRange(new List<byte>{0,0,0,0}); /* new string('0', 8);*/
// //                    Console.WriteLine($"argsLen == 3 {Hex.Convert(res.ToArray())}");

// 			break;
// 		}
// 		case JObject _ when data["bytes"] != null:
// 			res.Add(0x0A);
// 			res.AddRange(ForgeArray(Hex.Parse(data["bytes"].Value<string>())));
// //                    Console.WriteLine($"Bytes {Hex.Convert(res.ToArray())}");
// 			break;
// 		case JObject _ when data["int"] != null:
// 			res.Add(0x00);
// 			res.AddRange(ForgeInt(data["int"].Value<int>()));
// //                    Console.WriteLine($"int {Hex.Convert(res.ToArray())}");
// 			break;
// 		case JObject _ when data["string"] != null:
// 			res.Add(0x01);
// 			res.AddRange(ForgeArray(Encoding.UTF8.GetBytes(data["string"].Value<string>())));
// //                    Console.WriteLine($"String {data["string"].Value<string>()} {Hex.Convert(res.ToArray())}");
// 			break;
// 		case JObject _:
// 			throw new ArgumentException($"Michelson forge error");
// 		default:
// 			throw new ArgumentException($"Michelson forge error");
// 	}

// 	return res;
// }
