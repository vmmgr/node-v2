package vm

import "C"
import (
	"fmt"
	"github.com/yoneyan/vm_mgr/node/db"
	"github.com/yoneyan/vm_mgr/node/etc"
	"github.com/yoneyan/vm_mgr/node/manage"
	pb "github.com/yoneyan/vm_mgr/proto/proto-go"
	"log"
	"strconv"
	"strings"
)

type CreateVMInformation struct {
	ID          int
	Name        string
	CPU         int
	Mem         int
	Storage     string
	StoragePath string
	CDROM       string
	Net         string
	VNC         int
	AutoStart   bool
}

func CreateAutoVMProcess(c *pb.VMData) (string, bool) {
	path := manage.GetMainStorage(c)
	fmt.Println("Path: " + path)
	if etc.FileCopy(etc.GetImagePath()+"/"+c.Image.GetPath(), path) == false {
		return "File Copy Failed...", false
	}
	s := strings.Split(c.Storage, ",")
	size, _ := strconv.Atoi(s[0])
	manage.ResizeStorage(&manage.Storage{Path: manage.GetMainStorage(c), Size: size})

	var r CreateVMInformation

	r.ID = int(c.GetOption().Id)
	r.Name = c.GetVmname()
	r.CPU = int(c.GetVcpu())
	r.Mem = int(c.GetVmem())
	r.Storage = c.GetStorage()
	r.CDROM = c.GetOption().CdromPath
	r.Net = c.GetVnet()
	r.VNC = etc.GenerateVNCPort()
	r.AutoStart = true
	r.StoragePath = manage.StorageProcess(c)

	info, result := CreateVMProcess(&r)
	if result == false {
		fmt.Println("Error: " + info)
	}
	return info, result
}

func CreateVMProcess(c *CreateVMInformation) (string, bool) {
	fmt.Println("----VMNewCreate")

	if manage.VMVncExistsCheck(c.VNC) {
		fmt.Println("A VM with the same vnc port exists. So, change vnc port of the VM.")
		return "same vnc port!!", false
	}
	if manage.VMExistsName(c.Name) {
		fmt.Println("A VM with the same name exists. So, change the name of the VM.")
		return "same vm name!!", false
	}

	if len(c.Net) != 0 {
		d := strings.Split(c.Net, ",")
		fmt.Println(d)
		var net []string
		fmt.Println(net)
		for i, a := range d {
			if i == 0 {
				net = append(net, a)
			} else {
				net = append(net, a)
				net = append(net, manage.GenerateMacAddresss())
			}
		}
		fmt.Println(net)
		c.Net = strings.Join(net, ",")
		fmt.Println(c.Net)
	} else {
		c.Net = ""
	}

	CreateVMDBProcess(c)
	err := RunQEMUCmd(CreateGenerateCmd(c))
	if err != nil {
		fmt.Println(err)
		log.Println("VMNewCreate Error!!")
		return "Error: RunQEMUCmd", false
	} else {
		db.VMDBStatusUpdate(c.ID, 1)
	}

	return "ok", true
}

func CreateVMDBProcess(c *CreateVMInformation) {
	dbdata := db.VM{
		Name:        c.Name,
		CPU:         c.CPU,
		Mem:         c.Mem,
		StoragePath: c.StoragePath,
		Storage:     c.Storage,
		Net:         c.Net,
		Vnc:         c.VNC,
		Socket:      etc.SocketGenerate(c.Name),
		Status:      0,
		AutoStart:   c.AutoStart,
	}
	if db.AddDBVM(dbdata) == false {
		fmt.Println("Error: Failed add vm database")
	}
}
