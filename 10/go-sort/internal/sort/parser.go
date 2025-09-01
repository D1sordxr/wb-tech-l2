package sort

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type parser struct{}

var (
	months = map[string]int{
		"jan": 1,
		"feb": 2,
		"mar": 3,
		"apr": 4,
		"may": 5,
		"jun": 6,
		"jul": 7,
		"aug": 8,
		"sep": 9,
		"oct": 10,
		"nov": 11,
		"dec": 12,
	}
	unitMultipliers = map[string]float64{
		"k": 1024,
		"m": 1024 * 1024,
		"g": 1024 * 1024 * 1024,
	}
)

func (*parser) ParseMonth(s string) (int, error) {
	normalized := strings.ToLower(strings.TrimSpace(s))
	if val, ok := months[normalized]; ok {
		return val, nil
	}

	return 0, fmt.Errorf("invalid month: %s", s)
}

func (*parser) ParseHumanNumber(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, errors.New("empty string")
	}

	lower := strings.ToLower(s)
	for unit, multiplier := range unitMultipliers {
		if strings.HasSuffix(lower, unit) {
			numPart := strings.TrimSuffix(lower, unit)
			val, err := strconv.ParseFloat(numPart, 64)
			if err != nil {
				return 0, err
			}
			return val * multiplier, nil
		}
	}

	return strconv.ParseFloat(lower, 64)
}

func (*parser) ParseFloat(s string) (float64, error) {
	return strconv.ParseFloat(strings.TrimSpace(s), 64)
}
