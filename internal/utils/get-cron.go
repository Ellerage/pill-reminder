package utils

import (
	"fmt"
	"time"
)

func GetCronFromMinutes(minutes uint8) string {
	cronStr := fmt.Sprintf("*/%d * * * *", minutes)

	return cronStr
}

func GetCronFromStringUTCTime(timeStr string) (string, error) {
	parsed, err := time.ParseInLocation("15:04", timeStr, time.UTC)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d %d * * *", parsed.Minute(), parsed.Hour()), nil
}
