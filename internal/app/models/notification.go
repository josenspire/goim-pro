package models

import (
	"goim-pro/internal/app/repos/base"
	tbl "goim-pro/internal/db"
)

const (
	NotifyTypeSystem = "SYSTEM"
	NotifyTypeUser   = "USER"

	MsgTypeContactRequest = "CT_REQ"
	MsgTypeContactAccept  = "CT_ACT"
	MsgTypeContactRefuse  = "CT_REF"
	MsgTypeContactDelete  = "CT_DEL"
)

const (
	StsPending  = "PENDING"
	StsSent     = "SENT"
	StsCancel   = "CANCEL"
	StsReceived = "RECEIVED"
)

/**
 * 消息类型
 */
const (
	ObjectNameTxtMsg     = "GM:TxtMsg"     // text
	ObjectNameImgMsg     = "GM:ImgMsg"     // IMG
	ObjectNameGIFMsg     = "GM:GIFMsg"     // GIF
	ObjectNameHQVCMsg    = "GM:HQVCMsg"    // 语音消息
	ObjectNameFileMsg    = "GM:FileMsg"    // 文件消息
	ObjectNameSightMsg   = "GM:SightMsg"   // 小视频消息
	ObjectNameLBSMsg     = "GM:LBSMsg"     // 位置消息
	ObjectNameInfoNtf    = "GM:InfoNtf"    // 小灰条提示消息，e.g：请在聊天中注意人身财产安全
	ObjectNameProfileNtf = "GM:ProfileNtf" // 资料变更通知消息
	ObjectNameContactNtf = "GM:ContactNtf" // 联系人(好友)通知消息
	ObjectNameCmdMsg     = "GM:CmdMsg"     // 命令消息

	ObjectNameCustomContact = "Custom:ContactNtf" // 单用户通知
)

type Notification struct {
	NtfId            string              `json:"ntfId" gorm:"column:ntfId; type:varchar(32); primary_key; not null"`
	MessageId        string              `json:"messageId" gorm:"column:messageId; type:varchar(32); primary_key; not null"`
	SenderId         string              `json:"senderId" gorm:"column:senderId; type:varchar(32); not null"`
	TargetId         string              `json:"targetId" gorm:"column:targetId; type:varchar(32); not null"`
	NotificationType string              `json:"ntfType" gorm:"column:ntfType; type:varchar(20); not null"`
	MsgType          string              `json:"msgType" gorm:"column:msgType; type:varchar(100); not null"`
	ObjectName       string              `json:"objectName" gorm:"column:objectName; type:varchar(50); not null"`
	IsNeedRemind     bool                `json:"remind" gorm:"column:remind; type:varchar(32); not null"`
	Status           string              `json:"status" gorm:"column:status; type:varchar(32); not null"`
	Message          NotificationMessage `gorm:"ForeignKey:MessageId"`
	Sender           User                `gorm:"ForeignKey:senderId"`
	Receiver         User                `gorm:"ForeignKey:targetId"`
	base.BaseModel
}

type NotificationMessage struct {
	MessageId    string `json:"messageId" gorm:"column:messageId; type:varchar(32); primary_key; not null"`
	MsgOperation int32  `json:"operation" gorm:"column:operation; type:int(2); not null"`
	MsgContent   string `json:"content" gorm:"column:content; type:varchar(255); not null"`
	MsgExtra     string `json:"extra" gorm:"column:extra; type:varchar(255); not null"`
	base.BaseModel
}

func (Notification) TableName() string {
	return tbl.TableNotifications
}

func (NotificationMessage) TableName() string {
	return tbl.TableNotificationMsgs
}
