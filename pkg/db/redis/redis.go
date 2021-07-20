package redsrv

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"goim-pro/config"
	"goim-pro/pkg/logs"
	"sync"
)

var logger = logs.GetLogger("INFO")
var redisOnce sync.Once
var redisClient IMyRedis

var (
	host     = "0.0.0.0"
	port     = "6767"
	password = ""
	dbNum    = 1
)

func NewRedis() IMyRedis {
	redisOnce.Do(func() {
		redisClient = connect()
	})
	return redisClient
}

func connect() *BaseClient {
	host = config.GetRedisDBHost()
	port = config.GetRedisDBPort()
	password = config.GetRedisDBPassword()
	dbNum = config.GetRedisDBNum()
	//dbKey = config.GetRedisDBKey()

	addrs := fmt.Sprintf("%s:%s", host, port)
	opts := &redis.Options{
		Addr:         addrs,
		DB:           dbNum,
		Password:     password,
		PoolSize:     10,
		MinIdleConns: 0,
		MaxConnAge:   0,
		PoolTimeout:  0,
		IdleTimeout:  0,
	}
	baseClient := newBaseClient(opts)
	return baseClient
}
