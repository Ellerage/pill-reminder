package utils

import (
	"time"
)

func GetNowDateTbilisi() string {
	loc, err := time.LoadLocation("Asia/Tbilisi")
	if err != nil {
		panic(err)
	}

	now := time.Now().In(loc).Format("2006-01-02")

	return now
}

func GetNowTimeTbilisi() string {
	loc, err := time.LoadLocation("Asia/Tbilisi")

	if err != nil {
		panic(err)
	}

	now := time.Now().In(loc).Format("15:04")

	return now
}
