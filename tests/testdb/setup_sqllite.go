package testsdb

import (
	"pill-reminder/internal/db"

	"github.com/jmoiron/sqlx"
)

var SqlLiteClient *sqlx.DB

func SetupSQLite() (*sqlx.DB, func()) {
	dbClient, err := sqlx.Open("sqlite3", ":memory:")

	if err != nil {
		panic(err)
	}

	db.SetupSqlLiteTables(dbClient)

	SqlLiteClient = dbClient

	return dbClient, func() { dbClient.Close() }
}
