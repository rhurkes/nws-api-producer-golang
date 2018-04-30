package helpers

import (
	"strconv"
	"strings"
)

func NormalizeString(input string, preserveCase bool) string {
	textWithoutBreaks := strings.Replace(input, "\n", " ", -1)
	trimmedText := strings.TrimSpace(textWithoutBreaks)

	for strings.Contains(trimmedText, "  ") {
		trimmedText = strings.Replace(trimmedText, "  ", " ", -1)
	}

	if preserveCase {
		return trimmedText
	}

	return strings.ToLower(trimmedText)
}

func NormalizeFloat(input string) float32 {
	inputString := NormalizeString(input, false)
	num, err := strconv.ParseFloat(inputString, 32)
	if err != nil {
		num = 0
	}

	return float32(num)
}
