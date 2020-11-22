package v0

import (
	"os"
)

func fileExistsCheck(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	} else {
		return true
	}
}
