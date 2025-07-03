package tgbotapi

import (
	"log"
	"log/slog"
	cronnotifier "pill-reminder/internal/cron-notifier"
	"pill-reminder/internal/service"
	"pill-reminder/internal/utils/enums"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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
