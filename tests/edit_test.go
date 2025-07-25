package tests

import (
	"fmt"
	"pill-reminder/internal/utils/enums"
	"pill-reminder/tests/seeds"
	"pill-reminder/tests/utils"
	"strconv"
	"testing"

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
}
