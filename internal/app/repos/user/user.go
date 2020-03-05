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
	UserId      string `json:"userId" gorm:"column:userId; type:varchar(32); primary_key; not null"`
	Telephone   string `json:"telephone" gorm:"column:telephone; type:varchar(11)"`
	Email       string `json:"email" gorm:"column:email; type:varchar(100)"`
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
	QueryByEmailAndPassword(email string, password string) (user *User, err error)
	QueryByTelephoneAndPassword(telephone string, password string) (user *User, err error)
	RemoveUserByUserId(userId string, isForce bool) error
	ResetPasswordByTelephone(telephone string, newPassword string) error
	ResetPasswordByEmail(email string, newPassword string) error
	FindByUserId(userId string) (*User, error)
	FindOneUser(user *User) (*User, error)
	FindOneAndUpdateProfile(user *User, profile map[string]interface{}) error
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
	} else {
		isTelExist = false
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
	} else {
		isEmailExist = false
	}

	return isTelExist || isEmailExist, nil
}

func (u *User) QueryByEmailAndPassword(email string, enPassword string) (*User, error) {
	var user = &User{}
	var err error
	db := mysqlDB.First(&user, "email = ? and password = ?", email, enPassword)
	if db.RecordNotFound() {
		err = utils.ErrAccountOrPwdInvalid
	} else {
		err = db.Error
	}
	return user, err
}

func (u *User) QueryByTelephoneAndPassword(telephone string, enPassword string) (*User, error) {
	var user = &User{}
	var err error
	db := mysqlDB.First(user, "telephone = ? and password = ?", telephone, enPassword)
	if db.RecordNotFound() {
		err = utils.ErrAccountOrPwdInvalid
	} else {
		err = db.Error
	}
	return user, err
}

func (u *User) RemoveUserByUserId(userId string, isForce bool) (err error) {
	_db := mysqlDB
	if isForce {
		_db = mysqlDB.Unscoped()
	}
	_db = _db.Delete(&User{}, "userId = ?", userId)
	if _db.RecordNotFound() {
		logger.Warningln("remove user fail, userId not found")
	} else if _db.Error != nil {
		logger.Errorf("error happened to remove user: %v", _db.Error)
		err = _db.Error
	}
	return
}

func (u *User) ResetPasswordByTelephone(telephone string, newPassword string) (err error) {
	db := mysqlDB.Model(&User{}).Where("telephone = ?", telephone).Update("password", newPassword)
	if err = db.Error; err != nil {
		logger.Errorf("error happened to reset password by telephone: %v", err)
	}
	return
}

func (u *User) ResetPasswordByEmail(email string, newPassword string) (err error) {
	db := mysqlDB.Model(&User{}).Where("email = ?", email).Update("password", newPassword)
	if err = db.Error; err != nil {
		logger.Errorf("error happened to reset password by email: %v", err)
	}
	return
}

func (u *User) FindByUserId(userId string) (user *User, err error) {
	user = &User{}
	db := mysqlDB.First(user, "userId = ?", userId)
	if db.RecordNotFound() {
		err = utils.ErrInvalidUserId
	} else if err = db.Error; err != nil {
		logger.Errorf("error happend to get user by userId: %v", err)
	}
	return
}

func (u *User) FindOneUser(us *User) (user *User, err error) {
	user = &User{}
	db := mysqlDB.Where(us).First(&user)
	if db.RecordNotFound() {
		err = utils.ErrUserNotExists
	} else if err = db.Error; err != nil {
		logger.Errorf("error happend to query user information: %v", err)
	}
	return
}

func (u *User) FindOneAndUpdateProfile(us *User, profile map[string]interface{}) (err error) {
	db := mysqlDB.Table("users").Where(us).Update(profile)
	if err = db.Error; err != nil {
		logger.Errorf("error happened to update user profile: %v", err)
	}
	return
}
