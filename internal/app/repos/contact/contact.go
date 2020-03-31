package contact

import (
	"bytes"
	"fmt"
	"github.com/jinzhu/gorm"
	"goim-pro/internal/app/repos/base"
	"goim-pro/pkg/db"
	"goim-pro/pkg/logs"
)

type Contact struct {
	// Composite primary key - userId+contactId
	UserId    string `json:"userId" gorm:"column:userId; type:varchar(32); primary_key; not null"`
	ContactId string `json:"contactId" gorm:"column:contactId; type:varchar(32); primary_key; not null"`
	RemarkProfile
	base.BaseModel
}

type RemarkProfile struct {
	RemarkName  string `json:"remarkName" gorm:"column:remarkName; type:varchar(16)"`
	Telephone   string `json:"telephone" gorm:"column:telephone; type:varchar(255)"` // can support multiple tel, split by `;`
	Description string `json:"description" gorm:"column:description; type:varchar(255)"`
	Tags        string `json:"tags" gorm:"column:tags; type:varchar(255)"` // can support multiple tag, split by `;`
}

type IContactRepo interface {
	IsExistContact(userId, contactId string) (isExist bool, err error)
	FindOne(condition *Contact) (contact *Contact, err error)
	InsertContacts(newContacts ...*Contact) (err error)
	RemoveContactsByIds(userId string, contactIds ...string) (err error)
	FindOneAndUpdateRemark(ct *Contact, remarkInfo map[string]interface{}) (err error)
}

var logger = logs.GetLogger("ERROR")
var mysqlDB *gorm.DB

func (Contact) TableName() string {
	return tbl.TableContacts
}

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
	db := mysqlDB.Table(cta.TableName()).Where(ct).Update(remarkProfile)
	if err = db.Error; err != nil {
		logger.Errorf("error happened to update user profile: %v", err)
	}
	return
}
