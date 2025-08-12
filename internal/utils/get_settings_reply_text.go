package utils

import (
	"fmt"
	"log/slog"
	"pill-reminder/internal/i18n"
	"pill-reminder/internal/model"
)

func GetSettingsReplyText(user model.UserNotificationSettings) string {
	userTimeToNotify, err := GetUserTimeFromUTC(user.TimeToNotify, user.Timezone)
	if err != nil {
		userTimeToNotify = ""
		slog.Error(err.Error())
	}

	text := fmt.Sprintf(
		"<b>Time to notify:</b> %s\n<b>Remind interval:</b> %s\n<b>Timezone:</b> %s",
		userTimeToNotify,
		i18n.Plural(int(user.RemindInterval), "Minute"),
		user.Timezone,
	)

	return text
}
