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

type FollowupHandler struct {
	tgBot          TgBot
	pillDayService PillDayService
}

func makeFollowupReminderHandle(deps FollowupHandler) asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		var payload model.FollowUpReminderPayload

		if err := json.Unmarshal(t.Payload(), &payload); err != nil {
			slog.Error(err.Error())
		}

		deps.tgBot.SendMessage(payload.ChatId, i18n.GetText("reminderNotification"), &enums.SendMessageButtons{Take: true, Edit: true})

		return nil
	}
}
