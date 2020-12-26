package rpc_test

import (
	"testing"

	"github.com/goat-systems/go-tezos/v4/rpc"
	"github.com/stretchr/testify/assert"
)

func Test_Integration_Block(t *testing.T) {
	r, err := rpc.New(HOST)
	if ok := assert.Nil(t, err, "Failed to generate RPC client."); !ok {
		t.FailNow()
	}

	for i := 0; i < 20; i++ {
		getRandomBlock(r, 0, t)
	}
}

// TODO: func Test_Integration_EndorsingPower(t *testing.T) {}

func Test_Integration_Hash(t *testing.T) {
	r, err := rpc.New(HOST)
	if ok := assert.Nil(t, err, "Failed to generate RPC client."); !ok {
		t.FailNow()
	}

	if _, _, err := r.Hash(&rpc.BlockIDHead{}); err != nil {
		t.FailNow()
	}
}

func Test_Integration_Header(t *testing.T) {
	r, err := rpc.New(HOST)
	if ok := assert.Nil(t, err, "Failed to generate RPC client."); !ok {
		t.FailNow()
	}

	if _, _, err := r.Header(&rpc.BlockIDHead{}); err != nil {
		t.FailNow()
	}
}

func Test_Integration_HeaderRaw(t *testing.T) {
	r, err := rpc.New(HOST)
	if ok := assert.Nil(t, err, "Failed to generate RPC client."); !ok {
		t.FailNow()
	}

	if _, _, err := r.HeaderRaw(&rpc.BlockIDHead{}); err != nil {
		t.FailNow()
	}
}

func Test_Integration_HeaderShell(t *testing.T) {
	r, err := rpc.New(HOST)
	if ok := assert.Nil(t, err, "Failed to generate RPC client."); !ok {
		t.FailNow()
	}

	if _, _, err := r.HeaderShell(&rpc.BlockIDHead{}); err != nil {
		t.FailNow()
	}
}

func Test_Integration_HeaderProtocolData(t *testing.T) {
	r, err := rpc.New(HOST)
	if ok := assert.Nil(t, err, "Failed to generate RPC client."); !ok {
		t.FailNow()
	}

	if _, _, err := r.HeaderProtocolData(&rpc.BlockIDHead{}); err != nil {
		t.FailNow()
	}
}

func Test_Integration_HeaderProtocolDataRaw(t *testing.T) {
	r, err := rpc.New(HOST)
	if ok := assert.Nil(t, err, "Failed to generate RPC client."); !ok {
		t.FailNow()
	}

	if _, _, err := r.HeaderProtocolDataRaw(&rpc.BlockIDHead{}); err != nil {
		t.FailNow()
	}
}
