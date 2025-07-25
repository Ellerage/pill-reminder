package tests

import (
	"log/slog"
	"pill-reminder/internal/i18n"
	utilscommon "pill-reminder/internal/utils"
	"pill-reminder/internal/utils/enums"
	"pill-reminder/tests/seeds"
	"pill-reminder/tests/utils"
	"testing"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStartAppAndNotifications(t *testing.T) {
	modules, teardown := utils.Setup(t)

	nowTime := time.Now()
	timeToNotify := nowTime.Add(1 * time.Minute).Format("15:04")
	loc := time.Now().Location().String()

	timeToNotifyUTC, err := utilscommon.GetUTCFromUserTime(timeToNotify, &loc)
	if err != nil {
		slog.Error(err.Error())
	}

	remindIntervalMinutes := uint8(1)
	userSeed := seeds.UserSeed(modules.DB, seeds.UserParams{TimeToNotify: &timeToNotifyUTC, RemindIntervalMinutes: &remindIntervalMinutes})
	userChatId := userSeed.ChatId

	users, err := modules.UserService.GetAll()
	if err != nil {
		slog.Error(err.Error())
	}

	err = modules.ReminderQueue.Start(users)
	if err != nil {
		panic(err)
	}

	utils.StartScheduleHandlers(utils.StartScheduleHandlersParams{
		ReminderQueue:        modules.ReminderQueue,
		ReminderQueueService: modules.ReminderQueueService,
		Bot:                  modules.Bot,
		PillDayService:       modules.PillDayService,
		UserService:          modules.UserService,
	})

	expected := []string{
		i18n.GetText("firstNotification"),
		i18n.GetText("reminderNotification"),
	}

	for i := range 2 {
		select {
		case ch := <-modules.BotAPI.SendCalls:
			msg, ok := ch.(tgbotapi.MessageConfig)
			require.True(t, ok, "expected MessageConfig, got %T", ch)
			assert.Equal(t, userChatId, msg.ChatID)
			assert.Contains(t, expected[i], msg.Text)
		case <-time.After(time.Minute * 1):
			t.Fatal("Timeout")
		}
	}

	message := utils.GenerateMessage(userChatId, string(enums.ActionTake))
	err = modules.Bot.HandleMessage(message)
	if err != nil {
		t.Fatal(err.Error())
	}

	pillDay, err := utils.GetPillDayByChatId(modules.DB, userChatId)
	if err != nil {
		t.Fatal(err.Error())
	}

	assert.True(t, pillDay.HasTimeOfTaking())

	select {
	case ch := <-modules.BotAPI.SendCalls:
		msg, ok := ch.(tgbotapi.MessageConfig)
		require.True(t, ok, "expected MessageConfig, got %T", ch)
		assert.Equal(t, userChatId, msg.ChatID)
		assert.Contains(t, i18n.GetText("checked"), msg.Text)
	case <-time.After(time.Minute * 1):
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
