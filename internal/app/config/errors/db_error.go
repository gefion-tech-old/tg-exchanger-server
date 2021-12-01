package errors

import "errors"

var (
	ErrAlreadyRegistered = errors.New("user with this chat_id is already registered")
	ErrNotRegistered     = errors.New("user with this chat_id is not registered")
)
