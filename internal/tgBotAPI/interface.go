package tgbotapi

import tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type BotAPI interface {
	Send(tg.Chattable) (tg.Message, error)
	GetUpdatesChan(config tg.UpdateConfig) tg.UpdatesChannel
}
