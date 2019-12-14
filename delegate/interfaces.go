package delegate

type TezosDelegateService interface {
	GetDelegations(delegatePhk string) ([]string, error)
	GetDelegationsAtCycle(delegatePhk string, cycle int) ([]string, error)
	GetReport(delegatePhk string, cycle int, fee float64, earned bool) (*DelegateReport, error)
	// GetPayments(minimum int) []Payment
	GetFrozenBalanceRewards(delegatePhk string, cycle int) (string, error)
	GetScheduledRewards(delegatePhk string, cycle int) (string, error)
	GetDelegate(delegatePhk string) (Delegate, error)
	GetStakingBalanceAtCycle(delegateAddr string, cycle int) (string, error)
	GetBakingRights(cycle int) (BakingRights, error)
	GetBakingRightsForDelegate(cycle int, delegatePhk string, priority int) (BakingRights, error)
	GetBakingRightsAtLevel(level int) (BakingRights, error)
	GetBakingRightsForDelegateAtHashLevel(blockHash string, level int, delegatePhk string) (BakingRights, error)
	GetEndorsingRightsForDelegate(cycle int, delegatePhk string) (EndorsingRights, error)
	GetEndorsingRightsForDelegateAtHashLevel(blockHash string, level int, delegatePhk string) (EndorsingRights, error)
	GetEndorsingRights(cycle int) (EndorsingRights, error)
	GetEndorsingRightsAtLevel(level int) (EndorsingRights, error)
	GetAllDelegatesByHash(hash string) ([]string, error)
	GetAllDelegates() ([]string, error)
	GetStakingBalance(delegateAddr string, cycle int) (float64, error)
}
