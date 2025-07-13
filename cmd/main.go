package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"pill-reminder/configs"
	"pill-reminder/internal/db"
	"pill-reminder/internal/i18n"
	"pill-reminder/internal/logger"
	reminderqueue "pill-reminder/internal/reminder-queue"
	"pill-reminder/internal/repository"
	"pill-reminder/internal/schedulehandlers"
	"pill-reminder/internal/service"
	tgbotapi "pill-reminder/internal/tgBotAPI"
	"syscall"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	ctx, close := context.WithCancel(context.Background())
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// System Init
	cfg := configs.InitConfig()
	logger.Init()
	i18n.Init()
	mongo := db.ConnectMongo(db.ConnectMongoOptions{
		Uri:    cfg.MONGO_URL,
		DBName: cfg.MONGO_DB_NAME,
	})

	redis := db.ConnectRedis(ctx, db.ConnectRedisOptions{
		Addr:     cfg.REDIS_URL,
		Port:     cfg.REDIS_PORT,
		Password: cfg.REDIS_PASSWORD,
		DB:       cfg.REMINDER_DB,
	})

	// Services
	pillDayService := service.NewPillDayService(repository.NewPillDayRepo(mongo))
	userService := service.NewUserService(repository.NewUserRepo(mongo))
	reminderQueueService := service.NewReminderQueueService(repository.NewQueueRepository(redis))

	reminderQueue := reminderqueue.NewReminderQueue(
		reminderqueue.ReminderQueueDeps{
			ReminderQueueService: reminderQueueService,
			RedisConnectionOptions: reminderqueue.RedisConnectionOptions{
				RedisAddr: cfg.REDIS_URL,
				RedisPort: cfg.REDIS_PORT,
				RedisPwd:  cfg.REDIS_PASSWORD,
				RedisDB:   cfg.ASYNCQ_DB,
			},
		},
	)

	err := reminderQueue.Scheduler.Ping()
	if err != nil {
		slog.Error(err.Error())
	}

	// TG bot API
	botAPI, err := tg.NewBotAPI(cfg.BOT_TOKEN)
	if err != nil {
		log.Fatalf("Failed to init bot: %v", err)
	}
	botAPI.Debug = false

	botService := tgbotapi.NewBotService(tgbotapi.BotServiceParams{
		Timezone:        cfg.TIMEZONE,
		API:             botAPI,
		UserService:     userService,
		PillDayService:  pillDayService,
		ReminderQueue:   reminderQueue,
		ReminderService: reminderQueueService,
	})

	// Initialize schedule for users
	users, err := userService.GetAll()
	if err != nil {
		slog.Error(err.Error())
	}

	reminderQueue.Start(users)

	go func() {
		if err := reminderQueue.Scheduler.Run(); err != nil {
			panic(err)
		}
	}()

	go func() {
		botService.RegisterMessageListener(ctx)
	}()

	// Schedule Event handlers
	go func() {
		schedulehandlers.RegisterHandlers(schedulehandlers.HandlersParams{Server: reminderQueue.Server, Scheduler: reminderQueue.Scheduler, ReminderQueueService: reminderQueueService, TgBot: botService})
	}()

	sig := <-sigCh
	slog.Info("signal: " + sig.String())

	defer close()
	reminderQueue.Scheduler.Shutdown()
	reminderQueue.Server.Shutdown()
	reminderQueue.Client.Close()
}
