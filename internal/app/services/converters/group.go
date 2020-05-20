package converters

import (
	protos "goim-pro/api/protos/salty"
	"goim-pro/internal/app/constants"
	"goim-pro/internal/app/models"
	"goim-pro/pkg/utils"
)

func ConvertProto2EntityForGroupProfile() {

}

// ConvertEntity2ProtoForGroupProfile - convert group from models schema to proto schema
func ConvertEntity2ProtoForGroupProfile(groupProfile *models.Group) *protos.GroupProfile {
	return &protos.GroupProfile{
		GroupId:     groupProfile.GroupId,
		Name:        groupProfile.Name,
		CreateTime:  utils.ParseTimeToTimestamp(groupProfile.CreatedAt),
		OwnerUserId: groupProfile.OwnerUserId,
		Avatar:      groupProfile.Avatar,
		Notice:      groupProfile.Notice,
		Members:     convertEntity2ProtoForMemberProfiles(groupProfile.Members),
	}
}

func convertEntity2ProtoForMemberProfiles(members []models.Member) (protoMembers []*protos.GroupMemberProfile) {
	protoMembers = make([]*protos.GroupMemberProfile, len(members))
	for i, member := range members {
		protoMembers[i] = &protos.GroupMemberProfile{
			GroupId:     member.GroupId,
			Alias:       member.Alias,
			Role:        consts.MemberRoleStringMapping[member.Role],
			JoinTime:    0,
			UserProfile: ConvertEntity2ProtoForUserProfile(&member.User.UserProfile),
		}
	}
	return protoMembers
}
