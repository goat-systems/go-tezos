package gotezos

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
)

// Bootstrap is a structure representing the bootstrapped response
type Bootstrap struct {
	Block     string
	Timestamp time.Time
}

// Bootstrap gets the current node bootstrap
func (t *GoTezos) Bootstrap() (Bootstrap, error) {
	resp, err := t.get("/monitor/bootstrapped")
	if err != nil {
		return Bootstrap{}, errors.Wrap(err, "could not get Bootstrap")
	}

	var bootstrap Bootstrap
	err = json.Unmarshal(resp, &bootstrap)
	if err != nil {
		return bootstrap, errors.Wrap(err, "could not unmarshal Bootsrap")
	}

	return bootstrap, nil
}

// Commit gets the current commit the node is running
func (t *GoTezos) Commit() (string, error) {
	resp, err := t.get("/monitor/commit_hash")
	if err != nil {
		return "", errors.Wrap(err, "could not get commit hash")
	}

	var commit string
	err = json.Unmarshal(resp, &commit)
	if err != nil {
		return commit, errors.Wrap(err, "could unmarshal commit hash")
	}

	return commit, nil
}
