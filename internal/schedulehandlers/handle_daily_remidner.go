package schedulehandlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"pill-reminder/internal/i18n"
	"pill-reminder/internal/model"
	"pill-reminder/internal/utils/enums"

	"github.com/hibiken/asynq"
)

type DailyReminderHandler struct {
	reminderQueue        ReminderQueue
	reminderQueueService ReminderQueueService
	tgBot                TgBot
	pillDayService       PillDayService
}

func makeDailyReminderHandler(deps DailyReminderHandler) asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		var payload model.DailyReminderPayload

		if err := json.Unmarshal(t.Payload(), &payload); err != nil {
			slog.Error(err.Error())
		}

		isTaken, err := deps.pillDayService.IsTakenToday(payload.ChatId)

		if err != nil {
			slog.Error(err.Error())
		}

		if isTaken {
			return nil
		}

		cronId, err := deps.reminderQueue.RegisterSchedule(payload.RemindInterval, "reminder:followup", t.Payload())
		if err != nil {
			slog.Error(err.Error())
		}

		errCreating := deps.reminderQueueService.CreateOrUpdate(payload.ChatId, cronId, "Followup")

		if errCreating != nil {
			slog.Error(errCreating.Error())
		}

		deps.tgBot.SendMessage(payload.ChatId, i18n.GetText("firstNotification"), &enums.SendMessageButtons{Edit: true, Take: true, Delay: true}, nil)

		return nil
	}
}
