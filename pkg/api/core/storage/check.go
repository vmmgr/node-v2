package storage

import (
	pb "github.com/vmmgr/node/proto/proto-go"
	"os"
)

func storageCmdCheck(data *pb.StorageData) bool {
	if data.GroupID == 0 || data.MaxSize < 1 || data.Name == "" {
		return false
	}
	return true
}

func fileExistsCheck(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	} else {
		return true
	}
}
