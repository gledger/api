package db

import (
	"database/sql"

	"github.com/gledger/api"
	"github.com/pkg/errors"
)

func SaveAccount(exec func(query string, args ...interface{}) (sql.Result, error)) func(gledger.Account) error {
	return func(a gledger.Account) error {
		_, err := exec(`INSERT INTO accounts VALUES ($1, $2, $3, $4, now(), now())`, a.UUID, a.Name, a.Type, a.Active)
		return errors.Wrap(err, "error writing account")
	}
}

func AllAccounts(query func(query string, args ...interface{}) (*sql.Rows, error)) func() ([]gledger.Account, error) {
	return func() ([]gledger.Account, error) {
		var accounts []gledger.Account

		rows, err := query(
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
			err := rows.Scan(&a.UUID, &a.Name, &a.Type, &a.Active, &b)
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

func ReadAccount(queryRow func(string, ...interface{}) *sql.Row) func(string) (gledger.Account, error) {
	return func(uuid string) (gledger.Account, error) {
		var account gledger.Account
		var b sql.NullInt64

		err := queryRow(
			`SELECT
				account_uuid, name, type, active, sum(transactions.amount)
			FROM accounts
			LEFT JOIN transactions USING(account_uuid)
			WHERE account_uuid = $1
			GROUP BY accounts.account_uuid`,
			uuid,
		).Scan(
			&account.UUID,
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
