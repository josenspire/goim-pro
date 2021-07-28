package user

import (
	"github.com/jinzhu/gorm"
	"goim-pro/internal/app/models"
	tbl "goim-pro/internal/db"
	"goim-pro/pkg/logs"
)

var logger = logs.GetLogger("ERROR")
var mysqlDB *gorm.DB

type UserImpl models.User

type IUserRepo interface {
	IsTelephoneOrEmailRegistered(telephone string, email string) (bool, error)
	Register(newUser *models.User) error
	QueryByEmailAndPassword(email string, password string) (user *models.User, err error)
	QueryByTelephoneAndPassword(telephone string, password string) (user *models.User, err error)
	RemoveUserByUserId(userId string, isForce bool) error
	ResetPasswordByTelephone(telephone string, newPassword string) error
	ResetPasswordByEmail(email string, newPassword string) error
	FindByUserId(userId string) (*models.User, error)
	FindOneUser(condition interface{}) (*models.User, error)
	FindOneAndUpdateProfile(condition interface{}, profile interface{}) error
}

func NewUserRepo(db *gorm.DB) IUserRepo {
	mysqlDB = db
	return &UserImpl{}
}

func (u *UserImpl) Register(newUser *models.User) (err error) {
	_db := mysqlDB.Create(&newUser)
	if _db.Error != nil {
		err = _db.Error
		logger.Errorf("create user error: %s", err.Error())
	}
	return
}

func (u *UserImpl) IsTelephoneOrEmailRegistered(telephone string, email string) (bool, error) {
	var isTelExist bool = true
	var err error
	if telephone != "" {
		err = mysqlDB.First(&models.User{}, "telephone = ?", telephone).Error
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
		err = mysqlDB.First(&models.User{}, "email = ?", email).Error
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

func (u *UserImpl) QueryByEmailAndPassword(email string, enPassword string) (*models.User, error) {
	var user = &models.User{}
	var err error
	db := mysqlDB.First(&user, "email = ? and password = ?", email, enPassword)
	if db.RecordNotFound() {
		return nil, nil
	}
	if err = db.Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserImpl) QueryByTelephoneAndPassword(telephone string, enPassword string) (*models.User, error) {
	var user = &models.User{}
	var err error
	db := mysqlDB.First(user, "telephone = ? and password = ?", telephone, enPassword)
	if db.RecordNotFound() {
		return nil, nil
	}
	if err = db.Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserImpl) RemoveUserByUserId(userId string, isForce bool) (err error) {
	_db := mysqlDB
	if isForce {
		_db = mysqlDB.Unscoped()
	}
	_db = _db.Delete(&models.User{}, "userId = ?", userId)
	if _db.RecordNotFound() {
		logger.Warningln("remove user fail, userId not found")
	} else if _db.Error != nil {
		logger.Errorf("error happened to remove user: %v", _db.Error)
		err = _db.Error
	}
	return
}

func (u *UserImpl) ResetPasswordByTelephone(telephone string, newPassword string) (err error) {
	db := mysqlDB.Model(&models.User{}).Where("telephone = ?", telephone).Update("password", newPassword)
	if err = db.Error; err != nil {
		logger.Errorf("error happened to reset password by telephone: %v", err)
	}
	return
}

func (u *UserImpl) ResetPasswordByEmail(email string, newPassword string) (err error) {
	db := mysqlDB.Model(&models.User{}).Where("email = ?", email).Update("password", newPassword)
	if err = db.Error; err != nil {
		logger.Errorf("error happened to reset password by email: %v", err)
	}
	return
}

func (u *UserImpl) FindByUserId(userId string) (user *models.User, err error) {
	user = &models.User{}
	db := mysqlDB.First(user, "userId = ?", userId)
	if db.RecordNotFound() {
		return nil, nil
	}
	if err = db.Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserImpl) FindOneUser(condition interface{}) (user *models.User, err error) {
	user = &models.User{}
	db := mysqlDB.Where(condition).First(&user)
	if db.RecordNotFound() {
		return nil, nil
	}
	if err = db.Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserImpl) FindOneAndUpdateProfile(condition interface{}, profile interface{}) (err error) {
	db := mysqlDB.Table(tbl.TableUsers).Where(condition).Update(profile)
	if err = db.Error; err != nil {
		logger.Errorf("error happened to update user profile: %v", err)
	}
	return
}
