package models

import (
	"goim-pro/internal/app/repos/base"
	tbl "goim-pro/pkg/db"
)

type Contact struct {
	// Composite primary key - userId+contactId
	UserId    string `json:"userId" gorm:"column:userId; type:varchar(32); primary_key; not null"`
	ContactId string `json:"contactId" gorm:"column:contactId; type:varchar(32); primary_key; not null"`
	Status    string `json:"status" gorm:"column:status; type:varchar(32);"`
	RemarkProfile
	base.BaseModel

	User User `gorm:"ForeignKey:ContactId;AssociationForeignKey:UserId"`
}

type RemarkProfile struct {
	RemarkName  string `json:"remarkName" gorm:"column:remarkName; type:varchar(16)"`
	Telephone   string `json:"telephone" gorm:"column:telephone; type:varchar(255)"` // can support multiple tel, split by `;`
	Description string `json:"description" gorm:"column:description; type:varchar(255)"`
	Tags        string `json:"tags" gorm:"column:tags; type:varchar(255)"` // can support multiple tag, split by `;`
}

func (Contact) TableName() string {
	return tbl.TableContacts
}
