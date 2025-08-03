package utils

import "pill-reminder/internal/utils/enums"

func HandleCommandMessage(text string, command enums.ActionCommands) bool {
	return text == string(command) || command == enums.ActionToCommandMap[enums.Actions(text)]
}
