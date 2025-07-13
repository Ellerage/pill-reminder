package schedulehandlers

import (
	"log/slog"

	"github.com/hibiken/asynq"
)

type HandlersParams struct {
	Server               *asynq.Server
	Scheduler            *asynq.Scheduler
	ReminderQueueService ReminderQueueService
	TgBot                TgBot
	PillDayService       PillDayService
}

func RegisterHandlers(params HandlersParams) {
	mux := asynq.NewServeMux()

	mux.HandleFunc("reminder:daily", makeDailyReminderHandler(DailyReminderHandler{
		reminderQueueService: params.ReminderQueueService,
		scheduler:            params.Scheduler,
		tgBot:                params.TgBot,
		pillDayService:       params.PillDayService,
	}))

	mux.HandleFunc("reminder:followup", makeFollowupReminderHandle(FollowupHandler{
		tgBot:          params.TgBot,
		pillDayService: params.PillDayService,
	}))

	if err := params.Server.Run(mux); err != nil {
		slog.Error("asynq server failed", "err", err)
	}
}
