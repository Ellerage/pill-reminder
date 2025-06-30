package utils

import (
	"time"
)

var Timezone = "Asia/Tbilisi"

func GetNowDateTime(timezone *string) time.Time {
	if timezone == nil {
		defaultTz := Timezone
		timezone = &defaultTz
	}

	loc, err := time.LoadLocation(*timezone)

	if err != nil {
		panic(err)
	}

	now := time.Now().In(loc)

	return now
}

func GetFormattedNowDate(timezone *string) string {
	return GetNowDateTime(timezone).Format("2006-01-02")
}

func GetTimeFromStringWithServerTimezone(timeStr string, timezone *string) time.Time {
	if timezone == nil {
		defaultTz := Timezone
		timezone = &defaultTz
	}

	loc, err := time.LoadLocation(*timezone)
	if err != nil {
		panic(err)
	}

	now := time.Now().In(loc)

	parsed, err := time.ParseInLocation("15:04", timeStr, loc)
	if err != nil {
		panic(err)
	}

	return time.Date(
		now.Year(), now.Month(), now.Day(),
		parsed.Hour(), parsed.Minute(), 0, 0, loc,
	)
}

func ConvertTimeToTbilisi(timeStr, userTZ string) (time.Time, error) {
	userLoc, err := time.LoadLocation(userTZ)
	if err != nil {
		return time.Time{}, err
	}
	tbilisiLoc, err := time.LoadLocation(Timezone)
	if err != nil {
		return time.Time{}, err
	}

	now := time.Now().In(userLoc)
	parsed, err := time.ParseInLocation("15:04", timeStr, userLoc)
	if err != nil {
		return time.Time{}, err
	}

	userTime := time.Date(now.Year(), now.Month(), now.Day(),
		parsed.Hour(), parsed.Minute(), 0, 0, userLoc)

	return userTime.In(tbilisiLoc), nil
}

func IsValidTimezone(name string) bool {
	_, err := time.LoadLocation(name)
	return err == nil
}
