package tgbotapi

import (
	"fmt"
	"log/slog"
	"pill-reminder/internal/i18n"
	"pill-reminder/internal/model"
	"pill-reminder/internal/utils"
	"pill-reminder/internal/utils/enums"
	"regexp"
	"strconv"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *BotService) handleTimeEditing(message *tg.Message, userData HandleTimeEditing) {
	timeRegex := regexp.MustCompile(`^(?:[01]\d|2[0-3]):[0-5]\d$`)

	idleStatus := string(enums.UserStatusIdle)
	var toUpdate = model.UserUpdate{Status: &idleStatus}
	timeToNotify := userData.UserTimeToNotify

	isTime := timeRegex.MatchString(message.Text)
	isTimezone := utils.IsValidTimezone(message.Text)
	minutes, intParseError := strconv.ParseUint(message.Text, 10, 8)
	isMinutes := intParseError == nil

	if !isMinutes && !isTime && !isTimezone {
		b.SendMessage(message.Chat.ID, i18n.GetText("notValidTime"), nil)
		return
	}

	if isMinutes {
		remindInterval := uint8(minutes)
		cronStr := utils.GetCronFromMinutes(remindInterval)

		slog.Info(fmt.Sprintf("User - %d changes remind interval", message.Chat.ID))

		toUpdate.RemindInterval = &cronStr

		b.userService.Update(message.Chat.ID, toUpdate)

		b.reminderQueue.Register(message.Chat.ID, utils.GetDailyCronFromStringTime(userData.UserTimeToNotify), cronStr)
		b.SendMessage(message.Chat.ID, i18n.GetText("repeatIntervalTimeUpdated"), &enums.SendMessageButtons{Take: true, Edit: true})
		return
	}

	if isTime {
		slog.Info(fmt.Sprintf("User - %d changes time to notify", message.Chat.ID))

		timeToNotify = utils.GetUTCFromUserTime(message.Text, userData.UserTimezone)
		toUpdate.TimeToNotify = &timeToNotify
	}

	if isTimezone {
		slog.Info(fmt.Sprintf("User - %d changes timezone", message.Chat.ID))

		toUpdate.Timezone = &message.Text
	}

	b.userService.Update(message.Chat.ID, toUpdate)
	b.reminderQueue.Register(message.Chat.ID, utils.GetDailyCronFromStringTime(timeToNotify), userData.UserRepeatInterval)

	b.SendMessage(message.Chat.ID, i18n.GetText("firstAtDayNotificationTimeUpdated"), &enums.SendMessageButtons{Take: true, Edit: true})
}
