package stringutil

import "strings"

func ToUpper(text string) string {
	return strings.ToUpper(text)
}

func Join(first, second string) string {
	return first + " " + second
}
