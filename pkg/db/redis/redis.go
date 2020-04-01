package redsrv

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"goim-pro/config"
	"goim-pro/pkg/logs"
	"log"
	"os"
	"sync"
)

type RedisConnectionPool struct {
}

var logger = logs.GetLogger("INFO")
var redisOnce sync.Once
var redisInstance *RedisConnectionPool
var client *BaseClient

var (
	host     string = "0.0.0.0"
	port     string = "6767"
	password string = ""
	dbNum    int    = 1
	dbKey    string = "SaltyIMPro"
)

func NewRedisConnection() *RedisConnectionPool {
	redisOnce.Do(func() {
		redisInstance = &RedisConnectionPool{}
	})
	return redisInstance
}

func (rs *RedisConnectionPool) Connect() (err error) {
	host = config.GetRedisDBHost()
	port = config.GetRedisDBPort()
	password = config.GetRedisDBPassword()
	dbNum = config.GetRedisDBNum()
	dbKey = config.GetRedisDBKey()

	redis.SetLogger(log.New(os.Stderr, "redis: ", log.LstdFlags))

	uriAddr := fmt.Sprintf("%s:%s", host, port)
	client = NewBaseClient(uriAddr, password, dbNum)
	_, err = client.Ping()
	if err == nil {
		logger.Info("[redis] pong successfully!")
	}
	return
}

func (rs *RedisConnectionPool) GetRedisClient() *BaseClient {
	return client
}
