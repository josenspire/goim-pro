package user

import (
	"errors"
	"github.com/jinzhu/gorm"
	"goim-pro/config"
	"goim-pro/internal/app/repos/base"
	"goim-pro/pkg/logs"
	"goim-pro/pkg/utils"
)

type User struct {
	Password string `json:"password" gorm:"column:password; type:varchar(255); not null"`
	Role     string `json:"role" gorm:"column:role; type:ENUM('1', '10', '99'); default:'1'"`
	Status   string `json:"status" gorm:"column:status; type:ENUM('ACTIVE', 'INACTIVE'); default: 'ACTIVE'; not null"`
	UserProfile
	base.BaseModel
}

type UserProfile struct {
	UserID      string `json:"userID" gorm:"column:userID; type:varchar(32); primary_key; not null"`
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
	IsTelephoneOrEmailRegistered(telephone string, email string) (bool, error)
	Register(newUser *User) error
	LoginByEmail(email string, password string) (user *User, err error)
	LoginByTelephone(telephone string, password string) (user *User, err error)
	RemoveUserByUserID(userID string, isForce bool) error
}

var logger = logs.GetLogger("ERROR")
var crypto = utils.NewCrypto()
var mysqlDB *gorm.DB

func (User) TableName() string {
	return "users"
}

func NewUserRepo(db *gorm.DB) IUserRepo {
	mysqlDB = db
	return &User{}
}

// callbacks hock -- before create, encrypt password
func (u *User) BeforeCreate(scope *gorm.Scope) (err error) {
	if u.Password == "" {
		return errors.New("[aes] invalid password parameter")
	}
	encryptPassword, err := crypto.AESEncrypt(u.Password, config.GetApiSecretKey())
	if err != nil {
		logger.Errorf("[aes] encrypt password error: %s", err.Error())
		return
	}
	err = scope.SetColumn("password", encryptPassword)
	if err != nil {
		logger.Errorf("[aes] encrypt password error: %s", err.Error())
	}
	return
}

func (u *User) Register(newUser *User) (err error) {
	_db := mysqlDB.Create(&newUser)
	if _db.Error != nil {
		err = _db.Error
		logger.Errorf("create user error: %s", err.Error())
	}
	return
}

func (u *User) IsTelephoneOrEmailRegistered(telephone string, email string) (bool, error) {
	var isTelExist bool = true
	var err error
	if telephone != "" {
		err = mysqlDB.First(&User{}, "telephone = ?", telephone).Error
		if err == gorm.ErrRecordNotFound {
			isTelExist = false
			err = nil
		} else if err != nil {
			logger.Errorf("query user by telephone error: %s", err.Error())
			return isTelExist, err
		}
	}

	if isTelExist {
		return true, nil
	}

	var isEmailExist bool = true
	if email != "" {
		err = mysqlDB.First(&User{}, "email = ?", email).Error
		if err == gorm.ErrRecordNotFound {
			isEmailExist = false
			err = nil
		} else if err != nil {
			logger.Errorf("query user by email error: %s", err.Error())
			return isEmailExist, err
		}
	}

	return isTelExist || isEmailExist, nil
}

func (u *User) LoginByEmail(email string, password string) (user *User, err error) {
	db := mysqlDB.First(user, "email = ? and password = ?", email, password)
	if db.RecordNotFound() {
		err = utils.ErrAccountOrPswInvalid
	} else {
		err = db.Error
	}
	return
}

func (u *User) LoginByTelephone(telephone string, password string) (user *User, err error) {
	db := mysqlDB.First(user, "telephone = ? and password = ?", telephone, password)
	if db.RecordNotFound() {
		err = utils.ErrAccountOrPswInvalid
	} else {
		err = db.Error
	}
	return
}

func (u *User) RemoveUserByUserID(userID string, isForce bool) (err error) {
	_db := mysqlDB
	if isForce {
		_db = mysqlDB.Unscoped()
	}
	_db = _db.Delete(&User{}, "userID = ?", userID)
	if _db.RecordNotFound() {
		logger.Warningln("remove user fail, userID not found")
	} else if _db.Error != nil {
		logger.Errorf("error happened to remove user: %v\n", _db.Error)
		err = _db.Error
	}
	return
}
