package snapshot

import (
	"encoding/json"
	"strconv"

	"github.com/pkg/errors"

	"github.com/DefinitelyNotAGoat/go-tezos/block"
	tzc "github.com/DefinitelyNotAGoat/go-tezos/client"
	"github.com/DefinitelyNotAGoat/go-tezos/cycle"
	"github.com/DefinitelyNotAGoat/go-tezos/network"
)

// SnapshotService is a struct wrapper for snap shot functions
type SnapshotService struct {
	cycleService cycle.TezosCycleService
	tzclient     tzc.TezosClient
	blockService block.TezosBlockService
	constants    network.Constants
}

// Snapshot is a SnapShot on the Tezos Network.
type Snapshot struct {
	Cycle           int
	Number          int
	AssociatedHash  string
	AssociatedBlock int
}

// SnapshotQuery is a Snapshot returned by the Tezos RPC API.
type SnapshotQuery struct {
	RandomSeed   string `json:"random_seed"`
	RollSnapShot int    `json:"roll_snapshot"`
}

// NewSnapshotService returns a new SnapShotService
func NewSnapshotService(cycleService cycle.TezosCycleService, tzclient tzc.TezosClient, blockService block.TezosBlockService, constants network.Constants) *SnapshotService {
	return &SnapshotService{
		cycleService: cycleService,
		tzclient:     tzclient,
		blockService: blockService,
		constants:    constants,
	}
}

// Get takes a cycle number and returns a helper structure describing a snap shot on the tezos network.
func (s *SnapshotService) Get(cycle int) (Snapshot, error) {

	var snapShotQuery SnapshotQuery
	var snap Snapshot

	currentCycle, err := s.cycleService.GetCurrent()
	if err != nil {
		return snap, errors.Wrapf(err, "could not get snapshot at cycle '%d'", cycle)
	}

	if cycle > currentCycle+s.constants.PreservedCycles-1 {
		return snap, errors.Errorf("could not get snapshot at cycle '%d', cycle requested is in the future", cycle)
	}

	snap.Cycle = cycle
	strCycle := strconv.Itoa(cycle)

	query := "/chains/main/blocks/"
	if cycle < currentCycle {
		block := strconv.Itoa(cycle*s.constants.BlocksPerCycle + 1)
		query = query + block + "/context/raw/json/cycle/" + strCycle
	} else {
		query = query + "head/context/raw/json/cycle/" + strCycle
	}

	resp, err := s.tzclient.Get(query, nil)
	if err != nil {
		return snap, errors.Wrapf(err, "could not get snapshot '%s'", query)
	}

	snapShotQuery, err = snapShotQuery.unmarshalJSON(resp)
	if err != nil {
		return snap, errors.Wrapf(err, "could not get snapshot '%s'", query)
	}

	snap.Number = snapShotQuery.RollSnapShot

	snap.AssociatedBlock = ((cycle - s.constants.PreservedCycles - 2) * s.constants.BlocksPerCycle) + (snap.Number+1)*s.constants.BlocksPerRollSnapshot
	if snap.AssociatedBlock < 1 {
		snap.AssociatedBlock = 1
	}

	block, err := s.blockService.Get(snap.AssociatedBlock)
	if err != nil {
		return snap, errors.Wrapf(err, "could not get snapshot '%s'", query)
	}
	snap.AssociatedHash = block.Hash

	return snap, nil
}

// unmarshalJSON unmarshals the bytes received as a parameter, into the type SnapShotQuery.
func (sq *SnapshotQuery) unmarshalJSON(v []byte) (SnapshotQuery, error) {
	snapShotQuery := SnapshotQuery{}
	err := json.Unmarshal(v, &snapShotQuery)
	if err != nil {
		return snapShotQuery, errors.Wrap(err, "could not unmarshal bytes into SnapShotQuery")
	}
	return snapShotQuery, nil
}
