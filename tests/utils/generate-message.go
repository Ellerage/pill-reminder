package utils

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GenerateMessage(chatId int64, text string) *tgbotapi.Message {
	var message = &tgbotapi.Message{}

	if chatId != 0 {
		message.Chat = &tgbotapi.Chat{ID: chatId}
	}

	if text != "" {
		message.Text = text
	}

	return message
}
