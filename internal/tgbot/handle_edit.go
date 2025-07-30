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

func (b *BotService) handleTimeEditing(message *tg.Message, userData HandleTimeEditing) error {
	idleStatus := string(enums.UserStatusIdle)
	var toUpdate = model.UserUpdate{Status: &idleStatus}
	timeToNotify := userData.UserTimeToNotify

	messageToSend := ""

	isTime := timeRegex.MatchString(message.Text)
	isTimezone := utils.IsValidTimezone(message.Text)
	minutes, intParseError := strconv.ParseInt(message.Text, 10, 64)
	isReminderInterval := intParseError == nil

	if !isReminderInterval && !isTime && !isTimezone {
		return utils.ErrInvalidTimeEditInput
	}

	var parseErr error
	if isReminderInterval {
		toUpdate.RemindInterval = &minutes
		messageToSend = i18n.GetText("repeatIntervalTimeUpdated")
	} else if isTime {
		timeToNotify, parseErr = utils.GetUTCFromUserTime(message.Text, *userData.UserTimezone)

		toUpdate.TimeToNotify = &timeToNotify
		messageToSend = i18n.GetText("firstAtDayNotificationTimeUpdated")
	} else if isTimezone {
		timeToNotify, parseErr = utils.GetUTCFromUserTime(timeToNotify, message.Text)

		toUpdate.TimeToNotify = &timeToNotify
		toUpdate.Timezone = &message.Text
		messageToSend = i18n.GetText("timezoneWasChanged")
	}

	if parseErr != nil {
		return parseErr
	}

	dailyCronId, followUpCronId, err := b.reminderService.GetCronIdByChatId(message.Chat.ID)
	if err != nil {
		return err
	}

	// TODO: add unregister by slice
	if dailyCronId != "" {
		if err := b.reminderQueue.Unregister(dailyCronId); err != nil {
			slog.Warn(err.Error())
		}
	}

	if followUpCronId != "" {
		if err := b.reminderQueue.Unregister(followUpCronId); err != nil {
			slog.Warn(err.Error())
		}
	}

	if err := b.userService.Update(message.Chat.ID, toUpdate); err != nil {
		return err
	}

	cronStr, err := utils.GetCronFromStringUTCTime(timeToNotify)
	if err != nil {
		return err
	}

	cronId, cronRegisterErr := b.reminderQueue.RegisterSchedule(cronStr, enums.ReminderEventDaily, model.DailyReminderPayload{ChatId: message.Chat.ID})
	if cronRegisterErr != nil {
		return cronRegisterErr
	}

	if err := b.reminderService.CreateOrUpdate(message.Chat.ID, cronId, enums.ReminderTypeDaily); err != nil {
		return err
	}

	err = b.SendMessage(message.Chat.ID, messageToSend, &enums.SendMessageButtons{Take: true, Edit: true}, nil)
	if err != nil {
		return err
	}

	return nil
}
