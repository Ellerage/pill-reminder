package schedulehandlers

import "pill-reminder/internal/utils/enums"

type ReminderQueueService interface {
	CreateOrUpdate(chatId int64, cronId string, notificationType string) error
}

type TgBot interface {
	SendMessage(chatID int64, message string, buttons *enums.SendMessageButtons)
}

type PillDayService interface {
	IsTakenToday(chatId int64) (bool, error)
}
