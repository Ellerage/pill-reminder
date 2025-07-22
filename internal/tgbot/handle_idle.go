package tgbot

import (
	"log/slog"
	"pill-reminder/internal/i18n"
	"pill-reminder/internal/model"
	"pill-reminder/internal/utils/enums"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *BotService) cleanReminder(chatId int64) error {
	cronId, err := b.reminderService.GetFollowupCronIdByChatId(chatId)
	if err != nil {
		return err
	}

	errRemoveByChatId := b.reminderQueue.Unregister(cronId)
	if errRemoveByChatId != nil {
		return errRemoveByChatId
	}

	_, deleteErr := b.reminderService.DeleteByChatId(chatId, true)
	if deleteErr != nil {
		return deleteErr
	}

	return nil
}

func (b *BotService) handleIdleMessages(message *tg.Message) {
	if message.Text == string(enums.ActionTake) {
		err := b.pillDayService.MarkAsTakenNow(message.Chat.ID)

		if err != nil {
			slog.Error(err.Error())
			b.SendMessage(message.Chat.ID, i18n.GetText("tryAgain"), &enums.SendMessageButtons{Take: true, Edit: true})
			return
		}

		cleanReminderErr := b.cleanReminder(message.Chat.ID)
		if cleanReminderErr != nil {
			slog.Error(cleanReminderErr.Error())
		}

		b.SendMessage(message.Chat.ID, i18n.GetText("checked"), &enums.SendMessageButtons{Take: true, Edit: true})
	}

	if message.Text == string(enums.ActionEdit) {
		status := string(enums.UserStatusEditing)
		err := b.userService.Update(message.Chat.ID, model.UserUpdate{Status: &status})
		// TODO: Should I clean up schedule on changing status to edit?
		if err != nil {
			slog.Error(err.Error())
		}

		b.SendMessage(message.Chat.ID, i18n.GetText("enterNewTime"), &enums.SendMessageButtons{Take: true, Edit: true})
	}

	if message.Text == string(enums.ActionDelay) {
		cleanReminderErr := b.cleanReminder(message.Chat.ID)
		if cleanReminderErr != nil {
			slog.Error(cleanReminderErr.Error())
		}

		_, err := b.reminderQueue.RegisterDelayed(message.Chat.ID)
		if err != nil {
			slog.Error(err.Error())
		}

		b.SendMessage(message.Chat.ID, i18n.GetText("delayReminder"), &enums.SendMessageButtons{Take: true, Edit: true})
	}
}
