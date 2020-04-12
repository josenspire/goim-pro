package usersrv

import (
	"github.com/stretchr/testify/mock"
	"goim-pro/internal/app/models"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) IsTelephoneOrEmailRegistered(telephone string, email string) (bool, error) {
	args := m.Called(telephone, email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepo) Register(newUser *models.User) error {
	newUser.UserId = newUser.Telephone
	args := m.Called(newUser)
	return args.Error(0)
}

func (m *MockUserRepo) RemoveUserByUserId(userId string, isForce bool) error {
	args := m.Called(userId, isForce)
	return args.Error(0)
}

func (m *MockUserRepo) QueryByTelephoneAndPassword(telephone string, enPassword string) (user *models.User, err error) {
	args := m.Called(telephone, enPassword)
	return (args.Get(0)).(*models.User), args.Error(1)
}

func (m *MockUserRepo) QueryByEmailAndPassword(email string, enPassword string) (user *models.User, err error) {
	args := m.Called(email, enPassword)
	return (args.Get(0)).(*models.User), args.Error(1)
}

func (m *MockUserRepo) ResetPassword(email string, enPassword string) (user *models.User, err error) {
	args := m.Called(email, enPassword)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepo) ResetPasswordByTelephone(telephone string, newPassword string) error {
	args := m.Called(telephone, newPassword)
	return args.Error(0)
}

func (m *MockUserRepo) ResetPasswordByEmail(email string, newPassword string) error {
	args := m.Called(email, newPassword)
	return args.Error(0)
}

func (m *MockUserRepo) FindByUserId(userId string) (user *models.User, err error) {
	args := m.Called(userId)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepo) FindOneUser(condition map[string]interface{}) (user *models.User, err error) {
	args := m.Called(condition)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepo) FindOneAndUpdateProfile(condition, profile map[string]interface{}) (err error) {
	args := m.Called(condition, profile)
	return args.Error(0)
}
