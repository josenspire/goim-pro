package viper

import (
	"fmt"
	"github.com/spf13/viper"
)

var MyViper *viper.Viper

func init() {
	MyViper = viper.New()
	//MyViper.AddConfigPath("../configs")
	MyViper.SetConfigFile("application.json")
	MyViper.SetConfigType("json")

	err := MyViper.ReadInConfig()
	if err != nil {
		fmt.Println("application config loaded...")
	} else {
		fmt.Printf("error reading application config file: %v\n", err)
	}
}
