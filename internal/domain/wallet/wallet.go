package wallet

type Owner struct {
	id        string
	firstname string
	lastname  string
}

func NewOwner(id string, firstname string, lastname string) Owner {
	return Owner{
		id:        id,
		firstname: firstname,
		lastname:  lastname,
	}
}

func (o Owner) ID() string {
	return o.id
}

func (o Owner) Firstname() string {
	return o.firstname
}

func (o Owner) Lastname() string {
	return o.lastname
}

//---------------------------------------------------------------------------//

type Wallet struct {
	id       string
	balances uint64
	owner    Owner
}

func NewWallet(id string, balances uint64, owner Owner) Wallet {
	return Wallet{
		id:       id,
		balances: balances,
		owner:    owner,
	}
}

func (w *Wallet) Owner() Owner {
	return w.owner
}

func (w *Wallet) ID() string {
	return w.id
}

func (w *Wallet) Balance() uint64 {
	return w.balances
}

func (w *Wallet) Deposit(amount uint64) error {
	if amount <= 0 {
		return WalletErrInvalidAmount
	}

	updatedBalance := w.balances + amount
	if updatedBalance < w.balances {
		return WalletErrBalanceWillOverflow
	}

	w.balances = updatedBalance

	return nil
}

func (w *Wallet) Withdraw(amount uint64) error {
	if amount <= 0 {
		return WalletErrInvalidAmount
	}

	updatedBalance := w.balances - amount
	if updatedBalance > w.balances {
		return WalletErrBalanceWillUnderflow
	}

	w.balances = updatedBalance

	return nil
}
