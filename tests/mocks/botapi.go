package mocks

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotAPI struct {
	SendCalls chan tg.Chattable
}

var messageCh = make(chan tg.Chattable, 10)

func NewBotAPI() *BotAPI {
	return &BotAPI{SendCalls: messageCh}
}

func (a *BotAPI) ClearMessages() {
	a.SendCalls = make(chan tg.Chattable, 10)
}

func (a *BotAPI) Send(v tg.Chattable) (tg.Message, error) {
	a.SendCalls <- v
	return tg.Message{}, nil
}

func (a *BotAPI) GetUpdatesChan(config tg.UpdateConfig) tg.UpdatesChannel {
	return make(tg.UpdatesChannel)
}
