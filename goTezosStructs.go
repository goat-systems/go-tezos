package goTezos

import (
	"encoding/json"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/jamesruan/sodium"
	gocache "github.com/patrickmn/go-cache"
)

const MUTEZ = 1000000

type ResponseRaw struct {
	Bytes []byte
}

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

//unMarshals the bytes received as a parameter, into the type NetworkConstants.
func unMarshalNetworkConstants(v []byte) (NetworkConstants, error) {
	var networkConstants NetworkConstants

	err := json.Unmarshal(v, &networkConstants)
	if err != nil {
		log.Println("Could not get unMarshal bytes into NetworkConstants: " + err.Error())
		return networkConstants, err
	}
	return networkConstants, nil
}

//An unMarshaled representation of a block returned by the Tezos RPC API.
type Block struct {
	Protocol   string               `json:"protocol"`
	ChainID    string               `json:"chain_id"`
	Hash       string               `json:"hash"`
	Header     StructHeader         `json:"header"`
	Metadata   StructMetadata       `json:"metadata"`
	Operations [][]StructOperations `json:"operations"`
}

//An unMarshaled representation of a header in a block returned by the Tezos RPC API.
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

//An unMarshaled representation of Metadata in a block returned by the Tezos RPC API.
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

//An unMarshaled representation of a TestChainStatus found in the Metadata of a block returned by the Tezos RPC API.
type StructTestChainStatus struct {
	Status string `json:"status"`
}

//An unMarshaled representation of a MaxOperationListLength found in the Metadata of a block returned by the Tezos RPC API.
type StructMaxOperationListLength struct {
	MaxSize int `json:"max_size"`
	MaxOp   int `json:"max_op,omitempty"`
}

//An unMarshaled representation of a Level found in the Metadata of a block returned by the Tezos RPC API.
type StructLevel struct {
	Level                int  `json:"level"`
	LevelPosition        int  `json:"level_position"`
	Cycle                int  `json:"cycle"`
	CyclePosition        int  `json:"cycle_position"`
	VotingPeriod         int  `json:"voting_period"`
	VotingPeriodPosition int  `json:"voting_period_position"`
	ExpectedCommitment   bool `json:"expected_commitment"`
}

//An unMarshaled representation of BalanceUpdates found in the Metadata of a block returned by the Tezos RPC API.
type StructBalanceUpdates struct {
	Kind     string `json:"kind"`
	Contract string `json:"contract,omitempty"`
	Change   string `json:"change"`
	Category string `json:"category,omitempty"`
	Delegate string `json:"delegate,omitempty"`
	Level    int    `json:"level,omitempty"`
}

//An unMarshaled representation of Operations found in a block returned by the Tezos RPC API.
type StructOperations struct {
	Protocol  string           `json:"protocol"`
	ChainID   string           `json:"chain_id"`
	Hash      string           `json:"hash"`
	Branch    string           `json:"branch"`
	Contents  []StructContents `json:"contents"`
	Signature string           `json:"signature"`
}

//An unMarshaled representation of Contents found in a operation of a block returned by the Tezos RPC API.
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

//An unMarshaled representation of Metadata found in the Contents in a operation of a block returned by the Tezos RPC API.
type ContentsMetadata struct {
	BalanceUpdates []StructBalanceUpdates `json:"balance_updates"`
	Slots          []int                  `json:"slots"`
}

//unMarshals the bytes received as a parameter, into the type Block.
func unMarshalBlock(v []byte) (Block, error) {
	var block Block

	err := json.Unmarshal(v, &block)
	if err != nil {
		log.Println("Could not get unMarshal bytes into block: " + err.Error())
		return block, err
	}
	return block, nil
}

//An easy representation of a SnapShot on the Tezos Network.
type SnapShot struct {
	Cycle           int
	Number          int
	AssociatedHash  string
	AssociatedBlock int
}

//An unMarshaled representation of a SnapShot returned by the Tezos RPC API.
type SnapShotQuery struct {
	RandomSeed   string `json:"random_seed"`
	RollSnapShot int    `json:"roll_snapshot"`
}

//unMarshals the bytes received as a parameter, into the type SnapShotQuery.
func unMarshalSnapShotQuery(v []byte) (SnapShotQuery, error) {
	var snapShotQuery SnapShotQuery

	err := json.Unmarshal(v, &snapShotQuery)
	if err != nil {
		log.Println("Could not unmarhel SnapShotQuery: " + err.Error())
		return snapShotQuery, err
	}
	return snapShotQuery, nil
}

//An unMarshaled representation of a FrozenBalanceRewards query returned by the Tezos RPC API.
type FrozenBalanceRewards struct {
	Deposits string `json:"deposits"`
	Fees     string `json:"fees"`
	Rewards  string `json:"rewards"`
}

//unMarshals the bytes received as a parameter, into the type SnapShotQuery.
func unMarshalFrozenBalanceRewards(v []byte) (FrozenBalanceRewards, error) {
	var frozenBalanceRewards FrozenBalanceRewards

	err := json.Unmarshal(v, &frozenBalanceRewards)
	if err != nil {
		log.Println("Could not unmarhel frozenBalanceRewards: " + err.Error())
		return frozenBalanceRewards, err
	}
	return frozenBalanceRewards, nil
}

//unMarshals the bytes received as a parameter, into the type string.
func unMarshalString(v []byte) (string, error) {
	var str string

	err := json.Unmarshal(v, &str)
	if err != nil {
		log.Println("Could not unMarshal to string " + err.Error())
		return str, err
	}
	return str, nil
}

//unMarshals the bytes received as a parameter, into the type an array of strings.
func unMarshalStringArray(v []byte) ([]string, error) {
	var strs []string

	err := json.Unmarshal(v, &strs)
	if err != nil {
		log.Println("Could not unMarshal to strings " + err.Error())
		return strs, err
	}
	return strs, nil
}

//A helper structure to build out a transfer operation to post to the Tezos RPC
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

//A helper structure to build out the contents of a a transfer operation to post to the Tezos RPC
type Conts struct {
	Contents []TransOp `json:"contents"`
	Branch   string    `json:"branch"`
}

func (c Conts) String() string {
	res, _ := json.Marshal(c)
	return string(res)
}

// A complete transfer request
type Transfer struct {
	Conts
	Protocol  string `json:"protocol"`
	Signature string `json:"signature"`
}

//An unmarshalled representation of a delegate
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

//An unmarshalled representation of frozen balance by cycle
type FrozenBalanceByCycle struct {
	Cycle   int    `json:"cycle"`
	Deposit string `json:"deposit"`
	Fees    string `json:"fees"`
	Rewards string `json:"rewards"`
}

//Unmarshalls bytes into StructDelegate
func unMarshalDelegate(v []byte) (Delegate, error) {
	var delegate Delegate

	err := json.Unmarshal(v, &delegate)
	if err != nil {
		return delegate, err
	}
	return delegate, nil
}

//An unmarshalled representation of frozen balance
type FrozenBalance struct {
	Deposits string `json:"deposits"`
	Fees     string `json:"fees"`
	Rewards  string `json:"rewards"`
}

//Unmarshalls bytes into frozen balance
func unMarshalFrozenBalance(v []byte) (FrozenBalance, error) {
	var frozenBalance FrozenBalance

	err := json.Unmarshal(v, &frozenBalance)
	if err != nil {
		return frozenBalance, err
	}
	return frozenBalance, nil
}

//A representation of baking rights on the network
type Baking_Rights []struct {
	Level         int       `json:"level"`
	Delegate      string    `json:"delegate"`
	Priority      int       `json:"priority"`
	EstimatedTime time.Time `json:"estimated_time"`
}

//Unmarhsels bytes into Baking_Rights
func unMarshalBakingRights(v []byte) (Baking_Rights, error) {
	var bakingRights Baking_Rights

	err := json.Unmarshal(v, &bakingRights)
	if err != nil {
		return bakingRights, err
	}
	return bakingRights, nil
}

//A representation of endorsing rights on the network
type Endorsing_Rights []struct {
	Level         int       `json:"level"`
	Delegate      string    `json:"delegate"`
	Slots         []int     `json:"slots"`
	EstimatedTime time.Time `json:"estimated_time"`
}

//Unmarhsels bytes into Endorsing_Rights
func unMarshalEndorsingRights(v []byte) (Endorsing_Rights, error) {
	var endorsingRights Endorsing_Rights

	err := json.Unmarshal(v, &endorsingRights)
	if err != nil {
		return endorsingRights, err
	}
	return endorsingRights, nil
}

//A helper sturcture to represent a delegate and their delegations by a range of cycles
type DelegationServiceRewards struct {
	DelegatePhk    string         `json:"delegate"`
	RewardsByCycle []CycleRewards `json:"cycles"`
}

//A structure representing rewards for a delegate and their delegations in a cycle
type CycleRewards struct {
	Cycle        int               `json:"cycle"`
	TotalRewards string            `json:"total_rewards"`
	Delegations  []ContractRewards `json:"delegations"`
}

//A structure representing a delegations gross rewards and share
type ContractRewards struct {
	DelegationPhk string  `json:"delegation"`
	Share         float64 `json:"share"`
	GrossRewards  string  `json:"rewards"`
	Balance       float64 `json:"balance"`
}

//A structure representing baking rights for a specific delegate between cycles
type BRights struct {
	Delegate string    `json:"delegate"`
	Cycles   []BCycles `json:"cycles"`
}

//A structure representing endorsing rights for a specific delegate between cycles
type ERights struct {
	Delegate string    `json:"delegate"`
	Cycles   []ECycles `json:"cycles"`
}

//A structure representing the baking rights in a specific cycle
type BCycles struct {
	Cycle        int           `json:"cycle"`
	BakingRights Baking_Rights `json:"baking_rights"`
}

//A structure representing the endorsing rights in a specific cycle
type ECycles struct {
	Cycle           int              `json:"cycle"`
	EndorsingRights Endorsing_Rights `json:"endorsing_rights"`
}

//Wallet needed for signing operations
type Wallet struct {
	Address  string
	Mnemonic string
	Seed     []byte
	Kp       sodium.SignKP
	Sk       string
	Pk       string
}

//Struct used to define transactions in a batch operation.
type Payment struct {
	Address string
	Amount  float64
}

type TezClientWrapper struct {
	healthy bool // isHealthy
	client  *TezosRPCClient
}

/*
 * GoTezos manages multiple Clients
 * each Client represents a Connection to a Tezos Node
 * GoTezos manages failover if one Node is down, there
 * are 2 Strategies:
 * failover: always use the same unless it is down -> go to the next - default
 * random: send to each Node equally
 */
type GoTezos struct {
	clientLock       sync.Mutex
	RpcClients       []*TezClientWrapper
	ActiveRPCCient   *TezClientWrapper
	Constants        NetworkConstants
	balancerStrategy string
	rand             *rand.Rand
	logger           *log.Logger
	cache            *gocache.Cache
	debug            bool
}


// Generic error from RPC. Returns an array/slice of error objects
type RPCGenericError struct {
	Kind	string	`json:"kind"`
	Error	string	`json:"error"`
}

func unMarshalRPCGenericErrors(v []byte) ([]RPCGenericError, error) {
	
	var r []RPCGenericError

	err := json.Unmarshal(v, &r)
	if err != nil {
		return r, err
	}
	
	return r, nil
}
