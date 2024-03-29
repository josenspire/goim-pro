package group

import (
	"bytes"
	"fmt"
	"github.com/jinzhu/gorm"
	. "goim-pro/internal/app/models"
	tbl "goim-pro/internal/db"
	"goim-pro/pkg/logs"
	"goim-pro/pkg/utils"
	"time"
)

type GroupImpl Group

type MemberImpl Member

type IGroupRepo interface {
	// group
	CreateGroup(groupProfile *Group) (newGroup *Group, err error)
	RemoveGroupByGroupId(groupId string, isForce bool) (err error)
	FindOneGroup(condition interface{}) (groupProfile *Group, err error)
	CountGroup(condition interface{}) (count int, err error)
	FindOneGroupAndUpdate(condition interface{}, updated interface{}) (newProfile *Group, err error)

	// member
	FindOneMember(condition interface{}) (memberProfile *Member, err error)
	InsertMembers(members ...*Member) (err error)
	RemoveMembers(groupId string, memberIds []string, isForce bool) (deleteCount int64, err error)
	FindOneMemberAndUpdate(condition interface{}, updated interface{}) (newProfile *Member, err error)
}

var logger = logs.GetLogger("ERROR")
var mysqlDB *gorm.DB

func NewGroupRepo(db *gorm.DB) IGroupRepo {
	mysqlDB = db
	return &GroupImpl{}
}

// CreateGroup - create group with group profile and members
func (i *GroupImpl) CreateGroup(groupProfile *Group) (newGroup *Group, err error) {
	// Don't update associations having primary key, but will save reference
	_db := mysqlDB.Create(groupProfile)
	if _db.Error != nil {
		err = _db.Error
		logger.Errorf("create group error: %s", err.Error())

		return nil, err
	}
	return _db.Value.(*Group), nil
}

// group(1)->(n)members(1)->(1)user
func (i *GroupImpl) FindOneGroup(condition interface{}) (*Group, error) {
	var groupProfile = Group{}
	if err := mysqlDB.First(&groupProfile, condition).Preload("User").Related(&groupProfile.Members, "Members").Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &groupProfile, nil
}

func (i *GroupImpl) CountGroup(condition interface{}) (count int, err error) {
	//db := mysqlDB.Model(&Group{}).Where(condition).Count(&count)
	db := mysqlDB.Table(tbl.TableGroups).Where(condition).Count(&count)
	if db.RecordNotFound() {
		return 0, nil
	}
	if err = db.Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (i *GroupImpl) FindOneGroupAndUpdate(condition interface{}, updated interface{}) (newProfile *Group, err error) {
	newProfile = &Group{}
	db := mysqlDB.Model(&Group{}).Where(condition).Update(updated).First(newProfile)
	if db.RecordNotFound() {
		return nil, nil
	}
	if err = db.Error; err != nil {
		logger.Errorf("error happened to update group profile: %s", err.Error())
		return nil, err
	}
	return newProfile, nil
}

func (i *GroupImpl) RemoveGroupByGroupId(groupId string, isForce bool) (err error) {
	_db := mysqlDB
	if isForce {
		_db = mysqlDB.Unscoped()
	}
	_db = _db.Delete(&Group{}, "groupId = ?", groupId)
	if _db.RecordNotFound() {
		logger.Warningln("remove group fail, groupId not found")
	} else if _db.Error != nil {
		logger.Errorf("error happened to remove group: %s", _db.Error)
		err = _db.Error
	}
	return
}

func (i *GroupImpl) InsertMembers(members ...*Member) (err error) {
	var buffer bytes.Buffer
	sql := "INSERT INTO `members` (`groupId`, `userId`, `alias`, `role`, `status`, `createdAt`, `updatedAt`) values"
	if _, err := buffer.WriteString(sql); err != nil {
		return err
	}
	for i, e := range members {
		nowDateTime := utils.TimeFormat(time.Now(), utils.MysqlDateTimeFormat)
		if i == len(members)-1 {
			buffer.WriteString(fmt.Sprintf("('%s', '%s', '%s', '%s', '%s', '%s', '%s');", e.GroupId, e.UserId, e.Alias, "1", "NORMAL", nowDateTime, nowDateTime))
		} else {
			buffer.WriteString(fmt.Sprintf("('%s', '%s', '%s', '%s', '%s', '%s', '%s'),", e.GroupId, e.UserId, e.Alias, "1", "NORMAL", nowDateTime, nowDateTime))
		}
	}
	if err := mysqlDB.Exec(buffer.String()).Error; err != nil {
		logger.Errorf("exec insert members error: %s", err.Error())
		return err
	}
	return nil
}

func (i *GroupImpl) FindOneMember(condition interface{}) (memberProfile *Member, err error) {
	memberProfile = &Member{}
	db := mysqlDB.Where(condition).First(&memberProfile)
	if db.RecordNotFound() {
		return nil, nil
	}
	if err = db.Error; err != nil {
		logger.Errorf("error happened to query group information: %s", err.Error())
		return nil, err
	}
	return memberProfile, nil
}

func (i *GroupImpl) RemoveMembers(groupId string, memberIds []string, isForce bool) (deleteCount int64, err error) {
	_db := mysqlDB
	if isForce {
		_db = mysqlDB.Unscoped()
	}
	_db = _db.Delete(&Member{}, "groupId = ? and userId IN (?)", groupId, memberIds)
	if _db.RecordNotFound() {
		logger.Warningln("remove group member fail, group or member not found")
	} else if _db.Error != nil {
		logger.Errorf("error happened to remove group member: %s", _db.Error)
		err = _db.Error
	}
	return _db.RowsAffected, nil
}

func (i *GroupImpl) FindOneMemberAndUpdate(condition interface{}, updated interface{}) (newProfile *Member, err error) {
	newProfile = &Member{}
	db := mysqlDB.Model(&Member{}).Where(condition).Update(updated).First(newProfile)
	if db.RecordNotFound() {
		return nil, nil
	}
	if err = db.Error; err != nil {
		logger.Errorf("error happened to update group profile: %s", err.Error())
		return nil, err
	}
	return newProfile, nil
}
