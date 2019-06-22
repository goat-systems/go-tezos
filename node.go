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
func (n *NodeService) MonitorHeads(chain string, heads chan StructHeader, errc chan error, done chan bool) {
	query := "/monitor/heads/" + chain

	responses := make(chan []byte)
	errChan := make(chan error)

	go n.gt.client.StreamGet(query, nil, responses, errChan, done)

	defer close(errChan)
	defer close(responses)

	for {
		err := <-errChan
		res := <-responses
		if err != nil {
			errc <- err
			heads <- StructHeader{}
			return
		}

		blockHeader, err := UnmarshalBlockHeader(res)
		if err != nil {
			errc <- errors.Wrapf(err, "could not monitor chain heads '%s'", query)
			heads <- StructHeader{}
			return
		}

		errc <- nil
		heads <- blockHeader

		if blockHeader.Level == 0 {
			return
		}
	}
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
