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

		mysqlDbEnableLogMode := GetMysqlDbEnableLogMode()
		convey.So(mysqlDbEnableLogMode, convey.ShouldBeTrue)
	})
}
