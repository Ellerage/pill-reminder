package utils

import (
	"fmt"
	"pill-reminder/internal/i18n"
	"pill-reminder/internal/model"
)

func GetSettingsReplyText(user model.UserNotificationSettings) string {
	text := fmt.Sprintf(
		"<b>Time to notify:</b> %s\n<b>Remind interval:</b> %s\n<b>Timezone:</b> %s",
		user.TimeToNotify,
		i18n.Plural(int(user.RemindInterval), "Minute"),
		user.Timezone,
	)

	return text
}
