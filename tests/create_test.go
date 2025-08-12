package tests

import (
	"pill-reminder/tests/seeds"
	"pill-reminder/tests/utils"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

func TestCreateFlow(t *testing.T) {
	modules, teardown := utils.Setup(t)

	fakeChatId := gofakeit.Int64()
	message := utils.GenerateMessage(fakeChatId, "/start")

	modules.Bot.HandleMessage(message)

	user, err := seeds.GetUserByChatId(t, modules.DB, fakeChatId)
	if err != nil {
		t.Fatal(err.Error())
	}

	assert.Equal(t, fakeChatId, user.ChatId)

	t.Cleanup(teardown)
}
