package service

import "pill-reminder/internal/model"

type ReminderQueueService struct {
	notificationQueueRepository ReminderQueueRepository
}

func NewReminderQueueService(repo ReminderQueueRepository) *ReminderQueueService {
	return &ReminderQueueService{notificationQueueRepository: repo}
}

func (s *ReminderQueueService) GetAll(filters *model.GetAllFilters) []model.QueueReminder {
	return s.notificationQueueRepository.GetAll(filters)
}

func (s *ReminderQueueService) Create(chatId int64, cronId string, notificationType string) error {
	return s.notificationQueueRepository.Create(chatId, cronId, notificationType)
}

func (s *ReminderQueueService) DeleteByChatId(chatId int64, filters model.DeleteFilters) (int64, error) {
	return s.notificationQueueRepository.DeleteByChatId(chatId, filters)
}
