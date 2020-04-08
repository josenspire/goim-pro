package contact

import (
	"bytes"
	"fmt"
	"github.com/jinzhu/gorm"
	"goim-pro/internal/app/models"
	"goim-pro/pkg/db"
	"goim-pro/pkg/logs"
	"goim-pro/pkg/utils"
	"time"
)

type Contact models.Contact

type IContactRepo interface {
	IsExistContact(userId, contactId string) (isExist bool, err error)
	FindOne(condition *Contact) (contact *Contact, err error)
	FindAll(condition map[string]interface{}) (contacts []Contact, err error)
	InsertContacts(newContacts ...*Contact) (err error)
	RemoveContactsByIds(userId string, contactIds ...string) (err error)
	FindOneAndUpdateRemark(ct *Contact, remarkInfo map[string]interface{}) (err error)
}

var logger = logs.GetLogger("ERROR")
var mysqlDB *gorm.DB

func NewContactRepo(db *gorm.DB) IContactRepo {
	mysqlDB = db
	return &Contact{}
}

func (cta *Contact) IsExistContact(userId, contactId string) (isExists bool, err error) {
	db := mysqlDB.First(&Contact{}, "userId = ? and contactId = ?", userId, contactId)
	if db.RecordNotFound() {
		return false, nil
	} else if err = db.Error; err != nil {
		logger.Errorf("checking contact error: %v", err)
	}
	return true, err
}

func (cta *Contact) FindOne(condition *Contact) (contact *Contact, err error) {
	contact = &Contact{}
	db := mysqlDB.Where(condition).First(&contact)
	if db.RecordNotFound() {
		return nil, nil
	}
	if err = db.Error; err != nil {
		logger.Errorf("error happened to query user information: %v", err)
	}
	return
}

// TODO: should verify
func (cta *Contact) FindAll(condition map[string]interface{}) (contacts []Contact, err error) {
	// SELECT * FROM `contacts`  WHERE `contacts`.`deletedAt` IS NULL AND ((`contacts`.`UserId` = '01E07SG858N3CGV5M1APVQKZYR'))
	// SELECT * FROM `users`  WHERE `users`.`deletedAt` IS NULL AND ((`userId` IN ('01E2JVWZTG60NG2SXFYNEPNMCB','01E2JXMC98SZXMGEGVTDECSD78')))
	err = mysqlDB.Preload("User").Find(&contacts, condition).Error
	return
}

func (cta *Contact) InsertContacts(newContacts ...*Contact) (err error) {
	// BatchSave 批量插入数据
	var buffer bytes.Buffer
	sql := "INSERT INTO `contacts` (`userId`, `contactId`, `remarkName`, `telephone`, `description`, `tags`, `createdAt`, `updatedAt`) values"
	if _, err := buffer.WriteString(sql); err != nil {
		return err
	}

	for i, e := range newContacts {
		nowDateTime := utils.TimeFormat(time.Now(), utils.MysqlDateTimeFormat)

		if i == len(newContacts)-1 {
			buffer.WriteString(fmt.Sprintf("('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s');", e.UserId, e.ContactId, e.RemarkName, e.Telephone, e.Description, e.Tags, nowDateTime, nowDateTime))
		} else {
			buffer.WriteString(fmt.Sprintf("('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s'),", e.UserId, e.ContactId, e.RemarkName, e.Telephone, e.Description, e.Tags, nowDateTime, nowDateTime))
		}
	}
	if err = mysqlDB.Exec(buffer.String()).Error; err != nil {
		logger.Errorf("exec insert contacts error: %v", err)
	}
	return
}

// remove user's contacts by contact ids, by force
func (cta *Contact) RemoveContactsByIds(userId string, contactIds ...string) (err error) {
	_db := mysqlDB.Unscoped().Delete(&Contact{}, "userId = ? and contactId IN (?)", userId, contactIds)
	if _db.Error != nil {
		logger.Errorf("error happened to remove user: %v", _db.Error)
		err = _db.Error
	}
	return
}

func (cta *Contact) FindOneAndUpdateRemark(ct *Contact, remarkProfile map[string]interface{}) (err error) {
	db := mysqlDB.Table(tbl.TableContacts).Where(ct).Update(remarkProfile)
	if err = db.Error; err != nil {
		logger.Errorf("error happened to update user profile: %v", err)
	}
	return
}
