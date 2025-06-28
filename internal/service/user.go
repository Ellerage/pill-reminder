package service

import (
	"pill-reminder/internal/model"
	"pill-reminder/internal/repository"
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

func (s *UserService) Create(toCreate model.User) error {
	return s.userRepo.Create(toCreate)
}
