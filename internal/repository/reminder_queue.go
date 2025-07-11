package repository

import (
	"context"
	"fmt"
	"log/slog"
	"pill-reminder/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ReminderQueueRepository struct {
	db *mongo.Database
}

func NewQueueRepository(db *mongo.Database) *ReminderQueueRepository {
	return &ReminderQueueRepository{db: db}
}

func (repo *ReminderQueueRepository) GetAll(filters *model.GetAllFilters) []model.QueueReminder {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var notificationsQueue []model.QueueReminder
	filterOptions := bson.M{}

	if filters.ChatId != nil {
		filterOptions["chatId"] = filters.ChatId
	}

	if filters.ReminderType != nil {
		filterOptions["ReminderType"] = filters.ReminderType
	}

	cursor, err := repo.db.Collection("notification-queue").Find(ctx, filterOptions)
	if err != nil {
		slog.Error(err.Error())
	}

	decodeErr := cursor.All(ctx, &notificationsQueue)

	if decodeErr != nil {
		slog.Error(err.Error())
	}

	return notificationsQueue
}

func (repo *ReminderQueueRepository) Create(chatId int64, cronId string, notificationType string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	queue := model.QueueReminder{ChatId: chatId, CronId: cronId, ReminderType: notificationType}

	_, err := repo.db.Collection("notification-queue").InsertOne(ctx, queue)

	if err != nil {
		slog.Error(err.Error())
	}

	return err
}

func (repo *ReminderQueueRepository) DeleteByChatId(chatId int64, filters model.DeleteFilters) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"chatId": chatId}

	if filters.ReminderType != nil {
		filter["ReminderType"] = filters.ReminderType
	}

	result, err := repo.db.Collection("notification-queue").DeleteMany(ctx, filter)

	if err != nil {
		slog.Error(err.Error())
		return 0, err
	}

	slog.Info(fmt.Sprintf("Removed crons for chat id: %d. Amount: %d", chatId, result.DeletedCount))

	return result.DeletedCount, nil
}
