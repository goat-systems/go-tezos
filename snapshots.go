package gotezos

import (
	"encoding/json"
	"strconv"

	"github.com/pkg/errors"
)

// SnapShotService is a struct wrapper for snap shot functions
type SnapShotService struct {
	gt *GoTezos
}

// SnapShot is a SnapShot on the Tezos Network.
type SnapShot struct {
	Cycle                int
	Number               int
	AssociatedBlockHash  string
	AssociatedBlockLevel int
}

// SnapShotQuery is a SnapShot returned by the Tezos RPC API.
type SnapShotQuery struct {
	RandomSeed   string `json:"random_seed"`
	RollSnapShot int    `json:"roll_snapshot"`
}

// NewSnapShotService returns a new SnapShotService
func (gt *GoTezos) newSnapShotService() *SnapShotService {
	return &SnapShotService{gt: gt}
}

// Get takes a cycle number and returns a helper structure describing a snap shot on the tezos network.
func (s *SnapShotService) Get(cycle int) (SnapShot, error) {

	var snapShotQuery SnapShotQuery
	var snap SnapShot

	currentCycle, err := s.gt.Cycle.GetCurrent()
	if err != nil {
		return snap, errors.Wrapf(err, "could not get snapshot at cycle '%d'", cycle)
	}

	if cycle > currentCycle+s.gt.Constants.PreservedCycles-1 {
		return snap, errors.Errorf("could not get snapshot at cycle '%d', cycle requested is in the future", cycle)
	}

	snap.Cycle = cycle
	strCycle := strconv.Itoa(cycle)

	query := "/chains/main/blocks/"
	if cycle < currentCycle {
		block := strconv.Itoa(cycle*s.gt.Constants.BlocksPerCycle + 1)
		query = query + block + "/context/raw/json/cycle/" + strCycle
	} else {
		query = query + "head/context/raw/json/cycle/" + strCycle
	}

	resp, err := s.gt.Get(query, nil)
	if err != nil {
		return snap, errors.Wrapf(err, "could not get snapshot '%s'", query)
	}

	snapShotQuery, err = snapShotQuery.unmarshalJSON(resp)
	if err != nil {
		return snap, errors.Wrapf(err, "could not get snapshot '%s'", query)
	}

	snap.Number = snapShotQuery.RollSnapShot

	snap.AssociatedBlockLevel = ((cycle - s.gt.Constants.PreservedCycles - 2) * s.gt.Constants.BlocksPerCycle) + (snap.Number+1)*s.gt.Constants.BlocksPerRollSnapshot
	if snap.AssociatedBlockLevel < 1 {
		snap.AssociatedBlockLevel = 1
	}

	block, err := s.gt.Block.Get(snap.AssociatedBlockLevel)
	if err != nil {
		return snap, errors.Wrapf(err, "could not get snapshot '%s'", query)
	}
	snap.AssociatedBlockHash = block.Hash

	return snap, nil
}

// GetAll gets a list of all known snapshots to the network
func (s *SnapShotService) GetAll() ([]SnapShot, error) {
	var snapShotArray []SnapShot
	currentCycle, err := s.gt.Cycle.GetCurrent()
	if err != nil {
		return snapShotArray, errors.Wrap(err, "could not get all snapshots")
	}
	for i := 7; i <= currentCycle; i++ {
		snapShot, err := s.Get(i)
		if err != nil {
			return snapShotArray, errors.Wrap(err, "could not get all snapshots")
		}
		snapShotArray = append(snapShotArray, snapShot)
	}

	return snapShotArray, nil
}

// unmarshalJSON unmarshals the bytes received as a parameter, into the type SnapShotQuery.
func (sq *SnapShotQuery) unmarshalJSON(v []byte) (SnapShotQuery, error) {
	snapShotQuery := SnapShotQuery{}
	err := json.Unmarshal(v, &snapShotQuery)
	if err != nil {
		return snapShotQuery, errors.Wrap(err, "could not unmarshal bytes into SnapShotQuery")
	}
	return snapShotQuery, nil
}
