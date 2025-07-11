package model

type QueueReminder struct {
	ChatId       int64  `bson:"chatId"`
	CronId       string `bson:"cronId"`
	ReminderType string `bson:"type"` // Daily, Followup
}

type DailyReminderPayload struct {
	RemindInterval string
	ChatId         int64
}

type FollowUpReminderPayload struct {
	ChatId int64
}

type GetAllFilters struct {
	ChatId       *int64
	ReminderType *string
}

type DeleteFilters struct {
	ReminderType *string
}
