package main

import (
	"database/sql"

	"github.com/pkg/errors"

	"github.com/zombor/gledger"
)

func saveTransaction(db *sql.DB) func(gledger.Transaction) error {
	return func(t gledger.Transaction) error {
		_, err := db.Exec(
			`INSERT INTO transactions VALUES ($1, $2, $3, $4, $5, $6, $7, now(), now())`,
			t.Uuid, t.AccountUuid, t.OccurredAt, t.Payee, t.Amount, t.Cleared, t.Reconciled,
		)

		return errors.Wrap(err, "error writing transaction")
	}
}
