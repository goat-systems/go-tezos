package goTezos

import (
	"strconv"
	"strings"
)

//Takes a cycle number and returns a helper structure describing a snap shot on the tezos network.
func GetSnapShot(cycle int) (SnapShot, error) {
	var snapShotQuery SnapShotQuery
	var snap SnapShot
	var get string

	currentCycle, err := GetCurrentCycle()
	if err != nil {
		return snap, err
	}

	snap.Cycle = cycle
	strCycle := strconv.Itoa(cycle)

	if cycle < currentCycle {
		block, err := GetBlockAtLevel(cycle * 4096)
		if err != nil {
			return snap, err
		}
		get = "/chains/main/blocks/" + block.Hash + "/context/raw/json/cycle/" + strCycle

	} else {
		get = "/chains/main/blocks/head/context/raw/json/cycle/" + strCycle
	}

	byts, err := TezosRPCGet(get)
	if err != nil {
		logger.Println("Could not get snap shot: " + err.Error())
		return snap, err
	}
	snapShotQuery, err = unMarshelSnapShotQuery(byts)
	if err != nil {
		logger.Println("Could not get snap shot: " + err.Error())
		return snap, err
	}

	snap.Number = snapShotQuery.RollSnapShot
	snap.AssociatedBlock = ((cycle - 7) * 4096) + (snapShotQuery.RollSnapShot+1)*256
	snap.AssociatedHash, _ = GetBlockHashAtLevel(snap.AssociatedBlock)

	return snap, nil
}

//Gets a list of all known snapshots to the network
func GetAllCurrentSnapShots() ([]SnapShot, error) {
	var snapShotArray []SnapShot
	currentCycle, err := GetCurrentCycle()
	if err != nil {
		return snapShotArray, err
	}
	for i := 7; i <= currentCycle; i++ {
		snapShot, err := GetSnapShot(i)
		if err != nil {
			return snapShotArray, err
		}
		snapShotArray = append(snapShotArray, snapShot)
	}

	return snapShotArray, nil
}

//Returns the head block from the Tezos RPC.
func GetChainHead() (Block, error) {
	var block Block
	byts, err := TezosRPCGet("/chains/main/blocks/head")
	if err != nil {
		logger.Println("Could not get /chains/main/blocks/head: " + err.Error())
		return block, err
	}
	block, err = unMarshelBlock(byts)
	if err != nil {
		logger.Println("Could not get block head: " + err.Error())
	}

	return block, nil
}

//Returns the level, and the hash, of the head block.
func GetBlockLevelHead() (int, string, error) {
	block, err := GetChainHead()
	if err != nil {
		return block.Header.Level, block.Hash, err
	}
	return block.Header.Level, block.Hash, err
}

//Returns the hash of a block at a specific level.
func GetBlockHashAtLevel(level int) (string, error) {
	head, headHash, err := GetBlockLevelHead()
	if err != nil {
		return "", err
	}

	diffStr := strconv.Itoa(head - level)
	getBlockByLevel := "/chains/main/blocks/" + headHash + "~" + diffStr

	s, err := TezosRPCGet(getBlockByLevel)
	if err != nil {
		return "", err
	}

	block, err := unMarshelBlock(s)
	if err != nil {
		return "", err
	}

	return block.Hash, nil
}

//Returns a Block at a specific level
func GetBlockAtLevel(level int) (Block, error) {
	var block Block
	head, headHash, err := GetBlockLevelHead()
	if err != nil {
		return block, err
	}

	diffStr := strconv.Itoa(head - level)
	getBlockByLevel := "/chains/main/blocks/" + headHash + "~" + diffStr

	s, err := TezosRPCGet(getBlockByLevel)
	if err != nil {
		return block, err
	}

	block, err = unMarshelBlock(s)
	if err != nil {
		return block, err
	}

	return block, nil
}

//Returns a Block by the identifier hash.
func GetBlockByHash(hash string) (Block, error) {
	var block Block

	getBlockByLevel := "/chains/main/blocks/" + hash

	s, err := TezosRPCGet(getBlockByLevel)
	if err != nil {
		return block, err
	}
	block, err = unMarshelBlock(s)
	if err != nil {
		return block, err
	}
	return block, nil
}

//Gets the balance of a public key hash at a specific snapshot for a cycle.
func GetAccountBalanceAtSnapshot(tezosAddr string, cycle int) (float64, error) {
	snapShot, err := GetSnapShot(cycle)
	if err != nil {
		return 0, err
	}

	hash, err := GetBlockHashAtLevel(snapShot.AssociatedBlock)
	if err != nil {
		return 0, err
	}

	balanceCmdStr := "/chains/main/blocks/" + hash + "/context/contracts/" + tezosAddr + "/balance"
	s, err := TezosRPCGet(balanceCmdStr)
	if err != nil {
		return 0, err
	}

	strBalance, err := unMarshelString(s)
	if err != nil {
		return 0, err
	}

	floatBalance, err := strconv.ParseFloat(strBalance, 64)
	if err != nil {
		return 0, err
	}

	return floatBalance / 1000000, nil
}

//Gets the balance of a public key hash at a specific snapshot for a cycle.
func GetAccountBalance(tezosAddr string) (float64, error) {

	balanceCmdStr := "/chains/main/blocks/head/context/contracts/" + tezosAddr + "/balance"
	s, err := TezosRPCGet(balanceCmdStr)
	if err != nil {
		return 0, err
	}

	strBalance, err := unMarshelString(s)
	if err != nil {
		return 0, err
	}

	floatBalance, err := strconv.ParseFloat(strBalance, 64)
	if err != nil {
		return 0, err
	}

	return floatBalance / 1000000, nil
}

//Gets the staking balance for a delegate at a specific snapshot for a cycle.
func GetDelegateStakingBalance(delegateAddr string, cycle int) (float64, error) {
	var snapShot SnapShot
	var err error
	var hash string

	snapShot, err = GetSnapShot(cycle)
	if err != nil {
		return 0, err
	}

	hash, err = GetBlockHashAtLevel(snapShot.AssociatedBlock)
	if err != nil {
		return 0, err
	}

	rpcCall := "/chains/main/blocks/" + hash + "/context/delegates/" + delegateAddr + "/staking_balance"

	s, err := TezosRPCGet(rpcCall)
	if err != nil {
		return 0, err
	}

	strBalance, err := unMarshelString(s)
	if err != nil {
		return 0, err
	}

	floatBalance, err := strconv.ParseFloat(strBalance, 64) //TODO error checking
	if err != nil {
		return 0, err
	}

	return floatBalance / 1000000, nil
}

//Gets the current cycle of the chain
func GetCurrentCycle() (int, error) {
	block, err := GetChainHead()
	if err != nil {
		return 0, err
	}
	var cycle int
	cycle = block.Header.Level / 4096

	return cycle, nil
}

//Get the balance of an address at a specific hash
func GetAccountBalanceAtBlock(tezosAddr string, hash string) (int, error) {
	var balance string
	balanceCmdStr := "/chains/main/blocks/" + hash + "/context/contracts/" + tezosAddr + "/balance"

	byts, err := TezosRPCGet(balanceCmdStr)
	if err != nil {
		return 0, err
	}
	balance, err = unMarshelString(byts)
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
func GetChainId() (string, error) {
	chainIdCmd := "/chains/main/chain_id"
	bytes, err := TezosRPCGet(chainIdCmd)
	if err != nil {
		return "", err
	}

	chainId, err := unMarshelString(bytes)
	if err != nil {
		return "", err
	}

	return chainId, nil
}
