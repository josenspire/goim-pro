package converters

import (
	protos "goim-pro/api/protos/salty"
	"goim-pro/internal/app/repos/contact"
	"strings"
)

// convert proto to entity for contact remark profile
func ConvertProto2EntityForRemarkProfile(profile *protos.ContactRemark) *contact.RemarkProfile {
	return &contact.RemarkProfile{
		RemarkName:  profile.RemarkName,
		Telephone:   strings.Join(profile.Telephones, ","),
		Description: profile.Description,
		Tags:        strings.Join(profile.Tags, ""),
	}
}
