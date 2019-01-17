package goTezos

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

//Takes a cycle number and returns a helper structure describing a snap shot on the tezos network.
func (this *GoTezos) GetSnapShot(cycle int) (SnapShot, error) {
	
	var snapShotQuery SnapShotQuery
	var snap SnapShot
	var get string
	

	// Check cache first
	if cachedSS, exists := this.cache.Get(strconv.Itoa(cycle)); exists {
		if this.debug {
			this.logger.Printf("DEBUG: GetSnapShot %d (Cached)\n", cycle)
		}
		return cachedSS.(SnapShot), nil
	}
	
	if this.debug {
		this.logger.Printf("GetSnapShot %d\n", cycle)
	}

	currentCycle, err := this.GetCurrentCycle()
	if err != nil {
		return snap, err
	}
	
	// Sanity
	if cycle > currentCycle + this.Constants.PreservedCycles - 1 {
		return snap, fmt.Errorf("Unable to fetch snapshot for cycle %d; Cycle does not exist.", cycle)
	}
	
	snap.Cycle = cycle
	strCycle := strconv.Itoa(cycle)

	if cycle < currentCycle {
		block, err := this.GetBlockAtLevel(cycle*this.Constants.BlocksPerCycle + 1)
		if err != nil {
			return snap, err
		}
		get = "/chains/main/blocks/" + block.Hash + "/context/raw/json/cycle/" + strCycle

	} else {
		get = "/chains/main/blocks/head/context/raw/json/cycle/" + strCycle
	}

	resp, err := this.GetResponse(get, "{}")
	if err != nil {
		this.logger.Println("Could not get snap shot: " + err.Error())
		return snap, err
	}
	snapShotQuery, err = unMarshalSnapShotQuery(resp.Bytes)
	if err != nil {
		this.logger.Println("Could not get snap shot: " + err.Error())
		return snap, err
	}

	snap.Number = snapShotQuery.RollSnapShot

	snap.AssociatedBlock = ((cycle - this.Constants.PreservedCycles) * this.Constants.BlocksPerCycle) + (snapShotQuery.RollSnapShot+1)*256
	if snap.AssociatedBlock < 1 {
		snap.AssociatedBlock = 1
	}
	snap.AssociatedHash, _ = this.GetBlockHashAtLevel(snap.AssociatedBlock)

	// Cache for future
	// Can be a longer cache since old snapshots don't change
	this.cache.Set(strconv.Itoa(cycle), snap, 10 * time.Minute)
	
	return snap, nil
}

//Gets a list of all known snapshots to the network
func (this *GoTezos) GetAllCurrentSnapShots() ([]SnapShot, error) {
	var snapShotArray []SnapShot
	currentCycle, err := this.GetCurrentCycle()
	if err != nil {
		return snapShotArray, err
	}
	for i := 7; i <= currentCycle; i++ {
		snapShot, err := this.GetSnapShot(i)
		if err != nil {
			return snapShotArray, err
		}
		snapShotArray = append(snapShotArray, snapShot)
	}

	return snapShotArray, nil
}

//Returns the head block from the Tezos RPC.
func (this *GoTezos) GetChainHead() (Block, error) {

	var block Block
	
	// Check cache
	if cachedBlock, exists := this.cache.Get("head"); exists {
		if this.debug {
			this.logger.Println("DEBUG: GetChainHead() (Cached)")
		}
		return cachedBlock.(Block), nil
	}
	
	if this.debug {
		this.logger.Println("DEBUG: GetChainHead()")
	}

	resp, err := this.GetResponse("/chains/main/blocks/head", "{}")
	if err != nil {
		this.logger.Println("Could not get /chains/main/blocks/head: " + err.Error())
		return block, err
	}

	block, err = unMarshalBlock(resp.Bytes)
	if err != nil {
		this.logger.Println("Could not get block head: " + err.Error())
		return block, err
	}
	
	// Cache. Not for too long since the head can change every minute
	this.cache.Set("head", block, 10 * time.Second)
	return block, nil
}

func (this *GoTezos) GetNetworkConstants() (NetworkConstants, error) {
	networkConstants := NetworkConstants{}
	resp, err := this.GetResponse("/chains/main/blocks/head/context/constants", "{}")
	if err != nil {
		this.logger.Println("Could not get /chains/main/blocks/head/context/constants: " + err.Error())
		return networkConstants, err
	}
	networkConstants, err = unMarshalNetworkConstants(resp.Bytes)
	if err != nil {
		this.logger.Println("Could not get network constants: " + err.Error())
		return networkConstants, err
	}

	return networkConstants, nil
}

func (this *GoTezos) GetNetworkVersions() ([]NetworkVersion, error) {
	
	networkVersions := make([]NetworkVersion, 0)
	
	resp, err := this.GetResponse("/network/versions", "{}")
	if err != nil {
		this.logger.Println("Could not get /network/versions: " + err.Error())
		return networkVersions, err
	}
	
	nvs, err := unMarshalNetworkVersion(resp.Bytes)
	if err != nil {
		this.logger.Println("Could not get network version: " + err.Error())
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

func (this *GoTezos) GetBranchProtocol() (string, error) {
	block, err := this.GetChainHead()
	if err != nil {
		return "", err
	}
	return block.Protocol, nil
}

func (this *GoTezos) GetBranchHash() (string, error) {
	block, err := this.GetChainHead()
	if err != nil {
		return "", err
	}
	return block.Hash, nil
}

//Returns the level, and the hash, of the head block.
func (this *GoTezos) GetBlockLevelHead() (int, string, error) {
	block, err := this.GetChainHead()
	if err != nil {
		return block.Header.Level, block.Hash, err
	}
	return block.Header.Level, block.Hash, err
}

//Returns the hash of a block at a specific level.
func (this *GoTezos) GetBlockHashAtLevel(level int) (string, error) {
	block, err := this.GetBlockAtLevel(level)
	if err != nil {
		return "", err
	}

	return block.Hash, nil
}

//Returns a Block at a specific level
func (this *GoTezos) GetBlockAtLevel(level int) (Block, error) {
	
	var block Block
	
	// Check cache for block at this level
	if cachedBlock, exists := this.cache.Get(strconv.Itoa(level)); exists {
		if this.debug {
			this.logger.Printf("DEBUG: GetBlockAtLevel %d (Cached)\n", level)
		}
		return cachedBlock.(Block), nil
	}
	
	// Get current head block
	headLevel, headHash, err := this.GetBlockLevelHead()
	if err != nil {
		return block, err
	}
	
	if this.debug {
		this.logger.Printf("DEBUG: GetBlockAtLevel %d\n", level)
	}
	
	diffStr := strconv.Itoa(headLevel - level)
	getBlockByLevel := "/chains/main/blocks/" + headHash + "~" + diffStr
	fmt.Println(getBlockByLevel)

	resp, err := this.GetResponse(getBlockByLevel, "{}")
	if err != nil {
		return block, err
	}

	block, err = unMarshalBlock(resp.Bytes)
	if err != nil {
		return block, err
	}
	
	// Cache for future
	// Can be a longer cache since old blocks don't change
	this.cache.Set(strconv.Itoa(level), block, 10 * time.Minute)
	
	return block, nil
}

//Returns a Block by the identifier hash.
func (this *GoTezos) GetBlockByHash(hash string) (Block, error) {
	
	var block Block
	hashPrefix := hash[:15]
	
	// Check cache
	if cachedBlock, exists := this.cache.Get(hashPrefix); exists {
		if this.debug {
			this.logger.Println("DEBUG: GetBlockByHash (Cached)")
		}
		return cachedBlock.(Block), nil
	}
	
	if this.debug {
		this.logger.Println("DEBUG: GetBlockByHash")
	}
	
	getBlockByLevel := "/chains/main/blocks/" + hash

	resp, err := this.GetResponse(getBlockByLevel, "{}")
	if err != nil {
		return block, err
	}
	
	block, err = unMarshalBlock(resp.Bytes)
	if err != nil {
		return block, err
	}
	
	// Cache for future
	// Can be a longer cache since old blocks don't change
	this.cache.Set(hashPrefix, block, 10 * time.Minute)
	
	return block, nil
}

//Returns list of operations in block of head
func (this *GoTezos) GetBlockOperationHashesHead() (OperationHashes, error) {
	
	var operations OperationHashes
	
	// Get head hash
	blockHash, err := this.GetBranchHash()
	if err != nil {
		this.logger.Printf("Could not get block hash at head: %s\n", err)
		return operations, err
	}
	
	// Pass hash to helper
	return this.GetBlockOperationHashes(blockHash)
}

//Returns list of operations in block at specific level
func (this *GoTezos) GetBlockOperationHashesAtLevel(level int) (OperationHashes, error) {
	
	var operations OperationHashes
	
	blockHash, err := this.GetBlockHashAtLevel(level)
	if err != nil {
		this.logger.Printf("Could not get block hash at level %d: %s\n", level, err)
		return operations, err
	}
	
	// Pass hash to helper
	return this.GetBlockOperationHashes(blockHash)
}


func (this *GoTezos) GetBlockOperationHashes(blockHash string) (OperationHashes, error) {
	
	var operations OperationHashes
	
	resp, err := this.GetResponse("/chains/main/blocks/" + blockHash + "/operation_hashes", "{}")
	if err != nil {
		this.logger.Println("Could not get block operation_hashes: " + err.Error())
		return operations, err
	}
	
	operations, err = unMarshalOperationHashes(resp.Bytes)
	if err != nil {
		this.logger.Println("Could not decode operation hashes: " + err.Error())
		return operations, err
	}
	
	return operations, nil
}

//Gets the balance of a public key hash at a specific snapshot for a cycle.
func (this *GoTezos) GetAccountBalanceAtSnapshot(tezosAddr string, cycle int) (float64, error) {
	snapShot, err := this.GetSnapShot(cycle)
	if err != nil {
		return 0, err
	}

	hash, err := this.GetBlockHashAtLevel(snapShot.AssociatedBlock)
	if err != nil {
		return 0, err
	}

	balanceCmdStr := "/chains/main/blocks/" + hash + "/context/contracts/" + tezosAddr + "/balance"
	resp, err := this.GetResponse(balanceCmdStr, "{}")
	if err != nil {
		return 0, err
	}

	strBalance, err := unMarshalString(resp.Bytes)
	if err != nil {
		return 0, err
	}

	floatBalance, err := strconv.ParseFloat(strBalance, 64)
	if err != nil {
		return 0, err
	}

	return floatBalance / MUTEZ, nil
}

//Gets the balance of a public key hash at a specific snapshot for a cycle.
func (this *GoTezos) GetAccountBalance(tezosAddr string) (float64, error) {

	balanceCmdStr := "/chains/main/blocks/head/context/contracts/" + tezosAddr + "/balance"
	resp, err := this.GetResponse(balanceCmdStr, "{}")
	if err != nil {
		return 0, err
	}

	strBalance, err := unMarshalString(resp.Bytes)
	if err != nil {
		return 0, err
	}

	floatBalance, err := strconv.ParseFloat(strBalance, 64)
	if err != nil {
		return 0, err
	}

	return floatBalance / MUTEZ, nil
}

//Gets the staking balance for a delegate at a specific snapshot for a cycle.
func (this *GoTezos) GetDelegateStakingBalance(delegateAddr string, cycle int) (float64, error) {
	var snapShot SnapShot
	var err error
	var hash string

	snapShot, err = this.GetSnapShot(cycle)
	if err != nil {
		return 0, err
	}

	hash, err = this.GetBlockHashAtLevel(snapShot.AssociatedBlock)
	if err != nil {
		return 0, err
	}

	rpcCall := "/chains/main/blocks/" + hash + "/context/delegates/" + delegateAddr + "/staking_balance"

	resp, err := this.GetResponse(rpcCall, "{}")
	if err != nil {
		return 0, err
	}

	strBalance, err := unMarshalString(resp.Bytes)
	if err != nil {
		return 0, err
	}

	floatBalance, err := strconv.ParseFloat(strBalance, 64) //TODO error checking
	if err != nil {
		return 0, err
	}

	return floatBalance / MUTEZ, nil
}

//Gets the current cycle of the chain
func (this *GoTezos) GetCurrentCycle() (int, error) {
	block, err := this.GetChainHead()
	if err != nil {
		return 0, err
	}
	
	var cycle int
	cycle = block.Metadata.Level.Cycle
	
	return cycle, nil
}

//Get the balance of an address at a specific hash
func (this *GoTezos) GetAccountBalanceAtBlock(tezosAddr string, hash string) (int, error) {
	var balance string
	balanceCmdStr := "/chains/main/blocks/" + hash + "/context/contracts/" + tezosAddr + "/balance"

	resp, err := this.GetResponse(balanceCmdStr, "{}")
	if err != nil {
		return 0, err
	}
	balance, err = unMarshalString(resp.Bytes)
	if err != nil {
		return 0, err
	}

	var returnBalance int
	if strings.Contains(balance, "No service found at this URL") {
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

//Gets the ID of the chain with the most fitness
func (this *GoTezos) GetChainId() (string, error) {
	chainIdCmd := "/chains/main/chain_id"
	resp, err := this.GetResponse(chainIdCmd, "{}")
	if err != nil {
		return "", err
	}

	chainId, err := unMarshalString(resp.Bytes)
	if err != nil {
		return "", err
	}

	return chainId, nil
}
