package goTezos

import (
	"fmt"
	"strconv"
	"strings"
)

//Takes a cycle number and returns a helper structure describing a snap shot on the tezos network.
func (this *GoTezos) GetSnapShot(cycle int) (SnapShot, error) {
	
	var snapShotQuery SnapShotQuery
	var snap SnapShot
	var get string
	
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
	
	// Check cache
	if chainHeadCache.Hash == "" {
		
		var block Block
		
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
	
		chainHeadCache = block
	}
			
	return chainHeadCache, nil
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
	head, headHash, err := this.GetBlockLevelHead()
	if err != nil {
		return block, err
	}

	diffStr := strconv.Itoa(head - level)
	getBlockByLevel := "/chains/main/blocks/" + headHash + "~" + diffStr

	resp, err := this.GetResponse(getBlockByLevel, "{}")
	if err != nil {
		return block, err
	}

	block, err = unMarshalBlock(resp.Bytes)
	if err != nil {
		return block, err
	}

	return block, nil
}

//Returns a Block by the identifier hash.
func (this *GoTezos) GetBlockByHash(hash string) (Block, error) {
	var block Block

	getBlockByLevel := "/chains/main/blocks/" + hash

	resp, err := this.GetResponse(getBlockByLevel, "{}")
	if err != nil {
		return block, err
	}
	block, err = unMarshalBlock(resp.Bytes)
	if err != nil {
		return block, err
	}
	return block, nil
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
