package cerrors

import "errors"

var (
	ErrAlreadyExists    = errors.New("record with the passed parameters already exists")
	ErrRecordNotFound   = errors.New("record with the passed parameters is not found")
	ErrInvalidCondition = errors.New("invalid select condition")
)
