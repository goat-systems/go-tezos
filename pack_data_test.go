package gotezos

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestPackString(t *testing.T) {
	str := "123123123"
	serialized := PackStr(str)
	assert.Equal(t, "050100000009313233313233313233", serialized)
}

func TestPackUint64(t *testing.T) {
	serialized, err := PackUint64(1200000)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "0500809f49", serialized)
	serialized, err = PackUint64(2147483748)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "0500e480808008", serialized)
}

func TestPackInt64(t *testing.T) {
	serialized := PackInt64(-10)
	assert.Equal(t, "05004a", serialized)
	serialized = PackInt64(-11110)
	assert.Equal(t, "0500e6ad01", serialized)
	serialized = PackInt64(11110)
	assert.Equal(t, "0500a6ad01", serialized)
	serialized = PackInt64(0)
	assert.Equal(t, "050000", serialized)
	serialized = PackInt64(2147483748)
	assert.Equal(t, "0500a481808010", serialized)
}

func TestPackAddress(t *testing.T) {
	address := "tz1WfxXzNgcsMrQBV6YHChmyQgNVvfo44ihi"
	serialized, err := PackAddress(address)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "050a0000001600007906bec05de5c0bbaf5f6062fc33096ec29e9f30", serialized)

	address = "KT1TPBnfPq7XpCjL11HTzAeBkyxrSxUuskS9"
	serialized, err = PackAddress(address)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "050a0000001601ce399e401363480455445f652d8bbaed24b52b3400", serialized)
}
