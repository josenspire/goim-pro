package redsrv

import (
	"github.com/stretchr/testify/mock"
	"time"
)

type MockCmdable struct {
	mock.Mock
}

func (m *MockCmdable) RHGetAll(key string) (strVal []string, err error) {
	panic("implement me")
}

func (m *MockCmdable) RPing() (result string, err error) {
	panic("implement me")
}

func (m *MockCmdable) RHSet(key string, valueMap map[string]interface{}) (err error) {
	args := m.Called(key, valueMap)
	return args.Error(0)
}

func (m *MockCmdable) RHGet(key string, fields ...string) (valueMap map[string]interface{}, err error) {
	args := m.Called(key, fields)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockCmdable) RGet(key string) (strVal string) {
	args := m.Called(key)
	return args.String(0)
}

func (m *MockCmdable) RSet(key string, value string, expiresTime time.Duration) (err error) {
	args := m.Called(key, value, expiresTime)
	return args.Error(0)
}

func (m *MockCmdable) RDel(key string) (resultInt64 int64) {
	args := m.Called(key)
	return int64(args.Int(0))
}
