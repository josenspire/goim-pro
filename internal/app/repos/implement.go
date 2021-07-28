package repos

import (
	. "goim-pro/internal/db/mysql"
)

type GormRepoImpl struct {
	db *BaseMysql
}

func NewMysqlRepo(db *BaseMysql) IGormRepo {
	return &GormRepoImpl{db: db}
}

func (i *GormRepoImpl) InsertOne(target interface{}) (value interface{}, err error) {
	_db := i.db.Create(target)
	if err = _db.Error; err != nil {
		return nil, err
	}
	return _db.Value, nil
}
