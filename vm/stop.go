package vm

import (
	"fmt"
	"github.com/mattn/go-pipeline"
	"github.com/vmmgr/node/db"
	"strconv"

	//"github.com/yoneyan/vm_mgr/node/manage"
	"log"
)

func vmStop(id int) error {
	if d, err := db.SearchDBVM(db.VM{ID: id}); err != nil {
		return fmt.Errorf("VM is not found")
	} else {
		if d.Status == 0 {
			log.Println("Power Off State")
		} else if d.Status == 1 {
			log.Println("Power On State")
		} else {
			log.Println("Error: Power State")
			return fmt.Errorf("Error: Power state ")
		}
	}

	//ps axf | grep test|grep qemu  | grep -v grep | awk '{print "kill -9 " $1}' | sudo sh
	log.Println("-----VMStop Command-----")
	if out, err := pipeline.CombinedOutput(
		[]string{"ps", "axf"},
		[]string{"grep", strconv.Itoa(id) + ".sock"},
		[]string{"grep", "qemu"},
		[]string{"grep", "-v", "grep"},
		[]string{"awk", "{print \"kill -9 \" $1}"},
		[]string{"sudo", "sh"},
	); err != nil {
		log.Println("already stop")
	} else {
		log.Printf("%s \n", out)
	}

	if r := db.UpdateDBVM(db.VM{ID: id, Status: 0}); r.Error != nil {
		return fmt.Errorf("Error: DB Error. but vm stopped ")
	} else {
		return nil
	}
}
