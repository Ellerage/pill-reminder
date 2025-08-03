package repository

import (
	"context"
	"fmt"
	"log/slog"
	"pill-reminder/internal/utils/enums"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type ReminderQueueRepository struct {
	db *redis.Client
}

func NewQueueRepository(db *redis.Client) *ReminderQueueRepository {
	return &ReminderQueueRepository{db: db}
}

func (repo *ReminderQueueRepository) GetCronIdByChatId(chatId int64) (string, string, string, error) {
	result := repo.db.MGet(
		context.Background(),
		fmt.Sprintf("%d:%s", chatId, enums.ReminderTypeDaily),
		fmt.Sprintf("%d:%s", chatId, enums.ReminderTypeFollowup),
		fmt.Sprintf("%d:%s", chatId, enums.ReminderTypeDelayed),
	)
	cronIds, err := result.Result()
	if err != nil {
		return "", "", "", err
	}

	var dailyCronID, followupTaskID, delayedTaskId string

	if cronIds[0] != nil {
		dailyCronID = cronIds[0].(string)
	}
	if cronIds[1] != nil {
		followupTaskID = cronIds[1].(string)
	}
	if cronIds[2] != nil {
		delayedTaskId = cronIds[2].(string)
	}

	return dailyCronID, followupTaskID, delayedTaskId, nil
}

func (repo *ReminderQueueRepository) CreateOrUpdate(chatId int64, cronId string, notificationType enums.ReminderType) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	status := repo.db.Set(ctx, fmt.Sprintf("%d:%s", chatId, notificationType), cronId, time.Hour*24)

	err := status.Err()

	if err != nil {
		return err
	}

	return nil
}

func (repo *ReminderQueueRepository) DeleteByChatId(chatId int64, onlyFollowUp bool) (int64, error) {
	toDelete := []string{
		fmt.Sprintf("%d:%s", chatId, enums.ReminderTypeFollowup),
		fmt.Sprintf("%d:%s", chatId, enums.ReminderTypeDelayed),
	}

	if !onlyFollowUp {
		toDelete = append(toDelete, fmt.Sprintf("%d:%s", chatId, enums.ReminderTypeDaily))
	}

	result := repo.db.Del(context.TODO(), toDelete...)

	deleted, err := result.Result()
	if err != nil {
		return 0, err
	}

	deletedKeys := strings.Join(toDelete, " | ")
	slog.Info(fmt.Sprintf("Removed crons for chat id: %d. Amount: %d. Keys: %s", chatId, deleted, deletedKeys))

	return deleted, nil
}
