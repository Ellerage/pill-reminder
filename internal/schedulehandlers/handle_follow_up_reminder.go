package schedulehandlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"pill-reminder/internal/i18n"
	"pill-reminder/internal/model"
	"pill-reminder/internal/utils/enums"
	"time"

	"github.com/hibiken/asynq"
)

type FollowupHandler struct {
	tgBot                TgBot
	pillDayService       PillDayService
	reminderQueue        ReminderQueue
	reminderQueueService ReminderQueueService
	userService          UserService
}

func makeFollowupReminderHandle(deps FollowupHandler) asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		var parsed model.FollowUpReminderPayload

		if err := json.Unmarshal(t.Payload(), &parsed); err != nil {
			slog.Error(err.Error())
		}

		isTaken, err := deps.pillDayService.IsTakenToday(parsed.ChatId)
		if err != nil {
			return err
		}

		if isTaken {
			slog.Info("Pill was already taken")
			return nil
		}

		user, err := deps.userService.GetByChatId(parsed.ChatId)
		if err != nil {
			slog.Error(err.Error())
			return err
		}

		id, err := deps.reminderQueue.RegisterFollowup(parsed.ChatId, time.Duration(user.RemindInterval)*time.Minute)
		if err != nil {
			slog.Error(err.Error())
			return err
		}

		err = deps.reminderQueueService.CreateOrUpdate(parsed.ChatId, id, enums.ReminderTypeFollowup)
		if err != nil {
			slog.Error(err.Error())
			return nil
		}

		err = deps.tgBot.SendMessage(parsed.ChatId, i18n.GetText("reminderNotification"), &enums.SendMessageButtons{Take: true, Delay: true}, nil)
		if err != nil {
			slog.Error(err.Error())
			return err
		}

		return nil
	}
}
