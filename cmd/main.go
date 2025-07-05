package main

import (
	"log"
	"log/slog"
	"pill-reminder/configs"
	cronnotifier "pill-reminder/internal/cron-notifier"
	"pill-reminder/internal/db"
	"pill-reminder/internal/i18n"
	"pill-reminder/internal/logger"
	"pill-reminder/internal/repository"
	"pill-reminder/internal/service"
	tgbotapi "pill-reminder/internal/tgBotAPI"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	cfg := configs.InitConfig()
	logger.Init()
	mongo := db.Connect(db.ConnectMongoOptions{
		Uri:    cfg.MONGO_URL,
		DBName: cfg.MONGO_DB_NAME,
	})

	i18n.Init()

	pillDayService := service.NewPillDayService(repository.NewPillDayRepo(mongo))
	userService := service.NewUserService(repository.NewUserRepo(mongo))

	cronNotifier, err := cronnotifier.NewCronNotifier(
		cronnotifier.NotifierParams{
			Timezone:       cfg.TIMEZONE,
			PillDayService: pillDayService,
			Notifier:       nil,
		},
	)

	if err != nil {
		log.Fatalf("Failed to init CronNotifier: %v", err)
	}

	botAPI, err := tg.NewBotAPI(cfg.BOT_TOKEN)
	if err != nil {
		log.Fatalf("Failed to init bot: %v", err)
	}
	botAPI.Debug = false

	botService := tgbotapi.NewBotService(tgbotapi.BotServiceParams{
		Timezone:       cfg.TIMEZONE,
		API:            botAPI,
		UserService:    userService,
		PillDayService: pillDayService,
		CronNotifier:   cronNotifier,
	})

	cronNotifier.SetNotifier(botService)

	users, err := userService.GetAll()
	if err != nil {
		slog.Error(err.Error())
	}

	cronNotifier.Start(users)

	go botService.RegisterMessageListener()

	select {}
}
