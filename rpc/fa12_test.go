package rpc_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goat-systems/go-tezos/v4/rpc"
	"github.com/stretchr/testify/assert"
)

func Test_GetFA12Balance(t *testing.T) {
	type input struct {
		hanler              http.Handler
		getFA12BalanceInput rpc.GetFA12BalanceInput
	}

	type want struct {
		err      bool
		contains string
		balance  string
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"handles failure to validate input",
			input{
				gtGoldenHTTPMock(blankHandler),
				rpc.GetFA12BalanceInput{},
			},
			want{
				true,
				"invalid input",
				"0",
			},
		},
		{
			"handles failure to get cycle",
			input{
				gtGoldenHTTPMock(mockCycleFailed(blankHandler)),
				rpc.GetFA12BalanceInput{
					Cycle:        1,
					ChainID:      "some_chainid",
					Source:       "some_source",
					FA12Contract: "some_fa1.2_contract",
					OwnerAddress: "some_address",
				},
			},
			want{
				true,
				"could not get cycle",
				"0",
			},
		},
		{
			"handles failure to get counter",
			input{
				gtGoldenHTTPMock(mockCycleSuccessful(mockHandler(&requestResultPair{regContractCounter, []byte(`junk`)}, blankHandler))),
				rpc.GetFA12BalanceInput{
					Cycle:        1,
					ChainID:      "some_chainid",
					Source:       "some_source",
					FA12Contract: "some_fa1.2_contract",
					OwnerAddress: "some_address",
				},
			},
			want{
				true,
				"failed to unmarshal counter",
				"0",
			},
		},
		{
			"handles failure to run_operation",
			input{
				gtGoldenHTTPMock(mockCycleSuccessful(mockHandler(
					&requestResultPair{regContractCounter, []byte(`"100"`)},
					runOperationHandlerMock(
						[]byte(`junk`),
						blankHandler,
					),
				))),
				rpc.GetFA12BalanceInput{
					Cycle:        1,
					ChainID:      "some_chainid",
					Source:       "some_source",
					FA12Contract: "some_fa1.2_contract",
					OwnerAddress: "some_address",
					Testnet:      true,
				},
			},
			want{
				true,
				"failed to unmarshal operation",
				"0",
			},
		},
		{
			"handles failure to parse balance",
			input{
				gtGoldenHTTPMock(mockCycleSuccessful(mockHandler(
					&requestResultPair{regContractCounter, []byte(`"100"`)},
					runOperationHandlerMock(
						[]byte(`{"contents":[{"kind":"transaction","source":"tz1S82rGFZK8cVbNDpP1Hf9VhTUa4W8oc2WV","fee":"0","counter":"553001","gas_limit":"1040000","storage_limit":"60000","amount":"0","destination":"KT1Njyz94x2pNJGh5uMhKj24VB9JsGCdkySN","parameters":{"entrypoint":"default","value":[{"prim":"PUSH","args":[{"prim":"mutez"},{"int":"0"}]},{"prim":"NONE","args":[{"prim":"key_hash"}]},{"prim":"CREATE_CONTRACT","args":[[{"prim":"parameter","args":[{"prim":"nat"}]},{"prim":"storage","args":[{"prim":"unit"}]},{"prim":"code","args":[[{"prim":"FAILWITH"}]]}]]},{"prim":"DIP","args":[[{"prim":"DIP","args":[[{"prim":"LAMBDA","args":[{"prim":"pair","args":[{"prim":"address"},{"prim":"unit"}]},{"prim":"pair","args":[{"prim":"list","args":[{"prim":"operation"}]},{"prim":"unit"}]},[{"prim":"CAR"},{"prim":"CONTRACT","args":[{"prim":"nat"}]},{"prim":"IF_NONE","args":[[{"prim":"PUSH","args":[{"prim":"string"},{"string":"a"}]},{"prim":"FAILWITH"}],[]]},{"prim":"PUSH","args":[{"prim":"address"},{"string":"tz1MQehPikysuVYN5hTiKTnrsFidAww7rv3z"}]},{"prim":"PAIR"},{"prim":"DIP","args":[[{"prim":"PUSH","args":[{"prim":"address"},{"string":"KT1U8kHUKJVcPkzjHLahLxmnXGnK1StAeNjK"}]},{"prim":"CONTRACT","args":[{"prim":"pair","args":[{"prim":"address"},{"prim":"contract","args":[{"prim":"nat"}]}]}],"annots":["%getBalance"]},{"prim":"IF_NONE","args":[[{"prim":"PUSH","args":[{"prim":"string"},{"string":"b"}]},{"prim":"FAILWITH"}],[]]},{"prim":"PUSH","args":[{"prim":"mutez"},{"int":"0"}]}]]},{"prim":"TRANSFER_TOKENS"},{"prim":"DIP","args":[[{"prim":"NIL","args":[{"prim":"operation"}]}]]},{"prim":"CONS"},{"prim":"DIP","args":[[{"prim":"UNIT"}]]},{"prim":"PAIR"}]]}]]},{"prim":"APPLY"},{"prim":"DIP","args":[[{"prim":"PUSH","args":[{"prim":"address"},{"string":"KT1Njyz94x2pNJGh5uMhKj24VB9JsGCdkySN"}]},{"prim":"CONTRACT","args":[{"prim":"lambda","args":[{"prim":"unit"},{"prim":"pair","args":[{"prim":"list","args":[{"prim":"operation"}]},{"prim":"unit"}]}]}]},{"prim":"IF_NONE","args":[[{"prim":"PUSH","args":[{"prim":"string"},{"string":"c"}]},{"prim":"FAILWITH"}],[]]},{"prim":"PUSH","args":[{"prim":"mutez"},{"int":"0"}]}]]},{"prim":"TRANSFER_TOKENS"},{"prim":"DIP","args":[[{"prim":"NIL","args":[{"prim":"operation"}]}]]},{"prim":"CONS"}]]},{"prim":"CONS"},{"prim":"DIP","args":[[{"prim":"UNIT"}]]},{"prim":"PAIR"}]},"metadata":{"balance_updates":[],"operation_result":{"status":"backtracked","storage":{"prim":"Unit"},"consumed_gas":"26984","storage_size":"46"},"internal_operation_results":[{"kind":"origination","source":"KT1Njyz94x2pNJGh5uMhKj24VB9JsGCdkySN","nonce":0,"balance":"0","script":{"code":[{"prim":"parameter","args":[{"prim":"nat"}]},{"prim":"storage","args":[{"prim":"unit"}]},{"prim":"code","args":[[{"prim":"FAILWITH"}]]}],"storage":{"prim":"Unit"}},"result":{"status":"backtracked","big_map_diff":[],"balance_updates":[{"kind":"contract","contract":"tz1S82rGFZK8cVbNDpP1Hf9VhTUa4W8oc2WV","change":"-32000"},{"kind":"contract","contract":"tz1S82rGFZK8cVbNDpP1Hf9VhTUa4W8oc2WV","change":"-257000"}],"originated_contracts":["KT1TTaLErwMQAVRNB1sVXf9NUdXHLCrnNpUV"],"consumed_gas":"10696","storage_size":"32","paid_storage_size_diff":"32"}},{"kind":"transaction","source":"KT1Njyz94x2pNJGh5uMhKj24VB9JsGCdkySN","nonce":1,"amount":"0","destination":"KT1Njyz94x2pNJGh5uMhKj24VB9JsGCdkySN","parameters":{"entrypoint":"default","value":[{"prim":"PUSH","args":[{"prim":"address"},{"bytes":"01cf0e19e55c5c34ec644b3b1c46c5fe3d8feb96c600"}]},{"prim":"PAIR"},[{"prim":"CAR"},{"prim":"CONTRACT","args":[{"prim":"nat"}]},{"prim":"IF_NONE","args":[[{"prim":"PUSH","args":[{"prim":"string"},{"string":"a"}]},{"prim":"FAILWITH"}],[]]},{"prim":"PUSH","args":[{"prim":"address"},{"bytes":"000013687bbb1298ad36a4bdb2c5da2f126a59aac007"}]},{"prim":"PAIR"},{"prim":"DIP","args":[[{"prim":"PUSH","args":[{"prim":"address"},{"bytes":"01d676a1a3e2e602bbb478bb188265fdec8f09124d00"}]},{"prim":"CONTRACT","args":[{"prim":"pair","args":[{"prim":"address"},{"prim":"contract","args":[{"prim":"nat"}]}]}],"annots":["%getBalance"]},{"prim":"IF_NONE","args":[[{"prim":"PUSH","args":[{"prim":"string"},{"string":"b"}]},{"prim":"FAILWITH"}],[]]},{"prim":"PUSH","args":[{"prim":"mutez"},{"int":"0"}]}]]},{"prim":"TRANSFER_TOKENS"},{"prim":"DIP","args":[[{"prim":"NIL","args":[{"prim":"operation"}]}]]},{"prim":"CONS"},{"prim":"DIP","args":[[{"prim":"UNIT"}]]},{"prim":"PAIR"}]]},"result":{"status":"backtracked","storage":{"prim":"Unit"},"consumed_gas":"77868","storage_size":"46"}},{"kind":"transaction","source":"KT1Njyz94x2pNJGh5uMhKj24VB9JsGCdkySN","nonce":2,"amount":"0","destination":"KT1U8kHUKJVcPkzjHLahLxmnXGnK1StAeNjK","parameters":{"entrypoint":"getBalance","value":{"prim":"Pair","args":[{"bytes":"000013687bbb1298ad36a4bdb2c5da2f126a59aac007"},{"bytes":"01cf0e19e55c5c34ec644b3b1c46c5fe3d8feb96c600"}]}},"result":{"status":"backtracked","storage":{"prim":"Pair","args":[{"int":"5118"},{"int":"1670000"}]},"consumed_gas":"95277","storage_size":"3007"}},{"kind":"transaction","source":"KT1U8kHUKJVcPkzjHLahLxmnXGnK1StAeNjK","nonce":3,"amount":"0","destination":"KT1TTaLErwMQAVRNB1sVXf9NUdXHLCrnNpUV","parameters":{"entrypoint":"default","value":{"int":"1546544"}},"result":{"status":"failed"}}]}}]}`),
						blankHandler,
					),
				))),
				rpc.GetFA12BalanceInput{
					Cycle:        1,
					ChainID:      "some_chainid",
					Source:       "some_source",
					FA12Contract: "some_fa1.2_contract",
					OwnerAddress: "some_address",
					Testnet:      true,
				},
			},
			want{
				true,
				"failed to parse balance",
				"0",
			},
		},
		{
			"is successful",
			input{
				gtGoldenHTTPMock(mockCycleSuccessful(mockHandler(
					&requestResultPair{regContractCounter, []byte(`"100"`)},
					runOperationHandlerMock(
						[]byte(`{"contents":[{"kind":"transaction","source":"tz1S82rGFZK8cVbNDpP1Hf9VhTUa4W8oc2WV","fee":"0","counter":"553001","gas_limit":"1040000","storage_limit":"60000","amount":"0","destination":"KT1Njyz94x2pNJGh5uMhKj24VB9JsGCdkySN","parameters":{"entrypoint":"default","value":[{"prim":"PUSH","args":[{"prim":"mutez"},{"int":"0"}]},{"prim":"NONE","args":[{"prim":"key_hash"}]},{"prim":"CREATE_CONTRACT","args":[[{"prim":"parameter","args":[{"prim":"nat"}]},{"prim":"storage","args":[{"prim":"unit"}]},{"prim":"code","args":[[{"prim":"FAILWITH"}]]}]]},{"prim":"DIP","args":[[{"prim":"DIP","args":[[{"prim":"LAMBDA","args":[{"prim":"pair","args":[{"prim":"address"},{"prim":"unit"}]},{"prim":"pair","args":[{"prim":"list","args":[{"prim":"operation"}]},{"prim":"unit"}]},[{"prim":"CAR"},{"prim":"CONTRACT","args":[{"prim":"nat"}]},{"prim":"IF_NONE","args":[[{"prim":"PUSH","args":[{"prim":"string"},{"string":"a"}]},{"prim":"FAILWITH"}],[]]},{"prim":"PUSH","args":[{"prim":"address"},{"string":"tz1MQehPikysuVYN5hTiKTnrsFidAww7rv3z"}]},{"prim":"PAIR"},{"prim":"DIP","args":[[{"prim":"PUSH","args":[{"prim":"address"},{"string":"KT1U8kHUKJVcPkzjHLahLxmnXGnK1StAeNjK"}]},{"prim":"CONTRACT","args":[{"prim":"pair","args":[{"prim":"address"},{"prim":"contract","args":[{"prim":"nat"}]}]}],"annots":["%getBalance"]},{"prim":"IF_NONE","args":[[{"prim":"PUSH","args":[{"prim":"string"},{"string":"b"}]},{"prim":"FAILWITH"}],[]]},{"prim":"PUSH","args":[{"prim":"mutez"},{"int":"0"}]}]]},{"prim":"TRANSFER_TOKENS"},{"prim":"DIP","args":[[{"prim":"NIL","args":[{"prim":"operation"}]}]]},{"prim":"CONS"},{"prim":"DIP","args":[[{"prim":"UNIT"}]]},{"prim":"PAIR"}]]}]]},{"prim":"APPLY"},{"prim":"DIP","args":[[{"prim":"PUSH","args":[{"prim":"address"},{"string":"KT1Njyz94x2pNJGh5uMhKj24VB9JsGCdkySN"}]},{"prim":"CONTRACT","args":[{"prim":"lambda","args":[{"prim":"unit"},{"prim":"pair","args":[{"prim":"list","args":[{"prim":"operation"}]},{"prim":"unit"}]}]}]},{"prim":"IF_NONE","args":[[{"prim":"PUSH","args":[{"prim":"string"},{"string":"c"}]},{"prim":"FAILWITH"}],[]]},{"prim":"PUSH","args":[{"prim":"mutez"},{"int":"0"}]}]]},{"prim":"TRANSFER_TOKENS"},{"prim":"DIP","args":[[{"prim":"NIL","args":[{"prim":"operation"}]}]]},{"prim":"CONS"}]]},{"prim":"CONS"},{"prim":"DIP","args":[[{"prim":"UNIT"}]]},{"prim":"PAIR"}]},"metadata":{"balance_updates":[],"operation_result":{"status":"backtracked","storage":{"prim":"Unit"},"consumed_gas":"26984","storage_size":"46"},"internal_operation_results":[{"kind":"origination","source":"KT1Njyz94x2pNJGh5uMhKj24VB9JsGCdkySN","nonce":0,"balance":"0","script":{"code":[{"prim":"parameter","args":[{"prim":"nat"}]},{"prim":"storage","args":[{"prim":"unit"}]},{"prim":"code","args":[[{"prim":"FAILWITH"}]]}],"storage":{"prim":"Unit"}},"result":{"status":"backtracked","big_map_diff":[],"balance_updates":[{"kind":"contract","contract":"tz1S82rGFZK8cVbNDpP1Hf9VhTUa4W8oc2WV","change":"-32000"},{"kind":"contract","contract":"tz1S82rGFZK8cVbNDpP1Hf9VhTUa4W8oc2WV","change":"-257000"}],"originated_contracts":["KT1TTaLErwMQAVRNB1sVXf9NUdXHLCrnNpUV"],"consumed_gas":"10696","storage_size":"32","paid_storage_size_diff":"32"}},{"kind":"transaction","source":"KT1Njyz94x2pNJGh5uMhKj24VB9JsGCdkySN","nonce":1,"amount":"0","destination":"KT1Njyz94x2pNJGh5uMhKj24VB9JsGCdkySN","parameters":{"entrypoint":"default","value":[{"prim":"PUSH","args":[{"prim":"address"},{"bytes":"01cf0e19e55c5c34ec644b3b1c46c5fe3d8feb96c600"}]},{"prim":"PAIR"},[{"prim":"CAR"},{"prim":"CONTRACT","args":[{"prim":"nat"}]},{"prim":"IF_NONE","args":[[{"prim":"PUSH","args":[{"prim":"string"},{"string":"a"}]},{"prim":"FAILWITH"}],[]]},{"prim":"PUSH","args":[{"prim":"address"},{"bytes":"000013687bbb1298ad36a4bdb2c5da2f126a59aac007"}]},{"prim":"PAIR"},{"prim":"DIP","args":[[{"prim":"PUSH","args":[{"prim":"address"},{"bytes":"01d676a1a3e2e602bbb478bb188265fdec8f09124d00"}]},{"prim":"CONTRACT","args":[{"prim":"pair","args":[{"prim":"address"},{"prim":"contract","args":[{"prim":"nat"}]}]}],"annots":["%getBalance"]},{"prim":"IF_NONE","args":[[{"prim":"PUSH","args":[{"prim":"string"},{"string":"b"}]},{"prim":"FAILWITH"}],[]]},{"prim":"PUSH","args":[{"prim":"mutez"},{"int":"0"}]}]]},{"prim":"TRANSFER_TOKENS"},{"prim":"DIP","args":[[{"prim":"NIL","args":[{"prim":"operation"}]}]]},{"prim":"CONS"},{"prim":"DIP","args":[[{"prim":"UNIT"}]]},{"prim":"PAIR"}]]},"result":{"status":"backtracked","storage":{"prim":"Unit"},"consumed_gas":"77868","storage_size":"46"}},{"kind":"transaction","source":"KT1Njyz94x2pNJGh5uMhKj24VB9JsGCdkySN","nonce":2,"amount":"0","destination":"KT1U8kHUKJVcPkzjHLahLxmnXGnK1StAeNjK","parameters":{"entrypoint":"getBalance","value":{"prim":"Pair","args":[{"bytes":"000013687bbb1298ad36a4bdb2c5da2f126a59aac007"},{"bytes":"01cf0e19e55c5c34ec644b3b1c46c5fe3d8feb96c600"}]}},"result":{"status":"backtracked","storage":{"prim":"Pair","args":[{"int":"5118"},{"int":"1670000"}]},"consumed_gas":"95277","storage_size":"3007"}},{"kind":"transaction","source":"KT1U8kHUKJVcPkzjHLahLxmnXGnK1StAeNjK","nonce":3,"amount":"0","destination":"KT1TTaLErwMQAVRNB1sVXf9NUdXHLCrnNpUV","parameters":{"entrypoint":"default","value":{"int":"1546544"}},"result":{"status":"failed","errors":[{"kind":"temporary","id":"proto.006-PsCARTHA.michelson_v1.runtime_error","contract_handle":"KT1TTaLErwMQAVRNB1sVXf9NUdXHLCrnNpUV","contract_code":[{"prim":"parameter","args":[{"prim":"nat"}]},{"prim":"storage","args":[{"prim":"unit"}]},{"prim":"code","args":[[{"prim":"FAILWITH"}]]}]},{"kind":"temporary","id":"proto.006-PsCARTHA.michelson_v1.script_rejected","location":7,"with":{"prim":"Pair","args":[{"int":"1546544"},{"prim":"Unit"}]}}]}}]}}]}`),
						blankHandler,
					),
				))),
				rpc.GetFA12BalanceInput{
					Cycle:        1,
					ChainID:      "some_chainid",
					Source:       "some_source",
					FA12Contract: "some_fa1.2_contract",
					OwnerAddress: "some_address",
					Testnet:      true,
				},
			},
			want{
				false,
				"",
				"1546544",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input.hanler)
			defer server.Close()

			r, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, balance, err := r.GetFA12Balance(tt.input.getFA12BalanceInput)
			checkErr(t, tt.want.err, tt.want.contains, err)
			assert.Equal(t, tt.want.balance, balance)
		})
	}
}

func Test_GetFA12Supply(t *testing.T) {
	type input struct {
		hanler             http.Handler
		getFA12SupplyInput rpc.GetFA12SupplyInput
	}

	type want struct {
		err      bool
		contains string
		balance  string
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"handles failure to validate input",
			input{
				gtGoldenHTTPMock(blankHandler),
				rpc.GetFA12SupplyInput{},
			},
			want{
				true,
				"invalid input",
				"0",
			},
		},
		{
			"handles failure to get cycle",
			input{
				gtGoldenHTTPMock(mockCycleFailed(blankHandler)),
				rpc.GetFA12SupplyInput{
					Cycle:        10,
					Source:       "some_source",
					FA12Contract: "some_fa1.2_contract",
					ChainID:      "some_chainid",
				},
			},
			want{
				true,
				"could not get cycle",
				"0",
			},
		},
		{
			"handles failure to get counter",
			input{
				gtGoldenHTTPMock(mockCycleSuccessful(mockHandler(
					&requestResultPair{regContractCounter, []byte(`junk`)},
					blankHandler,
				))),
				rpc.GetFA12SupplyInput{
					Source:       "some_source",
					FA12Contract: "some_fa1.2_contract",
					ChainID:      "some_chainid",
				},
			},
			want{
				true,
				"could not get fa1.2 supply for contract",
				"0",
			},
		},
		{
			"handles failure to run_operation",
			input{
				gtGoldenHTTPMock(mockCycleSuccessful(mockHandler(
					&requestResultPair{regContractCounter, []byte(`"100"`)},
					runOperationHandlerMock(
						[]byte(`junk`),
						blankHandler,
					),
				))),
				rpc.GetFA12SupplyInput{
					Cycle:        10,
					Source:       "some_source",
					FA12Contract: "some_fa1.2_contract",
					Testnet:      true,
					ChainID:      "some_chainid",
				},
			},
			want{
				true,
				"failed to unmarshal operation",
				"0",
			},
		},
		{
			"handles failure to parse supply",
			input{
				gtGoldenHTTPMock(mockCycleSuccessful(mockHandler(
					&requestResultPair{regContractCounter, []byte(`"100"`)},
					runOperationHandlerMock(
						[]byte(`{"contents":[{"kind":"transaction","source":"tz1S82rGFZK8cVbNDpP1Hf9VhTUa4W8oc2WV","fee":"0","counter":"553001","gas_limit":"1040000","storage_limit":"60000","amount":"0","destination":"KT1Njyz94x2pNJGh5uMhKj24VB9JsGCdkySN","parameters":{"entrypoint":"default","value":[{"prim":"PUSH","args":[{"prim":"mutez"},{"int":"0"}]},{"prim":"NONE","args":[{"prim":"key_hash"}]},{"prim":"CREATE_CONTRACT","args":[[{"prim":"parameter","args":[{"prim":"nat"}]},{"prim":"storage","args":[{"prim":"unit"}]},{"prim":"code","args":[[{"prim":"FAILWITH"}]]}]]},{"prim":"DIP","args":[[{"prim":"DIP","args":[[{"prim":"LAMBDA","args":[{"prim":"pair","args":[{"prim":"address"},{"prim":"unit"}]},{"prim":"pair","args":[{"prim":"list","args":[{"prim":"operation"}]},{"prim":"unit"}]},[{"prim":"CAR"},{"prim":"CONTRACT","args":[{"prim":"nat"}]},{"prim":"IF_NONE","args":[[{"prim":"PUSH","args":[{"prim":"string"},{"string":"a"}]},{"prim":"FAILWITH"}],[]]},{"prim":"PUSH","args":[{"prim":"address"},{"string":"tz1MQehPikysuVYN5hTiKTnrsFidAww7rv3z"}]},{"prim":"PAIR"},{"prim":"DIP","args":[[{"prim":"PUSH","args":[{"prim":"address"},{"string":"KT1U8kHUKJVcPkzjHLahLxmnXGnK1StAeNjK"}]},{"prim":"CONTRACT","args":[{"prim":"pair","args":[{"prim":"address"},{"prim":"contract","args":[{"prim":"nat"}]}]}],"annots":["%getBalance"]},{"prim":"IF_NONE","args":[[{"prim":"PUSH","args":[{"prim":"string"},{"string":"b"}]},{"prim":"FAILWITH"}],[]]},{"prim":"PUSH","args":[{"prim":"mutez"},{"int":"0"}]}]]},{"prim":"TRANSFER_TOKENS"},{"prim":"DIP","args":[[{"prim":"NIL","args":[{"prim":"operation"}]}]]},{"prim":"CONS"},{"prim":"DIP","args":[[{"prim":"UNIT"}]]},{"prim":"PAIR"}]]}]]},{"prim":"APPLY"},{"prim":"DIP","args":[[{"prim":"PUSH","args":[{"prim":"address"},{"string":"KT1Njyz94x2pNJGh5uMhKj24VB9JsGCdkySN"}]},{"prim":"CONTRACT","args":[{"prim":"lambda","args":[{"prim":"unit"},{"prim":"pair","args":[{"prim":"list","args":[{"prim":"operation"}]},{"prim":"unit"}]}]}]},{"prim":"IF_NONE","args":[[{"prim":"PUSH","args":[{"prim":"string"},{"string":"c"}]},{"prim":"FAILWITH"}],[]]},{"prim":"PUSH","args":[{"prim":"mutez"},{"int":"0"}]}]]},{"prim":"TRANSFER_TOKENS"},{"prim":"DIP","args":[[{"prim":"NIL","args":[{"prim":"operation"}]}]]},{"prim":"CONS"}]]},{"prim":"CONS"},{"prim":"DIP","args":[[{"prim":"UNIT"}]]},{"prim":"PAIR"}]},"metadata":{"balance_updates":[],"operation_result":{"status":"backtracked","storage":{"prim":"Unit"},"consumed_gas":"26984","storage_size":"46"},"internal_operation_results":[{"kind":"origination","source":"KT1Njyz94x2pNJGh5uMhKj24VB9JsGCdkySN","nonce":0,"balance":"0","script":{"code":[{"prim":"parameter","args":[{"prim":"nat"}]},{"prim":"storage","args":[{"prim":"unit"}]},{"prim":"code","args":[[{"prim":"FAILWITH"}]]}],"storage":{"prim":"Unit"}},"result":{"status":"backtracked","big_map_diff":[],"balance_updates":[{"kind":"contract","contract":"tz1S82rGFZK8cVbNDpP1Hf9VhTUa4W8oc2WV","change":"-32000"},{"kind":"contract","contract":"tz1S82rGFZK8cVbNDpP1Hf9VhTUa4W8oc2WV","change":"-257000"}],"originated_contracts":["KT1TTaLErwMQAVRNB1sVXf9NUdXHLCrnNpUV"],"consumed_gas":"10696","storage_size":"32","paid_storage_size_diff":"32"}},{"kind":"transaction","source":"KT1Njyz94x2pNJGh5uMhKj24VB9JsGCdkySN","nonce":1,"amount":"0","destination":"KT1Njyz94x2pNJGh5uMhKj24VB9JsGCdkySN","parameters":{"entrypoint":"default","value":[{"prim":"PUSH","args":[{"prim":"address"},{"bytes":"01cf0e19e55c5c34ec644b3b1c46c5fe3d8feb96c600"}]},{"prim":"PAIR"},[{"prim":"CAR"},{"prim":"CONTRACT","args":[{"prim":"nat"}]},{"prim":"IF_NONE","args":[[{"prim":"PUSH","args":[{"prim":"string"},{"string":"a"}]},{"prim":"FAILWITH"}],[]]},{"prim":"PUSH","args":[{"prim":"address"},{"bytes":"000013687bbb1298ad36a4bdb2c5da2f126a59aac007"}]},{"prim":"PAIR"},{"prim":"DIP","args":[[{"prim":"PUSH","args":[{"prim":"address"},{"bytes":"01d676a1a3e2e602bbb478bb188265fdec8f09124d00"}]},{"prim":"CONTRACT","args":[{"prim":"pair","args":[{"prim":"address"},{"prim":"contract","args":[{"prim":"nat"}]}]}],"annots":["%getBalance"]},{"prim":"IF_NONE","args":[[{"prim":"PUSH","args":[{"prim":"string"},{"string":"b"}]},{"prim":"FAILWITH"}],[]]},{"prim":"PUSH","args":[{"prim":"mutez"},{"int":"0"}]}]]},{"prim":"TRANSFER_TOKENS"},{"prim":"DIP","args":[[{"prim":"NIL","args":[{"prim":"operation"}]}]]},{"prim":"CONS"},{"prim":"DIP","args":[[{"prim":"UNIT"}]]},{"prim":"PAIR"}]]},"result":{"status":"backtracked","storage":{"prim":"Unit"},"consumed_gas":"77868","storage_size":"46"}},{"kind":"transaction","source":"KT1Njyz94x2pNJGh5uMhKj24VB9JsGCdkySN","nonce":2,"amount":"0","destination":"KT1U8kHUKJVcPkzjHLahLxmnXGnK1StAeNjK","parameters":{"entrypoint":"getBalance","value":{"prim":"Pair","args":[{"bytes":"000013687bbb1298ad36a4bdb2c5da2f126a59aac007"},{"bytes":"01cf0e19e55c5c34ec644b3b1c46c5fe3d8feb96c600"}]}},"result":{"status":"backtracked","storage":{"prim":"Pair","args":[{"int":"5118"},{"int":"1670000"}]},"consumed_gas":"95277","storage_size":"3007"}},{"kind":"transaction","source":"KT1U8kHUKJVcPkzjHLahLxmnXGnK1StAeNjK","nonce":3,"amount":"0","destination":"KT1TTaLErwMQAVRNB1sVXf9NUdXHLCrnNpUV","parameters":{"entrypoint":"default","value":{"int":"1546544"}},"result":{"status":"failed"}}]}}]}`),
						blankHandler,
					),
				))),
				rpc.GetFA12SupplyInput{
					Cycle:        10,
					Source:       "some_source",
					FA12Contract: "some_fa1.2_contract",
					Testnet:      true,
					ChainID:      "some_chainid",
				},
			},
			want{
				true,
				"failed to parse supply",
				"0",
			},
		},
		{
			"is successful",
			input{
				gtGoldenHTTPMock(mockCycleSuccessful(mockHandler(
					&requestResultPair{regContractCounter, []byte(`"100"`)},
					runOperationHandlerMock(
						[]byte(`{"contents":[{"kind":"transaction","source":"tz1S82rGFZK8cVbNDpP1Hf9VhTUa4W8oc2WV","fee":"0","counter":"553001","gas_limit":"1040000","storage_limit":"60000","amount":"0","destination":"KT1Njyz94x2pNJGh5uMhKj24VB9JsGCdkySN","parameters":{"entrypoint":"default","value":[{"prim":"PUSH","args":[{"prim":"mutez"},{"int":"0"}]},{"prim":"NONE","args":[{"prim":"key_hash"}]},{"prim":"CREATE_CONTRACT","args":[[{"prim":"parameter","args":[{"prim":"nat"}]},{"prim":"storage","args":[{"prim":"unit"}]},{"prim":"code","args":[[{"prim":"FAILWITH"}]]}]]},{"prim":"DIP","args":[[{"prim":"DIP","args":[[{"prim":"LAMBDA","args":[{"prim":"pair","args":[{"prim":"address"},{"prim":"unit"}]},{"prim":"pair","args":[{"prim":"list","args":[{"prim":"operation"}]},{"prim":"unit"}]},[{"prim":"CAR"},{"prim":"CONTRACT","args":[{"prim":"nat"}]},{"prim":"IF_NONE","args":[[{"prim":"PUSH","args":[{"prim":"string"},{"string":"a"}]},{"prim":"FAILWITH"}],[]]},{"prim":"PUSH","args":[{"prim":"unit"},{"prim":"Unit"}]},{"prim":"PAIR"},{"prim":"DIP","args":[[{"prim":"PUSH","args":[{"prim":"address"},{"string":"KT1U8kHUKJVcPkzjHLahLxmnXGnK1StAeNjK"}]},{"prim":"CONTRACT","args":[{"prim":"pair","args":[{"prim":"unit"},{"prim":"contract","args":[{"prim":"nat"}]}]}],"annots":["%getTotalSupply"]},{"prim":"IF_NONE","args":[[{"prim":"PUSH","args":[{"prim":"string"},{"string":"b"}]},{"prim":"FAILWITH"}],[]]},{"prim":"PUSH","args":[{"prim":"mutez"},{"int":"0"}]}]]},{"prim":"TRANSFER_TOKENS"},{"prim":"DIP","args":[[{"prim":"NIL","args":[{"prim":"operation"}]}]]},{"prim":"CONS"},{"prim":"DIP","args":[[{"prim":"UNIT"}]]},{"prim":"PAIR"}]]}]]},{"prim":"APPLY"},{"prim":"DIP","args":[[{"prim":"PUSH","args":[{"prim":"address"},{"string":"KT1Njyz94x2pNJGh5uMhKj24VB9JsGCdkySN"}]},{"prim":"CONTRACT","args":[{"prim":"lambda","args":[{"prim":"unit"},{"prim":"pair","args":[{"prim":"list","args":[{"prim":"operation"}]},{"prim":"unit"}]}]}]},{"prim":"IF_NONE","args":[[{"prim":"PUSH","args":[{"prim":"string"},{"string":"c"}]},{"prim":"FAILWITH"}],[]]},{"prim":"PUSH","args":[{"prim":"mutez"},{"int":"0"}]}]]},{"prim":"TRANSFER_TOKENS"},{"prim":"DIP","args":[[{"prim":"NIL","args":[{"prim":"operation"}]}]]},{"prim":"CONS"}]]},{"prim":"CONS"},{"prim":"DIP","args":[[{"prim":"UNIT"}]]},{"prim":"PAIR"}]},"metadata":{"balance_updates":[],"operation_result":{"status":"backtracked","storage":{"prim":"Unit"},"consumed_gas":"26958","storage_size":"46"},"internal_operation_results":[{"kind":"origination","source":"KT1Njyz94x2pNJGh5uMhKj24VB9JsGCdkySN","nonce":0,"balance":"0","script":{"code":[{"prim":"parameter","args":[{"prim":"nat"}]},{"prim":"storage","args":[{"prim":"unit"}]},{"prim":"code","args":[[{"prim":"FAILWITH"}]]}],"storage":{"prim":"Unit"}},"result":{"status":"backtracked","big_map_diff":[],"balance_updates":[{"kind":"contract","contract":"tz1S82rGFZK8cVbNDpP1Hf9VhTUa4W8oc2WV","change":"-32000"},{"kind":"contract","contract":"tz1S82rGFZK8cVbNDpP1Hf9VhTUa4W8oc2WV","change":"-257000"}],"originated_contracts":["KT1RHhgQZtcbYK8cVvCBymUPqTmKS5PxQbPf"],"consumed_gas":"10696","storage_size":"32","paid_storage_size_diff":"32"}},{"kind":"transaction","source":"KT1Njyz94x2pNJGh5uMhKj24VB9JsGCdkySN","nonce":1,"amount":"0","destination":"KT1Njyz94x2pNJGh5uMhKj24VB9JsGCdkySN","parameters":{"entrypoint":"default","value":[{"prim":"PUSH","args":[{"prim":"address"},{"bytes":"01b73fd2d304bc41891cee0388f1a9e59706952cbd00"}]},{"prim":"PAIR"},[{"prim":"CAR"},{"prim":"CONTRACT","args":[{"prim":"nat"}]},{"prim":"IF_NONE","args":[[{"prim":"PUSH","args":[{"prim":"string"},{"string":"a"}]},{"prim":"FAILWITH"}],[]]},{"prim":"PUSH","args":[{"prim":"unit"},{"prim":"Unit"}]},{"prim":"PAIR"},{"prim":"DIP","args":[[{"prim":"PUSH","args":[{"prim":"address"},{"bytes":"01d676a1a3e2e602bbb478bb188265fdec8f09124d00"}]},{"prim":"CONTRACT","args":[{"prim":"pair","args":[{"prim":"unit"},{"prim":"contract","args":[{"prim":"nat"}]}]}],"annots":["%getTotalSupply"]},{"prim":"IF_NONE","args":[[{"prim":"PUSH","args":[{"prim":"string"},{"string":"b"}]},{"prim":"FAILWITH"}],[]]},{"prim":"PUSH","args":[{"prim":"mutez"},{"int":"0"}]}]]},{"prim":"TRANSFER_TOKENS"},{"prim":"DIP","args":[[{"prim":"NIL","args":[{"prim":"operation"}]}]]},{"prim":"CONS"},{"prim":"DIP","args":[[{"prim":"UNIT"}]]},{"prim":"PAIR"}]]},"result":{"status":"backtracked","storage":{"prim":"Unit"},"consumed_gas":"77824","storage_size":"46"}},{"kind":"transaction","source":"KT1Njyz94x2pNJGh5uMhKj24VB9JsGCdkySN","nonce":2,"amount":"0","destination":"KT1U8kHUKJVcPkzjHLahLxmnXGnK1StAeNjK","parameters":{"entrypoint":"getTotalSupply","value":{"prim":"Pair","args":[{"prim":"Unit"},{"bytes":"01b73fd2d304bc41891cee0388f1a9e59706952cbd00"}]}},"result":{"status":"backtracked","storage":{"prim":"Pair","args":[{"int":"5118"},{"int":"1670000"}]},"consumed_gas":"94016","storage_size":"3007"}},{"kind":"transaction","source":"KT1U8kHUKJVcPkzjHLahLxmnXGnK1StAeNjK","nonce":3,"amount":"0","destination":"KT1RHhgQZtcbYK8cVvCBymUPqTmKS5PxQbPf","parameters":{"entrypoint":"default","value":{"int":"1670000"}},"result":{"status":"failed","errors":[{"kind":"temporary","id":"proto.006-PsCARTHA.michelson_v1.runtime_error","contract_handle":"KT1RHhgQZtcbYK8cVvCBymUPqTmKS5PxQbPf","contract_code":[{"prim":"parameter","args":[{"prim":"nat"}]},{"prim":"storage","args":[{"prim":"unit"}]},{"prim":"code","args":[[{"prim":"FAILWITH"}]]}]},{"kind":"temporary","id":"proto.006-PsCARTHA.michelson_v1.script_rejected","location":7,"with":{"prim":"Pair","args":[{"int":"1670000"},{"prim":"Unit"}]}}]}}]}}]}`),
						blankHandler,
					),
				))),
				rpc.GetFA12SupplyInput{
					Cycle:        10,
					Source:       "some_source",
					FA12Contract: "some_fa1.2_contract",
					Testnet:      true,
					ChainID:      "some_chainid",
				},
			},
			want{
				false,
				"",
				"1670000",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input.hanler)
			defer server.Close()

			r, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, balance, err := r.GetFA12Supply(tt.input.getFA12SupplyInput)
			checkErr(t, tt.want.err, tt.want.contains, err)
			assert.Equal(t, tt.want.balance, balance)
		})
	}
}

func Test_GetFA12Allowance(t *testing.T) {
	type input struct {
		hanler                http.Handler
		getFA12AllowanceInput rpc.GetFA12AllowanceInput
	}

	type want struct {
		err      bool
		contains string
		balance  string
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"handles failure to validate input",
			input{
				gtGoldenHTTPMock(blankHandler),
				rpc.GetFA12AllowanceInput{},
			},
			want{
				true,
				"invalid input",
				"0",
			},
		},
		{
			"handles failure to get cycle",
			input{
				gtGoldenHTTPMock(mockCycleFailed(blankHandler)),
				rpc.GetFA12AllowanceInput{
					Cycle:          1,
					ChainID:        "some_chainid",
					Source:         "some_source",
					FA12Contract:   "some_fa1.2_contract",
					OwnerAddress:   "some_address",
					SpenderAddress: "some_address",
				},
			},
			want{
				true,
				"could not get cycle",
				"0",
			},
		},
		{
			"handles failure to get counter",
			input{
				gtGoldenHTTPMock(mockCycleSuccessful(mockHandler(
					&requestResultPair{regContractCounter, []byte(`junk`)},
					blankHandler,
				))),
				rpc.GetFA12AllowanceInput{
					Cycle:          1,
					ChainID:        "some_chainid",
					Source:         "some_source",
					FA12Contract:   "some_fa1.2_contract",
					OwnerAddress:   "some_address",
					SpenderAddress: "some_address",
				},
			},
			want{
				true,
				"failed to unmarshal counter",
				"0",
			},
		},
		{
			"handles failure to run_operation",
			input{
				gtGoldenHTTPMock(mockCycleSuccessful(mockHandler(
					&requestResultPair{regContractCounter, []byte(`"100"`)},
					runOperationHandlerMock(
						[]byte(`junk`),
						blankHandler,
					),
				))),
				rpc.GetFA12AllowanceInput{
					Cycle:          1,
					ChainID:        "some_chainid",
					Source:         "some_source",
					FA12Contract:   "some_fa1.2_contract",
					OwnerAddress:   "some_address",
					Testnet:        true,
					SpenderAddress: "some_address",
				},
			},
			want{
				true,
				"failed to unmarshal operation",
				"0",
			},
		},
		{
			"handles failure to parse allowance",
			input{
				gtGoldenHTTPMock(mockCycleSuccessful(mockHandler(
					&requestResultPair{regContractCounter, []byte(`"100"`)},
					runOperationHandlerMock(
						[]byte(`{"contents":[{"kind":"transaction","source":"tz1S82rGFZK8cVbNDpP1Hf9VhTUa4W8oc2WV","fee":"0","counter":"553001","gas_limit":"1040000","storage_limit":"60000","amount":"0","destination":"KT1Njyz94x2pNJGh5uMhKj24VB9JsGCdkySN","parameters":{"entrypoint":"default","value":[{"prim":"PUSH","args":[{"prim":"mutez"},{"int":"0"}]},{"prim":"NONE","args":[{"prim":"key_hash"}]},{"prim":"CREATE_CONTRACT","args":[[{"prim":"parameter","args":[{"prim":"nat"}]},{"prim":"storage","args":[{"prim":"unit"}]},{"prim":"code","args":[[{"prim":"FAILWITH"}]]}]]},{"prim":"DIP","args":[[{"prim":"DIP","args":[[{"prim":"LAMBDA","args":[{"prim":"pair","args":[{"prim":"address"},{"prim":"unit"}]},{"prim":"pair","args":[{"prim":"list","args":[{"prim":"operation"}]},{"prim":"unit"}]},[{"prim":"CAR"},{"prim":"CONTRACT","args":[{"prim":"nat"}]},{"prim":"IF_NONE","args":[[{"prim":"PUSH","args":[{"prim":"string"},{"string":"a"}]},{"prim":"FAILWITH"}],[]]},{"prim":"PUSH","args":[{"prim":"address"},{"string":"tz1MQehPikysuVYN5hTiKTnrsFidAww7rv3z"}]},{"prim":"PAIR"},{"prim":"DIP","args":[[{"prim":"PUSH","args":[{"prim":"address"},{"string":"KT1U8kHUKJVcPkzjHLahLxmnXGnK1StAeNjK"}]},{"prim":"CONTRACT","args":[{"prim":"pair","args":[{"prim":"address"},{"prim":"contract","args":[{"prim":"nat"}]}]}],"annots":["%getBalance"]},{"prim":"IF_NONE","args":[[{"prim":"PUSH","args":[{"prim":"string"},{"string":"b"}]},{"prim":"FAILWITH"}],[]]},{"prim":"PUSH","args":[{"prim":"mutez"},{"int":"0"}]}]]},{"prim":"TRANSFER_TOKENS"},{"prim":"DIP","args":[[{"prim":"NIL","args":[{"prim":"operation"}]}]]},{"prim":"CONS"},{"prim":"DIP","args":[[{"prim":"UNIT"}]]},{"prim":"PAIR"}]]}]]},{"prim":"APPLY"},{"prim":"DIP","args":[[{"prim":"PUSH","args":[{"prim":"address"},{"string":"KT1Njyz94x2pNJGh5uMhKj24VB9JsGCdkySN"}]},{"prim":"CONTRACT","args":[{"prim":"lambda","args":[{"prim":"unit"},{"prim":"pair","args":[{"prim":"list","args":[{"prim":"operation"}]},{"prim":"unit"}]}]}]},{"prim":"IF_NONE","args":[[{"prim":"PUSH","args":[{"prim":"string"},{"string":"c"}]},{"prim":"FAILWITH"}],[]]},{"prim":"PUSH","args":[{"prim":"mutez"},{"int":"0"}]}]]},{"prim":"TRANSFER_TOKENS"},{"prim":"DIP","args":[[{"prim":"NIL","args":[{"prim":"operation"}]}]]},{"prim":"CONS"}]]},{"prim":"CONS"},{"prim":"DIP","args":[[{"prim":"UNIT"}]]},{"prim":"PAIR"}]},"metadata":{"balance_updates":[],"operation_result":{"status":"backtracked","storage":{"prim":"Unit"},"consumed_gas":"26984","storage_size":"46"},"internal_operation_results":[{"kind":"origination","source":"KT1Njyz94x2pNJGh5uMhKj24VB9JsGCdkySN","nonce":0,"balance":"0","script":{"code":[{"prim":"parameter","args":[{"prim":"nat"}]},{"prim":"storage","args":[{"prim":"unit"}]},{"prim":"code","args":[[{"prim":"FAILWITH"}]]}],"storage":{"prim":"Unit"}},"result":{"status":"backtracked","big_map_diff":[],"balance_updates":[{"kind":"contract","contract":"tz1S82rGFZK8cVbNDpP1Hf9VhTUa4W8oc2WV","change":"-32000"},{"kind":"contract","contract":"tz1S82rGFZK8cVbNDpP1Hf9VhTUa4W8oc2WV","change":"-257000"}],"originated_contracts":["KT1TTaLErwMQAVRNB1sVXf9NUdXHLCrnNpUV"],"consumed_gas":"10696","storage_size":"32","paid_storage_size_diff":"32"}},{"kind":"transaction","source":"KT1Njyz94x2pNJGh5uMhKj24VB9JsGCdkySN","nonce":1,"amount":"0","destination":"KT1Njyz94x2pNJGh5uMhKj24VB9JsGCdkySN","parameters":{"entrypoint":"default","value":[{"prim":"PUSH","args":[{"prim":"address"},{"bytes":"01cf0e19e55c5c34ec644b3b1c46c5fe3d8feb96c600"}]},{"prim":"PAIR"},[{"prim":"CAR"},{"prim":"CONTRACT","args":[{"prim":"nat"}]},{"prim":"IF_NONE","args":[[{"prim":"PUSH","args":[{"prim":"string"},{"string":"a"}]},{"prim":"FAILWITH"}],[]]},{"prim":"PUSH","args":[{"prim":"address"},{"bytes":"000013687bbb1298ad36a4bdb2c5da2f126a59aac007"}]},{"prim":"PAIR"},{"prim":"DIP","args":[[{"prim":"PUSH","args":[{"prim":"address"},{"bytes":"01d676a1a3e2e602bbb478bb188265fdec8f09124d00"}]},{"prim":"CONTRACT","args":[{"prim":"pair","args":[{"prim":"address"},{"prim":"contract","args":[{"prim":"nat"}]}]}],"annots":["%getBalance"]},{"prim":"IF_NONE","args":[[{"prim":"PUSH","args":[{"prim":"string"},{"string":"b"}]},{"prim":"FAILWITH"}],[]]},{"prim":"PUSH","args":[{"prim":"mutez"},{"int":"0"}]}]]},{"prim":"TRANSFER_TOKENS"},{"prim":"DIP","args":[[{"prim":"NIL","args":[{"prim":"operation"}]}]]},{"prim":"CONS"},{"prim":"DIP","args":[[{"prim":"UNIT"}]]},{"prim":"PAIR"}]]},"result":{"status":"backtracked","storage":{"prim":"Unit"},"consumed_gas":"77868","storage_size":"46"}},{"kind":"transaction","source":"KT1Njyz94x2pNJGh5uMhKj24VB9JsGCdkySN","nonce":2,"amount":"0","destination":"KT1U8kHUKJVcPkzjHLahLxmnXGnK1StAeNjK","parameters":{"entrypoint":"getBalance","value":{"prim":"Pair","args":[{"bytes":"000013687bbb1298ad36a4bdb2c5da2f126a59aac007"},{"bytes":"01cf0e19e55c5c34ec644b3b1c46c5fe3d8feb96c600"}]}},"result":{"status":"backtracked","storage":{"prim":"Pair","args":[{"int":"5118"},{"int":"1670000"}]},"consumed_gas":"95277","storage_size":"3007"}},{"kind":"transaction","source":"KT1U8kHUKJVcPkzjHLahLxmnXGnK1StAeNjK","nonce":3,"amount":"0","destination":"KT1TTaLErwMQAVRNB1sVXf9NUdXHLCrnNpUV","parameters":{"entrypoint":"default","value":{"int":"1546544"}},"result":{"status":"failed"}}]}}]}`),
						blankHandler,
					),
				))),
				rpc.GetFA12AllowanceInput{
					Cycle:          1,
					ChainID:        "some_chainid",
					Source:         "some_source",
					FA12Contract:   "some_fa1.2_contract",
					OwnerAddress:   "some_address",
					Testnet:        true,
					SpenderAddress: "some_address",
				},
			},
			want{
				true,
				"failed to parse allowance",
				"0",
			},
		},
		{
			"is successful",
			input{
				gtGoldenHTTPMock(mockCycleSuccessful(mockHandler(
					&requestResultPair{regContractCounter, []byte(`"100"`)},
					runOperationHandlerMock(
						[]byte(`{"contents":[{"kind":"transaction","source":"tz1S82rGFZK8cVbNDpP1Hf9VhTUa4W8oc2WV","fee":"0","counter":"553001","gas_limit":"1040000","storage_limit":"60000","amount":"0","destination":"KT1Njyz94x2pNJGh5uMhKj24VB9JsGCdkySN","parameters":{"entrypoint":"default","value":[{"prim":"PUSH","args":[{"prim":"mutez"},{"int":"0"}]},{"prim":"NONE","args":[{"prim":"key_hash"}]},{"prim":"CREATE_CONTRACT","args":[[{"prim":"parameter","args":[{"prim":"nat"}]},{"prim":"storage","args":[{"prim":"unit"}]},{"prim":"code","args":[[{"prim":"FAILWITH"}]]}]]},{"prim":"DIP","args":[[{"prim":"DIP","args":[[{"prim":"LAMBDA","args":[{"prim":"pair","args":[{"prim":"address"},{"prim":"unit"}]},{"prim":"pair","args":[{"prim":"list","args":[{"prim":"operation"}]},{"prim":"unit"}]},[{"prim":"CAR"},{"prim":"CONTRACT","args":[{"prim":"nat"}]},{"prim":"IF_NONE","args":[[{"prim":"PUSH","args":[{"prim":"string"},{"string":"a"}]},{"prim":"FAILWITH"}],[]]},{"prim":"PUSH","args":[{"prim":"address"},{"string":"tz1S82rGFZK8cVbNDpP1Hf9VhTUa4W8oc2WV"}]},{"prim":"PUSH","args":[{"prim":"address"},{"string":"tz1MQehPikysuVYN5hTiKTnrsFidAww7rv3z"}]},{"prim":"PAIR"},{"prim":"PAIR"},{"prim":"DIP","args":[[{"prim":"PUSH","args":[{"prim":"address"},{"string":"KT1UpGMT5arFH2wo7WczhntnawEisdZnsMzc"}]},{"prim":"CONTRACT","args":[{"prim":"pair","args":[{"prim":"pair","args":[{"prim":"address"},{"prim":"address"}]},{"prim":"contract","args":[{"prim":"nat"}]}]}],"annots":["%getAllowance"]},{"prim":"IF_NONE","args":[[{"prim":"PUSH","args":[{"prim":"string"},{"string":"b"}]},{"prim":"FAILWITH"}],[]]},{"prim":"PUSH","args":[{"prim":"mutez"},{"int":"0"}]}]]},{"prim":"TRANSFER_TOKENS"},{"prim":"DIP","args":[[{"prim":"NIL","args":[{"prim":"operation"}]}]]},{"prim":"CONS"},{"prim":"DIP","args":[[{"prim":"UNIT"}]]},{"prim":"PAIR"}]]}]]},{"prim":"APPLY"},{"prim":"DIP","args":[[{"prim":"PUSH","args":[{"prim":"address"},{"string":"KT1Njyz94x2pNJGh5uMhKj24VB9JsGCdkySN"}]},{"prim":"CONTRACT","args":[{"prim":"lambda","args":[{"prim":"unit"},{"prim":"pair","args":[{"prim":"list","args":[{"prim":"operation"}]},{"prim":"unit"}]}]}]},{"prim":"IF_NONE","args":[[{"prim":"PUSH","args":[{"prim":"string"},{"string":"c"}]},{"prim":"FAILWITH"}],[]]},{"prim":"PUSH","args":[{"prim":"mutez"},{"int":"0"}]}]]},{"prim":"TRANSFER_TOKENS"},{"prim":"DIP","args":[[{"prim":"NIL","args":[{"prim":"operation"}]}]]},{"prim":"CONS"}]]},{"prim":"CONS"},{"prim":"DIP","args":[[{"prim":"UNIT"}]]},{"prim":"PAIR"}]},"metadata":{"balance_updates":[],"operation_result":{"status":"backtracked","storage":{"prim":"Unit"},"consumed_gas":"27331","storage_size":"46"},"internal_operation_results":[{"kind":"origination","source":"KT1Njyz94x2pNJGh5uMhKj24VB9JsGCdkySN","nonce":0,"balance":"0","script":{"code":[{"prim":"parameter","args":[{"prim":"nat"}]},{"prim":"storage","args":[{"prim":"unit"}]},{"prim":"code","args":[[{"prim":"FAILWITH"}]]}],"storage":{"prim":"Unit"}},"result":{"status":"backtracked","big_map_diff":[],"balance_updates":[{"kind":"contract","contract":"tz1S82rGFZK8cVbNDpP1Hf9VhTUa4W8oc2WV","change":"-32000"},{"kind":"contract","contract":"tz1S82rGFZK8cVbNDpP1Hf9VhTUa4W8oc2WV","change":"-257000"}],"originated_contracts":["KT19NoRBvBAHNPRS9tXv9mCTggSHitiqgJjV"],"consumed_gas":"10696","storage_size":"32","paid_storage_size_diff":"32"}},{"kind":"transaction","source":"KT1Njyz94x2pNJGh5uMhKj24VB9JsGCdkySN","nonce":1,"amount":"0","destination":"KT1Njyz94x2pNJGh5uMhKj24VB9JsGCdkySN","parameters":{"entrypoint":"default","value":[{"prim":"PUSH","args":[{"prim":"address"},{"bytes":"0108b4ae29b4e42627b4503d384b9b32bbd869654200"}]},{"prim":"PAIR"},[{"prim":"CAR"},{"prim":"CONTRACT","args":[{"prim":"nat"}]},{"prim":"IF_NONE","args":[[{"prim":"PUSH","args":[{"prim":"string"},{"string":"a"}]},{"prim":"FAILWITH"}],[]]},{"prim":"PUSH","args":[{"prim":"address"},{"bytes":"0000471c8882bcf12586e640b7efa46c6ea1e0f4da9e"}]},{"prim":"PUSH","args":[{"prim":"address"},{"bytes":"000013687bbb1298ad36a4bdb2c5da2f126a59aac007"}]},{"prim":"PAIR"},{"prim":"PAIR"},{"prim":"DIP","args":[[{"prim":"PUSH","args":[{"prim":"address"},{"bytes":"01ddeff433fcf03091285c26534ea0bbea08536f9f00"}]},{"prim":"CONTRACT","args":[{"prim":"pair","args":[{"prim":"pair","args":[{"prim":"address"},{"prim":"address"}]},{"prim":"contract","args":[{"prim":"nat"}]}]}],"annots":["%getAllowance"]},{"prim":"IF_NONE","args":[[{"prim":"PUSH","args":[{"prim":"string"},{"string":"b"}]},{"prim":"FAILWITH"}],[]]},{"prim":"PUSH","args":[{"prim":"mutez"},{"int":"0"}]}]]},{"prim":"TRANSFER_TOKENS"},{"prim":"DIP","args":[[{"prim":"NIL","args":[{"prim":"operation"}]}]]},{"prim":"CONS"},{"prim":"DIP","args":[[{"prim":"UNIT"}]]},{"prim":"PAIR"}]]},"result":{"status":"backtracked","storage":{"prim":"Unit"},"consumed_gas":"68699","storage_size":"46"}},{"kind":"transaction","source":"KT1Njyz94x2pNJGh5uMhKj24VB9JsGCdkySN","nonce":2,"amount":"0","destination":"KT1UpGMT5arFH2wo7WczhntnawEisdZnsMzc","parameters":{"entrypoint":"getAllowance","value":{"prim":"Pair","args":[{"prim":"Pair","args":[{"bytes":"000013687bbb1298ad36a4bdb2c5da2f126a59aac007"},{"bytes":"0000471c8882bcf12586e640b7efa46c6ea1e0f4da9e"}]},{"bytes":"0108b4ae29b4e42627b4503d384b9b32bbd869654200"}]}},"result":{"status":"backtracked","storage":{"prim":"Pair","args":[{"int":"18449"},{"int":"1670000"}]},"consumed_gas":"67459","storage_size":"2441"}},{"kind":"transaction","source":"KT1UpGMT5arFH2wo7WczhntnawEisdZnsMzc","nonce":3,"amount":"0","destination":"KT19NoRBvBAHNPRS9tXv9mCTggSHitiqgJjV","parameters":{"entrypoint":"default","value":{"int":"1000000"}},"result":{"status":"failed","errors":[{"kind":"temporary","id":"proto.006-PsCARTHA.michelson_v1.runtime_error","contract_handle":"KT19NoRBvBAHNPRS9tXv9mCTggSHitiqgJjV","contract_code":[{"prim":"parameter","args":[{"prim":"nat"}]},{"prim":"storage","args":[{"prim":"unit"}]},{"prim":"code","args":[[{"prim":"FAILWITH"}]]}]},{"kind":"temporary","id":"proto.006-PsCARTHA.michelson_v1.script_rejected","location":7,"with":{"prim":"Pair","args":[{"int":"1000000"},{"prim":"Unit"}]}}]}}]}}]}`),
						blankHandler,
					),
				))),
				rpc.GetFA12AllowanceInput{
					Cycle:          1,
					ChainID:        "some_chainid",
					Source:         "some_source",
					FA12Contract:   "some_fa1.2_contract",
					OwnerAddress:   "some_address",
					Testnet:        true,
					SpenderAddress: "some_address",
				},
			},
			want{
				false,
				"",
				"1000000",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input.hanler)
			defer server.Close()

			r, err := rpc.New(server.URL)
			assert.Nil(t, err)

			_, balance, err := r.GetFA12Allowance(tt.input.getFA12AllowanceInput)
			checkErr(t, tt.want.err, tt.want.contains, err)
			assert.Equal(t, tt.want.balance, balance)
		})
	}
}
