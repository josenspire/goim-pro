package user

import (
	"github.com/jinzhu/gorm"
	"goim-pro/config"
	"goim-pro/internal/app/constants"
	"goim-pro/internal/app/repos/base"
	"goim-pro/pkg/logs"
	"goim-pro/pkg/utils"
)

type User struct {
	UserID   uint64 `json:"userID" gorm:"column:userID; primary_key; not null"`
	Password string `json:"password" gorm:"column:password; type:varchar(255); not null"`
	Role     string `json:"role" gorm:"column:role; type:ENUM('1', '10', '99'); default:'1'"`
	Status   string `json:"status" gorm:"column:status; type:ENUM('ACTIVE', 'INACTIVE'); default: 'ACTIVE'; not null"`
	UserProfile
	base.BaseModel
}

type UserProfile struct {
	Telephone   string `json:"telephone" gorm:"column:telephone; type:varchar(11)"`
	Email       string `json:"email" gorm:"column:email; type:varchar(100)"`
	Username    string `json:"username" gorm:"column:username; type:varchar(18)"`
	Nickname    string `json:"nickname" gorm:"column:nickname; type:varchar(16)"`
	Avatar      string `json:"avatar" gorm:"column:avatar; type:varchar(255)"`
	Description string `json:"description" gorm:"column:description; type:varchar(255)"`
	Sex         string `json:"sex" gorm:"column:sex; type: ENUM('MALE', 'FEMALE'); default:'FEMALE'"`
	Birthday    int64  `json:"birthday" gorm:"column:birthday; type: bigint"`
	Location    string `json:"location" gorm:"column:location; type: varchar(255)"`
}

type IUserRepo interface {
	IsTelephoneRegistered(telephone string) (bool, error)
	Register(newUser *User) error
	RemoveUserByUserId(userID uint64) error
}

var logger = logs.GetLogger("ERROR")
var crypto = utils.NewCrypto()
var mysqlDB *gorm.DB

func (User) TableName() string {
	return "users"
}

func NewUserModel() *User {
	return &User{}
}

func NewUserRepo(db *gorm.DB) IUserRepo {
	mysqlDB = db
	return &User{}
}

// callbacks hock -- before create, encrypt password
func (u *User) BeforeCreate(scope *gorm.Scope) error {
	if utils.IsEmptyString(u.Password) {
		return nil
	}
	encryptPassword, err := crypto.AESEncrypt(u.Password, config.GetApiSecretKey())
	if err != nil {
		logger.Errorf("[aes] encrypt password error: %v", err)
		return err
	}
	err = scope.SetColumn("password", encryptPassword)
	if err != nil {
		logger.Errorf("[aes] encrypt password error: %v", err)
		return err
	}
	return nil
}

func (u *User) Register(newUser *User) (err error) {
	_db := mysqlDB.Create(&newUser)
	if _db.Error != nil {
		err = _db.Error
		logger.Errorf("create user error: %v", err)
	}
	return
}

func (u *User) IsTelephoneRegistered(telephone string) (bool, error) {
	err := mysqlDB.Where("telephone = ? and status = 'ACTIVE'", telephone).Find(&u).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	if err != nil {
		logger.Errorf("query user by telephone error: %v", err)
		return false, err
	}
	return true, nil
}

func (u *User) RemoveUserByUserId(userID uint64) (err error) {
	_user := &User{
		UserID: userID,
		Status: constants.StatusActive,
	}
	_db := mysqlDB.Set("gorm:delete_option", "UPDATE users SET users.status = 'INACTIVE'").Delete(_user)
	if _db.RecordNotFound() {
		logger.Warningln("remove user fail, userID not found")
	} else if _db.Error != nil {
		logger.Errorf("error happened to remove user: %v\n", _db.Error)
		err = _db.Error
	}
	return _db.Error
}
