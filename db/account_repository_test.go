package db

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/twinj/uuid"
	"github.com/zombor/gledger"

	"github.com/stretchr/testify/suite"
)

type AccountRepositorySuite struct {
	suite.Suite

	db *sql.DB
}

func (s *AccountRepositorySuite) SetupSuite() {
	pg, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	s.Require().NoError(err)
	s.db = pg
	_, err = s.db.Exec(`TRUNCATE accounts CASCADE`)
	s.Require().NoError(err)
}

func (s *AccountRepositorySuite) SetupTest() {
	_, err := s.db.Exec(`TRUNCATE accounts CASCADE`)
	s.Require().NoError(err)
}

func (s *AccountRepositorySuite) Test_ReadAccount() {
	u := uuid.NewV4().String()
	s.db.Exec(`INSERT INTO accounts VALUES ($1, 'name', 'type', 't', now(), now())`, u)

	a, err := ReadAccount(s.db)(u)
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
	s.db.Exec(`INSERT INTO accounts VALUES ($1, 'name', 'type', 't', now(), now())`, u)
	s.db.Exec(`INSERT INTO transactions VALUES ($1, $2, now(), 'payee', 10, 'f', 'f', now(), now())`, uuid.NewV4(), u)
	s.db.Exec(`INSERT INTO transactions VALUES ($1, $2, now(), 'payee', -5, 'f', 'f', now(), now())`, uuid.NewV4(), u)

	a, err := ReadAccount(s.db)(u)
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
	s.db.Exec(`INSERT INTO accounts VALUES ($1, 'name', 'type', 't', now(), now())`, u)

	as, err := AllAccounts(s.db)()
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
	s.db.Exec(`INSERT INTO accounts VALUES ($1, 'name', 'type', 't', now(), now())`, u)
	s.db.Exec(`INSERT INTO transactions VALUES ($1, $2, now(), 'payee', 10, 'f', 'f', now(), now())`, uuid.NewV4(), u)
	s.db.Exec(`INSERT INTO transactions VALUES ($1, $2, now(), 'payee', -5, 'f', 'f', now(), now())`, uuid.NewV4(), u)

	as, err := AllAccounts(s.db)()
	s.NoError(err)
	s.Equal([]gledger.Account{gledger.Account{
		UUID:    u,
		Name:    "name",
		Type:    "type",
		Active:  true,
		Balance: 5,
	}}, as)
}

func Test_AccountRepository(t *testing.T) {
	suite.Run(t, new(AccountRepositorySuite))
}
