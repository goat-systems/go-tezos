package gt

import (
	"github.com/DefinitelyNotAGoat/go-tezos/v2/account"
	"github.com/DefinitelyNotAGoat/go-tezos/v2/block"
	tzc "github.com/DefinitelyNotAGoat/go-tezos/v2/client"
	"github.com/DefinitelyNotAGoat/go-tezos/v2/contracts"
	"github.com/DefinitelyNotAGoat/go-tezos/v2/cycle"
	"github.com/DefinitelyNotAGoat/go-tezos/v2/delegate"
	"github.com/DefinitelyNotAGoat/go-tezos/v2/network"
	"github.com/DefinitelyNotAGoat/go-tezos/v2/node"
	"github.com/DefinitelyNotAGoat/go-tezos/v2/operations"
	"github.com/DefinitelyNotAGoat/go-tezos/v2/snapshot"
	"github.com/pkg/errors"
)

// GoTezos is the driver of the library, it inludes the several RPC services
// like Block, SnapSHot, Cycle, Account, Delegate, Operations, Contract, and Network
type GoTezos struct {
	Client    tzc.TezosClient
	Constants network.Constants
	Block     block.TezosBlockService
	Snapshot  snapshot.TezosSnapshotService
	Cycle     cycle.TezosCycleService
	Account   account.TezosAccountService
	Delegate  delegate.TezosDelegateService
	Network   network.TezosNetworkService
	Operation operations.TezosOperationsService
	Contract  contracts.TezosContractsService
	Node      node.TezosNodeService
}

// NewGoTezos is a constructor that returns a GoTezos object
func NewGoTezos(URL string) (*GoTezos, error) {
	gotezos := GoTezos{}

	gotezos.Client = tzc.NewClient(URL)
	gotezos.Network = network.NewNetworkService(gotezos.Client)
	var err error
	gotezos.Constants, err = gotezos.Network.GetConstants()
	if err != nil {
		return &gotezos, errors.Wrap(err, "could not get network constants")
	}
	gotezos.Block = block.NewBlockService(gotezos.Client)
	gotezos.Cycle = cycle.NewCycleService(gotezos.Block)
	gotezos.Snapshot = snapshot.NewSnapshotService(
		gotezos.Cycle,
		gotezos.Client,
		gotezos.Block,
		gotezos.Constants,
	)
	gotezos.Account = account.NewAccountService(
		gotezos.Client,
		gotezos.Block,
		gotezos.Snapshot,
	)
	gotezos.Delegate = delegate.NewDelegateService(
		gotezos.Client,
		gotezos.Block,
		gotezos.Snapshot,
		gotezos.Account,
		gotezos.Constants,
	)
	gotezos.Operation = operations.NewOperationService(gotezos.Block, gotezos.Client)
	gotezos.Contract = contracts.NewContractService(gotezos.Client)
	gotezos.Node = node.NewNodeService(gotezos.Client)

	return &gotezos, nil
}
