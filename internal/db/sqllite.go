package db

import (
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

func SetupSqlLite() (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite", "./pill-reminder.sqlite")
	if err != nil {
		return nil, err
	}

	SetupSqlLiteTables(db)

	return db, err
}

func SetupSqlLiteTables(db *sqlx.DB) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
			chatId integer NOT NULL,
			timezone text NOT NULL,
			timeToNotify text NOT NULL,
			status text NOT NULL,
			remindInterval integer NOT NULL
		)`)

	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS pillDays (
			date text NOT NULL,
			timeOfTaking text,
			ChatId integer NOT NULL
		)`)

	if err != nil {
		panic(err)
	}
}
