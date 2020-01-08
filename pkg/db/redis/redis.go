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

type RedisService struct {
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
var redisInstance *RedisService
var client *redis.Client

func init() {
	host = config.GetRedisDBHost()
	port = config.GetRedisDBPort()
	password = config.GetRedisDBPassword()
	dbNum = config.GetRedisDBNum()
	dbKey = config.GetRedisDBKey()

	redis.SetLogger(log.New(os.Stderr, "redis: ", log.LstdFlags))
}

func NewRedisService() *RedisService {
	redisOnce.Do(func() {
		redisInstance = &RedisService{}
	})
	return redisInstance
}

func (rs *RedisService) Connect() (err error) {
	uriAddr := fmt.Sprintf("%s:%s", host, port)
	client = redis.NewClient(&redis.Options{
		Addr:     uriAddr,
		Password: password,
		DB:       dbNum,
	})
	_, err = client.Ping().Result()
	if err != nil {
		logger.Errorf("[redis] ping redis fail: %v", err)
	} else {
		logger.Infoln("[redis] pong successful")
	}
	return
}

func GetRedisClient() *redis.Client {
	return client
}