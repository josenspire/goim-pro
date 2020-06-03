package errmsg

import "errors"

var (
	// error
	ErrIllegalOperation       = errors.New("illegal operation")
	ErrInvalidParameters      = errors.New("bad request, invalid parameters")
	ErrOperationForbidden     = errors.New("illegal operation, user do not have permission")
	ErrRepeatOperation        = errors.New("do not repeat the operation")
	ErrSystemUncheckException = errors.New("server uncheck exception")

	// user
	ErrTelephoneExists             = errors.New("telephone has been registered")
	ErrTelephoneNotExists          = errors.New("telephone has not been registered")
	ErrAccountSecurityVerification = errors.New("account need security verification")
	ErrAccountAlreadyExists        = errors.New("account is already exists")
	ErrAccountNotExists            = errors.New("account not exists")
	ErrUserNotExists               = errors.New("user not exists")
	ErrVerificationCodeExpired     = errors.New("verification code has expired")
	ErrInvalidVerificationCode     = errors.New("invalid verification code")
	ErrInvalidUserId               = errors.New("invalid user id")
	ErrAccountOrPwdInvalid         = errors.New("account or password is incorrect")
	ErrPwdInvalid                  = errors.New("account's password is incorrect")
	ErrRepeatPassword              = errors.New("the new password cannot be the same as the old one")

	// contact
	ErrContactNotExists       = errors.New("contact not exists")
	ErrContactAlreadyExists   = errors.New("user is already your contact")
	ErrContactRepeatOperation = errors.New("please do not repeat the operation")
	ErrIllegalRequestContacts = errors.New("illegal request contacts")
	ErrInvalidContact         = errors.New("invalid contacts")

	// group
	ErrGroupReachedLimit       = errors.New("groups has reached the limit")
	ErrGroupMemberReachedLimit = errors.New("group's member has reached the limit")
	ErrGroupNotExists          = errors.New("group not exists")
	ErrNotGroupMembers         = errors.New("user did not joined the group")
	ErrRepeatedlyJoinGroup     = errors.New("user has joined the group")
)
