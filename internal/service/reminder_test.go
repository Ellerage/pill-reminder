package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockRepo is a testify mock for ReminderQueueRepository
type MockReminderRepo struct {
	mock.Mock
}

func (m *MockReminderRepo) GetCronIdByChatId(chatId int64) (string, string, error) {
	args := m.Called(chatId)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *MockReminderRepo) GetFollowupCronIdByChatId(chatId int64) (string, error) {
	args := m.Called(chatId)
	return args.String(0), args.Error(1)
}

func (m *MockReminderRepo) CreateOrUpdate(chatId int64, cronId string, notificationType string) error {
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
	r.On("GetCronIdByChatId", chatID).Return(expectedMain, expectedFollowup, nil)

	mainId, fupId, err := svc.GetCronIdByChatId(chatID)
	assert.NoError(t, err)
	assert.Equal(t, expectedMain, mainId)
	assert.Equal(t, expectedFollowup, fupId)

	r.AssertExpectations(t)
}

func TestGetCronIdByChatId_Error(t *testing.T) {
	r := new(MockReminderRepo)
	svc := NewReminderQueueService(r)

	chatID := int64(102)
	expectedErr := errors.New("db fail")
	r.On("GetCronIdByChatId", chatID).Return("", "", expectedErr)

	mainId, fupId, err := svc.GetCronIdByChatId(chatID)
	assert.Error(t, err)
	assert.Empty(t, mainId)
	assert.Empty(t, fupId)

	r.AssertExpectations(t)
}

func TestGetFollowupCronIdByChatId_Success(t *testing.T) {
	r := new(MockReminderRepo)
	svc := NewReminderQueueService(r)

	chatID := int64(201)
	expectedFup := "followup"
	r.On("GetFollowupCronIdByChatId", chatID).Return(expectedFup, nil)

	fupId, err := svc.GetFollowupCronIdByChatId(chatID)
	assert.NoError(t, err)
	assert.Equal(t, expectedFup, fupId)

	r.AssertExpectations(t)
}

func TestGetFollowupCronIdByChatId_Error(t *testing.T) {
	r := new(MockReminderRepo)
	svc := NewReminderQueueService(r)

	chatID := int64(202)
	expectedErr := errors.New("not found")
	r.On("GetFollowupCronIdByChatId", chatID).Return("", expectedErr)

	fupId, err := svc.GetFollowupCronIdByChatId(chatID)
	assert.Error(t, err)
	assert.Empty(t, fupId)

	r.AssertExpectations(t)
}

func TestCreateOrUpdate_Success(t *testing.T) {
	r := new(MockReminderRepo)
	svc := NewReminderQueueService(r)

	chatID := int64(301)
	cronID := "cron123"
	nType := "main"
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
	nType := "followup"
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
