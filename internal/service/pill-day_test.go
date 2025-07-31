package service

import (
	"pill-reminder/internal/model"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MockPillDayRepo struct {
	mock.Mock
}

var ErrNoDocuments = mongo.ErrNoDocuments

func (m *MockPillDayRepo) Create(chatId int64, timeOfTaking *time.Time) error {
	args := m.Called(chatId, timeOfTaking)
	return args.Error(0)
}

func (m *MockPillDayRepo) MarkAsTakenNow(chatId int64) error {
	args := m.Called(chatId)
	return args.Error(0)
}

func (m *MockPillDayRepo) IsTakenToday(chatId int64) (bool, error) {
	args := m.Called(chatId)
	return args.Bool(0), args.Error(1)
}

func (m *MockPillDayRepo) GetByDateAndChatId(chatId int64, time time.Time) (*model.PillDay, error) {
	args := m.Called(chatId, time)

	return args.Get(0).(*model.PillDay), args.Error(1)
}

func (m *MockPillDayRepo) UpdateTimeByDate(chatId int64, time time.Time) error {
	args := m.Called(chatId, time)

	return args.Error(0)
}

func (m *MockPillDayRepo) UnsetTodayByChatId(chatId int64) error {
	args := m.Called(chatId)

	return args.Error(0)
}

func Test_Create_Success(t *testing.T) {
	mockRepo := new(MockPillDayRepo)
	service := NewPillDayService(mockRepo)

	chatId := int64(123)
	now := time.Now()

	mockRepo.On("Create", chatId, &now).Return(nil)

	err := service.Create(chatId, &now)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	mockRepo.AssertExpectations(t)
}

func Test_MarkAsTakenNow_CreateIfNoDocument(t *testing.T) {
	mockRepo := new(MockPillDayRepo)
	service := NewPillDayService(mockRepo)

	chatId := int64(123)

	mockRepo.On("GetByDateAndChatId", chatId, mock.AnythingOfType("time.Time")).Return((*model.PillDay)(nil), ErrNoDocuments)
	mockRepo.On("Create", chatId, mock.AnythingOfType("*time.Time")).Return(nil)

	err := service.MarkAsTakenNow(chatId)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	mockRepo.AssertExpectations(t)
}

func Test_MarkAsTakenNow_UpdateIfExists(t *testing.T) {
	mockRepo := new(MockPillDayRepo)
	service := NewPillDayService(mockRepo)

	chatId := int64(123)
	pillDay := &model.PillDay{}

	mockRepo.On("GetByDateAndChatId", chatId, mock.AnythingOfType("time.Time")).Return(pillDay, nil)
	mockRepo.On("UpdateTimeByDate", chatId, mock.AnythingOfType("time.Time")).Return(nil)

	err := service.MarkAsTakenNow(chatId)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	mockRepo.AssertExpectations(t)
}

func Test_IsTakenToday_NoDocument(t *testing.T) {
	mockRepo := new(MockPillDayRepo)
	service := NewPillDayService(mockRepo)

	chatId := int64(123)

	mockRepo.On("GetByDateAndChatId", chatId, mock.AnythingOfType("time.Time")).Return((*model.PillDay)(nil), ErrNoDocuments)
	mockRepo.On("Create", chatId, (*time.Time)(nil)).Return(nil)

	taken, err := service.IsTakenToday(chatId)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if taken {
		t.Errorf("expected not taken, got taken")
	}

	mockRepo.AssertExpectations(t)
}
