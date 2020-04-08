package group

import (
	"bytes"
	"fmt"
	"github.com/jinzhu/gorm"
	"goim-pro/internal/app/repos/base"
	tbl "goim-pro/pkg/db"
	"goim-pro/pkg/logs"
	"goim-pro/pkg/utils"
	"time"
)

// user conversation group
type Group struct {
	GroupId     string    `json:"groupId" gorm:"column:groupId; type:varchar(32); primary_key; not null"`
	CreatedBy   string    `json:"createdBy" gorm:"column:createdBy; type:varchar(32); not null"`
	OwnerUserId string    `json:"ownerUserId" gorm:"column:ownerUserId; type:varchar(32); not null"`
	Name        string    `json:"name" gorm:"column:name; type:varchar(100); not null; default: 'NewGroup'"` // TODO: should check out default group name
	Avatar      string    `json:"avatar" gorm:"column:avatar; type:varchar(255); default: ''"`
	Notice      string    `json:"notice" gorm:"column:notice; type:varchar(255); default: ''"`
	Members     []*Member `gorm:"ForeignKey:UserId;"` // foreign key
	base.BaseModel
}

// conversation group members
type Member struct {
	UserId string `json:"userId" gorm:"column:userId; type:varchar(32); primary_key; not null"`
	Alias  string `json:"alias" gorm:"column:alias; type:varchar(16)"`
	Role   string `json:"role" gorm:"column:role; type:ENUM('1', '10', '50', '99'); default: '1'; not null"`
	Status string `json:"status" gorm:"column:status; type:ENUM('NORMAL', 'MUTE'); default: 'NORMAL'; not null"`
	base.BaseModel
}

type IGroupRepo interface {
	InsertOne(groupProfile *Group) (newGroup *Group, err error)
	InsertMembers(members ...*Member) (newMembers *[]Member, err error)
}

var logger = logs.GetLogger("ERROR")
var mysqlDB *gorm.DB

func (Group) TableName() string {
	return tbl.TableGroups
}

func (Member) TableName() string {
	return tbl.TableMembers
}

func NewGroupRepo(db *gorm.DB) IGroupRepo {
	mysqlDB = db
	return &Group{}
}

func NewGroup(userId, groupName string) *Group {
	if groupName == "" {
		groupName = "新建群组"
	}
	return &Group{
		GroupId:     userId,
		CreatedBy:   userId,
		OwnerUserId: userId,
		Name:        groupName,
		Avatar:      "",
		Notice:      "",
	}
}

func NewMember(memberId, alias string) *Member {
	return &Member{
		UserId: memberId,
		Alias:  alias,
		Role:   "1",
		Status: "NORMAL",
	}
}

func (gp *Group) InsertOne(groupProfile *Group) (newGroup *Group, err error) {
	_db := mysqlDB.Create(&groupProfile)
	if _db.Error != nil {
		err = _db.Error
		logger.Errorf("create group error: %s", err.Error())

		return nil, err
	}
	return groupProfile, nil
}

func (gp *Group) InsertMembers(members ...*Member) (newMembers *[]Member, err error) {
	var buffer bytes.Buffer
	sql := "INSERT INTO `members` (`userId`, `alias`, `role`, `status`, `createdAt`, `updatedAt`) values"
	if _, err := buffer.WriteString(sql); err != nil {
		return nil, err
	}
	for i, e := range members {
		nowDateTime := utils.TimeFormat(time.Now(), utils.MysqlDateTimeFormat)
		if i == len(members)-1 {
			buffer.WriteString(fmt.Sprintf("('%s', '%s', '%s', '%s', '%s', '%s');", e.UserId, e.Alias, "1", "NORMAL", nowDateTime, nowDateTime))
		} else {
			buffer.WriteString(fmt.Sprintf("('%s', '%s', '%s', '%s', '%s', '%s'),", e.UserId, e.Alias, "1", "NORMAL", nowDateTime, nowDateTime))
		}
	}
	if err = mysqlDB.Exec(buffer.String()).Error; err != nil {
		logger.Errorf("exec insert members error: %v", err)
	}
	return
}
