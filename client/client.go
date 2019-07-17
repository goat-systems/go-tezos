package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Client is a struct to represent the http or rpc client
type Client struct {
	URL       string
	netClient httpClient
}

// RPCGenericError is an Error helper for the RPC
type genericRPCError struct {
	Kind  string `json:"kind"`
	Error string `json:"error"`
}

// RPCGenericErrors and array of RPCGenericErrors
type genericRPCErrors []genericRPCError

// NewClient returns a new client
func NewClient(URL string) *Client {
	if URL[len(URL)-1] == '/' {
		URL = URL[:len(URL)-1]
	}
	if !strings.HasPrefix(URL, "http://") && !strings.HasPrefix(URL, "https://") {
		URL = fmt.Sprintf("http://%s", URL) //default to http
	}

	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	var netClient = &http.Client{
		Timeout:   10 * time.Second,
		Transport: netTransport,
	}

	return &Client{URL: URL, netClient: netClient}
}

func (c *Client) Post(path, args string) ([]byte, error) {
	var respBytes []byte
	resp, err := c.netClient.Post(c.URL+path, "application/json", bytes.NewBuffer([]byte(args)))
	if err != nil {
		return respBytes, err
	}

	respBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return respBytes, errors.Wrap(err, "could not post")
	}

	if resp.StatusCode != http.StatusOK {
		return respBytes, errors.Errorf("%d error: %s", resp.StatusCode, string(respBytes))
	}

	err = c.handleRPCError(respBytes)
	if err != nil {
		return nil, err
	}

	c.netClient.CloseIdleConnections()

	return respBytes, nil
}

func (c *Client) Get(path string, params map[string]string) ([]byte, error) {
	var bytes []byte

	req, err := http.NewRequest("GET", c.URL+path, nil)
	if err != nil {
		return bytes, err
	}
	req.Header.Set("Accept-Encoding", "utf8")

	q := req.URL.Query()
	if len(params) > 0 {
		for k, v := range params {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	resp, err := c.netClient.Do(req)
	if err != nil {
		return bytes, err
	}

	bytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return bytes, err
	}

	if resp.StatusCode != http.StatusOK {
		return bytes, errors.Errorf("%d error: %s", resp.StatusCode, string(bytes))
	}

	err = c.handleRPCError(bytes)
	if err != nil {
		return nil, err
	}

	c.netClient.CloseIdleConnections()

	return bytes, nil
}

func (c *Client) StreamGet(path string, params map[string]string,
	res chan []byte, errc chan error, done chan bool) {

	req, err := http.NewRequest("GET", c.URL+path, nil)
	if err != nil {
		errc <- errors.Wrapf(err, "could not create request '%s'", path)
		res <- []byte{}
		return
	}
	// VERY IMPORTANT ::: This one was a pain in my ass -.-
	req.Header.Set("Accept-Encoding", "utf8")

	q := req.URL.Query()
	if len(params) > 0 {
		for k, v := range params {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		errc <- errors.Wrapf(err, "could not start request '%s'", path)
		res <- []byte{}
		return
	}

	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	for {
		select {
		case <-done:
			errc <- nil
			res <- []byte{}
			return
		default:
			response, err := reader.ReadBytes('\n')
			if err != nil {
				errc <- errors.Wrapf(err, "could not read response from '%s'", path)
				res <- []byte{}
				return
			}
			errc <- nil
			res <- bytes.TrimSpace(response)
		}
	}
}

func (c *Client) handleRPCError(resp []byte) error {
	if strings.Contains(string(resp), "\"error\":") {
		rpcErrors := genericRPCErrors{}
		rpcErrors, err := rpcErrors.unmarshalJSON(resp)
		if err != nil {
			return err
		}
		return errors.Errorf("rpc error (%s): %s", rpcErrors[0].Kind, rpcErrors[0].Error)
	}
	return nil
}

// UnmarshalJSON unmarhsels bytes into RPCGenericErrors
func (ge *genericRPCErrors) unmarshalJSON(v []byte) (genericRPCErrors, error) {
	r := genericRPCErrors{}

	err := json.Unmarshal(v, &r)
	if err != nil {
		return r, errors.Wrap(err, "could not unmarshal bytes into genericRPCErrors")
	}
	return r, nil
}
