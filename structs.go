package gotezos

import (
	"encoding/json"
	"log"
	"math/rand"
	"sync"
	"time"

	gocache "github.com/patrickmn/go-cache"
)

// MUTEZ is a helper for balance devision
const MUTEZ = 1000000

// ResponseRaw represents a raw RPC/HTTP response
type ResponseRaw struct {
	Bytes []byte
}

// NetworkVersion represents the network version returned by the Tezos network.
type NetworkVersion struct {
	Name    string `json:"name"`
	Major   int    `json:"major"`
	Minor   int    `json:"minor"`
	Network string // Human readable network name
}

// NetworkVersions is an array of NetworkVersion
type NetworkVersions []NetworkVersion

// NetworkConstants represents the network constants returned by the Tezos network.
type NetworkConstants struct {
	ProofOfWorkNonceSize         int      `json:"proof_of_work_nonce_size"`
	NonceLength                  int      `json:"nonce_length"`
	MaxRevelationsPerBlock       int      `json:"max_revelations_per_block"`
	MaxOperationDataLength       int      `json:"max_operation_data_length"`
	MaxProposalsPerDelegate      int      `json:"max_proposals_per_delegate"`
	PreservedCycles              int      `json:"preserved_cycles"`
	BlocksPerCycle               int      `json:"blocks_per_cycle"`
	BlocksPerCommitment          int      `json:"blocks_per_commitment"`
	BlocksPerRollSnapshot        int      `json:"blocks_per_roll_snapshot"`
	BlocksPerVotingPeriod        int      `json:"blocks_per_voting_period"`
	TimeBetweenBlocks            []string `json:"time_between_blocks"`
	EndorsersPerBlock            int      `json:"endorsers_per_block"`
	HardGasLimitPerOperation     string   `json:"hard_gas_limit_per_operation"`
	HardGasLimitPerBlock         string   `json:"hard_gas_limit_per_block"`
	ProofOfWorkThreshold         string   `json:"proof_of_work_threshold"`
	TokensPerRoll                string   `json:"tokens_per_roll"`
	MichelsonMaximumTypeSize     int      `json:"michelson_maximum_type_size"`
	SeedNonceRevelationTip       string   `json:"seed_nonce_revelation_tip"`
	OriginationSize              int      `json:"origination_size"`
	BlockSecurityDeposit         string   `json:"block_security_deposit"`
	EndorsementSecurityDeposit   string   `json:"endorsement_security_deposit"`
	BlockReward                  string   `json:"block_reward"`
	EndorsementReward            string   `json:"endorsement_reward"`
	CostPerByte                  string   `json:"cost_per_byte"`
	HardStorageLimitPerOperation string   `json:"hard_storage_limit_per_operation"`
}

// Block is a block returned by the Tezos RPC API.
type Block struct {
	Protocol   string               `json:"protocol"`
	ChainID    string               `json:"chain_id"`
	Hash       string               `json:"hash"`
	Header     StructHeader         `json:"header"`
	Metadata   StructMetadata       `json:"metadata"`
	Operations [][]StructOperations `json:"operations"`
}

// StructHeader is a header in a block returned by the Tezos RPC API.
type StructHeader struct {
	Level            int       `json:"level"`
	Proto            int       `json:"proto"`
	Predecessor      string    `json:"Predecessor"`
	Timestamp        time.Time `json:"timestamp"`
	ValidationPass   int       `json:"validation_pass"`
	OperationsHash   string    `json:"operations_hash"`
	Fitness          []string  `json:"fitness"`
	Context          string    `json:"context"`
	Priority         int       `json:"priority"`
	ProofOfWorkNonce string    `json:"proof_of_work_nonce"`
	Signature        string    `json:"signature"`
}

// StructMetadata is the Metadata in a block returned by the Tezos RPC API.
type StructMetadata struct {
	Protocol               string                         `json:"protocol"`
	NextProtocol           string                         `json:"next_protocol"`
	TestChainStatus        StructTestChainStatus          `json:"test_chain_status"`
	MaxOperationsTTL       int                            `json:"max_operations_ttl"`
	MaxOperationDataLength int                            `json:"max_operation_data_length"`
	MaxBlockHeaderLength   int                            `json:"max_block_header_length"`
	MaxOperationListLength []StructMaxOperationListLength `json:"max_operation_list_length"`
	Baker                  string                         `json:"baker"`
	Level                  StructLevel                    `json:"level"`
	VotingPeriodKind       string                         `json:"voting_period_kind"`
	NonceHash              interface{}                    `json:"nonce_hash"`
	ConsumedGas            string                         `json:"consumed_gas"`
	Deactivated            []string                       `json:"deactivated"`
	BalanceUpdates         []StructBalanceUpdates         `json:"balance_updates"`
}

// StructTestChainStatus is the TestChainStatus found in the Metadata of a block returned by the Tezos RPC API.
type StructTestChainStatus struct {
	Status string `json:"status"`
}

// StructMaxOperationListLength is the MaxOperationListLength found in the Metadata of a block returned by the Tezos RPC API.
type StructMaxOperationListLength struct {
	MaxSize int `json:"max_size"`
	MaxOp   int `json:"max_op,omitempty"`
}

// StructLevel the Level found in the Metadata of a block returned by the Tezos RPC API.
type StructLevel struct {
	Level                int  `json:"level"`
	LevelPosition        int  `json:"level_position"`
	Cycle                int  `json:"cycle"`
	CyclePosition        int  `json:"cycle_position"`
	VotingPeriod         int  `json:"voting_period"`
	VotingPeriodPosition int  `json:"voting_period_position"`
	ExpectedCommitment   bool `json:"expected_commitment"`
}

// StructBalanceUpdates is the BalanceUpdates found in the Metadata of a block returned by the Tezos RPC API.
type StructBalanceUpdates struct {
	Kind     string `json:"kind"`
	Contract string `json:"contract,omitempty"`
	Change   string `json:"change"`
	Category string `json:"category,omitempty"`
	Delegate string `json:"delegate,omitempty"`
	Level    int    `json:"level,omitempty"`
}

// StructOperations is the Operations found in a block returned by the Tezos RPC API.
type StructOperations struct {
	Protocol  string           `json:"protocol"`
	ChainID   string           `json:"chain_id"`
	Hash      string           `json:"hash"`
	Branch    string           `json:"branch"`
	Contents  []StructContents `json:"contents"`
	Signature string           `json:"signature"`
}

// StructContents is the Contents found in a operation of a block returned by the Tezos RPC API.
type StructContents struct {
	Kind             string           `json:"kind"`
	Source           string           `json:"source"`
	Fee              string           `json:"fee"`
	Counter          string           `json:"counter"`
	GasLimit         string           `json:"gas_limit"`
	StorageLimit     string           `json:"storage_limit"`
	Amount           string           `json:"amount"`
	Destination      string           `json:"destination"`
	Delegate         string           `json:"delegate"`
	Phk              string           `json:"phk"`
	Secret           string           `json:"secret"`
	Level            int              `json:"level"`
	ManagerPublicKey string           `json:"managerPubkey"`
	Balance          string           `json:"balance"`
	Metadata         ContentsMetadata `json:"metadata"`
}

// ContentsMetadata is the Metadata found in the Contents in a operation of a block returned by the Tezos RPC API.
type ContentsMetadata struct {
	BalanceUpdates []StructBalanceUpdates `json:"balance_updates"`
	Slots          []int                  `json:"slots"`
}

// SnapShot is a SnapShot on the Tezos Network.
type SnapShot struct {
	Cycle           int
	Number          int
	AssociatedHash  string
	AssociatedBlock int
}

// SnapShotQuery is a SnapShot returned by the Tezos RPC API.
type SnapShotQuery struct {
	RandomSeed   string `json:"random_seed"`
	RollSnapShot int    `json:"roll_snapshot"`
}

// FrozenBalanceRewards is a FrozenBalanceRewards query returned by the Tezos RPC API.
type FrozenBalanceRewards struct {
	Deposits string `json:"deposits"`
	Fees     string `json:"fees"`
	Rewards  string `json:"rewards"`
}

// TransOp is a helper structure to build out a transfer operation to post to the Tezos RPC
type TransOp struct {
	Kind         string `json:"kind"`
	Amount       string `json:"amount"`
	Source       string `json:"source"`
	Destination  string `json:"destination"`
	StorageLimit string `json:"storage_limit"`
	GasLimit     string `json:"gas_limit"`
	Fee          string `json:"fee"`
	Counter      string `json:"counter"`
}

// Conts is helper structure to build out the contents of a a transfer operation to post to the Tezos RPC
type Conts struct {
	Contents []TransOp `json:"contents"`
	Branch   string    `json:"branch"`
}

// Transfer a complete transfer request
type Transfer struct {
	Conts
	Protocol  string `json:"protocol"`
	Signature string `json:"signature"`
}

// Delegate is representation of a delegate on the Tezos Network
type Delegate struct {
	Balance              string                 `json:"balance"`
	FrozenBalance        string                 `json:"frozen_balance"`
	FrozenBalanceByCycle []FrozenBalanceByCycle `json:"frozen_balance_by_cycle"`
	StakingBalance       string                 `json:"staking_balance"`
	DelegateContracts    []string               `json:"delegated_contracts"`
	DelegatedBalance     string                 `json:"delegated_balance"`
	Deactivated          bool                   `json:"deactivated"`
	GracePeriod          int                    `json:"grace_period"`
}

// FrozenBalanceByCycle a representation of frozen balance by cycle on the Tezos network
type FrozenBalanceByCycle struct {
	Cycle   int    `json:"cycle"`
	Deposit string `json:"deposit"`
	Fees    string `json:"fees"`
	Rewards string `json:"rewards"`
}

// FrozenBalance is representation of frozen balance on the Tezos network
type FrozenBalance struct {
	Deposits string `json:"deposits"`
	Fees     string `json:"fees"`
	Rewards  string `json:"rewards"`
}

// BakingRights a representation of baking rights on the Tezos network
type BakingRights []struct {
	Level         int       `json:"level"`
	Delegate      string    `json:"delegate"`
	Priority      int       `json:"priority"`
	EstimatedTime time.Time `json:"estimated_time"`
}

// EndorsingRights is a representation of endorsing rights on the Tezos network
type EndorsingRights []struct {
	Level         int       `json:"level"`
	Delegate      string    `json:"delegate"`
	Slots         []int     `json:"slots"`
	EstimatedTime time.Time `json:"estimated_time"`
}

type DelegateReport struct {
	DelegatePhk      string
	Cycle            int
	Delegations      []DelegationReport
	CycleRewards     string
	TotalFeeRewards  string
	SelfBakedRewards string
	TotalRewards     string
}

type DelegationReport struct {
	DelegationPhk string
	Share         float64
	GrossRewards  string
	Fee           string
	NetRewards    string
}

// BRights is a structure representing baking rights for a specific delegate between cycles
type BRights struct {
	Delegate string    `json:"delegate"`
	Cycles   []BCycles `json:"cycles"`
}

// ERights is a structure representing endorsing rights for a specific delegate between cycles
type ERights struct {
	Delegate string    `json:"delegate"`
	Cycles   []ECycles `json:"cycles"`
}

// BCycles is a structure representing the baking rights in a specific cycle
type BCycles struct {
	Cycle        int          `json:"cycle"`
	BakingRights BakingRights `json:"baking_rights"`
}

// ECycles is a structure representing the endorsing rights in a specific cycle
type ECycles struct {
	Cycle           int             `json:"cycle"`
	EndorsingRights EndorsingRights `json:"endorsing_rights"`
}

// Payment is a helper struct for transfers
type Payment struct {
	Address string
	Amount  float64
}

// TezClientWrapper is a wrapper for the TezosRPCClient
type TezClientWrapper struct {
	healthy bool // isHealthy
	client  *TezosRPCClient
}

// RPCGenericError is an Error helper for the RPC
type RPCGenericError struct {
	Kind  string `json:"kind"`
	Error string `json:"error"`
}

// RPCGenericErrors and array of RPCGenericErrors
type RPCGenericErrors []RPCGenericError

// OperationHashes slice
type OperationHashes []string

// GoTezos manages multiple Clients
// each Client represents a Connection to a Tezos Node
// GoTezos manages failover if one Node is down, there
// are 2 Strategies:
// failover: always use the same unless it is down -> go to the next - default
// random: send to each Node equally
type GoTezos struct {
	clientLock       sync.Mutex
	RPCClients       []*TezClientWrapper
	ActiveRPCCient   *TezClientWrapper
	Constants        NetworkConstants
	Versions         []NetworkVersion
	balancerStrategy string
	rand             *rand.Rand
	logger           *log.Logger
	cache            *gocache.Cache
	debug            bool
}

//Wallet needed for signing operations
type Wallet struct {
	Address  string
	Mnemonic string
	Seed     []byte
	Kp       KeyPair
	Sk       string
	Pk       string
}

// Key Pair Storage
type KeyPair struct {
	PrivKey []byte
	PubKey  []byte
}

// UnmarshalJSON unmarshals the bytes received as a parameter, into the type NetworkVersion.
func (nvs *NetworkVersions) UnmarshalJSON(v []byte) (NetworkVersions, error) {
	networkVersions := NetworkVersions{}
	err := json.Unmarshal(v, &networkVersions)
	if err != nil {
		return networkVersions, err
	}
	return networkVersions, nil
}

// UnmarshalJSON unmarshals bytes received as a parameter, into the type NetworkConstants.
func (nc *NetworkConstants) UnmarshalJSON(v []byte) (NetworkConstants, error) {
	networkConstants := NetworkConstants{}
	err := json.Unmarshal(v, &networkConstants)
	if err != nil {
		log.Println("Could not get unMarshal bytes into NetworkConstants: " + err.Error())
		return networkConstants, err
	}
	return networkConstants, nil
}

// UnmarshalJSON unmarshals the bytes received as a parameter, into the type Block.
func (b *Block) UnmarshalJSON(v []byte) (Block, error) {
	block := Block{}
	err := json.Unmarshal(v, &block)
	if err != nil {
		log.Println("Could not get unMarshal bytes into block: " + err.Error())
		return block, err
	}
	return block, nil
}

// UnmarshalJSON unmarshals the bytes received as a parameter, into the type SnapShotQuery.
func (sq *SnapShotQuery) UnmarshalJSON(v []byte) (SnapShotQuery, error) {
	snapShotQuery := SnapShotQuery{}
	err := json.Unmarshal(v, &snapShotQuery)
	if err != nil {
		log.Println("Could not unmarhel SnapShotQuery: " + err.Error())
		return snapShotQuery, err
	}
	return snapShotQuery, nil
}

// UnmarshalJSON unmarshals the bytes received as a parameter, into the type SnapShotQuery.
func (fb *FrozenBalanceRewards) UnmarshalJSON(v []byte) (FrozenBalanceRewards, error) {
	frozenBalance := FrozenBalanceRewards{}
	err := json.Unmarshal(v, &frozenBalance)
	if err != nil {
		log.Println("Could not unmarhel frozenBalanceRewards: " + err.Error())
		return frozenBalance, err
	}
	return frozenBalance, nil
}

// unmarshalString unmarshals the bytes received as a parameter, into the type string.
func unmarshalString(v []byte) (string, error) {
	var str string
	err := json.Unmarshal(v, &str)
	if err != nil {
		log.Println("Could not unMarshal to string " + err.Error())
		return str, err
	}
	return str, nil
}

// UnmarshalJSON unmarshals the bytes received as a parameter, into the type an array of strings.
func unMarshalStringArray(v []byte) ([]string, error) {
	var strs []string

	err := json.Unmarshal(v, &strs)
	if err != nil {
		log.Println("Could not unMarshal to strings " + err.Error())
		return strs, err
	}
	return strs, nil
}

// UnmarshalJSON unmarshalls bytes into StructDelegate
func (d *Delegate) UnmarshalJSON(v []byte) (Delegate, error) {
	delegate := Delegate{}
	err := json.Unmarshal(v, &delegate)
	if err != nil {
		return delegate, err
	}
	return delegate, nil
}

// UnmarshalJSON unmarshalls bytes into frozen balance
func (fb *FrozenBalance) UnmarshalJSON(v []byte) (FrozenBalance, error) {
	frozenBalance := FrozenBalance{}
	err := json.Unmarshal(v, &frozenBalance)
	if err != nil {
		return frozenBalance, err
	}
	return frozenBalance, nil
}

// UnmarshalJSON unmarhsels bytes into Baking_Rights
func (br *BakingRights) UnmarshalJSON(v []byte) (BakingRights, error) {
	bakingRights := BakingRights{}
	err := json.Unmarshal(v, &bakingRights)
	if err != nil {
		return bakingRights, err
	}
	return bakingRights, nil
}

// UnmarshalJSON unmarhsels bytes into Endorsing_Rights
func (er *EndorsingRights) UnmarshalJSON(v []byte) (EndorsingRights, error) {
	endorsingRights := EndorsingRights{}
	err := json.Unmarshal(v, &endorsingRights)
	if err != nil {
		return endorsingRights, err
	}
	return endorsingRights, nil
}

func (c Conts) String() string {
	res, _ := json.Marshal(c)
	return string(res)
}

// UnmarshalJSON unmarhsels bytes into OperationHashes
func (oh *OperationHashes) UnmarshalJSON(v []byte) (OperationHashes, error) {

	// RPC returns slice of slice
	// Will flatten to single slice of ops for easy use
	ops := [][]string{}
	operationHashes := OperationHashes{}

	err := json.Unmarshal(v, &ops)
	if err != nil {
		return operationHashes, err
	}

	// flatten
	for _, i := range ops {
		for _, j := range i {
			operationHashes = append(*oh, j)
		}
	}
	return operationHashes, nil
}

// UnmarshalJSON unmarhsels bytes into RPCGenericErrors
func (ge *RPCGenericErrors) UnmarshalJSON(v []byte) (RPCGenericErrors, error) {
	r := RPCGenericErrors{}

	err := json.Unmarshal(v, &r)
	if err != nil {
		return r, err
	}
	return r, nil
}
