package utils

import "errors"

var (
	ErrInvalidParameters = errors.New("invalid parameters")

	ErrAccountOrPswInvalid = errors.New("account or password is incorrect")
)
