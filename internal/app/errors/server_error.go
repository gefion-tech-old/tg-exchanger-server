package errors

import "errors"

var (
	ErrInvalidBody       = errors.New("could not read data from the request body")
	ErrInvalidPathParams = errors.New("params in path is invalid")
	ErrNotEnoughRights   = errors.New("not enough rights to make this request")
	ErrTokenInvalid      = errors.New("token is invalid")
	ErrUnauthorized      = errors.New("unauthorized")

	ErrFailedToInitializeStruct = errors.New("failed to initialize structure")
)
