package viper

import (
	"github.com/spf13/viper"
	"goim-pro/pkg/logs"
	"os"
	"strings"
)

var MyViper *viper.Viper
var logger = logs.GetLogger("INFO")

func init() {
	MyViper = viper.New()
	MyViper.AddConfigPath("./config")
	MyViper.SetConfigName("application")
	MyViper.SetConfigType("json")
	readInConfig(MyViper, "default")

	profileName := os.Getenv("APP_ENV")
	//profileName = "PROD"
	if profileName != "" {
		profileName = strings.ToLower(profileName)
		profileViper := viper.New()
		profileViper.AddConfigPath("./config")
		profileViper.SetConfigName("application-" + profileName)
		readInConfig(profileViper, profileName)

		_ = MyViper.MergeConfigMap(profileViper.AllSettings())
	}
}

func readInConfig(myViper *viper.Viper, appEnv string) {
	err := myViper.ReadInConfig()
	if err != nil {
		logger.Errorf("error in reading application [%s] config file: %s", appEnv, err.Error())
	} else {
		logger.Infof("application [%s] config loaded...", appEnv)
	}
}
