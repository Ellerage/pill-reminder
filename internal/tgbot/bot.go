package tgbot

import (
	"context"
	"log/slog"
	"pill-reminder/internal/utils/enums"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotService struct {
	timezone        string
	api             BotAPI
	userService     UserService
	pillDayService  PillDayService
	reminderService ReminderService
	reminderQueue   ReminderQueue
}

type BotServiceParams struct {
	Timezone        string
	API             BotAPI
	UserService     UserService
	PillDayService  PillDayService
	ReminderService ReminderService
	ReminderQueue   ReminderQueue
}

type MessageOptions struct {
	ParseMode *string
}

func NewBotService(params BotServiceParams) *BotService {
	return &BotService{
		timezone:        params.Timezone,
		api:             params.API,
		userService:     params.UserService,
		pillDayService:  params.PillDayService,
		reminderQueue:   params.ReminderQueue,
		reminderService: params.ReminderService,
	}
}

func (b *BotService) RegisterMessageListener(ctx context.Context) {
	u := tg.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	slog.Info("Listening for new messages...")

	for {
		select {
		case update := <-updates:
			if update.Message != nil {
				b.HandleMessage(update.Message)
			}
		case <-ctx.Done():
			slog.Info("Bot listener stopped")
			return
		}
	}
}

func (b *BotService) SendMessage(chatID int64, message string, buttons *enums.SendMessageButtons, options *MessageOptions) {
	msg := tg.NewMessage(chatID, message)
	replyButtons := make([]tg.KeyboardButton, 0, 2)

	if options != nil {
		if options.ParseMode != nil {
			msg.ParseMode = *options.ParseMode
		}
	}

	if buttons != nil {
		if buttons.Take {
			replyButtons = append(replyButtons, tg.NewKeyboardButton(string(enums.ActionTake)))

		}
		if buttons.Edit {
			replyButtons = append(replyButtons, tg.NewKeyboardButton(string(enums.ActionEdit)))
		}
		if buttons.Delay {
			replyButtons = append(replyButtons, tg.NewKeyboardButton(string(enums.ActionDelay)))
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
