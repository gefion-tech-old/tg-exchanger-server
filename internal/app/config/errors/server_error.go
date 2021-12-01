package errors

import "errors"

var (
	ErrInvalidBody     = errors.New("could not read data from the request body")
	ErrNotEnoughRights = errors.New("not enough rights to make this request")
	ErrTokenInvalid    = errors.New("token is invalid")
)
