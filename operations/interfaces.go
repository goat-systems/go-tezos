package operations

import (
	"github.com/DefinitelyNotAGoat/go-tezos/v2/account"
	"github.com/DefinitelyNotAGoat/go-tezos/v2/delegate"
)

type TezosOperationsService interface {
	CreateBatchPayment(payments []delegate.Payment, wallet account.Wallet, paymentFee int, gaslimit int, batchSize int) ([]string, error)
	InjectOperation(op string) ([]byte, error)
	GetBlockOperationHashes(id interface{}) ([]string, error)
	ForgeOperationBytes(contents string) (string, error)
	SignOperationBytes(operationBytes string, wallet account.Wallet) (string, error)
	SignEndorsementBytes(operationBytes, chainID string, wallet account.Wallet) (string, error)
	DecodeSignature(sig string) (string, string, error)
	PreApplyOperations(opstring string) error
}
