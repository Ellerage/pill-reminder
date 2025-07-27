package tests

import (
	"pill-reminder/internal/i18n"
	"pill-reminder/internal/utils/enums"
	"pill-reminder/tests/seeds"
	"pill-reminder/tests/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

	// Mark as taken
	message := utils.GenerateMessage(user.ChatId, string(enums.ActionTake))
	err = modules.Bot.HandleMessage(message)
	if err != nil {
		t.Fatal(err.Error())
	}

	// Check updated
	pillDayActual, err := utils.GetPillDayByChatId(modules.DB, user.ChatId)
	if err != nil {
		t.Fatal(err.Error())
	}
	assert.True(t, pillDayActual.HasTimeOfTaking())

	// check reply for Checked message
	utils.ValidateReplyMessage(t, modules.BotAPI.SendCalls, user.ChatId, time.Second*10, i18n.GetText("checked"))

	// Check that notifications was stopped
	select {
	case <-modules.BotAPI.SendCalls:
		t.Fatal("Notifications didn't stop")
	case <-time.After(2 * time.Minute):
	}

	t.Cleanup(func() {
		teardown()
	})
}
