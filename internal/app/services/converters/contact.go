package converters

import (
	"encoding/json"
	"fmt"
	protos "goim-pro/api/protos/salty"
	"goim-pro/internal/app/constants"
	"goim-pro/internal/app/models"
	"goim-pro/pkg/utils"
	"strings"
)

// convert proto to entity for contact remark profile
func ConvertProto2EntityForRemarkProfile(profile *protos.ContactRemark) models.RemarkProfile {
	remarkProfile := models.RemarkProfile{
		RemarkName:  profile.RemarkName,
		Description: profile.Description,
	}
	if profile.Telephones != nil && len(profile.Telephones) > 0 {
		remarkProfile.Telephone = strings.Join(profile.Telephones, ",")
	}
	if profile.Tags != nil && len(profile.Tags) > 0 {
		remarkProfile.Tags = strings.Join(profile.Tags, ",")
	}
	return remarkProfile
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
				Sex:         consts.UserSexStringMapping[userProfile.Sex],
				Birthday:    userProfile.Birthday,
				Location:    userProfile.Location,
			},
			RemarkInfo: &protos.ContactRemark{
				RemarkName:  contact.RemarkName,
				Description: contact.Description,
				Telephones:  strings.Split(contact.Telephone, ","),
				Tags:        strings.Split(contact.Tags, ","),
			},
			SortId: utils.GetStringLetters(userProfile.Nickname),
		}
	}
	return protoContacts
}

func ConvertEntity2ProtoForNotificationMsg(notifications []models.Notification) (protoContactNotificationMsg []*protos.ContactOperationMessage) {
	protoContactNotificationMsg = make([]*protos.ContactOperationMessage, len(notifications))
	for i, item := range notifications {
		var content = make(map[string]interface{})
		if err := json.Unmarshal([]byte(item.Message.MsgContent), &content); err != nil {
			fmt.Println(err)
			return nil
		}
		sender := item.Sender
		protoContactNotificationMsg[i] = &protos.ContactOperationMessage{
			Common: &protos.MessageCommon{
				MessageId:    item.MessageId,
				CreatedTime:  utils.ParseTimeToTimestamp(item.CreatedAt),
				IsNeedRemind: item.IsNeedRemind,
				SortId:       "",
			},
			TriggerProfile: &protos.UserProfile{
				UserId:      sender.UserId,
				Telephone:   sender.Telephone,
				Email:       sender.Email,
				Nickname:    sender.Nickname,
				Avatar:      sender.Avatar,
				Description: sender.Description,
				Sex:         consts.UserSexStringMapping[sender.Sex],
				Birthday:    sender.Birthday,
				Location:    sender.Location,
			},
			AddReason:    content["addReason"].(string),
			RejectReason: content["rejectReason"].(string),
			Type:         protos.ContactOperationMessage_OperationType(item.Message.MsgOperation),
		}
	}
	return protoContactNotificationMsg
}
