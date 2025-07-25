package schedulehandlers

import (
	"pill-reminder/internal/model"
	"pill-reminder/internal/tgbot"
	"pill-reminder/internal/utils/enums"
)

type ReminderQueueService interface {
	CreateOrUpdate(chatId int64, cronId string, notificationType string) error
}

type TgBot interface {
	SendMessage(chatID int64, message string, buttons *enums.SendMessageButtons, options *tgbot.MessageOptions)
}

type PillDayService interface {
	IsTakenToday(chatId int64) (bool, error)
}

type ReminderQueue interface {
	RegisterSchedule(cronSpec string, taskType string, taskPayload any) (string, error)
}

type UserService interface {
	GetByChatId(chatId int64) (*model.User, error)
}
