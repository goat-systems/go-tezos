package gotezos

// IFace is an interface mocking a GoTezos object.
type IFace interface {
	ActiveChains() (ActiveChains, error)
	BakingRights(input BakingRightsInput) (*BakingRights, error)
	Balance(input BalanceInput) (int, error)
	BallotList(blockhash string) (BallotList, error)
	Ballots(blockhash string) (Ballots, error)
	BigMap(input BigMapInput) ([]byte, error)
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
	CurrentPeriodKind(blockhash string) (string, error)
	CurrentProposal(blockhash string) (string, error)
	CurrentQuorum(blockhash string) (int, error)
	Cycle(cycle int) (Cycle, error)
	Delegate(blockhash, delegate string) (Delegate, error)
	Delegates(input DelegatesInput) ([]string, error)
	DelegatedContracts(input DelegatedContractsInput) ([]string, error)
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
	Proposals(blockhash string) (Proposals, error)
	StakingBalance(input StakingBalanceInput) (int, error)
	UserActivatedProtocolOverrides() (UserActivatedProtocolOverrides, error)
	Version() (Version, error)
}
