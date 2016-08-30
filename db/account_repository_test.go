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

type AccountRepositorySuite struct {
	suite.Suite

	db *sql.DB
	tx *sql.Tx
}

func (s *AccountRepositorySuite) Test_ReadAccount() {
	u := uuid.NewV4().String()
	s.mustExec(`INSERT INTO accounts VALUES ($1, 'name', 'type', 't', now(), now())`, u)

	a, err := ReadAccount(s.tx.QueryRow)(u)
	s.NoError(err)
	s.Equal(gledger.Account{
		UUID:    u,
		Name:    "name",
		Type:    "type",
		Active:  true,
		Balance: 0,
	}, a)
}

func (s *AccountRepositorySuite) Test_ReadAccount_WithTransactions() {
	u := uuid.NewV4().String()
	eu := uuid.NewV4().String()
	s.mustExec(`INSERT INTO accounts VALUES ($1, 'name', 'type', 't', now(), now())`, u)
	s.mustExec(`INSERT INTO envelopes VALUES ($1, 'envelope name', now(), now(), 'expense')`, eu)
	s.mustExec(`INSERT INTO transactions VALUES ($1, $2, now(), 'payee', 10, 'f', 'f', now(), now(), $3)`, uuid.NewV4(), u, eu)
	s.mustExec(`INSERT INTO transactions VALUES ($1, $2, now(), 'payee', -5, 'f', 'f', now(), now(), $3)`, uuid.NewV4(), u, eu)

	a, err := ReadAccount(s.tx.QueryRow)(u)
	s.NoError(err)
	s.Equal(gledger.Account{
		UUID:    u,
		Name:    "name",
		Type:    "type",
		Active:  true,
		Balance: 5,
	}, a)
}

func (s *AccountRepositorySuite) Test_AllAccounts() {
	u := uuid.NewV4().String()
	s.mustExec(`INSERT INTO accounts VALUES ($1, 'name', 'type', 't', now(), now())`, u)

	as, err := AllAccounts(s.tx.Query)()
	s.NoError(err)
	s.Equal([]gledger.Account{gledger.Account{
		UUID:    u,
		Name:    "name",
		Type:    "type",
		Active:  true,
		Balance: 0,
	}}, as)
}

func (s *AccountRepositorySuite) Test_AllAccounts_WithTransactions() {
	u := uuid.NewV4().String()
	eu := uuid.NewV4().String()
	s.mustExec(`INSERT INTO accounts VALUES ($1, 'name', 'type', 't', now(), now())`, u)
	s.mustExec(`INSERT INTO envelopes VALUES ($1, 'envelope name', now(), now(), 'expense')`, eu)
	s.mustExec(`INSERT INTO transactions VALUES ($1, $2, now(), 'payee', 10, 'f', 'f', now(), now(), $3)`, uuid.NewV4(), u, eu)
	s.mustExec(`INSERT INTO transactions VALUES ($1, $2, now(), 'payee', -5, 'f', 'f', now(), now(), $3)`, uuid.NewV4(), u, eu)

	as, err := AllAccounts(s.tx.Query)()
	s.NoError(err)
	s.Equal([]gledger.Account{gledger.Account{
		UUID:    u,
		Name:    "name",
		Type:    "type",
		Active:  true,
		Balance: 5,
	}}, as)
}

func (s *AccountRepositorySuite) SetupSuite() {
	pg, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	s.Require().NoError(err)
	s.db = pg
}

func (s *AccountRepositorySuite) SetupTest() {
	var err error
	s.tx, err = s.db.Begin()
	s.NoError(err)
}

func (s *AccountRepositorySuite) TearDownTest() {
	err := s.tx.Rollback()
	s.NoError(err)
}

func (s *AccountRepositorySuite) mustExec(query string, args ...interface{}) {
	_, err := s.tx.Exec(query, args...)
	s.NoError(err)
}

func Test_AccountRepository(t *testing.T) {
	suite.Run(t, new(AccountRepositorySuite))
}
