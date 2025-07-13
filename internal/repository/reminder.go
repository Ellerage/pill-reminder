package repository

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

type ReminderQueueRepository struct {
	db *redis.Client
}

func NewQueueRepository(db *redis.Client) *ReminderQueueRepository {
	return &ReminderQueueRepository{db: db}
}

func (repo *ReminderQueueRepository) GetCronIdByChatId(chatId int64) (string, string, error) {
	result := repo.db.MGet(context.Background(), fmt.Sprintf("%d:Daily", chatId), fmt.Sprintf("%d:Followup", chatId))
	cronIds, err := result.Result()

	if err != nil {
		return "", "", err
	}

	var dailyCronID, followupCronID string

	if cronIds[0] != nil {
		dailyCronID = cronIds[0].(string)
	}

	if cronIds[1] != nil {
		followupCronID = cronIds[1].(string)
	}

	return dailyCronID, followupCronID, nil
}

func (repo *ReminderQueueRepository) GetFollowupCronIdByChatId(chatId int64) (string, error) {
	result, err := repo.db.Get(context.Background(), fmt.Sprintf("%d:Followup", chatId)).Result()

	if err != nil {
		return "", err
	}

	return result, nil
}

func (repo *ReminderQueueRepository) CreateOrUpdate(chatId int64, cronId string, notificationType string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	status := repo.db.Set(ctx, fmt.Sprintf("%d:%s", chatId, notificationType), cronId, 0)

	err := status.Err()

	if err != nil {
		return err
	}

	return nil
}

func (repo *ReminderQueueRepository) DeleteByChatId(chatId int64, onlyFollowUp bool) (int64, error) {
	toDelete := []string{fmt.Sprintf("%d:%s", chatId, "Followup")}

	if !onlyFollowUp {
		toDelete = append(toDelete, fmt.Sprintf("%d:%s", chatId, "Daily"))
	}

	result := repo.db.Del(context.TODO(), toDelete...)

	deleted, err := result.Result()

	if err != nil {
		return 0, err
	}

	slog.Info(fmt.Sprintf("Removed crons for chat id: %d. Amount: %d", chatId, deleted))

	return deleted, nil
}
