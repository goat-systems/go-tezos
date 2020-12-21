package rpc

import (
	"encoding/json"
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
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regBigMap, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get big map",
				[]byte{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regBigMap, []byte("success")}, blankHandler)),
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
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regContracts, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get contracts",
				[]string{},
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regContracts, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get contracts: failed to parse json",
				[]string{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regContracts, []byte(`[
				"tz1RSbL8J3PPcsx2W6y37tjTu6kUhsPLavid",
				"tz1W5BsELGH47PkzFPRHu1BmMCdEbWVrrCoz"
			 ]`)}, blankHandler)),
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
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regContract, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get contract 'KT1DrJV8vhkdLEj76h1H9Q4irZDqAkMPo1Qf'",
				Contract{},
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regContract, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get contract 'KT1DrJV8vhkdLEj76h1H9Q4irZDqAkMPo1Qf': failed to parse json",
				Contract{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regContract, readResponse(contract)}, blankHandler)),
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
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regContractBalance, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get balance for contract 'tz1U8sXoQWGUMQrfZeAYwAzMZUvWwy7mfpPQ'",
				"",
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regContractBalance, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get balance for contract 'tz1U8sXoQWGUMQrfZeAYwAzMZUvWwy7mfpPQ': failed to parse json",
				"",
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regContractBalance, readResponse(balance)}, blankHandler)),
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
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regContractCounter, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get counter for contract 'tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc'",
				0,
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regContractCounter, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get counter for contract 'tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc': failed to parse json",
				0,
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regContractCounter, readResponse(counter)}, blankHandler)),
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
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regContractDelegate, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get delegate for contract 'tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc'",
				"",
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regContractDelegate, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get delegate for contract 'tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc': failed to parse json",
				"",
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regContractDelegate, []byte(`"tz1WCd2jm4uSt4vntk4vSuUWoZQGhLcDuR9q"`)}, blankHandler)),
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
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regContractEntrypoints, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get entrypoints for contract 'KT1DrJV8vhkdLEj76h1H9Q4irZDqAkMPo1Qf'",
				map[string]struct{}{},
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regContractEntrypoints, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get entrypoints for contract 'KT1DrJV8vhkdLEj76h1H9Q4irZDqAkMPo1Qf': failed to parse json",
				map[string]struct{}{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regContractEntrypoints, readResponse(contractEntrypoints)}, blankHandler)),
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
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regContractEntrypoint, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get entrypoint 'xtzToToken' for contract 'KT1DrJV8vhkdLEj76h1H9Q4irZDqAkMPo1Qf'",
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regContractEntrypoints, []byte(`"tz1WCd2jm4uSt4vntk4vSuUWoZQGhLcDuR9q"`)}, blankHandler)),
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

func Test_ContractManagerKey(t *testing.T) {
	type want struct {
		err         bool
		errContains string
		manager     string
	}

	cases := []struct {
		name  string
		input http.Handler
		want  want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regContractManagerKey, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get manager for contract 'tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc'",
				"",
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regContractManagerKey, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get manager for contract 'tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc': failed to parse json",
				"",
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regContractManagerKey, []byte(`"tz1WCd2jm4uSt4vntk4vSuUWoZQGhLcDuR9q"`)}, blankHandler)),
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

			manager, err := rpc.ContractManagerKey(ContractManagerKeyInput{
				Blockhash:  mockBlockHash,
				ContractID: mockAddressTz1,
			})
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.manager, manager)
		})
	}
}

func Test_ContractScript(t *testing.T) {
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
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regContractScript, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get script for contract 'KT1DrJV8vhkdLEj76h1H9Q4irZDqAkMPo1Qf'",
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regContractScript, []byte(`"tz1WCd2jm4uSt4vntk4vSuUWoZQGhLcDuR9q"`)}, blankHandler)),
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

			_, err = rpc.ContractScript(ContractScriptInput{
				Blockhash:  mockBlockHash,
				ContractID: "KT1DrJV8vhkdLEj76h1H9Q4irZDqAkMPo1Qf",
			})
			checkErr(t, tt.want.err, tt.want.errContains, err)
		})
	}
}

//TODO: There's not a public sapling resource to test against to get mock data
// func Test_ContractSaplingDiff(t *testing.T) {
// 	type want struct {
// 		err         bool
// 		errContains string
// 		delegate    string
// 	}

// 	cases := []struct {
// 		name  string
// 		input http.Handler
// 		want  want
// 	}{
// 		{
// 			"handles rpc failure",
// 			gtGoldenHTTPMock(contractDelegateHandlerMock(readResponse(rpcerrors), blankHandler)),
// 			want{
// 				true,
// 				"failed to get delegate for contract 'tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc'",
// 				"",
// 			},
// 		},
// 		{
// 			"handles failure to unmarshal",
// 			gtGoldenHTTPMock(contractDelegateHandlerMock([]byte(`bad_data`), blankHandler)),
// 			want{
// 				true,
// 				"failed to get delegate for contract 'tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc': failed to parse json",
// 				"",
// 			},
// 		},
// 		{
// 			"is successful",
// 			gtGoldenHTTPMock(contractDelegateHandlerMock([]byte(`"tz1WCd2jm4uSt4vntk4vSuUWoZQGhLcDuR9q"`), blankHandler)),
// 			want{
// 				false,
// 				"",
// 				"tz1WCd2jm4uSt4vntk4vSuUWoZQGhLcDuR9q",
// 			},
// 		},
// 	}

// 	for _, tt := range cases {
// 		t.Run(tt.name, func(t *testing.T) {
// 			server := httptest.NewServer(tt.input)
// 			defer server.Close()

// 			rpc, err := New(server.URL)
// 			assert.Nil(t, err)

// 			saplingDiff, err := rpc.ContractSaplingDiff(ContractSaplingDiffInput{
// 				Blockhash:  mockBlockHash,
// 				ContractID: "KT1DrJV8vhkdLEj76h1H9Q4irZDqAkMPo1Qf",
// 			})
// 			checkErr(t, tt.want.err, tt.want.errContains, err)
// 			assert.Equal(t, tt.want.delegate, saplingDiff)
// 		})
// 	}
// }

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

	goldenMessage := &json.RawMessage{}
	*goldenMessage = storageJSON

	type want struct {
		err         bool
		containsErr string
		micheline   *json.RawMessage
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regContractStorage, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get storage for contract",
				nil,
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regContractStorage, storageJSON}, blankHandler)),
			want{
				false,
				"",
				goldenMessage,
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
				ContractID: "KT1LfoE9EbpdsfUzowRckGUfikGcd5PyVKg",
				Blockhash:  "BLzGD63HA4RP8Fh5xEtvdQSMKa2WzJMZjQPNVUc4Rqy8Lh5BEY1",
			})
			checkErr(t, tt.want.err, tt.containsErr, err)
			assert.Equal(t, tt.want.micheline, micheline)
		})
	}
}

func Test_Delegates(t *testing.T) {
	goldenDelegates := getResponse(delegatedcontracts).([]string)

	type want struct {
		wantErr     bool
		containsErr string
		delegates   []string
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regDelegates, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get delegates",
				[]string{},
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regDelegates, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get delegates: failed to parse json",
				[]string{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regDelegates, readResponse(delegatedcontracts)}, blankHandler)),
			want{
				false,
				"",
				goldenDelegates,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			rpc, err := New(server.URL)
			assert.Nil(t, err)

			delegates, err := rpc.Delegates(DelegatesInput{
				Blockhash: mockBlockHash,
			})
			checkErr(t, tt.wantErr, tt.containsErr, err)

			assert.Equal(t, tt.want.delegates, delegates)
		})
	}
}

func Test_Delegate(t *testing.T) {
	goldenDelegate := getResponse(delegate).(Delegate)

	type want struct {
		err         bool
		containsErr string
		delegate    Delegate
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want        want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regDelegate, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get delegate",
				Delegate{},
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regDelegate, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get delegate 'tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc': failed to parse json",
				Delegate{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regDelegate, readResponse(delegate)}, blankHandler)),
			want{
				false,
				"",
				goldenDelegate,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			rpc, err := New(server.URL)
			assert.Nil(t, err)

			delegate, err := rpc.Delegate(DelegateInput{
				Blockhash: mockBlockHash,
				Delegate:  "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
			})
			checkErr(t, tt.want.err, tt.want.containsErr, err)
			assert.Equal(t, tt.want.delegate, delegate)
		})
	}
}

func Test_DelegateBalance(t *testing.T) {
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
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regDelegateBalance, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get delegate 'tz1U8sXoQWGUMQrfZeAYwAzMZUvWwy7mfpPQ' balance",
				"",
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regDelegateBalance, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get delegate 'tz1U8sXoQWGUMQrfZeAYwAzMZUvWwy7mfpPQ' balance: failed to parse json",
				"",
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regDelegateBalance, readResponse(balance)}, blankHandler)),
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

			balance, err := rpc.DelegateBalance(DelegateBalanceInput{
				Delegate:  "tz1U8sXoQWGUMQrfZeAYwAzMZUvWwy7mfpPQ",
				Blockhash: mockBlockHash,
			})
			checkErr(t, tt.want.wantErr, tt.want.containsErr, err)
			assert.Equal(t, tt.want.balance, balance)
		})
	}
}

func Test_DelegateDeactivated(t *testing.T) {
	type want struct {
		wantErr     bool
		containsErr string
		status      bool
	}

	cases := []struct {
		name  string
		input http.Handler
		want  want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regDelegateDeactivated, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get delegate 'tz1U8sXoQWGUMQrfZeAYwAzMZUvWwy7mfpPQ' activation status",
				false,
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regDelegateDeactivated, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get delegate 'tz1U8sXoQWGUMQrfZeAYwAzMZUvWwy7mfpPQ' activation status: failed to parse json",
				false,
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regDelegateDeactivated, []byte(`true`)}, blankHandler)),
			want{
				false,
				"",
				true,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input)
			defer server.Close()

			rpc, err := New(server.URL)
			assert.Nil(t, err)

			deactivated, err := rpc.DelegateDeactivated(DelegateDeactivatedInput{
				Delegate:  "tz1U8sXoQWGUMQrfZeAYwAzMZUvWwy7mfpPQ",
				Blockhash: mockBlockHash,
			})
			checkErr(t, tt.want.wantErr, tt.want.containsErr, err)
			assert.Equal(t, tt.want.status, deactivated)
		})
	}
}

func Test_DelegateDelegatedBalance(t *testing.T) {
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
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regDelegateDelegatedBalance, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get delegate 'tz1U8sXoQWGUMQrfZeAYwAzMZUvWwy7mfpPQ' delegated balance",
				"",
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regDelegateDelegatedBalance, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get delegate 'tz1U8sXoQWGUMQrfZeAYwAzMZUvWwy7mfpPQ' delegated balance: failed to parse json",
				"",
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regDelegateDelegatedBalance, readResponse(balance)}, blankHandler)),
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

			balance, err := rpc.DelegateDelegatedBalance(DelegateBalanceInput{
				Delegate:  "tz1U8sXoQWGUMQrfZeAYwAzMZUvWwy7mfpPQ",
				Blockhash: mockBlockHash,
			})
			checkErr(t, tt.want.wantErr, tt.want.containsErr, err)
			assert.Equal(t, tt.want.balance, balance)
		})
	}
}

func Test_DelegateDelegatedContracts(t *testing.T) {
	goldenDelegations := getResponse(delegatedcontracts).([]string)

	type want struct {
		err                bool
		containsErr        string
		delegatedContracts []string
	}

	cases := []struct {
		name  string
		input http.Handler
		want  want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regDelegateDelegatedContracts, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get delegate 'tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc' delegated contracts",
				[]string{},
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regDelegateDelegatedContracts, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get delegate 'tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc' delegated contracts: failed to parse json",
				[]string{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regDelegateDelegatedContracts, readResponse(delegatedcontracts)}, blankHandler)),
			want{
				false,
				"",
				goldenDelegations,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input)
			defer server.Close()

			rpc, err := New(server.URL)
			assert.Nil(t, err)

			delegatedContracts, err := rpc.DelegateDelegatedContracts(DelegateDelegatedContractsInput{
				Blockhash: mockBlockHash,
				Delegate:  "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
			})
			checkErr(t, tt.want.err, tt.want.containsErr, err)
			assert.Equal(t, tt.want.delegatedContracts, delegatedContracts)
		})
	}
}

func Test_DelegateFrozenBalance(t *testing.T) {
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
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regDelegateFrozenBalance, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get delegate 'tz1U8sXoQWGUMQrfZeAYwAzMZUvWwy7mfpPQ' frozen balance",
				"",
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regDelegateFrozenBalance, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get delegate 'tz1U8sXoQWGUMQrfZeAYwAzMZUvWwy7mfpPQ' frozen balance: failed to parse json",
				"",
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regDelegateFrozenBalance, readResponse(balance)}, blankHandler)),
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

			balance, err := rpc.DelegateFrozenBalance(DelegateFrozenBalanceInput{
				Delegate:  "tz1U8sXoQWGUMQrfZeAYwAzMZUvWwy7mfpPQ",
				Blockhash: mockBlockHash,
			})
			checkErr(t, tt.want.wantErr, tt.want.containsErr, err)
			assert.Equal(t, tt.want.balance, balance)
		})
	}
}

func Test_DelegateFrozenBalanceBalanceAtCycle(t *testing.T) {
	goldenFrozenBalanceByCycle := getResponse(frozenbalanceByCycle).([]FrozenBalanceByCycle)

	type want struct {
		wantErr              bool
		containsErr          string
		frozenBalanceByCycle []FrozenBalanceByCycle
	}

	cases := []struct {
		name  string
		input http.Handler
		want  want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regDelegateFrozenBalanceByCycle, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get delegate 'tz1U8sXoQWGUMQrfZeAYwAzMZUvWwy7mfpPQ' frozen balance at cycle",
				[]FrozenBalanceByCycle{},
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regDelegateFrozenBalanceByCycle, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get delegate 'tz1U8sXoQWGUMQrfZeAYwAzMZUvWwy7mfpPQ' frozen balance at cycle: failed to parse json",
				[]FrozenBalanceByCycle{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regDelegateFrozenBalanceByCycle, readResponse(frozenbalanceByCycle)}, blankHandler)),
			want{
				false,
				"",
				goldenFrozenBalanceByCycle,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input)
			defer server.Close()

			rpc, err := New(server.URL)
			assert.Nil(t, err)

			frozenBalance, err := rpc.DelegateFrozenBalanceByCycle(DelegateFrozenBalanceByCycleInput{
				Delegate:  "tz1U8sXoQWGUMQrfZeAYwAzMZUvWwy7mfpPQ",
				Blockhash: mockBlockHash,
			})
			checkErr(t, tt.want.wantErr, tt.want.containsErr, err)
			assert.Equal(t, tt.want.frozenBalanceByCycle, frozenBalance)
		})
	}
}

func Test_DelegateGracePeriod(t *testing.T) {
	type want struct {
		wantErr     bool
		containsErr string
		period      int
	}

	cases := []struct {
		name  string
		input http.Handler
		want  want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regDelegateGracePeriod, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get delegate 'tz1U8sXoQWGUMQrfZeAYwAzMZUvWwy7mfpPQ' grace period",
				0,
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regDelegateGracePeriod, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get delegate 'tz1U8sXoQWGUMQrfZeAYwAzMZUvWwy7mfpPQ' grace period: failed to parse json",
				0,
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regDelegateGracePeriod, []byte(`10`)}, blankHandler)),
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

			period, err := rpc.DelegateGracePeriod(DelegateGracePeriodInput{
				Delegate:  "tz1U8sXoQWGUMQrfZeAYwAzMZUvWwy7mfpPQ",
				Blockhash: mockBlockHash,
			})
			checkErr(t, tt.want.wantErr, tt.want.containsErr, err)
			assert.Equal(t, tt.want.period, period)
		})
	}
}

func Test_DelegateStakingBalance(t *testing.T) {
	type want struct {
		wantErr        bool
		containsErr    string
		stakingBalance string
	}

	cases := []struct {
		name  string
		input http.Handler
		want  want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regDelegateStakingBalance, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get delegate 'tz1U8sXoQWGUMQrfZeAYwAzMZUvWwy7mfpPQ' staking balance",
				"",
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regDelegateStakingBalance, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get delegate 'tz1U8sXoQWGUMQrfZeAYwAzMZUvWwy7mfpPQ' staking balance: failed to parse json",
				"",
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regDelegateStakingBalance, readResponse(balance)}, blankHandler)),
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

			stakingBalance, err := rpc.DelegateStakingBalance(DelegateStakingBalanceInput{
				Delegate:  "tz1U8sXoQWGUMQrfZeAYwAzMZUvWwy7mfpPQ",
				Blockhash: mockBlockHash,
			})
			checkErr(t, tt.want.wantErr, tt.want.containsErr, err)
			assert.Equal(t, tt.want.stakingBalance, stakingBalance)
		})
	}
}

func Test_DelegateVotingPower(t *testing.T) {
	type want struct {
		wantErr     bool
		containsErr string
		period      int
	}

	cases := []struct {
		name  string
		input http.Handler
		want  want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regDelegateVotingPower, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get delegate 'tz1U8sXoQWGUMQrfZeAYwAzMZUvWwy7mfpPQ' voting power",
				0,
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regDelegateVotingPower, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get delegate 'tz1U8sXoQWGUMQrfZeAYwAzMZUvWwy7mfpPQ' voting power: failed to parse json",
				0,
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regDelegateVotingPower, []byte(`10`)}, blankHandler)),
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

			period, err := rpc.DelegateVotingPower(DelegateVotingPowerInput{
				Delegate:  "tz1U8sXoQWGUMQrfZeAYwAzMZUvWwy7mfpPQ",
				Blockhash: mockBlockHash,
			})
			checkErr(t, tt.want.wantErr, tt.want.containsErr, err)
			assert.Equal(t, tt.want.period, period)
		})
	}
}

func Test_Nonces(t *testing.T) {
	type want struct {
		wantErr     bool
		containsErr string
		nonces      Nonces
	}

	cases := []struct {
		name  string
		input http.Handler
		want  want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regNonces, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get nonces at level '1000000'",
				Nonces{},
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regNonces, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get nonces at level '1000000': failed to parse json",
				Nonces{},
			},
		},
		{
			"is successful with nonce",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regNonces, []byte(`{"nonce":"some_nonce"}`)}, blankHandler)),
			want{
				false,
				"",
				Nonces{
					Nonce: "some_nonce",
				},
			},
		},
		{
			"is successful with hash",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regNonces, []byte(`{"hash":"some_hash"}`)}, blankHandler)),
			want{
				false,
				"",
				Nonces{
					Hash: "some_hash",
				},
			},
		},
		{
			"is successful with forgotten",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regNonces, []byte(`{}`)}, blankHandler)),
			want{
				false,
				"",
				Nonces{
					Forgotten: true,
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

			nonces, err := rpc.Nonces(NoncesInput{
				Blockhash: mockBlockHash,
				Level:     1000000,
			})
			checkErr(t, tt.want.wantErr, tt.want.containsErr, err)
			assert.Equal(t, tt.want.nonces, nonces)
		})
	}
}

func Test_RawBytes(t *testing.T) {
	type want struct {
		err         bool
		containsErr string
		rawBytes    []byte
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regRawBytes, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to raw at bytes",
				[]byte{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regRawBytes, []byte(`some_raw_bytes`)}, blankHandler)),
			want{
				false,
				"",
				[]byte(`some_raw_bytes`),
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			rpc, err := New(server.URL)
			assert.Nil(t, err)

			rawBytes, err := rpc.RawBytes(RawBytesInput{
				Blockhash: "BLzGD63HA4RP8Fh5xEtvdQSMKa2WzJMZjQPNVUc4Rqy8Lh5BEY1",
			})
			checkErr(t, tt.want.err, tt.containsErr, err)
			assert.Equal(t, tt.want.rawBytes, rawBytes)
		})
	}
}

func Test_Seed(t *testing.T) {
	type want struct {
		wantErr     bool
		containsErr string
		seed        string
	}

	cases := []struct {
		name  string
		input http.Handler
		want  want
	}{
		{
			"handles rpc failure",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regSeed, readResponse(rpcerrors)}, blankHandler)),
			want{
				true,
				"failed to get seed",
				"",
			},
		},
		{
			"handles failure to unmarshal",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regSeed, []byte(`junk`)}, blankHandler)),
			want{
				true,
				"failed to get seed: failed to parse json",
				"",
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockHandler(&requestResultPair{regSeed, []byte(`"45a8ec878c9bd348359e6c84def1e98d615dfda7878706f0af8a93afbe1f3435"`)}, blankHandler)),
			want{
				false,
				"",
				"45a8ec878c9bd348359e6c84def1e98d615dfda7878706f0af8a93afbe1f3435",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input)
			defer server.Close()

			rpc, err := New(server.URL)
			assert.Nil(t, err)

			seed, err := rpc.Seed(SeedInput{
				Blockhash: mockBlockHash,
			})
			checkErr(t, tt.want.wantErr, tt.want.containsErr, err)
			assert.Equal(t, tt.want.seed, seed)
		})
	}
}
