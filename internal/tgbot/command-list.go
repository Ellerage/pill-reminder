package tgbot

import (
	"pill-reminder/internal/i18n"
	"pill-reminder/internal/utils/enums"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var commandsInfo = make([]tgbotapi.BotCommand, 0, 5)

func GetBotCommands() []tgbotapi.BotCommand {
	if len(commandsInfo) == 0 {
		commandsInfo = []tgbotapi.BotCommand{
			{Command: string(enums.ActionCommandCreate), Description: i18n.GetText("availableCommandsCreate")},
			{Command: string(enums.ActionCommandTake), Description: i18n.GetText("availableCommandsTake")},
			{Command: string(enums.ActionCommandUndo), Description: i18n.GetText("availableCommandsUndo")},
			{Command: string(enums.ActionCommandEdit), Description: i18n.GetText("availableCommandsEdit")},
			{Command: string(enums.ActionCommandSettings), Description: i18n.GetText("availableCommandsSettings")},
		}

		return commandsInfo
	}

	return commandsInfo
}
