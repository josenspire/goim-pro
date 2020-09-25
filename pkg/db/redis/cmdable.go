// custom cmdable
// extends redis.UniversalClient
// custom implement methods which is start with `R`

package redsrv

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"time"
)

type IMyRedis interface {
	redis.UniversalClient

	RPing() (result string, err error)

	RHSet(key string, valueMap map[string]interface{}) (err error)
	RHGet(key string, fields ...string) (valueMap map[string]interface{}, err error)
	RGet(key string) (strVal string)
	RSet(key string, value string, expiresTime time.Duration) (err error)
	RDel(key string) (resultInt64 int64)
	RHGetAll(key string) (strVal []string, err error)
}

type BaseClient struct {
	redis.Client
}

// new base client with redis options
func newBaseClient(opts *redis.Options) *BaseClient {
	baseClient := &BaseClient{
		*redis.NewClient(opts),
	}
	return baseClient
}

// redis connection testing: ping
func (c *BaseClient) RPing() (result string, err error) {
	return c.Ping().Result()
}

// get single string value by key
func (c *BaseClient) RGet(key string) (strVal string) {
	return c.Get(key).Val()
}

// set single string value by key
func (c *BaseClient) RSet(key string, value string, expiresTime time.Duration) (err error) {
	return c.Set(key, value, expiresTime).Err()
}

// del single record by key, return int64 as result: 0, 1
func (c *BaseClient) RDel(key string) (resultInt64 int64) {
	return c.Del(key).Val()
}

// set hash record, input: key, mapValue; return error
func (c *BaseClient) RHSet(key string, valueMap map[string]interface{}) (err error) {
	for field, value := range valueMap {
		err = c.HSet(key, field, value).Err()
	}
	return
}

// get hash record, input: key; return map value, error
func (c *BaseClient) RHGet(key string, fields ...string) (valueMap map[string]interface{}, err error) {
	valueMap = make(map[string]interface{})
	for _, field := range fields {
		var result interface{}
		val, err := c.HGet(key, fmt.Sprintf("%s", field)).Result()
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

func (c *BaseClient) RHGetAll(key string) (strVal []string, err error) {

	_ = c.HKeys(key)

	return nil, nil
}