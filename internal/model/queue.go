package model

type QueueReminder struct {
	ChatId       int64
	CronId       string
	ReminderType string // Daily, Followup
}

type DailyReminderPayload struct {
	RemindInterval string
	ChatId         int64
}

type FollowUpReminderPayload struct {
	ChatId int64
}

type GetAllQueueReminderFilters struct {
	ChatId       *int64
	ReminderType *string
}

type DeleteFilters struct {
	ReminderType *string
}
