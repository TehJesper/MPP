package domain

import (
	"strings"
)

func IsValidClass(class string, classes []string) bool {
	classLower := strings.ToLower(class)
	for _, r := range classes {
		if strings.ToLower(r) == classLower {
			return true
		}
	}
	return false
}
