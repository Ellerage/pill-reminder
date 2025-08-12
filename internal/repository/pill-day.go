package repository

import (
	"context"
	"database/sql"
	"pill-reminder/internal/model"
	"pill-reminder/internal/utils"
	"time"

	"github.com/jmoiron/sqlx"
)

type PillDayRepo struct {
	db *sqlx.DB
}

func NewPillDayRepo(db *sqlx.DB) *PillDayRepo {
	return &PillDayRepo{db: db}
}

func (repo *PillDayRepo) GetByDateAndChatId(chatId int64, date time.Time) (*model.PillDay, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	formattedDate := date.Format("2006-01-02")

	row := repo.db.QueryRowContext(ctx, `
		SELECT date, timeOfTaking, chatId
		FROM pillDays
		WHERE date = ? AND chatId = ?
	`, formattedDate, chatId)

	var result model.PillDay
	err := row.Scan(
		&result.Date,
		&result.TimeOfTaking,
		&result.ChatId,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.ErrNotFound
		}
		return nil, err
	}

	return &result, nil
}

func (repo *PillDayRepo) Create(chatId int64, timeOfTaking *time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var formattedTime *string

	if timeOfTaking != nil {
		str := timeOfTaking.Format("15:04")
		formattedTime = &str
	} else {
		formattedTime = nil
	}

	_, err := repo.db.ExecContext(ctx,
		"INSERT INTO pillDays (date, timeOfTaking, chatId) VALUES (?, ?, ?)", utils.GetFormattedNowDate(), formattedTime, chatId)

	return err
}

func (repo *PillDayRepo) UpdateTimeByDate(chatId int64, dateTime time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := repo.db.ExecContext(ctx,
		"UPDATE pillDays SET timeOfTaking = ? WHERE date = ? AND chatId = ?",
		dateTime.Format("15:04"), dateTime.Format("2006-01-02"), chatId)

	return err
}

func (repo *PillDayRepo) UnsetTodayByChatId(chatId int64) error {
	dateTime := utils.GetNowDateTime()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := repo.db.ExecContext(ctx,
		`UPDATE pillDays SET timeOfTaking = NULL WHERE date = ? AND chatId = ?`,
		dateTime.Format("2006-01-02"), chatId,
	)

	return err
}
