package vm

import (
	"fmt"
	"github.com/yoneyan/vm_mgr/node/db"
	"time"
)

func StartupProcess() {
	data := db.VMDBGetAll()
	var autostart []int
	for i, _ := range data {
		db.VMDBStatusUpdate(data[i].ID, 0)
		fmt.Println("Status 0  VMID: %d", data[i].ID)
		if data[i].AutoStart {
			autostart = append(autostart, data[i].ID)
		}
	}
	fmt.Printf(" AutoStartVMID: ")
	for i, _ := range autostart {
		time.Sleep(time.Second * 1)
		info, result := StartVMProcess(autostart[i])
		if result {
			fmt.Printf("Start VMID: %d", i)
		} else {
			fmt.Println(info)
			fmt.Printf("Failed start VMID: %d", i)
		}

	}
	fmt.Println("Start process is end!!")
}
