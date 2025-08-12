package repository

import (
	"pill-reminder/internal/db"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func SetupSQLite(t *testing.T) (*sqlx.DB, func()) {
	t.Helper()

	dbClient, err := sqlx.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	db.SetupSqlLiteTables(dbClient)
	require.NoError(t, err)

	teardown := func() {
		_ = dbClient.Close()
	}

	return dbClient, teardown
}
