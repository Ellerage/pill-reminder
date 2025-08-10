package tests

import (
	"log/slog"
	"pill-reminder/internal/utils/enums"
	"pill-reminder/tests/seeds"
	"pill-reminder/tests/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarkAsTakenFlow(t *testing.T) {
	modules, teardown := utils.Setup(t)
	user := seeds.UserSeed(modules.DB, nil)
	chatId := user.ChatId

	takeMessage := utils.GenerateMessage(chatId, string(enums.ActionTake))

	err := modules.Bot.HandleMessage(takeMessage)
	if err != nil {
		t.Fatal(err)
	}

	actualPillDay, err := seeds.FindPillDayByChatId(t, modules.DB, chatId)
	if err != nil {
		slog.Error(err.Error())
	}

	assert.NotNil(t, actualPillDay.Date)
	assert.NotEmpty(t, actualPillDay.Date)
	t.Cleanup(teardown)
}
