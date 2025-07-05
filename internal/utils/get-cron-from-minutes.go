package utils

import "fmt"

// Return cron format - every ? minutes
func GetCronFromMinutes(minutes uint8) string {
	cronStr := fmt.Sprintf("*/%d * * * *", minutes)

	return cronStr
}
