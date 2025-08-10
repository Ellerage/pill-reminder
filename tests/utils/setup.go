package utils

import (
	"context"
	"fmt"
	reminderqueue "pill-reminder/internal/reminder-queue"
	"pill-reminder/internal/repository"
	"pill-reminder/internal/schedulehandlers"
	"pill-reminder/internal/service"
	"pill-reminder/internal/tgbot"
	"pill-reminder/tests/mocks"
	testsdb "pill-reminder/tests/testdb"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type Return struct {
	Bot                  *tgbot.BotService
	UserService          *service.UserService
	PillDayService       *service.PillDayService
	ReminderQueueService *service.ReminderQueueService
	ReminderQueue        *reminderqueue.ReminderQueue
	DB                   *sqlx.DB
	Redis                *redis.Client
	BotAPI               *mocks.BotAPI
}

var userService *service.UserService
var reminderQueueService *service.ReminderQueueService
var pillDayService *service.PillDayService
var reminderQueue *reminderqueue.ReminderQueue
var bot *tgbot.BotService
var botApi *mocks.BotAPI

func Setup(t *testing.T) (Return, func()) {
	return Return{
			Bot:                  bot,
			PillDayService:       pillDayService,
			UserService:          userService,
			ReminderQueueService: reminderQueueService,
			ReminderQueue:        reminderQueue,
			DB:                   testsdb.SqlLiteClient,
			Redis:                testsdb.RedisClient,
			BotAPI:               botApi,
		}, func() {
			botApi.ClearMessages()
			CleanupDB()
		}
}

func Init() func() {
	fmt.Println("Init services")

	userService = service.NewUserService(repository.NewUserRepo(testsdb.SqlLiteClient))
	reminderQueueService = service.NewReminderQueueService(repository.NewQueueRepository(testsdb.RedisClient))
	pillDayService = service.NewPillDayService(repository.NewPillDayRepo(testsdb.SqlLiteClient))
	reminderQueue = reminderqueue.NewReminderQueue(
		reminderqueue.ReminderQueueDeps{
			ReminderQueueService: reminderQueueService,
			RedisConnectionOptions: reminderqueue.RedisConnectionOptions{
				RedisAddr: "localhost",
				RedisPort: testsdb.RedisPort,
				RedisPwd:  "",
				RedisDB:   0,
			},
		},
	)

	botApi = mocks.NewBotAPI()
	bot = tgbot.NewBotService(tgbot.BotServiceParams{
		Timezone:        "UTC",
		API:             botApi,
		UserService:     userService,
		PillDayService:  pillDayService,
		ReminderService: reminderQueueService,
		ReminderQueue:   reminderQueue,
	})

	reminderQueue.Scheduler.Ping()

	go func() {
		err := schedulehandlers.RegisterHandlers(
			schedulehandlers.HandlersParams{
				Server:               reminderQueue.Server,
				ReminderQueueService: reminderQueueService,
				TgBot:                bot,
				PillDayService:       pillDayService,
				UserService:          userService,
				ReminderQueue:        reminderQueue,
			},
		)

		if err != nil {
			panic(err)
		}
	}()

	go func() {
		if err := reminderQueue.Scheduler.Run(); err != nil {
			panic(err)
		}
	}()

	return func() {
		reminderQueue.Scheduler.Shutdown()
		time.Sleep(2 * time.Second)

		reminderQueue.Server.Shutdown()
		time.Sleep(2 * time.Second)

		reminderQueue.Client.Close()
	}
}

func CleanupDB() {
	FlushRedis(testsdb.RedisClient)
}

func FlushRedis(client *redis.Client) error {
	return client.FlushDB(context.Background()).Err()
}
