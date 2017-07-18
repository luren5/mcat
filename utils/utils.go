package utils

import "strings"

func Remove0x(s string) string {
	return strings.TrimPrefix(s, "0x")
}
