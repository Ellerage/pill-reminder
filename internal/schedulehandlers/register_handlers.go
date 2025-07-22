package schedulehandlers

import (
	"log/slog"

	"github.com/hibiken/asynq"
)

type HandlersParams struct {
	Server               *asynq.Server
	ReminderQueueService ReminderQueueService
	TgBot                TgBot
	PillDayService       PillDayService
	ReminderQueue        ReminderQueue
	UserService          UserService
}

func RegisterHandlers(params HandlersParams) error {
	mux := asynq.NewServeMux()

	mux.HandleFunc("reminder:daily", makeDailyReminderHandler(DailyReminderHandler{
		reminderQueueService: params.ReminderQueueService,
		reminderQueue:        params.ReminderQueue,
		tgBot:                params.TgBot,
		pillDayService:       params.PillDayService,
	}))

	mux.HandleFunc("reminder:followup", makeFollowupReminderHandle(FollowupHandler{
		tgBot:          params.TgBot,
		pillDayService: params.PillDayService,
	}))

	mux.HandleFunc("reminder:delayed", makeDelayedReminderHandler(DelayedReminderHandler{
		reminderQueueService: params.ReminderQueueService,
		reminderQueue:        params.ReminderQueue,
		tgBot:                params.TgBot,
		pillDayService:       params.PillDayService,
		userService:          params.UserService,
	}))

	err := params.Server.Run(mux)

	if err != nil {
		slog.Error("asynq server failed", "err", err)
		return err
	}

	return nil
}
