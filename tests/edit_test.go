package tests

import (
	"fmt"
	"pill-reminder/internal/i18n"
	"pill-reminder/internal/utils/enums"
	"pill-reminder/tests/seeds"
	"pill-reminder/tests/utils"
	"strconv"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

func TestEnterEditingState(t *testing.T) {
	modules, teardown := utils.Setup(t)
	defer teardown()

	userSeed := seeds.UserSeed(modules.DB, seeds.UserParams{})
	userChatId := userSeed.ChatId

	messageInit := utils.GenerateMessage(userChatId, string(enums.ActionEdit))
	modules.Bot.HandleMessage(messageInit)

	user2, err := utils.GetUserByChatId(modules.DB, userChatId)
	if err != nil {
		t.Fatal(err.Error())
	}
	assert.Equal(t, string(enums.UserStatusEditing), user2.Status)

	utils.ValidateReplyMessage(t, modules.BotAPI.SendCalls, userChatId, time.Second*5, i18n.GetText("enterNewTime"))
}

func TestEditTimeToNotify(t *testing.T) {
	modules, teardown := utils.Setup(t)
	defer teardown()

	status := string(enums.UserStatusEditing)
	userSeed := seeds.UserSeed(modules.DB, seeds.UserParams{Status: &status})
	userChatId := userSeed.ChatId

	newTimeToNotify := gofakeit.Date().Format("15:04")
	messageNewTime := utils.GenerateMessage(userChatId, newTimeToNotify)

	modules.Bot.HandleMessage(messageNewTime)

	updatedUser, err := utils.GetUserByChatId(modules.DB, userChatId)
	if err != nil {
		t.Fatal(err.Error())
	}

	assert.Equal(t, newTimeToNotify, updatedUser.TimeToNotify)
	utils.ValidateReplyMessage(t, modules.BotAPI.SendCalls, userChatId, time.Second*5, i18n.GetText("firstAtDayNotificationTimeUpdated"))
}

func TestEditTimeAndReminderNotifications(t *testing.T) {
	modules, teardown := utils.Setup(t)
	defer teardown()

	utils.InitScheduleForAllUsers(utils.StartScheduleHandlersParams{
		ReminderQueue:        modules.ReminderQueue,
		ReminderQueueService: modules.ReminderQueueService,
		Bot:                  modules.Bot,
		PillDayService:       modules.PillDayService,
		UserService:          modules.UserService,
	})

	status := string(enums.UserStatusEditing)
	userSeed := seeds.UserSeed(modules.DB, seeds.UserParams{Status: &status})
	userChatId := userSeed.ChatId

	newTimeToNotify := utils.GetNowUTCTime(1 * time.Minute)
	messageNewTime := utils.GenerateMessage(userChatId, newTimeToNotify)

	modules.Bot.HandleMessage(messageNewTime)

	utils.ValidateReplyMessage(t, modules.BotAPI.SendCalls, userChatId, time.Second*5, i18n.GetText("firstAtDayNotificationTimeUpdated"))

	utils.ValidateReplyMessage(t, modules.BotAPI.SendCalls, userChatId, time.Minute*2, i18n.GetText("firstNotification"))

	utils.ValidateReplyMessage(t, modules.BotAPI.SendCalls, userChatId, time.Minute*2, i18n.GetText("reminderNotification"))
}

func TestEditRemindInterval(t *testing.T) {
	modules, teardown := utils.Setup(t)
	defer teardown()

	status := string(enums.UserStatusEditing)
	userSeed := seeds.UserSeed(modules.DB, seeds.UserParams{Status: &status})
	userChatId := userSeed.ChatId

	expectedRemindIntervalMinutes := strconv.Itoa(gofakeit.Minute())
	messageNewRemindInterval := utils.GenerateMessage(userChatId, expectedRemindIntervalMinutes)
	modules.Bot.HandleMessage(messageNewRemindInterval)

	updatedRemindIntervalUser, err := utils.GetUserByChatId(modules.DB, userChatId)
	if err != nil {
		t.Fatal(err.Error())
	}

	expected := fmt.Sprintf("*/%s * * * *", expectedRemindIntervalMinutes)

	assert.Equal(t, expected, updatedRemindIntervalUser.RemindInterval)
	utils.ValidateReplyMessage(t, modules.BotAPI.SendCalls, userChatId, time.Second*5, i18n.GetText("repeatIntervalTimeUpdated"))
}
