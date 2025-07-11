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
		err := b.pillDayService.MarkAsTakenNow(message.Chat.ID)

		if err != nil {
			slog.Error(err.Error())
			b.SendMessage(message.Chat.ID, i18n.GetText("tryAgain"), &enums.SendMessageButtons{Take: true, Edit: true})
			return
		}

		b.reminderQueue.RemoveByChatId(message.Chat.ID, true)

		b.SendMessage(message.Chat.ID, i18n.GetText("checked"), &enums.SendMessageButtons{Take: true, Edit: true})
	}

	if message.Text == string(enums.ActionEdit) {
		status := string(enums.UserStatusEditing)
		err := b.userService.Update(message.Chat.ID, model.UserUpdate{Status: &status})

		if err != nil {
			slog.Error(err.Error())
		}

		b.SendMessage(message.Chat.ID, i18n.GetText("enterNewTime"), &enums.SendMessageButtons{Take: true, Edit: true})
	}
}
