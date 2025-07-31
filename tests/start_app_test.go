package tests

import (
	"fmt"
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
)

func TestStartAppAndNotifications(t *testing.T) {
	modules, teardown := utils.Setup(t)

	nowTime := time.Now()
	timeToNotify := nowTime.Add(1 * time.Minute).Format("15:04")
	loc := time.Now().Location().String()

	timeToNotifyUTC, err := utilscommon.GetUTCFromUserTime(timeToNotify, loc)

	fmt.Println("timeToNotifyUTC", timeToNotifyUTC)
	if err != nil {
		slog.Error(err.Error())
	}

	userSeed := seeds.UserSeed(modules.DB, seeds.UserParams{TimeToNotify: &timeToNotifyUTC})
	userChatId := userSeed.ChatId

	utils.InitScheduleForAllUsers(utils.StartScheduleHandlersParams{ReminderQueue: modules.ReminderQueue, UserService: modules.UserService})

	expected := []string{
		i18n.GetText("firstNotification"),
		i18n.GetText("reminderNotification"),
	}

	for i := range 2 {
		utils.ValidateReplyMessage(t, modules.BotAPI.SendCalls, userChatId, time.Second*90, expected[i])
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

	utils.ValidateReplyMessage(t, modules.BotAPI.SendCalls, userChatId, time.Second*90, i18n.GetText("checked"))

	select {
	case v := <-modules.BotAPI.SendCalls:
		msg, _ := v.(tgbotapi.MessageConfig)
		fmt.Println(msg.Text)
		t.Fatal("Notifications didn't stop")
	case <-time.After(2 * time.Minute):
	}

	t.Cleanup(teardown)
}
