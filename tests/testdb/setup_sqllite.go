package testsdb

import (
	"pill-reminder/internal/db"

	"github.com/jmoiron/sqlx"
)

func SetupSQLite() (*sqlx.DB, func()) {
	dbClient, err := sqlx.Open("sqlite3", ":memory:")

	if err != nil {
		panic(err)
	}

	db.SetupSqlLiteTables(dbClient)

	return dbClient, func() { dbClient.Close() }
}
