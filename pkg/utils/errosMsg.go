package utils

import "errors"

var (
	// error
	ErrInvalidParameters = errors.New("invalid parameters")

	ErrAccountOrPwdInvalid = errors.New("account or password is incorrect")

	ErrPwdInvalid = errors.New("account's password is incorrect")

	ErrInvalidUserId = errors.New("invalid user id")

	ErrUserNotExists = errors.New("user not exists")
)
