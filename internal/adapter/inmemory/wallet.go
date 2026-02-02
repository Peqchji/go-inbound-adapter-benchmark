package inmemoryrepository

import (
	"github.com/Peqchji/go-inbound-adapter-benchmark/internal/client/database/inmemory"
	"github.com/Peqchji/go-inbound-adapter-benchmark/internal/domain/wallet"
	"github.com/Peqchji/go-inbound-adapter-benchmark/pkg"
)

type InMemoryWalletRepository struct {
	walletTable inmemory.IInMemoryDBTable
}

func NewInMemoryWalletRepository(
	walletTable inmemory.IInMemoryDBTable) *InMemoryWalletRepository {
	return &InMemoryWalletRepository{
		walletTable: walletTable,
	}
}

func (r *InMemoryWalletRepository) GetById(id string) pkg.Result[*wallet.Wallet] {
	result := r.walletTable.GetById(id)
	if result.Err != nil {
		return pkg.Result[*wallet.Wallet]{
			Res: nil,
			Err: result.Err,
		}
	}

	w, ok := result.Res.(WalletDTO)
	if !ok {
		return pkg.Result[*wallet.Wallet]{
			Res: nil,
			Err: inmemory.ErrNotFoundRecord,
		}
	}

	domainWallet := w.ToDomain()

	return pkg.Result[*wallet.Wallet]{
		Res: &domainWallet,
		Err: nil,
	}
}

func (r *InMemoryWalletRepository) Save(w wallet.Wallet) error {
	walletDto := WalletDTO{}.FromDomain(w)

	return r.walletTable.Save(walletDto)
}

func (r *InMemoryWalletRepository) GetAll() pkg.Result[[]*wallet.Wallet] {
	result := r.walletTable.GetAll()
	if result.Err != nil {
		return pkg.Result[[]*wallet.Wallet]{
			Res: nil,
			Err: result.Err,
		}
	}

	records := result.Res
	wallets := make([]*wallet.Wallet, len(records))

	for i, rec := range records {
		w, ok := rec.(WalletDTO)
		if !ok {
			return pkg.Result[[]*wallet.Wallet]{
				Res: nil,
				Err: inmemory.ErrNotFoundRecord,
			}
		}

		domainWallet := w.ToDomain()
		wallets[i] = &domainWallet
	}

	return pkg.Result[[]*wallet.Wallet]{
		Res: wallets,
		Err: nil,
	}
}
