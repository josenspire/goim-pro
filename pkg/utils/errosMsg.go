package utils

import "errors"

var (
	// error
	ErrIllegalOperation   = errors.New("illegal operation")
	ErrInvalidParameters  = errors.New("bad request, invalid parameters")
	ErrOperationForbidden = errors.New("illegal operation, user do not have permission")

	ErrAccountOrPwdInvalid     = errors.New("account or password is incorrect")
	ErrPwdInvalid              = errors.New("account's password is incorrect")
	ErrInvalidVerificationCode = errors.New("invalid verification code")
	ErrInvalidUserId           = errors.New("invalid user id")
	ErrUserNotExists           = errors.New("user not exists")

	// contact
	ErrContactNotExists       = errors.New("contact not exists")
	ErrContactAlreadyExists   = errors.New("user are already your contact")
	ErrIllegalRequestContacts = errors.New("illegal request contacts")
	ErrInvalidContact         = errors.New("invalid contacts")

	// group
	ErrGroupReachedLimit       = errors.New("groups has reached the limit")
	ErrGroupMemberReachedLimit = errors.New("group's number of member out of limit")
	ErrGroupNotExists          = errors.New("group not exists")
)
