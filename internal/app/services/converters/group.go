package converters

import (
	protos "goim-pro/api/protos/salty"
	"goim-pro/internal/app/constants"
	. "goim-pro/internal/app/repos/group"
	"goim-pro/pkg/utils"
)

func ConvertProto2EntityForGroupProfile() {

}

func ConvertEntity2ProtoForGroupProfile(groupProfile *Group) *protos.GroupProfile {
	return &protos.GroupProfile{
		GroupId:     groupProfile.GroupId,
		Name:        groupProfile.Name,
		CreateTime:  utils.ParseTimeToTimestamp(groupProfile.CreatedAt),
		OwnerUserId: groupProfile.OwnerUserId,
		Avatar:      groupProfile.Avatar,
		Notice:      groupProfile.Notice,
		Members:     ConvertEntity2ProtoForMemberProfiles(groupProfile.Members),
	}
}

// TODO: should fill with member profile
func ConvertEntity2ProtoForMemberProfiles(members []*Member) (protoMembers []*protos.GroupMemberProfile) {
	protoMembers = make([]*protos.GroupMemberProfile, len(members))
	for i, member := range members {
		protoMembers[i] = &protos.GroupMemberProfile{
			GroupId:     "",
			Alias:       member.Alias,
			Role:        constants.MemberRoleStringMapping[member.Role],
			JoinTime:    0,
			UserProfile: ConvertEntity2ProtoForUserProfile(nil),
		}
	}
	return protoMembers
}
