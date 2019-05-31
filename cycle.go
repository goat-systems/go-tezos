package gotezos

import "github.com/pkg/errors"

// CycleService is a struct wrapper for cycle functions
type CycleService struct {
	gt *GoTezos
}

// NewCycleService returns a new CycleService
func (gt *GoTezos) newCycleService() *CycleService {
	return &CycleService{gt: gt}
}

// GetCurrent gets the current cycle of the chain
func (s *CycleService) GetCurrent() (int, error) {
	block, err := s.gt.Block.GetHead()
	if err != nil {
		return 0, errors.Wrap(err, "could not get current cycle")
	}

	return block.Metadata.Level.Cycle, nil
}
