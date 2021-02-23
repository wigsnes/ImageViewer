package stringo

import "strings"

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if strings.Contains(strings.ToLower(a), strings.ToLower(b)) {
			return true
		}
	}
	return false
}

func StringContains(a string, b string) bool {
	if strings.Contains(strings.ToLower(a), strings.ToLower(b)) {
		return true
	}
	return false
}
