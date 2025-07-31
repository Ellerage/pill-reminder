package tgbot

import (
	"pill-reminder/internal/model"
	"pill-reminder/internal/utils/enums"
	"time"

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
	UndoAsTakenToday(chatId int64) error
}

type ReminderService interface {
	GetCronIdByChatId(chatId int64) (string, string, error)
	GetFollowupCronIdByChatId(chatId int64) (string, error)
	CreateOrUpdate(chatId int64, cronId string, notificationType enums.ReminderType) error
	DeleteByChatId(chatId int64, onlyFollowUp bool) (int64, error)
}

type ReminderQueue interface {
	RegisterSchedule(cronSpec string, taskType enums.QueueEventsEnum, taskPayload any) (string, error)
	RegisterDelayed(chatId int64) (string, error)
	RegisterFollowup(chatId int64, interval time.Duration) (string, error)
	Unregister(cronId string) error
}
