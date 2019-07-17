package account

import (
	"github.com/DefinitelyNotAGoat/go-tezos/block"
	"github.com/DefinitelyNotAGoat/go-tezos/snapshot"
)

type snapshotServiceMock struct{}

func (s *snapshotServiceMock) Get(cycle int) (snapshot.Snapshot, error) {
	return snapshot.Snapshot{}, nil
}

func (s *snapshotServiceMock) GetAll() ([]snapshot.Snapshot, error) {
	return []snapshot.Snapshot{}, nil
}

type blockServiceMock struct {
}

func (b *blockServiceMock) GetHead() (block.Block, error) {
	return block.Block{
		Metadata: block.Metadata{
			Level: block.Level{
				Cycle: 9,
			},
		},
	}, nil
}

func (b *blockServiceMock) Get(id interface{}) (block.Block, error) {
	return block.Block{
		Hash: "BMXVTnGN7rwaCE34yuAuKzTHaPgyCUBxuVkM2Bbfo5jZvrrbZrY",
	}, nil
}

func (b *blockServiceMock) IDToString(id interface{}) (string, error) {
	return "", nil
}

type clientMock struct {
	ReturnBody []byte
}

func (c *clientMock) Post(path, args string) ([]byte, error) {
	return c.ReturnBody, nil
}

func (c *clientMock) Get(path string, params map[string]string) ([]byte, error) {
	return c.ReturnBody, nil
}
