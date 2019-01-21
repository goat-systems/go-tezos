package goTezos

import (
	"strconv"
	"sync"
)

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
	chCycleRewards := make(chan CycleRewards, cycleEnd-cycleStart)
	wg := &sync.WaitGroup{}

	for cycleStart <= cycleEnd {
		wg.Add(1)
		go func() {
			delegations, _ := this.getCycleRewards(delegatePhk, cycleStart)
			// if err != nil {
			// 	return dgRewards, err
			// }
			chCycleRewards <- delegations

			wg.Done()
		}()

		cycleStart++
	}
	go func() {
		wg.Wait()
		close(chCycleRewards)
	}()

	for item := range chCycleRewards {
		cycleRewardsArray = append(cycleRewardsArray, item)
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
	rewards := new(FrozenBalanceRewards)

	get := "/chains/main/blocks/head/context/raw/json/contracts/index/" + delegatePhk + "/frozen_balance/" + strconv.Itoa(cycle) + "/"
	resp, err := this.GetResponse(get, "{}")
	if err != nil {
		return "", err
	}
	err = rewards.UnmarshalJSON(resp.Bytes)
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

	bigIntRewards, err := strconv.Atoi(totalRewards)
	if err != nil {
		return contractRewards, err
	}

	floatRewards := float64(bigIntRewards) / MUTEZ
	len := len(delegations)
	chRewards := make(chan ContractRewards, len)
	wg := &sync.WaitGroup{}

	for index, contract := range delegations {
		wg.Add(1)
		go func(delegation string, ch <-chan ContractRewards, index int, len int) {
			contractReward := ContractRewards{}
			contractReward.DelegationPhk = delegation

			share, balance, _ := this.GetShareOfContract(delegatePhk, delegation, cycle)
			// if err != nil {
			// 	return contractRewards, err
			// }

			contractReward.Share = share
			contractReward.Balance = balance

			bigIntGrossRewards := int((share * floatRewards) * MUTEZ)
			strGrossRewards := strconv.Itoa(bigIntGrossRewards)
			contractReward.GrossRewards = strGrossRewards

			chRewards <- contractReward

			wg.Done()
		}(contract, chRewards, index, len)
	}
	go func() {
		wg.Wait()
		close(chRewards)
	}()

	for item := range chRewards {
		contractRewards = append(contractRewards, item)
	}

	return contractRewards, nil
}

//Returns the share of a delegation for a specific cycle.
func (this *GoTezos) GetShareOfContract(delegatePhk, delegationPhk string, cycle int) (float64, float64, error) {
	stakingBalance, err := this.GetDelegateStakingBalance(delegatePhk, cycle)
	if err != nil {
		return 0, 0, err
	}

	delegationBalance, err := this.GetAccountBalanceAtSnapshot(delegationPhk, cycle)
	if err != nil {
		return 0, 0, err
	}

	return delegationBalance / stakingBalance, delegationBalance, nil
}

//RPC command to retrieve information about a delegate at the head block
func (this *GoTezos) GetDelegate(delegatePhk string) (Delegate, error) {
	var delegate *Delegate
	get := "/chains/main/blocks/head/context/delegates/" + delegatePhk
	resp, err := this.GetResponse(get, "{}")
	if err != nil {
		return *delegate, err
	}
	err = delegate.UnmarshalJSON(resp.Bytes)
	if err != nil {
		return *delegate, err
	}

	return *delegate, err
}

//RPC command to get the staking balance of a delegate at a specific cycle
func (this *GoTezos) GetStakingBalanceAtCycle(delegateAddr string, cycle int) (string, error) {
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
	var bakingRights *Baking_Rights
	get := "/chains/main/blocks/head/helpers/baking_rights?cycle=" + strconv.Itoa(cycle) + "?max_priority=4"
	resp, err := this.GetResponse(get, "{}")
	if err != nil {
		return *bakingRights, err
	}

	err = bakingRights.UnmarshalJSON(resp.Bytes)
	if err != nil {
		return *bakingRights, err
	}

	return *bakingRights, nil
}

func (this *GoTezos) GetBakingRightsForDelegate(cycle int, delegatePhk string, priority int) (Baking_Rights, error) {
	var bakingRights *Baking_Rights
	get := "/chains/main/blocks/head/helpers/baking_rights?cycle=" + strconv.Itoa(cycle) + "&max_priority=" + strconv.Itoa(priority) + "&delegate=" + delegatePhk
	resp, err := this.GetResponse(get, "{}")
	if err != nil {
		return *bakingRights, err
	}

	err = bakingRights.UnmarshalJSON(resp.Bytes)
	if err != nil {
		return *bakingRights, err
	}

	return *bakingRights, nil
}

func (this *GoTezos) GetBakingRightsForDelegateForCycles(cycleStart int, cycleEnd int, delegatePhk string, priority int) ([]Baking_Rights, error) {
	var bakingRights []Baking_Rights
	chRights := make(chan Baking_Rights, cycleEnd-cycleStart)
	wg := &sync.WaitGroup{}

	for cycleStart <= cycleEnd {
		wg.Add(1)
		go func() {
			get := "/chains/main/blocks/head/helpers/baking_rights?cycle=" + strconv.Itoa(cycleStart) + "&max_priority=" + strconv.Itoa(priority) + "&delegate=" + delegatePhk
			resp, _ := this.GetResponse(get, "{}")
			// if err != nil {
			// 	return bakingRights, err
			// }

			bakingRight := new(Baking_Rights)
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

//Gets the endorsing rights for a specific cycle
func (this *GoTezos) GetEndorsingRightsForDelegate(cycle int, delegatePhk string) (Endorsing_Rights, error) {
	var endorsingRights *Endorsing_Rights
	get := "/chains/main/blocks/head/helpers/endorsing_rights?cycle=" + strconv.Itoa(cycle) + "&delegate=" + delegatePhk
	resp, err := this.GetResponse(get, "{}")
	if err != nil {
		return *endorsingRights, err
	}

	err = endorsingRights.UnmarshalJSON(resp.Bytes)
	if err != nil {
		return *endorsingRights, err
	}

	return *endorsingRights, nil
}

func (this *GoTezos) GetEndorsingRightsForDelegateForCycles(cycleStart int, cycleEnd int, delegatePhk string) ([]Endorsing_Rights, error) {
	var endorsingRights []Endorsing_Rights
	chRights := make(chan Endorsing_Rights, cycleEnd-cycleStart)
	wg := &sync.WaitGroup{}

	for cycleStart <= cycleEnd {
		wg.Add(1)
		go func() {
			get := "/chains/main/blocks/head/helpers/endorsing_rights?cycle=" + strconv.Itoa(cycleStart) + "&delegate=" + delegatePhk
			resp, _ := this.GetResponse(get, "{}")
			// if err != nil {
			// 	return endorsingRights, err
			// }
			var endorsingRight *Endorsing_Rights
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

//Gets the endorsing rights for a specific cycle
func (this *GoTezos) GetEndorsingRights(cycle int) (Endorsing_Rights, error) {
	var endorsingRights *Endorsing_Rights
	get := "/chains/main/blocks/head/helpers/endorsing_rights?cycle=" + strconv.Itoa(cycle)
	resp, err := this.GetResponse(get, "{}")
	if err != nil {
		return *endorsingRights, err
	}

	err = endorsingRights.UnmarshalJSON(resp.Bytes)
	if err != nil {
		return *endorsingRights, err
	}

	return *endorsingRights, nil
}

//Retrieves a list of all tz1 addresses at a certain hash
func (this *GoTezos) GetAllDelegatesByHash(hash string) ([]string, error) {
	var delList []string
	get := "/chains/main/blocks/" + hash + "/context/delegates?active"
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
