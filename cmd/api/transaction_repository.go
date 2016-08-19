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

func transactionsForAccount(db *sql.DB) func(string) ([]gledger.Transaction, error) {
	return func(u string) ([]gledger.Transaction, error) {
		var transactions []gledger.Transaction

		rows, err := db.Query(
			`SELECT transaction_uuid, account_uuid, occurred_at, payee, amount, cleared, reconciled FROM transactions where account_uuid = $1`,
			u,
		)
		if err != nil {
			return transactions, errors.Wrapf(err, "error getting transactions for %s", u)
		}
		defer rows.Close()

		for rows.Next() {
			var t gledger.Transaction
			err := rows.Scan(&t.Uuid, &t.AccountUuid, &t.OccurredAt, &t.Payee, &t.Amount, &t.Cleared, &t.Reconciled)
			if err != nil {
				return transactions, errors.Wrapf(err, "error scanning getting transactions for %s", u)
			}

			transactions = append(transactions, t)
		}
		err = rows.Err()
		return transactions, errors.Wrapf(err, "error getting transactions for %s", u)
	}
}
