package gledger

import "github.com/twinj/uuid"

type Account struct {
	Uuid, Name, Type string
	Active           bool
}

type AccountService interface {
	Create(Account) (Account, error)
	All() ([]Account, error)
	Read(string) (Account, error)
}

func NewAccountService(
	s func(Account) error,
	a func() ([]Account, error),
	r func(string) (Account, error),
) accountService {
	return accountService{saveAccount: s, allAccounts: a, readAccount: r}
}

type accountService struct {
	saveAccount func(Account) error
	allAccounts func() ([]Account, error)
	readAccount func(string) (Account, error)
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

func (as accountService) Read(u string) (Account, error) {
	return as.readAccount(u)
}
