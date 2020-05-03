package storage

import "testing"

func TestCopy(t *testing.T) {
	fileCopy("/home/yonedayuto/Downloads/ubuntu-20.04-desktop-amd64.iso", "/home/yonedayuto/test.iso")
}
