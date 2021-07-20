package utils

import (
	"github.com/mozillazg/go-pinyin"
	"regexp"
	"strings"
)

var (
	numberPattern = "^[0-9]+$"
	letterPattern = "^[a-zA-Z]+$"
)

const (
	defaultLetter = "#"
)

/**
 *  Get hans first letter,
 *  1. if first chat is Letter that will return Directly
 *  2. if first chat is Number or can not Decode Pinyin succeed, will return default string "#"
 */
func GetStringLetters(zhStr string) (letters string) {
	if len(zhStr) == 0 {
		return ""
	}

	a := pinyin.NewArgs()
	a.Style = pinyin.FirstLetter

	for i, r := range zhStr {
		match, result := checkNumberOrLetterPattern(string(r), i == 0)
		if match {
			letters += result
			continue
		}

		hansLetter := pinyin.SinglePinyin(r, a)
		if len(hansLetter) == 0 {
			letters += defaultLetter
			continue
		}
		letters += hansLetter[0]
	}

	return strings.ToUpper(letters)
}

func checkNumberOrLetterPattern(str string, isFirstRune bool) (match bool, result string) {
	matchNum, err := regexp.MatchString(numberPattern, str)
	if err != nil || matchNum {
		var s = str
		if isFirstRune {
			s = defaultLetter + str
		}
		return true, s
	}
	matchLet, err := regexp.MatchString(letterPattern, str)
	if err != nil {
		return true, str
	}
	if matchLet {
		return true, str
	}
	return false, ""
}
