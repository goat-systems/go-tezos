package network

type TezosNetworkService interface {
	GetVersions() ([]Version, error)
	GetConstants() (Constants, error)
	GetChainID() (string, error)
	GetConnections() (Connections, error)
}
