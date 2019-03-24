package gotezos

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"strings"

	"github.com/patrickmn/go-cache"
)

// NewGoTezos is a constructor that returns a GoTezos object
func NewGoTezos() *GoTezos {
	a := GoTezos{}

	a.UseBalancerStrategyFailover()
	a.rand = rand.New(rand.NewSource(time.Now().Unix()))
	go func(a *GoTezos) {
		for {
			time.Sleep(15 * time.Second)
			a.checkUnhealthyClients()
		}
	}(&a)

	// TTL Cache
	// 5s default cache, 5m garbage collection
	a.cache = cache.New(5*time.Second, 5*time.Minute)

	// Default logger
	a.logger = log.New(os.Stderr, "", log.LstdFlags)

	return &a
}

// SetLogger sets the logging functionality
func (gt *GoTezos) SetLogger(log *log.Logger) {
	gt.logger = log
}

// Debug puts go-tezos into debug mode for logging
func (gt *GoTezos) Debug(d bool) {
	gt.debug = d
}

// AddNewClient adds an RPC Client to query the tezos network
func (gt *GoTezos) AddNewClient(client *TezosRPCClient) {

	gt.clientLock.Lock()
	gt.RPCClients = append(gt.RPCClients, &TezClientWrapper{true, client})
	gt.clientLock.Unlock()

	var err error
	gt.Constants, err = gt.GetNetworkConstants()
	if err != nil {
		fmt.Println("Could not get network constants, library will fail. Exiting .... ")
		os.Exit(0)
	}

	gt.Versions, err = gt.GetNetworkVersions()
	if err != nil {
		fmt.Println("Could not get network version, library will fail. Exiting .... ")
		os.Exit(0)
	}
}

// IsMainnet checks whether the current network used is the Mainnet
func (gt *GoTezos) IsMainnet() bool {
	if len(gt.Versions) > 0 {
		return gt.Versions[0].Network == "BETANET"
	}
	return false
}

// IsAlphanet checks whether the current network used is the IsAlphanet
func (gt *GoTezos) IsAlphanet() bool {
	if len(gt.Versions) > 0 {
		return gt.Versions[0].Network == "ALPHANET"
	}
	return false
}

// IsZeronet checks whether the current network used is the IsZeronet
func (gt *GoTezos) IsZeronet() bool {
	if len(gt.Versions) > 0 {
		return gt.Versions[0].Network == "ZERONET"
	}
	return false
}

// UseBalancerStrategyFailover tells the client side failover to use the balancer strategy
func (gt *GoTezos) UseBalancerStrategyFailover() {
	gt.balancerStrategy = "failover"
}

// UseBalancerStrategyRandom tells the client side failover to use the balancer random strategy
func (gt *GoTezos) UseBalancerStrategyRandom() {
	gt.balancerStrategy = "random"
}

func (gt *GoTezos) checkHealthStatus() {
	gt.clientLock.Lock()
	wg := sync.WaitGroup{}
	for _, a := range gt.RPCClients {
		wg.Add(1)
		go func(wg *sync.WaitGroup, client *TezClientWrapper) {
			res := client.client.Healthcheck()
			if client.healthy && res == false {
				gt.logger.Println("Client state switched to unhealthy", gt.ActiveRPCCient.client.Host+gt.ActiveRPCCient.client.Port)
			}
			if !client.healthy && res {
				gt.logger.Println("Client state switched to healthy", gt.ActiveRPCCient.client.Host+gt.ActiveRPCCient.client.Port)
			}
			client.healthy = res
			wg.Done()
		}(&wg, a)
	}
	wg.Wait()
	gt.clientLock.Unlock()
}

func (gt *GoTezos) checkUnhealthyClients() {
	gt.clientLock.Lock()
	wg := sync.WaitGroup{}
	for _, a := range gt.RPCClients {
		wg.Add(1)
		go func(wg *sync.WaitGroup, client *TezClientWrapper) {
			if client.healthy == false {
				res := client.client.Healthcheck()
				if !client.healthy && res {
					gt.logger.Println("Client state switched to healthy", gt.ActiveRPCCient.client.Host+gt.ActiveRPCCient.client.Port)
				}
				client.healthy = res
			}
			wg.Done()
		}(&wg, a)
	}
	wg.Wait()
	gt.clientLock.Unlock()
}

func (gt *GoTezos) getFirstHealthyClient() *TezClientWrapper {
	gt.clientLock.Lock()
	defer gt.clientLock.Unlock()
	for _, a := range gt.RPCClients {
		if a.healthy == true {
			return a
		}
	}
	return nil
}

func (gt *GoTezos) getRandomHealthyClient() *TezClientWrapper {
	gt.clientLock.Lock()
	defer gt.clientLock.Unlock()
	clients := []int{}
	for i := range gt.RPCClients {
		clients = append(clients, i)
	}
	for _, i := range gt.rand.Perm(len(clients)) {
		return gt.RPCClients[i]
	}
	return nil
}

func (gt *GoTezos) setActiveclient() error {
	if gt.balancerStrategy == "failover" {
		c := gt.getFirstHealthyClient()
		if c == nil {
			gt.checkHealthStatus()
			c = gt.getFirstHealthyClient()
			if c == nil {
				return NoClientError{}
			}
		}
		gt.ActiveRPCCient = c
	}

	if gt.balancerStrategy == "random" {
		c := gt.getRandomHealthyClient()
		if c == nil {
			gt.checkHealthStatus()
			c = gt.getRandomHealthyClient()
			if c == nil {
				return NoClientError{}
			}
			gt.ActiveRPCCient = c
		}
		gt.ActiveRPCCient = c

	}
	return nil
}

// GetResponse takes path endpoint and any arguments and returns the raw response of the query
func (gt *GoTezos) GetResponse(path string, args string) (ResponseRaw, error) {
	return gt.HandleResponse("GET", path, args)
}

// PostResponse takes path endpoint and any arguments and returns the raw response of the POST query
func (gt *GoTezos) PostResponse(path string, args string) (ResponseRaw, error) {
	return gt.HandleResponse("POST", path, args)
}

// HandleResponse takes the method (GET/POST ... etc), the query path, any arguments, and returns the raw response of the query
func (gt *GoTezos) HandleResponse(method string, path string, args string) (ResponseRaw, error) {
	e := gt.setActiveclient()
	if e != nil {
		gt.logger.Println("goTezos", "Could not find any healthy clients")
		return ResponseRaw{}, e
	}

	r, err := gt.ActiveRPCCient.client.GetResponse(method, path, args)
	if err != nil {
		gt.ActiveRPCCient.healthy = false
		gt.logger.Println(gt.ActiveRPCCient.client.Host+gt.ActiveRPCCient.client.Port, "Client state switched to unhealthy")

		// recurse call self, which will pick a new ActiveClient if defined
		return gt.HandleResponse(method, path, args)
	}

	// Received a HTTP 200 OK response, but payload could contain error message
	if strings.Contains(string(r.Bytes), "\"error\":") {

		rpcErrors := RPCGenericErrors{}
		rpcErrors, err := rpcErrors.UnmarshalJSON(r.Bytes)
		if err != nil {
			return r, err
		}

		// Just return the first error for now
		// TODO: Return all errors
		return r, fmt.Errorf("RPC Error (%s): %s", rpcErrors[0].Kind, rpcErrors[0].Error)
	}

	return r, nil
}

// NoClientError is a helper structure for error handling
type NoClientError struct {
}

func (gt NoClientError) Error() string {
	return "GoTezos did not find any healthy Tezos Node"
}
