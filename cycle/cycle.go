package cycle

import (
	"github.com/DefinitelyNotAGoat/go-tezos/block"
	"github.com/pkg/errors"
)

// CycleService is a struct wrapper for cycle functions
type CycleService struct {
	blockService block.TezosBlockService
}

// NewCycleService returns a new CycleService
func NewCycleService(blockService block.TezosBlockService) *CycleService {
	return &CycleService{blockService: blockService}
}

// GetCurrent gets the current cycle of the chain
func (s *CycleService) GetCurrent() (int, error) {
	block, err := s.blockService.GetHead()
	if err != nil {
		return 0, errors.Wrap(err, "could not get current cycle")
	}

	return block.Metadata.Level.Cycle, nil
}
