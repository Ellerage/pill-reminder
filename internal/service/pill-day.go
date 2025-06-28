package service

import (
	"errors"
	"fmt"
	"pill-reminder/internal/model"
	"pill-reminder/internal/repository"
	"pill-reminder/internal/utils"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type PillDayService struct {
	pillDayRepo *repository.PillDayRepo
}

func NewPillDayService(repo *repository.PillDayRepo) *PillDayService {
	return &PillDayService{pillDayRepo: repo}
}

func (s *PillDayService) GetByDate(date time.Time) (*model.PillDay, error) {
	return s.pillDayRepo.GetByDate(date)
}

func (s *PillDayService) Create(timeOfTaking time.Time) error {
	return s.pillDayRepo.Create(&timeOfTaking)
}

func (s *PillDayService) UpdateTimeByDate(date time.Time, newTime time.Time) error {
	return s.pillDayRepo.UpdateTimeByDate(date, newTime)
}

func (s *PillDayService) MarkAsTakenNow() error {
	dateTime := utils.GetNowDateTime(nil)

	_, err := s.pillDayRepo.GetByDate(dateTime)

	var resultError error

	if errors.Is(err, mongo.ErrNoDocuments) {
		resultError = s.pillDayRepo.Create(&dateTime)
	} else if err == nil {
		resultError = s.pillDayRepo.UpdateTimeByDate(dateTime, dateTime)
	} else {
		fmt.Println(err)
	}

	return resultError
}

func (s *PillDayService) IsTakenToday() (bool, error) {
	date := utils.GetNowDateTime(nil)

	pillDay, err := s.pillDayRepo.GetByDate(date)

	if errors.Is(err, mongo.ErrNoDocuments) {
		s.pillDayRepo.Create(nil)

		return false, err
	} else {
		return pillDay.HasTimeOfTaking(), err
	}

}
