package repos

import (
	"github.com/jinzhu/gorm"
	"goim-pro/pkg/logs"
)

type GormRepoImpl struct {
	db *gorm.DB
}

func NewMysqlRepo(db *gorm.DB) IGormRepo {
	return &GormRepoImpl{db: db}
}

func (i *GormRepoImpl) InsertOne(target interface{}) (value interface{}, err error) {
	_db := i.db.Create(target)
	if err = _db.Error; err != nil {
		return nil, err
	}
	return _db.Value, nil
}
