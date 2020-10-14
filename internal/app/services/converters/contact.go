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
		//userProfile := contact.User
		var content = make(map[string]interface{})
		if err := json.Unmarshal([]byte(item.Message.MsgContent), &content); err != nil {
			fmt.Println(err)
			return nil
		}
		protoContactNotificationMsg[i] = &protos.ContactOperationMessage{
			Common: &protos.MessageCommon{
				MessageId:    item.MessageId,
				CreatedTime:  utils.ParseTimeToTimestamp(item.CreatedAt),
				IsNeedRemind: item.IsNeedRemind,
				SortId:       "",
			},
			TriggerProfile: &protos.UserProfile{
				UserId:      "",
				Telephone:   "",
				Email:       "",
				Nickname:    "",
				Avatar:      "",
				Description: "",
				Sex:         0,
				Birthday:    0,
				Location:    "",
			},
			AddReason:    content["addReason"].(string),
			RejectReason: content["rejectReason"].(string),
			Type:         protos.ContactOperationMessage_OperationType(content["operationType"].(float64)),
		}
	}
	return protoContactNotificationMsg
}
