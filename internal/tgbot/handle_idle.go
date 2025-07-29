package tgbot

import (
	"errors"
	"fmt"
	"log/slog"
	"pill-reminder/internal/i18n"
	"pill-reminder/internal/model"
	"pill-reminder/internal/utils"
	"pill-reminder/internal/utils/enums"
	"strings"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/redis/go-redis/v9"
)

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

func (b *BotService) handleIdleMessages(message *tg.Message) error {
	if message.Text == string(enums.ActionTake) {
		err := b.pillDayService.MarkAsTakenNow(message.Chat.ID)
		if err != nil {
			return err
		}

		cleanReminderErr := b.cleanReminder(message.Chat.ID)
		if cleanReminderErr != nil {
			return cleanReminderErr
		}

		err = b.SendMessage(message.Chat.ID, i18n.GetText("checked"), &enums.SendMessageButtons{Take: true, Edit: true}, nil)
		if err != nil {
			return err
		}

		return nil
	}

	if message.Text == string(enums.ActionDelay) {
		cleanReminderErr := b.cleanReminder(message.Chat.ID)
		if cleanReminderErr != nil {
			slog.Error(cleanReminderErr.Error())
			return cleanReminderErr
		}

		_, err := b.reminderQueue.RegisterDelayed(message.Chat.ID)
		if err != nil {
			return err
		}

		err = b.SendMessage(message.Chat.ID, i18n.GetText("delayReminder"), &enums.SendMessageButtons{Take: true, Edit: true}, nil)
		if err != nil {
			return err
		}

		return nil
	}

	if message.Text == string(enums.ActionEdit) {
		status := string(enums.UserStatusEditing)
		err := b.userService.Update(message.Chat.ID, model.UserUpdate{Status: &status})
		// TODO: Should I clean up schedule on changing status to edit?
		if err != nil {
			return err
		}

		err = b.SendMessage(message.Chat.ID, i18n.GetText("enterNewTime"), &enums.SendMessageButtons{Take: true, Edit: true}, nil)
		if err != nil {
			return err
		}

		return nil
	}

	if message.Text == string(enums.ActionMySetting) {
		user, err := b.userService.GetByChatId(message.Chat.ID)
		if err != nil {
			return err
		}

		minutes := strings.TrimPrefix(strings.Split(user.RemindInterval, " ")[0], "*/") + " Minutes"

		text := fmt.Sprintf(
			"<b>User ID:</b> %d\n<b>Time to notify:</b> %s\n<b>Remind interval:</b> %s\n<b>Timezone:</b> %s",
			user.ChatId,
			user.TimeToNotify,
			minutes,
			user.Timezone,
		)
		parseMode := "HTML"

		err = b.SendMessage(message.Chat.ID, text, &enums.SendMessageButtons{Take: true, Edit: true}, &MessageOptions{ParseMode: &parseMode})
		if err != nil {
			return err
		}

		return nil
	}

	return utils.ErrInvalidCommand
}
