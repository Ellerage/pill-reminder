package utils

import (
	"time"
)

type GetDateTimeOptions struct {
	Timezone *string
}

type GetFormattedNowDateTimeOptions struct {
	GetDateTimeOptions
	Format *string
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
