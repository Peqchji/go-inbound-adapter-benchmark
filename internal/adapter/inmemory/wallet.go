package inmemoryrepository

import (
	"github.com/Peqchji/go-inbound-adapter-benchmark/internal/client/database/inmemory"
	"github.com/Peqchji/go-inbound-adapter-benchmark/internal/domain/wallet"
	"github.com/Peqchji/go-inbound-adapter-benchmark/pkg"
)

type InMemoryWalletAdapter struct {
	walletTable inmemory.IInMemoryDBTable
}

func NewInMemoryWalletAdapter(
	walletTable inmemory.IInMemoryDBTable) *InMemoryWalletAdapter {
	return &InMemoryWalletAdapter{
		walletTable: walletTable,
	}
}

func (r *InMemoryWalletAdapter) GetById(id string) pkg.Result[wallet.Wallet] {
	result := r.walletTable.GetById(id)
	if result.Err != nil {
		return pkg.Result[wallet.Wallet]{
			Err: result.Err,
		}
	}

	// Assuming we stored *wallet.Wallet as RecordDTO directly in Save
	wPtr, ok := result.Res.(*wallet.Wallet)
	if !ok {
		return pkg.Result[wallet.Wallet]{
			Err: inmemory.ErrNotFoundRecord,
		}
	}

	return pkg.Result[wallet.Wallet]{
		Res: *wPtr,
		Err: nil,
	}
}

func (r *InMemoryWalletAdapter) Save(w wallet.Wallet) error {
	// Store pointer to wallet, as standard practice often points to heap, and also if ID() is pointer receiver
	return r.walletTable.Save(&w)
}

func (r *InMemoryWalletAdapter) GetAll() pkg.Result[[]wallet.Wallet] {
	result := r.walletTable.GetAll()
	if result.Err != nil {
		return pkg.Result[[]wallet.Wallet]{
			Err: result.Err,
		}
	}

	records := result.Res
	wallets := make([]wallet.Wallet, len(records))

	for i, rec := range records {
		wPtr, ok := rec.(*wallet.Wallet)
		if !ok {
			return pkg.Result[[]wallet.Wallet]{
				Err: inmemory.ErrNotFoundRecord,
			}
		}
		wallets[i] = *wPtr
	}

	return pkg.Result[[]wallet.Wallet]{
		Res: wallets,
		Err: nil,
	}
}
