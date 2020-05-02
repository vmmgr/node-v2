package vm

import (
	"fmt"
	"github.com/yoneyan/vm_mgr/node/db"
	"github.com/yoneyan/vm_mgr/node/manage"
	"strconv"
)

func StartVMProcess(id int) (string, bool) {
	fmt.Println("-----StartVMProcess-----")
	if manage.VMExistsID(id) == false {
		fmt.Println("VM Not Found!!")
		return "VM Not Found!!", false
	}
	data, err := db.VMDBGetData(id)
	if err != nil {
		fmt.Println("VM Data Not Found!!")
		return "VM Data Not Found!!", false
	}
	status, err := db.VMDBGetVMStatus(id)
	if status == 1 {
		fmt.Println("VM is power on!!")
		return "VM is power on!!", false
	} else if status > 1 || status < 0 {
		fmt.Println("VM status is error!! status: " + strconv.Itoa(status))
		return "VM status is error!! status: " + strconv.Itoa(status), false
	}
	var c CreateVMInformation
	c.Name = data.Name
	c.CPU = data.CPU
	c.Mem = data.Mem
	c.Net = data.Net
	c.VNC = data.Vnc
	c.StoragePath = data.StoragePath

	cmd := CreateGenerateCmd(&c)

	err = RunQEMUCmd(cmd)
	if err != nil {
		fmt.Println(err)

		return "QEMURunError!!", false
	}

	fmt.Println("Start End")
	result := db.VMDBStatusUpdate(id, 1)
	if result {
		return "ok", true
	} else {
		return "Status Update Error!!", false
	}

}
