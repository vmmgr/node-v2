package vm

import (
	"fmt"
	"github.com/vmmgr/node/db"
	"github.com/vmmgr/node/etc"
	"log"
	"strconv"
	"strings"
)

type data struct {
	id   int
	boot int
	vmDB db.VM
	iso  string
}

func StartUPVM() error {
	if vm, err := db.GetAllDBVM(); err != nil {
		return err
	} else {
		for _, d := range vm {
			if d.AutoStart == true {
				if err := startVM(data{id: d.ID}); err != nil {
					return err
				}
			}
		}
		return nil
	}
}

func startVM(d data) error {
	log.Println("-----StartVMProcess-----")
	db, err := db.SearchDBVM(db.VM{ID: d.id})
	if err != nil {
		return fmt.Errorf("VM is not found")
	}

	if db.Status == 1 {
		fmt.Println("VM is power on!!")
		return fmt.Errorf("This VM is Power ON ")
	} else if db.Status > 1 || db.Status < 0 {
		return fmt.Errorf("VM status is error!! status: %d ", db.Status)
	}

	if err := runQEMUCmd(generateVMCmd(d)); err != nil {
		return err
	}

	return nil
}

func generateVMCmd(d data) []string {
	var b string
	//boot
	if d.boot == 0 {
		b = "order=c"
	} else if d.boot == 1 {
		b = "order=d"
	} else if d.boot == 10 {
		b = "menu=on,strict=on"
	}

	//generate vm start command
	var args []string
	args = []string{
		"qemu-system-x86_64",
		//KVM
		"-enable-kvm",

		"-uuid",
		strconv.Itoa(d.vmDB.ID) + "-" + d.vmDB.Name,
		"-name",
		strconv.Itoa(d.vmDB.ID) + "-" + d.vmDB.Name,
		"-msg",
		"timestamp=on",

		"-boot",
		b,
		// VNC
		"-vnc",
		fmt.Sprintf("0.0.0.0:%data,websocket=%data", d.vmDB.ID, d.vmDB.ID+7000),

		// clock
		"-rtc",
		"base=utc,driftfix=slew",
		"-global",
		"kvm-pit.lost_tick_policy=delay",
		"-no-hpet",

		// CPU
		"-smp",
		fmt.Sprintf("%data,sockets=1,cores=%data,threads=1", d.vmDB.CPU, d.vmDB.CPU),
		"-cpu",
		"host",

		// Memory
		"-m",
		strconv.Itoa(d.vmDB.Mem),
		// Monitor
		"-monitor",
		etc.SocketGenerate(d.id),
	}

	//Storage
	for _, tmp := range generateStorageCmd(d) {
		args = append(args, tmp)
	}

	//NIC
	for _, tmp := range generateNICCmd(d.vmDB) {
		args = append(args, tmp)
	}

	return args
}

func generateStorageCmd(d data) []string {
	var args []string
	index := 0

	//ISO
	if d.iso != "" {
		iso := strings.Split(d.iso, ",")
		for _, a := range iso {
			args = append(args, "-drive")
			args = append(args, "file="+a+",index="+strconv.Itoa(index)+",media=cdrom")
			index++
		}
	}

	//DISK
	id := strings.Split(d.vmDB.Storage, ",")
	for _, tmp := range id {
		sid, _ := strconv.Atoi(tmp)
		if storage, err := db.SearchDBStorage(db.Storage{ID: sid}); err == nil {
			//virtio
			if storage.Type == 1 {
				args = append(args, "-drive")
				if storage.Type > 10 {
					args = append(args, "file="+storage.Path+",index="+strconv.Itoa(index)+",media=disk,if=virtio")
				} else {
					args = append(args, "file="+etc.GetDiskPath(storage.Type)+"/"+strconv.Itoa(storage.ID)+",index="+strconv.Itoa(index)+",media=disk,if=virtio")
				}
			}
			index++
		}
	}
	return args
}

func generateNICCmd(data db.VM) []string {
	var args []string
	id := strings.Split(data.NIC, ",")
	for _, tmp := range id {
		nid, _ := strconv.Atoi(tmp)
		if nic, err := db.SearchDBNIC(db.NIC{ID: nid}); err == nil {
			net, err := db.SearchDBNet(db.Net{ID: nic.NetID})
			if err != nil {
				break
			}
			if nic.Driver == 1 {
				//virtio
				args = append(args, "-nic")
				args = append(args, "bridge,br=br"+strconv.Itoa(net.VLAN)+",mac="+nic.MacAddress+",model=virtio")
			} else if nic.Driver == 2 {
				//e1000
				args = append(args, "-nic")
				args = append(args, "bridge,br=br"+strconv.Itoa(net.VLAN)+",mac="+nic.MacAddress+",model=e1000")
			}
		}
	}

	return args
}
