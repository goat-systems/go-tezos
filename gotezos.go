package gotezos

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// MUTEZ is mutez on the tezos network
const MUTEZ = 1000000

var (
	regRPCError = regexp.MustCompile(`\s?\{\s?\"(kind|error)\"\s?:\s?\"[^,]+\"\s?,\s?\"(kind|error)"\s?:\s?\"[^}]+\s?\}\s?`)
)

/*
GoTezos contains a client (http.Client), network contents, and the host of the node. Gives access to
RPC related functions.
*/
type GoTezos struct {
	client           client
	networkConstants *Constants
	host             string
}

/*
RPCError represents and RPC error
*/
type RPCError struct {
	Kind string `json:"kind"`
	Err  string `json:"error"`
}

func (r *RPCError) Error() string {
	return fmt.Sprintf("rpc error (%s): %s", r.Kind, r.Err)
}

/*
RPCErrors represents multiple RPCError(s).s
*/
type RPCErrors []RPCError

type rpcOptions struct {
	Key   string
	Value string
}

type client interface {
	Do(req *http.Request) (*http.Response, error)
	CloseIdleConnections()
}

/*
New returns a pointer to a GoTezos and initializes the library with the host's Tezos netowrk constants.


Parameters:

	host:
		A Tezos node.
*/
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

	block, err := gt.Head()
	if err != nil {
		return gt, errors.Wrap(err, "could not initialize library with network constants")
	}

	constants, err := gt.Constants(block.Hash)
	if err != nil {
		return gt, errors.Wrap(err, "could not initialize library with network constants")
	}
	gt.networkConstants = constants

	return gt, nil
}

/*
SetClient overrides GoTezos's client. *http.Client satisfies the client interface.

Parameters:

	client:
		A pointer to an http.Client.
*/
func (t *GoTezos) SetClient(client *http.Client) {
	t.client = client
}

/*
SetConstants overrides GoTezos's networkConstants.

Parameters:

	constants:
		Tezos Network Constants.
*/
func (t *GoTezos) SetConstants(constants Constants) {
	t.networkConstants = &constants
}

func (t *GoTezos) post(path string, body []byte, opts ...rpcOptions) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s", t.host, path), bytes.NewBuffer(body))
	if err != nil {
		return nil, errors.Wrap(err, "failed to construct request")
	}

	constructQueryParams(req, opts...)

	return t.do(req)
}

func (t *GoTezos) get(path string, opts ...rpcOptions) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", t.host, path), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to construct request")
	}

	constructQueryParams(req, opts...)

	return t.do(req)
}

func (t *GoTezos) delete(path string, opts ...rpcOptions) ([]byte, error) {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s%s", t.host, path), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to construct request")
	}

	constructQueryParams(req, opts...)

	return t.do(req)
}

func (t *GoTezos) do(req *http.Request) ([]byte, error) {
	req.Header.Set("Content-Type", "application/json")
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
		return byts, err
	}

	t.client.CloseIdleConnections()

	return byts, nil
}

func constructQueryParams(req *http.Request, opts ...rpcOptions) {
	q := req.URL.Query()
	for _, opt := range opts {
		q.Add(opt.Key, opt.Value)
	}

	req.URL.RawQuery = q.Encode()
}

func handleRPCError(resp []byte) error {
	if regRPCError.Match(resp) {
		rpcErrors := RPCErrors{}
		err := json.Unmarshal(resp, &rpcErrors)
		if err != nil {
			return nil
		}
		return &rpcErrors[0]
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
