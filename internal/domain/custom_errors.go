package domain

import "errors"

var ErrInvalidCredentials = errors.New("invalid credentials")

var ErrTaskNotFound = errors.New("task not found")
