package models

import (
	"goim-pro/internal/app/repos/base"
	tbl "goim-pro/pkg/db"
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

func (User) TableName() string {
	return tbl.TableUsers
}
