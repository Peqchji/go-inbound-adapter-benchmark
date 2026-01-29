package inmemoryrepository

import "github.com/Peqchji/go-inbound-adapter-benchmark/internal/domain/wallet"

type WalletDTO struct {
	walletId       string
	balances       uint64
	ownerID        string
	ownerFirstname string
	ownerLastname  string
}

func (r WalletDTO) ID() string {
	return r.walletId
}

func (r WalletDTO) ToDomain() wallet.Wallet {
	owner := wallet.NewOwner(
		r.ownerID,
		r.ownerFirstname,
		r.ownerLastname,
	)

	return wallet.NewWallet(
		r.walletId,
		r.balances,
		owner,
	)
}

func (r WalletDTO) FromDomain(w wallet.Wallet) *WalletDTO {
	return &WalletDTO{
		walletId:       w.ID(),
		balances:       w.Balance(),
		ownerID:        w.Owner().ID(),
		ownerFirstname: w.Owner().Firstname(),
		ownerLastname:  w.Owner().Lastname(),
	}
}
