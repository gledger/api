package gledger

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_TransactionCreateAssignsUuid_IfNotExists(t *testing.T) {
	a := Transaction{Payee: "test"}
	var saveCalled bool
	saveErr := errors.New("boom")
	saveTransaction := func(txn Transaction) error {
		assert.Equal(t, a.Payee, txn.Payee)
		assert.NotEmpty(t, txn.UUID)
		saveCalled = true
		return saveErr
	}

	svc := transactionService{saveTransaction: saveTransaction}
	_, err := svc.Create(a)
	assert.Equal(t, saveErr, err)
	assert.True(t, saveCalled)
}
