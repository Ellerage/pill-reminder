package service

import (
	"errors"
	"log/slog"
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

	_, err := s.pillDayRepo.GetByDateAndChatId(chatId, dateTime)

	var resultError error

	if errors.Is(err, mongo.ErrNoDocuments) {
		resultError = s.pillDayRepo.Create(chatId, &dateTime)
	} else if err == nil {
		resultError = s.pillDayRepo.UpdateTimeByDate(chatId, dateTime)
	} else {
		slog.Error(err.Error())
	}

	return resultError
}

func (s *PillDayService) IsTakenToday(chatId int64) (bool, error) {
	date := utils.GetNowDateTime()

	pillDay, err := s.pillDayRepo.GetByDateAndChatId(chatId, date)

	if errors.Is(err, mongo.ErrNoDocuments) {
		s.pillDayRepo.Create(chatId, nil)

		return false, nil
	} else {
		return pillDay.HasTimeOfTaking(), err
	}

}
