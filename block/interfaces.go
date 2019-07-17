package block

type TezosBlockService interface {
	GetHead() (Block, error)
	Get(id interface{}) (Block, error)
	IDToString(id interface{}) (string, error)
}
