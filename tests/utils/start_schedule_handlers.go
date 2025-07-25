package utils

import (
	"log/slog"
	reminderqueue "pill-reminder/internal/reminder-queue"
	"pill-reminder/internal/schedulehandlers"
	"pill-reminder/internal/service"
	"pill-reminder/internal/tgbot"
)

type StartScheduleHandlersParams struct {
	ReminderQueue        *reminderqueue.ReminderQueue
	ReminderQueueService *service.ReminderQueueService
	Bot                  *tgbot.BotService
	PillDayService       *service.PillDayService
	UserService          *service.UserService
}

func InitScheduleForAllUsers(modules StartScheduleHandlersParams) {
	users, err := modules.UserService.GetAll()
	if err != nil {
		slog.Error(err.Error())
	}

	err = modules.ReminderQueue.Start(users)
	if err != nil {
		panic(err)
	}

	StartScheduleHandlers(StartScheduleHandlersParams{
		ReminderQueue:        modules.ReminderQueue,
		ReminderQueueService: modules.ReminderQueueService,
		Bot:                  modules.Bot,
		PillDayService:       modules.PillDayService,
		UserService:          modules.UserService,
	})
}

func StartScheduleHandlers(modules StartScheduleHandlersParams) {
	go func() {
		err := schedulehandlers.RegisterHandlers(
			schedulehandlers.HandlersParams{
				Server:               modules.ReminderQueue.Server,
				ReminderQueueService: modules.ReminderQueueService,
				TgBot:                modules.Bot,
				PillDayService:       modules.PillDayService,
				UserService:          modules.UserService,
				ReminderQueue:        modules.ReminderQueue,
			},
		)

		if err != nil {
			panic(err)
		}
	}()

	go func() {
		if err := modules.ReminderQueue.Scheduler.Run(); err != nil {
			panic(err)
		}
	}()
}
