package gledger

import (
	"time"

	"github.com/twinj/uuid"
)

type Transaction struct {
	UUID, AccountUUID, Payee string
	OccurredAt               time.Time
	Amount, RollingTotal     int64
	Cleared, Reconciled      bool
}

type TransactionService interface {
	Create(Transaction) (Transaction, error)
	AllForAccount(string) ([]Transaction, error)
}

func NewTransactionService(
	s func(Transaction) error,
	r func(string) ([]Transaction, error),
) transactionService {
	return transactionService{saveTransaction: s, transactionsForAccount: r}
}

type transactionService struct {
	saveTransaction        func(Transaction) error
	transactionsForAccount func(string) ([]Transaction, error)
}

func (ts transactionService) Create(t Transaction) (Transaction, error) {
	if t.UUID == "" {
		t.UUID = uuid.NewV4().String()
	}

	return t, ts.saveTransaction(t)
}

func (ts transactionService) AllForAccount(u string) ([]Transaction, error) {
	return ts.transactionsForAccount(u)
}
