package reminderqueue

type ReminderQueueService interface {
	CreateOrUpdate(chatId int64, cronId string, notificationType string) error
}
