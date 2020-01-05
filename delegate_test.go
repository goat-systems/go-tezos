package gotezos

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Delegations(t *testing.T) {
	var goldenDelegations []string
	json.Unmarshal(mockDelegations, &goldenDelegations)

	type want struct {
		wantErr         bool
		containsErr     string
		wantDelegations []string
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"returns rpc error",
			gtGoldenHTTPMock(delegationsHandlerMock(mockRPCError, blankHandler)),
			want{
				true,
				"could not get delegations for",
				[]string{},
			},
		},
		{
			"fails to unmarshal",
			gtGoldenHTTPMock(delegationsHandlerMock([]byte(`junk`), blankHandler)),
			want{
				true,
				"could not unmarshal delegations for",
				[]string{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(delegationsHandlerMock(mockDelegations, blankHandler)),
			want{
				false,
				"",
				goldenDelegations,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			gt, err := New(server.URL)
			assert.Nil(t, err)

			delegations, err := gt.Delegations(mockHash, "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc")
			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tt.want.containsErr)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, tt.want.wantDelegations, delegations)
		})
	}
}

func Test_DelegationsAtCycle(t *testing.T) {

	var goldenDelegations []string
	json.Unmarshal(mockDelegations, &goldenDelegations)

	type want struct {
		wantErr         bool
		containsErr     string
		wantDelegations []string
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"failed to get cycle",
			gtGoldenHTTPMock(mockCycleFailed(blankHandler)),
			want{
				true,
				"could not get delegations for",
				[]string{},
			},
		},
		{
			"failed to get delegations",
			gtGoldenHTTPMock(mockCycleSuccessful(delegationsHandlerMock([]byte(`junk`), blankHandler))),
			want{
				true,
				"could not get delegations at cycle",
				[]string{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockCycleSuccessful(delegationsHandlerMock(mockDelegations, blankHandler))),
			want{
				false,
				"",
				goldenDelegations,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			gt, err := New(server.URL)
			assert.Nil(t, err)

			delegations, err := gt.DelegationsAtCycle(10, "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc")
			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tt.want.containsErr)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, tt.want.wantDelegations, delegations)
		})
	}
}

func Test_FrozenBalance(t *testing.T) {

	var goldenFrozenBalance FrozenBalance
	json.Unmarshal(mockFrozenBalance, &goldenFrozenBalance)

	type want struct {
		wantErr           bool
		containsErr       string
		wantFrozenBalance FrozenBalance
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"failed to get cycle",
			gtGoldenHTTPMock(mockCycleFailed(blankHandler)),
			want{
				true,
				"could not get frozen balance at cycle",
				FrozenBalance{},
			},
		},
		{
			"returns rpc error",
			gtGoldenHTTPMock(mockCycleSuccessful(frozenBalanceHandlerMock(mockRPCError, blankHandler))),
			want{
				true,
				"could not get frozen balance at cycle",
				FrozenBalance{},
			},
		},
		{
			"fails to unmarshal",
			gtGoldenHTTPMock(mockCycleSuccessful(frozenBalanceHandlerMock([]byte(`junk`), blankHandler))),
			want{
				true,
				"could not unmarshal frozen balance at cycle",
				FrozenBalance{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockCycleSuccessful(frozenBalanceHandlerMock(mockFrozenBalance, blankHandler))),
			want{
				false,
				"",
				goldenFrozenBalance,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			gt, err := New(server.URL)
			assert.Nil(t, err)

			frozenBalance, err := gt.FrozenBalance(10, "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc")
			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tt.want.containsErr)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, tt.want.wantFrozenBalance, frozenBalance)
		})
	}
}

func Test_Delegate(t *testing.T) {

	var goldenDelegate Delegate
	json.Unmarshal(mockDelegate, &goldenDelegate)

	type want struct {
		wantErr      bool
		containsErr  string
		wantDelegate Delegate
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"returns rpc error",
			gtGoldenHTTPMock(delegateHandlerMock(mockRPCError, blankHandler)),
			want{
				true,
				"could not get delegate",
				Delegate{},
			},
		},
		{
			"fails to unmarshal",
			gtGoldenHTTPMock(delegateHandlerMock([]byte(`junk`), blankHandler)),
			want{
				true,
				"could not unmarshal delegate",
				Delegate{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(delegateHandlerMock(mockDelegate, blankHandler)),
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

			gt, err := New(server.URL)
			assert.Nil(t, err)

			delegate, err := gt.Delegate(mockHash, "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc")
			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tt.want.containsErr)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, tt.want.wantDelegate, delegate)
		})
	}
}

func Test_StakingBalance(t *testing.T) {

	var goldenStakingBalance string
	json.Unmarshal(mockStakingBalance, &goldenStakingBalance)

	type want struct {
		wantErr            bool
		containsErr        string
		wantStakingBalance string
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"returns rpc error",
			gtGoldenHTTPMock(stakingBalanceHandlerMock(mockRPCError, blankHandler)),
			want{
				true,
				"could not get staking balance",
				"",
			},
		},
		{
			"fails to unmarshal",
			gtGoldenHTTPMock(stakingBalanceHandlerMock([]byte(`junk`), blankHandler)),
			want{
				true,
				"could not unmarshal staking balance",
				"",
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(stakingBalanceHandlerMock(mockStakingBalance, blankHandler)),
			want{
				false,
				"",
				goldenStakingBalance,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			gt, err := New(server.URL)
			assert.Nil(t, err)

			stakingBalance, err := gt.StakingBalance(mockHash, "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc")
			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tt.want.containsErr)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, tt.want.wantStakingBalance, stakingBalance)
		})
	}
}

func Test_StakingBalanceAtCycle(t *testing.T) {

	var goldenStakingBalance string
	json.Unmarshal(mockStakingBalance, &goldenStakingBalance)

	type want struct {
		wantErr            bool
		containsErr        string
		wantStakingBalance string
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"failed to get cycle",
			gtGoldenHTTPMock(mockCycleFailed(blankHandler)),
			want{
				true,
				"could not get staking balance for",
				"",
			},
		},
		{
			"failed to get staking balance",
			gtGoldenHTTPMock(mockCycleSuccessful(stakingBalanceHandlerMock([]byte(`junk`), blankHandler))),
			want{
				true,
				"could not unmarshal staking balance",
				"",
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockCycleSuccessful(stakingBalanceHandlerMock(mockStakingBalance, blankHandler))),
			want{
				false,
				"",
				goldenStakingBalance,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			gt, err := New(server.URL)
			assert.Nil(t, err)

			stakingBalance, err := gt.StakingBalanceAtCycle(10, "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc")
			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tt.want.containsErr)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, tt.want.wantStakingBalance, stakingBalance)
		})
	}
}

func Test_BakingRights(t *testing.T) {

	var goldenBakingRights BakingRights
	json.Unmarshal(mockBakingRights, &goldenBakingRights)

	type want struct {
		wantErr          bool
		containsErr      string
		wantBakingRights BakingRights
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"returns rpc error",
			gtGoldenHTTPMock(bakingRightsHandlerMock(mockRPCError, blankHandler)),
			want{
				true,
				"could not get baking rights",
				BakingRights{},
			},
		},
		{
			"fails to unmarshal",
			gtGoldenHTTPMock(bakingRightsHandlerMock([]byte(`junk`), blankHandler)),
			want{
				true,
				"could not unmarshal baking rights",
				BakingRights{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(bakingRightsHandlerMock(mockBakingRights, blankHandler)),
			want{
				false,
				"",
				goldenBakingRights,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			gt, err := New(server.URL)
			assert.Nil(t, err)

			bakingRights, err := gt.BakingRights(mockHash, 0)
			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tt.want.containsErr)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, tt.want.wantBakingRights, bakingRights)
		})
	}
}

func Test_BakingRightsAtCycle(t *testing.T) {

	var goldenBakingRights BakingRights
	json.Unmarshal(mockBakingRights, &goldenBakingRights)

	type want struct {
		wantErr          bool
		containsErr      string
		wantBakingRights BakingRights
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"failed to get cycle",
			gtGoldenHTTPMock(mockCycleFailed(blankHandler)),
			want{
				true,
				"could not get baking rights",
				BakingRights{},
			},
		},
		{
			"returns rpc error",
			gtGoldenHTTPMock(mockCycleSuccessful(bakingRightsHandlerMock(mockRPCError, blankHandler))),
			want{
				true,
				"could not get baking rights",
				BakingRights{},
			},
		},
		{
			"fails to unmarshal",
			gtGoldenHTTPMock(mockCycleSuccessful(bakingRightsHandlerMock([]byte(`junk`), blankHandler))),
			want{
				true,
				"could not unmarshal baking rights",
				BakingRights{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(mockCycleSuccessful(bakingRightsHandlerMock(mockBakingRights, blankHandler))),
			want{
				false,
				"",
				goldenBakingRights,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			gt, err := New(server.URL)
			assert.Nil(t, err)

			bakingRights, err := gt.BakingRightsAtCycle(10, 0)
			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tt.want.containsErr)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, tt.want.wantBakingRights, bakingRights)
		})
	}
}
