package domain

import "errors"

var ErrInvalidCredentials = errors.New("invalid credentials")

var ErrTaskNotFound = errors.New("task not found")
var ErrSessionNotFound = errors.New("session not found")
var ErrUserNotFound = errors.New("user not found")
