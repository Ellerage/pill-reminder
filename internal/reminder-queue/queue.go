package reminderqueue

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"pill-reminder/internal/model"
	"pill-reminder/internal/utils"
	"pill-reminder/internal/utils/enums"
	"strconv"
	"time"

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
		Addr:     net.JoinHostPort(deps.RedisConnectionOptions.RedisAddr, strconv.Itoa(deps.RedisConnectionOptions.RedisPort)),
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
		panic(serverPingErr)
	}

	scheduler := asynq.NewScheduler(opt, nil)
	schedulerPingErr := scheduler.Ping()

	if schedulerPingErr != nil {
		panic(schedulerPingErr)
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

		data := model.DailyReminderPayload{ChatId: user.ChatId, RemindInterval: user.RemindInterval}

		id, err := q.RegisterSchedule(cronStr, enums.ReminderEventDaily, data)

		if err != nil {
			slog.Error(err.Error())
			errs = append(errs, err)
		}

		if err := q.deps.ReminderQueueService.CreateOrUpdate(user.ChatId, id, enums.ReminderTypeDaily); err != nil {
			slog.Error(err.Error())
			errs = append(errs, err)
		}

		slog.Info(fmt.Sprintf("ChatId: %d, CronId: %s, CronStr: %s", user.ChatId, id, cronStr))
	}

	return errors.Join(errs...)
}

func (q *ReminderQueue) RegisterSchedule(cronSpec string, taskType enums.QueueEventsEnum, taskPayload any) (string, error) {
	data, err := json.Marshal(taskPayload)
	if err != nil {
		return "", err
	}

	cronId, registerErr := q.Scheduler.Register(cronSpec, asynq.NewTask(string(taskType), data))
	if registerErr != nil {
		return "", registerErr
	}

	slog.Info(fmt.Sprintf("Register Schedule spec: %s. Task type %s. Id: %s", cronSpec, string(taskType), cronId))

	return cronId, nil
}

func (q *ReminderQueue) RegisterDelayed(chatId int64) (string, error) {
	payload := model.DelayedReminderPayload{
		ChatId: chatId,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	timeToWait := asynq.ProcessIn(60 * time.Minute)
	info, err := q.Client.Enqueue(asynq.NewTask(string(enums.ReminderEventDelayed), data), timeToWait)
	if err != nil {
		return "", err
	}

	return info.ID, err
}

func (q *ReminderQueue) Unregister(cronId string) error {
	err := q.Scheduler.Unregister(cronId)
	if err != nil {
		return err
	}

	slog.Info(fmt.Sprintf("Unregistered Schedule %s", cronId))

	return nil
}
