package utils

import "errors"

var ErrInvalidTimeEditInput = errors.New("invalid input")
var ErrUserAlreadyExist = errors.New("user already exist")
var ErrAlreadyTakenToday = errors.New("pill already taken")
var ErrInvalidCommand = errors.New("invalid command")
