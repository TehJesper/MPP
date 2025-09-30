package domain

import (
	"strings"
)

func IsValidRace(race string, races []string) bool {
    raceLower := strings.ToLower(race)
    for _, r := range races {
        if strings.ToLower(r) == raceLower {
            return true
        }
    }
    return false
}
