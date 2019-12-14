package configs

import (
	viper "goim-pro/configs/viper"
)

var (
	appConfigMap   map[string]string
	mysqlConfigMap map[string]string
)

func init() {
	appConfigMap = viper.MyViper.GetStringMapString("appProfile")
	mysqlConfigMap = viper.MyViper.GetStringMapString("mysqlProfile")
}

var (
	defaultAppHost     = "127.0.0.1"
	defaultAppPort     = "9090"
	defaultAppLogLevel = "DEBUG"

	defaultMysqlUserName      = "root"
	defaultMysqlPassword      = "Password1!"
	defaultMysqlUri           = "127.0.0.1"
	defaultMysqlPort          = "3306"
	defaultMysqlName          = "goim"
	defaultMysqlEngine        = "InnoDB"
	defaultMysqlMaxIdleConns  = "10"
	defaultMysqlMaxOpenConns  = "30"
	defaultMysqlEnableLogMode = "true"
)

type AppProfile struct {
	AppHost     string
	AppPort     int
	AppLogLevel string
}

type MysqlProfile struct {
	DBUserName      string
	DBPassword      string
	DBUri           string
	DBPort          int
	DBName          string
	DBEngine        string
	DBMaxIdleConns  int
	DBMaxOpenConns  int
	DBEnableLogMode bool
}

/* app config **/
func GetAppHost() (appHost string) {
	appHost = appConfigMap["appHost"]
	return buildConfigParameter(appHost, defaultAppHost)
}
func GetAppPort() (appPort string) {
	appPort = appConfigMap["appPort"]
	return buildConfigParameter(appPort, defaultAppPort)
}
func GetAppLogLevel() (appLogLevel string) {
	appLogLevel = appConfigMap["appLogLevel"]
	return buildConfigParameter(appLogLevel, defaultAppLogLevel)
}

/* app config **/
func GetMysqlDbUserName() (dbUserName string) {
	dbUserName = mysqlConfigMap["dbUserName"]
	return buildConfigParameter(dbUserName, defaultMysqlUserName)
}
func GetMysqlDbPassword() (dbPassword string) {
	dbPassword = mysqlConfigMap["dbPassword"]
	return buildConfigParameter(dbPassword, defaultMysqlPassword)
}
func GetMysqlDbUri() (dbUri string) {
	dbUri = mysqlConfigMap["dbUri"]
	return buildConfigParameter(dbUri, defaultMysqlUri)
}
func GetMysqlDbPort() (dbPort string) {
	dbPort = mysqlConfigMap["dbPort"]
	return buildConfigParameter(dbPort, defaultMysqlPort)
}
func GetMysqlDbName() (dbName string) {
	dbName = mysqlConfigMap["dbName"]
	return buildConfigParameter(dbName, defaultMysqlName)
}
func GetMysqlDbEngine() (dbEngine string) {
	dbEngine = mysqlConfigMap["dbEngine"]
	return buildConfigParameter(dbEngine, defaultMysqlEngine)
}
func GetMysqlDbMaxIdleConns() (dbMaxIdleConns string) {
	dbMaxIdleConns = mysqlConfigMap["dbMaxIdleConns"]
	return buildConfigParameter(dbMaxIdleConns, defaultMysqlMaxIdleConns)
}
func GetMysqlDbMaxOpenConns() (dbMaxOpenConns string) {
	dbMaxOpenConns = mysqlConfigMap["dbMaxOpenConns"]
	return buildConfigParameter(dbMaxOpenConns, defaultMysqlMaxOpenConns)
}
func GetMysqlDbEnableLogMode() (dbEnableLogMode string) {
	dbEnableLogMode = mysqlConfigMap["dbEnableLogMode"]
	return buildConfigParameter(dbEnableLogMode, defaultMysqlEnableLogMode)
}

func buildConfigParameter(originalValue string, defaultValue string) string {
	if len(originalValue) > 0 {
		return originalValue
	}
	return defaultValue
}

//func _GetAppConfigProfile() (appProfile *AppProfile, err error) {
//	appConfigMap := myViper.GetStringMap("appProfile")
//
//	bytes, err := json.Marshal(appConfigMap)
//	if err != nil {
//		fmt.Printf("marshal app config fail: %v\n", err)
//		return nil, err
//	}
//	err = json.Unmarshal(bytes, &appProfile)
//	if err != nil {
//		fmt.Printf("unmarshal app profile fail: %v\n", err)
//		return nil, err
//	}
//	return
//}
//
//func _GetMysqlConfigProfile() (mysqlProfile *MysqlProfile, err error) {
//	mysqlConfigMap := myViper.GetStringMap("mysqlProfile")
//
//	bytes, err := json.Marshal(mysqlConfigMap)
//	if err != nil {
//		fmt.Printf("marshal mysql config fail: %v\n", err)
//		return nil, err
//	}
//	err = json.Unmarshal(bytes, &mysqlProfile)
//	if err != nil {
//		fmt.Printf("unmarshal mysql profile fail: %v\n", err)
//		return nil, err
//	}
//	return
//}