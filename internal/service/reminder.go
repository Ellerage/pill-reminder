package service

type ReminderQueueService struct {
	notificationQueueRepository ReminderQueueRepository
}

func NewReminderQueueService(repo ReminderQueueRepository) *ReminderQueueService {
	return &ReminderQueueService{notificationQueueRepository: repo}
}

func (s *ReminderQueueService) GetCronIdByChatId(chatId int64) (string, string, error) {
	return s.notificationQueueRepository.GetCronIdByChatId(chatId)
}

func (s *ReminderQueueService) GetFollowupCronIdByChatId(chatId int64) string {
	return s.notificationQueueRepository.GetFollowupCronIdByChatId(chatId)
}

func (s *ReminderQueueService) CreateOrUpdate(chatId int64, cronId string, notificationType string) error {
	return s.notificationQueueRepository.CreateOrUpdate(chatId, cronId, notificationType)
}

func (s *ReminderQueueService) DeleteByChatId(chatId int64, onlyFollowup bool) (int64, error) {
	return s.notificationQueueRepository.DeleteByChatId(chatId, onlyFollowup)
}
