package redsrv

import (
	"github.com/stretchr/testify/mock"
	"time"
)

type MockCmdable struct {
	mock.Mock
}

func (m *MockCmdable) Ping() (result string, err error) {
	panic("implement me")
}

func (m *MockCmdable) HSet(key string, valueMap map[string]interface{}) (err error) {
	panic("implement me")
}

func (m *MockCmdable) HGet(key string, fields ...string) (valueMap map[string]interface{}, err error) {
	panic("implement me")
}

func (m *MockCmdable) Get(key string) (strVal string) {
	args := m.Called(key)
	return args.String(0)
}

func (m *MockCmdable) Set(key string, value string, expiresTime time.Duration) (err error) {
	args := m.Called(key, value, expiresTime)
	return args.Error(0)
}

func (m *MockCmdable) Del(key string) (resultInt64 int64) {
	args := m.Called(key)
	return int64(args.Int(0))
}
