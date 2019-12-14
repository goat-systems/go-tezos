package block

type TezosBlockService interface {
	GetHead() (Block, error)
	Get(id interface{}) (Block, error)
	ForgeBlockHeader(operation string) (string, error)
	IDToString(id interface{}) (string, error)
}
