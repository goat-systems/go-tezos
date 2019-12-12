package delegate

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/DefinitelyNotAGoat/go-tezos/account"
	gotezos "github.com/DefinitelyNotAGoat/go-tezos/v2"
)

// DelegateService is a struct wrapper for delegate related functions
type DelegateService struct {
	gt *gotezos.GoTezos
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



// NewDelegateService returns a new DelegateService
func NewDelegateService(gt *gotezos.GoTezos) *DelegateService {
	return &DelegateService{
		gt: gt,
	}
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

// getShareOfContract returns the share of a delegation for a specific cycle.
func (d *DelegateService) getShareOfContract(delegationPhk, associatedBlockHash string, stakingBalance float64) (float64, float64, error) {
	delegationBalance, err := d.accountService.GetBalanceAtBlock(delegationPhk, associatedBlockHash)
	if err != nil {
		return 0, 0, errors.Errorf("could not get share of contract %s: %v", delegationPhk, err)
	}

	return delegationBalance / stakingBalance, delegationBalance, nil
}
