package ntft

import (
	"bytes"
	"fmt"
	"github.com/jinzhu/gorm"
	. "goim-pro/internal/app/models"
	"goim-pro/pkg/logs"
	"goim-pro/pkg/utils"
	"time"
)

type NotificationImpl Notification

type INotificationRepo interface {
	InsertOne(notification *Notification) (*Notification, error)
	InsertMany(notification ...*Notification) (err error)
	FindAll(condition interface{}) (notification []Notification, err error)

	// message
	InsertMessages(messages ...*NotificationMessage) (err error)
}

var logger = logs.GetLogger("ERROR")
var mysqlDB *gorm.DB

func NewNotificationRepo(db *gorm.DB) INotificationRepo {
	mysqlDB = db
	return &NotificationImpl{}
}

func (n *NotificationImpl) InsertOne(notification *Notification) (*Notification, error) {
	_db := mysqlDB.Create(notification)
	if _db.Error != nil {
		err := _db.Error
		logger.Errorf("create notification error: %s", err.Error())

		return nil, err
	}
	return _db.Value.(*Notification), nil
}

func (n *NotificationImpl) InsertMany(notification ...*Notification) (err error) {
	// BatchSave 批量插入数据
	var buffer bytes.Buffer
	sql := "INSERT INTO `notifications` (`notifyId`, `type`, `fromUserId`, `toUserId`, `remind`, `status`, `createdAt`, `updatedAt`) values"
	if _, err := buffer.WriteString(sql); err != nil {
		return err
	}

	for i, e := range notification {
		nowDateTime := utils.TimeFormat(time.Now(), utils.MysqlDateTimeFormat)

		if i == len(notification)-1 {
			buffer.WriteString(fmt.Sprintf("('%s', '%s', '%s', '%v', '%s', '%s', '%s', '%s');", e.NotifyId, e.MessageId, e.NotificationType, e.IsNeedRemind, e.FromUserId, e.ToUserId, nowDateTime, nowDateTime))
		} else {
			buffer.WriteString(fmt.Sprintf("('%s', '%s', '%s', '%v', '%s', '%s', '%s', '%s'),", e.NotifyId, e.MessageId, e.NotificationType, e.IsNeedRemind, e.FromUserId, e.ToUserId, nowDateTime, nowDateTime))
		}
	}
	if err = mysqlDB.Exec(buffer.String()).Error; err != nil {
		logger.Errorf("exec insert notification notification error: %v", err)
	}
	return
}

func (n *NotificationImpl) FindAll(condition interface{}) (notifications []Notification, err error) {
	db := mysqlDB.Preload("Message").Find(&notifications, condition)
	if db.RecordNotFound() {
		return nil, nil
	}
	if err = db.Error; err != nil {
		return nil, err
	}
	return notifications, nil
}

func (n *NotificationImpl) InsertMessages(messages ...*NotificationMessage) (err error) {
	var buffer bytes.Buffer
	sql := "INSERT INTO `notificationMsgs` (`messageId`, `msgType`, `content`, `createdAt`, `updatedAt`) values"
	if _, err := buffer.WriteString(sql); err != nil {
		return err
	}
	for i, e := range messages {
		nowDateTime := utils.TimeFormat(time.Now(), utils.MysqlDateTimeFormat)
		if i == len(messages)-1 {
			buffer.WriteString(fmt.Sprintf("('%s', '%s', '%v', '%s', '%s');", e.MessageId, e.MsgType, e.MsgContent, nowDateTime, nowDateTime))
		} else {
			buffer.WriteString(fmt.Sprintf("('%s', '%s', '%v', '%s', '%s'),", e.MessageId, e.MsgType, e.MsgContent, nowDateTime, nowDateTime))
		}
	}
	if err := mysqlDB.Exec(buffer.String()).Error; err != nil {
		logger.Errorf("exec insert notification message error: %s", err.Error())
		return err
	}
	return nil
}
