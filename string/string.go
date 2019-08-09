package gorc

import "strings"

// StringBuilder return concatenated string from given strings
func StringBuilder(strs ...string) string {
	var sb strings.Builder
	for _, str := range strs {
		sb.WriteString(str)
	}
	return sb.String()
}
