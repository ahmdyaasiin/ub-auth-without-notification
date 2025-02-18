package internal

import (
	"fmt"
	"strings"
)

func GetSubstringBetween(prefix, suffix, input string) (string, error) {
	firstPart := strings.Split(input, prefix)
	if len(firstPart) < 2 {
		return "", fmt.Errorf("failed to find prefix: %s", prefix)
	}

	secondPart := strings.Split(firstPart[1], suffix)
	if len(secondPart) < 2 {
		return "", fmt.Errorf("failed to find suffix: %s", suffix)
	}

	return secondPart[0], nil
}

func PascalCase(input string) string {
	words := strings.Split(input, " ")

	var result []string
	for _, word := range words {
		r := strings.ToUpper(string(word[0])) + strings.ToLower(word[1:]) + " "
		result = append(result, r)
	}

	return strings.TrimRight(strings.Join(result, ""), " ")
}
