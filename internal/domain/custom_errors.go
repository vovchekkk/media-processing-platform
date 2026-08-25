package domain

import "errors"

var ErrInvalidCredentials = errors.New("invalid credentials")

var ErrFailedToGenerateUUID = errors.New("failed to generate uuid")
var ErrTaskNotFound = errors.New("task not found")
