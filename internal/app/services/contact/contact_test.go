package contactsrv

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	protos "goim-pro/api/protos/salty"
	usersrv "goim-pro/internal/app/services/user"

	"goim-pro/pkg/utils"
	"testing"
)

func Test_contactService_DeleteContact(t *testing.T) {
	mu := &usersrv.MockUserRepo{}

	mc := &MockContactRepo{}
	mc.On("IsExistContact", "TEST001", "TEST002").Return(true, nil)

	mc.On("RemoveContactsByIds", "TEST001", "TEST002").Return(nil)
	mc.On("RemoveContactsByIds", "TEST002", "TEST001").Return(nil)

	Convey("Test_DeleteContact", t, func() {
		var ctx context.Context
		var req *protos.GrpcReq
		cs := &contactService{
			userRepo: mu,
			contactRepo: mc,
		}

		Convey("should_delete_successfully", func() {
			delContactReq := &protos.DeleteContactReq{
				UserId: "TEST002",
			}
			any, _ := utils.MarshalMessageToAny(delContactReq)
			req = &protos.GrpcReq{
				Data: any,
				Token: "TEST001",
			}
			actualResp, err := cs.DeleteContact(ctx, req)
			ShouldBeNil(err)
			So(actualResp.Message, ShouldEqual, "contact deleted successfully")
		})
	})
}
