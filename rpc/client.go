package rpc

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

// MUTEZ is mutez on the tezos network
const MUTEZ = 1000000

var (
	regRPCError = regexp.MustCompile(`\s?\{\s?\"(kind|error)\"\s?:\s?\"[^,]+\"\s?,\s?\"(kind|error)"\s?:\s?\"[^}]+\s?\}\s?`)
)

/*
Client contains a client (http.Client), network contents, and the host of the node. Gives access to
RPC related functions.
*/
type Client struct {
	client           *resty.Client
	chain            string
	networkConstants *Constants
	host             string
}

/*
Error represents and RPC error
*/
type Error struct {
	Kind string `json:"kind"`
	Err  string `json:"error"`
}

func (r *Error) Error() string {
	return fmt.Sprintf("rpc error (%s): %s", r.Kind, r.Err)
}

/*
Errors represents multiple RPCError(s).s
*/
type Errors []Error

type rpcOptions struct {
	Key   string
	Value string
}

func queryParams(options ...rpcOptions) map[string]string {
	m := make(map[string]string)
	for _, opt := range options {
		m[opt.Key] = opt.Value
	}
	return m
}

/*
New returns a pointer to a Client and initializes the rpc configuration with the host's Tezos netowrk constants.
*/
func New(host string) (*Client, error) {
	c := &Client{
		client: resty.New(),
		host:   cleanseHost(host),
		chain:  "main",
	}

	_, constants, err := c.Constants(ConstantsInput{BlockID: &BlockIDHead{}})
	if err != nil {
		return c, errors.Wrap(err, "failed to initialize library with network constants")
	}
	c.networkConstants = &constants

	return c, nil
}

// SetChain sets the chain for the rpc
func (c *Client) SetChain(chain string) {
	c.chain = chain
}

// CurrentContstants returns the constants used on the client
func (c *Client) CurrentContstants() Constants {
	return *c.networkConstants
}

/*
OverrideClient overrides underlying network client.
Can allow you to create middleware as needed: https://github.com/go-resty/resty#request-and-response-middleware
*/
func (c *Client) OverrideClient(client *resty.Client) {
	c.client = client
}

/*
SetConstants overrides GoTezos's networkConstants.
*/
func (c *Client) SetConstants(constants Constants) {
	c.networkConstants = &constants
}

func (c *Client) post(path string, body interface{}, opts ...rpcOptions) (*resty.Response, error) {
	resp, err := c.client.R().
		SetQueryParams(queryParams(opts...)).
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(fmt.Sprintf("%s%s", c.host, path))

	if err != nil {
		return resp, err
	}

	err = handleRPCError(resp.Body())

	return resp, err
}

func (c *Client) get(path string, opts ...rpcOptions) (*resty.Response, error) {
	resp, err := c.client.R().SetQueryParams(queryParams(opts...)).Get(fmt.Sprintf("%s%s", c.host, path))
	if err != nil {
		return resp, err
	}
	err = handleRPCError(resp.Body())

	return resp, err
}

func handleRPCError(resp []byte) error {
	if regRPCError.Match(resp) {
		rpcErrors := Errors{}
		err := json.Unmarshal(resp, &rpcErrors)
		if err != nil {
			return nil
		}
		return &rpcErrors[0]
	}

	return nil
}

func cleanseHost(host string) string {
	if len(host) == 0 {
		return ""
	}
	if host[len(host)-1] == '/' {
		host = host[:len(host)-1]
	}
	if !strings.HasPrefix(host, "http://") && !strings.HasPrefix(host, "https://") {
		host = fmt.Sprintf("http://%s", host) //default to http
	}
	return host
}
