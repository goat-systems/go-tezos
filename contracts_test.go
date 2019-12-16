package gotezos

import (
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ContractStorage", func() {
	It("gets RPC error", func() {
		server := httptest.NewServer(gtGoldenHTTPMock(storageHandlerMock(mockRPCError, blankHandler)))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		storage, err := gt.ContractStorage("BLzGD63HA4RP8Fh5xEtvdQSMKa2WzJMZjQPNVUc4Rqy8Lh5BEY1", "KT1LfoE9EbpdsfUzowRckGUfikGcd5PyVKg")
		Expect(err).NotTo(Succeed())
		Expect(err.Error()).To(MatchRegexp("could not get storage"))
		Expect(storage).To(Equal(mockRPCError))
	})

	It("is successful", func() {
		server := httptest.NewServer(gtGoldenHTTPMock(storageHandlerMock([]byte(`"Hello Tezos!"`), blankHandler)))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		storage, err := gt.ContractStorage("BLzGD63HA4RP8Fh5xEtvdQSMKa2WzJMZjQPNVUc4Rqy8Lh5BEY1", "KT1LfoE9EbpdsfUzowRckGUfikGcd5PyVKg")
		Expect(err).To(Succeed())
		Expect(storage).To(Equal([]byte(`"Hello Tezos!"`)))
	})
})
