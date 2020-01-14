package redsrv

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"goim-pro/config"
	"goim-pro/pkg/logs"
	"log"
	"os"
	"sync"
	"time"
)

type RedisConnectionPool struct {
}

var (
	host     string = ""
	port     string = ""
	password string = ""
	dbNum    int    = 1
	dbKey    string = ""
)

var logger = logs.GetLogger("INFO")
var redisOnce sync.Once
var redisInstance *RedisConnectionPool
var client *redis.Client

func init() {
	host = config.GetRedisDBHost()
	port = config.GetRedisDBPort()
	password = config.GetRedisDBPassword()
	dbNum = config.GetRedisDBNum()
	dbKey = config.GetRedisDBKey()

	redis.SetLogger(log.New(os.Stderr, "redis: ", log.LstdFlags))
}

func NewRedisConnection() *RedisConnectionPool {
	redisOnce.Do(func() {
		redisInstance = &RedisConnectionPool{}
	})
	return redisInstance
}

func (rs *RedisConnectionPool) Connect() (err error) {
	uriAddr := fmt.Sprintf("%s:%s", host, port)
	client = redis.NewClient(&redis.Options{
		Addr:        uriAddr,
		Password:    password,
		DB:          dbNum,
		DialTimeout: time.Second * 10,
	})
	_, err = client.Ping().Result()
	if err == nil {
		logger.Infoln("[redis] pong successful")
	}
	return
}

func (rs *RedisConnectionPool) GetRedisClient() *redis.Client {
	return client
}
