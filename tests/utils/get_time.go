package utils

import (
	utilscommon "pill-reminder/internal/utils"
	"time"
)

func GetMinuteAheadNowUTC() (string, error) {
	nowTime := time.Now()
	timeToNotify := nowTime.Add(1 * time.Minute).Format("15:04")
	loc := time.Now().Location().String()

	timeToNotifyUTC, err := utilscommon.GetUTCFromUserTime(timeToNotify, &loc)
	if err != nil {
		return "", err
	}

	return timeToNotifyUTC, nil
}
