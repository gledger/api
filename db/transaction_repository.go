package db

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/zombor/gledger"
)

func SaveTransaction(db *sql.DB) func(gledger.Transaction) error {
	return func(t gledger.Transaction) error {
		_, err := db.Exec(
			`INSERT INTO transactions VALUES ($1, $2, $3, $4, $5, $6, $7, now(), now())`,
			t.Uuid, t.AccountUuid, t.OccurredAt, t.Payee, t.Amount, t.Cleared, t.Reconciled,
		)

		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == ForeignKeyViolation {
				return notFoundError(fmt.Sprintf("Account %s not found", t.AccountUuid))
			}
		}

		return errors.Wrap(err, "error writing transaction")
	}
}

func TransactionsForAccount(db *sql.DB) func(string) ([]gledger.Transaction, error) {
	return func(u string) ([]gledger.Transaction, error) {
		var transactions []gledger.Transaction

		var uuid string
		err := db.QueryRow(`SELECT account_uuid FROM accounts WHERE account_uuid = $1`, u).Scan(&uuid)
		if err == sql.ErrNoRows {
			return transactions, notFoundError(fmt.Sprintf("Account %s not found", u))
		}

		rows, err := db.Query(
			`SELECT
				transaction_uuid,
				account_uuid,
				occurred_at,
				payee,
				amount,
				sum(amount) OVER (PARTITION BY account_uuid ORDER BY occurred_at, created_at),
				cleared,
				reconciled
			FROM transactions where account_uuid = $1`,
			u,
		)
		if err != nil {
			return transactions, errors.Wrapf(err, "error getting transactions for %s", u)
		}
		defer rows.Close()

		for rows.Next() {
			var t gledger.Transaction
			err := rows.Scan(&t.Uuid, &t.AccountUuid, &t.OccurredAt, &t.Payee, &t.Amount, &t.RollingTotal, &t.Cleared, &t.Reconciled)
			if err != nil {
				return transactions, errors.Wrapf(err, "error scanning getting transactions for %s", u)
			}

			transactions = append(transactions, t)
		}
		err = rows.Err()
		return transactions, errors.Wrapf(err, "error getting transactions for %s", u)
	}
}
