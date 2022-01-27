package cerrors

import "errors"

var (
	ErrFailedToGetAllMerchants = errors.New("failed to get list of all merchants")
	ErrFailedToDecodeParams    = errors.New("failed to decode merchant optional params")
)
