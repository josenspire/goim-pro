package cserr

type TError struct {
	Code   int32  // error code
	Detail string // error detail message
}

func (e *TError) Error() string {
	return e.Detail
}

func NewTError(code int32, err error) *TError {
	return &TError{
		Code:   code,
		Detail: err.Error(),
	}
}
