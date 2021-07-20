package notflogs

import (
	"bytes"
	"github.com/jinzhu/gorm"
	"goim-pro/internal/app/models"
	"goim-pro/pkg/logs"
)

type NotificationLogImpl models.Contact

type INotificationLogRepo interface {
	InsertMany(notification ...*models.NotificationLog) (err error)
	FindAll(condition interface{}) (notification []models.NotificationLog, err error)
}

var logger = logs.GetLogger("ERROR")
var mysqlDB *gorm.DB

func NewNotification(db *gorm.DB) INotificationLogRepo {
	mysqlDB = db
	return &NotificationLogImpl{}
}

func (n *NotificationLogImpl) InsertMany(notification ...*models.NotificationLog) (err error) {
	// BatchSave 批量插入数据
	var buffer bytes.Buffer
	sql := "INSERT INTO `notifications` (`id`, `messageId`, `type`, `remind`, `senderId`, `targetId`, `createdAt`, `updatedAt`) values"
	if _, err := buffer.WriteString(sql); err != nil {
		return err
	}

	//for i, e := range notification {
	//	nowDateTime := utils.TimeFormat(time.Now(), utils.MysqlDateTimeFormat)
	//
	//	if i == len(notification)-1 {
	//		buffer.WriteString(fmt.Sprintf("('%s', '%s', '%s', '%v', '%s', '%s', '%s', '%s');", e.Id, e.MessageId, e.NotificationType, e.IsNeedRemind, e.SenderId, e.TargetId, nowDateTime, nowDateTime))
	//	} else {
	//		buffer.WriteString(fmt.Sprintf("('%s', '%s', '%s', '%v', '%s', '%s', '%s', '%s'),", e.Id, e.MessageId, e.NotificationType, e.IsNeedRemind, e.SenderId, e.TargetId, nowDateTime, nowDateTime))
	//	}
	//}
	//if err = mysqlDB.Exec(buffer.String()).Error; err != nil {
	//	logger.Errorf("exec insert notification notification error: %v", err)
	//}
	return
}

func (n *NotificationLogImpl) FindAll(condition interface{}) (notifications []models.NotificationLog, err error) {
	db := mysqlDB.Preload("Message").Find(&notifications, condition)
	if db.RecordNotFound() {
		return nil, nil
	}
	if err = db.Error; err != nil {
		return nil, err
	}
	return notifications, nil
}
