package utils

import "errors"

var (
	ErrInvalidTimeEditInput    = errors.New("invalid input")
	ErrUserAlreadyExist        = errors.New("user already exist")
	ErrAlreadyTakenToday       = errors.New("pill already taken")
	ErrInvalidCommand          = errors.New("invalid command")
	ErrInvalidReminderInterval = errors.New("invalid reminder interval")
)
