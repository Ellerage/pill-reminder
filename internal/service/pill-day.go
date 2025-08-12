package service

import (
	"errors"
	"pill-reminder/internal/utils"
	"time"
)

type PillDayService struct {
	pillDayRepo PillDayRepository
}

func NewPillDayService(repo PillDayRepository) *PillDayService {
	return &PillDayService{pillDayRepo: repo}
}

func (s *PillDayService) Create(chatId int64, timeOfTaking *time.Time) error {
	return s.pillDayRepo.Create(chatId, timeOfTaking)
}

func (s *PillDayService) MarkAsTakenNow(chatId int64) error {
	dateTime := utils.GetNowDateTime()

	pillDay, err := s.pillDayRepo.GetByDateAndChatId(chatId, dateTime)
	if errors.Is(err, utils.ErrNotFound) {
		return s.pillDayRepo.Create(chatId, &dateTime)
	}
	if err != nil {
		return err
	}

	if pillDay.HasTimeOfTaking() {
		return utils.ErrAlreadyTakenToday
	}

	return s.pillDayRepo.UpdateTimeByDate(chatId, dateTime)
}

func (s *PillDayService) UndoAsTakenToday(chatId int64) error {
	return s.pillDayRepo.UnsetTodayByChatId(chatId)
}

func (s *PillDayService) IsTakenToday(chatId int64) (bool, error) {
	date := utils.GetNowDateTime()

	pillDay, err := s.pillDayRepo.GetByDateAndChatId(chatId, date)
	if errors.Is(err, utils.ErrNotFound) {
		if err := s.pillDayRepo.Create(chatId, nil); err != nil {
			return false, err
		}
	}

	if err != nil {
		return false, err
	}

	return pillDay.HasTimeOfTaking(), err
}
