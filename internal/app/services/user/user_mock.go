package usersrv

import (
	"github.com/stretchr/testify/mock"
	"goim-pro/internal/app/repos/user"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) IsTelephoneOrEmailRegistered(telephone string, email string) (bool, error) {
	args := m.Called(telephone, email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepo) Register(newUser *user.User) error {
	newUser.UserId = newUser.Telephone
	args := m.Called(newUser)
	return args.Error(0)
}

func (m *MockUserRepo) RemoveUserByUserId(userId string, isForce bool) error {
	panic("implement me")
}

func (m *MockUserRepo) LoginByTelephone(telephone string, enPassword string) (user *user.User, err error) {
	panic("implement me")
}

func (m *MockUserRepo) LoginByEmail(email string, enPassword string) (user *user.User, err error) {
	panic("implement me")
}
