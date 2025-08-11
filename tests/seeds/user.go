package seeds

import (
	"database/sql"
	"errors"
	"log/slog"
	"pill-reminder/internal/model"
	"pill-reminder/internal/utils"
	"pill-reminder/internal/utils/enums"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/jmoiron/sqlx"
)

type UserParams struct {
	TimeToNotify          *string
	Status                *string
	RemindIntervalMinutes *int64
}

func UserSeed(db *sqlx.DB, initial *UserParams) model.User {
	remindInterval := int64(1)

	user := model.User{
		ChatId:         gofakeit.Int64(),
		Timezone:       "UTC",
		TimeToNotify:   gofakeit.Date().Format("15:04"),
		Status:         string(enums.UserStatusIdle),
		RemindInterval: remindInterval,
	}

	if initial != nil {
		if initial.RemindIntervalMinutes != nil {
			user.RemindInterval = *initial.RemindIntervalMinutes
		}

		if initial.TimeToNotify != nil {
			user.TimeToNotify = *initial.TimeToNotify
		}

		if initial.Status != nil {
			user.Status = *initial.Status
		}
	}

	_, err := db.Exec("INSERT INTO users (chatId, timezone, timeToNotify, status, remindInterval) VALUES (?, ?, ?, ?, ?)",
		user.ChatId,
		user.Timezone,
		user.TimeToNotify,
		user.Status,
		user.RemindInterval,
	)

	if err != nil {
		slog.Debug(err.Error())
	}

	return user
}

func GetUserByChatId(t *testing.T, db *sqlx.DB, chatId int64) (*model.User, error) {
	row := db.QueryRow("SELECT chatId, timezone, timeToNotify, status, remindInterval FROM users WHERE chatId = ?", chatId)

	var result model.User
	err := row.Scan(
		&result.ChatId,
		&result.Timezone,
		&result.TimeToNotify,
		&result.Status,
		&result.RemindInterval,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.ErrNotFound
		}

		t.Fatal(err)
		return nil, err
	}

	return &result, nil
}
