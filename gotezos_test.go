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
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			if strings.Contains(req.URL.String(), "/chains/main/blocks/head") {
				rw.Write(mockBlockRandom)
			}
			if strings.Contains(req.URL.String(), "/context/constants") {
				rw.Write(mockConstants)
			}
		}))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())
		Expect(gt.Constants).To(Equal(&expectedConstants))
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
		Expect(gt.Constants).To(BeNil())
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
		Expect(gt.Constants).To(BeNil())
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

		Expect(gt.Constants).To(Equal(&constants))
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

		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			if strings.Contains(req.URL.String(), "/chains/main/blocks/head") {
				rw.Write(mockBlockRandom)
			} else if strings.Contains(req.URL.String(), "/context/constants") {
				rw.Write(mockConstants)
			} else {
				Expect(req.Method).To(Equal(http.MethodPost))
				body, _ := ioutil.ReadAll(req.Body)
				Expect(body).To(Equal(body))
				Expect(req.URL.String()).To(Equal(post))
				rw.Write(want)
			}
		}))
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

		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			if strings.Contains(req.URL.String(), "/chains/main/blocks/head") {
				rw.Write(mockBlockRandom)
			} else if strings.Contains(req.URL.String(), "/context/constants") {
				rw.Write(mockConstants)
			} else {
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
			}
		}))
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

		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			if strings.Contains(req.URL.String(), "/chains/main/blocks/head") {
				rw.Write(mockBlockRandom)
			} else if strings.Contains(req.URL.String(), "/context/constants") {
				rw.Write(mockConstants)
			} else {
				Expect(req.Method).To(Equal(http.MethodGet))
				Expect(req.URL.String()).To(Equal(post))
				rw.Write(want)
			}
		}))
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

		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			if strings.Contains(req.URL.String(), "/chains/main/blocks/head") {
				rw.Write(mockBlockRandom)
			} else if strings.Contains(req.URL.String(), "/context/constants") {
				rw.Write(mockConstants)
			} else {
				Expect(req.Method).To(Equal(http.MethodGet))
				Expect(req.URL.String()).To(Equal("/some/endpoint?my_key=my_val&other_key=other_val"))

				firstKey, ok := req.URL.Query()[params[0].key]
				Expect(ok).To(Equal(true))
				Expect(firstKey).To(Equal([]string{params[0].value}))

				secondKey, ok := req.URL.Query()[params[1].key]
				Expect(ok).To(Equal(true))
				Expect(secondKey).To(Equal([]string{params[1].value}))

				rw.Write(want)
			}
		}))
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

		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			if strings.Contains(req.URL.String(), "/chains/main/blocks/head") {
				rw.Write(mockBlockRandom)
			} else if strings.Contains(req.URL.String(), "/context/constants") {
				rw.Write(mockConstants)
			} else {
				Expect(req.Method).To(Equal(method))
				Expect(req.URL.String()).To(Equal(path))
				rw.Write(want)
			}
		}))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		req, err := http.NewRequest(method, fmt.Sprintf("%s%s", server.URL, path), nil)
		Expect(err).To(Succeed())

		p, err := gt.do(req)
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

		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			if strings.Contains(req.URL.String(), "/chains/main/blocks/head") {
				rw.Write(mockBlockRandom)
			} else if strings.Contains(req.URL.String(), "/context/constants") {
				rw.Write(mockConstants)
			} else {
				Expect(req.Method).To(Equal(http.MethodGet))
				Expect(req.URL.String()).To(Equal("/some/endpoint?my_key=my_val&other_key=other_val"))

				firstKey, ok := req.URL.Query()[params[0].key]
				Expect(ok).To(Equal(true))
				Expect(firstKey).To(Equal([]string{params[0].value}))

				secondKey, ok := req.URL.Query()[params[1].key]
				Expect(ok).To(Equal(true))
				Expect(secondKey).To(Equal([]string{params[1].value}))

				rw.Write(want)
			}
		}))
		defer server.Close()

		gt, err := New(server.URL)
		Expect(err).To(Succeed())

		p, err := gt.get(post, params...)
		Expect(err).To(Succeed())
		Expect(p).To(Equal(want))
	})
})
