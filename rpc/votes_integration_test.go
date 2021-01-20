package rpc_test

import (
	"testing"

	"github.com/goat-systems/go-tezos/v4/rpc"
	"github.com/stretchr/testify/assert"
)

const HOST = "https://mainnet-tezos.giganode.io"

func Test_Integration_BallotList(t *testing.T) {
	r, err := rpc.New(HOST)
	if ok := assert.Nil(t, err, "Failed to generate RPC client."); !ok {
		t.FailNow()
	}

	if _, _, err := r.BallotList(&rpc.BlockIDHead{}); err != nil {
		t.FailNow()
	}
}

func Test_Integration_Ballots(t *testing.T) {
	r, err := rpc.New(HOST)
	if ok := assert.Nil(t, err, "Failed to generate RPC client."); !ok {
		t.FailNow()
	}

	if _, _, err := r.Ballots(&rpc.BlockIDHead{}); err != nil {
		t.FailNow()
	}
}

func Test_Integration_CurrentPeriodKind(t *testing.T) {
	r, err := rpc.New(HOST)
	if ok := assert.Nil(t, err, "Failed to generate RPC client."); !ok {
		t.FailNow()
	}

	if _, _, err := r.CurrentPeriodKind(&rpc.BlockIDHead{}); err != nil {
		t.FailNow()
	}
}

func Test_Integration_CurrentProposal(t *testing.T) {
	r, err := rpc.New(HOST)
	if ok := assert.Nil(t, err, "Failed to generate RPC client."); !ok {
		t.FailNow()
	}

	if _, _, err := r.CurrentProposal(&rpc.BlockIDHead{}); err != nil {
		t.FailNow()
	}
}

func Test_Integration_CurrentQuorum(t *testing.T) {
	r, err := rpc.New(HOST)
	if ok := assert.Nil(t, err, "Failed to generate RPC client."); !ok {
		t.FailNow()
	}

	if _, _, err := r.CurrentQuorum(&rpc.BlockIDHead{}); err != nil {
		t.FailNow()
	}
}

func Test_Integration_Listings(t *testing.T) {
	r, err := rpc.New(HOST)
	if ok := assert.Nil(t, err, "Failed to generate RPC client."); !ok {
		t.FailNow()
	}

	if _, _, err := r.Listings(&rpc.BlockIDHead{}); err != nil {
		t.FailNow()
	}
}
