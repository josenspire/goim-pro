package config

import (
	viper "goim-pro/config/viper"
	"strconv"
)

var (
	appConfigMap   map[string]string
	mysqlConfigMap map[string]string
	redisConfigMap map[string]string
)

func init() {
	appConfigMap = viper.MyViper.GetStringMapString("appProfile")
	mysqlConfigMap = viper.MyViper.GetStringMapString("mysqlProfile")
	redisConfigMap = viper.MyViper.GetStringMapString("redisProfile")
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
	defaultMysqlEnableLogMode = true

	defaultAppSecretKey = "U0FMVFktSU0tUFJP"

	defaultRedisHost     = "127.0.0.1"
	defaultRedisPort     = "6379"
	defaultRedisPassword = ""
	defaultRedisNum      = "1"
	defaultRedisKey      = "SaltyIMPro"
)

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

/* api config **/
func GetApiSecretKey() (apiSecretKey string) {
	apiSecretKey = viper.MyViper.GetString("apiSecretKey")
	if len(apiSecretKey) == 0 {
		apiSecretKey = defaultAppSecretKey
	}
	return apiSecretKey
}

/* app config **/
func GetMysqlDBUserName() (dbUserName string) {
	dbUserName = mysqlConfigMap["dbUserName"]
	return buildConfigParameter(dbUserName, defaultMysqlUserName)
}
func GetMysqlDBPassword() (dbPassword string) {
	dbPassword = mysqlConfigMap["dbPassword"]
	return buildConfigParameter(dbPassword, defaultMysqlPassword)
}
func GetMysqlDBUri() (dbUri string) {
	dbUri = mysqlConfigMap["dbUri"]
	return buildConfigParameter(dbUri, defaultMysqlUri)
}
func GetMysqlDBPort() (dbPort string) {
	dbPort = mysqlConfigMap["dbPort"]
	return buildConfigParameter(dbPort, defaultMysqlPort)
}
func GetMysqlDBName() (dbName string) {
	dbName = mysqlConfigMap["dbName"]
	return buildConfigParameter(dbName, defaultMysqlName)
}
func GetMysqlDBEngine() (dbEngine string) {
	dbEngine = mysqlConfigMap["dbEngine"]
	return buildConfigParameter(dbEngine, defaultMysqlEngine)
}
func GetMysqlDBMaxIdleConns() (dbMaxIdleConns string) {
	dbMaxIdleConns = mysqlConfigMap["dbMaxIdleConns"]
	return buildConfigParameter(dbMaxIdleConns, defaultMysqlMaxIdleConns)
}
func GetMysqlDBMaxOpenConns() (dbMaxOpenConns string) {
	dbMaxOpenConns = mysqlConfigMap["dbMaxOpenConns"]
	return buildConfigParameter(dbMaxOpenConns, defaultMysqlMaxOpenConns)
}
func GetMysqlDBEnableLogMode() bool {
	dbEnableLogModeStr := mysqlConfigMap["dbEnableLogMode"]
	return buildBoolParameter(dbEnableLogModeStr, defaultMysqlEnableLogMode)
}

/* app config **/
func GetRedisDBHost() (redisHost string) {
	redisHost = redisConfigMap["dbHost"]
	return buildConfigParameter(redisHost, defaultRedisHost)
}
func GetRedisDBPort() (redisPort string) {
	redisPort = redisConfigMap["dbPort"]
	return buildConfigParameter(redisPort, defaultRedisPort)
}
func GetRedisDBPassword() (redisPassword string) {
	redisPassword = redisConfigMap["dbPassword"]
	return buildConfigParameter(redisPassword, defaultRedisPassword)
}
func GetRedisDBNum() (redisDBNum int) {
	redisDBStr := redisConfigMap["dbNum"]
	redisDBNum, _ = strconv.Atoi(buildConfigParameter(redisDBStr, defaultRedisNum))
	return
}
func GetRedisDBKey() (redisDBKey string) {
	redisDBKey = redisConfigMap["dbKey"]
	return buildConfigParameter(redisDBKey, defaultRedisKey)
}

func buildConfigParameter(originalValue string, defaultValue string) string {
	if len(originalValue) > 0 {
		return originalValue
	}
	return defaultValue
}

func buildBoolParameter(originalValue string, defaultValue bool) bool {
	if len(originalValue) > 0 {
		return originalValue == "true"
	}
	return defaultValue
}
