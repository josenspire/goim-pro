package contact

import (
	"fmt"
	"goim-pro/internal/app/models"
	mysqlsrv "goim-pro/pkg/db/mysql"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestContactImpl_IsExistContact(t *testing.T) {
	mysqlDB := mysqlsrv.NewMysqlConnection()
	_ = mysqlDB.Connect()
	NewContactRepo(mysqlsrv.NewMysqlConnection().GetMysqlInstance())

	newContact := &models.Contact{
		UserId:    "TEST001",
		ContactId: "TEST002",
	}

	ct := &ContactImpl{}
	err := ct.InsertContacts(newContact)
	if err != nil {
		t.FailNow()
	}

	Convey("Test_IsExistContact", t, func() {
		Convey("should_return_true_result", func() {
			isExist, err := ct.IsContactExists("TEST001", "TEST002")
			ShouldBeNil(err)
			ShouldBeTrue(isExist)
		})
	})

	_ = ct.RemoveContactsByIds("TEST001", "TEST002")

}

func TestContactImpl_FindOne(t *testing.T) {
	mysqlDB := mysqlsrv.NewMysqlConnection()
	_ = mysqlDB.Connect()
	NewContactRepo(mysqlsrv.NewMysqlConnection().GetMysqlInstance())

	newContact := &models.Contact{
		UserId:    "TEST001",
		ContactId: "TEST002",
	}

	ct := &ContactImpl{}
	err := ct.InsertContacts(newContact)
	if err != nil {
		t.FailNow()
	}

	Convey("Test_FindOne", t, func() {
		Convey("should_return_one_record_as_result", func() {
			contact, err := ct.FindOne(map[string]interface{}{
				"UserId":    "TEST001",
				"ContactId": "TEST002",
			})
			ShouldBeNil(err)
			ShouldNotBeNil(contact)
		})
	})

	_ = ct.RemoveContactsByIds("TEST001", "TEST002")
}

func TestContact_FindAll(t *testing.T) {
	mysqlDB := mysqlsrv.NewMysqlConnection()
	_ = mysqlDB.Connect()
	NewContactRepo(mysqlsrv.NewMysqlConnection().GetMysqlInstance())

	newContact1 := &models.Contact{
		UserId:    "01E07SG858N3CGV5M1APVQKZYR", // 13631210003
		ContactId: "01E2JVWZTG60NG2SXFYNEPNMCB", // 13631210007
	}
	newContact2 := &models.Contact{
		UserId:    "01E07SG858N3CGV5M1APVQKZYR",
		ContactId: "01E2JXMC98SZXMGEGVTDECSD78", // 13631210008
	}

	ct := &ContactImpl{}
	err := ct.InsertContacts(newContact1, newContact2)
	if err != nil {
		t.FailNow()
	}

	Convey("Test_FindAll", t, func() {
		Convey("should_find_all_contacts_with_profile", func() {
			condition := map[string]interface{}{
				"UserId": "01E07SG858N3CGV5M1APVQKZYR",
			}

			contacts, err := ct.FindAll(condition)
			fmt.Print(contacts)
			ShouldBeNil(err)
			So(len(contacts), ShouldEqual, 2)
		})
	})

	_ = ct.RemoveContactsByIds("01E07SG858N3CGV5M1APVQKZYR", "01E2JVWZTG60NG2SXFYNEPNMCB", "01E2JXMC98SZXMGEGVTDECSD78")
}

func TestContact_InsertContacts(t *testing.T) {
	mysqlDB := mysqlsrv.NewMysqlConnection()
	_ = mysqlDB.Connect()
	NewContactRepo(mysqlsrv.NewMysqlConnection().GetMysqlInstance())

	newContact1 := &models.Contact{
		UserId:    "01E07SG858N3CGV5M1APVQKZYR", // 13631210003
		ContactId: "01E2JVWZTG60NG2SXFYNEPNMCB", // 13631210007
	}
	newContact2 := &models.Contact{
		UserId:    "01E07SG858N3CGV5M1APVQKZYR",
		ContactId: "01E2JXMC98SZXMGEGVTDECSD78", // 13631210008
	}

	ct := &ContactImpl{}
	err := ct.InsertContacts(newContact1, newContact2)

	Convey("Test_InsertContacts", t, func() {
		Convey("should_insert_multiple_contacts", func() {
			ShouldBeNil(err)
		})
	})

	_ = ct.RemoveContactsByIds("01E07SG858N3CGV5M1APVQKZYR", "01E2JVWZTG60NG2SXFYNEPNMCB", "01E2JXMC98SZXMGEGVTDECSD78")
}

func TestContact_RemoveContactsByIds(t *testing.T) {
	mysqlDB := mysqlsrv.NewMysqlConnection()
	_ = mysqlDB.Connect()
	NewContactRepo(mysqlsrv.NewMysqlConnection().GetMysqlInstance())

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

	ct := &ContactImpl{}
	err := ct.InsertContacts(newContact1, newContact2, newContact3)
	if err != nil {
		fmt.Println(err)
	}
	Convey("Test_RemoveContacts", t, func() {
		Convey("should_remove_contacts_by_contactIds", func() {
			err := ct.RemoveContactsByIds("TEST001", "TEST002", "TEST003")
			ShouldBeNil(err)

			contact, err := ct.FindOne(map[string]interface{}{"UserId": "TEST001"})
			ShouldBeNil(err)
			ShouldBeNil(contact)

			contact2, err := ct.FindOne(map[string]interface{}{"UserId": "TEST003"})
			ShouldBeNil(err)
			So(contact2.UserId, ShouldEqual, "TEST003")

			_ = ct.RemoveContactsByIds("TEST003", "TEST004")

			err = ct.RemoveContactsByIds("TEST003", "TEST004")

			ShouldBeNil(err)
		})
	})

}

func TestContact_FindOneAndUpdateRemark(t *testing.T) {
	mysqlDB := mysqlsrv.NewMysqlConnection()
	_ = mysqlDB.Connect()
	NewContactRepo(mysqlsrv.NewMysqlConnection().GetMysqlInstance())

	newContact1 := &models.Contact{
		UserId:    "TEST001",
		ContactId: "TEST002",
	}

	ct := &ContactImpl{}
	_ = ct.InsertContacts(newContact1)

	Convey("Test_FindOneAndUpdateRemark", t, func() {
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
			err := ct.FindOneAndUpdateRemark(criteria, remarkProfile)
			ShouldBeNil(err)

			result, _ := ct.FindOne(criteria)
			So(result.RemarkName, ShouldEqual, "JAMES001")
			So(result.Telephone, ShouldEqual, "13631210001;13631210001")
			So(result.Description, ShouldEqual, "Crazy boy..")
			So(result.Tags, ShouldEqual, "Friend;Boy")
		})
	})

	_ = ct.RemoveContactsByIds("TEST001", "TEST002")
}
