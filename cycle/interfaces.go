package cycle

type TezosCycleService interface {
	GetCurrent() (int, error)
}
