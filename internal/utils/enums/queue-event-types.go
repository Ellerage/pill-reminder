package enums

type ReminderType string
type QueueEventsEnum string

const (
	ReminderEventDaily    QueueEventsEnum = "reminder:daily"
	ReminderEventFollowup QueueEventsEnum = "reminder:followup"
	ReminderEventDelayed  QueueEventsEnum = "reminder:delayed"
)

const (
	ReminderTypeDaily    ReminderType = "Daily"
	ReminderTypeFollowup ReminderType = "Followup"
)
