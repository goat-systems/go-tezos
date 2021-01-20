// +build integration

package rpc_test

import (
	"testing"

	"github.com/goat-systems/go-tezos/v4/rpc"
	"github.com/stretchr/testify/assert"
)

func Test_Integration_BakingRights(t *testing.T) {
	r, err := rpc.New(HOST)
	if ok := assert.Nil(t, err, "Failed to generate RPC client."); !ok {
		t.FailNow()
	}

	if _, _, err := r.BakingRights(rpc.BakingRightsInput{
		BlockID: &rpc.BlockIDHead{},
	}); err != nil {
		t.FailNow()
	}
}

func Test_Integration_CompletePrefix(t *testing.T) {
	r, err := rpc.New(HOST)
	if ok := assert.Nil(t, err, "Failed to generate RPC client."); !ok {
		t.FailNow()
	}

	if _, _, err := r.CompletePrefix(rpc.CompletePrefixInput{
		BlockID: &rpc.BlockIDHead{},
		Prefix:  "tz1",
	}); err != nil {
		t.FailNow()
	}
}

func Test_Integration_CurrentLevel(t *testing.T) {
	r, err := rpc.New(HOST)
	if ok := assert.Nil(t, err, "Failed to generate RPC client."); !ok {
		t.FailNow()
	}

	if _, _, err := r.CurrentLevel(rpc.CurrentLevelInput{
		BlockID: &rpc.BlockIDHead{},
	}); err != nil {
		t.FailNow()
	}
}

func Test_Integration_EndorsingRights(t *testing.T) {
	r, err := rpc.New(HOST)
	if ok := assert.Nil(t, err, "Failed to generate RPC client."); !ok {
		t.FailNow()
	}

	if _, _, err := r.EndorsingRights(rpc.EndorsingRightsInput{
		BlockID: &rpc.BlockIDHead{},
	}); err != nil {
		t.FailNow()
	}
}
