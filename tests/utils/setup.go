package utils

import (
	"context"
	"pill-reminder/internal/i18n"
	reminderqueue "pill-reminder/internal/reminder-queue"
	"pill-reminder/internal/repository"
	"pill-reminder/internal/service"
	"pill-reminder/internal/tgbot"
	"pill-reminder/tests/mocks"
	testsdb "pill-reminder/tests/testdb"
	"testing"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Return struct {
	Bot                  *tgbot.BotService
	UserService          *service.UserService
	PillDayService       *service.PillDayService
	ReminderQueueService *service.ReminderQueueService
	ReminderQueue        *reminderqueue.ReminderQueue
	DB                   *mongo.Database
	Redis                *redis.Client
	BotAPI               *mocks.BotAPI
}

func FlushRedis(client *redis.Client) error {
	return client.FlushDB(context.Background()).Err()
}

func ClearMongo(db *mongo.Database) error {
	ctx := context.Background()
	collections, err := db.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		return err
	}

	for _, coll := range collections {
		if err := db.Collection(coll).Drop(ctx); err != nil {
			return err
		}
	}
	return nil
}

func Setup(t *testing.T) (Return, func()) {
	i18n.Init()

	userService := service.NewUserService(repository.NewUserRepo(testsdb.MongoClient))
	reminderQueueService := service.NewReminderQueueService(repository.NewQueueRepository(testsdb.RedisClient))
	pillDayService := service.NewPillDayService(repository.NewPillDayRepo(testsdb.MongoClient))

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

	reminderQueue.Scheduler.Ping()

	messageCh := make(chan tg.Chattable, 10)
	botApi := mocks.NewBotAPI(messageCh)
	bot := tgbot.NewBotService(tgbot.BotServiceParams{
		Timezone:        "UTC",
		API:             botApi,
		UserService:     userService,
		PillDayService:  pillDayService,
		ReminderService: reminderQueueService,
		ReminderQueue:   reminderQueue,
	})

	t.Log("Setup success!")

	return Return{
			Bot:                  bot,
			PillDayService:       pillDayService,
			UserService:          userService,
			ReminderQueueService: reminderQueueService,
			ReminderQueue:        reminderQueue,
			DB:                   testsdb.MongoClient,
			Redis:                testsdb.RedisClient,
			BotAPI:               botApi,
		}, func() {
			reminderQueue.Scheduler.Shutdown()
			reminderQueue.Server.Shutdown()
			reminderQueue.Client.Close()

			FlushRedis(testsdb.RedisClient)
			ClearMongo(testsdb.MongoClient)
		}
}
