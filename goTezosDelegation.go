package goTezos

import "strconv"

//A function that retrieves a list of all currently delegated contracts for a delegate.
func GetDelegationsForDelegate(delegatePhk string) ([]string, error) {
	var rtnString []string
	getDelegations := "/chains/main/blocks/head/context/delegates/" + delegatePhk + "/delegated_contracts"
	s, err := TezosRPCGet(getDelegations)
	if err != nil {
		return rtnString, err
	}

	delegations, err := unMarshelStringArray(s)
	if err != nil {
		return rtnString, err
	}
	return delegations, nil
}

//A function that retrieves a list of all currently delegated contracts for a delegate at a specific cycle.
func GetDelegationsForDelegateByCycle(delegatePhk string, cycle int) ([]string, error) {
	var rtnString []string
	snapShot, err := GetSnapShot(cycle)
	if err != nil {
		return rtnString, err
	}

	hash, err := GetBlockHashAtLevel(snapShot.AssociatedBlock)
	if err != nil {
		return rtnString, err
	}
	getDelegations := "/chains/main/blocks/" + hash + "/context/delegates/" + delegatePhk + "/delegated_contracts"

	s, err := TezosRPCGet(getDelegations)
	if err != nil {
		return rtnString, err
	}

	delegations, err := unMarshelStringArray(s)
	if err != nil {
		return rtnString, err
	}

	return delegations, nil
}

//Gets the total rewards for a delegate earned and calculates the gross rewards earned by each delegation for multiple cycles. Also includes the share of each delegation.
func GetRewardsForDelegateForCycles(delegatePhk string, cycleStart int, cycleEnd int) (DelegationServiceRewards, error) {
	dgRewards := DelegationServiceRewards{}
	dgRewards.delegatePhk = delegatePhk
	var cycleRewardsArray []CycleRewards

	for cycleStart <= cycleEnd {
		delegations, err := getCycleRewards(delegatePhk, cycleStart)
		if err != nil {
			return dgRewards, err
		}
		cycleRewardsArray = append(cycleRewardsArray, delegations)
		cycleStart++
	}
	return dgRewards, nil
}

//Gets the total rewards for a delegate earned and calculates the gross rewards earned by each delegation for a single cycle. Also includes the share of each delegation.
func GetRewardsForDelegateCycle(delegatePhk string, cycle int) (DelegationServiceRewards, error) {
	dgRewards := DelegationServiceRewards{}
	dgRewards.delegatePhk = delegatePhk
	var cycleRewardsArray []CycleRewards

	delegations, err := getCycleRewards(delegatePhk, cycle)
	if err != nil {
		return dgRewards, err
	}
	cycleRewardsArray = append(cycleRewardsArray, delegations)
	return dgRewards, nil
}

//Get total rewards for a cycle for a delegate
func getCycleRewards(delegatePhk string, cycle int) (CycleRewards, error) {
	cycleRewards := CycleRewards{}
	cycleRewards.Cycle = cycle
	rewards, err := GetDelegateRewardsForCycle(delegatePhk, cycle)
	if err != nil {
		return cycleRewards, err
	}
	cycleRewards.TotalRewards = rewards
	contractRewards, err := getContractRewardsForDelegate(delegatePhk, cycleRewards.TotalRewards, cycle)
	if err != nil {
		return cycleRewards, err
	}
	cycleRewards.Delegations = contractRewards

	return cycleRewards, nil
}

//A function that gets the rewards earned by a delegate for a specific cycle.
func GetDelegateRewardsForCycle(delegatePhk string, cycle int) (string, error) {
	get := "/chains/main/blocks/head/context/raw/json/contracts/index/" + delegatePhk + "/frozen_balance/" + strconv.Itoa(cycle) + "/"
	s, err := TezosRPCGet(get)
	if err != nil {
		return "", err
	}
	rewards, err := unMarshelFrozenBalanceRewards(s)
	if err != nil {
		return rewards.Rewards, err
	}

	return rewards.Rewards, nil
}

//A private function to fill out delegation data like gross rewards and share.
func getContractRewardsForDelegate(delegatePhk, totalRewards string, cycle int) ([]ContractRewards, error) {
	var contractRewards []ContractRewards
	delegations, err := GetDelegationsForDelegateByCycle(delegatePhk, cycle)
	if err != nil {
		return contractRewards, err
	}
	for _, contract := range delegations {
		contractReward := ContractRewards{}
		contractReward.delegationPhk = contract
		bigIntRewards, err := strconv.Atoi(totalRewards)
		if err != nil {
			return contractRewards, err
		}
		floatRewards := float64(bigIntRewards) / 1000000
		share, err := GetShareOfContract(delegatePhk, contract, cycle)
		if err != nil {
			return contractRewards, err
		}
		contractReward.Share = share
		bigIntGrossRewards := int((share * floatRewards) * 1000000)
		strGrossRewards := strconv.Itoa(bigIntGrossRewards)
		contractReward.GrossRewards = strGrossRewards

		contractRewards = append(contractRewards, contractReward)
	}

	return contractRewards, nil
}

//Returns the share of a delegation for a specific cycle.
func GetShareOfContract(delegatePhk, delegationPhk string, cycle int) (float64, error) {
	stakingBalance, err := GetDelegateStakingBalance(delegatePhk, cycle)
	if err != nil {
		return 0, err
	}

	delegationBalance, err := GetAccountBalanceAtSnapshot(delegationPhk, cycle)
	if err != nil {
		return 0, err
	}

	return delegationBalance / stakingBalance, nil
}

//RPC command to retrieve information about a delegate at the head block
func GetDelegate(delegatePhk string) (Delegate, error) {
	var delegate Delegate
	get := "/chains/main/blocks/head/context/delegates/" + delegatePhk
	byts, err := TezosRPCGet(get)
	if err != nil {
		return delegate, err
	}
	delegate, err = unMarshelDelegate(byts)
	if err != nil {
		return delegate, err
	}

	return delegate, err
}

//RPC command to get the staking balance of a delegate at a specific cycle
func GetStakingBalanceAtCycle(cycle int, delegateAddr string) (string, error) {
	var balance string
	snapShot, err := GetSnapShot(cycle)
	if err != nil {
		return balance, err
	}
	get := "/chains/main/blocks/" + snapShot.AssociatedHash + "/context/delegates/" + delegateAddr + "/staking_balance"
	byts, err := TezosRPCGet(get)
	if err != nil {
		return balance, err
	}
	balance, err = unMarshelString(byts)
	if err != nil {
		return balance, err
	}

	return balance, nil
}

//Gets the baking rights for a specific cycle
func GetBakingRights(cycle int) (Baking_Rights, error) {
	var BakingRights Baking_Rights
	get := "/chains/main/blocks/head/helpers/baking_rights?cycle=" + strconv.Itoa(cycle) + "?max_priority=4"
	byts, err := TezosRPCGet(get)
	if err != nil {
		return BakingRights, err
	}

	BakingRights, err = unMarshelBakingRights(byts)
	if err != nil {
		return BakingRights, err
	}

	return BakingRights, nil
}

//Gets the endorsing rights for a specific cycle
func GetEndorsingRights(cycle int) (Endorsing_Rights, error) {
	var endorsingRights Endorsing_Rights
	get := "/chains/main/blocks/head/helpers/endorsing_rights?cycle=" + strconv.Itoa(cycle)
	byts, err := TezosRPCGet(get)
	if err != nil {
		return endorsingRights, err
	}

	endorsingRights, err = unMarshelEndorsingRights(byts)
	if err != nil {
		return endorsingRights, err
	}

	return endorsingRights, nil
}

//Retrieves a list of all tz1 addresses at a certain hash
func GetAllDelegatesByHash(hash string) ([]string, error) {
	var delList []string
	get := "/chains/main/" + hash + "/context/delegates"
	bytes, err := TezosRPCGet(get)
	if err != nil {
		return delList, err
	}
	delList, err = unMarshelStringArray(bytes)
	if err != nil {
		return delList, err
	}
	return delList, nil
}

//Retrieves a list of all tz1 addresses
func GetAllDelegates() ([]string, error) {
	var delList []string
	get := "/chains/main/head/context/delegates"
	bytes, err := TezosRPCGet(get)
	if err != nil {
		return delList, err
	}
	delList, err = unMarshelStringArray(bytes)
	if err != nil {
		return delList, err
	}
	return delList, nil
}
