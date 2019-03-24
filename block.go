package gotezos

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// GetSnapShot takes a cycle number and returns a helper structure describing a snap shot on the tezos network.
func (gt *GoTezos) GetSnapShot(cycle int) (SnapShot, error) {

	var snapShotQuery SnapShotQuery
	var snap SnapShot
	var get string

	// Check cache first
	if cachedSS, exists := gt.cache.Get(strconv.Itoa(cycle)); exists {
		if gt.debug {
			gt.logger.Printf("DEBUG: GetSnapShot %d (Cached)\n", cycle)
		}
		return cachedSS.(SnapShot), nil
	}

	if gt.debug {
		gt.logger.Printf("GetSnapShot %d\n", cycle)
	}

	currentCycle, err := gt.GetCurrentCycle()
	if err != nil {
		return snap, err
	}

	// Sanity
	if cycle > currentCycle+gt.Constants.PreservedCycles-1 {
		return snap, fmt.Errorf("unable to fetch snapshot for cycle %d; Cycle does not exist", cycle)
	}

	snap.Cycle = cycle
	strCycle := strconv.Itoa(cycle)

	if cycle < currentCycle {
		block, err := gt.GetBlockAtLevel(cycle*gt.Constants.BlocksPerCycle + 1)
		if err != nil {
			return snap, err
		}
		get = "/chains/main/blocks/" + block.Hash + "/context/raw/json/cycle/" + strCycle

	} else {
		get = "/chains/main/blocks/head/context/raw/json/cycle/" + strCycle
	}

	resp, err := gt.GetResponse(get, "{}")
	if err != nil {
		gt.logger.Println("Could not get snap shot: " + err.Error())
		return snap, err
	}

	snapShotQuery, err = snapShotQuery.UnmarshalJSON(resp.Bytes)
	if err != nil {
		gt.logger.Println("Could not get snap shot: " + err.Error())
		return snap, err
	}

	snap.Number = snapShotQuery.RollSnapShot

	snap.AssociatedBlock = ((cycle - gt.Constants.PreservedCycles - 2) * gt.Constants.BlocksPerCycle) + (snapShotQuery.RollSnapShot+1)*256
	if snap.AssociatedBlock < 1 {
		snap.AssociatedBlock = 1
	}
	snap.AssociatedHash, _ = gt.GetBlockHashAtLevel(snap.AssociatedBlock)

	// Cache for future
	// Can be a longer cache since old snapshots don't change
	gt.cache.Set(strconv.Itoa(cycle), snap, 10*time.Minute)

	return snap, nil
}

// GetAllCurrentSnapShots gets a list of all known snapshots to the network
func (gt *GoTezos) GetAllCurrentSnapShots() ([]SnapShot, error) {
	var snapShotArray []SnapShot
	currentCycle, err := gt.GetCurrentCycle()
	if err != nil {
		return snapShotArray, err
	}
	for i := 7; i <= currentCycle; i++ {
		snapShot, err := gt.GetSnapShot(i)
		if err != nil {
			return snapShotArray, err
		}
		snapShotArray = append(snapShotArray, snapShot)
	}

	return snapShotArray, nil
}

// GetChainHead returns the head block from the Tezos RPC.
func (gt *GoTezos) GetChainHead() (Block, error) {

	var block Block

	// Check cache
	if cachedBlock, exists := gt.cache.Get("head"); exists {
		if gt.debug {
			gt.logger.Println("DEBUG: GetChainHead() (Cached)")
		}
		return cachedBlock.(Block), nil
	}

	if gt.debug {
		gt.logger.Println("DEBUG: GetChainHead()")
	}

	resp, err := gt.GetResponse("/chains/main/blocks/head", "{}")
	if err != nil {
		gt.logger.Println("Could not get /chains/main/blocks/head: " + err.Error())
		return block, err
	}

	block, err = block.UnmarshalJSON(resp.Bytes)
	if err != nil {
		gt.logger.Println("Could not get block head: " + err.Error())
		return block, err
	}

	// Cache. Not for too long since the head can change every minute
	gt.cache.Set("head", block, 10*time.Second)
	return block, nil
}

// GetNetworkConstants gets the network constants for the Tezos network the client is using.
func (gt *GoTezos) GetNetworkConstants() (NetworkConstants, error) {
	networkConstants := NetworkConstants{}
	resp, err := gt.GetResponse("/chains/main/blocks/head/context/constants", "{}")
	if err != nil {
		gt.logger.Println("Could not get /chains/main/blocks/head/context/constants: " + err.Error())
		return networkConstants, err
	}
	networkConstants, err = networkConstants.UnmarshalJSON(resp.Bytes)
	if err != nil {
		gt.logger.Println("Could not get network constants: " + err.Error())
		return networkConstants, err
	}

	return networkConstants, nil
}

// GetNetworkVersions gets the network versions of Tezos network the client is using.
func (gt *GoTezos) GetNetworkVersions() ([]NetworkVersion, error) {

	networkVersions := make([]NetworkVersion, 0)

	resp, err := gt.GetResponse("/network/versions", "{}")
	if err != nil {
		gt.logger.Println("Could not get /network/versions: " + err.Error())
		return networkVersions, err
	}

	var nvs NetworkVersions
	nvs, err = nvs.UnmarshalJSON(resp.Bytes)
	if err != nil {
		gt.logger.Println("Could not get network version: " + err.Error())
		return networkVersions, err
	}

	// Extract just the network name and append to returning slice
	// 'range' operates on a copy of the struct so cannot update-in-place
	for _, v := range nvs {

		parts := strings.Split(v.Name, "_")
		if len(parts) == 3 {
			v.Network = parts[1]
		}

		networkVersions = append(networkVersions, v)
	}

	return networkVersions, nil
}

// GetBranchProtocol returns the current Tezos protocol hash
func (gt *GoTezos) GetBranchProtocol() (string, error) {
	block, err := gt.GetChainHead()
	if err != nil {
		return "", err
	}
	return block.Protocol, nil
}

// GetBranchHash returns the current branches hash
func (gt *GoTezos) GetBranchHash() (string, error) {
	block, err := gt.GetChainHead()
	if err != nil {
		return "", err
	}
	return block.Hash, nil
}

// GetBlockLevelHead returns the level, and the hash, of the head block.
func (gt *GoTezos) GetBlockLevelHead() (int, string, error) {
	block, err := gt.GetChainHead()
	if err != nil {
		return block.Header.Level, block.Hash, err
	}
	return block.Header.Level, block.Hash, err
}

// GetBlockHashAtLevel returns the hash of a block at a specific level.
func (gt *GoTezos) GetBlockHashAtLevel(level int) (string, error) {
	block, err := gt.GetBlockAtLevel(level)
	if err != nil {
		return "", err
	}

	return block.Hash, nil
}

// GetBlockAtLevel returns a Block at a specific level
func (gt *GoTezos) GetBlockAtLevel(level int) (Block, error) {

	var block Block

	// Check cache for block at level
	if cachedBlock, exists := gt.cache.Get(strconv.Itoa(level)); exists {
		if gt.debug {
			gt.logger.Printf("DEBUG: GetBlockAtLevel %d (Cached)\n", level)
		}
		return cachedBlock.(Block), nil
	}

	// Get current head block
	headLevel, headHash, err := gt.GetBlockLevelHead()
	if err != nil {
		return block, err
	}

	if gt.debug {
		gt.logger.Printf("DEBUG: GetBlockAtLevel %d\n", level)
	}

	diffStr := strconv.Itoa(headLevel - level)
	getBlockByLevel := "/chains/main/blocks/" + headHash + "~" + diffStr

	if gt.debug {
		gt.logger.Printf("DEBUG: - %s\n", getBlockByLevel)
	}

	resp, err := gt.GetResponse(getBlockByLevel, "{}")
	if err != nil {
		return block, err
	}

	block, err = block.UnmarshalJSON(resp.Bytes)
	if err != nil {
		return block, err
	}

	// Cache for future
	// Can be a longer cache since old blocks don't change
	gt.cache.Set(strconv.Itoa(level), block, 10*time.Minute)

	return block, nil
}

// GetBlockByHash returns a Block by the identifier hash.
func (gt *GoTezos) GetBlockByHash(hash string) (Block, error) {

	var block Block
	hashPrefix := hash[:15]

	// Check cache
	if cachedBlock, exists := gt.cache.Get(hashPrefix); exists {
		if gt.debug {
			gt.logger.Println("DEBUG: GetBlockByHash (Cached)")
		}
		return cachedBlock.(Block), nil
	}

	if gt.debug {
		gt.logger.Println("DEBUG: GetBlockByHash")
	}

	getBlockByLevel := "/chains/main/blocks/" + hash

	resp, err := gt.GetResponse(getBlockByLevel, "{}")
	if err != nil {
		return block, err
	}

	block, err = block.UnmarshalJSON(resp.Bytes)
	if err != nil {
		return block, err
	}

	// Cache for future
	// Can be a longer cache since old blocks don't change
	gt.cache.Set(hashPrefix, block, 10*time.Minute)

	return block, nil
}

// GetBlockOperationHashesHead returns list of operations in the head block
func (gt *GoTezos) GetBlockOperationHashesHead() (OperationHashes, error) {

	var operations OperationHashes

	// Get head hash
	blockHash, err := gt.GetBranchHash()
	if err != nil {
		gt.logger.Printf("Could not get block hash at head: %s\n", err)
		return operations, err
	}

	// Pass hash to helper
	return gt.GetBlockOperationHashes(blockHash)
}

// GetBlockOperationHashesAtLevel returns list of operations in block at specific level
func (gt *GoTezos) GetBlockOperationHashesAtLevel(level int) (OperationHashes, error) {

	var operations OperationHashes

	blockHash, err := gt.GetBlockHashAtLevel(level)
	if err != nil {
		gt.logger.Printf("Could not get block hash at level %d: %s\n", level, err)
		return operations, err
	}

	// Pass hash to helper
	return gt.GetBlockOperationHashes(blockHash)
}

// GetBlockOperationHashes returns all the operation hashes at specific block hash
func (gt *GoTezos) GetBlockOperationHashes(blockHash string) (OperationHashes, error) {

	var operations OperationHashes

	resp, err := gt.GetResponse("/chains/main/blocks/"+blockHash+"/operation_hashes", "{}")
	if err != nil {
		gt.logger.Println("Could not get block operation_hashes: " + err.Error())
		return operations, err
	}

	operations, err = operations.UnmarshalJSON(resp.Bytes)
	if err != nil {
		gt.logger.Println("Could not decode operation hashes: " + err.Error())
		return operations, err
	}

	return operations, nil
}

// GetAccountBalanceAtSnapshot gets the balance of a public key hash at a specific snapshot for a cycle.
func (gt *GoTezos) GetAccountBalanceAtSnapshot(tezosAddr string, cycle int) (float64, error) {
	snapShot, err := gt.GetSnapShot(cycle)
	if err != nil {
		return 0, err
	}

	hash, err := gt.GetBlockHashAtLevel(snapShot.AssociatedBlock)
	if err != nil {
		return 0, err
	}

	balanceCmdStr := "/chains/main/blocks/" + hash + "/context/contracts/" + tezosAddr + "/balance"
	resp, err := gt.GetResponse(balanceCmdStr, "{}")
	if err != nil {
		return 0, err
	}

	strBalance, err := unmarshalString(resp.Bytes)
	if err != nil {
		return 0, err
	}

	floatBalance, err := strconv.ParseFloat(strBalance, 64)
	if err != nil {
		return 0, err
	}

	return floatBalance / MUTEZ, nil
}

// GetAccountBalance gets the balance of a public key hash at a specific snapshot for a cycle.
func (gt *GoTezos) GetAccountBalance(tezosAddr string) (float64, error) {

	balanceCmdStr := "/chains/main/blocks/head/context/contracts/" + tezosAddr + "/balance"
	resp, err := gt.GetResponse(balanceCmdStr, "{}")
	if err != nil {
		return 0, err
	}

	strBalance, err := unmarshalString(resp.Bytes)
	if err != nil {
		return 0, err
	}

	floatBalance, err := strconv.ParseFloat(strBalance, 64)
	if err != nil {
		return 0, err
	}

	return floatBalance / MUTEZ, nil
}

// GetDelegateStakingBalance gets the staking balance for a delegate at a specific snapshot for a cycle.
func (gt *GoTezos) GetDelegateStakingBalance(delegateAddr string, cycle int) (float64, error) {
	var snapShot SnapShot
	var err error
	var hash string

	snapShot, err = gt.GetSnapShot(cycle)
	if err != nil {
		return 0, err
	}

	hash, err = gt.GetBlockHashAtLevel(snapShot.AssociatedBlock)
	if err != nil {
		return 0, err
	}

	rpcCall := "/chains/main/blocks/" + hash + "/context/delegates/" + delegateAddr + "/staking_balance"

	resp, err := gt.GetResponse(rpcCall, "{}")
	if err != nil {
		return 0, err
	}

	strBalance, err := unmarshalString(resp.Bytes)
	if err != nil {
		return 0, err
	}

	floatBalance, err := strconv.ParseFloat(strBalance, 64) //TODO error checking
	if err != nil {
		return 0, err
	}

	return floatBalance / MUTEZ, nil
}

// GetCurrentCycle gets the current cycle of the chain
func (gt *GoTezos) GetCurrentCycle() (int, error) {
	block, err := gt.GetChainHead()
	if err != nil {
		return 0, err
	}

	var cycle int
	cycle = block.Metadata.Level.Cycle

	return cycle, nil
}

// GetAccountBalanceAtBlock get the balance of an address at a specific hash
func (gt *GoTezos) GetAccountBalanceAtBlock(tezosAddr string, hash string) (int, error) {
	var balance string
	balanceCmdStr := "/chains/main/blocks/" + hash + "/context/contracts/" + tezosAddr + "/balance"

	resp, err := gt.GetResponse(balanceCmdStr, "{}")
	if err != nil {
		return 0, err
	}
	balance, err = unmarshalString(resp.Bytes)
	if err != nil {
		return 0, err
	}

	var returnBalance int
	if strings.Contains(balance, "No service found at gt URL") {
		returnBalance = 0
	}

	if len(balance) < 1 {
		returnBalance = 0
	} else {
		floatBalance, _ := strconv.Atoi(balance) //TODO error checking
		returnBalance = int(floatBalance)
	}

	return returnBalance, nil
}

// GetChainID gets the id of the chain with the most fitness
func (gt *GoTezos) GetChainID() (string, error) {
	get := "/chains/main/chain_id"
	resp, err := gt.GetResponse(get, "{}")
	if err != nil {
		return "", err
	}

	chainID, err := unmarshalString(resp.Bytes)
	if err != nil {
		return "", err
	}

	return chainID, nil
}
