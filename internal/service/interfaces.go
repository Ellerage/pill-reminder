package service

import (
	"pill-reminder/internal/model"
	"pill-reminder/internal/utils/enums"
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

type ReminderQueueRepository interface {
	GetFollowupCronIdByChatId(chatId int64) (string, error)
	GetCronIdByChatId(chatId int64) (string, string, error)
	CreateOrUpdate(chatId int64, cronId string, notificationType enums.ReminderType) error
	DeleteByChatId(chatId int64, onlyFollowup bool) (int64, error)
}
