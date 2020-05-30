package cserr

import protos "goim-pro/api/protos/salty"

type TError struct {
	Code   protos.StatusCode // error code
	Detail string            // error detail message
}

func (e *TError) Error() string {
	return e.Detail
}

func NewTError(code protos.StatusCode, err error) *TError {
	return &TError{
		Code:   code,
		Detail: err.Error(),
	}
}
