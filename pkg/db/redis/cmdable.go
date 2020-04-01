package redsrv

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"time"
)

type IBaseCmdable interface {
	Ping() (result string, err error)
	HSet(key string, valueMap map[string]interface{}) (err error)
	Get(key string) (strVal string)
	Set(key string, value string, expiresTime time.Duration) (err error)
	Del(key string) (resultInt64 int64)
}

type BaseClient struct {
	client *redis.Client
}

// new base client with redis options
func NewBaseClient(uriAddr string, password string, dbNum int) *BaseClient {
	baseClient := &BaseClient{}
	baseClient.client = redis.NewClient(
		&redis.Options{
			Addr:        uriAddr,
			Password:    password,
			DB:          dbNum,
			DialTimeout: time.Second * 10,
		})
	return baseClient
}

// redis connection testing: ping
func (bc *BaseClient) Ping() (result string, err error) {
	return bc.client.Ping().Result()
}

// get single string value by key
func (bc *BaseClient) Get(key string) (strVal string) {
	return bc.client.Get(key).Val()
}

// set single string value by key
func (bc *BaseClient) Set(key string, value string, expiresTime time.Duration) (err error) {
	return bc.client.Set(key, value, expiresTime).Err()
}

// del single record by key, return int64 as result: 0, 1
func (bc *BaseClient) Del(key string) (resultInt64 int64) {
	return bc.client.Del(key).Val()
}

// set hash record, input: key, mapValue; return error
func (bc *BaseClient) HSet(key string, valueMap map[string]interface{}) (err error) {
	for field, value := range valueMap {
		err = bc.client.HSet(key, field, value).Err()
	}
	return
}

// get hash record, input: key; return map value, error
func (bc *BaseClient) HGet(key string, fields ...string) (valueMap map[string]interface{}, err error) {
	valueMap = make(map[string]interface{})
	for _, field := range fields {
		var result interface{}
		val, err := bc.client.HGet(key, fmt.Sprintf("%s", field)).Result()
		if err == redis.Nil {
			valueMap[field] = result
			err = nil
		} else if err != nil {
			logger.Errorf("get hash record error: %s", err.Error())
			valueMap[field] = result
		}
		if val != "" {
			valueMap[field] = val
		} else {
			valueMap[field] = result
		}
	}
	return valueMap, err
}
