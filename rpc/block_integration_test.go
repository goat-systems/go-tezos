// +build integration

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

func Test_Integration_LiveBlocks(t *testing.T) {
	r, err := rpc.New(HOST)
	if ok := assert.Nil(t, err, "Failed to generate RPC client."); !ok {
		t.FailNow()
	}

	if _, _, err := r.LiveBlocks(&rpc.BlockIDHead{}); err != nil {
		t.FailNow()
	}
}

func Test_Integration_Metadata(t *testing.T) {
	r, err := rpc.New(HOST)
	if ok := assert.Nil(t, err, "Failed to generate RPC client."); !ok {
		t.FailNow()
	}

	if _, _, err := r.Metadata(&rpc.BlockIDHead{}); err != nil {
		t.FailNow()
	}
}

func Test_Integration_MinimalValidTime(t *testing.T) {
	r, err := rpc.New(HOST)
	if ok := assert.Nil(t, err, "Failed to generate RPC client."); !ok {
		t.FailNow()
	}

	if _, _, err := r.MinimalValidTime(rpc.MinimalValidTimeInput{
		BlockID: &rpc.BlockIDHead{},
	}); err != nil {
		t.FailNow()
	}
}

func Test_Integration_OperationHashes(t *testing.T) {
	r, err := rpc.New(HOST)
	if ok := assert.Nil(t, err, "Failed to generate RPC client."); !ok {
		t.FailNow()
	}

	if _, _, err := r.OperationHashes(rpc.OperationHashesInput{
		BlockID: &rpc.BlockIDHead{},
	}); err != nil {
		t.FailNow()
	}

	if _, _, err := r.OperationHashes(rpc.OperationHashesInput{
		BlockID:    &rpc.BlockIDHead{},
		ListOffset: "0",
	}); err != nil {
		t.FailNow()
	}

	if _, _, err := r.OperationHashes(rpc.OperationHashesInput{
		BlockID:         &rpc.BlockIDHead{},
		ListOffset:      "0",
		OperationOffset: "1",
	}); err != nil {
		t.FailNow()
	}
}

func Test_Integration_Operations(t *testing.T) {
	r, err := rpc.New(HOST)
	if ok := assert.Nil(t, err, "Failed to generate RPC client."); !ok {
		t.FailNow()
	}

	if _, _, err := r.Operations(rpc.OperationsInput{
		BlockID: &rpc.BlockIDHead{},
	}); err != nil {
		t.FailNow()
	}

	if _, _, err := r.Operations(rpc.OperationsInput{
		BlockID:    &rpc.BlockIDHead{},
		ListOffset: "0",
	}); err != nil {
		t.FailNow()
	}

	if _, _, err := r.Operations(rpc.OperationsInput{
		BlockID:         &rpc.BlockIDHead{},
		ListOffset:      "0",
		OperationOffset: "1",
	}); err != nil {
		t.FailNow()
	}
}

func Test_Integration_Protocols(t *testing.T) {
	r, err := rpc.New(HOST)
	if ok := assert.Nil(t, err, "Failed to generate RPC client."); !ok {
		t.FailNow()
	}

	if _, _, err := r.Protocols(&rpc.BlockIDHead{}); err != nil {
		t.FailNow()
	}
}

func Test_Integration_RequiredEndorsements(t *testing.T) {
	r, err := rpc.New(HOST)
	if ok := assert.Nil(t, err, "Failed to generate RPC client."); !ok {
		t.FailNow()
	}

	if _, _, err := r.RequiredEndorsements(rpc.RequiredEndorsementsInput{
		BlockID: &rpc.BlockIDHead{},
	}); err != nil {
		t.FailNow()
	}
}
