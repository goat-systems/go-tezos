package client

import (
	"io"
	"net/http"
)

type TezosClient interface {
	Post(path, args string) ([]byte, error)
	Get(path string, params map[string]string) ([]byte, error)
	StreamGet(path string, params map[string]string,
		res chan []byte, errc chan error, done chan bool)
}

// httpClient is an interface that exposes the HTTP methods for testing.
type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
	Post(url, contentType string, body io.Reader) (*http.Response, error)
	CloseIdleConnections()
}
