package service

import (
	"errors"
	"pill-reminder/internal/utils/enums"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockRepo is a testify mock for ReminderQueueRepository
type MockReminderRepo struct {
	mock.Mock
}

func (m *MockReminderRepo) GetCronIdByChatId(chatId int64) (string, string, string, error) {
	args := m.Called(chatId)
	return args.String(0), args.String(1), args.String(2), args.Error(3)
}

func (m *MockReminderRepo) CreateOrUpdate(chatId int64, cronId string, notificationType enums.ReminderType) error {
	args := m.Called(chatId, cronId, notificationType)
	return args.Error(0)
}

func (m *MockReminderRepo) DeleteByChatId(chatId int64, onlyFollowup bool) (int64, error) {
	args := m.Called(chatId, onlyFollowup)
	return args.Get(0).(int64), args.Error(1)
}

func TestGetCronIdByChatId_Success(t *testing.T) {
	r := new(MockReminderRepo)
	svc := NewReminderQueueService(r)

	chatID := int64(101)
	expectedMain := "mainId"
	expectedFollowup := "fupId"
	expectedDelay := "delayId"
	r.On("GetCronIdByChatId", chatID).Return(expectedMain, expectedFollowup, expectedDelay, nil)

	mainId, fupId, delay, err := svc.GetCronIdByChatId(chatID)
	assert.NoError(t, err)
	assert.Equal(t, expectedMain, mainId)
	assert.Equal(t, expectedFollowup, fupId)
	assert.Equal(t, expectedDelay, delay)

	r.AssertExpectations(t)
}

func TestGetCronIdByChatId_Error(t *testing.T) {
	r := new(MockReminderRepo)
	svc := NewReminderQueueService(r)

	chatID := int64(102)
	expectedErr := errors.New("db fail")
	r.On("GetCronIdByChatId", chatID).Return("", "", "", expectedErr)

	mainId, fupId, delId, err := svc.GetCronIdByChatId(chatID)
	assert.Error(t, err)
	assert.Empty(t, mainId)
	assert.Empty(t, fupId)
	assert.Empty(t, delId)

	r.AssertExpectations(t)
}

func TestCreateOrUpdate_Success(t *testing.T) {
	r := new(MockReminderRepo)
	svc := NewReminderQueueService(r)

	chatID := int64(301)
	cronID := "cron123"
	nType := enums.ReminderTypeDaily
	r.On("CreateOrUpdate", chatID, cronID, nType).Return(nil)

	err := svc.CreateOrUpdate(chatID, cronID, nType)
	assert.NoError(t, err)

	r.AssertExpectations(t)
}

func TestCreateOrUpdate_Error(t *testing.T) {
	r := new(MockReminderRepo)
	svc := NewReminderQueueService(r)

	chatID := int64(302)
	cronID := "cronErr"
	nType := enums.ReminderTypeFollowup
	expectedErr := errors.New("create error")
	r.On("CreateOrUpdate", chatID, cronID, nType).Return(expectedErr)

	err := svc.CreateOrUpdate(chatID, cronID, nType)
	assert.Error(t, err)

	r.AssertExpectations(t)
}

func TestDeleteByChatId_Success(t *testing.T) {
	r := new(MockReminderRepo)
	svc := NewReminderQueueService(r)

	chatID := int64(401)
	records := int64(3)
	onlyFup := true
	r.On("DeleteByChatId", chatID, onlyFup).Return(records, nil)

	count, err := svc.DeleteByChatId(chatID, onlyFup)
	assert.NoError(t, err)
	assert.Equal(t, records, count)

	r.AssertExpectations(t)
}

func TestDeleteByChatId_Error(t *testing.T) {
	r := new(MockReminderRepo)
	svc := NewReminderQueueService(r)

	chatID := int64(402)
	onlyFup := false
	expectedErr := errors.New("delete failed")
	r.On("DeleteByChatId", chatID, onlyFup).Return(int64(0), expectedErr)

	count, err := svc.DeleteByChatId(chatID, onlyFup)
	assert.Error(t, err)
	assert.Zero(t, count)

	r.AssertExpectations(t)
}
