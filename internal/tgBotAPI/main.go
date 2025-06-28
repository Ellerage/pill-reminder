package tgbotapi

import (
	"fmt"
	"log"
	"pill-reminder/configs"
	"pill-reminder/internal/service"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotAPIDeps struct {
	PillDayService *service.PillDayService
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

// func GetBot() *tgbotapi.BotAPI {
// 	if bot == nil {
// 		Init()
// 	}

// 	return bot
// }

func SendMessage(chatId int64, message string) {
	msg := tgbotapi.NewMessage(chatId, message)

	replyKeyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Take"),
			// tgbotapi.NewKeyboardButton("Delay 30 min"),
		),
	)
	msg.ReplyMarkup = replyKeyboard

	_, err = bot.Send(msg)

	if err != nil {
		fmt.Println("Ошибка отправки:", err)
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

// TODO: Make it more common and save userId/chatId in db
func handleMessage(deps BotAPIDeps, message *tgbotapi.Message) {
	if message.Chat.ID == deps.Config.MY_CHAT_ID && message.Text == "Take" {
		deps.PillDayService.MarkAsTakenNow()
	}
}
