package tgbot

import (
	"log/slog"
	"pill-reminder/internal/i18n"
	"pill-reminder/internal/model"
	"pill-reminder/internal/utils"
	"pill-reminder/internal/utils/enums"
	"regexp"
	"strconv"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var timeRegex = regexp.MustCompile(`^(?:[01]\d|2[0-3]):[0-5]\d$`)

func (b *BotService) handleTimeEditing(message *tg.Message, userData HandleTimeEditing) {
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
		timeToNotify, err := utils.GetUTCFromUserTime(message.Text, userData.UserTimezone)

		if err != nil {
			slog.Error(err.Error())
			return
		}

		toUpdate.TimeToNotify = &timeToNotify
		messageToSend = i18n.GetText("firstAtDayNotificationTimeUpdated")
	} else if isTimezone {
		timeToNotify, err := utils.GetUTCFromUserTime(userData.UserTimeToNotify, &message.Text)

		if err != nil {
			slog.Error(err.Error())
			return
		}

		toUpdate.TimeToNotify = &timeToNotify
		toUpdate.Timezone = &message.Text
		messageToSend = i18n.GetText("timezoneWasChanged")
	}

	dailyCronId, followUpCronId, err := b.reminderService.GetCronIdByChatId(message.Chat.ID)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	if dailyCronId != "" {
		if err := b.reminderQueue.Unregister(dailyCronId); err != nil {
			slog.Error(err.Error())
		}
	}

	if followUpCronId != "" {
		if err := b.reminderQueue.Unregister(followUpCronId); err != nil {
			slog.Error(err.Error())
		}
	}

	if err := b.userService.Update(message.Chat.ID, toUpdate); err != nil {
		slog.Error(err.Error())
		return
	}

	cronStr, err := utils.GetDailyCronFromStringTime(timeToNotify)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	cronId, cronRegisterErr := b.reminderQueue.Register(message.Chat.ID, cronStr, remindIntervalCron)

	if cronRegisterErr != nil {
		slog.Error(cronRegisterErr.Error())
		return
	}

	if err := b.reminderService.CreateOrUpdate(message.Chat.ID, cronId, "Daily"); err != nil {
		slog.Error(err.Error())
		return
	}

	b.SendMessage(message.Chat.ID, messageToSend, &enums.SendMessageButtons{Take: true, Edit: true})
}
