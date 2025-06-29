package tgbotapi

import (
	"log"
	"log/slog"
	"pill-reminder/configs"
	"pill-reminder/internal/model"
	"pill-reminder/internal/service"
	"pill-reminder/internal/utils"
	"pill-reminder/internal/utils/enums"
	"regexp"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotAPIDeps struct {
	PillDayService *service.PillDayService
	UserService    *service.UserService
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

func SendMessage(chatId int64, message string) {
	msg := tgbotapi.NewMessage(chatId, message)

	replyKeyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(string(enums.Take)),
			tgbotapi.NewKeyboardButton(string(enums.Edit)),
		),
	)
	msg.ReplyMarkup = replyKeyboard

	_, err = bot.Send(msg)

	if err != nil {
		slog.Error(err.Error())
	}
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

func handleMessage(deps BotAPIDeps, message *tgbotapi.Message) {
	var chatId = message.Chat.ID

	user, err := deps.UserService.GetByChatId(chatId)

	if err != nil {
		slog.Error(err.Error())
	}

	if user.Status == string(enums.UserStatusEditing) {
		timeRegex := regexp.MustCompile(`^(?:[01]\d|2[0-3]):[0-5]\d$`)

		idleStatus := string(enums.UserStatusIdle)

		if timeRegex.MatchString(message.Text) {
			deps.UserService.Update(chatId, model.UserUpdate{TimeToNotify: &message.Text, Status: &idleStatus})
			SendMessage(chatId, "Time was updated!")
		} else if utils.IsValidTimezone(message.Text) {
			deps.UserService.Update(chatId, model.UserUpdate{Timezone: &message.Text, Status: &idleStatus})
			SendMessage(chatId, "Timezone was updated!")
		} else {
			SendMessage(chatId, "Not valid time or timezone")
		}

		return
	}

	if user.Status == string(enums.UserStatusIdle) {
		if message.Text == string(enums.Take) {
			err := deps.PillDayService.MarkAsTakenNow(message.Chat.ID)

			if err != nil {
				slog.Error(err.Error())
				SendMessage(chatId, "Try again")
				return
			}

			SendMessage(chatId, "Checked!")
		}

		if message.Text == string(enums.Edit) {
			status := string(enums.UserStatusEditing)
			deps.UserService.Update(chatId, model.UserUpdate{Status: &status})

			SendMessage(chatId, "Enter new time to get notified - 15:04 format")
		}

		return
	}

}
