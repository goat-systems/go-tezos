package goTezos

/*
Author: DefinitelyNotAGoat/MagicAglet
Version: 0.0.1
Description: The Tezos API written in GO, for easy development.
License: MIT
*/

import (
	"log"
	"strconv"
	//"fmt"
)

/*
Description: Gets the snapshot number for a certain cycle and returns the block level
Param cycle (int): Takes a cycle number as an integer to query for that cycles snapshot
Returns struct SnapShot: A SnapShot Structure defined above.
*/
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

func GetBlockLevelHead() (int, string, error) {
	block, err := GetChainHead()
	if err != nil {
		return block.Header.Level, block.Hash, err
	}
	return block.Header.Level, block.Hash, err
}

/*
Description: Takes a block level, and returns the hash for that specific level
Param level (int): An integer representation of the block level to query
Returns (string): A string representation of the hash for the block level queried.
*/
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

	floatBalance, err := strconv.ParseFloat(strBalance, 64) //TODO error checking
	if err != nil {
		return 0, err
	}

	//fmt.Println(returnBalance)

	return floatBalance / 1000000, nil
}

/*
Description: Will get the staking balance of a delegate
Param delegateAddr (string): Takes a string representation of the address querying
Returns (float64): Returns a float64 representation of the balance for the account
*/
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
