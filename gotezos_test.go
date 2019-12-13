package gotezos

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var suiteName = "GoTezos"

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, fmt.Sprintf("%s Suite", suiteName))
}

var _ = Describe("New", func() {
	var expectedConstants Constants

	BeforeEach(func() {
		err := json.Unmarshal(mockConstants, &expectedConstants)
		Expect(err).To(Succeed())
	})

	It("is successful", func() {
		server := httptest.NewServer(gtGoldenHTTPMock(blankHandler))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())
		Expect(gt.networkConstants).To(Equal(&expectedConstants))
	})

	It("fails to fetch head block to get constants", func() {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			if strings.Contains(req.URL.String(), "/chains/main/blocks/head") {
				rw.Write([]byte(`some_junk_data`))
			}
		}))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).NotTo(BeNil())
		Expect(gt.networkConstants).To(BeNil())
	})

	It("fails to fetch constants", func() {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			if strings.Contains(req.URL.String(), "/chains/main/blocks/head") {
				rw.Write(mockBlockRandom)
			}
			if strings.Contains(req.URL.String(), "/context/constants") {
				rw.Write([]byte(`some_junk_data`))
			}
		}))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).NotTo(BeNil())
		Expect(gt.networkConstants).To(BeNil())
	})
})

var _ = Describe("GoTezos.SetClient", func() {
	It("sets the client", func() {
		gt := GoTezos{}

		client := &http.Client{}
		gt.SetClient(client)

		Expect(gt.client).To(Equal(client))
	})
})

var _ = Describe("GoTezos.SetConstants", func() {
	It("sets the client", func() {
		gt := GoTezos{}

		var constants Constants
		gt.SetConstants(constants)

		Expect(gt.networkConstants).To(Equal(&constants))
	})
})

var _ = Describe("GoTezos.post", func() {
	var expectedConstants Constants

	BeforeEach(func() {
		err := json.Unmarshal(mockConstants, &expectedConstants)
		Expect(err).To(Succeed())
	})

	It("posts basic request", func() {
		post := "/some/endpoint"
		body := []byte("some_body")
		want := []byte("success")

		server := httptest.NewServer(gtGoldenHTTPMock(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			Expect(r.Method).To(Equal(http.MethodPost))
			body, _ := ioutil.ReadAll(r.Body)
			Expect(body).To(Equal(body))
			Expect(r.URL.String()).To(Equal(post))
			w.Write(want)
		})))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		p, err := gt.post(post, body)
		Expect(err).To(Succeed())
		Expect(p).To(Equal(want))
	})

	It("posts with parameters", func() {
		post := "/some/endpoint"
		body := []byte("some_body")
		want := []byte("success")
		params := []params{
			{
				key:   "my_key",
				value: "my_val",
			},
			{
				key:   "other_key",
				value: "other_val",
			},
		}

		server := httptest.NewServer(gtGoldenHTTPMock(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			Expect(req.Method).To(Equal(http.MethodPost))
			body, _ := ioutil.ReadAll(req.Body)
			Expect(body).To(Equal(body))
			Expect(req.URL.String()).To(Equal("/some/endpoint?my_key=my_val&other_key=other_val"))

			firstKey, ok := req.URL.Query()[params[0].key]
			Expect(ok).To(Equal(true))
			Expect(firstKey).To(Equal([]string{params[0].value}))

			secondKey, ok := req.URL.Query()[params[1].key]
			Expect(ok).To(Equal(true))
			Expect(secondKey).To(Equal([]string{params[1].value}))

			rw.Write(want)
		})))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		p, err := gt.post(post, body, params...)
		Expect(err).To(Succeed())
		Expect(p).To(Equal(want))
	})
})

var _ = Describe("GoTezos.get", func() {
	var expectedConstants Constants

	BeforeEach(func() {
		err := json.Unmarshal(mockConstants, &expectedConstants)
		Expect(err).To(Succeed())
	})

	It("gets basic request", func() {
		post := "/some/endpoint"
		want := []byte("success")

		server := httptest.NewServer(gtGoldenHTTPMock(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			Expect(req.Method).To(Equal(http.MethodGet))
			Expect(req.URL.String()).To(Equal(post))
			rw.Write(want)
		})))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		p, err := gt.get(post)
		Expect(err).To(Succeed())
		Expect(p).To(Equal(want))
	})

	It("gets with parameters", func() {
		post := "/some/endpoint"
		want := []byte("success")
		params := []params{
			{
				key:   "my_key",
				value: "my_val",
			},
			{
				key:   "other_key",
				value: "other_val",
			},
		}

		server := httptest.NewServer(gtGoldenHTTPMock(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			Expect(req.Method).To(Equal(http.MethodGet))
			Expect(req.URL.String()).To(Equal("/some/endpoint?my_key=my_val&other_key=other_val"))

			firstKey, ok := req.URL.Query()[params[0].key]
			Expect(ok).To(Equal(true))
			Expect(firstKey).To(Equal([]string{params[0].value}))

			secondKey, ok := req.URL.Query()[params[1].key]
			Expect(ok).To(Equal(true))
			Expect(secondKey).To(Equal([]string{params[1].value}))

			rw.Write(want)
		})))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		p, err := gt.get(post, params...)
		Expect(err).To(Succeed())
		Expect(p).To(Equal(want))
	})
})

var _ = Describe("GoTezos.do", func() {
	var expectedConstants Constants

	BeforeEach(func() {
		err := json.Unmarshal(mockConstants, &expectedConstants)
		Expect(err).To(Succeed())
	})

	It("is successful", func() {
		method := http.MethodGet
		path := "/some/endpoint"
		want := []byte("success")

		server := httptest.NewServer(gtGoldenHTTPMock(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			Expect(req.Method).To(Equal(method))
			Expect(req.URL.String()).To(Equal(path))
			rw.Write(want)
		})))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		req, err := http.NewRequest(method, fmt.Sprintf("%s%s", server.URL, path), nil)
		Expect(err).To(Succeed())

		p, err := gt.do(req)
		Expect(err).To(Succeed())
		Expect(p).To(Equal(want))
	})

	It("returns errors if not 200 OK", func() {
		post := "/some/endpoint"
		want := []byte("fail")

		server := httptest.NewServer(gtGoldenHTTPMock(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			Expect(req.Method).To(Equal(http.MethodGet))
			Expect(req.URL.String()).To(Equal(post))

			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write(want)
		})))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", server.URL, post), nil)
		Expect(err).To(Succeed())

		p, err := gt.do(req)
		Expect(err).ToNot(BeNil())
		Expect(p).To(Equal(want))
	})

	It("returns rpc error", func() {
		post := "/some/endpoint"
		want := []byte(`[{"kind":"somekind","Error":"someerror"}]`)

		server := httptest.NewServer(gtGoldenHTTPMock(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			Expect(req.Method).To(Equal(http.MethodGet))
			Expect(req.URL.String()).To(Equal(post))

			rw.Write(want)
		})))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", server.URL, post), nil)
		Expect(err).To(Succeed())

		p, err := gt.do(req)
		Expect(err).ToNot(BeNil())
		Expect(p).To(Equal(want))
	})
})

var _ = Describe("handleRPCError", func() {
	It("found an rpc error", func() {
		resp := []byte(`[{"kind":"some_kind","error":"some_error"}]`)
		err := handleRPCError(resp)
		Expect(err).ToNot(BeNil())
		Expect(err.Error()).To(MatchRegexp("rpc error"))
	})

	It("failed to unmarshal rpc error", func() {
		resp := []byte(`error`)
		err := handleRPCError(resp)
		Expect(err).ToNot(BeNil())
		Expect(err.Error()).To(MatchRegexp("could not unmarshal rpc error"))
	})

	It("did not find an rpc error", func() {
		resp := []byte(`some other data`)
		err := handleRPCError(resp)
		Expect(err).To(BeNil())
	})
})

var _ = Describe("handleRPCError", func() {
	It("found an rpc error", func() {
		resp := []byte(`[{"kind":"some_kind","error":"some_error"}]`)
		err := handleRPCError(resp)
		Expect(err).ToNot(BeNil())
		Expect(err.Error()).To(MatchRegexp("rpc error"))
	})

	It("failed to unmarshal rpc error", func() {
		resp := []byte(`error`)
		err := handleRPCError(resp)
		Expect(err).ToNot(BeNil())
		Expect(err.Error()).To(MatchRegexp("could not unmarshal rpc error"))
	})

	It("did not find an rpc error", func() {
		resp := []byte(`some other data`)
		err := handleRPCError(resp)
		Expect(err).To(BeNil())
	})
})

var _ = Describe("constructQuery", func() {
	It("adds url parameters to http request", func() {
		req, err := http.NewRequest(http.MethodGet, "www.someurl.com/some/request", nil)
		Expect(err).To(Succeed())

		params := []params{
			{
				key:   "key",
				value: "val",
			},
			{
				key:   "key1",
				value: "val1",
			},
		}

		constructQueryParams(req, params...)
		Expect(req.URL.Query().Get(params[0].key)).To(Equal(params[0].value))
		Expect(req.URL.Query().Get(params[1].key)).To(Equal(params[1].value))
	})

	It("handles no parameters", func() {
		req, err := http.NewRequest(http.MethodGet, "www.someurl.com/some/request", nil)
		Expect(err).To(Succeed())
		constructQueryParams(req)
		Expect(req.URL.Query().Encode()).To(Equal(""))
	})
})

var _ = Describe("cleanseHost", func() {
	It("strips trailing / if missing", func() {
		host := cleanseHost("http://www.host.com/")
		Expect(host).To(Equal("http://www.host.com"))
	})

	It("handles missing http(s)://", func() {
		host := cleanseHost("www.host.com")
		Expect(host).To(Equal("http://www.host.com"))
	})
})

func gtGoldenHTTPMock(next http.Handler) http.Handler {
	var blockMock blockMock
	return constantsHandlerMock(
		mockConstants,
		blockMock.handler(
			mockBlockRandom,
			next,
		),
	)
}
