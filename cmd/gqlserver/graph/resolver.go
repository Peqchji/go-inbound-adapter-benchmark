package graph

import "github.com/Peqchji/go-inbound-adapter-benchmark/internal/domain/wallet"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	Service *wallet.WalletService
}
