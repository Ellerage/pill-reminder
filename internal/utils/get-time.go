package utils

import (
	"time"
)

type GetDateTimeFromOptions struct {
	Hours    *int
	Timezone *string
}

type GetFormattedNowDateTimeOptions struct {
	Timezone *string
	Format   *string
}

func GetNowDateTime(timezone *string) time.Time {
	if timezone == nil {
		defaultTz := "Asia/Tbilisi"
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

func GetFormattedNowTime(timezone *string) string {
	return GetNowDateTime(timezone).Format("15:04")
}

func GetFormattedNowDateTime(options GetFormattedNowDateTimeOptions) string {
	return GetNowDateTime(options.Timezone).Format(*options.Format)
}

func GetDateTimeFrom(options GetDateTimeFromOptions) time.Time {
	var timezone string

	if options.Timezone == nil {
		timezone = "Asia/Tbilisi"
	} else {
		timezone = *options.Timezone
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		panic(err)
	}

	now := time.Now().In(loc)

	return time.Date(
		now.Year(), now.Month(), now.Day(),
		*options.Hours, 0, 0, 0, loc,
	)

}

func ConvertTimeToTbilisi(timeStr, userTZ string) (time.Time, error) {
	userLoc, err := time.LoadLocation(userTZ)
	if err != nil {
		return time.Time{}, err
	}
	tbilisiLoc, err := time.LoadLocation("Asia/Tbilisi")
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
