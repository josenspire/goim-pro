package converters

import (
	protos "goim-pro/api/protos/salty"
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
