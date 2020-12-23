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
	_, head, err := r.Head()
	if ok := assert.Nil(t, err, "Random block generator failed to get current network height"); !ok {
		t.FailNow()
	}

	level := rand.Intn((head.Header.Level - from)) + from
	_, block, err := r.Block(level)
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
	_, err = r.BigMap(rpc.BigMapInput{
		Blockhash:        block.Hash,
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
		_, _, err = r.Constants(rpc.ConstantsInput{
			Blockhash: block.Hash,
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
	_, _, err = r.Contracts(rpc.ContractsInput{
		Blockhash: block.Hash,
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
		_, _, err := r.Contract(rpc.ContractInput{
			Blockhash:  block.Hash,
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
		_, _, err := r.ContractBalance(rpc.ContractBalanceInput{
			Blockhash:  block.Hash,
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
		_, _, err := r.ContractCounter(rpc.ContractCounterInput{
			Blockhash:  block.Hash,
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
		block := getRandomBlock(r, 1127225, t)
		_, _, err := r.ContractDelegate(rpc.ContractDelegateInput{
			Blockhash:  block.Hash,
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
	_, _, err = r.ContractEntrypoints(rpc.ContractEntrypointsInput{
		Blockhash:  block.Hash,
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
	_, _, err = r.ContractEntrypoint(rpc.ContractEntrypointInput{
		Blockhash:  block.Hash,
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
	_, _, err = r.ContractManagerKey(rpc.ContractManagerKeyInput{
		Blockhash:  block.Hash,
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
	_, err = r.ContractScript(rpc.ContractScriptInput{
		Blockhash:  block.Hash,
		ContractID: "KT1DrJV8vhkdLEj76h1H9Q4irZDqAkMPo1Qf",
	})
	if ok := assert.Nilf(t, err, "Failed at block '%d'", block.Header.Level); !ok {
		t.FailNow()
	}
}

// TODO
//func Test_Integration_ContractSaplingDiff(t *testing.T) {}

func Test_Integration_ContractStorage(t *testing.T) {
	r, err := rpc.New(HOST)
	if ok := assert.Nil(t, err, "Failed to generate RPC client."); !ok {
		t.FailNow()
	}

	block := getRandomBlock(r, 1269959, t)
	_, err = r.ContractStorage(rpc.ContractStorageInput{
		Blockhash:  block.Hash,
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
	_, _, err = r.Delegates(rpc.DelegatesInput{
		Blockhash: block.Hash,
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
	_, _, err = r.Delegate(rpc.DelegateInput{
		Blockhash: block.Hash,
		Delegate:  "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
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
	_, _, err = r.DelegateBalance(rpc.DelegateBalanceInput{
		Blockhash: block.Hash,
		Delegate:  "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
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
	_, _, err = r.DelegateDeactivated(rpc.DelegateDeactivatedInput{
		Blockhash: block.Hash,
		Delegate:  "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
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
	_, _, err = r.DelegateDelegatedBalance(rpc.DelegateDelegatedBalanceInput{
		Blockhash: block.Hash,
		Delegate:  "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
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
	_, _, err = r.DelegateDelegatedContracts(rpc.DelegateDelegatedContractsInput{
		Blockhash: block.Hash,
		Delegate:  "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
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
	_, _, err = r.DelegateFrozenBalance(rpc.DelegateFrozenBalanceInput{
		Blockhash: block.Hash,
		Delegate:  "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
	})
	if ok := assert.Nilf(t, err, "Failed at block '%d'", block.Header.Level); !ok {
		t.FailNow()
	}
}
