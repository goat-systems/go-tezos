package gotezos

import (
	"strconv"
	"sync"
)

// GetDelegationsForDelegate retrieves a list of all currently delegated contracts for a delegate.
func (gt *GoTezos) GetDelegationsForDelegate(delegatePhk string) ([]string, error) {
	rtnString := []string{}
	getDelegations := "/chains/main/blocks/head/context/delegates/" + delegatePhk + "/delegated_contracts"
	resp, err := gt.GetResponse(getDelegations, "{}")
	if err != nil {
		return rtnString, err
	}

	delegations, err := unMarshalStringArray(resp.Bytes)
	if err != nil {
		return rtnString, err
	}
	return delegations, nil
}

// GetDelegationsForDelegateByCycle retrieves a list of all currently delegated contracts for a delegate at a specific cycle.
func (gt *GoTezos) GetDelegationsForDelegateByCycle(delegatePhk string, cycle int) ([]string, error) {
	rtnString := []string{}
	snapShot, err := gt.GetSnapShot(cycle)
	if err != nil {
		return rtnString, err
	}

	hash, err := gt.GetBlockHashAtLevel(snapShot.AssociatedBlock)
	if err != nil {
		return rtnString, err
	}
	getDelegations := "/chains/main/blocks/" + hash + "/context/delegates/" + delegatePhk + "/delegated_contracts"

	resp, err := gt.GetResponse(getDelegations, "{}")
	if err != nil {
		return rtnString, err
	}

	delegations, err := unMarshalStringArray(resp.Bytes)
	if err != nil {
		return rtnString, err
	}

	return delegations, nil
}

// GetRewardsForDelegateForCycles gets the total rewards for a delegate earned
// and calculates the gross rewards earned by each delegation for multiple cycles.
// Also includes the share of each delegation.
func (gt *GoTezos) GetRewardsForDelegateForCycles(delegatePhk string, cycleStart int, cycleEnd int) (DelegationServiceRewards, error) {
	dgRewards := DelegationServiceRewards{}
	dgRewards.DelegatePhk = delegatePhk
	var cycleRewardsArray []CycleRewards

	for cycleStart <= cycleEnd {
		delegations, err := gt.GetCycleRewards(delegatePhk, cycleStart)
		if err != nil {
			return dgRewards, err
		}
		cycleRewardsArray = append(cycleRewardsArray, delegations)
		cycleStart++
	}
	dgRewards.RewardsByCycle = cycleRewardsArray
	return dgRewards, nil
}

// GetRewardsForDelegateCycle gets the total rewards for a delegate earned
// and calculates the gross rewards earned by each delegation for a single cycle.
// Also includes the share of each delegation.
func (gt *GoTezos) GetRewardsForDelegateCycle(delegatePhk string, cycle int) (DelegationServiceRewards, error) {
	dgRewards := DelegationServiceRewards{}
	dgRewards.DelegatePhk = delegatePhk

	delegations, err := gt.GetCycleRewards(delegatePhk, cycle)
	if err != nil {
		return dgRewards, err
	}
	dgRewards.RewardsByCycle = append(dgRewards.RewardsByCycle, delegations)
	return dgRewards, nil
}

// GetCycleRewards gets the total rewards for a cycle for a delegate
func (gt *GoTezos) GetCycleRewards(delegatePhk string, cycle int) (CycleRewards, error) {
	cycleRewards := CycleRewards{}
	cycleRewards.Cycle = cycle
	rewards, err := gt.GetDelegateRewardsForCycle(delegatePhk, cycle)
	if err != nil {
		return cycleRewards, err
	}

	if rewards == "" {
		rewards = "0"
	}
	cycleRewards.TotalRewards = rewards
	contractRewards, err := gt.getContractRewardsForDelegate(delegatePhk, cycleRewards.TotalRewards, cycle)
	if err != nil {
		return cycleRewards, err
	}
	cycleRewards.Delegations = contractRewards

	return cycleRewards, nil
}

// GetDelegateRewardsForCycle gets the rewards earned by a delegate for a specific cycle.
func (gt *GoTezos) GetDelegateRewardsForCycle(delegatePhk string, cycle int) (string, error) {
	rewards := FrozenBalanceRewards{}

	get := "/chains/main/blocks/head/context/raw/json/contracts/index/" + delegatePhk + "/frozen_balance/" + strconv.Itoa(cycle) + "/"
	resp, err := gt.GetResponse(get, "{}")
	if err != nil {
		return "", err
	}
	rewards, err = rewards.UnmarshalJSON(resp.Bytes)
	if err != nil {
		return rewards.Rewards, err
	}

	return rewards.Rewards, nil
}

//A private function to fill out delegation data like gross rewards and share.
func (gt *GoTezos) getContractRewardsForDelegate(delegatePhk, totalRewards string, cycle int) ([]ContractRewards, error) {

	var contractRewards []ContractRewards

	delegations, err := gt.GetDelegationsForDelegateByCycle(delegatePhk, cycle)
	if err != nil {
		return contractRewards, err
	}

	bigIntRewards, err := strconv.Atoi(totalRewards)
	if err != nil {
		return contractRewards, err
	}

	floatRewards := float64(bigIntRewards) / MUTEZ

	for _, contract := range delegations {

		contractReward := ContractRewards{}
		contractReward.DelegationPhk = contract

		share, balance, err := gt.GetShareOfContract(delegatePhk, contract, cycle)
		if err != nil {
			return contractRewards, err
		}

		contractReward.Share = share
		contractReward.Balance = balance

		bigIntGrossRewards := int((share * floatRewards) * MUTEZ)
		strGrossRewards := strconv.Itoa(bigIntGrossRewards)
		contractReward.GrossRewards = strGrossRewards

		contractRewards = append(contractRewards, contractReward)
	}

	return contractRewards, nil
}

// GetShareOfContract returns the share of a delegation for a specific cycle.
func (gt *GoTezos) GetShareOfContract(delegatePhk, delegationPhk string, cycle int) (float64, float64, error) {
	stakingBalance, err := gt.GetDelegateStakingBalance(delegatePhk, cycle)
	if err != nil {
		return 0, 0, err
	}

	delegationBalance, err := gt.GetAccountBalanceAtSnapshot(delegationPhk, cycle)
	if err != nil {
		return 0, 0, err
	}

	return delegationBalance / stakingBalance, delegationBalance, nil
}

// GetDelegate retrieves information about a delegate at the head block
func (gt *GoTezos) GetDelegate(delegatePhk string) (Delegate, error) {
	delegate := Delegate{}
	get := "/chains/main/blocks/head/context/delegates/" + delegatePhk
	resp, err := gt.GetResponse(get, "{}")
	if err != nil {
		return delegate, err
	}
	delegate, err = delegate.UnmarshalJSON(resp.Bytes)
	if err != nil {
		return delegate, err
	}

	return delegate, nil
}

// GetStakingBalanceAtCycle gets the staking balance of a delegate at a specific cycle
func (gt *GoTezos) GetStakingBalanceAtCycle(delegateAddr string, cycle int) (string, error) {
	balance := ""
	snapShot, err := gt.GetSnapShot(cycle)
	if err != nil {
		return balance, err
	}
	get := "/chains/main/blocks/" + snapShot.AssociatedHash + "/context/delegates/" + delegateAddr + "/staking_balance"
	resp, err := gt.GetResponse(get, "{}")
	if err != nil {
		return balance, err
	}
	balance, err = unmarshalString(resp.Bytes)
	if err != nil {
		return balance, err
	}

	return balance, nil
}

// GetBakingRights gets the baking rights for a specific cycle
func (gt *GoTezos) GetBakingRights(cycle int) (BakingRights, error) {
	bakingRights := BakingRights{}
	get := "/chains/main/blocks/head/helpers/baking_rights?cycle=" + strconv.Itoa(cycle) + "?max_priority=4"
	resp, err := gt.GetResponse(get, "{}")
	if err != nil {
		return bakingRights, err
	}

	bakingRights, err = bakingRights.UnmarshalJSON(resp.Bytes)
	if err != nil {
		return bakingRights, err
	}

	return bakingRights, nil
}

// GetBakingRightsForDelegate gets the baking rights for a delegate at a specific cycle with a certain priority level
func (gt *GoTezos) GetBakingRightsForDelegate(cycle int, delegatePhk string, priority int) (BakingRights, error) {
	bakingRights := BakingRights{}
	get := "/chains/main/blocks/head/helpers/baking_rights?cycle=" + strconv.Itoa(cycle) + "&max_priority=" + strconv.Itoa(priority) + "&delegate=" + delegatePhk
	resp, err := gt.GetResponse(get, "{}")
	if err != nil {
		return bakingRights, err
	}

	bakingRights, err = bakingRights.UnmarshalJSON(resp.Bytes)
	if err != nil {
		return bakingRights, err
	}

	return bakingRights, nil
}

// GetBakingRightsForDelegateForCycles gets the baking rights for a delegate at a range of cycles with a certain priority level
func (gt *GoTezos) GetBakingRightsForDelegateForCycles(cycleStart int, cycleEnd int, delegatePhk string, priority int) ([]BakingRights, error) {
	bakingRights := []BakingRights{}
	chRights := make(chan BakingRights, cycleEnd-cycleStart)
	wg := &sync.WaitGroup{}

	for cycleStart <= cycleEnd {
		wg.Add(1)
		go func() {
			get := "/chains/main/blocks/head/helpers/baking_rights?cycle=" + strconv.Itoa(cycleStart) + "&max_priority=" + strconv.Itoa(priority) + "&delegate=" + delegatePhk
			resp, _ := gt.GetResponse(get, "{}")
			// if err != nil {
			// 	return bakingRights, err
			// }

			bakingRight := new(BakingRights)
			bakingRight.UnmarshalJSON(resp.Bytes)
			// if err != nil {
			// 	return bakingRights, err
			// }
			chRights <- *bakingRight
			wg.Done()
		}()

		cycleStart++
	}
	go func() {
		wg.Wait()
		close(chRights)
	}()

	for item := range chRights {
		bakingRights = append(bakingRights, item)
	}

	return bakingRights, nil
}

// GetEndorsingRightsForDelegate gets the endorsing rights for a specific cycle
func (gt *GoTezos) GetEndorsingRightsForDelegate(cycle int, delegatePhk string) (EndorsingRights, error) {
	endorsingRights := EndorsingRights{}
	get := "/chains/main/blocks/head/helpers/endorsing_rights?cycle=" + strconv.Itoa(cycle) + "&delegate=" + delegatePhk
	resp, err := gt.GetResponse(get, "{}")
	if err != nil {
		return endorsingRights, err
	}

	endorsingRights, err = endorsingRights.UnmarshalJSON(resp.Bytes)
	if err != nil {
		return endorsingRights, err
	}

	return endorsingRights, nil
}

// GetEndorsingRightsForDelegateForCycles gets the endorsing rights for a delegate for a range of cycles
func (gt *GoTezos) GetEndorsingRightsForDelegateForCycles(cycleStart int, cycleEnd int, delegatePhk string) ([]EndorsingRights, error) {
	endorsingRights := []EndorsingRights{}
	chRights := make(chan EndorsingRights, cycleEnd-cycleStart)
	wg := &sync.WaitGroup{}

	for cycleStart <= cycleEnd {
		wg.Add(1)
		go func() {
			get := "/chains/main/blocks/head/helpers/endorsing_rights?cycle=" + strconv.Itoa(cycleStart) + "&delegate=" + delegatePhk
			resp, _ := gt.GetResponse(get, "{}")
			// if err != nil {
			// 	return endorsingRights, err
			// }
			endorsingRight := new(EndorsingRights)
			endorsingRight.UnmarshalJSON(resp.Bytes)
			// if err != nil {
			// 	return endorsingRights, err
			// }
			chRights <- *endorsingRight
			wg.Done()
		}()

		cycleStart++
	}

	go func() {
		wg.Wait()
		close(chRights)
	}()

	for item := range chRights {
		endorsingRights = append(endorsingRights, item)
	}

	return endorsingRights, nil
}

// GetEndorsingRights gets the endorsing rights for a specific cycle
func (gt *GoTezos) GetEndorsingRights(cycle int) (EndorsingRights, error) {
	endorsingRights := EndorsingRights{}
	get := "/chains/main/blocks/head/helpers/endorsing_rights?cycle=" + strconv.Itoa(cycle)
	resp, err := gt.GetResponse(get, "{}")
	if err != nil {
		return endorsingRights, err
	}

	endorsingRights, err = endorsingRights.UnmarshalJSON(resp.Bytes)
	if err != nil {
		return endorsingRights, err
	}

	return endorsingRights, nil
}

// GetAllDelegatesByHash gets a list of all tz1 addresses at a certain hash
func (gt *GoTezos) GetAllDelegatesByHash(hash string) ([]string, error) {
	delList := []string{}
	get := "/chains/main/blocks/" + hash + "/context/delegates?active"
	resp, err := gt.GetResponse(get, "{}")
	if err != nil {
		return delList, err
	}
	delList, err = unMarshalStringArray(resp.Bytes)
	if err != nil {
		return delList, err
	}
	return delList, nil
}

// GetAllDelegates a list of all tz1 addresses at the head block
func (gt *GoTezos) GetAllDelegates() ([]string, error) {
	delList := []string{}
	get := "/chains/main/blocks/head/context/delegates?active"
	resp, err := gt.GetResponse(get, "{}")
	if err != nil {
		return delList, err
	}
	delList, err = unMarshalStringArray(resp.Bytes)
	if err != nil {
		return delList, err
	}
	return delList, nil
}
