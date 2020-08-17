package gotezos

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
	"time"

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

func (r *Reveal) Forge_Prototype() ([]byte, error) {
	result := bytes.NewBuffer([]byte{})

	if kind, err := forgeNat(operationTags(r.Kind)); err == nil {
		result.Write(kind)
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge kind")
	}

	if source, err := forgeSource(r.Source); err == nil {
		result.Write(source)
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge source")
	}

	if fee, err := forgeNat(r.Fee.Big.Int64()); err == nil {
		result.Write(fee)
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge fee")
	}

	if counter, err := forgeNat(int64(r.Counter)); err == nil {
		result.Write(counter)
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge counter")
	}

	if gasLimit, err := forgeNat(r.GasLimit.Big.Int64()); err == nil {
		result.Write(gasLimit)
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge gas_limit")
	}

	if storageLimit, err := forgeNat(r.StorageLimit.Big.Int64()); err == nil {
		result.Write(storageLimit)
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge storage_limit")
	}

	if publicKey, err := forgePublicKey(r.PublicKey); err == nil {
		result.Write(publicKey)
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge public_key")
	}

	return result.Bytes(), nil
}

func (a *AccountActivation) Forge_Prototype() ([]byte, error) {
	result := bytes.NewBuffer([]byte{})

	if kind, err := forgeNat(operationTags(a.Kind)); err == nil {
		result.Write(kind)
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge kind")
	}

	if pkh, err := forgeActivationAddress(a.Pkh); err == nil {
		result.Write(pkh)
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge pkh")
	}

	if secret, err := hex.DecodeString(a.Secret); err == nil {
		result.Write(secret)
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge secret")
	}

	return result.Bytes(), nil
}

func (t *Transaction) Forge_Prototype() ([]byte, error) {
	result := bytes.NewBuffer([]byte{})

	if kind, err := forgeNat(operationTags(t.Kind)); err == nil {
		result.Write(kind)
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge kind")
	}

	if source, err := forgeSource(t.Source); err == nil {
		result.Write(source)
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge source")
	}

	if fee, err := forgeNat(t.Fee.Big.Int64()); err == nil {
		result.Write(fee)
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge fee")
	}

	if counter, err := forgeNat(int64(t.Counter)); err == nil {
		result.Write(counter)
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge counter")
	}

	if gasLimit, err := forgeNat(t.GasLimit.Big.Int64()); err == nil {
		result.Write(gasLimit)
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge gas_limit")
	}

	if storageLimit, err := forgeNat(t.StorageLimit.Big.Int64()); err == nil {
		result.Write(storageLimit)
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge storage_limit")
	}

	if amount, err := forgeNat(t.Amount.Big.Int64()); err == nil {
		result.Write(amount)
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge amount")
	}

	if destination, err := forgeAddress(t.Destination); err == nil {
		result.Write(destination)
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge destination")
	}

	if t.Parameters != nil {
		result.Write(forgeBool(true))
		result.Write(forgeEntrypoint(t.Parameters.Entrypoint))

		if micheline, err := forgeMicheline(&t.Parameters.Value); err == nil {
			result.Write(forgeArray(micheline, 4))
		} else {
			return []byte{}, errors.Wrap(err, "failed to forge parameters")
		}

	} else {
		result.Write(forgeBool(false))
	}

	return result.Bytes(), nil
}

func (o *Origination) Forge_Prototype() ([]byte, error) {
	result := bytes.NewBuffer([]byte{})

	if kind, err := forgeNat(operationTags(o.Kind)); err == nil {
		result.Write(kind)
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge kind")
	}

	if source, err := forgeSource(o.Source); err == nil {
		result.Write(source)
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge source")
	}

	if fee, err := forgeNat(o.Fee.Big.Int64()); err == nil {
		result.Write(fee)
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge fee")
	}

	if counter, err := forgeNat(int64(o.Counter)); err == nil {
		result.Write(counter)
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge counter")
	}

	if gasLimit, err := forgeNat(o.GasLimit.Big.Int64()); err == nil {
		result.Write(gasLimit)
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge gas_limit")
	}

	if storageLimit, err := forgeNat(o.StorageLimit.Big.Int64()); err == nil {
		result.Write(storageLimit)
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge storage_limit")
	}

	if balance, err := forgeNat(o.Balance.Big.Int64()); err == nil {
		result.Write(balance)
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge balance")
	}

	if o.Delegate != "" {
		result.Write(forgeBool(true))
		if delegate, err := forgeSource(o.Delegate); err == nil {
			result.Write(delegate)
		} else {
			return []byte{}, errors.Wrap(err, "failed to forge delegate")
		}
	} else {
		result.Write(forgeBool(false))
	}

	if script, err := forgeScript(o.Script); err == nil {
		result.Write(script)
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge script")
	}

	return result.Bytes(), nil
}

func (d *Delegation) Forge_Prototype() ([]byte, error) {
	result := bytes.NewBuffer([]byte{})

	if kind, err := forgeNat(operationTags(d.Kind)); err == nil {
		result.Write(kind)
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge kind")
	}

	if source, err := forgeSource(d.Source); err == nil {
		result.Write(source)
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge source")
	}

	if fee, err := forgeNat(d.Fee.Big.Int64()); err == nil {
		result.Write(fee)
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge fee")
	}

	if counter, err := forgeNat(int64(d.Counter)); err == nil {
		result.Write(counter)
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge counter")
	}

	if gasLimit, err := forgeNat(d.GasLimit.Big.Int64()); err == nil {
		result.Write(gasLimit)
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge gas_limit")
	}

	if storageLimit, err := forgeNat(d.StorageLimit.Big.Int64()); err == nil {
		result.Write(storageLimit)
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge storage_limit")
	}

	if d.Delegate != "" {
		result.Write(forgeBool(true))
		if source, err := forgeSource(d.Delegate); err == nil {
			result.Write(source)
		} else {
			return []byte{}, errors.Wrap(err, "failed to forge delegate")
		}
	} else {
		result.Write(forgeBool(false))
	}

	return result.Bytes(), nil
}

func (e *Endorsement) Forge_Prototype() ([]byte, error) {
	result := bytes.NewBuffer([]byte{})

	if kind, err := forgeNat(operationTags(e.Kind)); err == nil {
		result.Write(kind)
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge kind")
	}

	result.Write(forgeInt32(e.Level, 4))
	return result.Bytes(), nil
}

func (s *SeedNonceRevelation) Forge_Prototype() ([]byte, error) {
	result := bytes.NewBuffer([]byte{})

	if kind, err := forgeNat(operationTags(s.Kind)); err == nil {
		result.Write(kind)
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge kind")
	}

	result.Write(forgeInt32(s.Level, 4))

	if nonce, err := hex.DecodeString(s.Nonce); err == nil {
		result.Write(nonce)
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge nonce")
	}

	return result.Bytes(), nil
}

func (p *Proposal) Forge_Prototype() ([]byte, error) {
	result := bytes.NewBuffer([]byte{})

	if kind, err := forgeNat(operationTags(p.Kind)); err == nil {
		result.Write(kind)
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge kind")
	}

	if source, err := forgeSource(p.Source); err == nil {
		result.Write(source)
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge source")
	}

	result.Write(forgeInt32(p.Period, 4))

	buf := bytes.NewBuffer([]byte{})
	for _, proposal := range p.Proposals {
		if p, err := prefixAndBase58Encode(proposal, ProposalPrefix); err == nil {
			buf.Write([]byte(p))
		} else {
			return []byte{}, errors.Wrap(err, "failed to forge proposals")
		}
	}

	result.Write(forgeArray(buf.Bytes(), 4))
	return result.Bytes(), nil
}

func (b *Ballot) Forge_Prototype() ([]byte, error) {
	result := bytes.NewBuffer([]byte{})

	if kind, err := forgeNat(operationTags(b.Kind)); err == nil {
		result.Write(kind)
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge kind")
	}

	if source, err := forgeSource(b.Source); err == nil {
		result.Write(source)
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge source")
	}

	result.Write(forgeInt32(b.Period, 4))

	if p, err := prefixAndBase58Encode(b.Proposal, ProposalPrefix); err == nil {
		result.Write([]byte(p))
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge proposal")
	}

	result.Write([]byte(b.Ballot))

	return result.Bytes(), nil
}

func (d *DoubleEndorsementEvidence) Forge_Prototype() ([]byte, error) {
	result := bytes.NewBuffer([]byte{})

	if kind, err := forgeNat(operationTags(d.Kind)); err == nil {
		result.Write(kind)
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge kind")
	}

	if op1, err := d.Op1.forge(); err == nil {
		result.Write(op1)
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge op1")
	}

	if op2, err := d.Op2.forge(); err == nil {
		result.Write(op2)
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge op2")
	}

	return result.Bytes(), nil
}

func (d *DoubleBakingEvidence) Forge_Prototype() ([]byte, error) {
	result := bytes.NewBuffer([]byte{})

	if kind, err := forgeNat(operationTags(d.Kind)); err == nil {
		result.Write(kind)
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge kind")
	}

	if bh1, err := d.Bh1.forge(); err == nil {
		result.Write(bh1)
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge bh1")
	}

	if bh2, err := d.Bh2.forge(); err == nil {
		result.Write(bh2)
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge bh2")
	}

	return result.Bytes(), nil
}

func (i *InlinedEndorsement) forge() ([]byte, error) {
	result := bytes.NewBuffer([]byte{})
	if branch, err := prefixAndBase58Encode(i.Branch, BranchPrefix); err == nil {
		result.Write([]byte(branch))
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge branch")
	}

	if kind, err := forgeNat(operationTags(i.Operations.Kind)); err == nil {
		result.Write(kind)
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge operations kind")
	}

	result.Write(forgeInt32(i.Operations.Level, 4))

	if signature, err := prefixAndBase58Encode(i.Signature, SigPrefix); err == nil {
		result.Write([]byte(signature))
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge signature")
	}

	return forgeArray(result.Bytes(), 4), nil
}

func (b *BlockHeader) forge() ([]byte, error) {
	result := bytes.NewBuffer([]byte{})
	result.Write(forgeInt32(b.Level, 4))
	result.Write(forgeInt32(b.Proto, 1))

	if predecessor, err := prefixAndBase58Encode(b.Predecessor, BranchPrefix); err == nil {
		result.Write([]byte(predecessor))
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge predecessor")
	}

	ts := int(b.Timestamp.Sub(time.Date(1970, time.January, 1, 0, 0, 0, 0, nil)).Seconds())
	result.Write(forgeInt32(ts, 8))
	result.Write(forgeInt32(b.ValidationPass, 1))

	if operationHash, err := prefixAndBase58Encode(b.OperationsHash, OperationPrefix); err == nil {
		result.Write([]byte(operationHash))
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge operation_hash")
	}

	buf := bytes.NewBuffer([]byte{})
	for _, f := range b.Fitness {
		if fitness, err := hex.DecodeString(f); err == nil {
			buf.Write(forgeArray(fitness, 4))
		} else {
			return []byte{}, errors.Wrap(err, "failed to forge fitness")
		}
	}
	result.Write(forgeArray(buf.Bytes(), 4))

	if context, err := prefixAndBase58Encode(b.Context, ContextPrefix); err == nil {
		result.Write([]byte(context))
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge context")
	}

	result.Write(forgeInt32(b.Priority, 2))

	if proofOfWorkNonce, err := hex.DecodeString(b.ProofOfWorkNonce); err == nil {
		buf.Write(proofOfWorkNonce)
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge proof_of_work_nonce")
	}

	if signature, err := forgeSignature(b.Signature); err == nil {
		buf.Write(signature)
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge signature")
	}

	return forgeArray(result.Bytes(), 4), nil
}

func forgeSignature(value string) ([]byte, error) {
	var buf []byte
	if signature, err := prefixAndBase58Encode(value, SigPrefix); err == nil {
		buf = []byte(signature)
	} else {
		return []byte{}, errors.Wrap(err, "failed to forge context")
	}

	return append([]byte{0}, buf...), nil
}

func forgeInt32(value int, l int) []byte {
	return reverseBytes([]byte{byte(value)}[0:l])
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
	if len(address) != 36 {
		return []byte{}, fmt.Errorf("invalid length (%d!=36) source address", len(address))
	}
	prefix := address[0:3]
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
	buf := bytes.NewBuffer(reverseBytes([]byte{byte(len(value))}[0:l]))
	buf.Write(value)
	return buf.Bytes()
}

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

func forgePublicKey(value string) ([]byte, error) {
	if len(value) < 54 {
		return []byte{}, fmt.Errorf("invalid public key '%s'", value)
	}

	prefix := value[0:4]
	buf := base58.Decode(value)[4:]

	switch prefix {
	case "edpk":
		buf = append([]byte{0}, buf...)
	case "sppk":
		buf = append([]byte{1}, buf...)
	case "p2pk":
		buf = append([]byte{2}, buf...)
	default:
		return []byte{}, fmt.Errorf("invalid public key prefix '%s'", prefix)
	}

	return buf, nil
}

func forgeActivationAddress(value string) ([]byte, error) {
	buf := base58.Decode(value)
	if len(buf) < 3 {
		return []byte{}, fmt.Errorf("invalid activation address '%s'", value)
	}

	return buf[3:], nil
}

func forgeScript(script Script) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	if michline, err := forgeMicheline(script.Code); err == nil {
		buf.Write(forgeArray(michline, 4))
	} else {
		return []byte{}, err
	}

	if michline, err := forgeMicheline(script.Storage); err == nil {
		buf.Write(forgeArray(michline, 4))
	} else {
		return []byte{}, err
	}

	return buf.Bytes(), nil
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
		tmpBuf := bytes.NewBuffer([]byte{})
		for _, m := range micheline.MichelineMichelsonV1Expression {
			v, err := forgeMicheline(&m)
			if err != nil {
				return []byte{}, errors.New("failed to michline array \"int\"")
			}
			tmpBuf.Write(v)
		}

		buf.Write(forgeArray(tmpBuf.Bytes(), 4))
	} else if micheline.Prim != "" {
		argsLen := len(micheline.Args)
		annotsLen := len(micheline.Annots)

		buf.WriteByte(lenTags[argsLen][annotsLen > 0])
		buf.WriteByte(primTags(micheline.Prim))

		if argsLen > 0 {
			args := bytes.NewBuffer([]byte{})
			for _, arg := range micheline.Args {
				v, err := forgeMicheline(&arg)
				if err != nil {
					return []byte{}, errors.New("failed to michline array \"int\"")
				}
				args.Write(v)
			}

			if argsLen < 3 {
				buf.Write(args.Bytes())
			} else {
				buf.Write(forgeArray(args.Bytes(), 4))
			}
		}

		if annotsLen > 0 {
			buf.Write(forgeArray([]byte(strings.Join(micheline.Annots, " ")), 4))
		} else if argsLen == 3 {
			buf.Write([]byte{0, 0, 0, 0})
		}
	} else if micheline.Bytes != "" {
		buf.WriteByte(0x0A)
		bytes, err := hex.DecodeString(micheline.Bytes)
		if err != nil {
			return []byte{}, errors.New("failed to forge \"bytes\"")
		}

		buf.Write(forgeArray(bytes, 4))
	} else if micheline.Int != "" {
		buf.WriteByte(0x00)
		i, err := strconv.Atoi(micheline.Int)
		if err != nil {
			return []byte{}, errors.New("failed to forge \"int\"")
		}

		buf.Write(forgeInt(i))
	} else if micheline.String != "" {
		buf.WriteByte(0x01)
		buf.Write(forgeArray([]byte(micheline.String), 4))
	}

	return buf.Bytes(), nil
}
