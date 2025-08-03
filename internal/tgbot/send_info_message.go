package tgbot

import (
	"fmt"
	"pill-reminder/internal/i18n"
	"pill-reminder/internal/utils/enums"
	"strings"
)

func (b *BotService) SendInfoMessage(chatId int64) error {
	var commandsInfo = []string{}

	for _, command := range GetBotCommands() {
		commandsInfo = append(commandsInfo, fmt.Sprintf("%s - %s", command.Command, command.Description))
	}

	message := append([]string{i18n.GetText("availableCommandsTitle")}, commandsInfo...)

	err := b.SendMessage(chatId, strings.Join(message, "\n"), &enums.SendMessageButtons{}, &MessageOptions{ParseMode: "HTML"})

	return err
}
