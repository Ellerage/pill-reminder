package service

import (
	"errors"
	"pill-reminder/internal/model"
	"pill-reminder/internal/utils"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserService struct {
	userRepo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{userRepo: repo}
}

func (s *UserService) GetAll() ([]model.User, error) {
	return s.userRepo.GetAll()
}

func (s *UserService) GetByChatId(chatId int64) (*model.User, error) {
	return s.userRepo.GetByChatId(chatId)
}

func (s *UserService) Create(chatId int64, toCreate model.UserCreate) error {
	_, err := s.userRepo.GetByChatId(chatId)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return s.userRepo.Create(model.User{
			ChatId:         chatId,
			Timezone:       toCreate.Timezone,
			TimeToNotify:   toCreate.TimeToNotify,
			Status:         toCreate.Status,
			RemindInterval: utils.GetCronFromMinutes(toCreate.RemindInterval),
		})
	} else if err != nil {
		return err
	} else {
		return errors.New("user already exist")
	}
}

func (s *UserService) Update(chatId int64, toUpdate model.UserUpdate) error {
	return s.userRepo.Update(chatId, toUpdate)
}
