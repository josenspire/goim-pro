package contact


type ClassifyVO struct {
	UserId    string `json:"userId" gorm:"column:userId; type:varchar(32); primary_key; not null"`
	ContactId string `json:"contactId" gorm:"column:contactId; type:varchar(32); primary_key; not null"`

	Status string `json:"status" gorm:"column:status; type:varchar(32);"`
}

