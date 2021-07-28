package contact

import (
	"fmt"
	"github.com/jinzhu/gorm"
	. "github.com/smartystreets/goconvey/convey"
	protos "goim-pro/api/protos/salty"
	"goim-pro/internal/app/models"
	. "goim-pro/internal/app/repos/user"
	"goim-pro/internal/app/services/converters"
	mysqlsrv "goim-pro/internal/db/mysql"
	"testing"
)

var (
	db *gorm.DB
	userRepo IUserRepo
	ctRepo IContactRepo
)

func init() {
	mysqlsrv.NewMysql()
	db = mysqlsrv.GetMysql()
	userRepo = NewUserRepo(db)
	ctRepo = NewContactRepo(db)
}

func TestContactImpl_IsExistContact(t *testing.T) {
	_ = userRepo.Register(&models.User{
		Password:    "123",
		UserProfile: converters.ConvertProto2EntityForUserProfile(&protos.UserProfile{UserId: "TEST001", Telephone: "1234"}),
	})
	_ = userRepo.Register(&models.User{
		Password:    "123",
		UserProfile: converters.ConvertProto2EntityForUserProfile(&protos.UserProfile{UserId: "TEST002", Telephone: "12345"}),
	})
	Convey("Test_IsExistContact", t, func() {
		newContact := &models.Contact{
			UserId:    "TEST001",
			ContactId: "TEST002",
		}
		err := ctRepo.InsertContacts(newContact)
		if err != nil {
			t.FailNow()
		}
		Convey("should_return_true_result", func() {
			isExist, err := ctRepo.IsContactExists("TEST001", "TEST002")
			ShouldBeNil(err)
			ShouldBeTrue(isExist)
		})
	})

	_ = ctRepo.RemoveContactsByIds("TEST001", "TEST002")
	_ = userRepo.RemoveUserByUserId("TEST001", true) // remove demo user
	_ = userRepo.RemoveUserByUserId("TEST002", true) // remove demo user
}

func TestContactImpl_FindOne(t *testing.T) {
	_ = userRepo.Register(&models.User{
		Password:    "123",
		UserProfile: converters.ConvertProto2EntityForUserProfile(&protos.UserProfile{UserId: "TEST001", Telephone: "1234"}),
	})
	_ = userRepo.Register(&models.User{
		Password:    "123",
		UserProfile: converters.ConvertProto2EntityForUserProfile(&protos.UserProfile{UserId: "TEST002", Telephone: "12345"}),
	})
	Convey("Test_FindOne", t, func() {
		newContact := &models.Contact{
			UserId:    "TEST001",
			ContactId: "TEST002",
		}

		ct := &ContactImpl{}
		err := ct.InsertContacts(newContact)
		if err != nil {
			t.FailNow()
		}
		Convey("should_return_one_record_as_result", func() {
			contact, err := ct.FindOne(map[string]interface{}{
				"UserId":    "TEST001",
				"ContactId": "TEST002",
			})
			ShouldBeNil(err)
			ShouldNotBeNil(contact)
		})
	})

	_ = ctRepo.RemoveContactsByIds("TEST001", "TEST002")
	_ = userRepo.RemoveUserByUserId("TEST001", true) // remove demo user
	_ = userRepo.RemoveUserByUserId("TEST002", true) // remove demo user
}

func TestContact_FindAll(t *testing.T) {
	_ = userRepo.Register(&models.User{
		Password:    "123",
		UserProfile: converters.ConvertProto2EntityForUserProfile(&protos.UserProfile{UserId: "TEST001", Telephone: "1234"}),
	})
	_ = userRepo.Register(&models.User{
		Password:    "123",
		UserProfile: converters.ConvertProto2EntityForUserProfile(&protos.UserProfile{UserId: "TEST002", Telephone: "12345"}),
	})
	Convey("Test_FindAll", t, func() {
		newContact1 := &models.Contact{
			UserId:    "TEST001", // 13631210003
			ContactId: "TEST002", // 13631210007
		}
		newContact2 := &models.Contact{
			UserId:    "TEST001",
			ContactId: "TEST003", // 13631210008
		}

		ct := &ContactImpl{}
		err := ct.InsertContacts(newContact1, newContact2)
		if err != nil {
			t.FailNow()
		}
		Convey("should_find_all_contacts_with_profile", func() {
			condition := map[string]interface{}{
				"UserId": "TEST001",
			}
			contacts, err := ct.FindAll(condition)
			fmt.Print(contacts)
			ShouldBeNil(err)
			So(len(contacts), ShouldEqual, 2)
		})
	})

	_ = ctRepo.RemoveContactsByIds("TEST001", "TEST002")
	_ = userRepo.RemoveUserByUserId("TEST001", true) // remove demo user
	_ = userRepo.RemoveUserByUserId("TEST002", true) // remove demo user
}

func TestContact_InsertContacts(t *testing.T) {
	_ = userRepo.Register(&models.User{
		Password:    "123",
		UserProfile: converters.ConvertProto2EntityForUserProfile(&protos.UserProfile{UserId: "TEST001", Telephone: "1234"}),
	})
	_ = userRepo.Register(&models.User{
		Password:    "123",
		UserProfile: converters.ConvertProto2EntityForUserProfile(&protos.UserProfile{UserId: "TEST002", Telephone: "12345"}),
	})
	_ = userRepo.Register(&models.User{
		Password:    "12345",
		UserProfile: converters.ConvertProto2EntityForUserProfile(&protos.UserProfile{UserId: "TEST003", Telephone: "123456"}),
	})
	Convey("Test_InsertContacts", t, func() {
		newContact1 := &models.Contact{
			UserId:    "TEST001", // 13631210003
			ContactId: "TEST002", // 13631210007
		}
		newContact2 := &models.Contact{
			UserId:    "TEST001",
			ContactId: "TEST003", // 13631210008
		}
		err := ctRepo.InsertContacts(newContact1, newContact2)

		Convey("should_insert_multiple_contacts", func() {
			ShouldBeNil(err)
		})
	})

	_ = ctRepo.RemoveContactsByIds("TEST001", "TEST002", "TEST003")
	_ = userRepo.RemoveUserByUserId("TEST001", true) // remove demo user
	_ = userRepo.RemoveUserByUserId("TEST002", true) // remove demo user
}

func TestContact_RemoveContactsByIds(t *testing.T) {
	_ = userRepo.Register(&models.User{
		Password:    "123",
		UserProfile: converters.ConvertProto2EntityForUserProfile(&protos.UserProfile{UserId: "TEST001", Telephone: "1234"}),
	})
	_ = userRepo.Register(&models.User{
		Password:    "123",
		UserProfile: converters.ConvertProto2EntityForUserProfile(&protos.UserProfile{UserId: "TEST002", Telephone: "12345"}),
	})
	_ = userRepo.Register(&models.User{
		Password:    "12345",
		UserProfile: converters.ConvertProto2EntityForUserProfile(&protos.UserProfile{UserId: "TEST003", Telephone: "123456"}),
	})
	Convey("Test_RemoveContacts", t, func() {
		newContact1 := &models.Contact{
			UserId:    "TEST001",
			ContactId: "TEST002",
		}
		newContact2 := &models.Contact{
			UserId:    "TEST001",
			ContactId: "TEST003",
		}
		newContact3 := &models.Contact{
			UserId:    "TEST003",
			ContactId: "TEST004",
		}
		err := ctRepo.InsertContacts(newContact1, newContact2, newContact3)
		if err != nil {
			fmt.Println(err)
		}
		Convey("should_remove_contacts_by_contactIds", func() {
			err := ctRepo.RemoveContactsByIds("TEST001", "TEST002", "TEST003")
			ShouldBeNil(err)

			contact, err := ctRepo.FindOne(map[string]interface{}{"UserId": "TEST001"})
			ShouldBeNil(err)
			ShouldBeNil(contact)

			contact2, err := ctRepo.FindOne(map[string]interface{}{"UserId": "TEST003"})
			ShouldBeNil(err)
			So(contact2.UserId, ShouldEqual, "TEST003")

			_ = ctRepo.RemoveContactsByIds("TEST003", "TEST004")

			err = ctRepo.RemoveContactsByIds("TEST003", "TEST004")

			ShouldBeNil(err)
		})
	})
	_ = userRepo.RemoveUserByUserId("TEST001", true) // remove demo user
	_ = userRepo.RemoveUserByUserId("TEST002", true) // remove demo user
}

func TestContact_FindOneAndUpdateRemark(t *testing.T) {
	_ = userRepo.Register(&models.User{
		Password:    "123",
		UserProfile: converters.ConvertProto2EntityForUserProfile(&protos.UserProfile{UserId: "TEST001", Telephone: "12345"}),
	})
	_ = userRepo.Register(&models.User{
		Password:    "12345",
		UserProfile: converters.ConvertProto2EntityForUserProfile(&protos.UserProfile{UserId: "TEST002", Telephone: "123456"}),
	})
	Convey("Test_FindOneAndUpdateRemark", t, func() {
		newContact1 := &models.Contact{
			UserId:    "TEST001",
			ContactId: "TEST002",
		}
		_ = ctRepo.InsertContacts(newContact1)
		Convey("should_update_contact_remark_profile_successfully", func() {
			criteria := map[string]interface{}{
				"UserId":    "TEST001",
				"ContactId": "TEST002",
			}

			remarkProfile := map[string]interface{}{
				"RemarkName":  "JAMES001",
				"Telephone":   "13631210001;13631210001",
				"Description": "Crazy boy..",
				"Tags":        "Friend;Boy",
			}
			err := ctRepo.FindOneAndUpdateRemark(criteria, remarkProfile)
			ShouldBeNil(err)

			result, _ := ctRepo.FindOne(criteria)
			So(result.RemarkName, ShouldEqual, "JAMES001")
			So(result.Telephone, ShouldEqual, "13631210001;13631210001")
			So(result.Description, ShouldEqual, "Crazy boy..")
			So(result.Tags, ShouldEqual, "Friend;Boy")
		})
	})

	_ = ctRepo.RemoveContactsByIds("TEST001", "TEST002")
	_ = userRepo.RemoveUserByUserId("TEST001", true) // remove demo user
	_ = userRepo.RemoveUserByUserId("TEST002", true) // remove demo user
}
