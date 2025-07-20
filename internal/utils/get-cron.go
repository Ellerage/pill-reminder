package utils

import (
	"fmt"
	"time"
)

func GetCronFromMinutes(minutes uint8) string {
	cronStr := fmt.Sprintf("*/%d * * * *", minutes)

	return cronStr
}

func GetDailyCronFromStringTime(timeStr string) (string, error) {
	validTime, err := time.Parse("15:04", timeStr)

	if err != nil {
		return "", err
	}

	cron := fmt.Sprintf("%d %d * * *", validTime.Minute(), validTime.Hour())

	return cron, nil
}
