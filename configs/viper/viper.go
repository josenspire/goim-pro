package viper

import (
	"github.com/spf13/viper"
	"goim-pro/pkg/logs"
)

var MyViper *viper.Viper
var logger = logs.GetLogger("INFO")

func init() {
	MyViper = viper.New()
	//MyViper.AddConfigPath("../configs")
	MyViper.SetConfigFile("application.json")
	MyViper.SetConfigType("json")

	err := MyViper.ReadInConfig()
	if err != nil {
		logger.Info("application config loaded...")
	} else {
		logger.Infof("error reading application config file: %v\n", err)
	}
}
