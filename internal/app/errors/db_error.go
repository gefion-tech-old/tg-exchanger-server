package errors

import "errors"

var (
	ErrAlreadyRegistered = errors.New("user with this chat_id is already registered")
	ErrNotRegistered     = errors.New("user with this chat_id is not registered")

	ErrAlreadyExists    = errors.New("record with the passed parameters already exists")
	ErrRecordNotFound   = errors.New("record with the passed parameters is not found")
	ErrInvalidCondition = errors.New("invalid select condition")
)
