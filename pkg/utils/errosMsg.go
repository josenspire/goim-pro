package utils

import "errors"

var (
	ErrInvalidParameters = errors.New("invalid parameters")

	ErrAccountOrPwdInvalid = errors.New("account or password is incorrect")

	ErrPwdInvalid = errors.New("account's password is incorrect")
)
