package rpc

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	validator "github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"github.com/valyala/fastjson"
)

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
ContractManagerKeyInput is the input for the ContractManagerKey function.

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-contracts-contract-id-manager-key
*/
type ContractManagerKeyInput struct {
	// The block level of which you want to make the query
	Blockhash string `validate:"required"`
	// The contract ID of the contract delegate you wish to get.
	ContractID string `validate:"required"`
}

/*
ContractManagerKey accesses the manager of a contract.

Path:
	../<block_id>/context/contracts/<contract_id>/manager_key (GET)

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-contracts-contract-id-entrypoints
*/
func (c *Client) ContractManagerKey(input ContractManagerKeyInput) (string, error) {
	err := validator.New().Struct(input)
	if err != nil {
		return "", errors.Wrap(err, "failed to get manager: invalid input")
	}

	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/context/contracts/%s/manager_key", c.chain, input.Blockhash, input.ContractID))
	if err != nil {
		return "", errors.Wrapf(err, "failed to get manager for contract '%s'", input.ContractID)
	}

	var manager string
	err = json.Unmarshal(resp, &manager)
	if err != nil {
		return "", errors.Wrapf(err, "failed to get manager for contract '%s': failed to parse json", input.ContractID)
	}

	return manager, nil
}

/*
ContractScriptInput is the input for the ContractScript function.

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-contracts-contract-id-script
*/
type ContractScriptInput struct {
	// The block level of which you want to make the query
	Blockhash string
	// The cycle to get the balance at. If not provided Blockhash is required.
	Cycle int
	// The contract ID of the contract delegate you wish to get.
	ContractID string `validate:"required"`
}

/*
ContractScript accesses the code and data of the contract.

Path:
	../<block_id>/context/contracts/<contract_id>/script (GET)

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-contracts-contract-id-script
*/
func (c *Client) ContractScript(input ContractScriptInput) (*json.RawMessage, error) {
	hash, err := c.processContextRequest(input, input.Cycle, input.Blockhash)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get script")
	}

	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/context/contracts/%s/script", c.chain, hash, input.ContractID))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get script for contract '%s'", input.ContractID)
	}

	rawMessage := &json.RawMessage{}
	*rawMessage = resp

	return rawMessage, nil
}

/*
ContractSaplingDiffInput is the input for the ContractSaplingDiff function.

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-contracts-contract-id-single-sapling-get-diff
*/
type ContractSaplingDiffInput struct {
	// The block level of which you want to make the query. If not provided Cycle is required.
	Blockhash string
	// The cycle to get the balance at. If not provided Blockhash is required.
	Cycle int
	// The contract ID of the contract delegate you wish to get.
	ContractID string `validate:"required"`
	//  Commitments and ciphertexts are returned from the specified offset up to the most recent.
	OffsetCommitment int
	// Nullifiers are returned from the specified offset up to the most recent.
	OffsetNullifier int
}

/*
SingleSaplingDiff represents a a sapling diff for a contract.

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-contracts-contract-id-single-sapling-get-diff
*/
type SingleSaplingDiff struct {
	Root                      string                      `json:"root"`
	CommitmentsAndCiphertexts []CommitmentsAndCiphertexts `json:"commitments_and_ciphertexts"`
	Nullifiers                []string                    `json:"nullifiers"`
}

// UnmarshalJSON satisfies json.Marshler
func (s *SingleSaplingDiff) UnmarshalJSON(b []byte) error {
	v, err := fastjson.Parse(string(b))
	if err != nil {
		return errors.Wrap(err, "failed to parse json")
	}

	rootValue := v.GetStringBytes("root")
	if rootValue == nil {
		return errors.Wrap(err, "failed to parse json")
	}

	var singleSaplingDiff SingleSaplingDiff
	singleSaplingDiff.Root = string(rootValue)

	commitmentsAndCiphertexts := v.GetArray("commitments_and_ciphertexts")
	if commitmentsAndCiphertexts != nil {
		for _, value := range commitmentsAndCiphertexts {
			innerValues, err := value.Array()
			if err != nil {
				return errors.Wrap(err, "failed to parse json")
			}

			if len(innerValues) == 2 {
				var commitmentsAndCiphertexts CommitmentsAndCiphertexts
				commitmentsAndCiphertexts.Commitment = innerValues[0].String()
				var cipherText CipherText
				err = json.Unmarshal(innerValues[1].MarshalTo([]byte{}), &cipherText)
				if err != nil {
					return errors.Wrap(err, "failed to parse json")
				}

				commitmentsAndCiphertexts.CipherText = cipherText
				singleSaplingDiff.CommitmentsAndCiphertexts = append(singleSaplingDiff.CommitmentsAndCiphertexts, commitmentsAndCiphertexts)
			}
		}
	}

	nullifiers := v.GetArray("nullifiers")
	if nullifiers != nil {
		for _, value := range nullifiers {
			singleSaplingDiff.Nullifiers = append(singleSaplingDiff.Nullifiers, value.String())
		}
	}

	*s = singleSaplingDiff
	return nil
}

// CommitmentsAndCiphertexts is a group of a commitment and a CipherText
type CommitmentsAndCiphertexts struct {
	Commitment string
	CipherText CipherText
}

// CipherText is a sapling Cipher Text
type CipherText struct {
	CV         string
	EPK        string
	PayloadEnc string
	NonceEnc   string
	PayloadOut string
	NonceOut   string
}

/*
ContractSaplingDiff returns the root and a diff of a state starting from an optional offset which is zero by default.

###
NOTE: This function is not production ready because sapling contracts are not readily available yet.
###

Path:
	 ../<block_id>/context/contracts/<contract_id>/single_sapling_get_diff?[offset_commitment=<int64>]&[offset_nullifier=<int64>] (GET)

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-contracts-contract-id-single-sapling-get-diff
*/
func (c *Client) ContractSaplingDiff(input ContractSaplingDiffInput) (SingleSaplingDiff, error) {
	hash, err := c.processContextRequest(input, input.Cycle, input.Blockhash)
	if err != nil {
		return SingleSaplingDiff{}, errors.Wrap(err, "failed to get script")
	}

	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/context/contracts/%s/single_sapling_get_diff", c.chain, hash, input.ContractID), input.contructRPCOptions()...)
	if err != nil {
		return SingleSaplingDiff{}, errors.Wrapf(err, "failed to get single sapling diff for contract '%s'", input.ContractID)
	}

	var saplingDiff SingleSaplingDiff
	err = json.Unmarshal(resp, &saplingDiff)
	if err != nil {
		return SingleSaplingDiff{}, errors.Wrapf(err, "failed to get single sapling diff for contract '%s': failed to parse json", input.ContractID)
	}

	return saplingDiff, nil
}

func (c *ContractSaplingDiffInput) contructRPCOptions() []rpcOptions {
	if c.OffsetCommitment != 0 && c.OffsetNullifier != 0 {
		return []rpcOptions{
			{
				"offset_commitment",
				strconv.Itoa(c.OffsetCommitment),
			},
			{
				"offset_nullifier",
				strconv.Itoa(c.OffsetNullifier),
			},
		}
	} else if c.OffsetCommitment != 0 && c.OffsetNullifier == 0 {
		return []rpcOptions{
			{
				"offset_commitment",
				strconv.Itoa(c.OffsetCommitment),
			},
		}
	} else if c.OffsetCommitment != 0 && c.OffsetNullifier != 0 {
		return []rpcOptions{
			{
				"offset_nullifier",
				strconv.Itoa(c.OffsetNullifier),
			},
		}
	}

	return nil
}

/*
ContractStorageInput is the input for the ContractStorage function.

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-contracts-contract-id-storage
*/
type ContractStorageInput struct {
	// The block level of which you want to make the query. If not provided Cycle is required.
	Blockhash string
	// The cycle to get the balance at. If not provided Blockhash is required.
	Cycle int
	// The contract ID of the contract delegate you wish to get.
	ContractID string `validate:"required"`
}

/*
ContractStorage accesses the data of the contract.

Path:
	../<block_id>/context/contracts/<contract_id>/storage (GET)

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-contracts-contract-id-storage
*/
func (c *Client) ContractStorage(input ContractStorageInput) (*json.RawMessage, error) {
	hash, err := c.processContextRequest(input, input.Cycle, input.Blockhash)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get storage")
	}

	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/context/contracts/%s/storage", c.chain, hash, input.ContractID))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get storage for contract '%s'", input.ContractID)
	}

	rawMessage := &json.RawMessage{}
	*rawMessage = resp

	return rawMessage, nil
}

/*
DelegatesInput is the input for the Delegates function.

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-delegates
*/
type DelegatesInput struct {
	// The block level of which you want to make the query. If empty Cycle is required.
	Blockhash string
	// The cycle to get the balance at. If empty Blockhash is required.
	Cycle int
	// The block level of which you want to make the query.
	active bool
	// The cycle of which you want to make the query.
	inactive bool
}

/*
Delegates lists all registered delegates.

Path:
	../<block_id>/context/delegates (GET)

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-delegates
*/
func (c *Client) Delegates(input DelegatesInput) ([]string, error) {
	hash, err := c.processContextRequest(input, input.Cycle, input.Blockhash)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get delegates")
	}

	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/context/delegates", c.chain, hash), input.contructRPCOptions()...)
	if err != nil {
		return []string{}, errors.Wrap(err, "failed to get delegates")
	}

	var delegates []string
	err = json.Unmarshal(resp, &delegates)
	if err != nil {
		return []string{}, errors.Wrap(err, "failed to get delegates: failed to parse json")
	}

	return delegates, nil
}

func (d *DelegatesInput) contructRPCOptions() []rpcOptions {
	var opts []rpcOptions
	if d.active {
		opts = append(opts, rpcOptions{
			"active",
			"true",
		})
	} else {
		opts = append(opts, rpcOptions{
			"active",
			"false",
		})
	}

	if d.inactive {
		opts = append(opts, rpcOptions{
			"inactive",
			"true",
		})
	} else {
		opts = append(opts, rpcOptions{
			"inactive",
			"false",
		})
	}

	return opts
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
Delegate returns everything about a delegate.

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
		return Delegate{}, errors.Wrapf(err, "failed to get delegate '%s'", input.Delegate)
	}

	var delegate Delegate
	err = json.Unmarshal(resp, &delegate)
	if err != nil {
		return delegate, errors.Wrapf(err, "failed to get delegate '%s': failed to parse json", input.Delegate)
	}

	return delegate, nil
}

/*
DelegateBalanceInput is the input for the DelegateBalance function.

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-delegates-pkh-balance
*/
type DelegateBalanceInput struct {
	// The block level of which you want to make the query. If empty Cycle is required.
	Blockhash string
	// The cycle to get the balance at. If empty Blockhash is required.
	Cycle int
	// The delegate that you want to make the query.
	Delegate string `validate:"required"`
}

/*
DelegateBalance returns the full balance of a given delegate, including the frozen balances.

Path:
	../<block_id>/context/delegates/<pkh>/balance (GET)

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-delegates-pkh-balance
*/
func (c *Client) DelegateBalance(input DelegateBalanceInput) (string, error) {
	hash, err := c.processContextRequest(input, input.Cycle, input.Blockhash)
	if err != nil {
		return "", errors.Wrap(err, "failed to get delegate balance")
	}

	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/context/delegates/%s/balance", c.chain, hash, input.Delegate))
	if err != nil {
		return "", errors.Wrapf(err, "failed to get delegate '%s' balance", input.Delegate)
	}

	var balance string
	err = json.Unmarshal(resp, &balance)
	if err != nil {
		return "", errors.Wrapf(err, "failed to get delegate '%s' balance: failed to parse json", input.Delegate)
	}

	return balance, nil
}

/*
DelegateDeactivatedInput is the input for the DelegateDeactivated function.

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-delegates-pkh-deactivated
*/
type DelegateDeactivatedInput struct {
	// The block level of which you want to make the query. If empty Cycle is required.
	Blockhash string
	// The cycle to get the balance at. If empty Blockhash is required.
	Cycle int
	// The delegate that you want to make the query.
	Delegate string `validate:"required"`
}

/*
DelegateDeactivated tells whether the delegate is currently tagged as deactivated or not.

Path:
	../<block_id>/context/delegates/<pkh>/deactivated (GET)

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-delegates-pkh-deactivated
*/
func (c *Client) DelegateDeactivated(input DelegateDeactivatedInput) (bool, error) {
	hash, err := c.processContextRequest(input, input.Cycle, input.Blockhash)
	if err != nil {
		return false, errors.Wrap(err, "failed to get delegate activation status")
	}

	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/context/delegates/%s/deactivated", c.chain, hash, input.Delegate))
	if err != nil {
		return false, errors.Wrapf(err, "failed to get delegate '%s' activation status", input.Delegate)
	}

	var deactivated bool
	err = json.Unmarshal(resp, &deactivated)
	if err != nil {
		return false, errors.Wrapf(err, "failed to get delegate '%s' activation status: failed to parse json", input.Delegate)
	}

	return deactivated, nil
}

/*
DelegateDelegatedBalanceInput is the input for the DelegateDelegatedBalance function.

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-delegates-pkh-delegated-balance
*/
type DelegateDelegatedBalanceInput struct {
	// The block level of which you want to make the query. If empty Cycle is required.
	Blockhash string
	// The cycle to get the balance at. If empty Blockhash is required.
	Cycle int
	// The delegate that you want to make the query.
	Delegate string `validate:"required"`
}

/*
DelegateDelegatedBalance returns the balances of all the contracts that delegate to a given delegate.
This excludes the delegate's own balance and its frozen balances.

Path:
	../<block_id>/context/delegates/<pkh>/delegated_contracts (GET)

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-delegates-pkh-delegated-balance
*/
func (c *Client) DelegateDelegatedBalance(input DelegateBalanceInput) (string, error) {
	hash, err := c.processContextRequest(input, input.Cycle, input.Blockhash)
	if err != nil {
		return "", errors.Wrap(err, "failed to get delegate delegated balance")
	}

	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/context/delegates/%s/delegated_balance", c.chain, hash, input.Delegate))
	if err != nil {
		return "", errors.Wrapf(err, "failed to get delegate '%s' delegated balance", input.Delegate)
	}

	var balance string
	err = json.Unmarshal(resp, &balance)
	if err != nil {
		return "", errors.Wrapf(err, "failed to get delegate '%s' delegated balance: failed to parse json", input.Delegate)
	}

	return balance, nil
}

/*
DelegateDelegatedContractsInput is the input for the DelegateDelegatedContracts function.

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-delegates-pkh-delegated-contracts
*/
type DelegateDelegatedContractsInput struct {
	// The block level of which you want to make the query. If empty Cycle is required.
	Blockhash string
	// The cycle to get the balance at. If empty Blockhash is required.
	Cycle int
	// The delegate that you want to make the query.
	Delegate string `validate:"required"`
}

/*
DelegateDelegatedContracts returns the list of contracts that delegate to a given delegate.

Path:
	../<block_id>/context/delegates/<pkh>/delegated_contracts (GET)

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-delegates-pkh-delegated-contracts
*/
func (c *Client) DelegateDelegatedContracts(input DelegateDelegatedContractsInput) ([]string, error) {
	hash, err := c.processContextRequest(input, input.Cycle, input.Blockhash)
	if err != nil {
		return []string{}, errors.Wrap(err, "failed to get delegate delegated contracts")
	}

	var resp []byte
	resp, err = c.get(fmt.Sprintf("/chains/%s/blocks/%s/context/delegates/%s/delegated_contracts", c.chain, hash, input.Delegate))
	if err != nil {
		return []string{}, errors.Wrapf(err, "failed to get delegate '%s' delegated contracts", input.Delegate)
	}

	var delegatedContracts []string
	err = json.Unmarshal(resp, &delegatedContracts)
	if err != nil {
		return []string{}, errors.Wrapf(err, "failed to get delegate '%s' delegated contracts: failed to parse json", input.Delegate)
	}

	return delegatedContracts, nil
}

/*
DelegateFrozenBalanceInput is the input for the DelegateFrozenBalance function.

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-delegates-pkh-frozen-balance
*/
type DelegateFrozenBalanceInput struct {
	// The block level of which you want to make the query. If empty Cycle is required.
	Blockhash string
	// The cycle to get the balance at. If empty Blockhash is required.
	Cycle int
	// The delegate that you want to make the query.
	Delegate string `validate:"required"`
}

/*
DelegateFrozenBalance returns the total frozen balances of a given delegate, this includes the
frozen deposits, rewards and fees.

Path:
	../<block_id>/context/delegates/<pkh>/frozen_balance (GET)

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-delegates-pkh-frozen-balance
*/
func (c *Client) DelegateFrozenBalance(input DelegateFrozenBalanceInput) (string, error) {
	hash, err := c.processContextRequest(input, input.Cycle, input.Blockhash)
	if err != nil {
		return "", errors.Wrap(err, "failed to get delegate frozen balance")
	}

	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/context/delegates/%s/frozen_balance", c.chain, hash, input.Delegate))
	if err != nil {
		return "", errors.Wrapf(err, "failed to get delegate '%s' frozen balance", input.Delegate)
	}

	var balance string
	err = json.Unmarshal(resp, &balance)
	if err != nil {
		return "", errors.Wrapf(err, "failed to get delegate '%s' frozen balance: failed to parse json", input.Delegate)
	}

	return balance, nil
}

/*
FrozenBalanceByCycle represents the frozen balance of a delegate at a cycle.

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-delegates-pkh-frozen-balance-by-cycle
*/
type FrozenBalanceByCycle struct {
	Cycle   int
	Deposit string
	Fees    string
	Rewards string
}

/*
DelegateFrozenBalanceByCycleInput is the input for the DelegateFrozenBalanceByCycle function.

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-delegates-pkh-frozen-balance-by-cycle
*/
type DelegateFrozenBalanceByCycleInput struct {
	// The block level of which you want to make the query. If empty Cycle is required.
	Blockhash string
	// The cycle to get the balance at. If empty Blockhash is required.
	Cycle int
	// The delegate that you want to make the query.
	Delegate string `validate:"required"`
}

/*
DelegateFrozenBalanceByCycle returns the frozen balances of a given delegate,
indexed by the cycle by which it will be unfrozen

Path:
	../<block_id>/context/delegates/<pkh>/frozen_balance_by_cycle (GET)

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-delegates-pkh-frozen-balance-by-cycle
*/
func (c *Client) DelegateFrozenBalanceByCycle(input DelegateFrozenBalanceByCycleInput) ([]FrozenBalanceByCycle, error) {
	hash, err := c.processContextRequest(input, input.Cycle, input.Blockhash)
	if err != nil {
		return []FrozenBalanceByCycle{}, errors.Wrap(err, "failed to get delegate frozen balance at cycle")
	}

	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/context/delegates/%s/frozen_balance_by_cycle", c.chain, hash, input.Delegate))
	if err != nil {
		return []FrozenBalanceByCycle{}, errors.Wrapf(err, "failed to get delegate '%s' frozen balance at cycle", input.Delegate)
	}

	var frozenBalanceAtCycle []FrozenBalanceByCycle
	err = json.Unmarshal(resp, &frozenBalanceAtCycle)
	if err != nil {
		return []FrozenBalanceByCycle{}, errors.Wrapf(err, "failed to get delegate '%s' frozen balance at cycle: failed to parse json", input.Delegate)
	}

	return frozenBalanceAtCycle, nil
}

/*
DelegateGracePeriodInput is the input for the DelegateGracePeriod function.

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-delegates-pkh-grace-period
*/
type DelegateGracePeriodInput struct {
	// The block level of which you want to make the query. If empty Cycle is required.
	Blockhash string
	// The cycle to get the balance at. If empty Blockhash is required.
	Cycle int
	// The delegate that you want to make the query.
	Delegate string `validate:"required"`
}

/*
DelegateGracePeriod returns the cycle by the end of which the delegate might be deactivated
if she fails to execute any delegate action. A deactivated delegate might be reactivated
(without loosing any rolls) by simply re-registering as a delegate. For deactivated delegates,
this value contains the cycle by which they were deactivated.


Path:
	../<block_id>/context/delegates/<pkh>/grace_period (GET)

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-delegates-pkh-grace-period
*/
func (c *Client) DelegateGracePeriod(input DelegateGracePeriodInput) (int, error) {
	hash, err := c.processContextRequest(input, input.Cycle, input.Blockhash)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get delegate grace period")
	}

	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/context/delegates/%s/grace_period", c.chain, hash, input.Delegate))
	if err != nil {
		return 0, errors.Wrapf(err, "failed to get delegate '%s' grace period", input.Delegate)
	}

	var gracePeriod int
	err = json.Unmarshal(resp, &gracePeriod)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to get delegate '%s' grace period: failed to parse json", input.Delegate)
	}

	return gracePeriod, nil
}

/*
DelegateStakingBalanceInput is the input for the DelegateStakingBalance function.

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-delegates-pkh-staking-balance
*/
type DelegateStakingBalanceInput struct {
	// The block level of which you want to make the query. If empty Cycle is required.
	Blockhash string
	// The cycle to get the balance at. If empty Blockhash is required.
	Cycle int
	// The delegate that you want to make the query.
	Delegate string `validate:"required"`
}

/*
DelegateStakingBalance returns the total amount of tokens delegated to a given delegate.
This includes the balances of all the contracts that delegate to it, but also the balance of
the delegate itself and its frozen fees and deposits. The rewards do not count in the delegated
balance until they are unfrozen.


Path:
	../<block_id>/context/delegates/<pkh>/staking_balance (GET)

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-delegates-pkh-staking-balance
*/
func (c *Client) DelegateStakingBalance(input DelegateStakingBalanceInput) (string, error) {
	hash, err := c.processContextRequest(input, input.Cycle, input.Blockhash)
	if err != nil {
		return "", errors.Wrap(err, "failed to get delegate staking balance")
	}

	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/context/delegates/%s/staking_balance", c.chain, hash, input.Delegate))
	if err != nil {
		return "", errors.Wrapf(err, "failed to get delegate '%s' staking balance", input.Delegate)
	}

	var stakingBalance string
	err = json.Unmarshal(resp, &stakingBalance)
	if err != nil {
		return "", errors.Wrapf(err, "failed to get delegate '%s' staking balance: failed to parse json", input.Delegate)
	}

	return stakingBalance, nil
}

/*
DelegateVotingPowerInput is the input for the DelegateVotingPower function.

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-delegates-pkh-voting-power
*/
type DelegateVotingPowerInput struct {
	// The block level of which you want to make the query. If empty Cycle is required.
	Blockhash string
	// The cycle to get the balance at. If empty Blockhash is required.
	Cycle int
	// The delegate that you want to make the query.
	Delegate string `validate:"required"`
}

/*
DelegateVotingPower returns the number of rolls in the vote listings for a given delegate

Path:
	../<block_id>/context/delegates/<pkh>/voting_power (GET)

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-delegates-pkh-voting-power
*/
func (c *Client) DelegateVotingPower(input DelegateVotingPowerInput) (int, error) {
	hash, err := c.processContextRequest(input, input.Cycle, input.Blockhash)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get delegate voting power")
	}

	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/context/delegates/%s/voting_power", c.chain, hash, input.Delegate))
	if err != nil {
		return 0, errors.Wrapf(err, "failed to get delegate '%s' voting power", input.Delegate)
	}

	var votingPower int
	err = json.Unmarshal(resp, &votingPower)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to get delegate '%s' voting power: failed to parse json", input.Delegate)
	}

	return votingPower, nil
}

/*
Nonces represents nonces in the RPC

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-nonces-block-level
*/
type Nonces struct {
	Nonce     string
	Hash      string
	Forgotten bool
}

// UnmarshalJSON satisfies json.Marshaler
func (n *Nonces) UnmarshalJSON(data []byte) error {
	v, err := fastjson.Parse(string(data))
	if err != nil {
		return errors.Wrap(err, "failed to parse json")
	}

	*n = Nonces{Forgotten: true}

	if nonce := v.Get("nonce"); nonce != nil {
		*n = Nonces{Nonce: strings.Trim(nonce.String(), "\""), Forgotten: false}
	}

	if hash := v.Get("hash"); hash != nil {
		*n = Nonces{Hash: strings.Trim(hash.String(), "\""), Forgotten: false}
	}

	return nil
}

/*
NoncesInput is the input for the Nonces function.

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-nonces-block-level
*/
type NoncesInput struct {
	// The block level of which you want to make the query. If empty Cycle is required.
	Blockhash string
	// The cycle to get the balance at. If empty Blockhash is required.
	Cycle int
	// The delegate that you want to make the query.
	Level int `validate:"required"`
}

/*
Nonces returns the number of rolls in the vote listings for a given delegate

Path:
	../<block_id>/context/nonces/<block_level> (GET)

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-delegates-pkh-voting-power
*/
func (c *Client) Nonces(input NoncesInput) (Nonces, error) {
	hash, err := c.processContextRequest(input, input.Cycle, input.Blockhash)
	if err != nil {
		return Nonces{}, errors.Wrapf(err, "failed to get nonces at level '%d'", input.Level)
	}

	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/context/nonces/%d", c.chain, hash, input.Level))
	if err != nil {
		return Nonces{}, errors.Wrapf(err, "failed to get nonces at level '%d'", input.Level)
	}

	var nonces Nonces
	err = json.Unmarshal(resp, &nonces)
	if err != nil {
		return Nonces{}, errors.Wrapf(err, "failed to get nonces at level '%d': failed to parse json", input.Level)
	}

	return nonces, nil
}

/*
RawBytesInput is the input for the RawBytes function.

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-raw-bytes
*/
type RawBytesInput struct {
	// The block level of which you want to make the query. If empty Cycle is required.
	Blockhash string
	// The cycle to get the balance at. If empty Blockhash is required.
	Cycle int
	// The depth at which you want the raw bytes.
	Depth int
}

/*
RawBytes returns the raw context.

Path:
	../<block_id>/context/raw/bytes?[depth=<int>] (GET)

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-raw-bytes
*/
func (c *Client) RawBytes(input RawBytesInput) ([]byte, error) {
	hash, err := c.processContextRequest(input, input.Cycle, input.Blockhash)
	if err != nil {
		return []byte{}, errors.Wrap(err, "failed to raw at bytes")
	}

	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/context/raw/bytes", c.chain, hash), input.constructRPCOptions()...)
	if err != nil {
		return []byte{}, errors.Wrap(err, "failed to raw at bytes")
	}

	return resp, nil
}

func (r *RawBytesInput) constructRPCOptions() []rpcOptions {
	var opts []rpcOptions
	if r.Depth != 0 {
		opts = append(opts, rpcOptions{
			"depth",
			strconv.Itoa(r.Depth),
		})
	}
	return opts
}

/*
SaplingDiffInput is the input for the SaplingDiff function.

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-sapling-sapling-state-id-get-diff
*/
type SaplingDiffInput struct {
	// The block level of which you want to make the query. If not provided Cycle is required.
	Blockhash string
	// The cycle to get the balance at. If not provided Blockhash is required.
	Cycle int
	// The sapling state ID of the sapling you wish to get.
	SaplingStateID string `validate:"required"`
	//  Commitments and ciphertexts are returned from the specified offset up to the most recent.
	OffsetCommitment int
	// Nullifiers are returned from the specified offset up to the most recent.
	OffsetNullifier int
}

/*
SaplingDiff returns the root and a diff of a state starting from an optional offset which is zero by default.

###
TODO: Maybe just pass the bytes up until I can get example json to test with.
NOTE: This function is not production ready because sapling contracts are not readily available yet.
###

Path:
	../<block_id>/context/sapling/<sapling_state_id>/get_diff?[offset_commitment=<int64>]&[offset_nullifier=<int64>] (GET)

RPC:
	https://tezos.gitlab.io/008/rpc.html#get-block-id-context-sapling-sapling-state-id-get-diff
*/
func (c *Client) SaplingDiff(input SaplingDiffInput) (SingleSaplingDiff, error) {
	hash, err := c.processContextRequest(input, input.Cycle, input.Blockhash)
	if err != nil {
		return SingleSaplingDiff{}, errors.Wrap(err, "failed to get script")
	}

	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/context/sapling/%s/get_diff", c.chain, hash, input.SaplingStateID), input.contructRPCOptions()...)
	if err != nil {
		return SingleSaplingDiff{}, errors.Wrapf(err, "failed to get sapling diff for sapling '%s'", input.SaplingStateID)
	}

	var saplingDiff SingleSaplingDiff
	err = json.Unmarshal(resp, &saplingDiff)
	if err != nil {
		return SingleSaplingDiff{}, errors.Wrapf(err, "failed to get sapling diff for sapling '%s': failed to parse json", input.SaplingStateID)
	}

	return saplingDiff, nil
}

func (s *SaplingDiffInput) contructRPCOptions() []rpcOptions {
	if s.OffsetCommitment != 0 && s.OffsetNullifier != 0 {
		return []rpcOptions{
			{
				"offset_commitment",
				strconv.Itoa(s.OffsetCommitment),
			},
			{
				"offset_nullifier",
				strconv.Itoa(s.OffsetNullifier),
			},
		}
	} else if s.OffsetCommitment != 0 && s.OffsetNullifier == 0 {
		return []rpcOptions{
			{
				"offset_commitment",
				strconv.Itoa(s.OffsetCommitment),
			},
		}
	} else if s.OffsetCommitment != 0 && s.OffsetNullifier != 0 {
		return []rpcOptions{
			{
				"offset_nullifier",
				strconv.Itoa(s.OffsetNullifier),
			},
		}
	}

	return nil
}

/*
SeedInput is the input for the Seed function.

RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-context-seed
*/
type SeedInput struct {
	// The block level of which you want to make the query. If not provided Cycle is required.
	Blockhash string
	// The cycle to get the balance at. If not provided Blockhash is required.
	Cycle int
}

/*
Seed returns the seed of the cycle to which the block belongs.

Path:
	../<block_id>/context/seed (POST)

RPC:
	https://tezos.gitlab.io/008/rpc.html#post-block-id-context-seed
*/
func (c *Client) Seed(input SeedInput) (string, error) {
	hash, err := c.processContextRequest(input, input.Cycle, input.Blockhash)
	if err != nil {
		return "", errors.Wrap(err, "failed to get seed")
	}

	resp, err := c.get(fmt.Sprintf("/chains/%s/blocks/%s/context/seed", c.chain, hash))
	if err != nil {
		return "", errors.Wrap(err, "failed to get seed")
	}

	var seed string
	err = json.Unmarshal(resp, &seed)
	if err != nil {
		return "", errors.Wrapf(err, "failed to get seed: failed to parse json")
	}

	return seed, nil
}
