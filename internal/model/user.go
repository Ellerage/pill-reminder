package model

import "pill-reminder/internal/utils/enums"

type User struct {
	ChatId         int64  `db:"chatId"`
	Timezone       string `db:"timezone"`
	TimeToNotify   string `db:"timeToNotify"` // String in 15:04 format
	Status         string `db:"status"`
	RemindInterval int64  `db:"remindInterval"` // Time in minutes
}

type UserUpdate struct {
	Timezone       *string `db:"timezone,omitempty"`
	TimeToNotify   *string `db:"timeToNotify,omitempty"`
	Status         *string `db:"status,omitempty"`
	RemindInterval *int64  `db:"remindInterval, omitempty"`
}

type UserCreate struct {
	Timezone       string `db:"timezone"`
	TimeToNotify   string `db:"timeToNotify"`
	Status         string `db:"status"`
	RemindInterval int64  `db:"remindInterval"`
}

type UserNotificationSettings struct {
	Timezone       string `db:"timezone"`
	TimeToNotify   string `db:"timeToNotify"`
	RemindInterval int64  `db:"remindInterval"`
}

func (u *UserCreate) GetDefaultUser(timezone string) UserCreate {
	return UserCreate{
		Timezone:       timezone,
		TimeToNotify:   "00:00",
		Status:         string(enums.UserStatusInactive),
		RemindInterval: 20,
	}
}
