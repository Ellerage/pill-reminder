package service

import (
	"pill-reminder/internal/model"
	"time"
)

type PillDayRepository interface {
	Create(chatId int64, time *time.Time) error
	GetByDateAndChatId(chatID int64, date time.Time) (*model.PillDay, error)
	UpdateTimeByDate(chatID int64, dateTime time.Time) error
}

type UserRepository interface {
	GetAll() ([]model.User, error)
	GetByChatId(int64) (*model.User, error)
	Create(model.User) error
	Update(int64, model.UserUpdate) error
}

type ReminderCronRepository interface {
}
