package enums

type UserStatuses string

const (
	UserStatusIdle     UserStatuses = "Idle"
	UserStatusEditing  UserStatuses = "Editing"
	UserStatusInactive UserStatuses = "Inactive"
)
