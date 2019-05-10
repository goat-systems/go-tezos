package gotezos

import (
	"encoding/json"
	"fmt"
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

// Bootstrapped gets the current node bootstrap
func (n *NodeService) Bootstrapped() (Bootstrap, error) {
	var b Bootstrap
	query := "/monitor/bootstrapped"
	resp, err := n.gt.Get(query, nil)
	if err != nil {
		return b, fmt.Errorf("could not get bootstrap: %v", err)
	}

	b, err = unmarshallBootstrap(resp)
	if err != nil {
		return b, fmt.Errorf("could not get bootstrap: %v", err)
	}

	return b, nil
}

// CommitHash gets the current commit the node is running
func (n *NodeService) CommitHash() (string, error) {
	var c string
	query := "/monitor/commit_hash"
	resp, err := n.gt.Get(query, nil)
	if err != nil {
		return c, fmt.Errorf("could not get commit hash: %v", err)
	}

	c, err = unmarshalString(resp)
	if err != nil {
		return c, fmt.Errorf("could not get bootstrap: %v", err)
	}

	return c, nil
}

func unmarshallBootstrap(v []byte) (Bootstrap, error) {
	b := Bootstrap{}
	err := json.Unmarshal(v, &b)
	if err != nil {
		return b, err
	}

	return b, nil
}
