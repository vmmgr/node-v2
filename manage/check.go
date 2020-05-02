package manage

import (
	"fmt"
	"github.com/yoneyan/vm_mgr/node/db"
	"os"
	"strings"
)

func VMExistsName(name string) bool {
	_, err := db.VMDBGetVMID(name)
	if err != nil {
		return false
	}
	return true
}

func VMExistsID(id int) bool {
	_, err := db.VMDBGetData(id)
	if err != nil {
		return false
	}
	return true
}

func VMExistsCheck(name string, id int) bool {
	if VMExistsID(id) || VMExistsName(name) {
		return true
	} else {
		return false
	}
}

func VMVncExistsCheck(vnc int) bool {
	result := db.VMDBGetAll()
	for a, _ := range result {
		if result[a].Vnc == vnc {
			return true
		}
	}
	return false
}

func FileExistsCheck(filename string) bool {
	if f, err := os.Stat(filename); os.IsNotExist(err) || f.IsDir() {
		return false
	} else {
		return true
	}
}

func InputCheck(inputstoragepath, inputstorage string) bool {
	storagepath := len(strings.Split(inputstoragepath, ","))
	storage := len(strings.Split(inputstorage, ","))

	if storagepath != 0 && storage*2 != storagepath {
		fmt.Println("StoragePath or Storage InputError!!")
		return false
	}
	return true
}

func VMIDCheck(id int) bool {
	if id < 0 {
		fmt.Println("VMID Check NG.")
		return false
	} else {
		fmt.Println("VMID Check OK.")
		return true
	}
}
