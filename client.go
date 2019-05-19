package gotezos

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

// client is a struct to represent the http or rpc client
type client struct {
	URL       string
	netClient *http.Client
}

// newClient returns a new client
func newClient(URL string) *client {
	if URL[len(URL)-1] == '/' {
		URL = URL[:len(URL)-1]
	}

	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	var netClient = &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}

	return &client{URL: URL, netClient: netClient}
}

func (c *client) Post(path, args string) ([]byte, error) {
	var respBytes []byte
	resp, err := c.netClient.Post(c.URL+path, "application/json", bytes.NewBuffer([]byte(args)))
	if err != nil {
		return respBytes, err
	}

	if resp.StatusCode != http.StatusOK {
		return respBytes, fmt.Errorf("bad response code %d", resp.StatusCode)
	}

	respBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return respBytes, err
	}

	c.netClient.CloseIdleConnections()

	return respBytes, nil
}

func (c *client) Get(path string, params map[string]string) ([]byte, error) {
	var bytes []byte

	req, err := http.NewRequest("GET", c.URL+path, nil)
	if err != nil {
		return bytes, err
	}

	q := req.URL.Query()
	if len(params) > 0 {
		for k, v := range params {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	resp, err := c.netClient.Get(c.URL + path)
	if err != nil {
		return bytes, err
	}

	if resp.StatusCode != http.StatusOK {
		return bytes, fmt.Errorf("bad response code %d", resp.StatusCode)
	}

	bytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return bytes, err
	}

	c.netClient.CloseIdleConnections()

	return bytes, nil
}
