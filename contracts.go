package gotezos

import (
	"fmt"

	"github.com/pkg/errors"
)

// MichelsonType is a string name for a michelson data type
type MichelsonType string

const (
	// MInt michelson type int
	MInt MichelsonType = "INT"
	// MString michelson type string
	MString MichelsonType = "STRING"
	// MBytes michelson type bytes
	MBytes MichelsonType = "BYTES"
	// MNat michelson type nat
	MNat MichelsonType = "NAT"
	// MBool michelson type bool
	MBool MichelsonType = "BOOL"
	// MUnit michelson type unit
	MUnit MichelsonType = "UNIT"
	// MList michelson type list
	MList MichelsonType = "LIST"
	// MPair michelson type pair
	MPair MichelsonType = "PAIR"
	// MOption michelson type option
	MOption MichelsonType = "OPTION"
	// MOr michelson type or
	MOr MichelsonType = "OR"
	// MSet michelson type set
	MSet MichelsonType = "SET"
	// MMap michelson type map
	MMap MichelsonType = "MAP"
	// MBigMap michelson type BigMap
	MBigMap MichelsonType = "BIGMAP"
)

/*
ContractStorage gets access the data of the contract.

Path:
	../<block_id>/context/contracts/<contract_id>/storage (GET)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-context-contracts-contract-id-storage

Parameters:

	blockhash:
		The hash of block (height) of which you want to make the query.

	KT1:
		The contract address.
*/
func (t *GoTezos) ContractStorage(blockhash string, KT1 string) ([]byte, error) {
	query := fmt.Sprintf("/chains/main/blocks/%s/context/contracts/%s/storage", blockhash, KT1)
	resp, err := t.get(query)
	if err != nil {
		return resp, errors.Wrap(err, "could not get storage '%s'")
	}
	return resp, nil
}

// MichelsonStorageTree is a tree representation of the storage parsed from a micheslon contract
type MichelsonStorageTree struct {
	Prim string
	Args []MichelsonData
}

// MichelsonData represents a michelson data type and exposes useful functions to compare them to go values
type MichelsonData interface {
	Equal(x interface{}) bool
	Type() string
}

// MichelsonInt is a michelson int
type MichelsonInt struct {
	Value int
}

// MichelsonString is a michelson string
type MichelsonString struct {
	Value string
}

// MichelsonBytes is a michelson bytes
type MichelsonBytes struct {
	Value []byte
}

// MichelsonNat is a michelson nat
type MichelsonNat struct {
	Value int
}

// MichelsonBool is a michelson bool
type MichelsonBool struct {
	Value bool
}

// MichelsonUnit is a michelson unit
type MichelsonUnit struct{}

// MichelsonList is a michelson list
type MichelsonList struct {
	Value []interface{}
}

// MichelsonPair is a michelson pair
type MichelsonPair struct {
	Value [2]interface{}
}

// MichelsonOption is a michelson option
type MichelsonOption struct {
	Value bool
}

// MichelsonOr is a michelson or
type MichelsonOr struct{}

// MicheslsonSet is a michelson set
type MicheslsonSet struct{}

// MichelsonMap is a michelson map
type MichelsonMap struct{}

// MichelsonBigMap is a michelson bigmap
type MichelsonBigMap struct{}

// Equal expects a golang integer and will compare it with a michelson integer
func (m *MichelsonInt) Equal(x interface{}) (bool, error) {
	y, ok := x.(int)
	if !ok {
		return false, errors.New("failed to compare non integer value with integer")
	}

	if m.Value == y {
		return true, nil
	}

	return false, nil
}

// Type returns the type of a MichelsonInt
func (m *MichelsonInt) Type(x interface{}) MichelsonType {
	return MInt
}

// Equal expects a golang string and will compare it with a michelson string
func (m *MichelsonString) Equal(x interface{}) (bool, error) {
	y, ok := x.(string)
	if !ok {
		return false, errors.New("failed to compare non string value with string")
	}

	if m.Value == y {
		return true, nil
	}

	return false, nil
}

// Type returns the type of a MichelsonString
func (m *MichelsonString) Type(x interface{}) MichelsonType {
	return MString
}

// Equal expects a golang bool and will compare it with a michelson bool
func (m *MichelsonBytes) Equal(x interface{}) (bool, error) {
	y, ok := x.([]byte)
	if !ok {
		return false, errors.New("failed to compare non byte value with byte")
	}

	if len(m.Value) != len(y) {
		return false, nil
	}

	for i := range m.Value {
		if m.Value[i] != y[i] {
			return false, nil
		}
	}

	return true, nil
}

// Type returns the type of a MichelsonString
func (m *MichelsonBytes) Type(x interface{}) MichelsonType {
	return MBytes
}

// Equal expects a positive golang integer and will compare it with a michelson nat
func (m *MichelsonNat) Equal(x interface{}) (bool, error) {
	y, ok := x.(int)
	if !ok || y < 0 {
		return false, errors.New("failed to compare non nat value with nat")
	}

	if m.Value == y {
		return true, nil
	}

	return true, nil
}

// Type returns the type of a MichelsonString
func (m *MichelsonNat) Type(x interface{}) MichelsonType {
	return MNat
}

// Equal expects a golang bool and will compare it with a michelson bool
func (m *MichelsonBool) Equal(x interface{}) (bool, error) {
	y, ok := x.(bool)
	if !ok {
		return false, errors.New("failed to compare non bool value with bool")
	}

	if m.Value == y {
		return true, nil
	}

	return true, nil
}

// Type returns the type of a MichelsonNat
func (m *MichelsonBool) Type(x interface{}) MichelsonType {
	return MBool
}

// Equal on a michelson type unit is not comparable to any go value
// and is here as a placeholder to satisfy the MichelsonData interface.
func (m *MichelsonUnit) Equal(x interface{}) (bool, error) {
	return false, errors.New("michelson type unit is not comparable")
}

// Type returns the type of a MichelsonNat
func (m *MichelsonUnit) Type(x interface{}) MichelsonType {
	return MUnit
}

// Equal expects an array of interfaces and will compare it with a michelson list
func (m *MichelsonList) Equal(x interface{}) (bool, error) {
	y, ok := x.([]interface{})
	if !ok {
		return false, errors.New("failed to compare non bool value with bool")
	}

	if len(y) != len(m.Value) {
		return false, nil
	}

	for i := range y {
		if m.Value[i] != y[i] {
			return false, nil
		}
	}

	return true, nil
}

// Type returns the type of a MichelsonNat
func (m *MichelsonList) Type(x interface{}) MichelsonType {
	return MList
}
