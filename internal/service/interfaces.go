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

type ReminderQueueRepository interface {
	GetAll(filters *model.GetAllFilters) []model.QueueReminder
	Create(chatId int64, cronId string, notificationType string) error
	DeleteByChatId(chatId int64, filters model.DeleteFilters) (int64, error)
}
