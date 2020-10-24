package vm

import (
	"fmt"
	"github.com/vmmgr/node/db"
	"log"
)

func deleteVMProcess(id int) error {
	if err := vmStop(id); err != nil {
		log.Println("VM state stop")
	}
	log.Println("Stop process end!!")

	if r := db.DeleteDBVM(db.VM{ID: id}); r.Error != nil {
		return fmt.Errorf("Error: delete store error !! ")
	}
	return nil
}
