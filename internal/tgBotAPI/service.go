package tgbotapi

import (
	"log"
	"log/slog"
	cronnotifier "pill-reminder/internal/cron-notifier"
	"pill-reminder/internal/model"
	"pill-reminder/internal/service"
	"pill-reminder/internal/utils"
	"pill-reminder/internal/utils/enums"
	"regexp"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotAPI interface {
	Send(tg.Chattable) (tg.Message, error)
	GetUpdatesChan(config tg.UpdateConfig) tg.UpdatesChannel
}

type BotService struct {
	timezone       string
	api            BotAPI
	userService    *service.UserService
	pillDayService *service.PillDayService
	cronNotifier   *cronnotifier.NotifierService
}

type BotServiceParams struct {
	Timezone       string
	API            BotAPI
	UserService    *service.UserService
	PillDayService *service.PillDayService
	CronNotifier   *cronnotifier.NotifierService
}

func NewBotService(params BotServiceParams) *BotService {
	return &BotService{
		timezone:       params.Timezone,
		api:            params.API,
		userService:    params.UserService,
		pillDayService: params.PillDayService,
		cronNotifier:   params.CronNotifier,
	}
}

func (b *BotService) RegisterMessageListener() {
	u := tg.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	log.Println("Listening for new messages...")

	for update := range updates {
		if update.Message != nil {
			b.handleMessage(update.Message)
		}
	}
}

func (b *BotService) SendMessage(chatID int64, message string) {
	msg := tg.NewMessage(chatID, message)
	msg.ReplyMarkup = tg.NewReplyKeyboard(
		tg.NewKeyboardButtonRow(
			tg.NewKeyboardButton(string(enums.ActionTake)),
			tg.NewKeyboardButton(string(enums.ActionEdit)),
		),
	)

	if _, err := b.api.Send(msg); err != nil {
		slog.Error("Send message", "err", err)
	}
}

func (b *BotService) handleMessage(message *tg.Message) {
	chatId := message.Chat.ID

	if message.Text == string(enums.ActionCreate) {
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

		return
	}

	user, err := b.userService.GetByChatId(chatId)

	if err != nil {
		slog.Error(err.Error())
	}

	if user.Status == string(enums.UserStatusEditing) || user.Status == string(enums.UserStatusInactive) {
		timeRegex := regexp.MustCompile(`^(?:[01]\d|2[0-3]):[0-5]\d$`)

		idleStatus := string(enums.UserStatusIdle)

		isTime := timeRegex.MatchString(message.Text)
		isTimezone := utils.IsValidTimezone(message.Text)

		if isTime || isTimezone {
			var toUpdate = model.UserUpdate{Status: &idleStatus}

			parsedTime := utils.GetTimeFromStringWithServerTimezone(message.Text, &user.Timezone)
			timeToNotify := parsedTime.Format("15:04")

			if isTime {
				toUpdate.TimeToNotify = &timeToNotify
			} else if isTimezone {
				toUpdate.Timezone = &message.Text
			}

			b.userService.Update(chatId, toUpdate)
			b.cronNotifier.AddOrUpdateCron(chatId, parsedTime)

			b.SendMessage(chatId, "Time was updated!")
		} else {
			b.SendMessage(chatId, "Not valid time or timezone")
		}

		return
	}

	if user.Status == string(enums.UserStatusIdle) {
		if message.Text == string(enums.ActionTake) {
			err := b.pillDayService.MarkAsTakenNow(message.Chat.ID)

			if err != nil {
				slog.Error(err.Error())
				b.SendMessage(chatId, "Try again")
				return
			}

			b.SendMessage(chatId, "Checked!")
		}

		if message.Text == string(enums.ActionEdit) {
			status := string(enums.UserStatusEditing)
			b.userService.Update(chatId, model.UserUpdate{Status: &status})

			b.SendMessage(chatId, "Enter new time to get notified - 15:04 format")
		}

		return
	}
}
