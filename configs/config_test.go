package configs

import (
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGetAppConfigProfile(t *testing.T) {
	convey.Convey("Subject: Get application profile", t, func() {
		appHost := GetAppHost()
		convey.So(appHost, convey.ShouldEqual, "127.0.0.1")

		appPort := GetAppPort()
		convey.So(appPort, convey.ShouldEqual, "9090")
	})
}

func TestGetMysqlConfigProfile(t *testing.T) {
	convey.Convey("Subject: Get mysql db profile", t, func() {
		mysqlDbUri := GetMysqlDbUri()
		convey.So(mysqlDbUri, convey.ShouldEqual, "127.0.0.1")

		mysqlDbPort := GetMysqlDbPort()
		convey.So(mysqlDbPort, convey.ShouldEqual, "3306")

		mysqlDbName := GetMysqlDbName()
		convey.So(mysqlDbName, convey.ShouldEqual, "goim")
	})
}


//func TestGetAppConfigProfile(t *testing.T) {
//	Convey("Subject: Get application configuration profile", t, func() {
//		appProfile, _ := GetAppConfigProfile()
//		fmt.Printf("read profile: %v\n", *appProfile)
//		So(appProfile.AppHost, ShouldEqual, "127.0.0.1")
//		So(appProfile.AppPort, ShouldEqual, 9090)
//		So(appProfile.AppLogLevel, ShouldEqual, "DEBUG")
//	})
//}
//
//func TestGetMysqlConfigProfile(t *testing.T) {
//	Convey("Subject Get mysql configuration profile", t, func() {
//		mysqlProfile, _ := GetMysqlConfigProfile()
//		fmt.Printf("read profile: %v\n", *mysqlProfile)
//		So(mysqlProfile.DBUri, ShouldEqual, "127.0.0.1")
//		So(mysqlProfile.DBPort, ShouldEqual, 3306)
//		So(mysqlProfile.DBName, ShouldEqual, "goim")
//	})
//}
