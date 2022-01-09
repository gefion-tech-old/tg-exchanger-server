package cerrors

import "errors"

var (
	ErrAlreadyRegistered = errors.New("user with this chat_id is already registered")
	ErrNotRegistered     = errors.New("user with this chat_id is not registered")
	ErrInvalidBody       = errors.New("could not read data from the request body")
	ErrInvalidPathParams = errors.New("params in path is invalid")
	ErrNotEnoughRights   = errors.New("not enough rights to make this request")
	ErrTokenInvalid      = errors.New("token is invalid")
	ErrUnauthorized      = errors.New("unauthorized")
)
