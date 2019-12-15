package gotezos

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

// Cycle is a Snapshot returned by the Tezos RPC API.
type Cycle struct {
	RandomSeed   string `json:"random_seed"`
	RollSnapshot int    `json:"roll_snapshot"`
	BlockHash    string `json:"-"`
}

// Cycle returns a cycle information
func (t *GoTezos) Cycle(cycle int) (Cycle, error) {
	head, err := t.Head()
	if err != nil {
		return Cycle{}, errors.Wrapf(err, "could not get cycle '%d'", cycle)
	}

	if cycle > head.Metadata.Level.Cycle+t.networkConstants.PreservedCycles-1 {
		return Cycle{}, errors.Errorf("could not get cycle '%d': request is in the future", cycle)
	}

	var c Cycle
	if cycle < head.Metadata.Level.Cycle {
		block, err := t.Block(cycle*t.networkConstants.BlocksPerCycle + 1)
		if err != nil {
			return Cycle{}, errors.Wrapf(err, "could not get cycle '%d'", cycle)
		}
		c, err = t.getCycleAtHash(block.Hash, cycle)
		if err != nil {
			return Cycle{}, errors.Wrapf(err, "could not get cycle '%d'", cycle)
		}

	} else {
		var err error
		c, err = t.getCycleAtHash(head.Hash, cycle)
		if err != nil {
			return Cycle{}, errors.Wrapf(err, "could not get cycle '%d'", cycle)
		}
	}

	level := ((cycle - t.networkConstants.PreservedCycles - 2) * t.networkConstants.BlocksPerCycle) + (c.RollSnapshot+1)*t.networkConstants.BlocksPerRollSnapshot
	if level < 1 {
		level = 1
	}

	block, err := t.Block(level)
	if err != nil {
		return c, errors.Wrapf(err, "could not get cycle '%d'", cycle)
	}

	c.BlockHash = block.Hash
	return c, nil
}

// getCycleAtHash returns a cycle information
func (t *GoTezos) getCycleAtHash(blockhash string, cycle int) (Cycle, error) {
	resp, err := t.get(fmt.Sprintf("/chains/main/blocks/%s/context/raw/json/cycle/%d", blockhash, cycle))
	if err != nil {
		return Cycle{}, errors.Wrapf(err, "could not get cycle at hash '%s'", blockhash)
	}

	var c Cycle
	err = json.Unmarshal(resp, &c)
	if err != nil {
		return c, errors.Wrapf(err, "could not unmarshal at cycle hash '%s'", blockhash)
	}

	return c, nil
}
