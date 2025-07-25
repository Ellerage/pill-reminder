package reminderqueue

import "pill-reminder/internal/utils/enums"

type ReminderQueueService interface {
	CreateOrUpdate(chatId int64, cronId string, notificationType enums.ReminderType) error
}
