package service

import (
	"errors"
	"pill-reminder/internal/model"
	"pill-reminder/internal/utils"
	"pill-reminder/internal/utils/enums"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) GetAll() ([]model.User, error) {
	args := m.Called()
	users, _ := args.Get(0).([]model.User)
	return users, args.Error(1)
}

func (m *MockUserRepo) GetByChatId(chatID int64) (*model.User, error) {
	args := m.Called(chatID)
	user, _ := args.Get(0).(*model.User)
	return user, args.Error(1)
}

func (m *MockUserRepo) Create(u model.User) error {
	args := m.Called(u)
	return args.Error(0)
}

func (m *MockUserRepo) Update(chatID int64, upd model.UserUpdate) error {
	args := m.Called(chatID, upd)
	return args.Error(0)
}

func TestUserService_GetAll_Success(t *testing.T) {
	mockRepo := new(MockUserRepo)
	svc := NewUserService(mockRepo)

	var fakeUsers []model.User
	for i := 0; i < 5; i++ {
		fakeUsers = append(fakeUsers, generateFakeUser())
	}

	mockRepo.
		On("GetAll").
		Return(fakeUsers, nil).
		Once()

	users, err := svc.GetAll()

	assert.NoError(t, err)
	assert.Equal(t, fakeUsers, users)

	mockRepo.AssertExpectations(t)
}

func TestUserService_GetByChatId_Success(t *testing.T) {
	mockRepo := new(MockUserRepo)
	svc := NewUserService(mockRepo)

	fakeUser := generateFakeUser()

	mockRepo.On("GetByChatId", fakeUser.ChatId).Return(&fakeUser, nil)

	user, err := svc.GetByChatId(fakeUser.ChatId)

	assert.NoError(t, err)
	assert.Equal(t, &fakeUser, user)

	mockRepo.AssertExpectations(t)
}

func TestUserService_Create_Success(t *testing.T) {
	mockRepo := new(MockUserRepo)
	svc := NewUserService(mockRepo)

	fakeUser := generateFakeUser()

	chatId := fakeUser.ChatId
	toCreate := model.UserCreate{
		Timezone:       fakeUser.Timezone,
		TimeToNotify:   fakeUser.TimeToNotify,
		Status:         fakeUser.Status,
		RemindInterval: gofakeit.Int64(),
	}

	expectedUser := model.User{
		ChatId:         chatId,
		Timezone:       toCreate.Timezone,
		TimeToNotify:   toCreate.TimeToNotify,
		Status:         toCreate.Status,
		RemindInterval: toCreate.RemindInterval,
	}

	mockRepo.
		On("GetByChatId", chatId).
		Return(nil, utils.ErrNotFound).
		Once()

	mockRepo.
		On("Create", expectedUser).
		Return(nil).
		Once()

	err := svc.Create(chatId, toCreate)

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestUserService_Update_Success(t *testing.T) {
	mockRepo := new(MockUserRepo)
	svc := NewUserService(mockRepo)

	fakeUserOld := generateFakeUser()
	timezone := gofakeit.TimeZone()
	timeToNotify := gofakeit.Date().Format("15:04")
	status := string(enums.UserStatusIdle)

	fakeUserToUpdate := model.UserUpdate{
		Timezone:     &timezone,
		TimeToNotify: &timeToNotify,
		Status:       &status,
	}

	mockRepo.On("Update", fakeUserOld.ChatId, fakeUserToUpdate).Return(nil).Once()

	err := svc.Update(fakeUserOld.ChatId, fakeUserToUpdate)

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

// ERROR CASES
func TestUserService_GetAll_Error(t *testing.T) {
	mockRepo := new(MockUserRepo)
	svc := NewUserService(mockRepo)

	expectedErr := errors.New("error")
	mockRepo.
		On("GetAll").
		Return(nil, expectedErr).
		Once()

	users, err := svc.GetAll()

	assert.Nil(t, users)
	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}

func TestUserService_GetByChatId_Error(t *testing.T) {
	mockRepo := new(MockUserRepo)
	svc := NewUserService(mockRepo)

	fakeUser := generateFakeUser()

	expectedErr := errors.New("error")

	mockRepo.
		On("GetByChatId", fakeUser.ChatId).
		Return(nil, expectedErr).
		Once()

	user, err := svc.GetByChatId(fakeUser.ChatId)

	assert.Nil(t, user)
	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}

func TestUserService_Create_Error(t *testing.T) {
	mockRepo := new(MockUserRepo)
	svc := NewUserService(mockRepo)

	fakeUser := generateFakeUser()

	mockRepo.On("GetByChatId", fakeUser.ChatId).Return(&fakeUser, nil).Once()

	toCreate := model.UserCreate{
		Timezone:       fakeUser.Timezone,
		TimeToNotify:   fakeUser.TimeToNotify,
		Status:         fakeUser.Status,
		RemindInterval: gofakeit.Int64(),
	}

	err := svc.Create(fakeUser.ChatId, toCreate)

	expectedError := errors.New("user already exist")

	assert.EqualError(t, err, expectedError.Error())

	mockRepo.AssertExpectations(t)
}

func TestUserService_Update_Error(t *testing.T) {
	mockRepo := new(MockUserRepo)
	svc := NewUserService(mockRepo)

	fakeUserOld := generateFakeUser()
	timezone := gofakeit.TimeZone()
	timeToNotify := gofakeit.Date().Format("15:04")
	status := string(enums.UserStatusIdle)

	fakeUserToUpdate := model.UserUpdate{
		Timezone:     &timezone,
		TimeToNotify: &timeToNotify,
		Status:       &status,
	}

	mockRepo.On("Update", fakeUserOld.ChatId, fakeUserToUpdate).Return(errors.New("error")).Once()

	err := svc.Update(fakeUserOld.ChatId, fakeUserToUpdate)

	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}

// UTILS
func generateFakeUser() model.User {
	return model.User{
		ChatId:       gofakeit.Int64(),
		Timezone:     gofakeit.TimeZone(),
		TimeToNotify: gofakeit.Date().Format("15:04"),
		Status:       string(enums.UserStatusIdle),
	}
}
