package service

import (
	"errors"
	"pill-reminder/internal/model"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserRepository interface {
	GetAll() ([]model.User, error)
	GetByChatId(int64) (*model.User, error)
	Create(model.User) error
	Update(int64, model.UserUpdate) error
}

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

func (s *UserService) Create(toCreate model.User) error {
	_, err := s.userRepo.GetByChatId(toCreate.ChatId)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return s.userRepo.Create(toCreate)
	} else if err != nil {
		return err
	} else {
		return errors.New("user already exist")
	}
}

func (s *UserService) Update(chatId int64, toUpdate model.UserUpdate) error {
	return s.userRepo.Update(chatId, toUpdate)
}
