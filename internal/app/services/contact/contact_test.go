package contactsrv

import (
	"context"
	protos "goim-pro/api/protos/salty"
	"goim-pro/internal/app/models"
	usersrv "goim-pro/internal/app/services/user"
	"goim-pro/pkg/utils"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	mockUser = &models.User{
		Password:    "1234567890",
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
	mu := &usersrv.MockUserRepo{}
	mu.On("FindByUserId", "TEST002").Return(mockUser, nil)

	mc := &MockContactRepo{}
	mc.On("IsContactExists", "TEST001", "TEST002").Return(true, nil)

	Convey("Test_DeleteContact", t, func() {
		var ctx context.Context
		var req *protos.GrpcReq
		cs := &contactService{
			userRepo:    mu,
			contactRepo: mc,
		}
		Convey("should_found_and_send_notification_to_contact", func() {
			reqContactReq := &protos.RequestContactReq{
				UserId: "TEST002",
				Reason: "请求添加好友！",
			}
			any, _ := utils.MarshalMessageToAny(reqContactReq)
			req = &protos.GrpcReq{
				Data: any,
			}
			actualResp, err := cs.RequestContact(ctx, req)
			So(err, ShouldBeNil)
			So(actualResp.Code, ShouldEqual, 200)
		})
	})
}

func Test_contactService_DeleteContact(t *testing.T) {
	mu := &usersrv.MockUserRepo{}

	mc := &MockContactRepo{}
	mc.On("IsContactExists", "TEST001", "TEST002").Return(true, nil)

	mc.On("RemoveContactsByIds", "TEST001", "TEST002").Return(nil)
	mc.On("RemoveContactsByIds", "TEST002", "TEST001").Return(nil)

	Convey("Test_DeleteContact", t, func() {
		var ctx context.Context
		var req *protos.GrpcReq
		cs := &contactService{
			userRepo:    mu,
			contactRepo: mc,
		}

		Convey("should_delete_successfully", func() {
			delContactReq := &protos.DeleteContactReq{
				UserId: "TEST002",
			}
			any, _ := utils.MarshalMessageToAny(delContactReq)
			req = &protos.GrpcReq{
				Data:  any,
				Token: "TEST001",
			}
			actualResp, err := cs.DeleteContact(ctx, req)
			ShouldBeNil(err)
			So(actualResp.Message, ShouldEqual, "contact deleted successfully")
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
	//	cs := &contactService{
	//		userRepo:    mu,
	//		contactRepo: mc,
	//	}
	//	Convey("should_find_all_user's_contacts_then_return", func() {
	//	})
	//})
}

func Test_requestContactParameterCalibration(t *testing.T) {
	type args struct {
		userId string
		req    *protos.RequestContactReq
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := requestContactParameterCalibration(tt.args.userId, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("requestContactParameterCalibration() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
