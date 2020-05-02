package vm

import (
	"fmt"
	"github.com/mattn/go-pipeline"
	"github.com/yoneyan/vm_mgr/node/db"
	"github.com/yoneyan/vm_mgr/node/manage"
)

func VMStop(id int) (string, bool) {
	if manage.VMExistsID(id) == false {
		fmt.Println("VMID Error")
		return "Not Found VMID", false
	}

	result, err := db.VMDBGetData(id)
	if result.Status == 1 {
		fmt.Println("Power On State")
	} else if result.Status == 0 {
		fmt.Println("Power Off State")
	} else {
		fmt.Println("Power State Error")
		return "Power State Error!!", false
	}

	fmt.Println(result)
	if err != nil {
		fmt.Println("Error!!")
	}

	fmt.Println(result.Name)

	//ps axf | grep test|grep qemu  | grep -v grep | awk '{print "kill -9 " $1}' | sudo sh
	fmt.Println("-----VMStop Command-----")
	out, err := pipeline.CombinedOutput(
		[]string{"ps", "axf"},
		[]string{"grep", result.Name + ".sock"},
		[]string{"grep", "qemu"},
		[]string{"grep", "-v", "grep"},
		[]string{"awk", "{print \"kill -9 \" $1}"},
		[]string{"sudo", "sh"},
	)
	if err != nil {
		fmt.Println("already stopped")
	}
	fmt.Printf("%s", out)

	if db.VMDBStatusUpdate(id, 0) {
		fmt.Println("Power Off state")
	} else {
		fmt.Println("state Error!!")
	}
	return "ok", true
}

func StopProcess() {
	data := db.VMDBGetAll()
	var status []int
	for i, _ := range data {
		fmt.Printf("Status 0  VMID: %d", data[i].ID)
		if data[i].Status == 1 {
			status = append(status, data[i].ID)
		}
	}
	fmt.Printf("AutoStartVMID: ")
	fmt.Println(status)

	for i, _ := range status {
		info, result := VMStop(status[i])
		{
			if result == false {
				fmt.Println(info)
				fmt.Printf("Failed stop VMID: %d", i)
			}
			fmt.Printf("Start VMID: %d", i)
		}
		fmt.Println()

		fmt.Println("Start process is end!!")
	}
}
