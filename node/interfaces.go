package node

type TezosNodeService interface {
	Bootstrapped() (Bootstrap, error)
	CommitHash() (string, error)
}
