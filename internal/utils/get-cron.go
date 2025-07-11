package utils

import (
	"fmt"
	"log/slog"
	"time"
)

// Return cron format - every ? minutes
func GetCronFromMinutes(minutes uint8) string {
	cronStr := fmt.Sprintf("*/%d * * * *", minutes)

	return cronStr
}

func GetDailyCronFromStringTime(timeStr string) string {
	validTime, err := time.Parse("15:04", timeStr)

	if err != nil {
		slog.Error(err.Error())
	}

	cron := fmt.Sprintf("%d %d * * *", validTime.Minute(), validTime.Hour())

	return cron
}
