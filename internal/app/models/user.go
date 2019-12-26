package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	"goim-pro/configs"
	"goim-pro/pkg/db"
)

type User struct {
	UserID   uint64 `json:"userId" gorm:"column:userId; primary_key; not null"`
	Password string `json:"password" gorm:"column:password; type:varchar(255); not null"`
	Role     string `json:"role" gorm:"column:role; type:ENUM('0', '10', '99'); default:'0'"`
	Status   string `json:"status" gorm:"column:status; type:ENUM('0', '1'); default: '0'; not null"`
	UserProfile
	BaseModel
}

type UserProfile struct {
	Telephone string `json:"telephone" gorm:"column:telephone; type:varchar(11)"`
	Email     string `json:"email" gorm:"column:email; type:varchar(100)"`
	Username  string `json:"username" gorm:"column:username; type:varchar(18)"`
	Nickname  string `json:"nickname" gorm:"column:nickname; type:varchar(16)"`
	Avatar    string `json:"avatar" gorm:"column:avatar; type:varchar(255)"`
	Signature string `json:"signature" gorm:"column:signature; type:varchar(255)"`
	Sex       string `json:"sex" gorm:"column:sex; type: ENUM('0', '1'); default:'0'"`
	Birthday  string `json:"birthday" gorm:"column:birthday; type: varchar(12)"`
	Location  string `json:"location" gorm:"column:location; type: varchar(255)"`
}

var aesSecretKey string
var mysqlDB *gorm.DB

func init() {
	aesSecretKey = configs.GetApiSecretKey()
	mysqlDB = db.GetMysqlConnection().GetMysqlDBInstance()
}

func (u *User) Register(newUser User) (err error) {
	isExist, err := u.IsTelephoneRegistered(newUser.Telephone)
	if err != nil {
		return
	}
	if isExist {
		return errors.New("user already register, please login")
	}
	mysqlDB.Create(&newUser)
	return
}

func (u *User) IsTelephoneRegistered(telephone string) (bool, error) {
	err := mysqlDB.Where("telephone = ?", telephone).Find(&u).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
