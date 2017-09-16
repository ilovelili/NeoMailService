// Package util utilities
package util

import "os"

// GetWorkingDirectory get working directory
func GetWorkingDirectory() string {
	pwd, _ := os.Getwd()
	return pwd
}
