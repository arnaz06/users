package mysql_test

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type mysqlSuite struct {
	suite.Suite
	db *sql.DB
	mg *migrate.Migrate
}

func (m *mysqlSuite) SetupSuite() {
	dsnDB := "users:users-pass@tcp(localhost:3306)/users?parseTime=1&loc=Asia%2FJakarta&charset=utf8mb4&collation=utf8mb4_unicode_ci"

	db, err := sql.Open("mysql", dsnDB)
	require.NoError(m.T(), err)
	require.NotNil(m.T(), db)
	m.mg, err = MigrateDB(db)
	require.NoError(m.T(), err)
	m.db = db
}

func (m *mysqlSuite) TearDownSuite() {
	require.NoError(m.T(), m.mg.Drop())
	require.NoError(m.T(), m.db.Close())
}
