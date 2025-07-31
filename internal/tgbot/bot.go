package tgbot

import (
	"context"
	"errors"
	"log/slog"
	"pill-reminder/internal/i18n"
	"pill-reminder/internal/utils"
	"pill-reminder/internal/utils/enums"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/v2/mongo"
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
				chatId := update.Message.Chat.ID
				err := b.HandleMessage(update.Message)

				switch {
				case errors.Is(err, mongo.ErrNoDocuments):
					b.SendMessage(chatId, i18n.GetText("noAccount"), &enums.SendMessageButtons{Create: true}, nil)
				case errors.Is(err, utils.ErrInvalidCommand):
					b.SendInfoMessage(chatId)
				case errors.Is(err, utils.ErrUserAlreadyExist):
					b.SendMessage(chatId, i18n.GetText("accountAlreadyExist"), nil, nil)
					b.SendInfoMessage(chatId)
				case errors.Is(err, utils.ErrAlreadyTakenToday):
					b.SendMessage(chatId, i18n.GetText("pillAlreadyTaken"), nil, nil)
				case errors.Is(err, utils.ErrInvalidTimeEditInput):
					b.SendMessage(chatId, i18n.GetText("ErrInvalidTimeEditInput"), nil, nil)
				case err != nil:
					b.SendMessage(update.Message.Chat.ID, i18n.GetText("tryAgain"), nil, nil)
				default:
				}
			}
		case <-ctx.Done():
			slog.Info("Bot listener stopped")
			return
		}
	}
}

func (b *BotService) SendMessage(chatID int64, message string, buttons *enums.SendMessageButtons, options *MessageOptions) error {
	msg := tg.NewMessage(chatID, message)
	replyButtons := make([]tg.KeyboardButton, 0, 10)

	if options != nil {
		if options.ParseMode != nil {
			msg.ParseMode = *options.ParseMode
		}
	}

	if buttons != nil {
		if buttons.Create {
			replyButtons = append(replyButtons, tg.NewKeyboardButton(string(enums.ActionCreate)))
		}
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
		slog.Error(err.Error())
		return err
	}

	return nil
}
