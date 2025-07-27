package utils

import (
	"testing"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ValidateReplyMessage(t *testing.T, SendCalls chan tgbotapi.Chattable, userChatId int64, timeout time.Duration, text string) {
	select {
	case ch := <-SendCalls:
		msg, ok := ch.(tgbotapi.MessageConfig)
		require.True(t, ok, "expected MessageConfig, got %T", ch)
		assert.Equal(t, userChatId, msg.ChatID)
		assert.Contains(t, text, msg.Text)
	case <-time.After(timeout):
		t.Fatal("Timeout")
	}
}
