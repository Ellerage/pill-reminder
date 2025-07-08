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

type CronNotifier interface {
	AddOrUpdateCron(chatId int64, time string, repeatCronStr string) error
}
