package manage

import (
	"github.com/jinzhu/gorm"
	protos "goim-pro/api/protos/salty"
	"goim-pro/config"
	cserr "goim-pro/internal/app/models/errors"
	"goim-pro/internal/app/repos/user"
	redsrv "goim-pro/pkg/db/redis"
	errmsg "goim-pro/pkg/errors"
	"goim-pro/pkg/logs"
	"goim-pro/pkg/utils"
)

var (
	logger = logs.GetLogger("INFO")

	myRedis redsrv.IMyRedis
	mysqlDB *gorm.DB
)

type ManageService struct {
}

func NewManageService(db *gorm.DB) *ManageService {
	mysqlDB = db

	return &ManageService{}
}

func (s *ManageService) ResetPasswordAndUpdateUserProfile(telephone, email, newPassword, deviceId string, osVersion protos.GrpcReq_OS) error {
	var err error

	tx := mysqlDB.Begin()
	userRepo := user.NewUserRepo(tx)

	enNewPassword := utils.NewSHA256(newPassword, config.GetApiSecretKey())
	if telephone != "" {
		if err = userRepo.ResetPasswordByTelephone(telephone, enNewPassword); err != nil {
			logger.Errorf("reset password error: %s", err.Error())
			tx.Rollback()
			return cserr.NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, err)
		}
	} else {
		if err = userRepo.ResetPasswordByEmail(email, enNewPassword); err != nil {
			logger.Errorf("reset password error: %s", err.Error())
			tx.Rollback()
			return cserr.NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, err)
		}
	}

	// update device profile
	var userCriteria = make(map[string]interface{})
	if !utils.IsEmptyStrings(telephone) {
		userCriteria["Telephone"] = telephone
	} else if !utils.IsEmptyStrings(email) {
		userCriteria["Email"] = email
	}
	if len(userCriteria) == 0 {
		tx.Rollback()
		return errmsg.ErrIllegalOperation
	}
	var updated = map[string]interface{}{
		"DeviceId": deviceId,
	}
	if err := userRepo.FindOneAndUpdateProfile(userCriteria, updated); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return err
}
