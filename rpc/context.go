package rpc

import (
	"encoding/json"
	"fmt"
	"strconv"

	validator "github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"github.com/valyala/fastjson"
)

// TODO: needs test
func (c *Client) processContextRequest(input interface{}, cycle int, blockhash string) (string, error) {
	err := validator.New().Struct(input)
	if err != nil {
		return "", errors.Wrap(err, "invalid input")
	}

	err = validateContext(cycle, blockhash)
	if err != nil {
		return "", err
	}

	hash, err := c.extractContext(cycle, blockhash)
	if err != nil {
		return "", err
	}

	return hash, nil
}

func validateContext(cycle int, blockhash string) error {
	if blockhash == "" && cycle <= 0 {
		return errors.New("invalid input: missing key cycle or blockhash")
	} else if blockhash != "" && cycle > 0 {
		return errors.New("invalid input: cannot have both cycle and blockhash")
	}

	return nil
}

func (c *Client) extractContext(cycle int, blockhash string) (string, error) {
	if cycle != 0 {
		snapshot, err := c.Cycle(cycle)
		if err != nil {
			return "", errors.Wrapf(err, "failed to get extract hash for cycle '%d'", cycle)
		}

		return snapshot.BlockHash, nil
	}

	return blockhash, nil
}

/*
BigMapInput is the input for the BigMap function.

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-big-maps-big-map-id-script-expr
*/
type BigMapInput struct {
	// The block level of which you want to make the query. If not provided Cycle is required.
	Blockhash string
	// The cycle to get the balance at. If not provided Blockhash is required.
	Cycle int
	// The ID of the BigMap you wish to query.
	BigMapID int `validate:"required"`
	// The key. Look at the forge package for functions that end with Expression to forge this.
	ScriptExpression string `validate:"required"`
}

/*
BigMap reads data from a big_map.

Path:
 	../<block_id>/context/big_maps/<big_map_id>/<script_expr> (GET)

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-big-maps-big-map-id-script-expr
*/
func (c *Client) BigMap(input BigMapInput) ([]byte, error) {
	hash, err := c.processContextRequest(input, input.Cycle, input.Blockhash)
	if err != nil {
		return []byte{}, errors.Wrapf(err, "failed to get big map '%d' value with key '%s'", input.BigMapID, input.ScriptExpression)
	}

	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/context/big_maps/%d/%s", c.chain, hash, input.BigMapID, input.ScriptExpression))
	if err != nil {
		return []byte{}, errors.Wrapf(err, "failed to get big map '%d' value with key '%s'", input.BigMapID, input.ScriptExpression)
	}

	return resp, nil
}

/*
Constants represents the network constants.

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-constants
*/
type Constants struct {
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
	HardGasLimitPerOperation     int      `json:"hard_gas_limit_per_operation,string"`
	HardGasLimitPerBlock         int      `json:"hard_gas_limit_per_block,string"`
	ProofOfWorkThreshold         string   `json:"proof_of_work_threshold"`
	TokensPerRoll                string   `json:"tokens_per_roll"`
	MichelsonMaximumTypeSize     int      `json:"michelson_maximum_type_size"`
	SeedNonceRevelationTip       string   `json:"seed_nonce_revelation_tip"`
	OriginationSize              int      `json:"origination_size"`
	BlockSecurityDeposit         int      `json:"block_security_deposit,string"`
	EndorsementSecurityDeposit   int      `json:"endorsement_security_deposit,string"`
	BlockReward                  IntArray `json:"block_reward"`
	EndorsementReward            IntArray `json:"endorsement_reward"`
	CostPerByte                  int      `json:"cost_per_byte,string"`
	HardStorageLimitPerOperation int      `json:"hard_storage_limit_per_operation,string"`
}

// IntArray implements json.Marshaler so that a slice of string-ints can be a slice of ints
type IntArray []int

// UnmarshalJSON satisfies json.Marshaler
func (i *IntArray) UnmarshalJSON(data []byte) error {
	var array []string
	if err := json.Unmarshal(data, &array); err != nil {
		return err
	}

	for _, element := range array {
		if num, err := strconv.Atoi(element); err == nil {
			*i = append(*i, num)
		}
	}

	return nil
}

// MarshalJSON satisfies json.Marshaler
func (i *IntArray) MarshalJSON() ([]byte, error) {
	var array []string
	for _, num := range *i {
		array = append(array, strconv.Itoa(num))
	}

	return json.Marshal(array)
}

/*
ConstantsInput is the input for the Constants function.

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-constants
*/
type ConstantsInput struct {
	// The block level of which you want to make the query. If not provided Cycle is required.
	Blockhash string
	// The cycle to get the balance at. If not provided Blockhash is required.
	Cycle int
}

/*
Constants gets all constants.

Path:
	../<block_id>/context/constants (GET)

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-constants
*/
func (c *Client) Constants(input ConstantsInput) (Constants, error) {
	hash, err := c.processContextRequest(input, input.Cycle, input.Blockhash)
	if err != nil {
		return Constants{}, errors.Wrap(err, "failed to get constants")
	}

	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/context/constants", c.chain, hash))
	if err != nil {
		return Constants{}, errors.Wrapf(err, "failed to get constants")
	}

	var constants Constants
	err = json.Unmarshal(resp, &constants)
	if err != nil {
		return constants, errors.Wrapf(err, "failed to get constants: failed to parse json")
	}

	return constants, nil
}

/* ########### TODO ########### */
/* https://tezos.gitlab.io/008/rpc.html#get-block-id-context-constants-errors */

/*
ContractsInput is the input for the Contracts function.

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-contracts
*/
type ContractsInput struct {
	// The block level of which you want to make the query. If not provided Cycle is required.
	Blockhash string
	// The cycle to get the balance at. If not provided Blockhash is required.
	Cycle int
}

/*
Contracts gets all constants.

Path:
	../<block_id>/context/contracts (GET)

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-contracts
*/
func (c *Client) Contracts(input ContractsInput) ([]string, error) {
	hash, err := c.processContextRequest(input, input.Cycle, input.Blockhash)
	if err != nil {
		return []string{}, errors.Wrap(err, "failed to get contracts")
	}

	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/context/contracts", c.chain, hash))
	if err != nil {
		return []string{}, errors.Wrapf(err, "failed to get contracts")
	}

	var contracts []string
	err = json.Unmarshal(resp, &contracts)
	if err != nil {
		return []string{}, errors.Wrapf(err, "failed to get contracts: failed to parse json")
	}

	return contracts, nil
}

/*
Contract represents a contract.

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-contracts-contract-id
*/
type Contract struct {
	Balance  string `json:"balance"`
	Delegate string `json:"delegate"`
	Script   struct {
		Code    *json.RawMessage
		Stroage *json.RawMessage
	} `json:"script,omitempty"`
	Counter string `json:"counter,omitempty"`
}

/*
ContractInput is the input for the Contract function.

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-contracts-contract-id
*/
type ContractInput struct {
	// The block level of which you want to make the query. If not provided Cycle is required.
	Blockhash string
	// The cycle to get the balance at. If not provided Blockhash is required.
	Cycle int
	// The contract ID of the contract you wish to get.
	ContractID string `validate:"required"`
}

/*
Contract accesses the complete status of a contract.

Path:
	../<block_id>/context/contracts/<contract_id> (GET)

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-contracts-contract-id
*/
func (c *Client) Contract(input ContractInput) (Contract, error) {
	hash, err := c.processContextRequest(input, input.Cycle, input.Blockhash)
	if err != nil {
		return Contract{}, errors.Wrap(err, "failed to get contract")
	}

	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/context/contracts/%s", c.chain, hash, input.ContractID))
	if err != nil {
		return Contract{}, errors.Wrapf(err, "failed to get contract '%s'", input.ContractID)
	}

	var contract Contract
	err = json.Unmarshal(resp, &contract)
	if err != nil {
		return Contract{}, errors.Wrapf(err, "failed to get contract '%s': failed to parse json", input.ContractID)
	}

	return contract, nil
}

/*
ContractBalanceInput is the input for the Balance function.

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-contracts-contract-id-balance
*/
type ContractBalanceInput struct {
	// The block level of which you want to make the query. If not provided Cycle is required.
	Blockhash string
	// The cycle to get the balance at. If not provided Blockhash is required.
	Cycle int
	// The contract ID of the contract balance you wish to get.
	ContractID string `validate:"required"`
}

/*
ContractBalance accesses the balance of a contract.

Path:
	../<block_id>/context/contracts/<contract_id>/balance (GET)

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-contracts-contract-id-balance
*/
func (c *Client) ContractBalance(input ContractBalanceInput) (string, error) {
	hash, err := c.processContextRequest(input, input.Cycle, input.Blockhash)
	if err != nil {
		return "", errors.Wrap(err, "failed to get balance")
	}

	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/context/contracts/%s/balance", c.chain, hash, input.ContractID))
	if err != nil {
		return "", errors.Wrapf(err, "failed to get balance for contract '%s'", input.ContractID)
	}

	var balance string
	if err = json.Unmarshal(resp, &balance); err != nil {
		return "", errors.Wrapf(err, "failed to get balance for contract '%s': failed to parse json", input.ContractID)
	}

	return balance, nil
}

/*
ContractCounterInput is the input for the Counter function.

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-contracts-contract-id-counter
*/
type ContractCounterInput struct {
	// The block level of which you want to make the query. If not provided Cycle is required.
	Blockhash string
	// The cycle to get the balance at. If not provided Blockhash is required.
	Cycle int
	// The contract ID of the contract counter you wish to get.
	ContractID string `validate:"required"`
}

/*
ContractCounter accesses the counter of a contract, if any.

Path:
	../<block_id>/context/contracts/<contract_id>/counter (GET)

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-contracts-contract-id-counter
*/
func (c *Client) ContractCounter(input ContractCounterInput) (int, error) {
	hash, err := c.processContextRequest(input, input.Cycle, input.Blockhash)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get counter")
	}

	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/context/contracts/%s/counter", c.chain, hash, input.ContractID))
	if err != nil {
		return 0, errors.Wrapf(err, "failed to get counter for contract '%s'", input.ContractID)
	}

	var strCounter string
	err = json.Unmarshal(resp, &strCounter)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to get counter for contract '%s': failed to parse json", input.ContractID)
	}

	counter, err := strconv.Atoi(strCounter)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to get counter for contract '%s': failed to convert to int", input.ContractID)
	}
	return counter, nil
}

/*
ContractDelegateInput is the input for the ContractDelegate function.

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-contracts-contract-id-delegate
*/
type ContractDelegateInput struct {
	// The block level of which you want to make the query. If not provided Cycle is required.
	Blockhash string
	// The cycle to get the balance at. If not provided Blockhash is required.
	Cycle int
	// The contract ID of the contract delegate you wish to get.
	ContractID string `validate:"required"`
}

/*
ContractDelegate accesses the delegate of a contract, if any.

Path:
	../<block_id>/context/contracts/<contract_id>/delegate (GET)

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-contracts-contract-id-delegate
*/
func (c *Client) ContractDelegate(input ContractDelegateInput) (string, error) {
	hash, err := c.processContextRequest(input, input.Cycle, input.Blockhash)
	if err != nil {
		return "", errors.Wrap(err, "failed to get delegate")
	}

	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/context/contracts/%s/delegate", c.chain, hash, input.ContractID))
	if err != nil {
		return "", errors.Wrapf(err, "failed to get delegate for contract '%s'", input.ContractID)
	}

	var delegate string
	err = json.Unmarshal(resp, &delegate)
	if err != nil {
		return "", errors.Wrapf(err, "failed to get delegate for contract '%s': failed to parse json", input.ContractID)
	}

	return delegate, nil
}

/*
ContractEntrypointsInput is the input for the ContractDelegate function.

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-contracts-contract-id-entrypoints
*/
type ContractEntrypointsInput struct {
	// The block level of which you want to make the query
	Blockhash string `validate:"required"`
	// The contract ID of the contract delegate you wish to get.
	ContractID string `validate:"required"`
}

/*
ContractEntrypoints return a map of entrypoints of the contract where
the entrypoints are the keys and the micheline is the value.

Path:
	../<block_id>/context/contracts/<contract_id>/entrypoints (GET)

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-contracts-contract-id-entrypoints
*/
func (c *Client) ContractEntrypoints(input ContractEntrypointsInput) (map[string]*json.RawMessage, error) {
	err := validator.New().Struct(input)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get entrypoints: invalid input")
	}

	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/context/contracts/%s/entrypoints", c.chain, input.Blockhash, input.ContractID))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get entrypoints for contract '%s'", input.ContractID)
	}

	entrypoints := make(map[string]*json.RawMessage)
	var p fastjson.Parser
	parsedJSON, err := p.Parse(string(resp))
	if err != nil {
		return entrypoints, errors.Wrapf(err, "failed to get entrypoints for contract '%s': failed to parse json", input.ContractID)
	}

	obj, err := parsedJSON.Object()
	if err != nil {
		return entrypoints, errors.Wrapf(err, "failed to get entrypoints for contract '%s': unrecognized json", input.ContractID)
	}

	if v := obj.Get("entrypoints"); v != nil {
		obj, err = v.Object()
		if err != nil {
			return entrypoints, errors.Wrapf(err, "failed to get entrypoints for contract '%s': unrecognized json", input.ContractID)
		}
		obj.Visit(func(key []byte, v *fastjson.Value) {
			rawMessage := &json.RawMessage{}
			*rawMessage = v.MarshalTo([]byte{})
			entrypoints[string(key)] = rawMessage
		})
	}

	return entrypoints, nil
}

/*
ContractEntrypointInput is the input for the ContractDelegate function.

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-contracts-contract-id-entrypoints-string
*/
type ContractEntrypointInput struct {
	// The block level of which you want to make the query
	Blockhash string `validate:"required"`
	// The contract ID of the contract delegate you wish to get.
	ContractID string `validate:"required"`
	// The entrypoint of the contract you wish to get.
	Entrypoint string `validate:"required"`
}

/*
ContractEntrypoint returns the type of the given entrypoint of the contract.

Path:
	../<block_id>/context/contracts/<contract_id>/entrypoints/<string> (GET)

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-contracts-contract-id-entrypoints
*/
func (c *Client) ContractEntrypoint(input ContractEntrypointInput) (*json.RawMessage, error) {
	err := validator.New().Struct(input)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get entrypoint: invalid input")
	}

	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/context/contracts/%s/entrypoints/%s", c.chain, input.Blockhash, input.ContractID, input.Entrypoint))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get entrypoint '%s' for contract '%s'", input.Entrypoint, input.ContractID)
	}

	rawMessage := &json.RawMessage{}
	*rawMessage = resp

	return rawMessage, nil
}

/*
Delegate represents the frozen delegate RPC on the tezos network.

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-delegates-pkh
*/
type Delegate struct {
	Balance              string `json:"balance"`
	FrozenBalance        string `json:"frozen_balance"`
	FrozenBalanceByCycle []struct {
		Cycle   int `json:"cycle"`
		Deposit int `json:"deposit,string"`
		Fees    int `json:"fees,string"`
		Rewards int `json:"rewards,string"`
	} `json:"frozen_balance_by_cycle"`
	StakingBalance    string   `json:"staking_balance"`
	DelegateContracts []string `json:"delegated_contracts"`
	DelegatedBalance  string   `json:"delegated_balance"`
	Deactivated       bool     `json:"deactivated"`
	GracePeriod       int      `json:"grace_period"`
}

/*
DelegateInput is the input for the Delegate function.

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-contracts-contract-id-delegate
*/
type DelegateInput struct {
	// The block level of which you want to make the query. If empty Cycle is required.
	Blockhash string
	// The cycle to get the balance at. If empty Blockhash is required.
	Cycle int
	// The delegate that you want to make the query.
	Delegate string `validate:"required"`
}

/*
Delegate accessed the delegate of a contract, if any.

Path:
	../<block_id>/context/contracts/<contract_id>/delegate (GET)

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-contracts-contract-id-delegate
*/
func (c *Client) Delegate(input DelegateInput) (Delegate, error) {
	hash, err := c.processContextRequest(input, input.Cycle, input.Blockhash)
	if err != nil {
		return Delegate{}, errors.Wrap(err, "failed to get delegate")
	}

	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/context/delegates/%s", c.chain, hash, input.Delegate))
	if err != nil {
		return Delegate{}, errors.Wrapf(err, "could not get delegate '%s'", input.Delegate)
	}

	var delegate Delegate
	err = json.Unmarshal(resp, &delegate)
	if err != nil {
		return delegate, errors.Wrapf(err, "failed to get delegate '%s': failed to parse json", input.Delegate)
	}

	return delegate, nil
}

/*
ContractStorageInput is the input for the client.ContractStorage() function.

Function:
	func (c *Client) ContractStorage(input ContractStorageInput) ([]byte, error)  {}
*/
type ContractStorageInput struct {
	// The block level of which you want to make the query. If not provided Cycle is required.
	Blockhash string
	// The cycle to get the balance at. If not provided Blockhash is required.
	Cycle int
	// The contract at to get the storage for
	Contract string `validate:"required"`
}

func (c *ContractStorageInput) validate() error {
	if c.Blockhash == "" && c.Cycle == 0 {
		return errors.New("invalid input: missing key cycle or blockhash")
	} else if c.Blockhash != "" && c.Cycle != 0 {
		return errors.New("invalid input: cannot have both cycle and blockhash")
	}

	err := validator.New().Struct(c)
	if err != nil {
		return errors.Wrap(err, "invalid input")
	}

	return nil
}

/*
ContractStorage gets access the data of the contract.

Path:
	../<block_id>/context/contracts/<contract_id>/storage (GET)

Link:
	https://tezos.gitlab.io/api/rpc.html#get-block-id-context-contracts-contract-id-storage
*/
func (c *Client) ContractStorage(input ContractStorageInput) ([]byte, error) {
	err := validator.New().Struct(input)
	if err != nil {
		return []byte{}, errors.Wrap(err, "invalid input")
	}

	query := fmt.Sprintf("/chains/%s/blocks/%s/context/contracts/%s/storage", c.chain, input.Blockhash, input.Contract)
	resp, err := c.get(query)
	if err != nil {
		return []byte{}, errors.Wrap(err, "could not get storage '%s'")
	}

	return resp, nil
}
