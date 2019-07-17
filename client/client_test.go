package client

import (
	"net/http"
	"testing"

	"gotest.tools/assert"
)

func Test_ClientPost(t *testing.T) {
	cases := []struct {
		post   string
		want   []byte
		client httpClient
	}{
		{
			post: "hard taco",
			want: []byte("hard taco"),
			client: &httpClientMock{
				ReturnStatus: http.StatusOK,
				ReturnBody:   []byte("hard taco"),
			},
		},
		{
			post: "soft taco",
			want: []byte("soft taco"),
			client: &httpClientMock{
				ReturnStatus: http.StatusOK,
				ReturnBody:   []byte("soft taco"),
			},
		},
		{
			post: "carne asadas tacos",
			want: []byte("carne asadas tacos"),
			client: &httpClientMock{
				ReturnStatus: http.StatusOK,
				ReturnBody:   []byte("carne asadas tacos"),
			},
		},
	}

	client := NewClient("http://127.0.0.1:8732")

	for _, tc := range cases {
		client.netClient = tc.client
		bytes, err := client.Post("/example/post", tc.post)
		assert.NilError(t, err)
		assert.Equal(t, string(bytes), string(tc.want))
	}
}

func Test_ClientPostBadStatus(t *testing.T) {
	cases := []struct {
		post   string
		want   []byte
		client httpClient
	}{
		{
			post: "hard taco",
			want: []byte("hard taco"),
			client: &httpClientMock{
				ReturnStatus: http.StatusInternalServerError,
			},
		},
		{
			post: "soft taco",
			want: []byte("soft taco"),
			client: &httpClientMock{
				ReturnStatus: http.StatusInternalServerError,
			},
		},
		{
			post: "carne asadas tacos",
			want: []byte("carne asadas tacos"),
			client: &httpClientMock{
				ReturnStatus: http.StatusInternalServerError,
			},
		},
	}

	client := NewClient("http://127.0.0.1:8732")
	netClient := &httpClientMock{}
	client.netClient = netClient

	for _, tc := range cases {
		_, err := client.Post("/example/post", tc.post)
		assert.Assert(t, err != nil)
	}
}

func Test_ClientGet(t *testing.T) {
	client := NewClient("http://127.0.0.1:8732")
	netClient := &httpClientMock{
		ReturnStatus: http.StatusOK,
		ReturnBody:   []byte("some GET value"),
	}
	client.netClient = netClient

	bytes, err := client.Get("/example/get", nil)
	assert.NilError(t, err)
	assert.Equal(t, string(bytes), "some GET value")
}

func Test_ClientGetBadStatus(t *testing.T) {
	client := NewClient("http://127.0.0.1:8732")
	netClient := &httpClientMock{ReturnStatus: http.StatusInternalServerError}
	client.netClient = netClient

	_, err := client.Get("/example/get", nil)
	assert.Assert(t, err != nil)
}
