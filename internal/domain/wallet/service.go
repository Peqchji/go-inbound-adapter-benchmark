package wallet

import "github.com/Peqchji/go-inbound-adapter-benchmark/pkg"

type WalletService struct {
	repo WalletRepository
}

func NewWalletService(repo WalletRepository) *WalletService {
	return &WalletService{repo: repo}
}

func (s *WalletService) GetWallet(id string) pkg.Result[Wallet] {
	return s.repo.GetById(id)
}

func (s *WalletService) CreateWallet(ownerID, firstname, lastname string) pkg.Result[Wallet] {
	owner := NewOwner(ownerID, firstname, lastname)
	wallet := NewWallet(ownerID, 0, owner)

	err := s.repo.Save(wallet)
	return pkg.Result[Wallet]{
		Res: wallet,
		Err: err,
	}
}
