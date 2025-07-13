package tgbotapi

import (
	"pill-reminder/internal/model"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotAPI interface {
	Send(tg.Chattable) (tg.Message, error)
	GetUpdatesChan(config tg.UpdateConfig) tg.UpdatesChannel
}

type UserService interface {
	GetByChatId(chatId int64) (*model.User, error)
	Create(chatId int64, toCreate model.UserCreate) error
	Update(chatId int64, toUpdate model.UserUpdate) error
}

type PillDayService interface {
	MarkAsTakenNow(chatId int64) error
}

type ReminderService interface {
	GetCronIdByChatId(chatId int64) (string, string, error)
	GetFollowupCronIdByChatId(chatId int64) string
	CreateOrUpdate(chatId int64, cronId string, notificationType string) error
	DeleteByChatId(chatId int64, onlyFollowUp bool) (int64, error)
}

type ReminderQueue interface {
	Register(chatId int64, cron string, repeatCronStr string) (string, error)
	Unregister(cronId string) error
}
