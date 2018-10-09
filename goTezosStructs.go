//Package goTezos exposes the Tezos RPC API in goLang.
package goTezos

import (
	"encoding/json"
	"log"
	"time"
)


//An unmarsheled representation of a block returned by the Tezos RPC API.
type Block struct {
	Protocol   string               `json:"protocol"`
	ChainID    string               `json:"chain_id"`
	Hash       string               `json:"hash"`
	Header     StructHeader         `json:"header"`
	Metadata   StructMetadata       `json:"metadata"`
	Operations [][]StructOperations `json:"operations"`
}

//An unmarsheled representation of a header in a block returned by the Tezos RPC API.
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
//An unmarsheled representation of Metadata in a block returned by the Tezos RPC API.
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

//An unmarsheled representation of a TestChainStatus found in the Metadata of a block returned by the Tezos RPC API.
type StructTestChainStatus struct {
	Status string `json:"status"`
}

//An unmarsheled representation of a MaxOperationListLength found in the Metadata of a block returned by the Tezos RPC API.
type StructMaxOperationListLength struct {
	MaxSize int `json:"max_size"`
	MaxOp   int `json:"max_op,omitempty"`
}

//An unmarsheled representation of a Level found in the Metadata of a block returned by the Tezos RPC API.
type StructLevel struct {
	Level                int  `json:"level"`
	LevelPosition        int  `json:"level_position"`
	Cycle                int  `json:"cycle"`
	CyclePosition        int  `json:"cycle_position"`
	VotingPeriod         int  `json:"voting_period"`
	VotingPeriodPosition int  `json:"voting_period_position"`
	ExpectedCommitment   bool `json:"expected_commitment"`
}

//An unmarsheled representation of BalanceUpdates found in the Metadata of a block returned by the Tezos RPC API.
type StructBalanceUpdates struct {
	Kind     string `json:"kind"`
	Contract string `json:"contract,omitempty"`
	Change   string `json:"change"`
	Category string `json:"category,omitempty"`
	Delegate string `json:"delegate,omitempty"`
	Level    int    `json:"level,omitempty"`
}

//An unmarsheled representation of Operations found in a block returned by the Tezos RPC API.
type StructOperations struct {
	Protocol  string           `json:"protocol"`
	ChainID   string           `json:"chain_id"`
	Hash      string           `json:"hash"`
	Branch    string           `json:"branch"`
	Contents  []StructContents `json:"contents"`
	Signature string           `json:"signature"`
}

//An unmarsheled representation of Contents found in a operation of a block returned by the Tezos RPC API.
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

//An unmarsheled representation of Metadata found in the Contents in a operation of a block returned by the Tezos RPC API.
type ContentsMetadata struct {
	BalanceUpdates []StructBalanceUpdates `json:"balance_updates"`
	Slots          []int                  `json:"slots"`
}

//Unmarshels the bytes received as a parameter, into the type Block. 
func unMarshelBlock(v []byte) (Block, error) {
	var block Block

	err := json.Unmarshal(v, &block)
	if err != nil {
		log.Println("Could not get unmarshel bytes into block: " + err.Error())
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

//An unmarsheled representation of a SnapShot returned by the Tezos RPC API.
type SnapShotQuery struct {
	RandomSeed   string `json:"random_seed"`
	RollSnapShot int    `json:"roll_snapshot"`
}

//Unmarshels the bytes received as a parameter, into the type SnapShotQuery. 
func unMarshelSnapShotQuery(v []byte) (SnapShotQuery, error) {
	var snapShotQuery SnapShotQuery

	err := json.Unmarshal(v, &snapShotQuery)
	if err != nil {
		log.Println("Could not unmarhel SnapShotQuery: " + err.Error())
		return snapShotQuery, err
	}
	return snapShotQuery, nil
}

//An unmarsheled representation of a FrozenBalanceRewards query returned by the Tezos RPC API.
type FrozenBalanceRewards struct {
	Deposits string `json:"deposits"`
	Fees     string `json:"fees"`
	Rewards  string `json:"rewards"`
}

//Unmarshels the bytes received as a parameter, into the type SnapShotQuery. 
func unMarshelFrozenBalanceRewards(v []byte) (FrozenBalanceRewards, error) {
	var frozenBalanceRewards FrozenBalanceRewards

	err := json.Unmarshal(v, &frozenBalanceRewards)
	if err != nil {
		log.Println("Could not unmarhel frozenBalanceRewards: " + err.Error())
		return frozenBalanceRewards, err
	}
	return frozenBalanceRewards, nil
}

//Unmarshels the bytes received as a parameter, into the type string. 
func unMarshelString(v []byte) (string, error) {
	var str string

	err := json.Unmarshal(v, &str)
	if err != nil {
		log.Println("Could not unmarshel to string " + err.Error())
		return str, err
	}
	return str, nil
}

//Unmarshels the bytes received as a parameter, into the type an array of strings. 
func unMarshelStringArray(v []byte) ([]string, error) {
	var strs []string

	err := json.Unmarshal(v, &strs)
	if err != nil {
		log.Println("Could not unmarshel to strings " + err.Error())
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

//A representation of delegations delegated to a delegate. 
type DelegatedContract struct {
	Address   string     
	Contracts []Contract 
	Delegate  bool       
	//  TimeStamp time.Time
	TotalPayout float64
	Fee         float64
}

//A representation of contracts by cycle for a delegation to a delegate. 
type Contract struct {
	Cycle           int
	Amount          float64
	RollInclusion   float64
	SharePercentage float64
	GrossPayout     float64
	NetPayout       float64
	Fee             float64
}

//A structure to hold a delegates details
type KnownAddress struct {
	Address string
	Alias   string
	Sk      string
}
