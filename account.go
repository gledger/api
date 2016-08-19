package gledger

import "github.com/twinj/uuid"

type Account struct {
	Uuid, Name, Type string
	Active           bool
}

type AccountService interface {
	Create(Account) (Account, error)
	All() ([]Account, error)
}

func NewAccountService(
	s func(Account) error,
	a func() ([]Account, error),
) accountService {
	return accountService{saveAccount: s, allAccounts: a}
}

type accountService struct {
	saveAccount func(Account) error
	allAccounts func() ([]Account, error)
}

func (as accountService) Create(a Account) (Account, error) {
	if a.Uuid == "" {
		a.Uuid = uuid.NewV4().String()
	}

	return a, as.saveAccount(a)
}

func (as accountService) All() ([]Account, error) {
	return as.allAccounts()
}
