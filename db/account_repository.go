package db

import (
	"database/sql"

	"github.com/pkg/errors"

	"github.com/zombor/gledger"
)

func SaveAccount(db *sql.DB) func(gledger.Account) error {
	return func(a gledger.Account) error {
		_, err := db.Exec(`INSERT INTO accounts VALUES ($1, $2, $3, $4, now(), now())`, a.Uuid, a.Name, a.Type, a.Active)
		return errors.Wrap(err, "error writing account")
	}
}

func AllAccounts(db *sql.DB) func() ([]gledger.Account, error) {
	return func() ([]gledger.Account, error) {
		var accounts []gledger.Account

		rows, err := db.Query(
			`SELECT
				account_uuid, name, type, active, sum(transactions.amount)
			FROM accounts
			LEFT JOIN transactions USING(account_uuid)
			GROUP BY accounts.account_uuid`,
		)
		if err != nil {
			return accounts, errors.Wrap(err, "error getting all accounts")
		}
		defer rows.Close()

		for rows.Next() {
			var a gledger.Account
			var b sql.NullInt64
			err := rows.Scan(&a.Uuid, &a.Name, &a.Type, &a.Active, &b)
			if err != nil {
				return accounts, errors.Wrap(err, "error scanning all accounts")
			}

			if b.Valid {
				a.Balance = b.Int64
			}

			accounts = append(accounts, a)
		}
		err = rows.Err()
		return accounts, errors.Wrap(err, "error getting all account rows")
	}
}

func ReadAccount(db *sql.DB) func(string) (gledger.Account, error) {
	return func(uuid string) (gledger.Account, error) {
		var account gledger.Account
		var b sql.NullInt64

		err := db.QueryRow(
			`SELECT
				account_uuid, name, type, active, sum(transactions.amount)
			FROM accounts
			LEFT JOIN transactions USING(account_uuid)
			WHERE account_uuid = $1
			GROUP BY accounts.account_uuid`,
			uuid,
		).Scan(
			&account.Uuid,
			&account.Name,
			&account.Type,
			&account.Active,
			&b,
		)

		if b.Valid {
			account.Balance = b.Int64
		}

		return account, errors.Wrapf(err, "error reading account %s", uuid)
	}
}
