package tgbot

import (
	"errors"
	"log/slog"
	"pill-reminder/internal/i18n"
	"pill-reminder/internal/model"
	"pill-reminder/internal/utils"
	"pill-reminder/internal/utils/enums"
	"time"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/redis/go-redis/v9"
)

func (b *BotService) handleIdleMessages(message *tg.Message) error {
	chatId := message.Chat.ID

	if utils.HandleCommandMessage(message.Text, enums.ActionCommandTake) {
		err := b.pillDayService.MarkAsTakenNow(chatId)
		if err != nil {
			return err
		}

		cleanReminderErr := b.cleanReminder(chatId)
		if cleanReminderErr != nil {
			return cleanReminderErr
		}

		err = b.SendMessage(chatId, i18n.GetText("checked"), &enums.SendMessageButtons{Take: true}, nil)
		if err != nil {
			return err
		}

		return nil
	}

	if utils.HandleCommandMessage(message.Text, enums.ActionCommandDelay) {
		cleanReminderErr := b.cleanReminder(chatId)
		if cleanReminderErr != nil {
			slog.Error(cleanReminderErr.Error())
			return cleanReminderErr
		}

		taskId, err := b.reminderQueue.RegisterDelayed(chatId)
		if err != nil {
			return err
		}

		err = b.reminderService.CreateOrUpdate(chatId, taskId, enums.ReminderTypeDelayed)
		if err != nil {
			return err
		}

		err = b.SendMessage(chatId, i18n.GetText("delayReminder"), &enums.SendMessageButtons{Take: true}, nil)
		if err != nil {
			return err
		}

		return nil
	}

	if utils.HandleCommandMessage(message.Text, enums.ActionCommandEdit) {
		status := string(enums.UserStatusEditing)
		err := b.userService.Update(chatId, model.UserUpdate{Status: &status})
		// TODO: Should I clean up schedule on changing status to edit?
		if err != nil {
			return err
		}

		err = b.SendMessage(chatId, i18n.GetText("enterNewTime"), &enums.SendMessageButtons{Take: true}, nil)
		if err != nil {
			return err
		}

		return nil
	}

	if utils.HandleCommandMessage(message.Text, enums.ActionCommandSettings) {
		user, err := b.userService.GetByChatId(chatId)
		if err != nil {
			return err
		}

		text := utils.GetSettingsReplyText(model.UserNotificationSettings{
			TimeToNotify:   user.TimeToNotify,
			Timezone:       user.Timezone,
			RemindInterval: user.RemindInterval,
		})

		err = b.SendMessage(chatId, text, &enums.SendMessageButtons{Take: true}, &MessageOptions{ParseMode: "HTML"})
		if err != nil {
			return err
		}

		return nil
	}

	if utils.HandleCommandMessage(message.Text, enums.ActionCommandUndo) {
		user, err := b.userService.GetByChatId(chatId)
		if err != nil {
			return err
		}

		now := utils.GetNowDateTime()
		userTime := utils.GetTimeFromString(user.TimeToNotify)

		if now.After(userTime) {
			// Register follow up
			taskId, err := b.reminderQueue.RegisterFollowup(chatId, time.Duration(user.RemindInterval)*time.Minute)

			if err != nil {
				return err
			}

			err = b.reminderService.CreateOrUpdate(chatId, taskId, enums.ReminderTypeFollowup)
			if err != nil {
				return err
			}
		}

		err = b.pillDayService.UndoAsTakenToday(chatId)
		if err != nil {
			return err
		}

		return b.SendMessage(chatId, i18n.GetText("undoneTaken"), nil, nil)
	}

	return utils.ErrInvalidCommand
}

func (b *BotService) cleanReminder(chatId int64) error {
	_, followUpTaskId, delayedTaskId, err := b.reminderService.GetCronIdByChatId(chatId)

	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}

	err = b.reminderQueue.Unregister(followUpTaskId, enums.ReminderTypeFollowup)
	if err != nil {
		slog.Warn(err.Error())
	}

	err = b.reminderQueue.Unregister(delayedTaskId, enums.ReminderTypeDelayed)
	if err != nil {
		slog.Warn(err.Error())
	}

	_, deleteErr := b.reminderService.DeleteByChatId(chatId, true)
	if deleteErr != nil {
		return deleteErr
	}

	return nil
}
