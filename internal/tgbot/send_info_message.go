package tgbot

import (
	"fmt"
	"pill-reminder/internal/i18n"
	"pill-reminder/internal/utils/enums"
	"strings"
)

func (b *BotService) SendInfoMessage(chatId int64) error {
	var commandsInfo = []string{
		fmt.Sprintf("%s - %s", enums.ActionCommandCreate, i18n.GetText("availableCommandsCreate")),
		fmt.Sprintf("%s - %s", enums.ActionCommandTake, i18n.GetText("availableCommandsTake")),
		fmt.Sprintf("%s - %s", enums.ActionCommandUndo, i18n.GetText("availableCommandsUndo")),
		fmt.Sprintf("%s - %s", enums.ActionCommandEdit, i18n.GetText("availableCommandsEdit")),
		fmt.Sprintf("%s - %s", enums.ActionCommandSettings, i18n.GetText("availableCommandsSettings")),
	}

	message := append([]string{i18n.GetText("availableCommandsTitle")}, commandsInfo...)

	err := b.SendMessage(chatId, strings.Join(message, "\n"), &enums.SendMessageButtons{}, &MessageOptions{ParseMode: "HTML"})

	return err
}
