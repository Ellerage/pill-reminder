package tgbot

import (
	"log/slog"
	"pill-reminder/internal/i18n"
	"pill-reminder/internal/model"
)

func (b *BotService) handleUserCreate(chatId int64) {
	var u model.UserCreate

	err := b.userService.Create(chatId, u.GetDefaultUser(b.timezone))

	if err != nil {
		slog.Error(err.Error())
	} else {
		b.SendMessage(chatId, i18n.GetText("initialTime"), nil)
	}
}
