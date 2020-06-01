package contactsrv

import (
	"github.com/stretchr/testify/mock"
	protos "goim-pro/api/protos/salty"
	consts "goim-pro/internal/app/constants"
	"goim-pro/internal/app/models"
	"goim-pro/internal/app/repos/user"
	mysqlsrv "goim-pro/pkg/db/mysql"
	redsrv "goim-pro/pkg/db/redis"
	errmsg "goim-pro/pkg/errors"
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
	mu.On("FindByUserId", "TEST003").Return(nil, nil)
	mu.On("FindByUserId", "TEST004").Return(mockUser, nil)

	mc := &MockContactRepo{}
	mc.On("IsContactExists", "TEST001", "TEST002").Return(false, nil)
	mc.On("IsContactExists", "TEST001", "TEST004").Return(true, nil)

	r := new(redsrv.MockCmdable)
	reqContactNotifyKey := "CT-REQ-TEST001-TEST002"
	r.On("RSet", reqContactNotifyKey, mock.Anything, consts.DefaultExpiresTime).Return(nil)

	cs := new(ContactService)
	userRepo = mu
	contactRepo = mc
	myRedis = r

	Convey("Test_DeleteContact", t, func() {
		Convey("should_found_and_send_notification_to_contact_success", func() {
			var userId = "TEST001"
			var contactId = "TEST002"
			var reqReason = "请求添加好友！"

			tErr := cs.RequestContact(userId, contactId, reqReason)
			So(tErr, ShouldBeNil)
		})
		Convey("should_not_found_and_return_error_when_given_unExists_contactId", func() {
			var userId = "TEST001"
			var contactId = "TEST003"
			var reqReason = "请求添加好友！"

			tErr := cs.RequestContact(userId, contactId, reqReason)
			So(tErr, ShouldNotBeNil)
			So(tErr.Code, ShouldEqual, protos.StatusCode_STATUS_BAD_REQUEST)
		})
		Convey("should_not_found_and_return_error_when_contact_already_exists", func() {
			var userId = "TEST001"
			var contactId = "TEST004"
			var reqReason = "请求添加好友！"

			tErr := cs.RequestContact(userId, contactId, reqReason)
			So(tErr, ShouldNotBeNil)
			So(tErr.Code, ShouldEqual, protos.StatusCode_STATUS_BAD_REQUEST)
			So(tErr.Detail, ShouldEqual, errmsg.ErrContactAlreadyExists.Error())
		})
	})
}

func Test_contactService_DeleteContact(t *testing.T) {
	mysqlDB = mysqlsrv.NewMysql()

	mu := &user.MockUserRepo{}
	mu.On("FindByUserId", "TEST002").Return(mockUser, nil)
	mu.On("FindByUserId", "TEST003").Return(mockUser, nil)

	mc := &MockContactRepo{}
	mc.On("IsContactExists", "TEST001", "TEST002").Return(true, nil)
	mc.On("IsContactExists", "TEST001", "TEST003").Return(false, nil)

	mc.On("RemoveContactsByIds", mock.Anything, mock.Anything).Return(nil)

	r := new(redsrv.MockCmdable)
	delContactNotifyKey := "CT-DEL-TEST001-TEST002"
	r.On("RSet", delContactNotifyKey, mock.Anything, consts.DefaultExpiresTime).Return(nil)

	cs := new(ContactService)
	userRepo = mu
	contactRepo = mc
	myRedis = r

	Convey("Test_DeleteContact", t, func() {
		Convey("should_remove_contact_successful_when_given_correct_contactId", func() {
			var userId = "TEST001"
			var contactId = "TEST002"

			tErr := cs.DeleteContact(userId, contactId)
			So(tErr, ShouldBeNil)
		})
		Convey("should_remove_contact_failed_when_given_unExists_contactId", func() {
			var userId = "TEST001"
			var contactId = "TEST003"

			tErr := cs.DeleteContact(userId, contactId)
			So(tErr, ShouldNotBeNil)
			So(tErr.Code, ShouldEqual, protos.StatusCode_STATUS_BAD_REQUEST)
			So(tErr.Detail, ShouldEqual, errmsg.ErrContactNotExists.Error())
		})
	})
}

func Test_contactService_GetContacts(t *testing.T) {
	mysqlDB = mysqlsrv.NewMysql()

	mc := &MockContactRepo{}
	var contacts = make([]models.Contact, 5)
	mc.On("FindAll", map[string]interface{}{"userId": "TEST001"}).Return(contacts, nil)
	mc.On("FindAll", map[string]interface{}{"userId": "TEST002"}).Return(nil, nil)

	cs := new(ContactService)
	contactRepo = mc

	Convey("Test_DeleteContact", t, func() {
		Convey("should_find_all_user's_contacts_then_return", func() {
			contacts, tErr := cs.GetContacts("TEST001")
			So(tErr, ShouldBeNil)
			So(len(contacts), ShouldEqual, 5)
		})
		Convey("should_not_found_contacts_when_given_invalid_userId", func() {
			contacts, tErr := cs.GetContacts("TEST002")
			So(tErr, ShouldBeNil)
			So(contacts, ShouldBeNil)
		})
	})
}

func Test_requestContactParameterCalibration(t *testing.T) {

}
