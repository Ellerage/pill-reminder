package tgbot

import (
	"pill-reminder/internal/utils"
	"pill-reminder/internal/utils/enums"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type HandleTimeEditing struct {
	UserTimezone       *string
	UserTimeToNotify   string
	UserRepeatInterval int64
}

func (b *BotService) HandleMessage(message *tg.Message) error {
	chatId := message.Chat.ID

	if utils.HandleCommandMessage(message.Text, enums.ActionCommandCreate) {
		err := b.handleUserCreate(chatId)
		return err
	}

	user, err := b.userService.GetByChatId(chatId)
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
