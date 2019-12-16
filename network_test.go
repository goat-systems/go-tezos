package gotezos

import (
	"encoding/json"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Versions", func() {
	It("returns rpc error", func() {
		server := httptest.NewServer(gtGoldenHTTPMock(versionsHandlerMock(mockRPCError, blankHandler)))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		versions, err := gt.Versions()
		Expect(err).NotTo(Succeed())
		Expect(err.Error()).To(MatchRegexp("could not get network versions"))
		Expect(versions).To(Equal(Versions{}))
	})

	It("fails to unmarshal", func() {
		server := httptest.NewServer(gtGoldenHTTPMock(versionsHandlerMock([]byte(`junk`), blankHandler)))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		versions, err := gt.Versions()
		Expect(err).NotTo(Succeed())
		Expect(err.Error()).To(MatchRegexp("could not unmarshal network versions"))
		Expect(versions).To(Equal(Versions{}))
	})

	It("is successful", func() {
		server := httptest.NewServer(gtGoldenHTTPMock(versionsHandlerMock(mockVersions, blankHandler)))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		var want Versions
		json.Unmarshal(mockVersions, &want)

		versions, err := gt.Versions()
		Expect(err).To(Succeed())
		Expect(versions).To(Equal(want))
	})
})

var _ = Describe("Constants", func() {
	var mock constantsMock

	BeforeEach(func() {
		mock = constantsMock{}
	})
	It("returns rpc error", func() {
		server := httptest.NewServer(gtGoldenHTTPMock(mock.handler(mockRPCError, blankHandler)))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		constants, err := gt.Constants("BLzGD63HA4RP8Fh5xEtvdQSMKa2WzJMZjQPNVUc4Rqy8Lh5BEY1")
		Expect(err).NotTo(Succeed())
		Expect(err.Error()).To(MatchRegexp("could not get network constants"))
		Expect(constants).To(Equal(Constants{}))
	})

	It("fails to unmarshal", func() {
		server := httptest.NewServer(gtGoldenHTTPMock(mock.handler([]byte(`junk`), blankHandler)))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		constants, err := gt.Constants("BLzGD63HA4RP8Fh5xEtvdQSMKa2WzJMZjQPNVUc4Rqy8Lh5BEY1")
		Expect(err).NotTo(Succeed())
		Expect(err.Error()).To(MatchRegexp("could not unmarshal network constants"))
		Expect(constants).To(Equal(Constants{}))
	})

	It("is successful", func() {
		server := httptest.NewServer(gtGoldenHTTPMock(mock.handler(mockConstants, blankHandler)))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		var want Constants
		json.Unmarshal(mockConstants, &want)

		constants, err := gt.Constants("BLzGD63HA4RP8Fh5xEtvdQSMKa2WzJMZjQPNVUc4Rqy8Lh5BEY1")
		Expect(err).To(Succeed())
		Expect(constants).To(Equal(want))
	})
})

var _ = Describe("ChainID", func() {
	It("returns rpc error", func() {
		server := httptest.NewServer(gtGoldenHTTPMock(chainIDHandlerMock(mockRPCError, blankHandler)))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		chainID, err := gt.ChainID()
		Expect(err).NotTo(Succeed())
		Expect(err.Error()).To(MatchRegexp("could not get chain id"))
		Expect(chainID).To(Equal(""))
	})

	It("fails to unmarshal", func() {
		server := httptest.NewServer(gtGoldenHTTPMock(chainIDHandlerMock([]byte(`junk`), blankHandler)))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		chainID, err := gt.ChainID()
		Expect(err).NotTo(Succeed())
		Expect(err.Error()).To(MatchRegexp("could unmarshal chain id"))
		Expect(chainID).To(Equal(""))
	})

	It("is successful", func() {
		server := httptest.NewServer(gtGoldenHTTPMock(chainIDHandlerMock(mockChainID, blankHandler)))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		var want string
		json.Unmarshal(mockChainID, &want)

		chainID, err := gt.ChainID()
		Expect(err).To(Succeed())
		Expect(chainID).To(Equal(want))
	})
})

var _ = Describe("Connections", func() {
	It("returns rpc error", func() {
		server := httptest.NewServer(gtGoldenHTTPMock(connectionsHandlerMock(mockRPCError, blankHandler)))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		connections, err := gt.Connections()
		Expect(err).NotTo(Succeed())
		Expect(err.Error()).To(MatchRegexp("could not get network connections"))
		Expect(connections).To(Equal(Connections{}))
	})

	It("fails to unmarshal", func() {
		server := httptest.NewServer(gtGoldenHTTPMock(connectionsHandlerMock([]byte(`junk`), blankHandler)))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		connections, err := gt.Connections()
		Expect(err).NotTo(Succeed())
		Expect(err.Error()).To(MatchRegexp("could not unmarshal network connections"))
		Expect(connections).To(Equal(Connections{}))
	})

	It("is successful", func() {
		server := httptest.NewServer(gtGoldenHTTPMock(connectionsHandlerMock(mockConnections, blankHandler)))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		var want Connections
		json.Unmarshal(mockConnections, &want)

		connections, err := gt.Connections()
		Expect(err).To(Succeed())
		Expect(connections).To(Equal(want))
	})
})

var _ = Describe("Bootsrap", func() {
	It("returns rpc error", func() {
		server := httptest.NewServer(gtGoldenHTTPMock(bootstrapHandlerMock(mockRPCError, blankHandler)))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		bootstrap, err := gt.Bootstrap()
		Expect(err).NotTo(Succeed())
		Expect(err.Error()).To(MatchRegexp("could not get bootstrap"))
		Expect(bootstrap).To(Equal(Bootstrap{}))
	})

	It("fails to unmarshal", func() {
		server := httptest.NewServer(gtGoldenHTTPMock(bootstrapHandlerMock([]byte(`junk`), blankHandler)))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		bootstrap, err := gt.Bootstrap()
		Expect(err).NotTo(Succeed())
		Expect(err.Error()).To(MatchRegexp("could not unmarshal bootstrap"))
		Expect(bootstrap).To(Equal(Bootstrap{}))
	})

	It("is successful", func() {
		server := httptest.NewServer(gtGoldenHTTPMock(bootstrapHandlerMock(mockBootstrap, blankHandler)))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		var want Bootstrap
		json.Unmarshal(mockBootstrap, &want)

		bootstrap, err := gt.Bootstrap()
		Expect(err).To(Succeed())
		Expect(bootstrap).To(Equal(want))
	})
})

var _ = Describe("Commit", func() {
	It("returns rpc error", func() {
		server := httptest.NewServer(gtGoldenHTTPMock(commitHandlerMock(mockRPCError, blankHandler)))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		commit, err := gt.Commit()
		Expect(err).NotTo(Succeed())
		Expect(err.Error()).To(MatchRegexp("could not get commit"))
		Expect(commit).To(Equal(""))
	})

	It("fails to unmarshal", func() {
		server := httptest.NewServer(gtGoldenHTTPMock(commitHandlerMock([]byte(`junk`), blankHandler)))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		commit, err := gt.Commit()
		Expect(err).NotTo(Succeed())
		Expect(err.Error()).To(MatchRegexp("could unmarshal commit"))
		Expect(commit).To(Equal(""))
	})

	It("is successful", func() {
		server := httptest.NewServer(gtGoldenHTTPMock(commitHandlerMock(mockCommit, blankHandler)))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		var want string
		json.Unmarshal(mockCommit, &want)

		commit, err := gt.Commit()
		Expect(err).To(Succeed())
		Expect(commit).To(Equal(want))
	})
})
