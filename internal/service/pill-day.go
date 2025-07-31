package service

import (
	"errors"
	"pill-reminder/internal/utils"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
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

	if errors.Is(err, mongo.ErrNoDocuments) {
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

	if errors.Is(err, mongo.ErrNoDocuments) {
		if err := s.pillDayRepo.Create(chatId, nil); err != nil {
			return false, err
		}

		return false, nil
	} else {
		return pillDay.HasTimeOfTaking(), err
	}

}
