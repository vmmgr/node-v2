package manage

import (
	"fmt"
	"github.com/vmmgr/node/db"
	"github.com/vmmgr/node/etc"
	pb "github.com/vmmgr/node/proto/proto-go"
	"log"
	"os/exec"
	"strconv"
)

type result struct {
	Data string
	Path string
	Err  error
}

type storage struct {
	path   string
	format int
	size   int
}

func RunStorageCmd(cmd []string) {
	out, _ := exec.Command("qemu-img", cmd...).Output()
	log.Println(string(out))
}

func GetMainStorage(data *pb.VMData) result {
	var basePath, path string
	sp := data.Storage[0].GetMode()
	r := db.AddDBStorage(db.Storage{
		GroupID: int(data.Storage[0].GetGroupID()),
		Name:    strconv.Itoa(int(data.Storage[0].GetGroupID())),
		Driver:  int(data.Storage[0].GetDriver()),
		Mode:    int(data.Storage[0].GetMode()),
		Path:    data.Storage[0].GetPath(),
		MaxSize: int(data.Storage[0].GetMaxSize()),
		Type:    0, Lock: 0})

	if sp == 10 {
		path = data.Storage[0].GetPath() + "/" + strconv.Itoa(r.ID) + ".img"
		log.Println("Storage: Manual Mode")
	} else {
		basePath = etc.GetDiskPath(int(data.Storage[0].GetMode()))
		log.Println("Storage: Auto Mode " + strconv.Itoa(int(data.Storage[0].GetMode())))
		if basePath == "" {
			return result{Path: "", Err: fmt.Errorf("Error: no diskpath on configfile ")}
		}
		log.Println("basePath: " + basePath)
		path = basePath + "/" + strconv.Itoa(r.ID) + ".img"
	}
	return result{Path: path, Err: nil}
}

func StorageProcess(data *pb.VMData) result {
	for _, data := range data.Storage {
		var basePath, path string

		r := db.AddDBStorage(db.Storage{
			GroupID: int(data.GetGroupID()),
			Name:    strconv.Itoa(int(data.GetGroupID())),
			Driver:  int(data.GetDriver()),
			Mode:    int(data.GetMode()),
			Path:    data.GetPath(),
			MaxSize: int(data.GetMaxSize()),
			Type:    0, Lock: 0})
		if data.Mode == 10 {
			path = data.GetPath() + "/" + strconv.Itoa(r.ID) + ".img"
			log.Println("Storage: Manual Mode")
		} else {
			basePath = etc.GetDiskPath(int(data.GetMode()))
			log.Println("Storage: Auto Mode " + strconv.Itoa(int(data.GetMode())))
			if basePath == "" {
				return result{Err: fmt.Errorf("Error: no diskpath on configfile ")}
			}
			log.Println("basePath: " + basePath)
			path = basePath + "/" + strconv.Itoa(r.ID) + ".img"
		}
		err := CreateStorage(storage{path: path, format: int(data.Mode), size: int(data.MaxSize)})
		if err != nil {
			return result{Err: err}
		}

	}
	return result{Err: nil}
}

//path, name string, format, size int
func CreateStorage(s storage) error {
	log.Println("----storage create----")
	var extension string
	if s.size < 0 {
		return fmt.Errorf("Wrong storage size !!")
	}
	//0: virtio 1:img
	if s.format == 1 {
		extension = "img"
	} else {
		extension = "virtio"
	}

	var cmd []string
	//qemu-img create [-f format] filename [size]
	cmdArray := []string{"create", "-f", extension, s.path, strconv.Itoa(s.size) + "M"}

	log.Println(cmdArray)
	cmd = append(cmd, cmdArray...)
	RunStorageCmd(cmd)

	return nil
}

func DeleteStorage(s storage) error {
	var cmd []string

	if FileExistsCheck(s.path) {
		cmd = append(cmd, "info")
		cmd = append(cmd, s.path+".img")
		RunStorageCmd(cmd)

		return nil
	}

	return fmt.Errorf("File not exists !! ")
}

func ResizeStorage(s storage) error {
	//qemu-img resize [filename] [size]

	var cmd []string

	cmd = append(cmd, "qemu-img")
	cmd = append(cmd, "resize")
	cmd = append(cmd, s.path)
	cmd = append(cmd, strconv.Itoa(s.size)+"M")

	RunStorageCmd(cmd)

	return nil
}

func InformationStorage(s storage) error {
	//qemu-img info [-f format] [filename]
	var cmd []string

	cmd = append(cmd, "qemu-img")
	cmd = append(cmd, "info")
	cmd = append(cmd, s.path)

	RunStorageCmd(cmd)
	return nil
}
