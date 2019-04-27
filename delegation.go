package gotezos

import (
	"strconv"
	"sync"
)

type delegationReportJob struct {
	delegatePhk   string
	delegationPhk string
	Fee           float64
	cycle         int
	cycleRewards  int
}

type delegationReportJobResult struct {
	report DelegationReport
	err    error
}

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

// GetRewardsForDelegateForCycle gets the total rewards for a delegate earned
// and calculates the gross rewards earned by each delegation for a single cycle.
// Also includes the share of each delegation.
func (gt *GoTezos) GetRewardsForDelegateForCycle(delegatePhk string, cycle int, fee float64) (DelegateReport, error) {
	report := DelegateReport{DelegatePhk: delegatePhk, Cycle: cycle}

	cycleRewards, err := gt.GetCycleRewardsForDelegate(delegatePhk, cycle)
	if err != nil {
		return report, err
	}
	report.CycleRewards = cycleRewards

	delegations, err := gt.GetDelegationsForDelegateByCycle(delegatePhk, cycle)
	if err != nil {
		return report, err
	}

	delegationReports, gross, err := gt.getDelegationReports(delegatePhk, delegations, cycle, cycleRewards, fee)
	if err != nil {
		return report, err
	}
	report.Delegations = delegationReports
	intRewards, _ := strconv.Atoi(cycleRewards)
	selfBakeRewards := strconv.Itoa(intRewards - gross)
	report.SelfBakedRewards = selfBakeRewards
	intFeeRewards := int(float64(gross) * fee)
	report.TotalFeeRewards = strconv.Itoa(intFeeRewards)
	report.TotalRewards = strconv.Itoa(intFeeRewards + intRewards)

	return report, nil
}

// GetPayments will convert a delegate report into payments for batch pay
func (dr *DelegateReport) GetPayments() []Payment {
	payments := []Payment{}
	for _, delegate := range dr.Delegations {
		payment := Payment{}
		payment.Address = delegate.DelegationPhk
		amount, _ := strconv.ParseFloat(delegate.NetRewards, 64)
		payment.Amount = amount
		payments = append(payments, payment)
	}
	return payments
}

func (gt *GoTezos) getDelegationReports(delegate string, delegations []string, cycle int, cycleRewards string, fee float64) ([]DelegationReport, int, error) {
	reports := []DelegationReport{}

	jobs := make(chan delegationReportJob, 1000)
	results := make(chan delegationReportJobResult, 1000)

	for w := 1; w <= 50; w++ {
		go gt.delegationReportWorker(jobs, results)
	}

	bigIntCycleRewards, err := strconv.Atoi(cycleRewards)
	if err != nil {
		return reports, 0, err
	}

	for _, delegation := range delegations {
		job := delegationReportJob{delegatePhk: delegate, delegationPhk: delegation, Fee: fee, cycle: cycle, cycleRewards: bigIntCycleRewards}
		jobs <- job
	}

	totalGross := 0
	for i := 0; i < len(delegations); i++ {
		result := <-results
		if result.err != nil {
			return reports, 0, result.err
		}
		reports = append(reports, result.report)
		gross, _ := strconv.Atoi(result.report.GrossRewards)
		totalGross = totalGross + gross
	}
	return reports, totalGross, nil
}

func (gt *GoTezos) delegationReportWorker(jobs <-chan delegationReportJob, results chan<- delegationReportJobResult) {
	for j := range jobs {
		result := delegationReportJobResult{}
		report := DelegationReport{}
		report.DelegationPhk = j.delegationPhk

		share, _, err := gt.GetShareOfContract(j.delegatePhk, j.delegationPhk, j.cycle)
		if err != nil {
			result.err = err
		}
		report.Share = share
		gross := share * float64(j.cycleRewards)
		intGross := int(gross)
		report.GrossRewards = strconv.Itoa(intGross)

		fee := j.Fee * gross
		intFee := int(fee)
		report.Fee = strconv.Itoa(intFee)

		intNetRewards := intGross - intFee
		report.NetRewards = strconv.Itoa(intNetRewards)
		result.report = report
		results <- result
	}
}

// GetCycleRewardsForDelegate gets the rewards earned by a delegate for a specific cycle.
func (gt *GoTezos) GetCycleRewardsForDelegate(delegatePhk string, cycle int) (string, error) {
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
