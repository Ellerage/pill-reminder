package reminderqueue

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"pill-reminder/internal/model"
	"pill-reminder/internal/utils"

	"github.com/hibiken/asynq"
)

type RedisConnectionOptions struct {
	RedisAddr string
	RedisPort int
	RedisPwd  string
	RedisDB   int
}

type ReminderQueueDeps struct {
	ReminderQueueService   ReminderQueueService
	RedisConnectionOptions RedisConnectionOptions
}

type ReminderQueue struct {
	Client    *asynq.Client
	Server    *asynq.Server
	Scheduler *asynq.Scheduler
	deps      ReminderQueueDeps
}

func NewReminderQueue(deps ReminderQueueDeps) *ReminderQueue {
	opt := asynq.RedisClientOpt{
		Addr:     fmt.Sprintf("%s:%d", deps.RedisConnectionOptions.RedisAddr, deps.RedisConnectionOptions.RedisPort),
		Password: deps.RedisConnectionOptions.RedisPwd,
		DB:       deps.RedisConnectionOptions.RedisDB,
	}

	client := asynq.NewClient(opt)
	server := asynq.NewServer(
		opt,
		asynq.Config{
			Concurrency: 10,
		},
	)
	serverPingErr := server.Ping()
	if serverPingErr != nil {
		slog.Error(serverPingErr.Error())
	}

	scheduler := asynq.NewScheduler(opt, nil)
	schedulerPingErr := scheduler.Ping()
	if schedulerPingErr != nil {
		slog.Error(schedulerPingErr.Error())
	}

	return &ReminderQueue{
		Client:    client,
		Server:    server,
		Scheduler: scheduler,
		deps:      deps,
	}
}

func (q *ReminderQueue) Start(users []model.User) error {
	var errs []error

	for _, user := range users {
		timeToNotify := utils.GetTimeFromString(user.TimeToNotify)
		cronStr := fmt.Sprintf("%d %d * * *", timeToNotify.Minute(), timeToNotify.Hour())

		data, err := json.Marshal(model.DailyReminderPayload{ChatId: user.ChatId, RemindInterval: user.RemindInterval})
		if err != nil {
			slog.Error(err.Error())
			errs = append(errs, err)
		}

		id, err := q.Scheduler.Register(cronStr, asynq.NewTask("reminder:daily", data))

		if err != nil {
			slog.Error(err.Error())
			errs = append(errs, err)
		}

		if err := q.deps.ReminderQueueService.CreateOrUpdate(user.ChatId, id, "Daily"); err != nil {
			slog.Error(err.Error())
			errs = append(errs, err)
		}

		slog.Info(fmt.Sprintf("ChatId: %d, CronId: %s, CronStr: %s", user.ChatId, id, cronStr))
	}

	return errors.Join(errs...)
}

func (q *ReminderQueue) Register(chatId int64, cron string, remindInterval string) (string, error) {
	data, err := json.Marshal(model.DailyReminderPayload{ChatId: chatId, RemindInterval: remindInterval})
	if err != nil {
		slog.Error(fmt.Sprintf("json.Marshal failed: %v", err))
	}

	id, registerError := q.Scheduler.Register(cron, asynq.NewTask("reminder:daily", data))
	if registerError != nil {
		return "", registerError
	}

	slog.Info(fmt.Sprintf("ChatId: %d, CronId: %s, CronStr: %s", chatId, id, cron))

	return id, registerError
}

func (q *ReminderQueue) Unregister(cronId string) error {
	err := q.Scheduler.Unregister(cronId)
	if err != nil {
		return err
	}

	slog.Info(fmt.Sprintf("Cron ID: %s - was removed", cronId))

	return nil
}
