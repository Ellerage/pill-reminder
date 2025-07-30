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
	user := seeds.UserSeed(modules.DB, seeds.UserParams{})
	chatId := user.ChatId

	takeMessage := utils.GenerateMessage(chatId, string(enums.ActionTake))

	modules.Bot.HandleMessage(takeMessage)

	actualPillDay, err := utils.GetPillDayByChatId(modules.DB, chatId)
	if err != nil {
		slog.Error(err.Error())
	}

	assert.NotNil(t, actualPillDay.Date)
	assert.NotEmpty(t, actualPillDay.Date)
	t.Cleanup(teardown)
}
