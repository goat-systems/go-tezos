package goTezos

import (
	"math/rand"
	"time"
	"sync"
	"log"
)

type TezClientWrapper struct {
	healthy bool // isHealthy
	client  *TezosRPCClient
}

/*

 * GoTezos manages multiple Clients
 * each Client represents a Connection to a Tezos Node
 * GoTezos manages failover if one Node is down, there
 * are 2 Strategies:
 * failover: always use the same unless it is down -> go to the next - default
 * random: send to each Node equally
 */
type GoTezos struct {
	clientLock sync.Mutex
	RpcClients []*TezClientWrapper
	ActiveRPCCient *TezClientWrapper
	balancerStrategy string
	rand *rand.Rand
	logger *log.Logger
}


func NewGoTezos() *GoTezos {
	a := GoTezos{}
	a.UseBalancerStrategyFailover()
	a.rand = rand.New(rand.NewSource(time.Now().Unix()))
	go func(a *GoTezos){
		for {
			time.Sleep(15*time.Second)
			a.checkUnhealthyClients()
		}
	}(&a)
	return &a
}



func (this *GoTezos) SetLogger(log *log.Logger) {
	this.logger = log
}

func (this *GoTezos) AddNewClient(client *TezosRPCClient) {
	this.clientLock.Lock()
	this.RpcClients = append(this.RpcClients, &TezClientWrapper{true, client})
	this.clientLock.Unlock()
}

func (this *GoTezos) UseBalancerStrategyFailover(){
	this.balancerStrategy = "failover"
}

func (this *GoTezos) UseBalancerStrategyRandom(){
	this.balancerStrategy = "random"
}

func (this *GoTezos) checkHealthStatus(){
	this.clientLock.Lock()
	wg := sync.WaitGroup{}
	for _,a := range this.RpcClients {
		wg.Add(1)
		go func(wg *sync.WaitGroup, client *TezClientWrapper){
			res := a.client.Healthcheck()
			if a.healthy && res == false {
				this.logger.Println("Client state switched to unhealthy", this.ActiveRPCCient.client.Host + this.ActiveRPCCient.client.Port)
			}
			if !a.healthy && res {
				this.logger.Println("Client state switched to healthy", this.ActiveRPCCient.client.Host + this.ActiveRPCCient.client.Port)
			}
			a.healthy = res
			wg.Done()
		}(&wg, a)
	}
	wg.Wait()
	this.clientLock.Unlock()
}

func (this *GoTezos) checkUnhealthyClients(){
	this.clientLock.Lock()
	wg := sync.WaitGroup{}
	for _,a := range this.RpcClients {
		wg.Add(1)
		go func(wg *sync.WaitGroup, client *TezClientWrapper){
			if a.healthy == false {
				res := a.client.Healthcheck()
				if !a.healthy && res {
					this.logger.Println("Client state switched to healthy", this.ActiveRPCCient.client.Host + this.ActiveRPCCient.client.Port)
				}
				a.healthy = res
			}
			wg.Done()
		}(&wg, a)
	}
	wg.Wait()
	this.clientLock.Unlock()
}

func (this *GoTezos) getFirstHealthyClient() *TezClientWrapper{
	this.clientLock.Lock()
	defer this.clientLock.Unlock()
	for _,a := range this.RpcClients {
		if a.healthy == true {
			return a
		}
	}
	return nil
}

func (this *GoTezos) getRandomHealthyClient() *TezClientWrapper{
	this.clientLock.Lock()
	defer this.clientLock.Unlock()
	clients := []int{}
	for i,_ := range this.RpcClients {
		clients = append(clients, i)
	}
	for _, i := range this.rand.Perm(len(clients)) {
		return this.RpcClients[i]
	}
	return nil
}


func (this *GoTezos) setActiveclient() error {
	if this.balancerStrategy == "failover" {
		c := this.getFirstHealthyClient()
		if c == nil {
			this.checkHealthStatus()
			c = this.getFirstHealthyClient()
			if c == nil {
				return NoClientError{}
			}
		}
		this.ActiveRPCCient = c
	}

	if this.balancerStrategy == "random" {
		c := this.getRandomHealthyClient()
		if c == nil {
			this.checkHealthStatus()
			c = this.getRandomHealthyClient()
			if c == nil {
				return NoClientError{}
			} else {
				this.ActiveRPCCient = c
			}
		}
		this.ActiveRPCCient = c

	}
	return nil
}


func (this *GoTezos) GetResponse(path string, args string) (ResponseRaw, error) {
	return this.HandleResponse("GET", path, args)
}


func (this *GoTezos) PostResponse(path string, args string) (ResponseRaw, error) {
	return this.HandleResponse("POST", path, args)
}


func (this *GoTezos) HandleResponse(method string, path string, args string) (ResponseRaw, error) {
	e := this.setActiveclient()
	if e != nil {
		this.logger.Println("goTezos", "Could not find any healthy clients")
		return ResponseRaw{}, e
	}

	r, err := this.ActiveRPCCient.client.GetResponse(method, path, args)
	if err != nil {
		this.ActiveRPCCient.healthy = false
		this.logger.Println(this.ActiveRPCCient.client.Host + this.ActiveRPCCient.client.Port, "Client state switched to unhealthy")
		return this.GetResponse(path, args)
	}
	return r, err
}
