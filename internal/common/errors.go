package common

import "errors"

var (
	ErrCoinNotFound  = errors.New("coin not found")
	ErrPriceNotFound = errors.New("response not found")
)
