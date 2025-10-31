package domain

import (
	"strings"
)

func IsValidRace(race string, races []string) bool {
	normalizedRace := strings.ToLower(strings.ReplaceAll(race, " ", "-"))

	for _, r := range races {
		normalizedR := strings.ToLower(strings.ReplaceAll(r, " ", "-"))

		if normalizedRace == normalizedR {
			return true
		}
	}

	return false
}
