package converters

import (
	protos "goim-pro/api/protos/salty"
	"goim-pro/internal/app/models"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConvertRegisterUserProfile(t *testing.T) {
	var pbUser1 = &protos.UserProfile{
		Telephone:   "13631210000",
		Email:       "123456@qq.com",
		Nickname:    "JAMES01",
		Avatar:      "www.baidu.com/1.png",
		Description: "Never settle",
		Sex:         0,
		Birthday:    1578903121862,
		Location:    "",
	}
	var expectation = &models.UserProfile{
		Telephone:   "13631210000",
		Email:       "123456@qq.com",
		Nickname:    "JAMES01",
		Avatar:      "www.baidu.com/1.png",
		Description: "Never settle",
		Sex:         "MALE",
		Birthday:    1578903121862,
		Location:    "",
	}
	Convey("Testing_ConvertRegisterUserProfile", t, func() {
		Convey("should return UserProfile entity data", func() {
			actual := ConvertProto2EntityForUserProfile(pbUser1)
			So(actual.Telephone, ShouldEqual, expectation.Telephone)
			So(actual.Sex, ShouldEqual, expectation.Sex)
			So(actual.Birthday, ShouldEqual, expectation.Birthday)
		})
	})
}

func TestConvertLoginResp(t *testing.T) {
	var userProfile = &models.UserProfile{
		Telephone:   "13631210000",
		Email:       "123456@qq.com",
		Nickname:    "JAMES01",
		Avatar:      "www.baidu.com/1.png",
		Description: "Never settle",
		Sex:         "MALE",
		Birthday:    1578903121862,
		Location:    "",
	}
	var pbUser = &protos.UserProfile{
		Telephone:   "13631210000",
		Email:       "123456@qq.com",
		Nickname:    "JAMES01",
		Avatar:      "www.baidu.com/1.png",
		Description: "Never settle",
		Sex:         protos.UserProfile_MALE,
		Birthday:    1578903121862,
		Location:    "",
	}

	Convey("Testing_ConvertLoginResp", t, func() {
		actual := ConvertEntity2ProtoForUserProfile(userProfile)
		So(actual.Telephone, ShouldEqual, pbUser.Telephone)
		So(actual.Sex, ShouldEqual, pbUser.Sex)
	})

}
