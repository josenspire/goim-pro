package contactsrv

import (
	"goim-pro/internal/app/models"
	"goim-pro/internal/app/repos/user"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	mockUser = &models.User{
		Password: "1234567890",
		UserProfile: models.UserProfile{
			UserId:      "TEST002",
			Telephone:   "13631210001",
			Email:       "123456@mai.com",
			Nickname:    "JAMES01",
			Avatar:      "",
			Description: "Good!",
			Sex:         "1",
			Birthday:    0,
			Location:    "",
		},
	}
	mockContact = &models.Contact{
		UserId:    "TEST001",
		ContactId: "TEST002",
	}
)

func Test_contactService_RequestContact(t *testing.T) {
	mu := &user.MockUserRepo{}
	mu.On("FindByUserId", "TEST002").Return(mockUser, nil)

	mc := &MockContactRepo{}
	mc.On("IsContactExists", "TEST001", "TEST002").Return(true, nil)

	cs := new(ContactService)
	userRepo = mu
	contactRepo = mc

	Convey("Test_DeleteContact", t, func() {
		Convey("should_found_and_send_notification_to_contact", func() {
			var userId = "TEST002"
			var contactId = "TEST002"
			var reqReason = "请求添加好友！"

			tErr := cs.RequestContact(userId, contactId, reqReason)
			So(tErr.Code, ShouldBeNil)
		})
	})
}

func Test_contactService_DeleteContact(t *testing.T) {
	mu := &user.MockUserRepo{}

	mc := &MockContactRepo{}
	mc.On("IsContactExists", "TEST001", "TEST002").Return(true, nil)

	mc.On("RemoveContactsByIds", "TEST001", "TEST002").Return(nil)
	mc.On("RemoveContactsByIds", "TEST002", "TEST001").Return(nil)

	cs := new(ContactService)
	userRepo = mu
	contactRepo = mc

	Convey("Test_DeleteContact", t, func() {
		Convey("should_delete_successfully", func() {
			var userId = "TEST001"
			var contactId = "TEST002"

			tErr := cs.DeleteContact(userId, contactId)
			So(tErr, ShouldBeNil)
		})
	})
}

func Test_contactService_GetContacts(t *testing.T) {
	//m := &MockContactRepo{}

	//var expectationContacts = make([]contact.ContactImpl, 3)
	//expectationContacts := [...]contact.ContactImpl{}
	//m.On("FindAll", map[string]interface{}{"UserId": "TEST002"}).Return(nil)
	//
	//Convey("Test_DeleteContact", t, func() {
	//	var ctx context.Context
	//	var req *protos.GrpcReq
	//	cs := &ContactService{
	//		userRepo:    mu,
	//		contactRepo: mc,
	//	}
	//	Convey("should_find_all_user's_contacts_then_return", func() {
	//	})
	//})
}

func Test_requestContactParameterCalibration(t *testing.T) {

}
