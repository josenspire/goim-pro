package converters

import (
	"goim-pro/api/protos"
	"goim-pro/internal/app/constants"
	"goim-pro/internal/app/repos/user"
)

// convert user register entity
func ConvertRegisterUserProfile(profile *protos.UserProfile) user.UserProfile {
	return user.UserProfile{
		Telephone: profile.GetTelephone(),
		Email:     profile.GetEmail(),
		Username:  profile.GetUsername(),
		Nickname:  profile.GetNickname(),
		Avatar:    profile.GetAvatar(),
		Signature: profile.GetSignature(),
		Sex:       constants.USER_SEX[int32(profile.GetSex())],
		Birthday:  profile.GetBirthday().String(),
		Location:  profile.GetLocation(),
	}
}
