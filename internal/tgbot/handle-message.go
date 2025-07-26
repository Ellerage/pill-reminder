package tgbot

import (
	"errors"
	"pill-reminder/internal/i18n"
	"pill-reminder/internal/utils/enums"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type HandleTimeEditing struct {
	UserTimezone       *string
	UserTimeToNotify   string
	UserRepeatInterval string
}

func (b *BotService) HandleMessage(message *tg.Message) error {
	chatId := message.Chat.ID

	if message.Text == string(enums.ActionCreate) {
		err := b.handleUserCreate(chatId)
		return err
	}

	user, err := b.userService.GetByChatId(chatId)

	if errors.Is(err, mongo.ErrNoDocuments) {
		b.SendMessage(message.Chat.ID, i18n.GetText("noAccount"), &enums.SendMessageButtons{Create: true}, nil)
		return nil
	}

	if err != nil {
		return err
	}

	if user.Status == string(enums.UserStatusEditing) || user.Status == string(enums.UserStatusInactive) {
		err := b.handleTimeEditing(message, HandleTimeEditing{UserTimezone: &user.Timezone, UserTimeToNotify: user.TimeToNotify, UserRepeatInterval: user.RemindInterval})

		return err
	}

	if user.Status == string(enums.UserStatusIdle) {
		err := b.handleIdleMessages(message)
		return err
	}

	return nil
}
