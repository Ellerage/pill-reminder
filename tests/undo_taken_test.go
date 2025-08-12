package tests

import (
	"log/slog"
	"pill-reminder/internal/i18n"
	"pill-reminder/internal/utils/enums"
	"pill-reminder/tests/seeds"
	"pill-reminder/tests/utils"
	"testing"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
)

func TestUndoTaken(t *testing.T) {
	t.Parallel()
	modules, teardown := utils.Setup(t)

	seedTimeToNotify, err := utils.GetMinuteAheadNowUTC()
	if err != nil {
		t.Fatal(err.Error())
	}

	user := seeds.UserSeed(modules.DB, &seeds.UserParams{
		TimeToNotify: &seedTimeToNotify,
	})

	utils.InitScheduleForAllUsers(utils.StartScheduleHandlersParams{
		ReminderQueue: modules.ReminderQueue,
		UserService:   modules.UserService,
	})

	// Mark as taken
	message := utils.GenerateMessage(user.ChatId, string(enums.ActionTake))
	err = modules.Bot.HandleMessage(message)
	if err != nil {
		t.Fatal(err.Error())
	}

	// Check updated
	pillDayActual, err := seeds.FindPillDayByChatId(t, modules.DB, user.ChatId)
	if err != nil {
		t.Fatal(err.Error())
	}
	assert.True(t, pillDayActual.HasTimeOfTaking())

	// check reply for Checked message
	utils.ValidateReplyMessage(t, modules.BotAPI.SendCalls, user.ChatId, time.Second*10, i18n.GetText("checked"))

	// Check that notifications was stopped
	select {
	case v := <-modules.BotAPI.SendCalls:
		msg, _ := v.(tgbotapi.MessageConfig)
		slog.Info(msg.Text)
		t.Fatal("Notifications didn't stop")
	case <-time.After(65 * time.Second):
	}

	// Undo action
	message = utils.GenerateMessage(user.ChatId, string(enums.ActionUndo))
	err = modules.Bot.HandleMessage(message)
	if err != nil {
		t.Fatal(err.Error())
	}
	utils.ValidateReplyMessage(t, modules.BotAPI.SendCalls, user.ChatId, 65*time.Second, i18n.GetText("undoneTaken"))

	// Check notification resumed
	utils.ValidateReplyMessage(t, modules.BotAPI.SendCalls, user.ChatId, 65*time.Second, i18n.GetText("reminderNotification"))

	t.Cleanup(teardown)
}
