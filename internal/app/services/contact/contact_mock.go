package contactsrv

import (
	"github.com/stretchr/testify/mock"
	. "goim-pro/internal/app/repos/contact"
)

type MockContactRepo struct {
	mock.Mock
}

func (m *MockContactRepo) IsExistContact(userId, contactId string) (isExist bool, err error) {
	args := m.Called(userId, contactId)
	return args.Bool(0), args.Error(1)
}

func (m *MockContactRepo) FindOne(condition *Contact) (contact *Contact, err error) {
	args := m.Called(condition)
	return args.Get(0).(*Contact), args.Error(1)
}

func (m *MockContactRepo) InsertContacts(newContacts ...*Contact) (err error) {
	args := m.Called(newContacts)
	return args.Error(0)
}

func (m *MockContactRepo) RemoveContactsByIds(userId string, contactIds ...string) (err error) {
	args := m.Called(userId, contactIds)
	return args.Error(0)
}

func (m *MockContactRepo) FindOneAndUpdateRemark(ct *Contact, remarkInfo map[string]interface{}) (err error) {

	return nil
}

func (m *MockContactRepo) FindAll(condition map[string]interface{}) ([]Contact, error) {
	return nil, nil
}
