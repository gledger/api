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
		Name: "Test Envelope",
	}
	err := SaveEnvelope(s.tx.Exec)(e)
	s.NoError(err)
	var u, n string
	err = s.tx.QueryRow(`SELECT envelope_uuid, name FROM envelopes WHERE envelope_uuid = $1`, e.UUID).Scan(&u, &n)
	if s.NoError(err) {
		s.Equal(e.UUID, u)
		s.Equal(e.Name, n)
	}
}

func (s *EnvelopeRepositorySuite) Test_AllEnvelopes() {
	e1 := gledger.Envelope{
		UUID: uuid.NewV4().String(),
		Name: "Test Envelope 1",
	}
	e2 := gledger.Envelope{
		UUID: uuid.NewV4().String(),
		Name: "Test Envelope 2",
	}
	_, err := s.tx.Exec(`INSERT INTO envelopes VALUES (
		$1, $2, now(), now()
	)`, e1.UUID, e1.Name)
	s.NoError(err)
	_, err = s.tx.Exec(`INSERT INTO envelopes VALUES (
		$1, $2, now(), now()
	)`, e2.UUID, e2.Name)
	s.NoError(err)

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
}

func (s *EnvelopeRepositorySuite) TearDownTest() {
	err := s.tx.Rollback()
	s.NoError(err)
}

func Test_EnvelopeRepository(t *testing.T) {
	suite.Run(t, new(EnvelopeRepositorySuite))
}
