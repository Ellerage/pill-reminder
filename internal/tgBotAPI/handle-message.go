package tgbotapi

import (
	"log/slog"
	"pill-reminder/internal/model"
	"pill-reminder/internal/utils"
	"pill-reminder/internal/utils/enums"
	"regexp"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *BotService) handleUserCreate(chatId int64) {
	err := b.userService.Create(model.User{
		ChatId:       chatId,
		Timezone:     b.timezone,
		TimeToNotify: "00:00",
		Status:       string(enums.UserStatusInactive),
	})

	if err != nil {
		slog.Error(err.Error())
	} else {
		b.SendMessage(chatId, "What's time you want to get reminders? Type it in 15:04 format")
	}
}

func (b *BotService) handleTimeEditing(message *tg.Message, timezone *string) {
	timeRegex := regexp.MustCompile(`^(?:[01]\d|2[0-3]):[0-5]\d$`)

	idleStatus := string(enums.UserStatusIdle)

	isTime := timeRegex.MatchString(message.Text)
	isTimezone := utils.IsValidTimezone(message.Text)

	if isTime || isTimezone {
		var toUpdate = model.UserUpdate{Status: &idleStatus}

		parsedTime := utils.GetTimeFromStringWithServerTimezone(message.Text, timezone)
		timeToNotify := parsedTime.Format("15:04")

		if isTime {
			toUpdate.TimeToNotify = &timeToNotify
		} else if isTimezone {
			toUpdate.Timezone = &message.Text
		}

		b.userService.Update(message.Chat.ID, toUpdate)
		b.cronNotifier.AddOrUpdateCron(message.Chat.ID, parsedTime)

		b.SendMessage(message.Chat.ID, "Time was updated!")
	} else {
		b.SendMessage(message.Chat.ID, "Not valid time or timezone")
	}
}

func (b *BotService) handleIdleMessages(message *tg.Message) {
	if message.Text == string(enums.ActionTake) {
		err := b.pillDayService.MarkAsTakenNow(message.Chat.ID)

		if err != nil {
			slog.Error(err.Error())
			b.SendMessage(message.Chat.ID, "Try again")
			return
		}

		b.SendMessage(message.Chat.ID, "Checked!")
	}

	if message.Text == string(enums.ActionEdit) {
		status := string(enums.UserStatusEditing)
		b.userService.Update(message.Chat.ID, model.UserUpdate{Status: &status})

		b.SendMessage(message.Chat.ID, "Enter new time to get notified - 15:04 format")
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
		b.handleTimeEditing(message, &user.Timezone)

		return
	}

	if user.Status == string(enums.UserStatusIdle) {
		b.handleIdleMessages(message)
		return
	}
}
