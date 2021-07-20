package consts

import protos "goim-pro/api/protos/salty"

var UserSexProtoMapping = map[protos.UserProfile_Sex]string{
	protos.UserProfile_NOT_SET: "MALE",
	protos.UserProfile_MALE:    "MALE",
	protos.UserProfile_FEMALE:  "FEMALE",
}

var UserSexStringMapping = map[string]protos.UserProfile_Sex{
	"MALE":   protos.UserProfile_MALE,
	"FEMALE": protos.UserProfile_FEMALE,
}

var MemberRoleProtoMapping = map[protos.GroupMemberProfile_GroupRole]string{
	protos.GroupMemberProfile_NONE:          "1",
	protos.GroupMemberProfile_ADMINISTRATOR: "99",
}

var MemberRoleStringMapping = map[string]protos.GroupMemberProfile_GroupRole {
	"1": protos.GroupMemberProfile_NONE,
	"99": protos.GroupMemberProfile_ADMINISTRATOR,
}