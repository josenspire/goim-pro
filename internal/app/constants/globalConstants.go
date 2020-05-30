package consts

type CodeType int32

const (
	StatusInactive = "INACTIVE"
	StatusActive   = "ACTIVE"

	CodeTypeRegister      CodeType = 1
	CodeTypeLogin         CodeType = 2
	CodeTypeResetPassword CodeType = 3

	MaximumNumberOfGroups       int = 10
	MaximumNumberOfGroupMembers int = 500
)
