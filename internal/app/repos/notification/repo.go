package ntf

import (
	"bytes"
	"fmt"
	"github.com/jinzhu/gorm"
	. "goim-pro/internal/app/models"
	tbl "goim-pro/pkg/db"
	"goim-pro/pkg/logs"
	"goim-pro/pkg/utils"
	"time"
)

type NotificationImpl Notification

type INotificationRepo interface {
	// notifications
	InsertOne(notification *Notification) (*Notification, error)
	FindAll(condition string) (notification []Notification, err error)
	FindOne(condition interface{}) (ntf *Notification, err error)
	UpdateOne(condition, updated interface{}) (err error)

	// messages
	InsertMessages(messages ...*NotificationMessage) (err error)
	UpdateOneMessage(condition, updated interface{}) (err error)
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

func (n *NotificationImpl) FindAll(condition string) (notifications []Notification, err error) {
	db := mysqlDB.Preload("Message").Preload("Sender").Find(&notifications, condition)
	//mysqlDB.Find(&notifications, "targetId = ? and createdAt >= ?", userId, fromDate).Preload("User").Related(&groupProfile.Members, "Members")
	//db := mysqlDB.Preload("Message").Where("targetId = ? and createdAt >= ?", userId, fromDate).Find(&notifications)
	if db.RecordNotFound() {
		return nil, nil
	}
	if err = db.Error; err != nil {
		return nil, err
	}

	return notifications, nil
}

func (n *NotificationImpl) FindOne(condition interface{}) (ntf *Notification, err error) {
	ntf = new(Notification)
	db := mysqlDB.Preload("Message").First(ntf, condition)
	if db.RecordNotFound() {
		return nil, nil
	}
	if err = db.Error; err != nil {
		return nil, err
	}
	return ntf, nil
}

func (n *NotificationImpl) UpdateOne(condition, updated interface{}) (err error) {
	db := mysqlDB.Table(tbl.TableNotifications).Where(condition).Update(updated)
	if err = db.Error; err != nil {
		logger.Errorf("error happened to update notifications: %v", err)
	}
	return
}

func (n *NotificationImpl) InsertMessages(messages ...*NotificationMessage) (err error) {
	var buffer bytes.Buffer
	sql := "INSERT INTO `notificationMsgs` (`messageId`, `content`, `extra`, `createdAt`, `updatedAt`) values"
	if _, err := buffer.WriteString(sql); err != nil {
		return err
	}
	for i, e := range messages {
		nowDateTime := utils.TimeFormat(time.Now(), utils.MysqlDateTimeFormat)
		if i == len(messages)-1 {
			buffer.WriteString(fmt.Sprintf("('%s', '%d', '%s', '%s', '%s', '%s');", e.MessageId, e.MsgOperation, e.MsgContent, e.MsgExtra, nowDateTime, nowDateTime))
		} else {
			buffer.WriteString(fmt.Sprintf("('%s', '%d', '%s', '%s', '%s', '%s'),", e.MessageId, e.MsgOperation, e.MsgContent, e.MsgExtra, nowDateTime, nowDateTime))
		}
	}
	if err := mysqlDB.Exec(buffer.String()).Error; err != nil {
		logger.Errorf("exec insert notification message error: %s", err.Error())
		return err
	}
	return nil
}

func (n *NotificationImpl) UpdateOneMessage(condition, updated interface{}) (err error) {
	db := mysqlDB.Table(tbl.TableNotificationMsgs).Where(condition).Update(updated)
	if err = db.Error; err != nil {
		logger.Errorf("error happened to update notification msg: %v", err)
	}
	return
}
