package account

type TezosAccountService interface {
	GetBalanceAtSnapshot(tezosAddr string, cycle int) (float64, error)
	GetBalance(tezosAddr string) (float64, error)
	GetBalanceAtBlock(tezosAddr string, id interface{}) (float64, error)
	CreateWallet(mnenomic string, password string) (Wallet, error)
	ImportWallet(address, public, secret string) (Wallet, error)
	ImportEncryptedWallet(pw, encKey string) (Wallet, error)
}
