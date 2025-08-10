package seeds

import (
	"database/sql"
	"pill-reminder/internal/model"
	"pill-reminder/internal/utils"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/jmoiron/sqlx"
)

type PillDayParams struct {
	Date         *string
	TimeOfTaking *string
	ChatId       *int64
}

func PillDaySeed(db *sqlx.DB, initial *PillDayParams) model.PillDay {
	timeOfTaking := gofakeit.Date().Format("15:04")

	pillDay := model.PillDay{
		ChatId:       gofakeit.Int64(),
		Date:         gofakeit.Date().Format("2006-01-02"),
		TimeOfTaking: &timeOfTaking,
	}

	if initial != nil {
		if initial.Date != nil {
			pillDay.Date = *initial.Date
		}

		if initial.TimeOfTaking != nil {
			pillDay.TimeOfTaking = initial.TimeOfTaking
		}

		if initial.ChatId != nil {
			pillDay.ChatId = *initial.ChatId
		}
	}

	db.Exec("INSERT INTO pillDays (date, timeOfTaking, chatId) VALUES (?, ?, ?)", pillDay.Date, pillDay.TimeOfTaking, pillDay.ChatId)

	return pillDay
}

func FindPillDayByChatId(t *testing.T, db *sqlx.DB, chatId int64) (*model.PillDay, error) {
	row := db.QueryRow("SELECT chatId, timeOfTaking, date FROM pillDays WHERE chatId = ?", chatId)

	var result model.PillDay
	err := row.Scan(
		&result.ChatId,
		&result.TimeOfTaking,
		&result.Date,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.ErrNotFound
		}

		t.Fatal(err)
		return nil, err
	}

	return &result, nil
}
