package utils

import (
	"fmt"
	"strconv"
)

func ParseAndValidateFloat(fieldName string, value string, min float64, max float64) (float64, error) {
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, err
	}
	if v < min || v > max {
		return 0, fmt.Errorf("%s out of bounds, should be between %.1f and %.1f", fieldName, min, max)
	}
	return v, nil
}

func ParseAndValidateInt(fieldName string, value string, min int64, max int64) (int64, error) {
	v, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, err
	}
	if v < min || v > max {
		return 0, fmt.Errorf("%s out of bounds, should be between %d and %d", fieldName, min, max)
	}
	return v, nil
}

func WordCount(s string) int {
	count := 0
	inWord := false
	for _, c := range s {
		if c == ' ' || c == '\n' || c == '\t' {
			if inWord {
				count += 1
				inWord = false
			}
		} else {
			inWord = true
		}
	}
	return count
}

