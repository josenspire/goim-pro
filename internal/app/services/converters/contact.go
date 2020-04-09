package converters

import (
	protos "goim-pro/api/protos/salty"
	"goim-pro/internal/app/constants"
	"goim-pro/internal/app/models"
	"strings"
)

// convert proto to entity for contact remark profile
func ConvertProto2EntityForRemarkProfile(profile *protos.ContactRemark) models.RemarkProfile {
	return models.RemarkProfile{
		RemarkName:  profile.RemarkName,
		Telephone:   strings.Join(profile.Telephones, ","),
		Description: profile.Description,
		Tags:        strings.Join(profile.Tags, ","),
	}
}

func ConvertEntity2ProtoForContacts(contacts []models.Contact) (protoContacts []*protos.ContactProfile) {
	protoContacts = make([]*protos.ContactProfile, len(contacts))
	for i, contact := range contacts {
		userProfile := contact.User
		protoContacts[i] = &protos.ContactProfile{
			UserProfile: &protos.UserProfile{
				UserId:      contact.UserId,
				Telephone:   contact.Telephone,
				Email:       userProfile.Email,
				Nickname:    userProfile.Nickname,
				Avatar:      userProfile.Avatar,
				Description: userProfile.Description,
				Sex:         constants.UserSexStringMapping[userProfile.Sex],
				Birthday:    userProfile.Birthday,
				Location:    userProfile.Location,
			},
			RemarkInfo: &protos.ContactRemark{
				RemarkName:  contact.RemarkName,
				Description: contact.Description,
				Telephones:  strings.Split(contact.Telephone, ","),
				Tags:        strings.Split(contact.Tags, ","),
			},
		}
	}
	return protoContacts
}
