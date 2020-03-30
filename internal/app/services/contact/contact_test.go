package contactsrv

import (
	usersrv "goim-pro/internal/app/services/user"
	"testing"
)

func Test_contactService_DeleteContact(t *testing.T) {
	mu := &usersrv.MockUserRepo{}
	mu.On("FindByUserId", "TEST001").Return(nil, nil)

	m := &MockContactRepo{}
	m.On("IsExistContact", "TEST001", "TEST002").Return(true, nil)

	m.On("RemoveContactsByIds", "TEST001", "TEST002").Return(nil)
	m.On("RemoveContactsByIds", "TEST002", "TEST001").Return(nil)
}
