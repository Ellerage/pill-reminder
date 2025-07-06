package tgbotapi

import (
	"log"
	"log/slog"
	"pill-reminder/internal/utils/enums"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotService struct {
	timezone       string
	api            BotAPI
	userService    UserService
	pillDayService PillDayService
	cronNotifier   CronNotifier
}

type BotServiceParams struct {
	Timezone       string
	API            BotAPI
	UserService    UserService
	PillDayService PillDayService
	CronNotifier   CronNotifier
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

func (b *BotService) SendMessage(chatID int64, message string, buttons *enums.SendMessageButtons) {
	msg := tg.NewMessage(chatID, message)
	replyButtons := make([]tg.KeyboardButton, 0, 2)

	if buttons != nil {
		if buttons.Take {
			replyButtons = append(replyButtons, tg.NewKeyboardButton(string(enums.ActionTake)))

		}
		if buttons.Edit {
			replyButtons = append(replyButtons, tg.NewKeyboardButton(string(enums.ActionEdit)))
		}
	}

	if len(replyButtons) > 0 {
		msg.ReplyMarkup = tg.NewReplyKeyboard(
			tg.NewKeyboardButtonRow(
				replyButtons...,
			),
		)
	}

	if _, err := b.api.Send(msg); err != nil {
		slog.Error("Send message", "err", err)
	}
}
