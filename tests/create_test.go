package tests

import (
	"pill-reminder/tests/utils"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

func TestCreateFlow(t *testing.T) {
	modules, teardown := utils.Setup(t)
	defer teardown()

	fakeChatId := gofakeit.Int64()
	message := utils.GenerateMessage(fakeChatId, "/start")

	modules.Bot.HandleMessage(message)

	user, err := utils.GetUserByChatId(modules.DB, fakeChatId)
	if err != nil {
		t.Fatal(err.Error())
	}

	assert.Equal(t, fakeChatId, user.ChatId)
}
