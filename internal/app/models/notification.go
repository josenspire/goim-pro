package models

import (
	"goim-pro/internal/app/repos/base"
	tbl "goim-pro/pkg/db"
)

//const (
//	TypeSystem         = "SYSTEM"
//	TypeContactRequest = "CT_REQ"
//	TypeContactAccept  = "CT_ACT"
//	TypeContactRefuse  = "CT_REF"
//)

type Notification struct {
	Id               string              `json:"id" gorm:"column:id; type:varchar(32); primary_key; not null"`
	MessageId        string              `json:"messageId" gorm:"column:messageId; type:varchar(32); not null"`
	NotificationType string              `json:"type" gorm:"column:type; type:varchar(100); not null"`
	IsNeedRemind     bool                `json:"remind" gorm:"column:remind; type:varchar(32); not null"`
	FromUserId       string              `json:"fromUserId" gorm:"column:fromUserId; type:varchar(32); not null"`
	ToUserId         string              `json:"toUserId" gorm:"column:toUserId; type:varchar(32); not null"`
	Message          NotificationMessage `gorm:"ForeignKey:messageId;"` // foreign key
	base.BaseModel
}

type NotificationMessage struct {
	MessageId    string `json:"messageId" gorm:"column:messageId; type:varchar(32); primary_key; not null"`
	MsgType      string `json:"type" gorm:"column:type; type:varchar(100); not null"`
	IsNeedRemind bool   `json:"isNeedRemind" gorm:"column:isNeedRemind; type:varchar(32); not null"`
	base.BaseModel
}

func (Notification) TableName() string {
	return tbl.TableNotifications
}

func (NotificationMessage) TableName() string {
	return tbl.TableNotificationMsgs
}
