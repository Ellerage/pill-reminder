package utils

import (
	"context"
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

func Setup(t *testing.T) (Return, func()) {
	sqlLiteClient, closeSQLLiteDB := testsdb.SetupSQLite()
	teardownRedis := testsdb.SetupRedis()

	userService := service.NewUserService(repository.NewUserRepo(sqlLiteClient))
	reminderQueueService := service.NewReminderQueueService(repository.NewQueueRepository(testsdb.RedisClient))
	pillDayService := service.NewPillDayService(repository.NewPillDayRepo(sqlLiteClient))
	reminderQueue := reminderqueue.NewReminderQueue(
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

	botApi := mocks.NewBotAPI()
	bot := tgbot.NewBotService(tgbot.BotServiceParams{
		Timezone:        "UTC",
		API:             botApi,
		UserService:     userService,
		PillDayService:  pillDayService,
		ReminderService: reminderQueueService,
		ReminderQueue:   reminderQueue,
	})

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

	return Return{
			Bot:                  bot,
			PillDayService:       pillDayService,
			UserService:          userService,
			ReminderQueueService: reminderQueueService,
			ReminderQueue:        reminderQueue,
			DB:                   sqlLiteClient,
			Redis:                testsdb.RedisClient,
			BotAPI:               botApi,
		}, func() {
			CleanupDB()

			reminderQueue.Scheduler.Shutdown()
			time.Sleep(2 * time.Second)

			reminderQueue.Server.Stop()
			time.Sleep(1 * time.Second)
			reminderQueue.Server.Shutdown()
			time.Sleep(2 * time.Second)

			reminderQueue.Client.Close()

			closeSQLLiteDB()
			teardownRedis()
		}
}

func Init() func() {

	return func() {

	}
}

func CleanupDB() {
	FlushRedis(testsdb.RedisClient)
}

func FlushRedis(client *redis.Client) error {
	return client.FlushDB(context.Background()).Err()
}
