package utils

import (
	"fmt"
	"testing"
)

type TestUserProfile struct {
	UserId      string
	Telephone   string
	Email       string
	Nickname    string
	Avatar      string
	Description string
	Sex         string
	Birthday    int64
	Location    string
}

func TestDeepEqual(t *testing.T) {

	//result := DeepEqual(user1, user2)
	//fmt.Println(result)
}

func TestRemoveMapProperties(t *testing.T) {
	profile := TestUserProfile{
		UserId:      "2",
		Telephone:   "13631210111",
		Email:       "29400123@qq.com",
		Nickname:    "TEST03",
		Avatar:      "",
		Description: "Never Settle ..",
		Sex:         "0",
		Birthday:    1578903121862,
		Location:    "",
	}

	profileMap := TransformStructToMap(profile)

	RemoveMapProperties(profileMap, "UserId", "Telephone", "Email")

	fmt.Println(profileMap)
}
