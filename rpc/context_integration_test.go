package rpc_test

import (
	"math/rand"
	"testing"

	"github.com/goat-systems/go-tezos/v4/forge"
	"github.com/goat-systems/go-tezos/v4/rpc"
	"github.com/stretchr/testify/assert"
)

const HOST = "https://mainnet-tezos.giganode.io"

func getRandomBlock(r *rpc.Client, from int, t *testing.T) *rpc.Block {
	_, head, err := r.Block(&rpc.BlockIDHead{})
	if ok := assert.Nil(t, err, "Random block generator failed to get current network height"); !ok {
		t.FailNow()
	}

	id := rpc.BlockIDLevel(rand.Intn((head.Header.Level - from)) + from)
	_, block, err := r.Block(&id)
	if ok := assert.Nil(t, err, "Random block generator failed to get block"); !ok {
		t.FailNow()
	}

	return block
}

func Test_Integration_BigMap(t *testing.T) {
	r, err := rpc.New(HOST)
	if ok := assert.Nil(t, err, "Failed to generate RPC client."); !ok {
		t.FailNow()
	}

	scriptExp, err := forge.ForgeAddressExpression("tz1UKGJ98tjySyWkFFtaPRXFjoYZrzb5rhPD")
	if ok := assert.Nil(t, err, "Failed to forge script expression."); !ok {
		t.FailNow()
	}

	block := getRandomBlock(r, 1208511, t)

	id := rpc.BlockIDHash(block.Hash)
	_, err = r.BigMap(rpc.BigMapInput{
		BlockID:          &id,
		BigMapID:         123, // tzBTC big_map ID
		ScriptExpression: scriptExp,
	})
	if ok := assert.Nilf(t, err, "Failed at block '%d'", block.Header.Level); !ok {
		t.FailNow()
	}
}

func Test_Integration_Constants(t *testing.T) {
	r, err := rpc.New(HOST)
	if ok := assert.Nil(t, err, "Failed to generate RPC client."); !ok {
		t.FailNow()
	}

	for i := 0; i < 5; i++ {
		block := getRandomBlock(r, 0, t)

		id := rpc.BlockIDHash(block.Hash)
		_, _, err = r.Constants(rpc.ConstantsInput{
			BlockID: &id,
		})
		if ok := assert.Nilf(t, err, "Failed at block '%d'", block.Header.Level); !ok {
			t.FailNow()
		}
	}
}

func Test_Integration_Contracts(t *testing.T) {
	r, err := rpc.New(HOST)
	if ok := assert.Nil(t, err, "Failed to generate RPC client."); !ok {
		t.FailNow()
	}

	block := getRandomBlock(r, 0, t)

	id := rpc.BlockIDHash(block.Hash)
	_, _, err = r.Contracts(rpc.ContractsInput{
		BlockID: &id,
	})
	if ok := assert.Nilf(t, err, "Failed at block '%d'", block.Header.Level); !ok {
		t.FailNow()
	}
}

func Test_Integration_Contract(t *testing.T) {
	r, err := rpc.New(HOST)
	if ok := assert.Nil(t, err, "Failed to generate RPC client."); !ok {
		t.FailNow()
	}

	for i := 0; i < 5; i++ {
		block := getRandomBlock(r, 1208511, t)

		id := rpc.BlockIDHash(block.Hash)
		_, _, err := r.Contract(rpc.ContractInput{
			BlockID:    &id,
			ContractID: "tz1ZbSrRrfhU8LYHELWNswx2JcFARXTGKKVk",
		})
		if ok := assert.Nilf(t, err, "Failed at block '%d'", block.Header.Level); !ok {
			t.FailNow()
		}
	}
}

func Test_Integration_ContractBalance(t *testing.T) {
	r, err := rpc.New(HOST)
	if ok := assert.Nil(t, err, "Failed to generate RPC client."); !ok {
		t.FailNow()
	}

	for i := 0; i < 5; i++ {
		block := getRandomBlock(r, 1127225, t)

		id := rpc.BlockIDHash(block.Hash)
		_, _, err := r.ContractBalance(rpc.ContractBalanceInput{
			BlockID:    &id,
			ContractID: "tz1ZbSrRrfhU8LYHELWNswx2JcFARXTGKKVk",
		})
		if ok := assert.Nilf(t, err, "Failed at block '%d'", block.Header.Level); !ok {
			t.FailNow()
		}
	}
}

func Test_Integration_ContractCounter(t *testing.T) {
	r, err := rpc.New(HOST)
	if ok := assert.Nil(t, err, "Failed to generate RPC client."); !ok {
		t.FailNow()
	}

	for i := 0; i < 5; i++ {
		block := getRandomBlock(r, 1127225, t)

		id := rpc.BlockIDHash(block.Hash)
		_, _, err := r.ContractCounter(rpc.ContractCounterInput{
			BlockID:    &id,
			ContractID: "tz1ZbSrRrfhU8LYHELWNswx2JcFARXTGKKVk",
		})
		if ok := assert.Nilf(t, err, "Failed at block '%d'", block.Header.Level); !ok {
			t.FailNow()
		}
	}
}

func Test_Integration_ContractDelegate(t *testing.T) {
	r, err := rpc.New(HOST)
	if ok := assert.Nil(t, err, "Failed to generate RPC client."); !ok {
		t.FailNow()
	}

	for i := 0; i < 5; i++ {
		block := getRandomBlock(r, 1275219, t)

		id := rpc.BlockIDHash(block.Hash)
		_, _, err := r.ContractDelegate(rpc.ContractDelegateInput{
			BlockID:    &id,
			ContractID: "tz1ZbSrRrfhU8LYHELWNswx2JcFARXTGKKVk",
		})
		if ok := assert.Nilf(t, err, "Failed at block '%d'", block.Header.Level); !ok {
			t.FailNow()
		}
	}
}

func Test_Integration_ContractEntrypoints(t *testing.T) {
	r, err := rpc.New(HOST)
	if ok := assert.Nil(t, err, "Failed to generate RPC client."); !ok {
		t.FailNow()
	}

	block := getRandomBlock(r, 1269959, t)

	id := rpc.BlockIDHash(block.Hash)
	_, _, err = r.ContractEntrypoints(rpc.ContractEntrypointsInput{
		BlockID:    &id,
		ContractID: "KT1DrJV8vhkdLEj76h1H9Q4irZDqAkMPo1Qf",
	})
	if ok := assert.Nilf(t, err, "Failed at block '%d'", block.Header.Level); !ok {
		t.FailNow()
	}
}

func Test_Integration_ContractEntrypoint(t *testing.T) {
	r, err := rpc.New(HOST)
	if ok := assert.Nil(t, err, "Failed to generate RPC client."); !ok {
		t.FailNow()
	}

	block := getRandomBlock(r, 1269959, t)

	id := rpc.BlockIDHash(block.Hash)
	_, _, err = r.ContractEntrypoint(rpc.ContractEntrypointInput{
		BlockID:    &id,
		ContractID: "KT1DrJV8vhkdLEj76h1H9Q4irZDqAkMPo1Qf",
		Entrypoint: "tokenToToken",
	})
	if ok := assert.Nilf(t, err, "Failed at block '%d'", block.Header.Level); !ok {
		t.FailNow()
	}
}

func Test_Integration_ContractManagerKey(t *testing.T) {
	r, err := rpc.New(HOST)
	if ok := assert.Nil(t, err, "Failed to generate RPC client."); !ok {
		t.FailNow()
	}

	block := getRandomBlock(r, 1269959, t)

	id := rpc.BlockIDHash(block.Hash)
	_, _, err = r.ContractManagerKey(rpc.ContractManagerKeyInput{
		BlockID:    &id,
		ContractID: "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
	})
	if ok := assert.Nilf(t, err, "Failed at block '%d'", block.Header.Level); !ok {
		t.FailNow()
	}
}

func Test_Integration_ContractScript(t *testing.T) {
	r, err := rpc.New(HOST)
	if ok := assert.Nil(t, err, "Failed to generate RPC client."); !ok {
		t.FailNow()
	}

	block := getRandomBlock(r, 1269959, t)

	id := rpc.BlockIDHash(block.Hash)
	_, err = r.ContractScript(rpc.ContractScriptInput{
		BlockID:    &id,
		ContractID: "KT1DrJV8vhkdLEj76h1H9Q4irZDqAkMPo1Qf",
	})
	if ok := assert.Nilf(t, err, "Failed at block '%d'", block.Header.Level); !ok {
		t.FailNow()
	}
}

// TODO: EDO sapling contracts are not readily available
//func Test_Integration_ContractSaplingDiff(t *testing.T) {}

func Test_Integration_ContractStorage(t *testing.T) {
	r, err := rpc.New(HOST)
	if ok := assert.Nil(t, err, "Failed to generate RPC client."); !ok {
		t.FailNow()
	}

	block := getRandomBlock(r, 1269959, t)

	id := rpc.BlockIDHash(block.Hash)
	_, err = r.ContractStorage(rpc.ContractStorageInput{
		BlockID:    &id,
		ContractID: "KT1DrJV8vhkdLEj76h1H9Q4irZDqAkMPo1Qf",
	})
	if ok := assert.Nilf(t, err, "Failed at block '%d'", block.Header.Level); !ok {
		t.FailNow()
	}
}

func Test_Integration_Delegates(t *testing.T) {
	r, err := rpc.New(HOST)
	if ok := assert.Nil(t, err, "Failed to generate RPC client."); !ok {
		t.FailNow()
	}

	block := getRandomBlock(r, 0, t)

	id := rpc.BlockIDHash(block.Hash)
	_, _, err = r.Delegates(rpc.DelegatesInput{
		BlockID: &id,
	})
	if ok := assert.Nilf(t, err, "Failed at block '%d'", block.Header.Level); !ok {
		t.FailNow()
	}
}

func Test_Integration_Delegate(t *testing.T) {
	r, err := rpc.New(HOST)
	if ok := assert.Nil(t, err, "Failed to generate RPC client."); !ok {
		t.FailNow()
	}

	block := getRandomBlock(r, 1000000, t)

	id := rpc.BlockIDHash(block.Hash)
	_, _, err = r.Delegate(rpc.DelegateInput{
		BlockID:  &id,
		Delegate: "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
	})
	if ok := assert.Nilf(t, err, "Failed at block '%d'", block.Header.Level); !ok {
		t.FailNow()
	}
}

func Test_Integration_DelegateBalance(t *testing.T) {
	r, err := rpc.New(HOST)
	if ok := assert.Nil(t, err, "Failed to generate RPC client."); !ok {
		t.FailNow()
	}

	block := getRandomBlock(r, 1000000, t)

	id := rpc.BlockIDHash(block.Hash)
	_, _, err = r.DelegateBalance(rpc.DelegateBalanceInput{
		BlockID:  &id,
		Delegate: "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
	})
	if ok := assert.Nilf(t, err, "Failed at block '%d'", block.Header.Level); !ok {
		t.FailNow()
	}
}

func Test_Integration_DelegateDeactivated(t *testing.T) {
	r, err := rpc.New(HOST)
	if ok := assert.Nil(t, err, "Failed to generate RPC client."); !ok {
		t.FailNow()
	}

	block := getRandomBlock(r, 1000000, t)

	id := rpc.BlockIDHash(block.Hash)
	_, _, err = r.DelegateDeactivated(rpc.DelegateDeactivatedInput{
		BlockID:  &id,
		Delegate: "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
	})
	if ok := assert.Nilf(t, err, "Failed at block '%d'", block.Header.Level); !ok {
		t.FailNow()
	}
}

func Test_Integration_DelegateDelegatedBalance(t *testing.T) {
	r, err := rpc.New(HOST)
	if ok := assert.Nil(t, err, "Failed to generate RPC client."); !ok {
		t.FailNow()
	}

	block := getRandomBlock(r, 1000000, t)

	id := rpc.BlockIDHash(block.Hash)
	_, _, err = r.DelegateDelegatedBalance(rpc.DelegateDelegatedBalanceInput{
		BlockID:  &id,
		Delegate: "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
	})
	if ok := assert.Nilf(t, err, "Failed at block '%d'", block.Header.Level); !ok {
		t.FailNow()
	}
}

func Test_Integration_DelegateDelegatedContracts(t *testing.T) {
	r, err := rpc.New(HOST)
	if ok := assert.Nil(t, err, "Failed to generate RPC client."); !ok {
		t.FailNow()
	}

	block := getRandomBlock(r, 1000000, t)

	id := rpc.BlockIDHash(block.Hash)
	_, _, err = r.DelegateDelegatedContracts(rpc.DelegateDelegatedContractsInput{
		BlockID:  &id,
		Delegate: "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
	})
	if ok := assert.Nilf(t, err, "Failed at block '%d'", block.Header.Level); !ok {
		t.FailNow()
	}
}

func Test_Integration_DelegateFrozenBalance(t *testing.T) {
	r, err := rpc.New(HOST)
	if ok := assert.Nil(t, err, "Failed to generate RPC client."); !ok {
		t.FailNow()
	}

	block := getRandomBlock(r, 1000000, t)

	id := rpc.BlockIDHash(block.Hash)
	_, _, err = r.DelegateFrozenBalance(rpc.DelegateFrozenBalanceInput{
		BlockID:  &id,
		Delegate: "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
	})
	if ok := assert.Nilf(t, err, "Failed at block '%d'", block.Header.Level); !ok {
		t.FailNow()
	}
}

func Test_Integration_DelegateFrozenBalanceByCycle(t *testing.T) {
	r, err := rpc.New(HOST)
	if ok := assert.Nil(t, err, "Failed to generate RPC client."); !ok {
		t.FailNow()
	}

	block := getRandomBlock(r, 1000000, t)

	id := rpc.BlockIDHash(block.Hash)
	_, _, err = r.DelegateFrozenBalanceByCycle(rpc.DelegateFrozenBalanceByCycleInput{
		BlockID:  &id,
		Delegate: "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
	})
	if ok := assert.Nilf(t, err, "Failed at block '%d'", block.Header.Level); !ok {
		t.FailNow()
	}
}

func Test_Integration_DelegateGracePeriod(t *testing.T) {
	r, err := rpc.New(HOST)
	if ok := assert.Nil(t, err, "Failed to generate RPC client."); !ok {
		t.FailNow()
	}

	block := getRandomBlock(r, 1000000, t)

	id := rpc.BlockIDHash(block.Hash)
	_, _, err = r.DelegateGracePeriod(rpc.DelegateGracePeriodInput{
		BlockID:  &id,
		Delegate: "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
	})
	if ok := assert.Nilf(t, err, "Failed at block '%d'", block.Header.Level); !ok {
		t.FailNow()
	}
}

func Test_Integration_DelegateStakingBalance(t *testing.T) {
	r, err := rpc.New(HOST)
	if ok := assert.Nil(t, err, "Failed to generate RPC client."); !ok {
		t.FailNow()
	}

	block := getRandomBlock(r, 1000000, t)

	id := rpc.BlockIDHash(block.Hash)
	_, _, err = r.DelegateStakingBalance(rpc.DelegateStakingBalanceInput{
		BlockID:  &id,
		Delegate: "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
	})
	if ok := assert.Nilf(t, err, "Failed at block '%d'", block.Header.Level); !ok {
		t.FailNow()
	}
}

func Test_Integration_DelegateVotingPower(t *testing.T) {
	r, err := rpc.New(HOST)
	if ok := assert.Nil(t, err, "Failed to generate RPC client."); !ok {
		t.FailNow()
	}

	block := getRandomBlock(r, 1000000, t)

	id := rpc.BlockIDHash(block.Hash)
	_, _, err = r.DelegateVotingPower(rpc.DelegateVotingPowerInput{
		BlockID:  &id,
		Delegate: "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
	})
	if ok := assert.Nilf(t, err, "Failed at block '%d'", block.Header.Level); !ok {
		t.FailNow()
	}
}

func Test_Integration_Nonces(t *testing.T) {
	r, err := rpc.New(HOST)
	if ok := assert.Nil(t, err, "Failed to generate RPC client."); !ok {
		t.FailNow()
	}

	block := getRandomBlock(r, 1000000, t)

	id := rpc.BlockIDHash(block.Hash)
	_, _, err = r.Nonces(rpc.NoncesInput{
		BlockID: &id,
		Level:   block.Header.Level - 1,
	})
	if ok := assert.Nilf(t, err, "Failed at block '%d'", block.Header.Level); !ok {
		t.FailNow()
	}
}

func Test_Integration_RawBytes(t *testing.T) {
	r, err := rpc.New(HOST)
	if ok := assert.Nil(t, err, "Failed to generate RPC client."); !ok {
		t.FailNow()
	}

	block := getRandomBlock(r, 1000000, t)

	id := rpc.BlockIDHash(block.Hash)
	_, err = r.RawBytes(rpc.RawBytesInput{
		BlockID: &id,
		Depth:   1,
	})
	if ok := assert.Nilf(t, err, "Failed at block '%d'", block.Header.Level); !ok {
		t.FailNow()
	}
}

// TODO: EDO sapling contracts are not readily available
// func Test_SaplingDiff(t *testing.T) {}

func Test_Integration_Seed(t *testing.T) {
	r, err := rpc.New(HOST)
	if ok := assert.Nil(t, err, "Failed to generate RPC client."); !ok {
		t.FailNow()
	}

	_, head, err := r.Block(&rpc.BlockIDHead{})
	if ok := assert.Nil(t, err, "Failed to generate RPC client."); !ok {
		t.FailNow()
	}

	id := rpc.BlockIDHash(head.Hash)
	_, _, err = r.Seed(rpc.SeedInput{
		BlockID: &id,
	})
	if ok := assert.Nilf(t, err, "Failed at block '%d'", head.Header.Level); !ok {
		t.FailNow()
	}
}
