package gotezos

import (
	"encoding/json"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Delegations", func() {
	It("returns rpc error", func() {
		server := httptest.NewServer(gtGoldenHTTPMock(delegationsHandlerMock(mockRPCError, blankHandler)))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		delegations, err := gt.Delegations(mockHash, "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc")
		Expect(err).NotTo(Succeed())
		Expect(err.Error()).To(MatchRegexp("could not get delegations for"))
		Expect(delegations).To(Equal([]string{}))
	})

	It("fails to unmarshal", func() {
		server := httptest.NewServer(gtGoldenHTTPMock(delegationsHandlerMock([]byte(`junk`), blankHandler)))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		delegations, err := gt.Delegations(mockHash, "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc")
		Expect(err).NotTo(Succeed())
		Expect(err.Error()).To(MatchRegexp("could not unmarshal delegations for"))
		Expect(delegations).To(Equal([]string{}))
	})

	It("is successful", func() {
		server := httptest.NewServer(gtGoldenHTTPMock(delegationsHandlerMock(mockDelegations, blankHandler)))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		var want []string
		json.Unmarshal(mockDelegations, &want)

		delegations, err := gt.Delegations(mockHash, "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc")
		Expect(err).To(Succeed())
		Expect(delegations).To(Equal(want))
	})
})

var _ = Describe("DelegationsAtCycle", func() {
	It("failed to get cycle", func() {
		server := httptest.NewServer(gtGoldenHTTPMock(mockCycleFailed(blankHandler)))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		delegations, err := gt.DelegationsAtCycle(10, "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc")
		Expect(err).NotTo(Succeed())
		Expect(err.Error()).To(MatchRegexp("could not get delegations for"))
		Expect(delegations).To(Equal([]string{}))
	})

	It("failed to get delegations", func() {
		server := httptest.NewServer(gtGoldenHTTPMock(mockCycleSuccessful(delegationsHandlerMock([]byte(`junk`), blankHandler))))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		delegations, err := gt.DelegationsAtCycle(10, "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc")
		Expect(err).NotTo(Succeed())
		Expect(err.Error()).To(MatchRegexp("could not get delegations at cycle"))
		Expect(delegations).To(Equal([]string{}))
	})

	It("is successful", func() {
		server := httptest.NewServer(gtGoldenHTTPMock(mockCycleSuccessful(delegationsHandlerMock(mockDelegations, blankHandler))))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		var want []string
		json.Unmarshal(mockDelegations, &want)

		delegations, err := gt.DelegationsAtCycle(10, "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc")
		Expect(err).To(Succeed())
		Expect(delegations).To(Equal(want))
	})
})

var _ = Describe("FrozenBalance", func() {
	It("failed to get cycle", func() {
		server := httptest.NewServer(gtGoldenHTTPMock(mockCycleFailed(blankHandler)))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		fbalance, err := gt.FrozenBalance(10, "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc")
		Expect(err).NotTo(Succeed())
		Expect(err.Error()).To(MatchRegexp("could not get frozen balance at cycle"))
		Expect(fbalance).To(Equal(FrozenBalance{}))
	})

	It("returns rpc error", func() {
		server := httptest.NewServer(gtGoldenHTTPMock(mockCycleSuccessful(frozenBalanceHandlerMock(mockRPCError, blankHandler))))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		fbalance, err := gt.FrozenBalance(10, "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc")
		Expect(err).NotTo(Succeed())
		Expect(err.Error()).To(MatchRegexp("could not get frozen balance at cycle"))
		Expect(fbalance).To(Equal(FrozenBalance{}))
	})

	It("fails to unmarshal", func() {
		server := httptest.NewServer(gtGoldenHTTPMock(mockCycleSuccessful(frozenBalanceHandlerMock([]byte(`junk`), blankHandler))))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		fbalance, err := gt.FrozenBalance(10, "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc")
		Expect(err).NotTo(Succeed())
		Expect(err.Error()).To(MatchRegexp("could not unmarshal frozen balance at cycle"))
		Expect(fbalance).To(Equal(FrozenBalance{}))
	})

	It("is successful", func() {
		server := httptest.NewServer(gtGoldenHTTPMock(mockCycleSuccessful(frozenBalanceHandlerMock(mockFrozenBalance, blankHandler))))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		fbalance, err := gt.FrozenBalance(10, "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc")
		Expect(err).To(Succeed())
		Expect(fbalance).To(Equal(FrozenBalance{
			Deposits: "15296000000",
			Fees:     "76724",
			Rewards:  "474800000",
		}))
	})
})

var _ = Describe("Delegate", func() {
	It("returns rpc error", func() {
		server := httptest.NewServer(gtGoldenHTTPMock(delegateHandlerMock(mockRPCError, blankHandler)))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		delegate, err := gt.Delegate(mockHash, "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc")
		Expect(err).NotTo(Succeed())
		Expect(err.Error()).To(MatchRegexp("could not get delegate"))
		Expect(delegate).To(Equal(Delegate{}))
	})

	It("fails to unmarshal", func() {
		server := httptest.NewServer(gtGoldenHTTPMock(delegateHandlerMock([]byte(`junk`), blankHandler)))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		delegate, err := gt.Delegate(mockHash, "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc")
		Expect(err).NotTo(Succeed())
		Expect(err.Error()).To(MatchRegexp("could not unmarshal delegate"))
		Expect(delegate).To(Equal(Delegate{}))
	})

	It("is successful", func() {
		server := httptest.NewServer(gtGoldenHTTPMock(delegateHandlerMock(mockDelegate, blankHandler)))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		var want Delegate
		json.Unmarshal(mockDelegate, &want)

		delegations, err := gt.Delegate(mockHash, "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc")
		Expect(err).To(Succeed())
		Expect(delegations).To(Equal(want))
	})
})

var _ = Describe("StakingBalance", func() {
	It("returns rpc error", func() {
		server := httptest.NewServer(gtGoldenHTTPMock(stakingBalanceHandlerMock(mockRPCError, blankHandler)))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		sbalance, err := gt.StakingBalance(mockHash, "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc")
		Expect(err).NotTo(Succeed())
		Expect(err.Error()).To(MatchRegexp("could not get staking balance"))
		Expect(sbalance).To(Equal(""))
	})

	It("fails to unmarshal", func() {
		server := httptest.NewServer(gtGoldenHTTPMock(stakingBalanceHandlerMock([]byte(`junk`), blankHandler)))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		sbalance, err := gt.StakingBalance(mockHash, "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc")
		Expect(err).NotTo(Succeed())
		Expect(err.Error()).To(MatchRegexp("could not unmarshal staking balance"))
		Expect(sbalance).To(Equal(""))
	})

	It("is successful", func() {
		server := httptest.NewServer(gtGoldenHTTPMock(stakingBalanceHandlerMock(mockStakingBalance, blankHandler)))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		var want string
		json.Unmarshal(mockStakingBalance, &want)

		sbalance, err := gt.StakingBalance(mockHash, "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc")
		Expect(err).To(Succeed())
		Expect(sbalance).To(Equal(want))
	})
})

var _ = Describe("StakingBalanceAtCycle", func() {
	It("failed to get cycle", func() {
		server := httptest.NewServer(gtGoldenHTTPMock(mockCycleFailed(blankHandler)))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		sbalance, err := gt.StakingBalanceAtCycle(10, "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc")
		Expect(err).NotTo(Succeed())
		Expect(err.Error()).To(MatchRegexp("could not get staking balance for"))
		Expect(sbalance).To(Equal(""))
	})

	It("failed to get staking balance", func() {
		server := httptest.NewServer(gtGoldenHTTPMock(mockCycleSuccessful(stakingBalanceHandlerMock([]byte(`junk`), blankHandler))))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		sbalance, err := gt.StakingBalanceAtCycle(10, "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc")
		Expect(err).NotTo(Succeed())
		Expect(err.Error()).To(MatchRegexp("could not get staking balance for"))
		Expect(sbalance).To(Equal(""))
	})

	It("is successful", func() {
		server := httptest.NewServer(gtGoldenHTTPMock(mockCycleSuccessful(stakingBalanceHandlerMock(mockStakingBalance, blankHandler))))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		var want string
		json.Unmarshal(mockStakingBalance, &want)

		delegations, err := gt.StakingBalanceAtCycle(10, "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc")
		Expect(err).To(Succeed())
		Expect(delegations).To(Equal(want))
	})
})

var _ = Describe("BakingRights", func() {
	It("returns rpc error", func() {
		server := httptest.NewServer(gtGoldenHTTPMock(bakingRightsHandlerMock(mockRPCError, blankHandler)))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		rights, err := gt.BakingRights(mockHash, 0)
		Expect(err).NotTo(Succeed())
		Expect(err.Error()).To(MatchRegexp("could not get baking rights"))
		Expect(rights).To(Equal(BakingRights{}))
	})

	It("fails to unmarshal", func() {
		server := httptest.NewServer(gtGoldenHTTPMock(bakingRightsHandlerMock([]byte(`junk`), blankHandler)))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		rights, err := gt.BakingRights(mockHash, 0)
		Expect(err).NotTo(Succeed())
		Expect(err.Error()).To(MatchRegexp("could not unmarshal baking rights"))
		Expect(rights).To(Equal(BakingRights{}))
	})

	It("is successful", func() {
		server := httptest.NewServer(gtGoldenHTTPMock(bakingRightsHandlerMock(mockBakingRights, blankHandler)))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		var want BakingRights
		json.Unmarshal(mockBakingRights, &want)

		rights, err := gt.BakingRights(mockHash, 0)
		Expect(err).To(Succeed())
		Expect(rights).To(Equal(want))
	})
})

var _ = Describe("BakingRightsAtCycle", func() {
	It("failed to get cycle", func() {
		server := httptest.NewServer(gtGoldenHTTPMock(mockCycleFailed(blankHandler)))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		rights, err := gt.BakingRightsAtCycle(10, 0)
		Expect(err).NotTo(Succeed())
		Expect(err.Error()).To(MatchRegexp("could not get baking rights"))
		Expect(rights).To(Equal(BakingRights{}))
	})

	It("returns rpc error", func() {
		server := httptest.NewServer(gtGoldenHTTPMock(mockCycleSuccessful(bakingRightsHandlerMock(mockRPCError, blankHandler))))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		rights, err := gt.BakingRightsAtCycle(10, 0)
		Expect(err).NotTo(Succeed())
		Expect(err.Error()).To(MatchRegexp("could not get baking rights"))
		Expect(rights).To(Equal(BakingRights{}))
	})

	It("fails to unmarshal", func() {
		server := httptest.NewServer(gtGoldenHTTPMock(mockCycleSuccessful(bakingRightsHandlerMock([]byte(`junk`), blankHandler))))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		rights, err := gt.BakingRightsAtCycle(10, 0)
		Expect(err).NotTo(Succeed())
		Expect(err.Error()).To(MatchRegexp("could not unmarshal baking rights"))
		Expect(rights).To(Equal(BakingRights{}))
	})

	It("is successful", func() {
		server := httptest.NewServer(gtGoldenHTTPMock(mockCycleSuccessful(bakingRightsHandlerMock(mockBakingRights, blankHandler))))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		var want BakingRights
		json.Unmarshal(mockBakingRights, &want)

		rights, err := gt.BakingRightsAtCycle(10, 0)
		Expect(err).To(Succeed())
		Expect(rights).To(Equal(want))
	})
})
