//Package goTezos exposes the Tezos RPC API in goLang.
package goTezos

import (
	"log"
	"strconv"
)

//Takes a cycle number and returns a helper structure describing a snap shot on the tezos network.
func GetSnapShot(cycle int) (SnapShot, error) {
	var snapShotQuery SnapShotQuery
	var snap SnapShot
	snap.Cycle = cycle
	strCycle := strconv.Itoa(cycle)

	get := "/chains/main/blocks/head/context/raw/json/cycle/" + strCycle

	byts, err := TezosRPCGet(get)
	if err != nil {
		log.Println("Could not get snap shot: " + err.Error())
		return snap, err
	}
	snapShotQuery, err = unMarshelSnapShotQuery(byts)
	if err != nil {
		log.Println("Could not get snap shot: " + err.Error())
		return snap, err
	}

	snap.Number = snapShotQuery.RollSnapShot
	snap.AssociatedBlock = ((cycle - 7) * 4096) + (snapShotQuery.RollSnapShot+1)*256
	snap.AssociatedHash, _ = GetBlockHashAtLevel(snap.AssociatedBlock)

	return snap, nil
}

//Returns the head block from the Tezos RPC.
func GetChainHead() (Block, error) {
	var block Block
	byts, err := TezosRPCGet("/chains/main/blocks/head")
	if err != nil {
		log.Println("Could not get /chains/main/blocks/head: " + err.Error())
		return block, err
	}
	block, err = unMarshelBlock(byts)
	if err != nil {
		log.Println("Could not get block head: " + err.Error())
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
