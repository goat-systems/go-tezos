package gotezos

import (
	"encoding/json"
	"github.com/pkg/errors"
	"time"
)

// NodeService is a service for node related functions
type NodeService struct {
	gt *GoTezos
}

// Bootstrap is a structure representing the bootstrapped response
type Bootstrap struct {
	Block     string
	Timestamp time.Time
}

// newNodeService returns a new NodeService
func (gt *GoTezos) newNodeService() *NodeService {
	return &NodeService{gt: gt}
}

// MonitorHeads gets the new Head of the chain every time it changes
func (n *NodeService) MonitorHeads(chain string, head chan StructHeader, done chan bool) error {
	query := "/monitor/heads/" + chain

	res := make(chan []byte)
	errorChan := make(chan error)

	go func() {
		err := n.gt.client.StreamGet(query, nil, res, done)
		if err != nil {
			errorChan <- err
		}
		close(res)
	}()

	for response := range res {
		select {
		case <-errorChan:
			return errors.Errorf("could not monitor chain heads '%s'", query)
		default:
			blockHeader, err := UnmarshalBlockHeader(response)
			if err != nil {
				return errors.Wrapf(err, "could not monitor chain heads '%s'", query)
			}

			head <- blockHeader
		}
	}

	return nil
}

// Bootstrapped gets the current node bootstrap
func (n *NodeService) Bootstrapped() (Bootstrap, error) {
	var b Bootstrap
	query := "/monitor/bootstrapped"
	resp, err := n.gt.Get(query, nil)
	if err != nil {
		return b, errors.Wrapf(err, "could not node bootstraped '%s'", query)
	}

	b, err = unmarshallBootstrap(resp)
	if err != nil {
		return b, errors.Wrapf(err, "could not node bootstraped '%s'", query)
	}

	return b, nil
}

// CommitHash gets the current commit the node is running
func (n *NodeService) CommitHash() (string, error) {
	var c string
	query := "/monitor/commit_hash"
	resp, err := n.gt.Get(query, nil)
	if err != nil {
		return c, errors.Wrapf(err, "could not node commit hash '%s'", query)
	}

	c, err = unmarshalString(resp)
	if err != nil {
		return c, errors.Wrapf(err, "could not node commit hash '%s'", query)
	}

	return c, nil
}

func unmarshallBootstrap(v []byte) (Bootstrap, error) {
	b := Bootstrap{}
	err := json.Unmarshal(v, &b)
	if err != nil {
		return b, errors.Wrap(err, "could not unmarshal bytes into Bootstrap")
	}

	return b, nil
}
