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

type ContactImpl models.Contact

type IContactRepo interface {
	IsContactExists(userId, contactId string) (isExist bool, err error)
	FindOne(condition map[string]interface{}) (contact *models.Contact, err error)
	FindAll(condition map[string]interface{}) (contacts []models.Contact, err error)
	InsertContacts(newContacts ...*models.Contact) (err error)
	RemoveContactsByIds(userId string, contactIds ...string) (err error)
	FindOneAndUpdateRemark(ct map[string]interface{}, remarkInfo map[string]interface{}) (err error)
}

var logger = logs.GetLogger("ERROR")
var mysqlDB *gorm.DB

func NewContactRepo(db *gorm.DB) IContactRepo {
	mysqlDB = db
	return &ContactImpl{}
}

func (cta *ContactImpl) IsContactExists(userId, contactId string) (isExists bool, err error) {
	db := mysqlDB.First(&models.Contact{}, "userId = ? and contactId = ?", userId, contactId)
	if db.RecordNotFound() {
		return false, nil
	}
	if err = db.Error; err != nil {
		logger.Errorf("checking contact error: %v", err)
		return false, err
	}
	return true, nil
}

func (cta *ContactImpl) FindOne(condition map[string]interface{}) (contact *models.Contact, err error) {
	contact = &models.Contact{}
	db := mysqlDB.Where(condition).First(contact)
	if db.RecordNotFound() {
		return nil, nil
	}
	if err = db.Error; err != nil {
		logger.Errorf("error happened to query user information: %v", err)
	}
	return
}

func (cta *ContactImpl) FindAll(condition map[string]interface{}) (contacts []models.Contact, err error) {
	// SELECT * FROM `contacts`  WHERE `contacts`.`deletedAt` IS NULL AND ((`contacts`.`UserId` = '01E07SG858N3CGV5M1APVQKZYR'))
	// SELECT * FROM `users`  WHERE `users`.`deletedAt` IS NULL AND ((`userId` IN ('01E2JVWZTG60NG2SXFYNEPNMCB','01E2JXMC98SZXMGEGVTDECSD78')))
	db := mysqlDB.Preload("User").Find(&contacts, condition)
	if db.RecordNotFound() {
		return nil, nil
	}
	if err = db.Error; err != nil {
		return nil, err
	}
	return contacts, nil
}

func (cta *ContactImpl) InsertContacts(newContacts ...*models.Contact) (err error) {
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
func (cta *ContactImpl) RemoveContactsByIds(userId string, contactIds ...string) (err error) {
	_db := mysqlDB.Unscoped().Delete(&models.Contact{}, "userId = ? and contactId IN (?)", userId, contactIds)
	if _db.Error != nil {
		logger.Errorf("error happened to remove user: %v", _db.Error)
		err = _db.Error
	}
	return
}

func (cta *ContactImpl) FindOneAndUpdateRemark(condition map[string]interface{}, remarkProfile map[string]interface{}) (err error) {
	db := mysqlDB.Table(tbl.TableContacts).Where(condition).Update(remarkProfile)
	if err = db.Error; err != nil {
		logger.Errorf("error happened to update user profile: %v", err)
	}
	return
}
