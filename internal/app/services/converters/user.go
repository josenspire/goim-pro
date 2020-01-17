package converters

import (
	protos "goim-pro/api/protos/salty"
	"goim-pro/internal/app/constants"
	"goim-pro/internal/app/repos/user"
)

// convert user register entity
func ConvertRegisterUserProfile(profile *protos.UserProfile) user.UserProfile {
	return user.UserProfile{
		UserID:      profile.GetUserID(),
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

func ConvertLoginResp(profile user.UserProfile) *protos.UserProfile {
	return &protos.UserProfile{
		UserID:      profile.UserID,
		Telephone:   profile.Telephone,
		Email:       profile.Email,
		Username:    profile.Username,
		Nickname:    profile.Nickname,
		Avatar:      profile.Avatar,
		Description: profile.Description,
		// TODO:
		//Sex:         constants.UserSex[profile.Sex],
		Birthday:    profile.Birthday,
		Location:    profile.Location,
	}
}
