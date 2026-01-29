package wallet

import "github.com/Peqchji/go-inbound-adapter-benchmark/pkg"

type WalletRepository interface {
	GetById(id string) 	pkg.Result[Wallet]
	Save(wallet Wallet) error
	GetAll() 			pkg.Result[[]Wallet]
}
