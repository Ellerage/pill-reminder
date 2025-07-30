package schedulehandlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"pill-reminder/internal/i18n"
	"pill-reminder/internal/model"
	"pill-reminder/internal/utils/enums"
	"time"

	"github.com/hibiken/asynq"
)

type DailyReminderHandler struct {
	reminderQueue        ReminderQueue
	reminderQueueService ReminderQueueService
	tgBot                TgBot
	pillDayService       PillDayService
	userService          UserService
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
			slog.Info("Pill was already taken today")
			return nil
		}

		user, err := deps.userService.GetByChatId(parsed.ChatId)
		if err != nil {
			return err
		}

		cronId, err := deps.reminderQueue.RegisterFollowup(parsed.ChatId, time.Duration(user.RemindInterval)*time.Minute)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		errCreating := deps.reminderQueueService.CreateOrUpdate(parsed.ChatId, cronId, enums.ReminderTypeFollowup)
		if errCreating != nil {
			return errCreating
		}

		err = deps.tgBot.SendMessage(parsed.ChatId, i18n.GetText("firstNotification"), &enums.SendMessageButtons{Edit: true, Take: true, Delay: true}, nil)
		if err != nil {
			return err
		}

		return nil
	}
}
