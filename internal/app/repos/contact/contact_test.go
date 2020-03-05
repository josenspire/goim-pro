package contact

import (
	. "github.com/smartystreets/goconvey/convey"
	mysqlsrv "goim-pro/pkg/db/mysql"
	"testing"
)

func TestContact_InsertContacts(t *testing.T) {
	mysqlDB := mysqlsrv.NewMysqlConnection()
	_ = mysqlDB.Connect()
	NewContactRepo(mysqlsrv.NewMysqlConnection().GetMysqlInstance())

	newContact1 := &Contact{
		UserId:        "TEST001",
		ContactId:     "TEST002",
	}
	newContact2 := &Contact{
		UserId:        "TEST003",
		ContactId:     "TEST004",
	}

	ct := &Contact{}
	err := ct.InsertContacts(newContact1, newContact2)

	Convey("Test_InsertContacts", t, func() {
		Convey("should_insert_multiple_contacts", func() {
			ShouldBeNil(err)
		})
	})
}
