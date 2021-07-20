package wuid

import "goim-pro/vendor/github.com/stretchr/testify/mock"

type MockWuid struct {
	mock.Mock
}

func (m *MockWuid) NewWUID() string {
	args := m.Called()
	return args.String(0)
}
