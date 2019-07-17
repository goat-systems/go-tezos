package operations

import (
	"github.com/DefinitelyNotAGoat/go-tezos/account"
	"github.com/DefinitelyNotAGoat/go-tezos/delegate"
)

type TezosOperationsService interface {
	CreateBatchPayment(payments []delegate.Payment, wallet account.Wallet, paymentFee int, gaslimit int, batchSize int) ([]string, error)
	InjectOperation(op string) ([]byte, error)
	GetBlockOperationHashes(id interface{}) ([]string, error)
}
