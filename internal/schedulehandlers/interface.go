package schedulehandlers

import "pill-reminder/internal/utils/enums"

type ReminderQueueService interface {
	Create(chatId int64, cronId string, notificationType string) error
}

type TgBot interface {
	SendMessage(chatID int64, message string, buttons *enums.SendMessageButtons)
}
