package wallet

import (
	"errors"

	"github.com/Peqchji/go-inbound-adapter-benchmark/internal/domain"
)

type WalletErr struct {
    domain.DomainErr
}

func NewWalletErr(msg string) *WalletErr {
    return &WalletErr{
        DomainErr: domain.DomainErr{
            DomainName: "Wallet",
            Err:        errors.New(msg),
        },
    }
}


var (
	WalletErrInvalidAmount = NewWalletErr("Invalid incoming amount")
	WalletErrBalanceWillOverflow = NewWalletErr("Incoming amount will causing overflow balances")
	WalletErrNegBalance = NewWalletErr("Incoming amount will causing negative balances")
)