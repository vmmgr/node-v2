package vm

import (
	"fmt"
	"github.com/vmmgr/node/db"
	"log"
	"strconv"
)

type data struct {
	id   int
	boot int
	vmDB db.VM
	iso  string
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

	if err := runQEMUCmd(generateVMCmd(db, d.boot)); err != nil {
		return err
	}

	return nil
}

func generateVMCmd(data db.VM, boot int) []string {
	var b string
	//boot
	if boot == 0 {
		b = "order=c"
	} else if boot == 1 {
		b = "order=d"
	} else if boot == 10 {
		b = "menu=on,strict=on"
	}

	//generate vm start command
	var args []string
	args = []string{
		"qemu-system-x86_64",
		//KVM
		"-enable-kvm",

		"-uuid",
		strconv.Itoa(data.ID) + "-" + data.Name,
		"-name",
		strconv.Itoa(data.ID) + "-" + data.Name,
		"-msg",
		"timestamp=on",

		"-boot",
		b,
		// VNC
		"-vnc",
		fmt.Sprintf("0.0.0.0:%data,websocket=%data", data.ID, data.ID+7000),

		// clock
		"-rtc",
		"base=utc,driftfix=slew",
		"-global",
		"kvm-pit.lost_tick_policy=delay",
		"-no-hpet",

		// CPU
		"-smp",
		fmt.Sprintf("%data,sockets=1,cores=%data,threads=1", data.CPU, data.CPU),
		"-cpu",
		"host",

		// Memory
		"-m",
		strconv.Itoa(data.Mem),
	}

	return args
}
