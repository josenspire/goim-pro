package converters

import (
	. "github.com/smartystreets/goconvey/convey"
	"goim-pro/internal/app/models"
	"goim-pro/internal/app/repos/base"
	"testing"
)

var groupProfile = &models.Group{
	GroupId:     "01E8BK7PVG6833B7MVG85CA1K2",
	CreatedBy:   "01E07SG858N3CGV5M1APVQKZYR",
	OwnerUserId: "01E07SG858N3CGV5M1APVQKZYR",
	Name:        "NEW_GROUP",
	Avatar:      "",
	Notice:      "Never Settle",
	Members:     nil,
	BaseModel:   base.BaseModel{},
}
var member1 = models.Member{
	UserId:  "01E59Z8HMG8SK8C65XV42M33QP",
	Alias:   "",
	Role:    "1",
	Status:  "NORMAL",
	GroupId: "01E8BK7PVG6833B7MVG85CA1K2",
	User: models.User{
		Password: "bWWsxgDSeAyJmNm+50mL3g==",
		Role:     "1",
		Status:   "ACTIVE",
		UserProfile: models.UserProfile{
			UserId:      "01E59Z8HMG8SK8C65XV42M33QP",
			Telephone:   "13631210001",
			Email:       "12345671@qq.com",
			Nickname:    "JAMES01",
			Avatar:      "https://www.baidu.com/avatar/header1.png",
			Description: "Never settle",
			Sex:         "MALE",
			Birthday:    1586251450678,
			Location:    "CHINA-ZHA",
		},
	},
}
var member2 = models.Member{
	UserId:  "01E59ZNYB8KDNW0W3NHGDZDD6V",
	Alias:   "",
	Role:    "1",
	Status:  "NORMAL",
	GroupId: "01E8BK7PVG6833B7MVG85CA1K2",
	User: models.User{
		Password: "bWWsxgDSeAyJmNm+50mL3g==",
		Role:     "1",
		Status:   "ACTIVE",
		UserProfile: models.UserProfile{
			UserId:      "01E59ZNYB8KDNW0W3NHGDZDD6V",
			Telephone:   "13631210003",
			Email:       "12345673@qq.com",
			Nickname:    "JAMES03",
			Avatar:      "https://www.baidu.com/avatar/header1.png",
			Description: "Never settle",
			Sex:         "MALE",
			Birthday:    1586251889316,
			Location:    "CHINA-ZHA",
		},
	},
}

func TestConvertEntity2ProtoForGroupProfile(t *testing.T) {
	groupProfile.Members = []models.Member{member1, member2}

	Convey("Test_ConvertEntity2ProtoForGroupProfile", t, func() {
		Convey("should_convert_group_&_member_&_userProfile_then_return_pb_entity_without_password", func() {
			pbProfile := ConvertEntity2ProtoForGroupProfile(groupProfile)
			So(pbProfile.Members[0].GroupId, ShouldEqual, "01E8BK7PVG6833B7MVG85CA1K2")
			So(pbProfile.Members[1].GroupId, ShouldEqual, "01E8BK7PVG6833B7MVG85CA1K2")

			So(pbProfile.Members[0].UserProfile.UserId, ShouldEqual, "01E59Z8HMG8SK8C65XV42M33QP")
			So(pbProfile.Members[1].UserProfile.UserId, ShouldEqual, "01E59ZNYB8KDNW0W3NHGDZDD6V")

			So(pbProfile.Members[0].UserProfile.Telephone, ShouldEqual, "13631210001")
			So(pbProfile.Members[1].UserProfile.Email, ShouldEqual, "12345673@qq.com")
		})
	})
}
