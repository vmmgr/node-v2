package vm

import (
	"fmt"
	"github.com/yoneyan/vm_mgr/node/db"
	"github.com/yoneyan/vm_mgr/node/manage"
)

func DeleteVMProcess(id int) (string, bool) {
	result := manage.VMExistsID(id)
	if result == false {
		fmt.Println("VMID Not Found!!")
		return "VMID Not Found!!", false
	}
	info, result := VMStop(id)
	if result == false {
		fmt.Println(info)
		fmt.Println("Already stopped!!")
	} else {
		fmt.Println("Stop process end!!")
	}

	db.DeleteDBVM(id)
	return "ok", true
}
