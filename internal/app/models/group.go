package models

import (
	"goim-pro/internal/app/repos/base"
	tbl "goim-pro/pkg/db"
)

// user conversation group
type Group struct {
	GroupId     string   `json:"groupId" gorm:"column:groupId; type:varchar(32); primary_key; not null"`
	CreatedBy   string   `json:"createdBy" gorm:"column:createdBy; type:varchar(32); not null"`
	OwnerUserId string   `json:"ownerUserId" gorm:"column:ownerUserId; type:varchar(32); not null"`
	Name        string   `json:"name" gorm:"column:name; type:varchar(100); not null; default: 'NewGroup'"`
	Avatar      string   `json:"avatar" gorm:"column:avatar; type:varchar(255); default: ''"`
	Notice      string   `json:"notice" gorm:"column:notice; type:varchar(255); default: ''"`
	Members     []Member `gorm:"ForeignKey:GroupId;AssociationForeignKey:GroupId"` // foreign key
	base.BaseModel
}

// conversation group members
type Member struct {
	UserId  string `json:"userId" gorm:"column:userId; type:varchar(32); primary_key; not null"`
	Alias   string `json:"alias" gorm:"column:alias; type:varchar(16)"`
	Role    string `json:"role" gorm:"column:role; type:ENUM('1', '10', '50', '99'); default: '1'; not null"`
	Status  string `json:"status" gorm:"column:status; type:ENUM('NORMAL', 'MUTE'); default: 'NORMAL'; not null"`
	GroupId string `json:"groupId" gorm:"column:groupId; type:varchar(32); not null"`
	base.BaseModel
}

func (Group) TableName() string {
	return tbl.TableGroups
}

func (Member) TableName() string {
	return tbl.TableMembers
}

func NewGroup(userId, groupName string, members []Member) *Group {
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
		Members:     members,
	}
}

func NewMember(memberId, alias string) Member {
	return Member{
		UserId: memberId,
		Alias:  alias,
		Role:   "1",
		Status: "NORMAL",
	}
}
