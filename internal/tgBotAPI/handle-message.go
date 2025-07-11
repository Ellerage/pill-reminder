package tgbotapi

import (
	"fmt"
	"log/slog"
	"pill-reminder/internal/utils/enums"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type HandleTimeEditing struct {
	UserTimezone       *string
	UserTimeToNotify   string
	UserRepeatInterval string
}

func (b *BotService) handleMessage(message *tg.Message) {
	chatId := message.Chat.ID

	slog.Info(fmt.Sprintf("Got Message from: %d", chatId))

	if message.Text == string(enums.ActionCreate) {
		b.handleUserCreate(chatId)
		return
	}

	user, err := b.userService.GetByChatId(chatId)

	if err != nil {
		slog.Error(err.Error())
	}

	if user.Status == string(enums.UserStatusEditing) || user.Status == string(enums.UserStatusInactive) {
		b.handleTimeEditing(message, HandleTimeEditing{UserTimezone: &user.Timezone, UserTimeToNotify: user.TimeToNotify, UserRepeatInterval: user.RemindInterval})

		return
	}

	if user.Status == string(enums.UserStatusIdle) {
		b.handleIdleMessages(message)
		return
	}
}
