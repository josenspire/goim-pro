package usersrv

import (
	"github.com/stretchr/testify/mock"
	. "goim-pro/internal/app/repos/user"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) IsTelephoneOrEmailRegistered(telephone string, email string) (bool, error) {
	args := m.Called(telephone, email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepo) Register(newUser *User) error {
	newUser.UserId = newUser.Telephone
	args := m.Called(newUser)
	return args.Error(0)
}

func (m *MockUserRepo) RemoveUserByUserId(userId string, isForce bool) error {
	panic("implement me")
}

func (m *MockUserRepo) QueryByTelephoneAndPassword(telephone string, enPassword string) (user *User, err error) {
	args := m.Called(telephone, enPassword)
	return (args.Get(0)).(*User), args.Error(1)
}

func (m *MockUserRepo) QueryByEmailAndPassword(email string, enPassword string) (user *User, err error) {
	args := m.Called(email, enPassword)
	return (args.Get(0)).(*User), args.Error(1)
}

func (m *MockUserRepo) ResetPassword(email string, enPassword string) (user *User, err error) {
	panic("implement me")
}

func (m *MockUserRepo) ResetPasswordByTelephone(telephone string, newPassword string) error {
	args := m.Called(telephone, newPassword)
	return args.Error(0)

}

func (m *MockUserRepo) ResetPasswordByEmail(email string, newPassword string) error {
	args := m.Called(email, newPassword)
	return args.Error(0)
}

func (m *MockUserRepo) FindByUserId(userId string) (user *User, err error) {
	args := m.Called(userId)
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockUserRepo) FindOneUser(us *User) (user *User, err error) {
	args := m.Called(us)
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockUserRepo) FindOneAndUpdateProfile(us *User, profile map[string]interface{}) (err error) {
	args := m.Called(us, profile)
	return args.Error(0)
}
