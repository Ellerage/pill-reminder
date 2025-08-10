package repository

import (
	"context"
	"database/sql"
	"fmt"
	"pill-reminder/internal/model"
	"pill-reminder/internal/utils"
	"pill-reminder/internal/utils/enums"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (repo *UserRepo) GetAll() ([]model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rows, err := repo.db.QueryContext(ctx, `SELECT chatId, timeToNotify, timezone, status, remindInterval FROM users WHERE status != ?`, string(enums.UserStatusInactive))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]model.User, 0)

	for rows.Next() {
		var user model.User
		err := rows.Scan(
			&user.ChatId,
			&user.TimeToNotify,
			&user.Timezone,
			&user.Status,
			&user.RemindInterval,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (repo *UserRepo) GetByChatId(chatId int64) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	row := repo.db.QueryRowContext(ctx, `SELECT chatId, timeToNotify, timezone, status, remindInterval FROM users WHERE chatId = ?`, chatId)

	var user model.User
	err := row.Scan(
		&user.ChatId,
		&user.TimeToNotify,
		&user.Timezone,
		&user.Status,
		&user.RemindInterval,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (repo *UserRepo) Create(toCreate model.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := repo.db.ExecContext(ctx,
		`
			INSERT INTO users (chatId, timezone, timeToNotify, status, remindInterval)
			VALUES (?, ?, ?, ?, ?)
		`,
		toCreate.ChatId,
		toCreate.Timezone,
		toCreate.TimeToNotify,
		toCreate.Status,
		toCreate.RemindInterval,
	)

	return err
}

func (repo *UserRepo) Update(chatId int64, toUpdate model.UserUpdate) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	setClauses := make([]string, 0)
	args := make([]interface{}, 0)

	if toUpdate.Status != nil {
		setClauses = append(setClauses, "status = ?")
		args = append(args, *toUpdate.Status)
	}
	if toUpdate.Timezone != nil {
		setClauses = append(setClauses, "timezone = ?")
		args = append(args, *toUpdate.Timezone)
	}
	if toUpdate.TimeToNotify != nil {
		setClauses = append(setClauses, "timeToNotify = ?")
		args = append(args, *toUpdate.TimeToNotify)
	}
	if toUpdate.RemindInterval != nil {
		setClauses = append(setClauses, "remindInterval = ?")
		args = append(args, *toUpdate.RemindInterval)
	}

	if len(setClauses) == 0 {
		return nil
	}

	args = append(args, chatId)

	query := fmt.Sprintf("UPDATE users SET %s WHERE chatId = ?", strings.Join(setClauses, ", "))

	_, err := repo.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
