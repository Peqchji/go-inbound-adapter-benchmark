package wallet

import (
	"math"
)

type Owner struct {
	id        string
	firstname string
	lastname  string
}

type Wallet struct {
	id       string
	balances float64
	owner    Owner
}

func (w *Wallet) GetBalance() float64 {
	return w.balances
}

func (w *Wallet) Deposit(amount float64) *WalletErr {
	if amount <= 0 {
		return WalletErrInvalidAmount
	}

	preCal := w.balances + amount
	if math.IsInf(preCal, 1) {
		return WalletErrBalanceWillOverflow
	}

	w.balances = preCal

	return nil
}

func (w *Wallet) Withdrawn(amount float64) *WalletErr {
	if amount <= 0 {
		return WalletErrInvalidAmount
	}

	preCal := w.balances - amount
	if preCal < 0 {
		return WalletErrNegBalance
	}

	w.balances = preCal

	return nil
}