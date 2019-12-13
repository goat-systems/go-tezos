package gotezos

import (
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cycle", func() {
	var blockmock blockMock
	BeforeEach(func() {
		blockmock = blockMock{}
	})

	It("failed to get head block", func() {
		server := httptest.NewServer(gtGoldenHTTPMock(blockmock.handler([]byte(`not_block_data`), blankHandler)))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		cycle, err := gt.Cycle(10)
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(MatchRegexp("could not get cycle '10': could not get head block"))
		Expect(cycle).To(Equal(Cycle{}))
	})

	It("failed to get cycle because cycle is in the future", func() {
		server := httptest.NewServer(gtGoldenHTTPMock(blockmock.handler(mockBlockRandom, blankHandler)))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		cycle, err := gt.Cycle(300)
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(MatchRegexp("request is in the future"))
		Expect(cycle).To(Equal(Cycle{}))
	})

	It("failed to get block less than cycle", func() {
		var oldHTTPBlock blockMock

		server := httptest.NewServer(
			gtGoldenHTTPMock(
				blockmock.handler(
					mockBlockRandom,
					oldHTTPBlock.handler(
						[]byte(`not_block_data`),
						blankHandler,
					),
				),
			),
		)
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		cycle, err := gt.Cycle(2)
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(MatchRegexp("could not get block"))
		Expect(cycle).To(Equal(Cycle{}))
	})

	It("failed to unmarshal cycle", func() {
		var oldHTTPBlock blockMock

		server := httptest.NewServer(
			gtGoldenHTTPMock(
				cycleHandlerMock(
					[]byte(`bad_cycle_data`),
					blockmock.handler(
						mockBlockRandom,
						oldHTTPBlock.handler(
							mockBlockRandom,
							blankHandler,
						),
					),
				),
			),
		)
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		cycle, err := gt.Cycle(2)
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(MatchRegexp("could not unmarshal at cycle hash"))
		Expect(cycle).To(Equal(Cycle{}))
	})

	It("failed to get cycle block level", func() {
		var oldHTTPBlock blockMock
		var blockAtLevel blockMock

		server := httptest.NewServer(
			gtGoldenHTTPMock(
				cycleHandlerMock(
					mockCycle,
					blockmock.handler(
						mockBlockRandom,
						oldHTTPBlock.handler(
							mockBlockRandom,
							blockAtLevel.handler(
								[]byte(`not_block_data`),
								blankHandler,
							),
						),
					),
				),
			),
		)
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		cycle, err := gt.Cycle(2)
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(MatchRegexp("could not get block"))
		Expect(cycle).To(Equal(Cycle{
			RandomSeed:   "04dca5c197fc2e18309b60844148c55fc7ccdbcb498bd57acd4ac29f16e22846",
			RollSnapshot: 4,
		}))
	})

	It("is successful", func() {
		var oldHTTPBlock blockMock
		var blockAtLevel blockMock

		server := httptest.NewServer(
			gtGoldenHTTPMock(
				cycleHandlerMock(
					mockCycle,
					blockmock.handler(
						mockBlockRandom,
						oldHTTPBlock.handler(
							mockBlockRandom,
							blockAtLevel.handler(
								mockBlockRandom,
								blankHandler,
							),
						),
					),
				),
			),
		)
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		cycle, err := gt.Cycle(2)
		Expect(err).To(BeNil())
		Expect(cycle).To(Equal(Cycle{
			RandomSeed:   "04dca5c197fc2e18309b60844148c55fc7ccdbcb498bd57acd4ac29f16e22846",
			RollSnapshot: 4,
			BlockHash:    "BLzGD63HA4RP8Fh5xEtvdQSMKa2WzJMZjQPNVUc4Rqy8Lh5BEY1",
		}))
	})
})
