package gotezos

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

// TezosRPCClient is a struct to represent the client to reach a Tezos node
type TezosRPCClient struct {
	Host        string
	Port        string
	logfunction func(level, msg string)
	logger      *log.Logger
	isWebClient bool
	httpClient  *http.Client
}

// NewTezosRPCClient creates a new RPC client using the specified hostname and port.
// Also acceptable is the hostname of a web-endpoint that supports https.
func NewTezosRPCClient(hostname string, port string) *TezosRPCClient {
	t := TezosRPCClient{}

	// Strip off posible trailing '/'
	hLen := len(hostname)
	if hostname[hLen-1] == '/' {
		hostname = hostname[:hLen-1]
	}

	// Strip off URI scheme
	if hostname[:8] == "https://" {
		hostname = hostname[8:]
		t.isWebClient = true
	} else if hostname[:7] == "http://" {
		hostname = hostname[7:]
	}

	t.Host = hostname
	t.Port = port
	t.logfunction = func(level, msg string) {
		fmt.Println(level + ": " + msg)
	}
	t.SetLogger(log.New(os.Stdout, hostname, 0))

	var netTransport = &http.Transport{ // TODO make gt as config option, but with defaults like this
		Dial: (&net.Dialer{
			Timeout: 3 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 3 * time.Second,
	}

	t.httpClient = &http.Client{
		Timeout:   time.Second * 3,
		Transport: netTransport,
	}

	return &t
}

// SetLogger set the logger for the RPC Client
func (gt *TezosRPCClient) SetLogger(log *log.Logger) {
	gt.logger = log
}

// IsWebClient tells the TezosRPCClient calling it that it is a web client
func (gt *TezosRPCClient) IsWebClient(b bool) {
	gt.isWebClient = b
}

// GetResponse gets the raw response using TezosRPCClient with the path and args to query
func (gt *TezosRPCClient) GetResponse(method string, path string, args string) (ResponseRaw, error) {

	var url string

	if gt.isWebClient {
		url = fmt.Sprintf("https://%s:%s%s", gt.Host, gt.Port, path)
	} else {
		url = fmt.Sprintf("http://%s:%s%s", gt.Host, gt.Port, path)
	}

	var jsonStr = []byte(args)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonStr))
	if err != nil {
		gt.logger.Println("Error in GetResponse: " + err.Error())
		return ResponseRaw{}, err
	}

	resp, err := gt.httpClient.Do(req)
	if err != nil {
		gt.logger.Println("Error in GetResponse: " + err.Error())
		return ResponseRaw{}, err
	}
	var b []byte
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		gt.logger.Println("Error in GetResponse - readAll bytes: " + err.Error())
		return ResponseRaw{}, err
	}
	defer resp.Body.Close()
	return ResponseRaw{b}, nil
}

// Healthcheck a function just to perform a query to see if an RPC Client's endpoint is alive (heartbeat)
func (gt *TezosRPCClient) Healthcheck() bool {
	_, err := gt.GetResponse("GET", "/chains/main/blocks", "")
	if err == nil {
		return true // healthy
	}
	return false // unhelaty
}
