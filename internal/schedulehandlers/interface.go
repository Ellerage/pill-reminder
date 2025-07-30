package schedulehandlers

import (
	"pill-reminder/internal/model"
	"pill-reminder/internal/tgbot"
	"pill-reminder/internal/utils/enums"
	"time"
)

type ReminderQueueService interface {
	CreateOrUpdate(chatId int64, cronId string, notificationType enums.ReminderType) error
}

type TgBot interface {
	SendMessage(chatID int64, message string, buttons *enums.SendMessageButtons, options *tgbot.MessageOptions) error
}

type PillDayService interface {
	IsTakenToday(chatId int64) (bool, error)
}

type ReminderQueue interface {
	RegisterSchedule(cronSpec string, taskType enums.QueueEventsEnum, taskPayload any) (string, error)
	RegisterFollowup(chatId int64, interval time.Duration) (string, error)
}

type UserService interface {
	GetByChatId(chatId int64) (*model.User, error)
}
