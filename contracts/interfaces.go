package contracts

type TezosContractsService interface {
	GetStorage(contract string) ([]byte, error)
}
