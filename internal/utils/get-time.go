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

func GetUTCFromUserTime(timeStr string, userTimeZone *string) string {
	var timezone string

	if userTimeZone != nil {
		timezone = *userTimeZone
	} else {
		timezone = TimezoneDefault
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		slog.Error(err.Error())
	}

	now := time.Now().In(loc)
	parsed, err := time.Parse("15:04", timeStr)

	if err != nil {
		slog.Error(err.Error())
	}

	return time.Date(
		now.Year(), now.Month(), now.Day(),
		parsed.Hour(), parsed.Minute(), 0, 0,
		loc,
	).UTC().Format("15:04")
}

func GetTimeFromString(str string) time.Time {
	parsed, err := time.Parse("15:04", str)

	if err != nil {
		slog.Error(err.Error())
	}

	return parsed
}

func IsValidTimezone(name string) bool {
	_, err := time.LoadLocation(name)
	return err == nil
}
