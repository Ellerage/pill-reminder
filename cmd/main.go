package main

import (
	"pill-reminder/configs"
	cronnotifier "pill-reminder/internal/cron-notifier"
	"pill-reminder/internal/db"
	"pill-reminder/internal/repository"
	"pill-reminder/internal/service"
	tgbotapi "pill-reminder/internal/tgBotAPI"
)

func main() {
	cfg := configs.InitConfig()
	db := db.Connect(cfg.MONGO_URL)
	tgbotapi.Init(cfg.BOT_TOKEN)

	pillDayRepo := repository.NewPillDayRepo(db)
	pillDayService := service.NewPillDayService(pillDayRepo)

	cronnotifier.RegisterCronNotifier(
		cronnotifier.NotifierDeps{
			Config:         cfg,
			PillDayService: pillDayService,
		},
	)

	tgbotapi.RegisterMessageListener(
		tgbotapi.BotAPIDeps{
			Config:         cfg,
			PillDayService: pillDayService,
		},
	)

	select {}
}
