package tgbot

import (
	"pill-reminder/internal/i18n"
	"pill-reminder/internal/model"
)

func (b *BotService) handleUserCreate(chatId int64) error {
	var u model.UserCreate

	err := b.userService.Create(chatId, u.GetDefaultUser(b.timezone))

	if err != nil {
		return err
	}

	err = b.SendMessage(chatId, i18n.GetText("initialTime"), nil, nil)
	if err != nil {
		return err
	}

	return nil
}
