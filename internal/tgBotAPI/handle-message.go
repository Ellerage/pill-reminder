package tgbotapi

import (
	"log/slog"
	"pill-reminder/internal/i18n"
	"pill-reminder/internal/model"
	"pill-reminder/internal/utils"
	"pill-reminder/internal/utils/enums"
	"regexp"
	"strconv"
	"time"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *BotService) handleUserCreate(chatId int64) {
	var u model.UserCreate

	err := b.userService.Create(chatId, u.GetDefaultUser(b.timezone))

	if err != nil {
		slog.Error(err.Error())
	} else {
		b.SendMessage(chatId, i18n.GetText("initialTime"))
	}
}

func (b *BotService) handleTimeEditing(message *tg.Message, timezone *string, userTimeToNotify string, userRepeatInterval string) {
	timeRegex := regexp.MustCompile(`^(?:[01]\d|2[0-3]):[0-5]\d$`)

	idleStatus := string(enums.UserStatusIdle)
	var toUpdate = model.UserUpdate{Status: &idleStatus}

	isTime := timeRegex.MatchString(message.Text)
	isTimezone := utils.IsValidTimezone(message.Text)
	minutes, parseErr := strconv.ParseUint(message.Text, 10, 8)

	timeToNotify, err := time.Parse("15:04", userTimeToNotify)

	if err != nil {
		slog.Error(err.Error())
	}

	if parseErr == nil {
		if minutes <= 60 {
			remindInterval := uint8(minutes)

			b.userService.Update(message.Chat.ID, toUpdate)
			b.cronNotifier.AddOrUpdateCron(message.Chat.ID, timeToNotify, utils.GetCronFromMinutes(remindInterval))
			b.SendMessage(message.Chat.ID, i18n.GetText("repeatIntervalTimeUpdated"))
			return
		} else {
			b.SendMessage(message.Chat.ID, i18n.GetText("shouldBeLessOneHour"))
			return
		}
	}

	if isTime || isTimezone {
		parsedTime := utils.GetTimeFromStringWithServerTimezone(message.Text, timezone)
		timeToNotify := parsedTime.Format("15:04")

		if isTime {
			toUpdate.TimeToNotify = &timeToNotify
		} else if isTimezone {
			toUpdate.Timezone = &message.Text
		}

		b.userService.Update(message.Chat.ID, toUpdate)
		b.cronNotifier.AddOrUpdateCron(message.Chat.ID, parsedTime, userRepeatInterval)

		b.SendMessage(message.Chat.ID, i18n.GetText("firstAtDayNotificationTimeUpdated"))
	} else {
		b.SendMessage(message.Chat.ID, i18n.GetText("notValidTime"))
	}
}

func (b *BotService) handleIdleMessages(message *tg.Message) {
	if message.Text == string(enums.ActionTake) {
		err := b.pillDayService.MarkAsTakenNow(message.Chat.ID)

		if err != nil {
			slog.Error(err.Error())
			b.SendMessage(message.Chat.ID, i18n.GetText("tryAgain"))
			return
		}

		b.SendMessage(message.Chat.ID, i18n.GetText("checked"))
	}

	if message.Text == string(enums.ActionEdit) {
		status := string(enums.UserStatusEditing)
		b.userService.Update(message.Chat.ID, model.UserUpdate{Status: &status})

		b.SendMessage(message.Chat.ID, i18n.GetText("enterNewTime"))
	}
}

func (b *BotService) handleMessage(message *tg.Message) {
	chatId := message.Chat.ID

	if message.Text == string(enums.ActionCreate) {
		b.handleUserCreate(chatId)
		return
	}

	user, err := b.userService.GetByChatId(chatId)

	if err != nil {
		slog.Error(err.Error())
	}

	if user.Status == string(enums.UserStatusEditing) || user.Status == string(enums.UserStatusInactive) {
		b.handleTimeEditing(message, &user.Timezone, user.TimeToNotify, user.RemindInterval)

		return
	}

	if user.Status == string(enums.UserStatusIdle) {
		b.handleIdleMessages(message)
		return
	}
}
