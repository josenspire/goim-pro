package contact

import (
	"bytes"
	"fmt"
	"github.com/jinzhu/gorm"
	"goim-pro/internal/app/models"
	. "goim-pro/internal/app/repos/user"
	"goim-pro/pkg/db"
	"goim-pro/pkg/logs"
)

type Contact models.Contact

type IContactRepo interface {
	IsExistContact(userId, contactId string) (isExist bool, err error)
	FindOne(condition *Contact) (contact *Contact, err error)
	FindAll(condition *Contact) (contacts []Contact, err error)
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
func (cta *Contact) FindAll(condition *Contact) (contacts []Contact, err error) {
	//err = mysqlDB.Where(condition).Joins("left join users on users.status = 'ACTIVE' and contacts.userId = users.userId ORDER BY contacts.createdAt DESC;").Scan(&contacts).Error
	//logger.Error(contacts)
	//return contacts, err
	var user User
	err = mysqlDB.Model(&user).Where(condition).Related(&contacts).Error
	logger.Info()
	return
}

func (cta *Contact) InsertContacts(newContacts ...*Contact) (err error) {
	// BatchSave 批量插入数据
	var buffer bytes.Buffer
	sql := "INSERT INTO `contacts` (`userId`, `contactId`, `remarkName`, `telephone`, `description`, `tags`) values"
	if _, err := buffer.WriteString(sql); err != nil {
		return err
	}
	for i, e := range newContacts {
		if i == len(newContacts)-1 {
			// TODO: should add createdAt dateTime
			buffer.WriteString(fmt.Sprintf("('%s', '%s', '%s', '%s', '%s', '%s');", e.UserId, e.ContactId, e.RemarkName, e.Telephone, e.Description, e.Tags))
		} else {
			buffer.WriteString(fmt.Sprintf("('%s', '%s', '%s', '%s', '%s', '%s'),", e.UserId, e.ContactId, e.RemarkName, e.Telephone, e.Description, e.Tags))
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
