package usersrv

import (
	"github.com/stretchr/testify/mock"
	"goim-pro/internal/app/repos/user"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) IsTelephoneRegistered(telephone string) (bool, error) {
	args := m.Called(telephone)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepo) Register(newUser *user.User) error {
	newUser.UserID = newUser.Telephone
	args := m.Called(newUser)
	return args.Error(0)
}

func (m *MockUserRepo) RemoveUserByUserID(userID string, isForce bool) error {
	panic("implement me")
}
