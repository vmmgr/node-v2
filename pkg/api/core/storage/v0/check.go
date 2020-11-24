package v0

import (
	"os"
)

func FileExistsCheck(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	} else {
		return true
	}
}
