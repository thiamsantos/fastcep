package address

import (
	"regexp"
	"strings"
)

var nonDigitsRegex = regexp.MustCompile(`\D+`)

func RemoveNonDigits(rawCep string) string {
	return nonDigitsRegex.ReplaceAllString(rawCep, "")
}

func LeftPadZero(rawCep string, length int) string {
	return strings.Repeat("0", length-len(rawCep)) + rawCep
}
