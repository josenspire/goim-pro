package converters

import (
	protos "goim-pro/api/protos/salty"
	"goim-pro/internal/app/constants"
	"goim-pro/internal/app/repos/user"
)

// convert user register entity
func ConvertRegisterUserProfile(profile *protos.UserProfile) user.UserProfile {
	return user.UserProfile{
		UserId:      profile.GetUserId(),
		Telephone:   profile.GetTelephone(),
		Email:       profile.GetEmail(),
		Nickname:    profile.GetNickname(),
		Avatar:      profile.GetAvatar(),
		Description: profile.GetDescription(),
		Sex:         constants.UserSexProtoMapping[profile.GetSex()],
		Birthday:    profile.GetBirthday(),
		Location:    profile.GetLocation(),
	}
}

func ConvertLoginResp(profile *user.UserProfile) *protos.UserProfile {
	return &protos.UserProfile{
		UserId:      profile.UserId,
		Telephone:   profile.Telephone,
		Email:       profile.Email,
		Nickname:    profile.Nickname,
		Avatar:      profile.Avatar,
		Description: profile.Description,
		Sex:         constants.UserSexStringMapping[profile.Sex],
		Birthday:    profile.Birthday,
		Location:    profile.Location,
	}
}
