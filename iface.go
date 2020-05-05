package gotezos

import "math/big"

// IFace is an interface mocking a GoTezos object.
type IFace interface {
	ActiveChains() (ActiveChains, error)
	BakingRights(input BakingRightsInput) (*BakingRights, error)
	Balance(blockhash, address string) (*big.Int, error)
	Block(id interface{}) (*Block, error)
	Blocks(input BlocksInput) ([][]string, error)
	Bootstrap() (Bootstrap, error)
	ChainID() (string, error)
	Checkpoint() (Checkpoint, error)
	Commit() (string, error)
	Connections() (Connections, error)
	Constants(blockhash string) (Constants, error)
	ContractStorage(blockhash string, KT1 string) ([]byte, error)
	Counter(blockhash, pkh string) (int, error)
	Cycle(cycle int) (Cycle, error)
	Delegate(blockhash, delegate string) (Delegate, error)
	Delegates(input DelegatesInput) ([]*string, error)
	DelegatedContracts(blockhash, delegate string) ([]*string, error)
	DelegatedContractsAtCycle(cycle int, delegate string) ([]*string, error)
	DeleteInvalidBlock(blockHash string) error
	EndorsingRights(input EndorsingRightsInput) (*EndorsingRights, error)
	FrozenBalance(cycle int, delegate string) (FrozenBalance, error)
	Head() (*Block, error)
	InjectionBlock(input InjectionBlockInput) ([]byte, error)
	InjectionOperation(input InjectionOperationInput) (string, error)
	InvalidBlock(blockHash string) (InvalidBlock, error)
	InvalidBlocks() ([]InvalidBlock, error)
	OperationHashes(blockhash string) ([][]string, error)
	PreapplyOperations(input PreapplyOperationsInput) ([]Operations, error)
	StakingBalance(blockhash, delegate string) (*big.Int, error)
	StakingBalanceAtCycle(cycle int, delegate string) (*big.Int, error)
	UserActivatedProtocolOverrides() (UserActivatedProtocolOverrides, error)
	Version() (Version, error)
}
