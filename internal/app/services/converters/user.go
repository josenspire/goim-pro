package converters

import (
	protos "goim-pro/api/protos/salty"
	"goim-pro/internal/app/constants"
	"goim-pro/internal/app/repos/user"
)

// convert user register entity
func ConvertRegisterUserProfile(profile *protos.UserProfile) user.UserProfile {
	return user.UserProfile{
		Telephone:   profile.GetTelephone(),
		Email:       profile.GetEmail(),
		Username:    profile.GetUsername(),
		Nickname:    profile.GetNickname(),
		Avatar:      profile.GetAvatar(),
		Description: profile.GetDescription(),
		Sex:         constants.USER_SEX[int32(profile.GetSex())],
		// TODO: should parse and format
		//Birthday:    profile.GetBirthday(),
		Location:    profile.GetLocation(),
	}
}
