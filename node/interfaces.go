package node

import "github.com/DefinitelyNotAGoat/go-tezos/block"

type TezosNodeService interface {
	MonitorHeads(chain string, heads chan block.Header, errc chan error, done chan bool)
	Bootstrapped() (Bootstrap, error)
	CommitHash() (string, error)
}
