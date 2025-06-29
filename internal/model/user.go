package model

type User struct {
	ChatId       int64  `bson:"chatId"`
	Timezone     string `bson:"timezone"`
	TimeToNotify string `bson:"timeToNotify"`
	Status       string `bson:"status"`
}

type UserUpdate struct {
	Timezone     *string `bson:"timezone,omitempty"`
	TimeToNotify *string `bson:"timeToNotify,omitempty"`
	Status       *string `bson:"status,omitempty"`
}
