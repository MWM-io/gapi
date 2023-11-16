package utils

import (
	"regexp"
	"strings"
)

func toTitleCase(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func GenerateOperationID(path string) string {
	re := regexp.MustCompile(`{([^}]+)}`)
	path = re.ReplaceAllString(path, "")

	pathParts := strings.Split(path, "/")
	for i, part := range pathParts {
		if part != "" {
			pathParts[i] = toTitleCase(strings.ToLower(part))
		}
	}
	operationID := strings.Join(pathParts, "")

	operationID = strings.ReplaceAll(operationID, "-", "")
	operationID = strings.ReplaceAll(operationID, "/", "")
	operationID = strings.ReplaceAll(operationID, "_", "")
	operationID = strings.ToLower(operationID[:1]) + operationID[1:]

	return operationID
}
