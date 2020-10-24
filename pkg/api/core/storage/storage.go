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

type Tmp struct {
	Info string
}

//Driver 0:qcow2(default) 1:img
func (t *Tmp) AddStorage(data *pb.StorageData) result {
	var basePath, path string
	if storageCmdCheck(data) == false {
		return result{Err: fmt.Errorf("Error: command wrong... ")}
	}
	d := getDriver(int(data.GetDriver()))
	manualPath := data.GetPath()
	if 0 < data.GetMode() && data.GetMode() < 10 {
		manualPath = ""
	}

	//add storage database
	r := db.AddDBStorage(db.Storage{
		GroupID: int(data.GetGroupID()),
		Name:    strconv.Itoa(int(data.GetGroupID())),
		Driver:  d.driver,
		Mode:    int(data.GetMode()),
		Path:    manualPath,
		MaxSize: int(data.GetMaxSize()),
		Type:    0, Lock: 0})

	if data.GetMode() == 10 {
		path = manualPath + "/" + strconv.Itoa(r.ID) + "."
		log.Println("AddStorage: Manual Mode")
	} else if 0 < data.GetMode() && data.GetMode() < 10 {
		basePath = etc.GetDiskPath(int(data.GetMode()))
		log.Println("AddStorage: Auto Mode " + strconv.Itoa(int(data.GetMode())))
		if basePath == "" {
			return result{Err: fmt.Errorf("Error: no diskpath on configfile ")}
		}
		log.Println("basePath: " + basePath)
		path = basePath + "/" + strconv.Itoa(r.ID) + "."
	} else {
		return result{Err: fmt.Errorf("Error: mode value error ")}
	}

	path += d.extension
	//apply vm image or not
	if data.GetImage() == "" {
		createStorageCmd(storage{path: path, format: d.extension, size: int(data.GetMaxSize())})
	} else {
		r := t.fileCopy(data.GetImage(), path)
		if r.Err != nil {
			return result{Info: r.Info, Err: r.Err}
		}
	}
	return result{Info: "OK", Err: nil}
}

func DeleteStorage(data *pb.StorageData) result {
	var path string
	dbData, err := db.SearchDBStorage(db.Storage{ID: int(data.GetID() / 100000)})
	if err != nil {
		return result{Err: fmt.Errorf("Error: store read error ")}
	}
	if 0 < dbData.Mode && dbData.Mode < 10 {
		basePath := etc.GetDiskPath(dbData.Mode)
		if basePath == "" {
			return result{Err: fmt.Errorf("Error: no diskpath on configfile ")}
		}
		path = basePath + "/" + strconv.Itoa(int(data.GetID())) + "."
	} else if dbData.Mode == 10 {
		path = dbData.Path + "/" + strconv.Itoa(int(data.GetID())) + "."
	}
	d := getDriver(int(data.GetDriver()))
	path += d.extension
	log.Println("Path: " + path)

	if cmdResult := deleteStorageCmd(storage{path: path}); cmdResult.Err != nil {
		return result{Err: fmt.Errorf("Error: delete failed !! ")}
	}

	if deleteResult := db.DeleteDBStorage(db.Storage{ID: int(data.GetID())}); deleteResult.Error != nil {
		return result{Err: fmt.Errorf("Error: delete store error !! ")}
	}
	return result{Info: "OK", Err: nil}
}

// #4 Test
func UpdateStorage(data *pb.StorageData) result {
	var path string
	dbData, err := db.SearchDBStorage(db.Storage{ID: int(data.GetID() / 100000)})
	if err != nil {
		return result{Info: "Error: store read error", Err: err}
	}
	if dbData.Lock == 1 {
		return result{Err: fmt.Errorf("Error: disk is locked ")}
	}

	if 0 < dbData.Mode && dbData.Mode < 10 {
		basePath := etc.GetDiskPath(dbData.Mode)
		if basePath == "" {
			return result{Err: fmt.Errorf("Error: no diskpath on configfile ")}
		}
		path = basePath + "/" + strconv.Itoa(int(data.GetID())) + "."
	} else if dbData.Mode == 10 {
		path = dbData.Path + "/" + strconv.Itoa(int(data.GetID())) + "."
	}
	d := getDriver(int(data.GetDriver()))
	path += d.extension
	log.Println("Path: " + path)

	if dbData.MaxSize < int(data.GetMaxSize()) {
		if cmdResult := resizeStorageCmd(storage{path: path, size: int(data.MaxSize)}); cmdResult.Err != nil {
			return result{Info: "Error: storage resize failed !! ", Err: cmdResult.Err}
		}
		dbData.MaxSize = int(data.GetMaxSize())
	}

	if data.GetPath() != "" {
		dbData.Path = data.GetPath()
	}
	if data.GetName() != "" {
		dbData.Name = data.GetName()
	}
	if data.GetGroupID() != 0 {
		dbData.GroupID = int(data.GetGroupID())
	}
	if data.GetMode() != 0 {
		dbData.Mode = int(data.GetMode())
	}
	if data.GetDriver() != 0 {
		dbData.Driver = int(data.GetDriver())
	}

	if r := db.UpdateDBStorage(dbData); r.Error != nil {
		return result{Err: err}
	} else {
		return result{Info: "OK", Err: nil}
	}
}
