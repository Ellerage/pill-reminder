package utils

import (
	"time"
)

func GetNowDateTbilisi() string {
	// TODO: get timezone from args
	loc, err := time.LoadLocation("Asia/Tbilisi")
	if err != nil {
		panic(err)
	}

	// TODO: return without format
	now := time.Now().In(loc).Format("2006-01-02")

	return now
}

func GetNowTimeTbilisi() string {
	// TODO: get timezone from args
	loc, err := time.LoadLocation("Asia/Tbilisi")

	if err != nil {
		panic(err)
	}

	// TODO: return without format
	now := time.Now().In(loc).Format("15:04")

	return now
}
