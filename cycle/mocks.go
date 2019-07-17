package cycle

import (
	"github.com/DefinitelyNotAGoat/go-tezos/block"
)

type blockServiceMock struct {
}

func (b *blockServiceMock) GetHead() (block.Block, error) {
	return block.Block{
		Metadata: block.Metadata{
			Level: block.Level{
				Cycle: 127,
			},
		},
	}, nil
}

func (b *blockServiceMock) Get(id interface{}) (block.Block, error) {
	return block.Block{}, nil
}

func (b *blockServiceMock) IDToString(id interface{}) (string, error) {
	return "", nil
}
