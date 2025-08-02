package tgbot

import (
	"pill-reminder/internal/utils/enums"
	"strings"
)

var titleInfo = []string{
	"<b>Available options:</b>",
}

var commandsInfo = []string{
	"/start - Create account",
	"Take - Mark day as taken",
	"Undo - Undo mark as taken today",
	"Edit - change notification settings",
	"Settings - show user settings",
}

func (b *BotService) SendInfoMessage(chatId int64) error {
	message := append(titleInfo, commandsInfo...)

	err := b.SendMessage(chatId, strings.Join(message, "\n"), &enums.SendMessageButtons{}, &MessageOptions{ParseMode: "HTML"})

	return err
}
