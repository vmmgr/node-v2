package storage

import (
	"fmt"
	"github.com/vmmgr/node/db"
	"github.com/vmmgr/node/etc"
	pb "github.com/vmmgr/node/proto/proto-go"
	"log"
	"strconv"
)

type result struct {
	Info string
	Data string
	Path string
	Err  error
}

type storage struct {
	path   string
	format string
	size   int
}

//Driver 0:qcow2(default) 1:img
func AddStorage(data *pb.StorageData) result {
	var basePath, path string
	driver := 0
	extension := "qcow2"

	if storageCmdCheck(data) == false {
		return result{Err: fmt.Errorf("Error: command wrong... ")}
	}

	if data.GetDriver() == 1 {
		driver = int(data.GetDriver())
		extension = "img"
	}

	//add storage database
	r := db.AddDBStorage(db.Storage{
		GroupID: int(data.GetGroupID()),
		Name:    strconv.Itoa(int(data.GetGroupID())),
		Driver:  driver,
		Mode:    int(data.GetMode()),
		Path:    data.GetPath(),
		MaxSize: int(data.GetMaxSize()),
		Type:    0, Lock: 0})

	if data.GetMode() == 10 {
		path = data.GetPath() + "/" + strconv.Itoa(r.ID) + ".img"
		log.Println("AddStorage: Manual Mode")
	} else {
		basePath = etc.GetDiskPath(int(data.GetMode()))
		log.Println("AddStorage: Auto Mode " + strconv.Itoa(int(data.GetMode())))
		if basePath == "" {
			return result{Err: fmt.Errorf("Error: no diskpath on configfile ")}
		}
		log.Println("basePath: " + basePath)
		path = basePath + "/" + strconv.Itoa(r.ID) + "."
	}

	path += extension

	//apply vm image or not
	if data.GetImage() == "" {
		createStorageCmd(storage{path: path, format: extension, size: int(data.GetMaxSize())})
	} else {
		fileCopy(data.GetImage(), path)
	}

	return result{Path: path, Err: nil}
}

//func DeleteStorage(data *pb.StorageData) result {
////	add storage database
//r := db.DeleteDBStorage(db.Storage{ID: int(data.GetID())})
//
//
//
//return result{Path: path, Err: nil}
//}
