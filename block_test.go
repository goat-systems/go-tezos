package gotezos

import (
	"encoding/json"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Head", func() {
	It("failed to unmarshal", func() {
		var blockmock blockMock
		server := httptest.NewServer(gtGoldenHTTPMock(blockmock.handler([]byte(`not_block_data`), blankHandler)))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		head, err := gt.Head()
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(MatchRegexp("could not get head block: invalid character"))
		Expect(head).To(Equal(Block{}))
	})

	It("is successful", func() {
		var blockmock blockMock
		server := httptest.NewServer(gtGoldenHTTPMock(blockmock.handler(mockBlockRandom, blankHandler)))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		var wantHead Block
		json.Unmarshal(mockBlockRandom, &wantHead)

		head, err := gt.Head()
		Expect(err).To(BeNil())
		Expect(head).To(Equal(wantHead))
	})
})

var _ = Describe("Block", func() {
	It("failed to unmarshal", func() {
		var blockmock blockMock
		server := httptest.NewServer(gtGoldenHTTPMock(blockmock.handler([]byte(`not_block_data`), blankHandler)))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		head, err := gt.Block(50)
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(MatchRegexp("could not get block '50': invalid character"))
		Expect(head).To(Equal(Block{}))
	})

	It("is successful", func() {
		var blockmock blockMock
		server := httptest.NewServer(gtGoldenHTTPMock(blockmock.handler(mockBlockRandom, blankHandler)))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		var wantHead Block
		json.Unmarshal(mockBlockRandom, &wantHead)

		head, err := gt.Block(50)
		Expect(err).To(BeNil())
		Expect(head).To(Equal(wantHead))
	})
})

var _ = Describe("OperationHashes", func() {
	It("failed to unmarshal", func() {
		server := httptest.NewServer(gtGoldenHTTPMock(opHashesHandlerMock([]byte(`junk`), blankHandler)))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		hashes, err := gt.OperationHashes("BLzGD63HA4RP8Fh5xEtvdQSMKa2WzJMZjQPNVUc4Rqy8Lh5BEY1")
		Expect(err).NotTo(Succeed())
		Expect(err.Error()).To(MatchRegexp("could not unmarshal operation hashes"))
		Expect(hashes).To(Equal([]string{}))
	})

	It("is successful", func() {
		server := httptest.NewServer(gtGoldenHTTPMock(opHashesHandlerMock(mockOpHashes, blankHandler)))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		hashes, err := gt.OperationHashes("BLzGD63HA4RP8Fh5xEtvdQSMKa2WzJMZjQPNVUc4Rqy8Lh5BEY1")
		Expect(err).To(Succeed())
		Expect(hashes).To(Equal([]string{
			"BLzGD63HA4RP8Fh5xEtvdQSMKa2WzJMZjQPNVUc4Rqy8Lh5BEY1",
			"BLzGD63HA4RP8Fh5xEtvdQSMKa2WzJMZjQPNVUc4Rqy8Lh5BEY1",
			"BLzGD63HA4RP8Fh5xEtvdQSMKa2WzJMZjQPNVUc4Rqy8Lh5BEY1",
		}))
	})
})

var _ = Describe("idToString", func() {
	It("uses integer id", func() {
		str, err := idToString(50)
		Expect(err).To(Succeed())
		Expect(str).To(Equal("50"))
	})

	It("uses integer string", func() {
		str, err := idToString("BLzGD63HA4RP8Fh5xEtvdQSMKa2WzJMZjQPNVUc4Rqy8Lh5BEY1")
		Expect(err).To(Succeed())
		Expect(str).To(Equal("BLzGD63HA4RP8Fh5xEtvdQSMKa2WzJMZjQPNVUc4Rqy8Lh5BEY1"))
	})

	It("uses bad id type", func() {
		str, err := idToString(45.433)
		Expect(err).NotTo(Succeed())
		Expect(err.Error()).To(Equal("id must be block level (int) or block hash (string)"))
		Expect(str).To(Equal(""))
	})
})
