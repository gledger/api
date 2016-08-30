package db

import (
	"database/sql"

	"github.com/gledger/api"
	"github.com/pkg/errors"
)

func SaveEnvelope(exec func(query string, args ...interface{}) (sql.Result, error)) func(gledger.Envelope) error {
	return func(e gledger.Envelope) error {
		_, err := exec(`INSERT INTO envelopes VALUES (
			$1, $2, now(), now(), $3
		)`, e.UUID, e.Name, e.Type)
		return errors.Wrapf(err, `could not create envelope %s`, e.UUID)
	}
}

func AllEnvelopes(query func(query string, args ...interface{}) (*sql.Rows, error)) func() ([]gledger.Envelope, error) {
	return func() (es []gledger.Envelope, err error) {
		rows, err := query(
			`SELECT envelopes.envelope_uuid, envelopes.name, envelopes.type, coalesce(sum(amount), 0)
				FROM envelopes
				LEFT JOIN transactions using(envelope_uuid)
				GROUP BY envelopes.envelope_uuid, envelopes.name, envelopes.type`,
		)
		if err != nil {
			return
		}
		defer rows.Close()

		for rows.Next() {
			var e gledger.Envelope
			err = rows.Scan(&e.UUID, &e.Name, &e.Type, &e.Balance)
			if err != nil {
				return
			}
			es = append(es, e)
		}
		err = rows.Err()
		return
	}
}
