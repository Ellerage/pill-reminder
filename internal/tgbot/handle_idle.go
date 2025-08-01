package tgbot

import (
	"errors"
	"fmt"
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

	if message.Text == string(enums.ActionTake) {
		err := b.pillDayService.MarkAsTakenNow(chatId)
		if err != nil {
			return err
		}

		cleanReminderErr := b.cleanReminder(chatId)
		if cleanReminderErr != nil {
			return cleanReminderErr
		}

		err = b.SendMessage(chatId, i18n.GetText("checked"), &enums.SendMessageButtons{Take: true, Edit: true}, nil)
		if err != nil {
			return err
		}

		return nil
	}

	if message.Text == string(enums.ActionDelay) {
		cleanReminderErr := b.cleanReminder(chatId)
		if cleanReminderErr != nil {
			slog.Error(cleanReminderErr.Error())
			return cleanReminderErr
		}

		_, err := b.reminderQueue.RegisterDelayed(chatId)
		if err != nil {
			return err
		}

		err = b.SendMessage(chatId, i18n.GetText("delayReminder"), &enums.SendMessageButtons{Take: true, Edit: true}, nil)
		if err != nil {
			return err
		}

		return nil
	}

	if message.Text == string(enums.ActionEdit) {
		status := string(enums.UserStatusEditing)
		err := b.userService.Update(chatId, model.UserUpdate{Status: &status})
		// TODO: Should I clean up schedule on changing status to edit?
		if err != nil {
			return err
		}

		err = b.SendMessage(chatId, i18n.GetText("enterNewTime"), &enums.SendMessageButtons{Take: true, Edit: true}, nil)
		if err != nil {
			return err
		}

		return nil
	}

	if message.Text == string(enums.ActionMySetting) {
		user, err := b.userService.GetByChatId(chatId)
		if err != nil {
			return err
		}

		timeToNotify, err := utils.GetUserTimeFromUTC(user.TimeToNotify, user.Timezone)
		if err != nil {
			return err
		}

		text := fmt.Sprintf(
			"<b>User ID:</b> %d\n<b>Time to notify:</b> %s\n<b>Remind interval:</b> %s\n<b>Timezone:</b> %s",
			user.ChatId,
			timeToNotify,
			fmt.Sprintf("%d Minutes", user.RemindInterval),
			user.Timezone,
		)
		parseMode := "HTML"

		err = b.SendMessage(chatId, text, &enums.SendMessageButtons{Take: true, Edit: true}, &MessageOptions{ParseMode: &parseMode})
		if err != nil {
			return err
		}

		return nil
	}

	if message.Text == string(enums.ActionUndo) {
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
	cronId, err := b.reminderService.GetFollowupCronIdByChatId(chatId)

	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}

	errRemoveByChatId := b.reminderQueue.Unregister(cronId)
	if errRemoveByChatId != nil {
		slog.Warn(errRemoveByChatId.Error())
	}

	_, deleteErr := b.reminderService.DeleteByChatId(chatId, true)
	if deleteErr != nil {
		return deleteErr
	}

	return nil
}
