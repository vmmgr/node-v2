package vm

import (
	"fmt"
	"github.com/vmmgr/node/db"
)

func DeleteVMProcess(id int) error {
	info, result := VMStop(id)
	if result == false {
		fmt.Println(info)
		fmt.Println("Already stopped!!")
	} else {
		fmt.Println("Stop process end!!")
	}

	if r := db.DeleteDBVM(db.VM{ID: id}); r.Error != nil {
		return fmt.Errorf("Error: delete db error !! ")
	}
	return nil
}
