package auth

import (
	"github.com/stretchr/testify/mock"
	"goim-pro/api/protos/salty"
	"goim-pro/internal/app/models/errors"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) ObtainSMSCode(telephone string, operationType com_salty_protos.SMSOperationType) (code string, tErr *cserr.TError) {
	args := m.Called(telephone, operationType)
	if tErr := args.Get(1); tErr != nil {
		return args.String(0), args.Get(1).(*cserr.TError)
	}
	return args.String(0), nil
}

func (m *MockAuthService) VerifySMSCode(telephone string, operationType com_salty_protos.SMSOperationType, codeStr string) (isPass bool, tErr *cserr.TError) {
	args := m.Called(telephone, operationType, codeStr)
	if tErr := args.Get(1); tErr != nil {
		return args.Bool(0), args.Get(1).(*cserr.TError)
	}
	return args.Bool(0), nil
}
