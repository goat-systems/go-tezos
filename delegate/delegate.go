package delegate

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/DefinitelyNotAGoat/go-tezos/account"
	"github.com/DefinitelyNotAGoat/go-tezos/block"
	tzc "github.com/DefinitelyNotAGoat/go-tezos/client"
	"github.com/DefinitelyNotAGoat/go-tezos/network"
	"github.com/DefinitelyNotAGoat/go-tezos/snapshot"
)

// DelegateService is a struct wrapper for delegate related functions
type DelegateService struct {
	tzclient        tzc.TezosClient
	blockService    block.TezosBlockService
	snapshotService snapshot.TezosSnapshotService
	accountService  account.TezosAccountService
	constants       network.Constants
}

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

// BakingRights a representation of baking rights on the Tezos network
type BakingRights []struct {
	Level         int       `json:"level"`
	Delegate      string    `json:"delegate"`
	Priority      int       `json:"priority"`
	EstimatedTime time.Time `json:"estimated_time"`
}

// EndorsingRights is a representation of endorsing rights on the Tezos network
type EndorsingRights []struct {
	Level         int       `json:"level"`
	Delegate      string    `json:"delegate"`
	Slots         []int     `json:"slots"`
	EstimatedTime time.Time `json:"estimated_time"`
}

// DelegateReport represents a rewards report for a delegate and all their delegations for a cycle
type DelegateReport struct {
	DelegatePhk      string
	Cycle            int
	Delegations      []DelegationReport
	CycleRewards     string
	TotalFeeRewards  string
	SelfBakedRewards string
	TotalRewards     string
}

// DelegationReport represents a rewards report for a delegation in DelegateReport
type DelegationReport struct {
	DelegationPhk string
	Share         float64
	GrossRewards  string
	Fee           string
	NetRewards    string
}

// Payment is a helper struct for transfers
type Payment struct {
	Address string
	Amount  float64
}

// Delegate is representation of a delegate on the Tezos Network
type Delegate struct {
	Balance              string                 `json:"balance"`
	FrozenBalance        string                 `json:"frozen_balance"`
	FrozenBalanceByCycle []frozenBalanceByCycle `json:"frozen_balance_by_cycle"`
	StakingBalance       string                 `json:"staking_balance"`
	DelegateContracts    []string               `json:"delegated_contracts"`
	DelegatedBalance     string                 `json:"delegated_balance"`
	Deactivated          bool                   `json:"deactivated"`
	GracePeriod          int                    `json:"grace_period"`
}

// FrozenBalanceByCycle a representation of frozen balance by cycle on the Tezos network
type frozenBalanceByCycle struct {
	Cycle   int    `json:"cycle"`
	Deposit string `json:"deposit"`
	Fees    string `json:"fees"`
	Rewards string `json:"rewards"`
}

// FrozenBalanceRewards is a FrozenBalanceRewards query returned by the Tezos RPC API.
type FrozenBalanceRewards struct {
	Deposits string `json:"deposits"`
	Fees     string `json:"fees"`
	Rewards  string `json:"rewards"`
}

// NewDelegateService returns a new DelegateService
func NewDelegateService(tzclient tzc.TezosClient, blockService block.TezosBlockService, snapshotService snapshot.TezosSnapshotService, accountService account.TezosAccountService, constants network.Constants) *DelegateService {
	return &DelegateService{
		tzclient:        tzclient,
		blockService:    blockService,
		snapshotService: snapshotService,
		accountService:  accountService,
		constants:       constants,
	}
}

// GetDelegations retrieves a list of all currently delegated contracts for a delegate.
func (d *DelegateService) GetDelegations(delegatePhk string) ([]string, error) {
	rtnString := []string{}
	query := "/chains/main/blocks/head/context/delegates/" + delegatePhk + "/delegated_contracts"
	resp, err := d.tzclient.Get(query, nil)
	if err != nil {
		return rtnString, errors.Wrapf(err, "could not get delegations for '%s'", query)
	}

	delegations, err := unmarshalStringArray(resp)
	if err != nil {
		return rtnString, errors.Wrapf(err, "could not get delegations for '%s'", query)
	}
	return delegations, nil
}

// GetDelegationsAtCycle retrieves a list of all currently delegated contracts for a delegate at a specific cycle.
func (d *DelegateService) GetDelegationsAtCycle(delegatePhk string, cycle int) ([]string, error) {
	rtnString := []string{}
	snapShot, err := d.snapshotService.Get(cycle)
	if err != nil {
		return rtnString, errors.Wrapf(err, "could not get delegations for %s at cycle %d", delegatePhk, cycle)
	}

	block, err := d.blockService.Get(snapShot.AssociatedBlock)
	if err != nil {
		return rtnString, errors.Wrapf(err, "could not get delegations for %s at cycle %d", delegatePhk, cycle)
	}
	query := "/chains/main/blocks/" + block.Hash + "/context/delegates/" + delegatePhk + "/delegated_contracts"

	resp, err := d.tzclient.Get(query, nil)
	if err != nil {
		return rtnString, errors.Wrapf(err, "could not get delegations '%s'", query)
	}

	delegations, err := unmarshalStringArray(resp)
	if err != nil {
		return rtnString, errors.Wrapf(err, "could not get delegations '%s'", query)
	}

	return delegations, nil
}

// GetReport gets the total rewards for a delegate earned
// and calculates the gross rewards earned by each delegation for a single cycle.
// Also includes the share of each delegation.
func (d *DelegateService) GetReport(delegatePhk string, cycle int, fee float64) (*DelegateReport, error) {
	report := DelegateReport{DelegatePhk: delegatePhk, Cycle: cycle}

	cycleRewards, err := d.GetRewards(delegatePhk, cycle)
	if err != nil {
		return &report, errors.Wrapf(err, "could not get delegate report for %s at cycle %d", delegatePhk, cycle)
	}
	report.CycleRewards = cycleRewards

	delegations, err := d.GetDelegationsAtCycle(delegatePhk, cycle)
	if err != nil {
		return &report, errors.Wrapf(err, "could not get delegate report for %s at cycle %d", delegatePhk, cycle)
	}

	delegationReports, gross, err := d.getDelegationReports(delegatePhk, delegations, cycle, cycleRewards, fee)
	if err != nil {
		return &report, errors.Wrapf(err, "could not get delegate report for %s at cycle %d", delegatePhk, cycle)
	}
	report.Delegations = delegationReports
	intRewards, _ := strconv.Atoi(cycleRewards)
	selfBakeRewards := strconv.Itoa(intRewards - gross)
	report.SelfBakedRewards = selfBakeRewards
	intFeeRewards := int(float64(gross) * fee)
	report.TotalFeeRewards = strconv.Itoa(intFeeRewards)
	report.TotalRewards = strconv.Itoa(intFeeRewards + intRewards)

	return &report, nil
}

// GetPayments will convert a delegate report into payments for batch pay with a minimum requirement in mutez
func (dr *DelegateReport) GetPayments(minimum int) []Payment {
	payments := []Payment{}
	for _, delegate := range dr.Delegations {
		payment := Payment{}
		payment.Address = delegate.DelegationPhk
		amount, _ := strconv.ParseFloat(delegate.NetRewards, 64)
		payment.Amount = amount

		if payment.Amount >= float64(minimum) && payment.Amount != 0 {
			payments = append(payments, payment)
		}
	}
	return payments
}

func (d *DelegateService) getDelegationReports(delegate string, delegations []string, cycle int, cycleRewards string, fee float64) ([]DelegationReport, int, error) {
	reports := []DelegationReport{}

	numberOfDelegators := len(delegations)

	jobs := make(chan delegationReportJob, numberOfDelegators)
	results := make(chan delegationReportJobResult, numberOfDelegators)

	bigIntCycleRewards, err := strconv.Atoi(cycleRewards)
	if err != nil {
		return reports, 0, errors.Wrap(err, "could not get delegation reports")
	}

	for _, delegation := range delegations {
		job := delegationReportJob{delegatePhk: delegate, delegationPhk: delegation, Fee: fee, cycle: cycle, cycleRewards: bigIntCycleRewards}
		jobs <- job
	}

	stakingBalance, err := d.GetStakingBalance(delegate, cycle)
	if err != nil {
		return reports, 0, errors.Errorf("could not get staking balance of delegate %s: %v", delegate, err)
	}

	snapShot, err := d.snapshotService.Get(cycle)
	if err != nil {
		return reports, 0, errors.Errorf("could not get snap shot at %d cycle: %v", cycle, err)
	}

	block, err := d.blockService.Get(snapShot.AssociatedBlock)
	if err != nil {
		return reports, 0, errors.Errorf("could not get associated snap shot block at %d cycle: %v", cycle, err)
	}

	for w := 1; w <= 50; w++ {
		go d.delegationReportWorker(jobs, results, block.Hash, stakingBalance)
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

func (d *DelegateService) delegationReportWorker(jobs <-chan delegationReportJob, results chan<- delegationReportJobResult,
	associatedBlockHash string, stakingBalance float64) {
	for j := range jobs {
		result := delegationReportJobResult{}
		report := DelegationReport{}
		report.DelegationPhk = j.delegationPhk

		share, _, err := d.getShareOfContract(j.delegationPhk, associatedBlockHash, stakingBalance)
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

// GetRewards gets the rewards earned by a delegate for a specific cycle.
func (d *DelegateService) GetRewards(delegatePhk string, cycle int) (string, error) {
	rewards := FrozenBalanceRewards{}
	level := (cycle+1)*(d.constants.BlocksPerCycle) + 1

	head, err := d.blockService.Get(level)
	if err != nil {
		return "", errors.Wrapf(err, "could not get rewards for %s at %d cycle", delegatePhk, cycle)
	}

	query := "/chains/main/blocks/" + head.Hash + "/context/raw/json/contracts/index/" + delegatePhk + "/frozen_balance/" + strconv.Itoa(cycle) + "/"
	resp, err := d.tzclient.Get(query, nil)
	if err != nil {
		return "", errors.Wrapf(err, "could not get rewards '%s'", query)
	}
	rewards, err = rewards.unmarshalJSON(resp)
	if err != nil {
		return rewards.Rewards, errors.Wrapf(err, "could not get rewards '%s'", query)
	}

	if rewards.Rewards == "" {
		return rewards.Rewards, errors.Errorf("could not get rewards '%s', empty string", query)
	}

	return rewards.Rewards, nil
}

// getShareOfContract returns the share of a delegation for a specific cycle.
func (d *DelegateService) getShareOfContract(delegationPhk, associatedBlockHash string, stakingBalance float64) (float64, float64, error) {
	delegationBalance, err := d.accountService.GetBalanceAtBlock(delegationPhk, associatedBlockHash)
	if err != nil {
		return 0, 0, errors.Errorf("could not get share of contract %s: %v", delegationPhk, err)
	}

	return delegationBalance / stakingBalance, delegationBalance, nil
}

// GetDelegate retrieves information about a delegate at the head block
func (d *DelegateService) GetDelegate(delegatePhk string) (Delegate, error) {
	delegate := Delegate{}
	get := "/chains/main/blocks/head/context/delegates/" + delegatePhk
	resp, err := d.tzclient.Get(get, nil)
	if err != nil {
		return delegate, errors.Wrapf(err, "could not get delegate '%s'", get)
	}
	delegate, err = delegate.unmarshalJSON(resp)
	if err != nil {
		return delegate, errors.Wrapf(err, "could not get delegate '%s'", get)
	}

	return delegate, nil
}

// GetStakingBalanceAtCycle gets the staking balance of a delegate at a specific cycle
func (d *DelegateService) GetStakingBalanceAtCycle(delegateAddr string, cycle int) (string, error) {
	balance := ""
	snapShot, err := d.snapshotService.Get(cycle)
	if err != nil {
		return balance, errors.Wrapf(err, "could not get staking balance for %s at cycle %d", delegateAddr, cycle)
	}
	query := "/chains/main/blocks/" + snapShot.AssociatedHash + "/context/delegates/" + delegateAddr + "/staking_balance"
	resp, err := d.tzclient.Get(query, nil)
	if err != nil {
		return balance, errors.Wrapf(err, "could not get staking balance '%s'", query)
	}
	balance, err = unmarshalString(resp)
	if err != nil {
		return balance, errors.Wrapf(err, "could not get staking balance '%s'", query)
	}

	return balance, nil
}

// GetBakingRights gets the baking rights for a specific cycle
func (d *DelegateService) GetBakingRights(cycle int) (BakingRights, error) {
	bakingRights := BakingRights{}

	snapShot, err := d.snapshotService.Get(cycle)
	if err != nil {
		return bakingRights, errors.Wrapf(err, "could not get baking rights for cycle %d", err)
	}

	params := make(map[string]string)
	params["cycle"] = strconv.Itoa(cycle)

	query := "/chains/main/blocks/" + snapShot.AssociatedHash + "/helpers/baking_rights"
	resp, err := d.tzclient.Get(query, params)
	if err != nil {
		return bakingRights, errors.Wrapf(err, "could not get baking rights '%s'", err)
	}

	bakingRights, err = bakingRights.unmarshalJSON(resp)
	if err != nil {
		return bakingRights, errors.Wrapf(err, "could not get baking rights '%s'", err)
	}

	return bakingRights, nil
}

// GetBakingRightsForDelegate gets the baking rights for a delegate at a specific cycle with a certain priority level
func (d *DelegateService) GetBakingRightsForDelegate(cycle int, delegatePhk string, priority int) (BakingRights, error) {
	bakingRights := BakingRights{}

	snapShot, err := d.snapshotService.Get(cycle)
	if err != nil {
		return bakingRights, errors.Wrapf(err, "could not get baking rights for delegate %s at cycle %d", delegatePhk, cycle)
	}

	params := make(map[string]string)
	params["cycle"] = strconv.Itoa(cycle)
	params["delegate"] = delegatePhk
	params["max_priority"] = strconv.Itoa(priority)

	query := "/chains/main/blocks/" + snapShot.AssociatedHash + "/helpers/baking_rights"
	resp, err := d.tzclient.Get(query, params)
	if err != nil {
		return bakingRights, errors.Wrapf(err, "could not get baking rights for delegate '%s'", query)
	}

	bakingRights, err = bakingRights.unmarshalJSON(resp)
	if err != nil {
		return bakingRights, errors.Wrapf(err, "could not get baking rights for delegate '%s'", query)
	}

	return bakingRights, nil
}

// GetEndorsingRightsForDelegate gets the endorsing rights for a specific cycle
func (d *DelegateService) GetEndorsingRightsForDelegate(cycle int, delegatePhk string) (EndorsingRights, error) {
	endorsingRights := EndorsingRights{}

	snapShot, err := d.snapshotService.Get(cycle)
	if err != nil {
		return endorsingRights, errors.Wrapf(err, "could not get endorsing rights for delegate %s at cycle %d", delegatePhk, cycle)
	}

	params := make(map[string]string)
	params["cycle"] = strconv.Itoa(cycle)
	params["delegate"] = delegatePhk

	query := "/chains/main/blocks/" + snapShot.AssociatedHash + "/helpers/endorsing_rights"
	resp, err := d.tzclient.Get(query, params)
	if err != nil {
		return endorsingRights, errors.Wrapf(err, "could not get endorsing rights for delegate '%s'", query)
	}

	endorsingRights, err = endorsingRights.unmarshalJSON(resp)
	if err != nil {
		return endorsingRights, errors.Wrapf(err, "could not get endorsing rights for delegate '%s'", query)
	}

	return endorsingRights, nil
}

// GetEndorsingRights gets the endorsing rights for a specific cycle
func (d *DelegateService) GetEndorsingRights(cycle int) (EndorsingRights, error) {
	endorsingRights := EndorsingRights{}

	snapShot, err := d.snapshotService.Get(cycle)
	if err != nil {
		return endorsingRights, errors.Wrapf(err, "could not get endorsing rights for cycle %d", cycle)
	}

	params := make(map[string]string)
	params["cycle"] = strconv.Itoa(cycle)

	get := "/chains/main/blocks/" + snapShot.AssociatedHash + "/helpers/endorsing_rights"
	resp, err := d.tzclient.Get(get, params)
	if err != nil {
		return endorsingRights, errors.Wrapf(err, "could not get endorsing rights for cycle '%s'", get)
	}

	endorsingRights, err = endorsingRights.unmarshalJSON(resp)
	if err != nil {
		return endorsingRights, errors.Wrapf(err, "could not get endorsing rights for cycle '%s'", get)
	}

	return endorsingRights, nil
}

// GetAllDelegatesByHash gets a list of all tz1 addresses at a certain hash
func (d *DelegateService) GetAllDelegatesByHash(hash string) ([]string, error) {
	delList := []string{}
	query := "/chains/main/blocks/" + hash + "/context/delegates"
	resp, err := d.tzclient.Get(query, nil)
	if err != nil {
		return delList, errors.Wrapf(err, "could not get all delegates '%s'", query)
	}
	delList, err = unmarshalStringArray(resp)
	if err != nil {
		return delList, errors.Wrapf(err, "could not get all delegates '%s'", query)
	}
	return delList, nil
}

// GetAllDelegates a list of all tz1 addresses at the head block
func (d *DelegateService) GetAllDelegates() ([]string, error) {
	delList := []string{}
	query := "/chains/main/blocks/head/context/delegates?active"
	resp, err := d.tzclient.Get(query, nil)
	if err != nil {
		return delList, errors.Wrapf(err, "could not get all delegates '%s'", query)
	}
	delList, err = unmarshalStringArray(resp)
	if err != nil {
		return delList, errors.Wrapf(err, "could not get all delegates '%s'", query)
	}
	return delList, nil
}

// GetStakingBalance gets the staking balance for a delegate at a specific snapshot for a cycle.
func (d *DelegateService) GetStakingBalance(delegateAddr string, cycle int) (float64, error) {

	snapShot, err := d.snapshotService.Get(cycle)
	if err != nil {
		return 0, errors.Wrapf(err, "could not get staking balance for %s at cycle %d", delegateAddr, cycle)
	}

	block, err := d.blockService.Get(snapShot.AssociatedBlock)
	if err != nil {
		return 0, errors.Wrapf(err, "could not get staking balance for %s at cycle %d", delegateAddr, cycle)
	}

	query := "/chains/main/blocks/" + block.Hash + "/context/delegates/" + delegateAddr + "/staking_balance"

	resp, err := d.tzclient.Get(query, nil)
	if err != nil {
		return 0, errors.Wrapf(err, "could not get staking balance '%s'", query)
	}

	strBalance, err := unmarshalString(resp)
	if err != nil {
		return 0, errors.Wrapf(err, "could not get staking balance '%s'", query)
	}

	floatBalance, err := strconv.ParseFloat(strBalance, 64) //TODO error checking
	if err != nil {
		return 0, errors.Wrapf(err, "could not get staking balance '%s'", query)
	}

	return floatBalance / account.MUTEZ, nil
}

// unmarshalJSON unmarshalls bytes into StructDelegate
func (d *Delegate) unmarshalJSON(v []byte) (Delegate, error) {
	delegate := Delegate{}
	err := json.Unmarshal(v, &delegate)
	if err != nil {
		return delegate, errors.Wrap(err, "could not unmarshal bytes into Delegate")
	}
	return delegate, nil
}

// UnmarshalJSON unmarhsels bytes into Baking_Rights
func (br *BakingRights) unmarshalJSON(v []byte) (BakingRights, error) {
	bakingRights := BakingRights{}
	err := json.Unmarshal(v, &bakingRights)
	if err != nil {
		return bakingRights, errors.Wrap(err, "could not unmarshal bytes into BakingRights")
	}
	return bakingRights, nil
}

// UnmarshalJSON unmarhsels bytes into Endorsing_Rights
func (er *EndorsingRights) unmarshalJSON(v []byte) (EndorsingRights, error) {
	endorsingRights := EndorsingRights{}
	err := json.Unmarshal(v, &endorsingRights)
	if err != nil {
		return endorsingRights, errors.Wrap(err, "could not unmarshal bytes into EndorsingRights")
	}
	return endorsingRights, nil
}

// UnmarshalJSON unmarshals the bytes received as a parameter, into the type SnapShotQuery.
func (fb *FrozenBalanceRewards) unmarshalJSON(v []byte) (FrozenBalanceRewards, error) {
	frozenBalance := FrozenBalanceRewards{}
	err := json.Unmarshal(v, &frozenBalance)
	if err != nil {
		return frozenBalance, errors.Wrap(err, "could not unmarshal bytes into FrozenBalanceRewards")
	}
	return frozenBalance, nil
}

// UnmarshalJSON unmarshals the bytes received as a parameter, into the type an array of strings.
func unmarshalStringArray(v []byte) ([]string, error) {
	var strs []string
	err := json.Unmarshal(v, &strs)
	if err != nil {
		return strs, errors.Wrap(err, "could not unmarshal bytes into []string")
	}
	return strs, nil
}

// unmarshalString unmarshals the bytes received as a parameter, into the type string.
func unmarshalString(v []byte) (string, error) {
	var str string
	err := json.Unmarshal(v, &str)
	if err != nil {
		return str, errors.Wrap(err, "could not unmarshal bytes to string")
	}
	return str, nil
}
