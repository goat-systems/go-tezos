package rpc

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_BigMap(t *testing.T) {
	type want struct {
		err         bool
		containsErr string
		result      []byte
	}

	cases := []struct {
		name  string
		input http.Handler
		want  want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(bigMapHandlerMock(readResponse(rpcerrors), blankHandler)),
			want{
				true,
				"failed to get big map",
				[]byte{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(bigMapHandlerMock([]byte("success"), blankHandler)),
			want{
				false,
				"",
				[]byte("success"),
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input)
			defer server.Close()

			rpc, err := New(server.URL)
			assert.Nil(t, err)

			result, err := rpc.BigMap(BigMapInput{
				Blockhash:        mockBlockHash,
				BigMapID:         101,
				ScriptExpression: "exprupozG51AtT7yZUy5sg6VbJQ4b9omAE1PKD2PXvqi2YBuZqoKG3",
			})
			checkErr(t, tt.want.err, tt.want.containsErr, err)
			assert.Equal(t, tt.want.result, result)
		})
	}
}

func Test_Constants(t *testing.T) {
	goldenConstants := getResponse(constants).(Constants)

	type want struct {
		err         bool
		containsErr string
		constants   Constants
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(newConstantsMock().handler(readResponse(rpcerrors), blankHandler)),
			want{
				true,
				"failed to get constants",
				Constants{},
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(newConstantsMock().handler([]byte(`junk`), blankHandler)),
			want{
				true,
				"failed to get constants: failed to parse json",
				Constants{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(newConstantsMock().handler(readResponse(constants), blankHandler)),
			want{
				false,
				"",
				goldenConstants,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			rpc, err := New(server.URL)
			assert.Nil(t, err)

			constants, err := rpc.Constants(ConstantsInput{Blockhash: mockBlockHash})
			checkErr(t, tt.want.err, tt.want.containsErr, err)

			assert.Equal(t, tt.want.constants, constants)
		})
	}
}

func Test_Contracts(t *testing.T) {
	type want struct {
		err         bool
		containsErr string
		contracts   []string
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(contractsHandlerMock(readResponse(rpcerrors), blankHandler)),
			want{
				true,
				"failed to get contracts",
				[]string{},
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(contractsHandlerMock([]byte(`junk`), blankHandler)),
			want{
				true,
				"failed to get contracts: failed to parse json",
				[]string{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(contractsHandlerMock([]byte(`[
				"tz1RSbL8J3PPcsx2W6y37tjTu6kUhsPLavid",
				"tz1W5BsELGH47PkzFPRHu1BmMCdEbWVrrCoz"
			 ]`), blankHandler)),
			want{
				false,
				"",
				[]string{
					"tz1RSbL8J3PPcsx2W6y37tjTu6kUhsPLavid",
					"tz1W5BsELGH47PkzFPRHu1BmMCdEbWVrrCoz",
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			rpc, err := New(server.URL)
			assert.Nil(t, err)

			contracts, err := rpc.Contracts(ContractsInput{Blockhash: mockBlockHash})
			checkErr(t, tt.want.err, tt.want.containsErr, err)

			assert.Equal(t, tt.want.contracts, contracts)
		})
	}
}

func Test_Contract(t *testing.T) {
	goldenContract := getResponse(contract).(Contract)

	type want struct {
		err         bool
		containsErr string
		contract    Contract
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(contractHandlerMock(readResponse(rpcerrors), blankHandler)),
			want{
				true,
				"failed to get contract 'KT1DrJV8vhkdLEj76h1H9Q4irZDqAkMPo1Qf'",
				Contract{},
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(contractHandlerMock([]byte(`junk`), blankHandler)),
			want{
				true,
				"failed to get contract 'KT1DrJV8vhkdLEj76h1H9Q4irZDqAkMPo1Qf': failed to parse json",
				Contract{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(contractHandlerMock(readResponse(contract), blankHandler)),
			want{
				false,
				"",
				goldenContract,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			rpc, err := New(server.URL)
			assert.Nil(t, err)

			contract, err := rpc.Contract(ContractInput{
				Blockhash:  mockBlockHash,
				ContractID: "KT1DrJV8vhkdLEj76h1H9Q4irZDqAkMPo1Qf",
			})
			checkErr(t, tt.want.err, tt.want.containsErr, err)

			assert.Equal(t, tt.want.contract, contract)
		})
	}
}

func Test_ContractBalance(t *testing.T) {
	type want struct {
		wantErr     bool
		containsErr string
		balance     string
	}

	cases := []struct {
		name  string
		input http.Handler
		want  want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(balanceHandlerMock(readResponse(rpcerrors), blankHandler)),
			want{
				true,
				"failed to get balance for contract 'tz1U8sXoQWGUMQrfZeAYwAzMZUvWwy7mfpPQ'",
				"",
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(balanceHandlerMock([]byte(`not_balance_data`), blankHandler)),
			want{
				true,
				"failed to get balance for contract 'tz1U8sXoQWGUMQrfZeAYwAzMZUvWwy7mfpPQ': failed to parse json",
				"",
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(balanceHandlerMock(readResponse(balance), blankHandler)),
			want{
				false,
				"",
				"1216660108948",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input)
			defer server.Close()

			rpc, err := New(server.URL)
			assert.Nil(t, err)

			balance, err := rpc.ContractBalance(ContractBalanceInput{
				ContractID: "tz1U8sXoQWGUMQrfZeAYwAzMZUvWwy7mfpPQ",
				Blockhash:  mockBlockHash,
			})
			checkErr(t, tt.want.wantErr, tt.want.containsErr, err)
			assert.Equal(t, tt.want.balance, balance)
		})
	}
}

func Test_ContractCounter(t *testing.T) {
	type want struct {
		err         bool
		errContains string
		counter     int
	}

	cases := []struct {
		name  string
		input http.Handler
		want  want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(counterHandlerMock(readResponse(rpcerrors), blankHandler)),
			want{
				true,
				"failed to get counter for contract 'tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc'",
				0,
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(counterHandlerMock([]byte(`bad_data`), blankHandler)),
			want{
				true,
				"failed to get counter for contract 'tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc': failed to parse json",
				0,
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(counterHandlerMock(readResponse(counter), blankHandler)),
			want{
				false,
				"",
				10,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input)
			defer server.Close()

			rpc, err := New(server.URL)
			assert.Nil(t, err)

			counter, err := rpc.ContractCounter(ContractCounterInput{
				Blockhash:  mockBlockHash,
				ContractID: mockAddressTz1,
			})
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.counter, counter)
		})
	}
}

func Test_ContractDelegate(t *testing.T) {
	type want struct {
		err         bool
		errContains string
		delegate    string
	}

	cases := []struct {
		name  string
		input http.Handler
		want  want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(contractDelegateHandlerMock(readResponse(rpcerrors), blankHandler)),
			want{
				true,
				"failed to get delegate for contract 'tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc'",
				"",
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(contractDelegateHandlerMock([]byte(`bad_data`), blankHandler)),
			want{
				true,
				"failed to get delegate for contract 'tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc': failed to parse json",
				"",
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(contractDelegateHandlerMock([]byte(`"tz1WCd2jm4uSt4vntk4vSuUWoZQGhLcDuR9q"`), blankHandler)),
			want{
				false,
				"",
				"tz1WCd2jm4uSt4vntk4vSuUWoZQGhLcDuR9q",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input)
			defer server.Close()

			rpc, err := New(server.URL)
			assert.Nil(t, err)

			delegate, err := rpc.ContractDelegate(ContractDelegateInput{
				Blockhash:  mockBlockHash,
				ContractID: mockAddressTz1,
			})
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.delegate, delegate)
		})
	}
}

func Test_ContractEntrypoints(t *testing.T) {
	type want struct {
		err         bool
		errContains string
		keys        map[string]struct{}
	}

	cases := []struct {
		name  string
		input http.Handler
		want  want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(contractEntrypointsHandlerMock(readResponse(rpcerrors), blankHandler)),
			want{
				true,
				"failed to get entrypoints for contract 'KT1DrJV8vhkdLEj76h1H9Q4irZDqAkMPo1Qf'",
				map[string]struct{}{},
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(contractEntrypointsHandlerMock([]byte(`bad_data`), blankHandler)),
			want{
				true,
				"failed to get entrypoints for contract 'KT1DrJV8vhkdLEj76h1H9Q4irZDqAkMPo1Qf': failed to parse json",
				map[string]struct{}{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(contractEntrypointsHandlerMock(readResponse(contractEntrypoints), blankHandler)),
			want{
				false,
				"",
				map[string]struct{}{
					"xtzToToken":              {},
					"updateTokenPoolInternal": {},
					"tokenToXtz":              {},
					"setBaker":                {},
					"default":                 {},
					"addLiquidity":            {},
					"updateTokenPool":         {},
					"tokenToToken":            {},
					"setManager":              {},
					"removeLiquidity":         {},
					"approve":                 {},
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input)
			defer server.Close()

			rpc, err := New(server.URL)
			assert.Nil(t, err)

			entrypoints, err := rpc.ContractEntrypoints(ContractEntrypointsInput{
				Blockhash:  mockBlockHash,
				ContractID: "KT1DrJV8vhkdLEj76h1H9Q4irZDqAkMPo1Qf",
			})
			checkErr(t, tt.want.err, tt.want.errContains, err)

			for entrypoint := range entrypoints {
				_, ok := tt.want.keys[entrypoint]
				assert.True(t, ok)
			}
		})
	}
}

func Test_ContractEntrypoint(t *testing.T) {
	type want struct {
		err         bool
		errContains string
	}

	cases := []struct {
		name  string
		input http.Handler
		want  want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(contractEntrypointHandlerMock(readResponse(rpcerrors), blankHandler)),
			want{
				true,
				"failed to get entrypoint 'xtzToToken' for contract 'KT1DrJV8vhkdLEj76h1H9Q4irZDqAkMPo1Qf'",
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(contractEntrypointHandlerMock([]byte(`"tz1WCd2jm4uSt4vntk4vSuUWoZQGhLcDuR9q"`), blankHandler)),
			want{
				false,
				"",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input)
			defer server.Close()

			rpc, err := New(server.URL)
			assert.Nil(t, err)

			_, err = rpc.ContractEntrypoint(ContractEntrypointInput{
				Blockhash:  mockBlockHash,
				ContractID: "KT1DrJV8vhkdLEj76h1H9Q4irZDqAkMPo1Qf",
				Entrypoint: "xtzToToken",
			})
			checkErr(t, tt.want.err, tt.want.errContains, err)
		})
	}
}

func Test_ContractStorage(t *testing.T) {
	storageJSON := []byte(`[
		{
		  "prim": "parameter",
		  "args": [
			{
			  "prim": "unit",
			  "annots": [
				"%abc"
			  ]
			}
		  ]
		},
		{
		  "prim": "storage",
		  "args": [
			{
			  "prim": "unit"
			}
		  ]
		},
		{
		  "prim": "code",
		  "args": [
			[
			  {
				"prim": "CDR"
			  },
			  {
				"prim": "NIL",
				"args": [
				  {
					"prim": "operation"
				  }
				]
			  },
			  {
				"prim": "PAIR"
			  }
			]
		  ]
		}
	  ]`)

	type want struct {
		err         bool
		containsErr string
		micheline   []byte
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"returns rpc error",
			gtGoldenHTTPMock(storageHandlerMock(readResponse(rpcerrors), blankHandler)),
			want{
				true,
				"could not get storage",
				[]byte{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(storageHandlerMock(storageJSON, blankHandler)),
			want{
				false,
				"",
				storageJSON,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			rpc, err := New(server.URL)
			assert.Nil(t, err)

			micheline, err := rpc.ContractStorage(ContractStorageInput{
				Contract:  "KT1LfoE9EbpdsfUzowRckGUfikGcd5PyVKg",
				Blockhash: "BLzGD63HA4RP8Fh5xEtvdQSMKa2WzJMZjQPNVUc4Rqy8Lh5BEY1",
			})
			checkErr(t, tt.want.err, tt.containsErr, err)
			assert.Equal(t, tt.want.micheline, micheline)
		})
	}
}
