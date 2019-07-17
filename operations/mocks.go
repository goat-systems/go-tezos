package operations

// import (
// 	"github.com/DefinitelyNotAGoat/go-tezos/block"
// )

// type clientMock struct {
// 	ReturnBodysPost [][]byte
// 	ReturnBodysGet  [][]byte
// }

// func (c *clientMock) Post(path, args string) ([]byte, error) {
// 	var returnBody []byte
// 	returnBody, c.ReturnBodysPost = c.ReturnBodysPost[0], c.ReturnBodysPost[1:]
// 	return returnBody, nil
// }

// func (c *clientMock) Get(path string, params map[string]string) ([]byte, error) {
// 	var returnBody []byte
// 	returnBody, c.ReturnBodysGet = c.ReturnBodysGet[0], c.ReturnBodysGet[1:]
// 	return returnBody, nil
// }

// type blockServiceMock struct {
// }

// func (b *blockServiceMock) GetHead() (block.Block, error) {
// 	return block.Block{
// 		Metadata: block.Metadata{
// 			Level: block.Level{
// 				Cycle: 9,
// 			},
// 		},
// 	}, nil
// }

// func (b *blockServiceMock) Get(id interface{}) (block.Block, error) {
// 	return block.Block{
// 		Hash: "BMXVTnGN7rwaCE34yuAuKzTHaPgyCUBxuVkM2Bbfo5jZvrrbZrY",
// 	}, nil
// }

// func (b *blockServiceMock) IDToString(id interface{}) (string, error) {
// 	return "", nil
// }
