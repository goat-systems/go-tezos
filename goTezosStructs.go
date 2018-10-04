package goTezos

import (
	"encoding/json"
	"log"
	"time"
)

/*
Author: DefinitelyNotAGoat/MagicAglet
Version: 0.0.1
Description: This file contains structures used for the goTezos lib
License: MIT
*/

type Block struct {
	Protocol   string               `json:"protocol"`
	ChainID    string               `json:"chain_id"`
	Hash       string               `json:"hash"`
	Header     StructHeader         `json:"header"`
	Metadata   StructMetadata       `json:"metadata"`
	Operations [][]StructOperations `json:"operations"`
}

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

type StructTestChainStatus struct {
	Status string `json:"status"`
}

type StructMaxOperationListLength struct {
	MaxSize int `json:"max_size"`
	MaxOp   int `json:"max_op,omitempty"`
}

type StructLevel struct {
	Level                int  `json:"level"`
	LevelPosition        int  `json:"level_position"`
	Cycle                int  `json:"cycle"`
	CyclePosition        int  `json:"cycle_position"`
	VotingPeriod         int  `json:"voting_period"`
	VotingPeriodPosition int  `json:"voting_period_position"`
	ExpectedCommitment   bool `json:"expected_commitment"`
}

type StructBalanceUpdates struct {
	Kind     string `json:"kind"`
	Contract string `json:"contract,omitempty"`
	Change   string `json:"change"`
	Category string `json:"category,omitempty"`
	Delegate string `json:"delegate,omitempty"`
	Level    int    `json:"level,omitempty"`
}

type StructOperations struct {
	Protocol  string           `json:"protocol"`
	ChainID   string           `json:"chain_id"`
	Hash      string           `json:"hash"`
	Branch    string           `json:"branch"`
	Contents  []StructContents `json:"contents"`
	Signature string           `json:"signature"`
}

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

type ContentsMetadata struct {
	BalanceUpdates []StructBalanceUpdates `json:"balance_updates"`
	Slots          []int                  `json:"slots"`
}

func unMarshelBlock(v []byte) (Block, error) {
	var block Block

	err := json.Unmarshal(v, &block)
	if err != nil {
		log.Println("Could not get unmarshel bytes into block: " + err.Error())
		return block, err
	}
	return block, nil
}

type SnapShot struct {
	Cycle           int
	Number          int
	AssociatedHash  string
	AssociatedBlock int
}

type SnapShotQuery struct {
	RandomSeed   string `json:"random_seed"`
	RollSnapShot int    `json:"roll_snapshot"`
}

func unMarshelSnapShotQuery(v []byte) (SnapShotQuery, error) {
	var snapShotQuery SnapShotQuery

	err := json.Unmarshal(v, &snapShotQuery)
	if err != nil {
		log.Println("Could not unmarhel SnapShotQuery: " + err.Error())
		return snapShotQuery, err
	}
	return snapShotQuery, nil
}

type FrozenBalanceRewards struct {
	Deposits string `json:"deposits"`
	Fees     string `json:"fees"`
	Rewards  string `json:"rewards"`
}

func unMarshelFrozenBalanceRewards(v []byte) (FrozenBalanceRewards, error) {
	var frozenBalanceRewards FrozenBalanceRewards

	err := json.Unmarshal(v, &frozenBalanceRewards)
	if err != nil {
		log.Println("Could not unmarhel frozenBalanceRewards: " + err.Error())
		return frozenBalanceRewards, err
	}
	return frozenBalanceRewards, nil
}

func unMarshelString(v []byte) (string, error) {
	var str string

	err := json.Unmarshal(v, &str)
	if err != nil {
		log.Println("Could not unmarshel to string " + err.Error())
		return str, err
	}
	return str, nil
}

func unMarshelStringArray(v []byte) ([]string, error) {
	var strs []string

	err := json.Unmarshal(v, &strs)
	if err != nil {
		log.Println("Could not unmarshel to strings " + err.Error())
		return strs, err
	}
	return strs, nil
}

// OLD

//import "time"

/*
Description: A way to repesent each Delegated Contact, and their share for each cycle
Address: A string value representing the delegated contracts address
Commitments: An array of a structure that holds the amount commited for a cycle, and the percentage share
Delegator: Is this contract the delegator?
*/
type DelegatedContract struct {
	Address   string     //Public Key Hash
	Contracts []Contract //Percentage of total delegation for profit share for each cycle participated
	Delegate  bool       //If this client is yourself or not.
	//  TimeStamp time.Time
	TotalPayout float64
	Fee         float64
}

/*
Description: A representation of the amount commited in a cycle, and the percentage share for that amount.
Cycle: The cycle number
Amount: XTZ value of the amount commited in the cycle
SharePercentage: The percentage value of the amount to all commitments made in that cycle
Payout: Amount of rewards to be paid out for the commitment
Timestamp: A timestamp to show when the commitment was made
*/
type Contract struct {
	Cycle           int
	Amount          float64
	RollInclusion   float64
	SharePercentage float64
	GrossPayout     float64
	NetPayout       float64
	Fee             float64
}

/*
Description: A structure to represent a known address.
Address: a string value representing the address
Alias: the alias assigned to the known address
Sk: The protection around the Sk, unencrypted, legder, etc
*/
type KnownAddress struct {
	Address string
	Alias   string
	Sk      string
}
