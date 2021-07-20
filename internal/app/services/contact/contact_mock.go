package contactsrv

import (
	"github.com/stretchr/testify/mock"
	"goim-pro/internal/app/models"
)

type MockContactRepo struct {
	mock.Mock
}

func (m *MockContactRepo) IsContactExists(userId, contactId string) (isExist bool, err error) {
	args := m.Called(userId, contactId)
	return args.Bool(0), args.Error(1)
}

func (m *MockContactRepo) FindOne(condition map[string]interface{}) (contact *models.Contact, err error) {
	args := m.Called(condition)
	return args.Get(0).(*models.Contact), args.Error(1)
}

func (m *MockContactRepo) InsertContacts(newContacts ...*models.Contact) (err error) {
	args := m.Called(newContacts)
	return args.Error(0)
}

func (m *MockContactRepo) RemoveContactsByIds(userId string, contactIds ...string) (err error) {
	args := m.Called(userId, contactIds)
	return args.Error(0)
}

func (m *MockContactRepo) FindOneAndUpdateRemark(condition interface{}, remarkInfo interface{}) (err error) {
	args := m.Called(condition, remarkInfo)
	return args.Error(0)
}

func (m *MockContactRepo) FindAll(condition map[string]interface{}) ([]models.Contact, error) {
	args := m.Called(condition)
	arg1 := args.Get(0)
	if arg1 == nil {
		return nil, args.Error(1)
	}
	return arg1.([]models.Contact), args.Error(1)
}
