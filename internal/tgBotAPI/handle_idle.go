package tgbotapi

import (
	"log/slog"
	"pill-reminder/internal/i18n"
	"pill-reminder/internal/model"
	"pill-reminder/internal/utils/enums"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *BotService) handleIdleMessages(message *tg.Message) {
	if message.Text == string(enums.ActionTake) {
		errMarkAsTaken := b.pillDayService.MarkAsTakenNow(message.Chat.ID)

		if errMarkAsTaken != nil {
			slog.Error(errMarkAsTaken.Error())
			b.SendMessage(message.Chat.ID, i18n.GetText("tryAgain"), &enums.SendMessageButtons{Take: true, Edit: true})
			return
		}

		cronId := b.reminderService.GetFollowupCronIdByChatId(message.Chat.ID)

		errRemoveByChatId := b.reminderQueue.Unregister(cronId)
		if errRemoveByChatId != nil {
			slog.Error(errRemoveByChatId.Error())
		}

		_, deleteErr := b.reminderService.DeleteByChatId(message.Chat.ID, true)

		if deleteErr != nil {
			slog.Error(deleteErr.Error())
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
}
