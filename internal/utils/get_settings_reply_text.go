package utils

import (
	"fmt"
	"pill-reminder/internal/model"
)

func GetSettingsReplyText(user model.UserNotificationSettings) string {
	text := fmt.Sprintf(
		"<b>Time to notify:</b> %s\n<b>Remind interval:</b> %s\n<b>Timezone:</b> %s",
		user.TimeToNotify,
		fmt.Sprintf("%d Minutes", user.RemindInterval),
		user.Timezone,
	)

	return text
}
