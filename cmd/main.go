package main

import (
	"log/slog"
	"pill-reminder/configs"
	cronnotifier "pill-reminder/internal/cron-notifier"
	"pill-reminder/internal/db"
	"pill-reminder/internal/repository"
	"pill-reminder/internal/service"
	tgbotapi "pill-reminder/internal/tgBotAPI"
)

func main() {
	cfg := configs.InitConfig()
	db := db.Connect(db.ConnectMongoOptions{Uri: cfg.MONGO_URL, DBName: cfg.MONGO_DB_NAME})
	tgbotapi.Init(cfg.BOT_TOKEN)

	pillDayRepo := repository.NewPillDayRepo(db)
	pillDayService := service.NewPillDayService(pillDayRepo)

	userRepo := repository.NewUserRepo(db)
	userService := service.NewUserService(userRepo)

	cronNotifier, err := cronnotifier.NewCronNotifier(
		cronnotifier.NotifierDeps{
			Config:         cfg,
			PillDayService: pillDayService,
			UserService:    userService,
		},
	)

	if err != nil {
		slog.Error(err.Error())
	}

	cronNotifier.Start()

	tgbotapi.RegisterMessageListener(
		tgbotapi.BotAPIDeps{
			Config:         cfg,
			PillDayService: pillDayService,
			UserService:    userService,
		},
	)

	select {}
}
