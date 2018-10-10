package goTezos

import (
	"errors"
	"math"
	"strconv"
)

//A function used to calculate gross/net rewards, share, and fees for a delegate in a range of cycles. The function takes in a list of delegated contracts, the cycles to calculate,
//the and rate of the delegate.
func CalculateAllContractsForCycles(delegatedContracts []DelegatedContract, cycleStart int, cycleEnd int, rate float64, spillage bool, delegateAddr string) ([]DelegatedContract, error) {
	var err error

	for cycleStart <= cycleEnd {
		delegatedContracts, err = CalculateAllContractsForCycle(delegatedContracts, cycleStart, rate, spillage, delegateAddr)
		if err != nil {
			return delegatedContracts, errors.New("Could not calculate all commitments for cycles " + strconv.Itoa(cycleStart) + "-" + strconv.Itoa(cycleEnd) + ":CalculateAllCommitmentsForCycle(delegatedContracts []DelegatedContract, cycle int, rate float64) failed: " + err.Error())
		}
		cycleStart = cycleStart + 1
	}
	return delegatedContracts, nil
}

//A function used to Calculate gross/net rewards, share, and fees for a delegate. The function takes in a list of delegated contracts, the cycle to calculate,
//the and rate of the delegate.
func CalculateAllContractsForCycle(delegatedContracts []DelegatedContract, cycle int, rate float64, spillage bool, delegateAddr string) ([]DelegatedContract, error) {
	var err error
	var balance float64
	delegationsForCycle, _ := GetDelegatedContractsForCycle(cycle, delegateAddr)

	for index, delegation := range delegatedContracts {
		balance, err = GetAccountBalanceAtSnapshot(delegation.Address, cycle)
		if err != nil {
			return delegatedContracts, errors.New("Could not calculate all commitments for cycle " + strconv.Itoa(cycle) + ":GetAccountBalanceAtSnapshot(tezosAddr string, cycle int) failed: " + err.Error())
		}
		if isDelegationInGroup(delegatedContracts[index].Address, delegationsForCycle, delegatedContracts[index].Delegate) {
			delegatedContracts[index].Contracts = append(delegatedContracts[index].Contracts, Contract{Cycle: cycle, Amount: balance})
		} else {
			delegatedContracts[index].Contracts = append(delegatedContracts[index].Contracts, Contract{Cycle: cycle, Amount: 0})
		}
	}

	delegatedContracts, err = CalculatePercentageSharesForCycle(delegatedContracts, cycle, rate, spillage, delegateAddr)
	if err != nil {
		return delegatedContracts, errors.New("func CalculateAllContractsForCycle(delegatedContracts []DelegatedContract, cycle int, rate float64, spillage bool, delegateAddr string) failed: " + err.Error())
	}
	return delegatedContracts, nil
}

//A helper function to see if a public key hash is included in an array of public key hashes.
func isDelegationInGroup(phk string, group []string, delegate bool) bool {
	if delegate {
		return true
	}
	for _, address := range group {
		if address == phk {
			return true
		}
	}
	return false
}

//A function that loops through an array of delegations, and calulates the share of each delegation for a cycle.
func CalculatePercentageSharesForCycle(delegatedContracts []DelegatedContract, cycle int, rate float64, spillage bool, delegateAddr string) ([]DelegatedContract, error) {
	var stakingBalance float64
	//var balance float64
	var err error

	spillAlert := false

	stakingBalance, err = GetDelegateStakingBalance(delegateAddr, cycle)
	if err != nil {
		return delegatedContracts, errors.New("func CalculateRollSpillage(delegatedContracts []DelegatedContract, delegateAddr string) failed: " + err.Error())
	}

	mod := math.Mod(stakingBalance, 10000)
	sum := stakingBalance - mod
	balanceCheck := stakingBalance - mod

	for index, delegation := range delegatedContracts {
		counter := 0
		for i, _ := range delegation.Contracts {
			if delegatedContracts[index].Contracts[i].Cycle == cycle {
				break
			}
			counter = counter + 1
		}
		balanceCheck = balanceCheck - delegatedContracts[index].Contracts[counter].Amount
		//fmt.Println(stakingBalance)
		if spillAlert {
			delegatedContracts[index].Contracts[counter].SharePercentage = 0
			delegatedContracts[index].Contracts[counter].RollInclusion = 0
		} else if balanceCheck < 0 && spillage {
			spillAlert = true
			delegatedContracts[index].Contracts[counter].SharePercentage = (delegatedContracts[index].Contracts[counter].Amount + stakingBalance) / sum
			delegatedContracts[index].Contracts[counter].RollInclusion = delegatedContracts[index].Contracts[counter].Amount + stakingBalance
		} else {
			delegatedContracts[index].Contracts[counter].SharePercentage = delegatedContracts[index].Contracts[counter].Amount / stakingBalance
			delegatedContracts[index].Contracts[counter].RollInclusion = delegatedContracts[index].Contracts[counter].Amount
		}
		delegatedContracts[index].Contracts[counter] = CalculatePayoutForContract(delegatedContracts[index].Contracts[counter], rate, delegatedContracts[index].Delegate, delegateAddr)
		delegatedContracts[index].Fee = delegatedContracts[index].Fee + delegatedContracts[index].Contracts[counter].Fee
	}

	return delegatedContracts, nil
}

//A function that retrieves a list of all delegated contracts for a delegate in a specific cycle.
func GetDelegatedContractsForCycle(cycle int, delegateAddr string) ([]string, error) {
	var rtnString []string
	snapShot, err := GetSnapShot(cycle)
	// fmt.Println(snapShot)
	if err != nil {
		return rtnString, errors.New("Could not get delegated contracts for cycle " + strconv.Itoa(cycle) + ": GetSnapShot(cycle int) failed: " + err.Error())
	}
	hash, err := GetBlockHashAtLevel(snapShot.AssociatedBlock)
	if err != nil {
		return rtnString, errors.New("Could not get delegated contracts for cycle " + strconv.Itoa(cycle) + ": GetBlockLevelHash(level int) failed: " + err.Error())
	}
	// fmt.Println(hash)
	getDelegatedContracts := "/chains/main/blocks/" + hash + "/context/delegates/" + delegateAddr + "/delegated_contracts"

	s, err := TezosRPCGet(getDelegatedContracts)
	if err != nil {
		return rtnString, errors.New("Could not get delegated contracts for cycle " + strconv.Itoa(cycle) + ": TezosRPCGet(arg string) failed: " + err.Error())
	}

	DelegatedContracts, err := unMarshelStringArray(s)
	if err != nil {
		return rtnString, errors.New("Could not get delegated contracts for cycle " + strconv.Itoa(cycle) + ": You have no contracts.")
	}

	return DelegatedContracts, nil
}

//A function that retrieves a list of all currently delegated contracts.
func GetAllDelegatedContracts(delegateAddr string) ([]string, error) {
	var rtnString []string
	delegatedContractsCmd := "/chains/main/blocks/head/context/delegates/" + delegateAddr + "/delegated_contracts"
	s, err := TezosRPCGet(delegatedContractsCmd)
	if err != nil {
		return rtnString, errors.New("Could not get delegated contracts: TezosRPCGet(arg string) failed: " + err.Error())
	}

	DelegatedContracts, err := unMarshelStringArray(s)
	if err != nil {
		return rtnString, err
	}
	//fmt.Println(rtnString)
	return DelegatedContracts, nil
}

//A function that retrieves a list of all delegated contracts for a range of cycles.
func GetDelegatedContractsBetweenContracts(cycleStart int, cycleEnd int, delegateAddr string) ([]string, error) {

	contracts, err := GetDelegatedContractsForCycle(cycleStart, delegateAddr)
	if err != nil {
		return contracts, err
	}

	cycleStart++

	for ; cycleStart <= cycleEnd; cycleStart++ {
		tmpContracts, err := GetDelegatedContractsForCycle(cycleStart, delegateAddr)
		if err != nil {
			return contracts, err
		}
		for _, tmpContract := range tmpContracts {

			found := false
			for _, mainContract := range contracts {
				if mainContract == tmpContract {
					found = true
				}
			}
			if !found {
				contracts = append(contracts, tmpContract)
			}
		}
	}
	return contracts, nil
}

//A function that calculates the gross and net rewards for a delegation(contract).
func CalculatePayoutForContract(contract Contract, rate float64, delegate bool, delegateAddr string) Contract {

	totalNodeRewards, _ := GetDelegateRewardsForCycle(contract.Cycle, delegateAddr)

	grossRewards := contract.SharePercentage * float64(totalNodeRewards)
	contract.GrossPayout = grossRewards
	fee := rate * grossRewards
	contract.Fee = fee
	var netRewards float64
	if delegate {
		netRewards = grossRewards
		contract.NetPayout = netRewards
		contract.Fee = 0
	} else {
		netRewards = grossRewards - fee
		contract.NetPayout = contract.NetPayout + netRewards
	}

	return contract
}

//A function that gets the rewards earned by a delegate for a specific cycle.
func GetDelegateRewardsForCycle(cycle int, delegate string) (float64, error) {
	var rtn float64
	rtn = 0
	get := "/chains/main/blocks/head/context/raw/json/contracts/index/" + delegate + "/frozen_balance/" + strconv.Itoa(cycle) + "/"
	s, err := TezosRPCGet(get)
	if err != nil {
		return rtn, errors.New("Could not get rewards for delegate: " + err.Error())
	}
	rewards, err := unMarshelFrozenBalanceRewards(s)
	if err != nil {
		return rtn, errors.New("Could not get rewards for delegate, could not parse.")
	}

	iRewards, err := strconv.Atoi(rewards.Rewards)
	if err != nil {
		return rtn, errors.New("Could not get rewards for delegate: " + err.Error())
	}

	rtn = float64(iRewards) / 1000000

	return rtn, nil
}

//A function that loops through an array of delegated contracts, and calculates the net total payout for each
//cycle included in the array, for each delegation.
func CalculateDelegateNetPayout(delegatedContracts []DelegatedContract) []DelegatedContract {
	var delegateIndex int

	for index, delegate := range delegatedContracts {
		if delegate.Delegate {
			delegateIndex = index
		}
	}

	for _, delegate := range delegatedContracts {
		if !delegate.Delegate {
			delegatedContracts[delegateIndex].TotalPayout = delegatedContracts[delegateIndex].TotalPayout + delegate.Fee
		}
	}
	return delegatedContracts
}

//A function that loops through an array of delegated contracts, and calculates the net total payout for each
//cycle included in the array, for each delegation.
func CalculateTotalPayout(delegatedContract DelegatedContract) DelegatedContract {
	for _, contract := range delegatedContract.Contracts {
		delegatedContract.TotalPayout = delegatedContract.TotalPayout + contract.NetPayout
	}
	return delegatedContract
}

//A function that loops through an array of delegated contracts, and calculates the net total payout for each
//cycle included in the array, for each delegation.
func CalculateAllTotalPayout(delegatedContracts []DelegatedContract) []DelegatedContract {
	for index, delegatedContract := range delegatedContracts {
		delegatedContracts[index] = CalculateTotalPayout(delegatedContract)
	}

	return delegatedContracts
}

func GetDelegate(delegatePhk string) (StructDelegate, error) {
	var delegate StructDelegate
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

//A helper function that tests the correctness of the shares calulated for delegations.
func checkPercentageSumForCycle(cycle int, delegatedContracts []DelegatedContract) float64 {
	var sum float64
	sum = 0
	for x := 0; x < len(delegatedContracts); x++ {
		counter := 0
		for y := 0; y < len(delegatedContracts[x].Contracts); y++ {
			if delegatedContracts[x].Contracts[y].Cycle == cycle {
				break
			}
			counter = counter + 1
		}

		sum = sum + delegatedContracts[x].Contracts[counter].SharePercentage
	}
	return sum
}

// func CalculateRollSpillage(delegatedContracts []DelegatedContract, delegateAddr string, cycle int) ([]DelegatedContract, error) {
// 	stakingBalance, err := GetDelegateStakingBalance(delegateAddr, cycle)
// 	if err != nil {
// 		return delegatedContracts, errors.New("func CalculateRollSpillage(delegatedContracts []DelegatedContract, delegateAddr string) failed: " + err.Error())
// 	}

// 	mod := math.Mod(stakingBalance, 10000)
// 	sum := mod * 10000

// 	for index, delegatedContract := range delegatedContracts {
// 		for i, contract := range delegatedContract.Contracts {
// 			if contract.Cycle == cycle {
// 				stakingBalance = stakingBalance - contract.Amount
// 				if stakingBalance < 0 {
// 					delegatedContracts[index].Contracts[i].SharePercentage = (contract.Amount - stakingBalance) / sum
// 				}
// 			}
// 		}
// 	}

// 	return delegatedContracts, nil
// }

//A helper function that reverses the order of []DelegatedContracts.
func sortDelegateContracts(delegatedContracts []DelegatedContract) []DelegatedContract {
	for i, j := 0, len(delegatedContracts)-1; i < j; i, j = i+1, j-1 {
		delegatedContracts[i], delegatedContracts[j] = delegatedContracts[j], delegatedContracts[i]
	}
	return delegatedContracts
}
