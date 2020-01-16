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
	defaultAppHost     = "0.0.0.0"
	defaultAppPort     = "9090"
	defaultAppLogLevel = "DEBUG"

	defaultMysqlUserName      = "root"
	defaultMysqlPassword      = "Password1!"
	defaultMysqlURI           = "0.0.0.0"
	defaultMysqlPort          = "3306"
	defaultMysqlName          = "goim"
	defaultMysqlEngine        = "InnoDB"
	defaultMysqlMaxIdleConns  = "10"
	defaultMysqlMaxOpenConns  = "30"
	defaultMysqlEnableLogMode = true

	defaultAppSecretKey = "U0FMVFktSU0tUFJP"

	defaultRedisHost     = "0.0.0.0"
	defaultRedisPort     = "6379"
	defaultRedisPassword = ""
	defaultRedisNum      = "1"
	defaultRedisKey      = "SaltyIMPro"
)

/* app config **/
func GetAppHost() (appHost string) {
	appHost = appConfigMap["apphost"]
	return buildConfigParameter(appHost, defaultAppHost)
}
func GetAppPort() (appPort string) {
	appPort = appConfigMap["appport"]
	return buildConfigParameter(appPort, defaultAppPort)
}
func GetAppLogLevel() (appLogLevel string) {
	appLogLevel = appConfigMap["apploglevel"]
	return buildConfigParameter(appLogLevel, defaultAppLogLevel)
}

/* api config **/
func GetApiSecretKey() (apiSecretKey string) {
	apiSecretKey = viper.MyViper.GetString("apisecretkey")
	if len(apiSecretKey) == 0 {
		apiSecretKey = defaultAppSecretKey
	}
	return apiSecretKey
}

/* app config **/
func GetMysqlDBUserName() (dbUserName string) {
	dbUserName = mysqlConfigMap["dbusername"]
	return buildConfigParameter(dbUserName, defaultMysqlUserName)
}
func GetMysqlDBPassword() (dbPassword string) {
	dbPassword = mysqlConfigMap["dbpassword"]
	return buildConfigParameter(dbPassword, defaultMysqlPassword)
}
func GetMysqlDBUri() (dbURI string) {
	dbURI = mysqlConfigMap["dburi"]
	return buildConfigParameter(dbURI, defaultMysqlURI)
}
func GetMysqlDBPort() (dbPort string) {
	dbPort = mysqlConfigMap["dbport"]
	return buildConfigParameter(dbPort, defaultMysqlPort)
}
func GetMysqlDBName() (dbName string) {
	dbName = mysqlConfigMap["dbname"]
	return buildConfigParameter(dbName, defaultMysqlName)
}
func GetMysqlDBEngine() (dbEngine string) {
	dbEngine = mysqlConfigMap["dbengine"]
	return buildConfigParameter(dbEngine, defaultMysqlEngine)
}
func GetMysqlDBMaxIdleConns() (dbMaxIdleConns string) {
	dbMaxIdleConns = mysqlConfigMap["dbmaxidleconns"]
	return buildConfigParameter(dbMaxIdleConns, defaultMysqlMaxIdleConns)
}
func GetMysqlDBMaxOpenConns() (dbMaxOpenConns string) {
	dbMaxOpenConns = mysqlConfigMap["dbmaxopenconns"]
	return buildConfigParameter(dbMaxOpenConns, defaultMysqlMaxOpenConns)
}
func GetMysqlDBEnableLogMode() bool {
	dbEnableLogModeStr := mysqlConfigMap["dbenablelogmode"]
	return buildBoolParameter(dbEnableLogModeStr, defaultMysqlEnableLogMode)
}

/* app config **/
func GetRedisDBHost() (redisHost string) {
	redisHost = redisConfigMap["dbhost"]
	return buildConfigParameter(redisHost, defaultRedisHost)
}
func GetRedisDBPort() (redisPort string) {
	redisPort = redisConfigMap["dbport"]
	return buildConfigParameter(redisPort, defaultRedisPort)
}
func GetRedisDBPassword() (redisPassword string) {
	redisPassword = redisConfigMap["dbpassword"]
	return buildConfigParameter(redisPassword, defaultRedisPassword)
}
func GetRedisDBNum() (redisDBNum int) {
	redisDBStr := redisConfigMap["dbnum"]
	redisDBNum, _ = strconv.Atoi(buildConfigParameter(redisDBStr, defaultRedisNum))
	return
}
func GetRedisDBKey() (redisDBKey string) {
	redisDBKey = redisConfigMap["dbkey"]
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
