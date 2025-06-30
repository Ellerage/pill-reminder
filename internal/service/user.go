package service

import (
	"errors"
	"pill-reminder/internal/model"
	"pill-reminder/internal/repository"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserService struct {
	userRepo *repository.UserRepo
}

func NewUserService(repo *repository.UserRepo) *UserService {
	return &UserService{userRepo: repo}
}

func (s *UserService) GetAll() ([]model.User, error) {
	return s.userRepo.GetAll()
}

func (s *UserService) GetByChatId(chatId int64) (model.User, error) {
	return s.userRepo.GetByChatId(chatId)
}

func (s *UserService) Create(toCreate model.User) error {
	_, err := s.userRepo.GetByChatId(toCreate.ChatId)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return s.userRepo.Create(toCreate)
	}

	return err
}

func (s *UserService) Update(chatId int64, toUpdate model.UserUpdate) error {
	return s.userRepo.Update(chatId, toUpdate)
}
