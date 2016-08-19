package gledger

import (
	"time"

	"github.com/twinj/uuid"
)

type Transaction struct {
	Uuid, AccountUuid, Payee string
	OccurredAt               time.Time
	Amount                   int64
	Cleared, Reconciled      bool
}

type TransactionService interface {
	Create(Transaction) (Transaction, error)
}

func NewTransactionService(
	s func(Transaction) error,
) transactionService {
	return transactionService{saveTransaction: s}
}

type transactionService struct {
	saveTransaction func(Transaction) error
}

func (ts transactionService) Create(t Transaction) (Transaction, error) {
	if t.Uuid == "" {
		t.Uuid = uuid.NewV4().String()
	}

	return t, ts.saveTransaction(t)
}
