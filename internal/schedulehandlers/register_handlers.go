package schedulehandlers

import (
	"log/slog"
	"pill-reminder/internal/utils/enums"

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

	mux.HandleFunc(string(enums.ReminderEventDaily), makeDailyReminderHandler(DailyReminderHandler{
		reminderQueueService: params.ReminderQueueService,
		reminderQueue:        params.ReminderQueue,
		tgBot:                params.TgBot,
		pillDayService:       params.PillDayService,
		userService:          params.UserService,
	}))

	mux.HandleFunc(string(enums.ReminderEventFollowup), makeFollowupReminderHandle(FollowupHandler{
		tgBot:                params.TgBot,
		pillDayService:       params.PillDayService,
		reminderQueue:        params.ReminderQueue,
		reminderQueueService: params.ReminderQueueService,
		userService:          params.UserService,
	}))

	mux.HandleFunc(string(enums.ReminderEventDelayed), makeDelayedReminderHandler(DelayedReminderHandler{
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
