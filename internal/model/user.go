package model

type User struct {
	ChatId       int64  `bson:"chatId"`
	Timezone     string `bson:"timezone"`
	TimeToNotify string `bson:"timeToNotify"`
}
