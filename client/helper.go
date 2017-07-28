package client

import "strings"

// Helper functions.

// Remove spaces etc.
func cleanFileName(fileName string) string {
	fileName = strings.Trim(fileName, " ")
	fileName = strings.Replace(fileName, " ", "-", -1)

	return fileName
}
