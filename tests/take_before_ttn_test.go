package tests

import (
	"pill-reminder/internal/i18n"
	"pill-reminder/internal/utils/enums"
	"pill-reminder/tests/seeds"
	"pill-reminder/tests/utils"
	"testing"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTakeBeforeTimeToNotify(t *testing.T) {
	modules, teardown := utils.Setup(t)

	seedTimeToNotify, err := utils.GetMinuteAheadNowUTC()
	if err != nil {
		t.Fatal(err.Error())
	}

	user := seeds.UserSeed(modules.DB, seeds.UserParams{
		TimeToNotify: &seedTimeToNotify,
	})

	utils.InitScheduleForAllUsers(utils.StartScheduleHandlersParams{
		ReminderQueue:        modules.ReminderQueue,
		ReminderQueueService: modules.ReminderQueueService,
		Bot:                  modules.Bot,
		PillDayService:       modules.PillDayService,
		UserService:          modules.UserService,
	})

	message := utils.GenerateMessage(user.ChatId, string(enums.ActionTake))
	err = modules.Bot.HandleMessage(message)
	if err != nil {
		t.Fatal(err.Error())
	}

	pillDayActual, err := utils.GetPillDayByChatId(modules.DB, user.ChatId)
	if err != nil {
		t.Fatal(err.Error())
	}

	assert.True(t, pillDayActual.HasTimeOfTaking())

	select {
	case ch := <-modules.BotAPI.SendCalls:
		msg, ok := ch.(tgbotapi.MessageConfig)
		require.True(t, ok, "expected MessageConfig, got %T", ch)
		assert.Equal(t, user.ChatId, msg.ChatID)
		assert.Contains(t, i18n.GetText("checked"), msg.Text)
	case <-time.After(time.Second * 10):
		t.Fatal("Timeout")
	}

	select {
	case <-modules.BotAPI.SendCalls:
		t.Fatal("Notifications didn't stop")
	case <-time.After(2 * time.Minute):
	}

	t.Cleanup(func() {
		teardown()
	})
}
