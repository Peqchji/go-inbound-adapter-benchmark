package wallet

type WalletRepository interface {
	GetById(id string) (Wallet, error)
	Save(wallet Wallet) error
	GetAll() ([]Wallet, error)
}