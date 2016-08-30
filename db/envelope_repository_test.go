package db

import (
	"database/sql"
	"os"
	"testing"

	"github.com/gledger/api"
	_ "github.com/lib/pq"
	"github.com/twinj/uuid"

	"github.com/stretchr/testify/suite"
)

type EnvelopeRepositorySuite struct {
	suite.Suite

	db *sql.DB
	tx *sql.Tx
}

func (s *EnvelopeRepositorySuite) Test_SaveEnvelope() {
	e := gledger.Envelope{
		UUID: uuid.NewV4().String(),
		Type: "income",
		Name: "Test Envelope",
	}
	err := SaveEnvelope(s.tx.Exec)(e)
	s.NoError(err)
	var u, n, t string
	err = s.tx.QueryRow(`SELECT envelope_uuid, name, type FROM envelopes WHERE envelope_uuid = $1`, e.UUID).Scan(&u, &n, &t)
	if s.NoError(err) {
		s.Equal(e.UUID, u)
		s.Equal(e.Name, n)
		s.Equal(e.Type, t)
	}
}

func (s *EnvelopeRepositorySuite) Test_AllEnvelopes() {
	e1 := gledger.Envelope{
		UUID:    uuid.NewV4().String(),
		Name:    "Test Envelope 1",
		Type:    "income",
		Balance: 10,
	}
	e2 := gledger.Envelope{
		UUID: uuid.NewV4().String(),
		Name: "Test Envelope 2",
		Type: "expense",
	}
	s.mustExec(`INSERT INTO envelopes VALUES (
		$1, $2, now(), now(), $3
	)`, e1.UUID, e1.Name, e1.Type)
	s.mustExec(`INSERT INTO envelopes VALUES (
		$1, $2, now(), now(), $3
	)`, e2.UUID, e2.Name, e2.Type)
	s.mustExec(`INSERT INTO transactions VALUES ($1, $2, now(), 'payee', 10, 'f', 'f', now(), now(), $3)`, uuid.NewV4(), "cadd0722-6fd1-47ff-b390-b53307cc8c01", e1.UUID)

	es, err := AllEnvelopes(s.tx.Query)()
	if s.NoError(err) {
		s.Equal([]gledger.Envelope{e1, e2}, es)
	}
}

func (s *EnvelopeRepositorySuite) SetupSuite() {
	pg, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	s.Require().NoError(err)
	s.db = pg
}

func (s *EnvelopeRepositorySuite) SetupTest() {
	var err error
	s.tx, err = s.db.Begin()
	s.NoError(err)

	s.mustExec(`INSERT INTO accounts VALUES ($1, 'account name', 'type', 't', now(), now())`, "cadd0722-6fd1-47ff-b390-b53307cc8c01")
}

func (s *EnvelopeRepositorySuite) TearDownTest() {
	err := s.tx.Rollback()
	s.NoError(err)
}

func (s *EnvelopeRepositorySuite) mustExec(query string, args ...interface{}) {
	_, err := s.tx.Exec(query, args...)
	s.NoError(err)
}

func Test_EnvelopeRepository(t *testing.T) {
	suite.Run(t, new(EnvelopeRepositorySuite))
}
