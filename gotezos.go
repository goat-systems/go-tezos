package gotezos

import (
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

// MUTEZ is mutez on the tezos network
const MUTEZ = 1000000

// GoTezos is the driver of the library, it inludes the several RPC services
// like Block, SnapSHot, Cycle, Account, Delegate, Operations, Contract, and Network
type GoTezos struct {
	client    *http.Client
	Constants *Constants
	host      string
}

// RPCGenericError is an Error helper for the RPC
type genericRPCError struct {
	Kind  string `json:"kind"`
	Error string `json:"error"`
}

// RPCGenericErrors and array of RPCGenericErrors
type genericRPCErrors []genericRPCError

type params struct {
	key   string
	value string
}

// New returns a new GoTezos
func New(host string) (*GoTezos, error) {
	gt := &GoTezos{
		client: &http.Client{
			Timeout: time.Second * 10,
			Transport: &http.Transport{
				Dial: (&net.Dialer{
					Timeout: 10 * time.Second,
				}).Dial,
				TLSHandshakeTimeout: 10 * time.Second,
			},
		},
		host: cleanseHost(host),
	}

	block, err := gt.HeadBlock()
	if err != nil {
		return gt, errors.Wrap(err, "could not initialize library with network constants")
	}
	constants, err := gt.NetworkConstants(block.Hash)
	if err != nil {
		return gt, errors.Wrap(err, "could not initialize library with network constants")
	}
	gt.Constants = &constants

	return gt, nil
}

// SetClient overides GoTezos's default http client
func (t *GoTezos) SetClient(client *http.Client) {
	t.client = client
}

// SetConstants overides GoTezos's network constants
func (t *GoTezos) SetConstants(constants Constants) {
	t.Constants = &constants
}

func (t *GoTezos) post(path string, body []byte, params ...params) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s", t.host, path), bytes.NewBuffer(body))
	if err != nil {
		return nil, errors.Wrap(err, "failed to construct request")
	}

	constructQuery(req, params...)

	return t.do(req)
}

func (t *GoTezos) get(path string, params ...params) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", t.host, path), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to construct request")
	}

	constructQuery(req, params...)

	return t.do(req)
}

func (t *GoTezos) do(req *http.Request) ([]byte, error) {
	resp, err := t.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to complete request")
	}

	byts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return byts, errors.Wrap(err, "could not read response body")
	}

	if resp.StatusCode != http.StatusOK {
		return byts, fmt.Errorf("response returned code %d with body %s", resp.StatusCode, string(byts))
	}

	err = handleRPCError(byts)
	if err != nil {
		return nil, err
	}

	t.client.CloseIdleConnections()

	return byts, nil
}

func constructQuery(req *http.Request, params ...params) {
	q := req.URL.Query()
	for _, param := range params {
		q.Add(param.key, param.value)
	}

	req.URL.RawQuery = q.Encode()
}

func handleRPCError(resp []byte) error {
	if strings.Contains(string(resp), "\"error\":") {
		rpcErrors := genericRPCErrors{}
		err := json.Unmarshal(resp, &rpcErrors)
		if err != nil {
			return errors.Wrap(err, "could not unmarshal rpc error")
		}
		return fmt.Errorf("rpc error (%s): %s", rpcErrors[0].Kind, rpcErrors[0].Error)
	}
	return nil
}

func cleanseHost(host string) string {
	if host[len(host)-1] == '/' {
		host = host[:len(host)-1]
	}
	if !strings.HasPrefix(host, "http://") && !strings.HasPrefix(host, "https://") {
		host = fmt.Sprintf("http://%s", host) //default to http
	}
	return host
}
