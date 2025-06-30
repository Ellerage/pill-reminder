package tgbotapi

import (
	"log"
	"log/slog"
	"pill-reminder/configs"
	cronnotifier "pill-reminder/internal/cron-notifier"
	"pill-reminder/internal/model"
	"pill-reminder/internal/service"
	"pill-reminder/internal/utils"
	"pill-reminder/internal/utils/enums"
	"regexp"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TgNotifier struct{}

func (TgNotifier) SendMessage(chatId int64, message string) {
	SendMessage(chatId, message)
}

type BotAPIDeps struct {
	PillDayService *service.PillDayService
	UserService    *service.UserService
	CronNotifier   *cronnotifier.CronNotifier
	Config         *configs.Config
}

var (
	bot  *tgbotapi.BotAPI
	once sync.Once
	err  error
)

func Init(token string) {
	once.Do(func() {
		bot, err = tgbotapi.NewBotAPI(token)
		if err != nil {
			log.Fatalf("failed to initialize Telegram bot: %v", err)
		}

		log.Println("TG bot init")
	})
}

func RegisterMessageListener(deps BotAPIDeps) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	log.Println("Listening for new messages...")

	for update := range updates {
		if update.Message != nil {
			handleMessage(deps, update.Message)
		}
	}
}

func SendMessage(chatId int64, message string) {
	msg := tgbotapi.NewMessage(chatId, message)

	replyKeyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(string(enums.ActionTake)),
			tgbotapi.NewKeyboardButton(string(enums.ActionEdit)),
		),
	)
	msg.ReplyMarkup = replyKeyboard

	_, err = bot.Send(msg)

	if err != nil {
		slog.Error(err.Error())
	}
}

func handleMessage(deps BotAPIDeps, message *tgbotapi.Message) {
	var chatId = message.Chat.ID

	if message.Text == string(enums.ActionCreate) {
		err := deps.UserService.Create(model.User{
			ChatId:       chatId,
			Timezone:     deps.Config.TIMEZONE,
			TimeToNotify: "00:00",
			Status:       string(enums.UserStatusInactive),
		})

		if err != nil {
			slog.Error(err.Error())
		} else {
			SendMessage(chatId, "What's time you want to get reminders? Type it in 15:04 format")
		}

		return
	}

	user, err := deps.UserService.GetByChatId(chatId)

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

			deps.UserService.Update(chatId, toUpdate)
			deps.CronNotifier.AddOrUpdateCron(chatId, parsedTime)

			SendMessage(chatId, "Time was updated!")
		} else {
			SendMessage(chatId, "Not valid time or timezone")
		}

		return
	}

	if user.Status == string(enums.UserStatusIdle) {
		if message.Text == string(enums.ActionTake) {
			err := deps.PillDayService.MarkAsTakenNow(message.Chat.ID)

			if err != nil {
				slog.Error(err.Error())
				SendMessage(chatId, "Try again")
				return
			}

			SendMessage(chatId, "Checked!")
		}

		if message.Text == string(enums.ActionEdit) {
			status := string(enums.UserStatusEditing)
			deps.UserService.Update(chatId, model.UserUpdate{Status: &status})

			SendMessage(chatId, "Enter new time to get notified - 15:04 format")
		}

		return
	}
}
