package gotezos

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

// MUTEZ is a helper for balance devision
const MUTEZ = 1000000

// GoTezos is the driver of the library, it inludes the several RPC services
// like Block, SnapSHot, Cycle, Account, Delegate, Operations, Contract, and Network
type GoTezos struct {
	client    *client
	Constants NetworkConstants
	Block     *BlockService
	SnapShot  *SnapShotService
	Cycle     *CycleService
	Account   *AccountService
	Delegate  *DelegateService
	Network   *NetworkService
	Operation *OperationService
	Contract  *ContractService
	Node      *NodeService
}

// ResponseRaw represents a raw RPC/HTTP response
type responseRaw struct {
	Bytes []byte
}

// RPCGenericError is an Error helper for the RPC
type genericRPCError struct {
	Kind  string `json:"kind"`
	Error string `json:"error"`
}

// RPCGenericErrors and array of RPCGenericErrors
type genericRPCErrors []genericRPCError

// NewGoTezos is a constructor that returns a GoTezos object
func NewGoTezos(URL ...string) (*GoTezos, error) {
	gt := GoTezos{}

	var url string
	if len(URL) > 1 {
		// RPC Address
		url = URL[0]
	} else {
		err := godotenv.Load()
		if err != nil {
			return &gt, errors.Wrap(err, "Error loading .env file")
		}

		url = os.Getenv("RPC_ADDRESS")
	}

	gt.Block = gt.newBlockService()
	gt.SnapShot = gt.newSnapShotService()
	gt.Cycle = gt.newCycleService()
	gt.Account = gt.newAccountService()
	gt.Delegate = gt.newDelegateService()
	gt.Network = gt.newNetworkService()
	gt.Operation = gt.newOperationService()
	gt.Contract = gt.newContractService()
	gt.Node = gt.newNodeService()

	gt.client = newClient(url)

	var err error
	gt.Constants, err = gt.Network.GetConstants()
	if err != nil {
		return &gt, errors.Wrap(err, "could not get network constants")
	}

	return &gt, nil
}

// SetHTTPClient allows you to pass your own Go http client, with your own settings
func (gt *GoTezos) SetHTTPClient(client *http.Client) {
	gt.client.netClient = client
}

// Get takes path endpoint and returns the response of the query
func (gt *GoTezos) Get(path string, params map[string]string) ([]byte, error) {
	resp, err := gt.client.Get(path, params)
	if err != nil {
		return nil, err
	}

	err = gt.handleRPCError(resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Post takes path endpoint and any arguments and returns the response of the POST
func (gt *GoTezos) Post(path string, args string) ([]byte, error) {
	resp, err := gt.client.Post(path, args)
	if err != nil {
		return nil, err
	}

	err = gt.handleRPCError(resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (gt *GoTezos) handleRPCError(resp []byte) error {
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
