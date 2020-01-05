package gotezos

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	expectedConstants := expectedConstants(t)

	cases := []struct {
		name          string
		inputHandler  http.Handler
		wantErr       bool
		wantConstants *Constants
	}{
		{
			"Successful",
			gtGoldenHTTPMock(blankHandler),
			false,
			expectedConstants,
		},
		{
			"fails to fetch head block to get constants",
			http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				if strings.Contains(req.URL.String(), "/chains/main/blocks/head") {
					rw.Write([]byte(`some_junk_data`))
				}
			}),
			true,
			nil,
		},
		{
			"fails to fetch constants",
			http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				if strings.Contains(req.URL.String(), "/chains/main/blocks/head") {
					rw.Write(mockBlockRandom)
				}
				if strings.Contains(req.URL.String(), "/context/constants") {
					rw.Write([]byte(`some_junk_data`))
				}
			}),
			true,
			nil,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHandler)
			defer server.Close()

			gt, err := New(server.URL)

			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, tt.wantConstants, gt.networkConstants)
		})
	}
}

func Test_SetClient(t *testing.T) {
	gt := GoTezos{}

	client := &http.Client{}
	gt.SetClient(client)

	assert.Equal(t, client, gt.client)
}

func Test_SetConstants(t *testing.T) {
	gt := GoTezos{}

	var constants Constants
	gt.SetConstants(constants)

	assert.Equal(t, constants, *gt.networkConstants)
}

func Test_post(t *testing.T) {
	type input struct {
		handler http.Handler
		body    []byte
		post    string
		params  []params
	}

	type want struct {
		err  bool
		resp []byte
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"posts basic request",
			input{
				gtGoldenHTTPMock(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, http.MethodPost, r.Method)
					body, _ := ioutil.ReadAll(r.Body)
					assert.Equal(t, []byte("some_body"), body)
					assert.Equal(t, "/some/endpoint", r.URL.String())
					w.Write([]byte("success"))
				})),
				[]byte("some_body"),
				"/some/endpoint",
				nil,
			},
			want{
				false,
				[]byte("success"),
			},
		},
		{
			"posts with parameters",
			input{
				gtGoldenHTTPMock(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, http.MethodPost, r.Method)
					body, _ := ioutil.ReadAll(r.Body)
					assert.Equal(t, []byte("some_body"), body)
					assert.Equal(t, "/some/endpoint?my_key=my_val&other_key=other_val", r.URL.String())

					firstKey, ok := r.URL.Query()["my_key"]
					assert.True(t, ok)
					assert.Equal(t, []string{"my_val"}, firstKey)

					secondKey, ok := r.URL.Query()["other_key"]
					assert.True(t, ok)
					assert.Equal(t, []string{"other_val"}, secondKey)

					w.Write([]byte("success"))
				})),
				[]byte("some_body"),
				"/some/endpoint",
				[]params{
					{
						key:   "my_key",
						value: "my_val",
					},
					{
						key:   "other_key",
						value: "other_val",
					},
				},
			},
			want{
				false,
				[]byte("success"),
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input.handler)
			defer server.Close()

			gt, err := New(server.URL)
			assert.Nil(t, err)

			p, err := gt.post(tt.input.post, tt.input.body, tt.input.params...)
			if tt.want.err {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, tt.want.resp, p)
		})
	}
}

func Test_get(t *testing.T) {
	type input struct {
		handler http.Handler
		get     string
		params  []params
	}

	type want struct {
		err  bool
		resp []byte
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"gets basic request",
			input{
				gtGoldenHTTPMock(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, http.MethodGet, r.Method)
					assert.Equal(t, "/some/endpoint", r.URL.String())
					w.Write([]byte("success"))
				})),
				"/some/endpoint",
				nil,
			},
			want{
				false,
				[]byte("success"),
			},
		},
		{
			"gets with parameters",
			input{
				gtGoldenHTTPMock(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, http.MethodGet, r.Method)
					assert.Equal(t, "/some/endpoint?my_key=my_val&other_key=other_val", r.URL.String())

					firstKey, ok := r.URL.Query()["my_key"]
					assert.True(t, ok)
					assert.Equal(t, []string{"my_val"}, firstKey)

					secondKey, ok := r.URL.Query()["other_key"]
					assert.True(t, ok)
					assert.Equal(t, []string{"other_val"}, secondKey)

					w.Write([]byte("success"))
				})),
				"/some/endpoint",
				[]params{
					{
						key:   "my_key",
						value: "my_val",
					},
					{
						key:   "other_key",
						value: "other_val",
					},
				},
			},
			want{
				false,
				[]byte("success"),
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input.handler)
			defer server.Close()

			gt, err := New(server.URL)
			assert.Nil(t, err)

			p, err := gt.get(tt.input.get, tt.input.params...)
			if tt.want.err {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, tt.want.resp, p)
		})
	}
}

func Test_do(t *testing.T) {
	type input struct {
		handler http.Handler
		method  string
		path    string
	}

	type want struct {
		err  bool
		resp []byte
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"is successful",
			input{
				gtGoldenHTTPMock(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, http.MethodGet, r.Method)
					assert.Equal(t, "/some/endpoint", r.URL.String())
					w.Write([]byte("success"))
				})),
				http.MethodGet,
				"/some/endpoint",
			},
			want{
				false,
				[]byte("success"),
			},
		},
		{
			"returns errors if not 200 OK",
			input{
				gtGoldenHTTPMock(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, http.MethodGet, r.Method)
					assert.Equal(t, "/some/endpoint", r.URL.String())

					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("fail"))
				})),
				http.MethodGet,
				"/some/endpoint",
			},
			want{
				true,
				[]byte("fail"),
			},
		},
		{
			"returns rpc error",
			input{
				gtGoldenHTTPMock(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, http.MethodGet, r.Method)
					assert.Equal(t, "/some/endpoint", r.URL.String())

					w.Write(mockRPCError)
				})),
				http.MethodGet,
				"/some/endpoint",
			},
			want{
				true,
				mockRPCError,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input.handler)
			defer server.Close()

			gt, err := New(server.URL)
			assert.Nil(t, err)

			req, err := http.NewRequest(tt.input.method, fmt.Sprintf("%s%s", server.URL, tt.input.path), nil)
			assert.Nil(t, err)

			p, err := gt.do(req)
			if tt.want.err {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, tt.want.resp, p)
		})
	}
}

func Test_handleRPCError(t *testing.T) {
	cases := []struct {
		name        string
		resp        []byte
		wantErr     bool
		errContents string
	}{
		{
			"found an rpc error",
			[]byte(`[{"kind":"some_kind","error":"some_error"}]`),
			true,
			"rpc error",
		},
		{
			"failed to unmarshal rpc error",
			[]byte(`error`),
			true,
			"could not unmarshal rpc error",
		},
		{
			"did not find an rpc error",
			[]byte(`some other data`),
			false,
			"",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := handleRPCError(tt.resp)
			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tt.errContents)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func Test_constructQuery(t *testing.T) {
	cases := []struct {
		name   string
		params []params
	}{
		{
			"adds url parameters to http request",
			[]params{
				{
					key:   "key",
					value: "val",
				},
				{
					key:   "key1",
					value: "val1",
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "www.someurl.com/some/request", nil)
			assert.Nil(t, err)

			constructQueryParams(req, tt.params...)
			if tt.params != nil {
				assert.Equal(t, tt.params[0].value, req.URL.Query().Get(tt.params[0].key))
				assert.Equal(t, tt.params[1].value, req.URL.Query().Get(tt.params[1].key))
			} else {
				assert.Equal(t, "", req.URL.Query().Encode())
			}
		})
	}
}

func Test_cleanseHost(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{
			"strips trailing / if missing",
			"http://www.host.com/",
			"http://www.host.com",
		},
		{
			"handles missing http(s)://",
			"www.host.com",
			"http://www.host.com",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			host := cleanseHost(tt.input)
			assert.Equal(t, tt.want, host)
		})
	}
}

func gtGoldenHTTPMock(next http.Handler) http.Handler {
	var blockMock blockMock
	var constantsMock constantsMock
	return constantsMock.handler(
		mockConstants,
		blockMock.handler(
			mockBlockRandom,
			next,
		),
	)
}

func expectedConstants(t *testing.T) *Constants {
	var expectedConstants Constants
	err := json.Unmarshal(mockConstants, &expectedConstants)
	assert.Nilf(t, err, "could no unmarhsal mock constants")

	return &expectedConstants
}
