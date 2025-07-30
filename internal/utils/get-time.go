package utils

import (
	"log/slog"
	"time"
	_ "time/tzdata"
)

const TimezoneDefault = "UTC"

func GetNowDateTime() time.Time {
	loc, err := time.LoadLocation(TimezoneDefault)
	if err != nil {
		panic(err)
	}

	now := time.Now().In(loc)

	return now
}

func GetFormattedNowDate() string {
	return GetNowDateTime().Format("2006-01-02")
}

func GetUTCFromUserTime(timeStr string, userTimeZone string) (string, error) {
	loc, err := time.LoadLocation(userTimeZone)
	if err != nil {
		return "", err
	}

	parsed, err := time.ParseInLocation("15:04", timeStr, loc)
	if err != nil {
		return "", err
	}

	now := time.Now().In(loc)

	return time.Date(
		now.Year(), now.Month(), now.Day(),
		parsed.Hour(), parsed.Minute(),
		0, 0,
		loc,
	).UTC().Format("15:04"), nil
}

func GetTimeFromString(str string) time.Time {
	parsed, err := time.ParseInLocation("15:04", str, time.UTC)
	if err != nil {
		slog.Error(err.Error())
	}

	return parsed
}

func IsValidTimezone(name string) bool {
	_, err := time.LoadLocation(name)
	return err == nil
}
