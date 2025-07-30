package tests

import (
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
	t.Cleanup(teardown)
}

func TestEditTimeToNotify(t *testing.T) {
	modules, teardown := utils.Setup(t)

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
	t.Cleanup(teardown)
}

func TestEditTimeAndReminderNotifications(t *testing.T) {
	modules, teardown := utils.Setup(t)

	utils.InitScheduleForAllUsers(utils.StartScheduleHandlersParams{
		ReminderQueue: modules.ReminderQueue,
		UserService:   modules.UserService,
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
	t.Cleanup(teardown)
}

func TestEditRemindInterval(t *testing.T) {
	modules, teardown := utils.Setup(t)

	status := string(enums.UserStatusEditing)
	userSeed := seeds.UserSeed(modules.DB, seeds.UserParams{Status: &status})
	userChatId := userSeed.ChatId

	expectedRemindIntervalMinutes := gofakeit.Minute()
	messageNewRemindInterval := utils.GenerateMessage(userChatId, strconv.Itoa(expectedRemindIntervalMinutes))
	modules.Bot.HandleMessage(messageNewRemindInterval)

	updatedRemindIntervalUser, err := utils.GetUserByChatId(modules.DB, userChatId)
	if err != nil {
		t.Fatal(err.Error())
	}

	assert.Equal(t, int64(expectedRemindIntervalMinutes), updatedRemindIntervalUser.RemindInterval)
	utils.ValidateReplyMessage(t, modules.BotAPI.SendCalls, userChatId, time.Second*5, i18n.GetText("repeatIntervalTimeUpdated"))
	t.Cleanup(teardown)
}
