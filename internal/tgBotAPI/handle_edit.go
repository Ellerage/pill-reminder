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
	remindIntervalCron := userData.UserRepeatInterval
	messageToSend := ""

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
		remindIntervalCron = utils.GetCronFromMinutes(remindInterval)

		toUpdate.RemindInterval = &remindIntervalCron
		messageToSend = i18n.GetText("repeatIntervalTimeUpdated")
	} else if isTime {
		slog.Info(fmt.Sprintf("User - %d changes time to notify", message.Chat.ID))

		timeToNotify = utils.GetUTCFromUserTime(message.Text, userData.UserTimezone)
		toUpdate.TimeToNotify = &timeToNotify
		messageToSend = i18n.GetText("firstAtDayNotificationTimeUpdated")
	} else if isTimezone {
		slog.Info(fmt.Sprintf("User - %d changes timezone", message.Chat.ID))

		toUpdate.Timezone = &message.Text
		messageToSend = i18n.GetText("timezoneWasChanged")
	}

	dailyCronId, followUpCronId, err := b.reminderService.GetCronIdByChatId(message.Chat.ID)
	if err == nil {
		if dailyCronId != "" {
			b.reminderQueue.Unregister(dailyCronId)
		}

		if followUpCronId != "" {
			b.reminderQueue.Unregister(followUpCronId)
		}
	} else {
		slog.Error(err.Error())
	}

	b.userService.Update(message.Chat.ID, toUpdate)

	cronId, cronRegisterErr := b.reminderQueue.Register(message.Chat.ID, utils.GetDailyCronFromStringTime(timeToNotify), remindIntervalCron)

	if cronRegisterErr != nil {
		slog.Error(cronRegisterErr.Error())
	}

	b.reminderService.CreateOrUpdate(message.Chat.ID, cronId, "Daily")
	b.SendMessage(message.Chat.ID, messageToSend, &enums.SendMessageButtons{Take: true, Edit: true})
}
