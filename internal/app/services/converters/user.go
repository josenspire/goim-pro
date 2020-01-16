package converters

import (
	protos "goim-pro/api/protos/salty"
	"goim-pro/internal/app/constants"
	"goim-pro/internal/app/repos/user"
	"goim-pro/pkg/logs"
	"strconv"
)

var logger = logs.GetLogger("ERROR")

// convert user register entity
func ConvertRegisterUserProfile(profile *protos.UserProfile) user.UserProfile {
	userID, err := strconv.ParseUint(profile.GetUserID(), 10, 64)
	if err != nil {
		logger.Errorf("parsing userID error: %s", err.Error())
	}
	return user.UserProfile{
		UserID:      userID,
		Telephone:   profile.GetTelephone(),
		Email:       profile.GetEmail(),
		Username:    profile.GetUsername(),
		Nickname:    profile.GetNickname(),
		Avatar:      profile.GetAvatar(),
		Description: profile.GetDescription(),
		Sex:         constants.USER_SEX[int32(profile.GetSex())],
		Birthday:    profile.GetBirthday(),
		Location:    profile.GetLocation(),
	}
}
