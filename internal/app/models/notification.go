package models

import (
	"goim-pro/internal/app/repos/base"
	tbl "goim-pro/pkg/db"
)

const (
	NotifyTypeSystem = "SYSTEM"
	NotifyTypeUser   = "USER"

	MsgTypeContactRequest = "CT_REQ"
	MsgTypeContactAccept  = "CT_ACT"
	MsgTypeContactRefuse  = "CT_REF"
)

const (
	StsPending  = "PENDING"
	StsSent     = "SENT"
	StsCancel   = "CANCEL"
	StsReceived = "RECEIVED"
)

type Notification struct {
	NotifyId         string              `json:"notifyId" gorm:"column:notifyId; type:varchar(32); primary_key; not null"`
	MessageId        string              `json:"messageId" gorm:"column:messageId; type:varchar(32); primary_key; not null"`
	NotificationType string              `json:"type" gorm:"column:type; type:varchar(100); not null"`
	FromUserId       string              `json:"fromUserId" gorm:"column:fromUserId; type:varchar(32); not null"`
	ToUserId         string              `json:"toUserId" gorm:"column:toUserId; type:varchar(32); not null"`
	IsNeedRemind     bool                `json:"remind" gorm:"column:remind; type:varchar(32); not null"`
	Status           string              `json:"status" gorm:"column:status; type:varchar(32); not null"`
	Message          NotificationMessage `gorm:"ForeignKey:messageId;"` // foreign key
	base.BaseModel
}

type NotificationMessage struct {
	MessageId    string `json:"messageId" gorm:"column:messageId; type:varchar(32); primary_key; not null"`
	MsgType      string `json:"type" gorm:"column:type; type:varchar(100); not null"`
	MsgOperation int32  `json:"operation" gorm:"column:operation; type:int(2); not null"`
	MsgContent   string `json:"content" gorm:"column:content; type:varchar(255); not null"`
	base.BaseModel
}

func (Notification) TableName() string {
	return tbl.TableNotifications
}

func (NotificationMessage) TableName() string {
	return tbl.TableNotificationMsgs
}
