package models

import (
	"goim-pro/internal/app/repos/base"
	tbl "goim-pro/pkg/db"
)

const (
	TypeSystem         = "SYSTEM"
	TypeContactRequest = "CT_REQ"
	TypeContactAccept  = "CT_ACT"
	TypeContactRefuse  = "CT_REF"
)

type NotificationLog struct {
	Id         string `json:"groupId" gorm:"column:groupId; type:varchar(32); primary_key; not null"`
	MsgType    string `json:"type" gorm:"column:type; type:varchar(100); not null"`
	FromUserId string `json:"fromUserId" gorm:"column:fromUserId; type:varchar(32); not null"`
	ToUserId   string `json:"toUserId" gorm:"column:toUserId; type:varchar(32); not null"`
	Message    string `json:"message" gorm:"column:message; type:varchar(255)"`
	base.BaseModel
}

func (NotificationLog) TableName() string {
	return tbl.TableNotificationLogs
}

func NewNotificationLog(id, msgType, fromUserId, toUserId, message string) *NotificationLog {
	return &NotificationLog{
		Id:         id,
		MsgType:    msgType,
		FromUserId: fromUserId,
		ToUserId:   toUserId,
		Message:    message,
	}
}

