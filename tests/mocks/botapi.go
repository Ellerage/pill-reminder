package mocks

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotAPI struct {
	SendCalls chan tg.Chattable
}

func NewBotAPI(messageCh chan tg.Chattable) *BotAPI {
	return &BotAPI{SendCalls: messageCh}
}

func (a *BotAPI) Send(v tg.Chattable) (tg.Message, error) {
	a.SendCalls <- v
	return tg.Message{}, nil
}

func (a *BotAPI) GetUpdatesChan(config tg.UpdateConfig) tg.UpdatesChannel {
	return make(tg.UpdatesChannel)
}
