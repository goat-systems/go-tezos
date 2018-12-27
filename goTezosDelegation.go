package goTezos

import "strconv"

//A function that retrieves a list of all currently delegated contracts for a delegate.
func (this *GoTezos) GetDelegationsForDelegate(delegatePhk string) ([]string, error) {
	var rtnString []string
	getDelegations := "/chains/main/blocks/head/context/delegates/" + delegatePhk + "/delegated_contracts"
	resp, err := this.GetResponse(getDelegations, "{}")
	if err != nil {
		return rtnString, err
	}

	delegations, err := unMarshalStringArray(resp.Bytes)
	if err != nil {
		return rtnString, err
	}
	return delegations, nil
}

//A function that retrieves a list of all currently delegated contracts for a delegate at a specific cycle.
func (this *GoTezos) GetDelegationsForDelegateByCycle(delegatePhk string, cycle int) ([]string, error) {
	var rtnString []string
	snapShot, err := this.GetSnapShot(cycle)
	if err != nil {
		return rtnString, err
	}

	hash, err := this.GetBlockHashAtLevel(snapShot.AssociatedBlock)
	if err != nil {
		return rtnString, err
	}
	getDelegations := "/chains/main/blocks/" + hash + "/context/delegates/" + delegatePhk + "/delegated_contracts"

	resp, err := this.GetResponse(getDelegations, "{}")
	if err != nil {
		return rtnString, err
	}

	delegations, err := unMarshalStringArray(resp.Bytes)
	if err != nil {
		return rtnString, err
	}

	return delegations, nil
}

//Gets the total rewards for a delegate earned and calculates the gross rewards earned by each delegation for multiple cycles. Also includes the share of each delegation.
func (this *GoTezos) GetRewardsForDelegateForCycles(delegatePhk string, cycleStart int, cycleEnd int) (DelegationServiceRewards, error) {
	dgRewards := DelegationServiceRewards{}
	dgRewards.DelegatePhk = delegatePhk
	var cycleRewardsArray []CycleRewards

	for cycleStart <= cycleEnd {
		delegations, err := this.getCycleRewards(delegatePhk, cycleStart)
		if err != nil {
			return dgRewards, err
		}
		cycleRewardsArray = append(cycleRewardsArray, delegations)
		cycleStart++
	}
	dgRewards.RewardsByCycle = cycleRewardsArray
	return dgRewards, nil
}

//Gets the total rewards for a delegate earned and calculates the gross rewards earned by each delegation for a single cycle. Also includes the share of each delegation.
func (this *GoTezos) GetRewardsForDelegateCycle(delegatePhk string, cycle int) (DelegationServiceRewards, error) {
	dgRewards := DelegationServiceRewards{}
	dgRewards.DelegatePhk = delegatePhk

	delegations, err := this.getCycleRewards(delegatePhk, cycle)
	if err != nil {
		return dgRewards, err
	}
	dgRewards.RewardsByCycle = append(dgRewards.RewardsByCycle, delegations)
	return dgRewards, nil
}

//Get total rewards for a cycle for a delegate
func (this *GoTezos) getCycleRewards(delegatePhk string, cycle int) (CycleRewards, error) {
	cycleRewards := CycleRewards{}
	cycleRewards.Cycle = cycle
	rewards, err := this.GetDelegateRewardsForCycle(delegatePhk, cycle)
	if err != nil {
		return cycleRewards, err
	}

	if rewards == "" {
		rewards = "0"
	}
	cycleRewards.TotalRewards = rewards
	contractRewards, err := this.getContractRewardsForDelegate(delegatePhk, cycleRewards.TotalRewards, cycle)
	if err != nil {
		return cycleRewards, err
	}
	cycleRewards.Delegations = contractRewards

	return cycleRewards, nil
}

//A function that gets the rewards earned by a delegate for a specific cycle.
func (this *GoTezos) GetDelegateRewardsForCycle(delegatePhk string, cycle int) (string, error) {
	get := "/chains/main/blocks/head/context/raw/json/contracts/index/" + delegatePhk + "/frozen_balance/" + strconv.Itoa(cycle) + "/"
	resp, err := this.GetResponse(get, "{}")
	if err != nil {
		return "", err
	}
	rewards, err := unMarshalFrozenBalanceRewards(resp.Bytes)
	if err != nil {
		return rewards.Rewards, err
	}

	return rewards.Rewards, nil
}

//A private function to fill out delegation data like gross rewards and share.
func (this *GoTezos) getContractRewardsForDelegate(delegatePhk, totalRewards string, cycle int) ([]ContractRewards, error) {
	var contractRewards []ContractRewards
	delegations, err := this.GetDelegationsForDelegateByCycle(delegatePhk, cycle)
	if err != nil {
		return contractRewards, err
	}
	for _, contract := range delegations {
		contractReward := ContractRewards{}
		contractReward.DelegationPhk = contract
		bigIntRewards, err := strconv.Atoi(totalRewards)
		if err != nil {
			return contractRewards, err
		}
		floatRewards := float64(bigIntRewards) / MUTEZ
		share, err := this.GetShareOfContract(delegatePhk, contract, cycle)
		if err != nil {
			return contractRewards, err
		}
		contractReward.Share = share
		bigIntGrossRewards := int((share * floatRewards) * MUTEZ)
		strGrossRewards := strconv.Itoa(bigIntGrossRewards)
		contractReward.GrossRewards = strGrossRewards

		contractRewards = append(contractRewards, contractReward)
	}

	return contractRewards, nil
}

//Returns the share of a delegation for a specific cycle.
func (this *GoTezos) GetShareOfContract(delegatePhk, delegationPhk string, cycle int) (float64, error) {
	stakingBalance, err := this.GetDelegateStakingBalance(delegatePhk, cycle)
	if err != nil {
		return 0, err
	}

	delegationBalance, err := this.GetAccountBalanceAtSnapshot(delegationPhk, cycle)
	if err != nil {
		return 0, err
	}

	return delegationBalance / stakingBalance, nil
}

//RPC command to retrieve information about a delegate at the head block
func (this *GoTezos) GetDelegate(delegatePhk string) (Delegate, error) {
	var delegate Delegate
	get := "/chains/main/blocks/head/context/delegates/" + delegatePhk
	resp, err := this.GetResponse(get, "{}")
	if err != nil {
		return delegate, err
	}
	delegate, err = unMarshalDelegate(resp.Bytes)
	if err != nil {
		return delegate, err
	}

	return delegate, err
}

//RPC command to get the staking balance of a delegate at a specific cycle
func (this *GoTezos) GetStakingBalanceAtCycle(cycle int, delegateAddr string) (string, error) {
	var balance string
	snapShot, err := this.GetSnapShot(cycle)
	if err != nil {
		return balance, err
	}
	get := "/chains/main/blocks/" + snapShot.AssociatedHash + "/context/delegates/" + delegateAddr + "/staking_balance"
	resp, err := this.GetResponse(get, "{}")
	if err != nil {
		return balance, err
	}
	balance, err = unMarshalString(resp.Bytes)
	if err != nil {
		return balance, err
	}

	return balance, nil
}

//Gets the baking rights for a specific cycle
func (this *GoTezos) GetBakingRights(cycle int) (Baking_Rights, error) {
	var bakingRights Baking_Rights
	get := "/chains/main/blocks/head/helpers/baking_rights?cycle=" + strconv.Itoa(cycle) + "?max_priority=4"
	resp, err := this.GetResponse(get, "{}")
	if err != nil {
		return bakingRights, err
	}

	bakingRights, err = unMarshalBakingRights(resp.Bytes)
	if err != nil {
		return bakingRights, err
	}

	return bakingRights, nil
}

func (this *GoTezos) GetBakingRightsForDelegate(cycle int, delegatePhk string, priority int) (Baking_Rights, error) {
	var bakingRights Baking_Rights
	get := "/chains/main/blocks/head/helpers/baking_rights?cycle=" + strconv.Itoa(cycle) + "&max_priority=" + strconv.Itoa(priority) + "&delegate=" + delegatePhk
	resp, err := this.GetResponse(get, "{}")
	if err != nil {
		return bakingRights, err
	}

	bakingRights, err = unMarshalBakingRights(resp.Bytes)
	if err != nil {
		return bakingRights, err
	}

	return bakingRights, nil
}

func (this *GoTezos) GetBakingRightsForDelegateForCycles(cycleStart int, cycleEnd int, delegatePhk string, priority int) ([]Baking_Rights, error) {
	var bakingRights []Baking_Rights
	for cycleStart <= cycleEnd {
		get := "/chains/main/blocks/head/helpers/baking_rights?cycle=" + strconv.Itoa(cycleStart) + "&max_priority=" + strconv.Itoa(priority) + "&delegate=" + delegatePhk
		resp, err := this.GetResponse(get, "{}")
		if err != nil {
			return bakingRights, err
		}

		bakingRight, err := unMarshalBakingRights(resp.Bytes)
		if err != nil {
			return bakingRights, err
		}
		bakingRights = append(bakingRights, bakingRight)
		cycleStart++
	}

	return bakingRights, nil
}

//Gets the endorsing rights for a specific cycle
func (this *GoTezos) GetEndorsingRightsForDelegate(cycle int, delegatePhk string) (Endorsing_Rights, error) {
	var endorsingRights Endorsing_Rights
	get := "/chains/main/blocks/head/helpers/endorsing_rights?cycle=" + strconv.Itoa(cycle) + "&delegate=" + delegatePhk
	resp, err := this.GetResponse(get, "{}")
	if err != nil {
		return endorsingRights, err
	}

	endorsingRights, err = unMarshalEndorsingRights(resp.Bytes)
	if err != nil {
		return endorsingRights, err
	}

	return endorsingRights, nil
}

func (this *GoTezos) GetEndorsingRightsForDelegateForCycles(cycleStart int, cycleEnd int, delegatePhk string) ([]Endorsing_Rights, error) {
	var endorsingRights []Endorsing_Rights
	for cycleStart <= cycleEnd {
		get := "/chains/main/blocks/head/helpers/endorsing_rights?cycle=" + strconv.Itoa(cycleStart) + "&delegate=" + delegatePhk
		resp, err := this.GetResponse(get, "{}")
		if err != nil {
			return endorsingRights, err
		}

		endorsingRight, err := unMarshalEndorsingRights(resp.Bytes)
		if err != nil {
			return endorsingRights, err
		}
		endorsingRights = append(endorsingRights, endorsingRight)
		cycleStart++
	}

	return endorsingRights, nil
}

//Gets the endorsing rights for a specific cycle
func (this *GoTezos) GetEndorsingRights(cycle int) (Endorsing_Rights, error) {
	var endorsingRights Endorsing_Rights
	get := "/chains/main/blocks/head/helpers/endorsing_rights?cycle=" + strconv.Itoa(cycle)
	resp, err := this.GetResponse(get, "{}")
	if err != nil {
		return endorsingRights, err
	}

	endorsingRights, err = unMarshalEndorsingRights(resp.Bytes)
	if err != nil {
		return endorsingRights, err
	}

	return endorsingRights, nil
}

//Retrieves a list of all tz1 addresses at a certain hash
func (this *GoTezos) GetAllDelegatesByHash(hash string) ([]string, error) {
	var delList []string
	get := "/chains/main/blocks" + hash + "/context/delegates?active"
	resp, err := this.GetResponse(get, "{}")
	if err != nil {
		return delList, err
	}
	delList, err = unMarshalStringArray(resp.Bytes)
	if err != nil {
		return delList, err
	}
	return delList, nil
}

//Retrieves a list of all tz1 addresses
func (this *GoTezos) GetAllDelegates() ([]string, error) {
	var delList []string
	get := "/chains/main/blocks/head/context/delegates?active"
	resp, err := this.GetResponse(get, "{}")
	if err != nil {
		return delList, err
	}
	delList, err = unMarshalStringArray(resp.Bytes)
	if err != nil {
		return delList, err
	}
	return delList, nil
}
