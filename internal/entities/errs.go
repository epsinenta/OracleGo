package entities

import "github.com/pkg/errors"

var (
	AlreadyRegisteredError     = errors.New("user already registered")
	SessionAlreadyStartedError = errors.New("session already started")
)
