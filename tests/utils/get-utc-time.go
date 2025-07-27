package utils

import (
	"log/slog"
	utilscommon "pill-reminder/internal/utils"
	"time"
)

func GetNowUTCTime(add time.Duration) string {
	nowTime := time.Now()
	timeToNotify := nowTime.Add(add).Format("15:04")
	loc := time.Now().Location().String()

	timeToNotifyUTC, err := utilscommon.GetUTCFromUserTime(timeToNotify, &loc)
	if err != nil {
		slog.Error(err.Error())
	}

	return timeToNotifyUTC
}
