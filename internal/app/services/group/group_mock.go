package groupsrv

import (
	"github.com/stretchr/testify/mock"
	"goim-pro/internal/app/models"
)

type MockGroup struct {
	mock.Mock
}

func (m *MockGroup) CreateGroup(groupProfile *models.Group) (newGroup *models.Group, err error) {
	args := m.Called(groupProfile)
	arg1 := args.Get(0)
	err = args.Error(1)
	if arg1 == nil {
		return nil, err
	}
	return arg1.(*models.Group), err
}

func (m *MockGroup) RemoveGroupByGroupId(groupId string, isForce bool) (err error) {
	args := m.Called(groupId, isForce)
	return args.Error(0)
}

func (m *MockGroup) FindOneGroup(condition interface{}) (groupProfile *models.Group, err error) {
	args := m.Called(condition)
	arg1 := args.Get(0)
	err = args.Error(1)
	if arg1 == nil {
		return nil, err
	}
	return arg1.(*models.Group), err
}

func (m *MockGroup) CountGroup(condition interface{}) (count int, err error) {
	args := m.Called(condition)
	return args.Int(0), args.Error(1)
}

func (m *MockGroup) FindOneGroupAndUpdate(condition interface{}, updated interface{}) (newProfile *models.Group, err error) {
	args := m.Called(condition, updated)
	arg1 := args.Get(0)
	err = args.Error(1)
	if arg1 == nil {
		return nil, err
	}
	return arg1.(*models.Group), err
}

func (m *MockGroup) FindOneMember(condition interface{}) (memberProfile *models.Member, err error) {
	args := m.Called(condition)
	arg1 := args.Get(0)
	err = args.Error(1)
	if arg1 == nil {
		return nil, err
	}
	return arg1.(*models.Member), err
}

func (m *MockGroup) InsertMembers(members ...*models.Member) (err error) {
	args := m.Called(members)
	return args.Error(0)
}

func (m *MockGroup) RemoveMembers(groupId string, memberIds []string, isForce bool) (deleteCount int64, err error) {
	args := m.Called(groupId, memberIds, isForce)
	return int64(args.Int(0)), args.Error(1)
}

func (m *MockGroup) FindOneMemberAndUpdate(condition interface{}, updated interface{}) (newProfile *models.Member, err error) {
	args := m.Called(condition, updated)
	arg1 := args.Get(0)
	err = args.Error(1)
	if arg1 == nil {
		return nil, err
	}
	return arg1.(*models.Member), err
}
