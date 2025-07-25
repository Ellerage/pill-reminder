package model

import "pill-reminder/internal/utils/enums"

type User struct {
	ChatId         int64  `bson:"chatId"`
	Timezone       string `bson:"timezone"`
	TimeToNotify   string `bson:"timeToNotify"` // 15:04
	Status         string `bson:"status"`
	RemindInterval string `bson:"remindInterval"` // "* * * * *"
}

type UserUpdate struct {
	Timezone       *string `bson:"timezone,omitempty"`
	TimeToNotify   *string `bson:"timeToNotify,omitempty"`
	Status         *string `bson:"status,omitempty"`
	RemindInterval *string `bson:"remindInterval, omitempty"`
}

type UserCreate struct {
	Timezone       string `bson:"timezone"`
	TimeToNotify   string `bson:"timeToNotify"`
	Status         string `bson:"status"`
	RemindInterval uint8  `bson:"remindInterval"`
}

func (u *UserCreate) GetDefaultUser(timezone string) UserCreate {
	return UserCreate{
		Timezone:       timezone,
		TimeToNotify:   "00:00",
		Status:         string(enums.UserStatusInactive),
		RemindInterval: 20,
	}
}
