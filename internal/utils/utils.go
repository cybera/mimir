package utils

import "strings"

func Chomp(s string) string {
	return strings.Trim(s, " \r\n")
}
