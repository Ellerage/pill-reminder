package utils

import (
	"log/slog"
	reminderqueue "pill-reminder/internal/reminder-queue"
	"pill-reminder/internal/service"
)

type StartScheduleHandlersParams struct {
	ReminderQueue *reminderqueue.ReminderQueue
	UserService   *service.UserService
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
}
