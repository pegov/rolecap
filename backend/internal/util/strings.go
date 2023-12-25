package util

import (
	"strings"
	"unicode"
)

func RemoveDoubleWhiteSpaces(s string) string {
	words := strings.Fields(s)
	return strings.Join(words, " ")
}

func TrimSpaceAndRemoveDoubleSpaces(s string) string {
	return RemoveDoubleWhiteSpaces(strings.TrimSpace(s))
}

func HasWhiteSpaceLeft(s string) bool {
	if len(s) == 0 {
		return false
	} else if len(s) == 1 {
		return unicode.IsSpace(rune(s[0]))
	} else {
		return unicode.IsSpace(rune(s[0]))
	}
}

func HasWhiteSpaceRight(s string) bool {
	if len(s) == 0 {
		return false
	} else if len(s) == 1 {
		return unicode.IsSpace(rune(s[0]))
	} else {
		return unicode.IsSpace(rune(s[len(s)-1]))
	}
}
