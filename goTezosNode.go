package goTezos

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

type TezosRPCClient struct {
	Host        string
	Port        string
	logfunction func(level, msg string)
	logger      *log.Logger
}

//Creates a new RPC Client to query the tezos network with
func NewTezosRPCClient(hostname string, port string) *TezosRPCClient {
	t := TezosRPCClient{}
	t.Host = hostname
	t.Port = port
	t.logfunction = func(level, msg string) {
		fmt.Println(level + ": " + msg)
	}
	t.SetLogger(log.New(os.Stdout, hostname, 0))
	return &t
}

//Set the logger for the RPC Client
func (this *TezosRPCClient) SetLogger(log *log.Logger) {
	this.logger = log
}

func (this *TezosRPCClient) GetResponse(method string, path string, args string) (ResponseRaw, error) {

	url := "http://" + this.Host + this.Port + "" + path
	var jsonStr = []byte(args)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonStr))
	if err != nil {
		this.logger.Println("Error in GetResponse: " + err.Error())
		return ResponseRaw{}, err
	}

	var netTransport = &http.Transport{ // TODO make this as config option, but with defaults like this
		Dial: (&net.Dialer{
			Timeout: 3 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 3 * time.Second,
	}

	var netClient = &http.Client{
		Timeout:   time.Second * 3,
		Transport: netTransport,
	}

	resp, err := netClient.Do(req)
	if err != nil {
		this.logger.Println("Error in GetResponse: " + err.Error())
		return ResponseRaw{}, err
	}
	var b []byte
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		this.logger.Println("Error in GetResponse - readAll bytes: " + err.Error())
		return ResponseRaw{}, err
	}
	netTransport.CloseIdleConnections()
	defer resp.Body.Close()
	return ResponseRaw{b}, nil
}

//A function just to perform a query to see if an RPC Client's endpoint is alive (heartbeat)
func (this *TezosRPCClient) Healthcheck() bool {
	_, err := this.GetResponse("GET", "/chains/main/blocks", "")
	if err == nil {
		return true // healthy
	}
	return false // unhelaty
}
