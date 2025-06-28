package service

import (
	"errors"
	"fmt"
	"pill-reminder/internal/model"
	"pill-reminder/internal/repository"
	"pill-reminder/internal/utils"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type PillDayService struct {
	pillDayRepo *repository.PillDayRepo
}

func NewPillDayService(repo *repository.PillDayRepo) *PillDayService {
	return &PillDayService{pillDayRepo: repo}
}

func (s *PillDayService) GetByDate(date string) (*model.PillDay, error) {
	return s.pillDayRepo.GetByDate(date)
}

func (s *PillDayService) Create(timeOfTaking *string) error {
	return s.pillDayRepo.Create(timeOfTaking)
}

func (s *PillDayService) UpdateTimeByDate(date string, newTime string) error {
	return s.pillDayRepo.UpdateTimeByDate(date, newTime)
}

func (s *PillDayService) MarkAsTakenNow() error {
	date := utils.GetNowDateTbilisi()
	time := utils.GetNowTimeTbilisi()

	_, err := s.pillDayRepo.GetByDate(date)

	var resultError error

	if errors.Is(err, mongo.ErrNoDocuments) {
		resultError = s.pillDayRepo.Create(&time)
	} else if err == nil {
		resultError = s.pillDayRepo.UpdateTimeByDate(date, time)
	} else {
		fmt.Println(err)
	}

	return resultError
}

func (s *PillDayService) IsTakenToday() (bool, error) {
	date := utils.GetNowDateTbilisi()

	pillDay, err := s.pillDayRepo.GetByDate(date)

	if errors.Is(err, mongo.ErrNoDocuments) {
		s.pillDayRepo.Create(nil)

		return false, err
	} else {
		return pillDay.HasTimeOfTaking(), err
	}

}
