package schedulehandlers

import (
	"context"
	"encoding/json"
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
		var parsed model.DailyReminderPayload

		if err := json.Unmarshal(t.Payload(), &parsed); err != nil {
			return err
		}

		isTaken, err := deps.pillDayService.IsTakenToday(parsed.ChatId)
		if err != nil {
			return err
		}

		if isTaken {
			return nil
		}

		payload := model.FollowUpReminderPayload{ChatId: parsed.ChatId}
		cronId, err := deps.reminderQueue.RegisterSchedule(parsed.RemindInterval, enums.ReminderEventFollowup, payload)
		if err != nil {
			return err
		}

		errCreating := deps.reminderQueueService.CreateOrUpdate(parsed.ChatId, cronId, enums.ReminderTypeFollowup)

		if errCreating != nil {
			return errCreating
		}

		deps.tgBot.SendMessage(payload.ChatId, i18n.GetText("firstNotification"), &enums.SendMessageButtons{Edit: true, Take: true, Delay: true}, nil)

		return nil
	}
}
